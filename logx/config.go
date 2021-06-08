package logx

type LogConfig struct {
	ServiceName        string `json:",optional"`
	Mode               string `json:",default=console,options=console|file|volume"`
	Path               string `json:",default=logs"`
	Level              string `json:",default=info,options=info|error|severe"`
	Compress           bool   `json:",optional"`
	KeepDays           int    `json:",optional"`
	StackArchiveMillis int    `json:",default=100"`
	NeedCpuMem         bool   `json:",default=true"`
	StyleName          string `json:",default=json"`

	// 内部：日志样式
	style int8 `json:",optional"`
}

const (
	_styleJson    = "json"
	_styleSdx     = "sdx"
	_styleSdxMini = "sdx-mini"
)

const (
	styleJson int8 = iota
	styleSdx
	styleSdxMini
)

var currConfig *LogConfig
