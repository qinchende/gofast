package dts

import "errors"

var (
	errNumOutOfRange = errors.New("dts: number out of range")
)

//recordPtr := reflect.New(recordType)
//recordVal := reflect.Indirect(recordPtr)
