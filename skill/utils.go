package skill

import (
	"os"
	"reflect"
	"runtime"
)

func NameOfFunc(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func ResolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			DebugPrint(`Environment variable PORT="%s"`, port)
			return ":" + port
		}
		DebugPrint("Environment variable PORT is undefined. Using port :8099 by default")
		return ":8099"
	case 1:
		return addr[0]
	default:
		panic("Too many parameters")
	}
}
