package timex

import (
	"time"
)

// 相对时间（本系统时间原点）
// 这样在很多地方就只需要存储一个 Duration 类型的值，占用8字节，避免了存储 time.Time 类型（占用24字节）。
// 当前时间年月日分别减去1之后的时间作为参考时间点，其它时间都和这个比得到相对时间差值。
// Use the long enough past time as start time, in case timex.Now() - lastTime equals 0.
//var initTime = time.Now().AddDate(-1, -1, -1)
var initTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)

// 当前相对原点的时差，因为全系统都是相对时间，你可以认为这个时差就是当前时间
func Now() time.Duration {
	return time.Now().Sub(initTime)
}

func ToDuration(tm *time.Time) time.Duration {
	return tm.Sub(initTime)
}

// 当前时间和传入的指定时间之间的时间差（这里比较绕，都是相对的概念）
func Since(d time.Duration) time.Duration {
	return Now() - d
}

// 时间差毫秒
func SinceMS(d time.Duration) int64 {
	return int64(Since(d) / time.Millisecond)
}

func DiffMS(a, b time.Duration) int64 {
	return int64((a - b) / time.Millisecond)
}

// 时间差秒
func SinceS(d time.Duration) int64 {
	return int64(Since(d) / time.Second)
}

func DiffS(a, b time.Duration) int64 {
	return int64((a - b) / time.Second)
}

// 将指定的相对时间转成 真实的 time.Time 类型
func ToTime(d time.Duration) time.Time {
	return initTime.Add(d)
}
