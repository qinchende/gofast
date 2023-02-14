package task

type TaskFunc func() bool

const (
	LiteRedisKeyPrefix  = "GF.LT."
	DefGroupIntervalS   = 60
	DefTaskPetStartTime = "00:00"
	DefTaskPetEndTime   = "23:59"
)
