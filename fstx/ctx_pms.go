package fstx

import "github.com/qinchende/gofast/fst"

func PmsParser(ctx *fst.Context) {
	ctx.ParseRequestData()
}

func BuildPmsOfJson(ctx *fst.Context) {
	ctx.GenPmsByJSONBody()
}

func BuildPmsOfXml(ctx *fst.Context) {
	ctx.GenPmsByXMLBody()
}

func BuildPmsOfForm(ctx *fst.Context) {
	ctx.GenPmsByFormBody()
}
