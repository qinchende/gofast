package httpx

import (
	"github.com/qinchende/gofast/logx"
	"net/http"
	"os"
)

const xForwardFor = "X-Forward-For"

// Returns the peer address, supports X-Forward-For
func GetRemoteAddr(r *http.Request) string {
	v := r.Header.Get(xForwardFor)
	if len(v) > 0 {
		return v
	}
	return r.RemoteAddr
}

func ResolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			logx.DebugPrint(`Environment variable PORT="%s"`, port)
			return ":" + port
		}
		logx.DebugPrint("Environment variable PORT is undefined. Using port :8099 by default")
		return ":8099"
	case 1:
		return addr[0]
	default:
		panic("Too many parameters")
	}
}
