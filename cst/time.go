package cst

import "time"

const (
	TimeFmtSaveRFC3339 = time.RFC3339 // "2006-01-02T15:04:05Z07:00" or "2006-01-02T15:04:05+08:00"
	TimeFmtSaveYmdHms  = "2006-01-02 15:04:05"
	TimeFmtSaveMdHms   = "01-02 15:04:05"
	TimeFmtSaveYmd     = "2006-01-02"
	TimeFmtSaveHms     = "15:04:05"
)
