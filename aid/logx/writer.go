// Copyright 2022 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
package logx

import (
	"io"
	"sync"
)

const (
	dateFormatYMD       = "2006-01-02"
	hoursPerDay         = 24
	bufferSize          = 100
	defaultDirMode      = 0755
	defaultFileMode     = 0600
	backupFileDelimiter = "-"
)

//
//type WriterCloser interface {
//	io.Writer
//	//Writeln(data string) (err error)
//	//WritelnBytes(data []byte) (err error)
//	//WritelnBuilder(sb *strings.Builder) (err error)
//}

type syncWriter struct {
	mu sync.Mutex
	lw io.Writer
}

type consoleWriter struct {
	*limitedExecutor
	lw io.Writer
}

type multiWriter struct {
	lws []io.Writer
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//type logWriter struct {
//	logger *log.Logger
//}
//
//func newLogWriter(logger *log.Logger) logWriter {
//	return logWriter{
//		logger: logger,
//	}
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (lw logWriter) Close() error {
//	return nil
//}
//
//func (lw logWriter) Write(data []byte) (int, error) {
//	err := lw.logger.LogBuilder(2, lang.B2S(data))
//	return len(data), err
//}

//func (lw logWriter) Writeln(data string) error {
//	err := lw.logger.LogBuilder(2, data)
//	return err
//}
//
//func (lw logWriter) WritelnBytes(bs []byte) error {
//	err := lw.logger.LogBuilder(2, string(bs))
//	return err
//}
//
//func (lw logWriter) WritelnBuilder(sb *strings.Builder) error {
//	err := lw.logger.LogBuilder(2, sb.String())
//	return err
//}
