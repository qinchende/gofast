package timex

import (
	"fmt"
	"time"
)

func ToStringMS(dur time.Duration) string {
	return fmt.Sprintf("%.1fms", float32(dur)/float32(time.Millisecond))
}
