package timex

import "time"

// 相对时间，这样在很多地方就只需要存储一个 Duration 类型的值，占用8字节，避免了存储 time.Time 类型。

// Use the long enough past time as start time, in case timex.Now() - lastTime equals 0.
var initTime = time.Now().AddDate(-1, -1, -1)

func Now() time.Duration {
	return time.Since(initTime)
}

func Since(d time.Duration) time.Duration {
	return time.Since(initTime) - d
}

func Time() time.Time {
	return initTime.Add(Now())
}

func ToTime(d time.Duration) time.Time {
	return initTime.Add(d)
}
