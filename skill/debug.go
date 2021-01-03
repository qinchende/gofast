package skill

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var isDebug bool = false
var _once sync.Once
// 每个程序只设置一次debug标志位，后面的设置都失效
func SetDebugStatus(yn bool) {
	_once.Do(func() {
		isDebug = yn
	})
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
