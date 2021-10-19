package stat

import (
	"os"
	"runtime"
	"time"

	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/executors"
)

var (
	LogInterval = time.Minute
	//writerLock  sync.Mutex
	//reportWriter Writer = nil
)

type ReqItem struct {
	Drop     bool          // 是否是一个丢弃的任务
	Duration time.Duration // 任务耗时
	//Description string
}

type (
	//Writer interface {
	//	Write(report *MetricInfo) error
	//}

	MetricInfo struct {
		Name          string  `json:"name"`
		Timestamp     int64   `json:"tm"`
		Pid           int     `json:"pid"`
		ReqsPerSecond float32 `json:"qps"`
		Drops         int     `json:"drops"`
		Average       float32 `json:"avg"`
		Median        float32 `json:"med"`
		Top90th       float32 `json:"t90"`
		Top99th       float32 `json:"t99"`
		Top99p9th     float32 `json:"t99p9"`
	}

	Metrics struct {
		executor  *executors.IntervalExecutor
		container *metricsContainer
	}
)

//func SetReportWriter(writer Writer) {
//	writerLock.Lock()
//	reportWriter = writer
//	writerLock.Unlock()
//}

func NewMetrics(name string) *Metrics {
	container := &metricsContainer{
		name: name,
		pid:  os.Getpid(),
	}

	return &Metrics{
		executor:  executors.NewIntervalExecutor(LogInterval, container),
		container: container,
	}
}

func (m *Metrics) AddItem(item ReqItem) {
	m.executor.Add(item)
}

func (m *Metrics) AddDrop() {
	m.executor.Add(ReqItem{
		Drop: true,
	})
}

func (m *Metrics) SetName(name string) {
	m.executor.Sync(func() {
		m.container.name = name
	})
}

type (
	tasksDurationPair struct {
		items    []ReqItem
		duration time.Duration
		drops    int
	}

	metricsContainer struct {
		name     string
		pid      int
		items    []ReqItem
		duration time.Duration
		drops    int
	}
)

// 添加新项
func (c *metricsContainer) AddItem(v interface{}) bool {
	if item, ok := v.(ReqItem); ok {
		if item.Drop {
			c.drops++
		} else {
			c.items = append(c.items, item)
			c.duration += item.Duration
		}
	}
	return false
}

// 执行任务
func (c *metricsContainer) Execute(v interface{}) {
	pair := v.(tasksDurationPair)
	items := pair.items
	duration := pair.duration
	drops := pair.drops
	size := len(items)
	report := &MetricInfo{
		Name:          c.name,
		Timestamp:     time.Now().Unix(),
		Pid:           c.pid,
		ReqsPerSecond: float32(size) / float32(LogInterval/time.Second),
		Drops:         drops,
	}

	if size > 0 {
		report.Average = float32(duration/time.Millisecond) / float32(size)

		fiftyPercent := size >> 1
		if fiftyPercent > 0 {
			top50pTasks := topK(items, fiftyPercent)
			medianTask := top50pTasks[0]
			report.Median = float32(medianTask.Duration) / float32(time.Millisecond)
			tenPercent := fiftyPercent / 5
			if tenPercent > 0 {
				top10pTasks := topK(items, tenPercent)
				task90th := top10pTasks[0]
				report.Top90th = float32(task90th.Duration) / float32(time.Millisecond)
				onePercent := tenPercent / 10
				if onePercent > 0 {
					top1pTasks := topK(top10pTasks, onePercent)
					task99th := top1pTasks[0]
					report.Top99th = float32(task99th.Duration) / float32(time.Millisecond)
					pointOnePercent := onePercent / 10
					if pointOnePercent > 0 {
						topPointOneTasks := topK(top1pTasks, pointOnePercent)
						task99Point9th := topPointOneTasks[0]
						report.Top99p9th = float32(task99Point9th.Duration) / float32(time.Millisecond)
					} else {
						report.Top99p9th = getTopDuration(top1pTasks)
					}
				} else {
					mostDuration := getTopDuration(top10pTasks)
					report.Top99th = mostDuration
					report.Top99p9th = mostDuration
				}
			} else {
				mostDuration := getTopDuration(items)
				report.Top90th = mostDuration
				report.Top99th = mostDuration
				report.Top99p9th = mostDuration
			}
		} else {
			mostDuration := getTopDuration(items)
			report.Median = mostDuration
			report.Top90th = mostDuration
			report.Top99th = mostDuration
			report.Top99p9th = mostDuration
		}
	}

	log(report)
}

func (c *metricsContainer) RemoveAll() interface{} {
	items := c.items
	duration := c.duration
	drops := c.drops
	c.items = nil
	c.duration = 0
	c.drops = 0

	return tasksDurationPair{
		items:    items,
		duration: duration,
		drops:    drops,
	}
}

func getTopDuration(items []ReqItem) float32 {
	top := topK(items, 1)
	if len(top) < 1 {
		return 0
	} else {
		return float32(top[0].Duration) / float32(time.Millisecond)
	}
}

func log(report *MetricInfo) {
	// writeReport(report)
	logx.Statf("(%s) - qps: %.1f/s, drops: %d, avg time: %.1fms, med: %.1fms, 90th: %.1fms, 99th: %.1fms, 99.9th: %.1fms, G: %d",
		report.Name, report.ReqsPerSecond, report.Drops, report.Average, report.Median,
		report.Top90th, report.Top99th, report.Top99p9th, runtime.NumGoroutine())
}

//func writeReport(report *MetricInfo) {
//	writerLock.Lock()
//	defer writerLock.Unlock()
//
//	if reportWriter != nil {
//		if err := reportWriter.Write(report); err != nil {
//			logx.Error(err)
//		}
//	}
//}
