package sdx

import (
	"github.com/qinchende/gofast/cst"
)

type JwtSession struct {
	values    cst.KV // map[string]interface{}
	sessIsNew bool   // Sid is new
}

// JwtSession 需要实现 sessionKeeper 所有接口
//var _ fst.SessionKeeper = &JwtSession{}
