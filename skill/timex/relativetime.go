package timex

import (
	"time"
)

// important: 这里的 time.Duration 全部是相对基准时间的偏移

// 相对时间（本系统时间原点）
// 这样在很多地方就只需要存储一个 Duration 类型的值，占用8字节，避免了存储 time.Time 类型（占用24字节）。
// 当前时间年月日分别减去1之后的时间作为参考时间点，其它时间都和这个比得到相对时间差值。
// Use the long enough past time as start time, in case timex.NowDur() - lastTime equals 0.
// Note：固定 1970-01-01 00:00:00 为基准时间，Duration都是相对这个时间的偏移
var initTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)

// 当前相对原点的时差，因为全系统都是相对时间，你可以认为这个时差就是当前时间
func NowDur() time.Duration {
	return time.Now().Sub(initTime)
}

func ToDur(tm *time.Time) time.Duration {
	return tm.Sub(initTime)
}

// 将指定的相对时间转成 真实的 time.Time 类型
func ToTime(d time.Duration) time.Time {
	return initTime.Add(d)
}

func ToS(d time.Duration) int64 {
	return int64(d / time.Second)
}

func ToMS(d time.Duration) int64 {
	return int64(d / time.Millisecond)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 两个时间相差多少秒
func DiffS(a, b time.Duration) int64 {
	return int64((a - b) / time.Second)
}

// 两个时间相差多少毫秒
func DiffMS(a, b time.Duration) int64 {
	return int64((a - b) / time.Millisecond)
}

// 当前时间和传入的指定时间之间的时间差（这里比较绕，这里的两个时间都是相对本框架的时间标准）
func NowDiffDur(d time.Duration) time.Duration {
	return NowDur() - d
}

// 和当前时间差多少秒
func NowDiffS(d time.Duration) int64 {
	return int64((NowDur() - d) / time.Second)
}

// 和当前时间差多少毫秒
func NowDiffMS(d time.Duration) int64 {
	return int64((NowDur() - d) / time.Millisecond)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 和当前时间差多少秒
func NowAddSDur(s int) time.Duration {
	return NowDur() + time.Duration(s)*time.Second
}

// 和当前时间差多少毫秒
func NowAddMSDur(ms int) time.Duration {
	return NowDur() + time.Duration(ms)*time.Millisecond
}
