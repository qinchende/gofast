package timex

import "time"

// 相对时间（本系统时间原点）
// 这样在很多地方就只需要存储一个 Duration 类型的值，占用8字节，避免了存储 time.Time 类型（占用24字节）。
// 当前时间年月日分别减去1之后的时间作为参考时间点，其它时间都和这个比得到相对时间差值。
// Note：固定 UTC 时间 2024-01-01 00:00:00 为基准时间，Duration都是相对这个时间的偏移
var sdxBaseTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // 自定义基准时间，这样偏移量数值会小一些

// 当前相对原点的时差，因为全系统都是相对时间，你可以认为这个时差就是当前时间
func SdxNowDur() time.Duration {
	return time.Now().Sub(sdxBaseTime)
}

func SdxToDur(tm *time.Time) time.Duration {
	return tm.Sub(sdxBaseTime)
}

// 将指定的相对时间转成 真实的 time.Time 类型
func SdxToTime(d time.Duration) time.Time {
	return sdxBaseTime.Add(d)
}

// 当前时间和传入的指定时间之间的时间差（这里比较绕，这里的两个时间都是相对本框架的时间标准）
func SdxNowDiffDur(d time.Duration) time.Duration {
	return SdxNowDur() - d
}

// 和当前时间差多少秒
func SdxNowDiffS(d time.Duration) int64 {
	return int64((SdxNowDur() - d) / time.Second)
}

// 和当前时间差多少毫秒
func SdxNowDiffMS(d time.Duration) int64 {
	return int64((SdxNowDur() - d) / time.Millisecond)
}

func SdxNowAddDur(d time.Duration) time.Duration {
	return SdxNowDur() + d
}

func SdxNowAddSDur(s int) time.Duration {
	return SdxNowDur() + time.Duration(s)*time.Second
}

func SdxNowAddMSDur(ms int) time.Duration {
	return SdxNowDur() + time.Duration(ms)*time.Millisecond
}
