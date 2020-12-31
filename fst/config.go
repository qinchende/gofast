package fst

import "gofast/skill"

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 当前环境配置
var FEnv *FConfig

type FConfig struct {
	CurrMode               string
	SecureJsonPrefix       string
	MaxMultipartMemory     int64
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	ForwardedByClientIP    bool
	UseRawPath             bool
	UnescapePathValues     bool
	RemoveExtraSlash       bool
	PrintRouteTrees        bool
	modeType               int8
}

func (fc *FConfig) initServerEnv() {
	//FuncMap: 				template.FuncMap{},
	//RedirectTrailingSlash:  true,
	//RedirectFixedPath:      false,
	//HandleMethodNotAllowed: false,
	//ForwardedByClientIP:    true,
	//AppEngine:              defaultAppEngine,
	//UseRawPath:             false,
	//RemoveExtraSlash:       false,
	//UnescapePathValues:     true,
	//MaxMultipartMemory:     defaultMultipartMemory,
	//trees:                  make(methodTrees, 0, 9),
	//delims:                 render.Delims{Left: "{{", Right: "}}"},
	//secureJsonPrefix:       "while(1);",
	if fc.SecureJsonPrefix == "" {
		fc.SecureJsonPrefix = "while(1);"
	}
	if fc.MaxMultipartMemory == 0 {
		fc.MaxMultipartMemory = 32 << 20 // 32 MB
	}

	FEnv = fc
	SetMode(fc.CurrMode)
	skill.SetDebugStatus(FEnv.modeType == modeDebug)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 当前运行处于啥模式：
const (
	modeDebug int8 = iota
	modeTest
	modeProduct
)

const (
	DebugMode   = "debug"
	TestMode    = "test"
	ProductMode = "product"
)

func IsDebugMode() bool {
	return FEnv.modeType == modeDebug
}

func SetMode(val string) {
	switch val {
	case DebugMode, "":
		FEnv.modeType = modeDebug
	case ProductMode:
		FEnv.modeType = modeProduct
	case TestMode:
		FEnv.modeType = modeTest
	default:
		panic("GoFast mode unknown: " + val)
	}
}
