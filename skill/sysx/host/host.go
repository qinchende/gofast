package host

import (
	"github.com/qinchende/gofast/skill/lang"
	"os"
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = lang.RandId()
	}
}

func Hostname() string {
	return hostname
}
