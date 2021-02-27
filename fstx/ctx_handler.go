package fstx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/mid"
)

func JwtAuthHandler(secret string) fst.CtxHandler {
	return mid.JwtAuthHandler(secret)
}
