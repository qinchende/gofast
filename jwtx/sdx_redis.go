package jwtx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/bytesconv"
	"github.com/qinchende/gofast/skill/json"
)

// 从 redis 中获取 当前 请求上下文的 session data.
func (ss *SdxSession) initCtxSess(ctx *fst.Context) {
	str, err := ss.Redis.Get(sdxSessKeyPrefix + ctx.Sess.Sid)
	if str == "" || err != nil {
		str = `{}`
	}

	err = json.Unmarshal(bytesconv.StringToBytes(str), &ctx.Sess.Values)
	if err != nil {
		fst.RaisePanicErr(err)
	}
}

//func InitRedis(ss *fst.CtxSession) {
//
//}

func SaveRedis(sdx *fst.CtxSession) {
	str, _ := json.Marshal(sdx.Values)
	_, _ = ss.Redis.Set(sdxSessKeyPrefix+sdx.Sid, str, ss.TTL)
}
