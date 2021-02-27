package jwtx

import (
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/fst"
)

type SdxSessConfig struct {
	sessKey string
	secret  string
}

type SdxSession struct {
	SdxSessConfig
	Redis *redis.GoRedisX
}

var ss *SdxSession

func InitSdxRedis(i *SdxSession) {
	ss = i
}

func SdxSessHandler(ctx *fst.Context) {
	//ss.Redis
	if ctx.Pms["tok"] == nil {

	}
}
