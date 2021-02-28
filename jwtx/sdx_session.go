package jwtx

import (
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/fst"
)

type SdxSessConfig struct {
	sessKey string
	secret  string

	sessTTL    int
	sessTTLNew int
}

type SdxSession struct {
	SdxSessConfig
	Redis *redis.GoRedisX
}

var ss *SdxSession

func InitSdxRedis(i *SdxSession) {
	ss = i
	if ss.sessTTL == 0 {
		ss.sessTTL = 3600 * 4 // 默认4个小时
	}
	if ss.sessTTLNew == 0 {
		ss.sessTTLNew = 180 // 默认三分钟
	}
}

func SdxSessHandler(ctx *fst.Context) {
	//ss.Redis
	if ctx.Pms["tok"] == nil {

	}
}

