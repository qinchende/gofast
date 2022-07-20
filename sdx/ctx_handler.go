package sdx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/mid"
)

func JwtAuthHandler(secret string) fst.CtxHandler {
	return mid.JwtAuthHandler(secret)
}

func PmsParser(ctx *fst.Context) {
	err := ctx.ParseRequestData()
	if err != nil {
		ctx.FaiMsg("PmsParser err: " + err.Error())
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
