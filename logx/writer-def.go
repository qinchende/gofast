package logx

import (
	"io"
	"log"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type WriterCloser interface {
	io.Writer
	io.Closer
	Writeln(data string) (err error)
}

// 自定义 logger
type logWriter struct {
	logger *log.Logger
}

func newLogWriter(logger *log.Logger) logWriter {
	return logWriter{
		logger: logger,
	}
}

func (lw logWriter) Close() error {
	return nil
}

func (lw logWriter) Write(data []byte) (int, error) {
	lw.logger.Print(data)
	return len(data), nil
}

func (lw logWriter) Writeln(data string) error {
	err := lw.logger.Output(2, data)
	return err
}
