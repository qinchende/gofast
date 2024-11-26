package httpx

import (
	"fmt"
	"github.com/qinchende/gofast/aid/logx"
	"net/http"
	"net/url"
	"os"
)

func GetRemoteAddr(r *http.Request) string {
	v := r.Header.Get(XForwardFor)
	if len(v) > 0 {
		return v
	}
	return r.RemoteAddr
}

func ResolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			logx.Debug().SendMsg("Environment variable PORT=" + port)
			return ":" + port
		}
		logx.Debug().SendMsg("Environment variable PORT is undefined. Using port :8099 by default")
		return ":8099"
	case 1:
		return addr[0]
	default:
		panic("Too many parameters")
	}
}

func CheckWriteHeaderCode(code int) {
	// Issue 22880: require valid WriteHeader status codes.
	// For now we only enforce that it's three digits.
	// In the future we might block things over 599 (600 and above aren't defined
	// at https://httpwg.org/specs/rfc7231.html#status.codes)
	// and we might block under 200 (once we have more mature 1xx support).
	// But for now any three digits.
	//
	// We used to send "HTTP/1.1 000 0" on the wire in responses but there's
	// no equivalent bogus thing we can realistically send in HTTP/2,
	// so we'll consistently panic instead and help people find their bugs
	// early. (We can't return an error from WriteHeader even if we wanted to.)
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}

// 拷贝集合对象
func CopyUrlValues(dst, src url.Values) {
	for k, vs := range src {
		dst[k] = append(dst[k], vs...)
	}
}
