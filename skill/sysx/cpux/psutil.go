package cpux

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"math"
)

func GetAllBusy(t cpu.TimesStat) (float64, float64) {
	busy := t.User + t.System + t.Nice + t.Iowait + t.Irq +
		t.Softirq + t.Steal
	return busy + t.Idle, busy
}

func BusyPercent(t1, t2 cpu.TimesStat) float64 {
	t1All, t1Busy := GetAllBusy(t1)
	t2All, t2Busy := GetAllBusy(t2)

	if t2Busy <= t1Busy {
		return 0
	}
	if t2All <= t1All {
		return 100
	}
	return math.Min(100, math.Max(0, (t2Busy-t1Busy)/(t2All-t1All)*100))
}
