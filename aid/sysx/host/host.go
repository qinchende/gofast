package host

import (
	"github.com/qinchende/gofast/aid/randx"
	"os"
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = randx.RandId()
	}
}

func Hostname() string {
	return hostname
}
