package skill

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var isDebug bool = false

func SetDebugStatus(yn bool) {
	isDebug = yn
}

func DebugPrint(format string, values ...interface{}) {
	if isDebug {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(os.Stdout, "[fst-debug] "+format, values...)
	}
}

func DebugPrintError(err error) {
	if err != nil && isDebug {
		fmt.Fprintf(os.Stderr, "[fst-debug] [ERROR] %v\n", err)
	}
}

func GetMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}
