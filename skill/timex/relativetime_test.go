package timex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRelativeTime(t *testing.T) {
	time.Sleep(time.Millisecond)
	now := NowDur()
	assert.True(t, now > 0)
	time.Sleep(time.Millisecond)
	assert.True(t, NowDiffDur(now) > 0)
}

func TestRelativeTime_Time(t *testing.T) {
	diff := time.Until(time.Now())
	if diff > 0 {
		assert.True(t, diff < time.Second)
	} else {
		assert.True(t, -diff < time.Second)
	}
}
