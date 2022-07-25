package sdx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/mid"
)

func JwtAuthHandler(secret string) fst.CtxHandler {
	return mid.JwtAuthHandler(secret)
}

func PmsParser(ctx *fst.Context) {
	if err := ctx.ParseRequestData(); err != nil {
		ctx.AbortFaiMsg("PmsParser err: " + err.Error())
	}
}

//func BuildPmsOfJson(ctx *fst.Context) {
//	ctx.GenPmsByJSONBody()
//}
//
//func BuildPmsOfXml(ctx *fst.Context) {
//	ctx.GenPmsByXMLBody()
//}
//
//func BuildPmsOfForm(ctx *fst.Context) {
//	ctx.GenPmsByFormBody()
//}
