package jsonx

import (
	"errors"
	"math"
)

const (
	maxJsonLength = math.MaxInt32 - 1 // 最大2GB
)

var (
	sErr            = errors.New("jsonx: json syntax error.")
	errJsonTooLarge = errors.New("jsonx: string too large.")
)
