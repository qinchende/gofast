package cpux

import (
	"errors"
	"github.com/qinchende/gofast/aid/logx"
	"strings"
	"time"

	"github.com/qinchende/gofast/aid/iox"
)

const (
	cpuTicks  = 100
	cpuFields = 8
)

var (
	preSystem uint64
	preTotal  uint64
	quota     float64
	cores     uint64
)

// if /proc not present, ignore the cpu calcuation, like wsl linux
func init() {
	cpus, err := perCpuUsage()
	if err != nil {
		logx.Error(err)
		return
	}

	cores = uint64(len(cpus))
	sets, err := cpuSets()
	if err != nil {
		logx.Error(err)
		return
	}

	quota = float64(len(sets))
	cq, err := cpuQuota()
	if err == nil {
		if cq != -1 {
			period, err := cpuPeriod()
			if err != nil {
				logx.Error(err)
				return
			}

			limit := float64(cq) / float64(period)
			if limit < quota {
				quota = limit
			}
		}
	}

	preSystem, err = systemCpuUsage()
	if err != nil {
		logx.Error(err)
		return
	}

	preTotal, err = totalCpuUsage()
	if err != nil {
		logx.Error(err)
		return
	}
}

func RefreshCpu() uint64 {
	total, err := totalCpuUsage()
	if err != nil {
		return 0
	}
	system, err := systemCpuUsage()
	if err != nil {
		return 0
	}

	var usage uint64
	cpuDelta := total - preTotal
	systemDelta := system - preSystem
	if cpuDelta > 0 && systemDelta > 0 {
		usage = uint64(float64(cpuDelta*cores*1e3) / (float64(systemDelta) * quota))
	}
	preSystem = system
	preTotal = total

	return usage
}

func cpuQuota() (int64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return 0, err
	}

	return cg.cpuQuotaUs()
}

func cpuPeriod() (uint64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return 0, err
	}

	return cg.cpuPeriodUs()
}

func cpuSets() ([]uint64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return nil, err
	}

	return cg.cpus()
}

func perCpuUsage() ([]uint64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return nil, err
	}

	return cg.acctUsagePerCpu()
}

func systemCpuUsage() (uint64, error) {
	lines, err := iox.ReadTextLines("/proc/stat", iox.WithoutBlank())
	if err != nil {
		return 0, err
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			if len(fields) < cpuFields {
				return 0, errors.New("bad format of cpu stats")
			}

			var totalClockTicks uint64
			for _, i := range fields[1:cpuFields] {
				v, err := parseUint(i)
				if err != nil {
					return 0, err
				}

				totalClockTicks += v
			}

			return (totalClockTicks * uint64(time.Second)) / cpuTicks, nil
		}
	}

	return 0, errors.New("bad stats format")
}

func totalCpuUsage() (usage uint64, err error) {
	var cg *cgroup
	if cg, err = currentCgroup(); err != nil {
		return
	}

	return cg.acctUsageAllCpus()
}
