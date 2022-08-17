// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"io"
	"log"
	"strings"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type WriterCloser interface {
	io.Writer
	io.Closer
	Writeln(data string) (err error)
	WritelnBytes(data []byte) (err error)
	WritelnBuilder(sb *strings.Builder) (err error)
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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

func (lw logWriter) WritelnBytes(bs []byte) error {
	err := lw.logger.Output(2, string(bs))
	return err
}

func (lw logWriter) WritelnBuilder(sb *strings.Builder) error {
	err := lw.logger.Output(2, sb.String())
	return err
}
