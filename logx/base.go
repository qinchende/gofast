package logx

import (
	"io"
	"os"
)

var (
	DefaultWriter      io.Writer = os.Stdout
	DefaultErrorWriter io.Writer = os.Stderr
)
