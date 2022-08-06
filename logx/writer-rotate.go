package logx

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/qinchende/gofast/skill/fs"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/timex"
)

const (
	dateFormat          = "2006-01-02"
	hoursPerDay         = 24
	bufferSize          = 100
	defaultDirMode      = 0755
	defaultFileMode     = 0600
	backupFileDelimiter = "-"
)

type (
	RotateRule interface {
		BackupFileName() string
		OutdatedFiles() []string
		MarkRotated()
		NeedRotate() bool
	}

	RotateLogger struct {
		filename string
		backup   string
		fp       *os.File
		channel  chan []byte
		done     chan lang.PlaceholderType
		rule     RotateRule
		compress bool
		keepDays int
		// can't use threading.RoutineGroup because of cycle import
		waitGroup sync.WaitGroup
		closeOnce sync.Once
		//mu        sync.Mutex // ensures atomic writes; protects the following fields
		//buf       []byte     // for accumulating text to write
	}

	DailyRotateRule struct {
		yearDay   int
		filename  string
		delimiter string
		days      int
		gzip      bool
	}
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func DefDailyRotateRule(filename, delimiter string, days int, gzip bool) RotateRule {
	return &DailyRotateRule{
		yearDay:   time.Now().YearDay(),
		filename:  filename,
		delimiter: delimiter,
		days:      days,
		gzip:      gzip,
	}
}

func (r *DailyRotateRule) BackupFileName() string {
	return fmt.Sprintf("%s%s%s", r.filename, r.delimiter, time.Now().Format(dateFormat))
}

func (r *DailyRotateRule) MarkRotated() {
	r.yearDay = time.Now().YearDay()
}

func (r *DailyRotateRule) OutdatedFiles() []string {
	if r.days <= 0 {
		return nil
	}

	var pattern string
	if r.gzip {
		pattern = fmt.Sprintf("%s%s*.gz", r.filename, r.delimiter)
	} else {
		pattern = fmt.Sprintf("%s%s*", r.filename, r.delimiter)
	}

	files, err := filepath.Glob(pattern)
	if err != nil {
		ErrorF("failed to delete outdated log files, error: %s", err)
		return nil
	}

	var buf strings.Builder
	boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay*r.days)).Format(dateFormat)
	fmt.Fprintf(&buf, "%s%s%s", r.filename, r.delimiter, boundary)
	if r.gzip {
		buf.WriteString(".gz")
	}
	boundaryFile := buf.String()

	var outDates []string
	for _, file := range files {
		if file < boundaryFile {
			outDates = append(outDates, file)
		}
	}

	return outDates
}

func (r *DailyRotateRule) NeedRotate() bool {
	return time.Now().YearDay() != r.yearDay
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 自动归档日志系统
func NewRotateLogger(filename string, rule RotateRule, compress bool) (*RotateLogger, error) {
	rl := &RotateLogger{
		filename: filename,
		channel:  make(chan []byte, bufferSize),
		done:     make(chan lang.PlaceholderType),
		rule:     rule,
		compress: compress,
	}
	if err := rl.initRotateLogger(); err != nil {
		return nil, err
	}

	rl.startWorker()
	return rl, nil
}

func (rl *RotateLogger) Close() error {
	var err error
	rl.closeOnce.Do(func() {
		close(rl.done)
		rl.waitGroup.Wait()
		if err = rl.fp.Sync(); err != nil {
			return
		}
		err = rl.fp.Close()
	})
	return err
}

func (rl *RotateLogger) Write(data []byte) (int, error) {
	select {
	case rl.channel <- data:
		return len(data), nil
	case <-rl.done:
		log.Println(data)
		return 0, errors.New("error: log file closed")
	}
}

// 每次调用都会在 data 后面自动判断并加上 \n
func (rl *RotateLogger) Writeln(data string) (err error) {
	// TODO：字符串转换成 字节切片，增加内存开销，会影响性能，类似这种问题一会想办法改进
	bs := []byte(data)

	if len(bs) == 0 || bs[len(bs)-1] != '\n' {
		bs = append(bs, '\n')
	}
	_, err = rl.Write(bs)
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// utils
func (rl *RotateLogger) getBackupFilename() string {
	if len(rl.backup) == 0 {
		return rl.rule.BackupFileName()
	} else {
		return rl.backup
	}
}

// TODO: 这里可以指定用 标准库中的 log 包来写数据
func (rl *RotateLogger) initRotateLogger() error {
	rl.backup = rl.rule.BackupFileName()

	if _, err := os.Stat(rl.filename); err != nil {
		basePath := path.Dir(rl.filename)
		if _, err = os.Stat(basePath); err != nil {
			if err = os.MkdirAll(basePath, defaultDirMode); err != nil {
				return err
			}
		}

		if rl.fp, err = os.Create(rl.filename); err != nil {
			return err
		}
	} else if rl.fp, err = os.OpenFile(rl.filename, os.O_APPEND|os.O_WRONLY, defaultFileMode); err != nil {
		return err
	}

	fs.CloseOnExec(rl.fp)

	return nil
}

func (rl *RotateLogger) maybeCompressFile(file string) {
	if !rl.compress {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			Stack(r)
		}
	}()
	compressLogFile(file)
}

func (rl *RotateLogger) maybeDeleteOutdatedFiles() {
	files := rl.rule.OutdatedFiles()
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			ErrorF("failed to remove outdated file: %s", file)
		}
	}
}

func (rl *RotateLogger) postRotate(file string) {
	go func() {
		// we cannot use threading.GoSafe here, because of import cycle.
		rl.maybeCompressFile(file)
		rl.maybeDeleteOutdatedFiles()
	}()
}

func (rl *RotateLogger) rotate() error {
	if rl.fp != nil {
		err := rl.fp.Close()
		rl.fp = nil
		if err != nil {
			return err
		}
	}

	_, err := os.Stat(rl.filename)
	if err == nil && len(rl.backup) > 0 {
		backupFilename := rl.getBackupFilename()
		err = os.Rename(rl.filename, backupFilename)
		if err != nil {
			return err
		}

		rl.postRotate(backupFilename)
	}

	rl.backup = rl.rule.BackupFileName()
	if rl.fp, err = os.Create(rl.filename); err == nil {
		fs.CloseOnExec(rl.fp)
	}

	return err
}

func (rl *RotateLogger) startWorker() {
	rl.waitGroup.Add(1)

	go func() {
		defer rl.waitGroup.Done()
		for {
			select {
			case bytes := <-rl.channel:
				rl.writeExec(bytes)
			case <-rl.done:
				return
			}
		}
	}()
}

// 检查标记，做好日志的拆分，自动判断 gzip 标记并压缩
func (rl *RotateLogger) writeExec(data []byte) {
	if rl.rule.NeedRotate() {
		if err := rl.rotate(); err != nil {
			log.Println(err)
		} else {
			rl.rule.MarkRotated()
		}
	}
	if rl.fp != nil {
		_, _ = rl.fp.Write(data)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func compressLogFile(file string) {
	start := timex.Now()
	InfoF("compressing log file: %s", file)
	if err := gzipFile(file); err != nil {
		ErrorF("compress error: %s", err)
	} else {
		InfoF("compressed log file: %s, took %s", file, timex.Since(start))
	}
}

func gzipFile(file string) error {
	in, err := os.Open(file)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(fmt.Sprintf("%s.gz", file))
	if err != nil {
		return err
	}
	defer out.Close()

	w := gzip.NewWriter(out)
	if _, err = io.Copy(w, in); err != nil {
		return err
	} else if err = w.Close(); err != nil {
		return err
	}

	return os.Remove(file)
}
