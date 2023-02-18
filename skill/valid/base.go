package valid

const (
	attrRequired = "required" // 必选 或者 可空
	attrDefault  = "def"
	attrEnum     = "enum"
	attrRange    = "range"
	attrLength   = "len"
	attrRegex    = "regex"
	attrMatch    = "match" // email,mobile,ipv4,ipv4:port,ipv6,id_card,url,file,base64,time,datetime

	// 常用关键字
	itemSeparator = "|"
	equalToken    = "="
)

type (
	FieldOpts struct {
		Range    *numRange
		Enum     []string
		Len      *numRange
		Regex    string
		Match    string
		DefValue string
		Required bool
	}

	numRange struct {
		left     float64
		right    float64
		lInclude bool
		rInclude bool
	}
)
