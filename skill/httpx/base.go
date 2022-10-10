package httpx

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
