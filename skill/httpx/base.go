package httpx

import (
	"github.com/qinchende/gofast/cst"
)

const XForwardFor = "X-Forward-For"

const (
	emptyJson         = "{}"
	maxMemory         = 32 << 20 // 32MB
	maxBodyLen        = 8 << 20  // 8MB
	separator         = ";"
	tokensInAttribute = 2
)

const (
	FormatJson = iota
	FormatUrlEncoding
	FormatXml
)

type RequestPet struct {
	Method     string    // GET or POST
	Url        string    // http(s)地址
	Headers    cst.WebKV // 请求头
	QueryArgs  cst.KV    // url上的参数
	BodyArgs   cst.KV    // body带的参数
	BodyFormat int8      // body数据的格式，比如 json|url-encoding|xml
}
