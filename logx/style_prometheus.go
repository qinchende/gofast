package logx

import (
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/timex"
)

type logPrometheusEntry struct {
	Timestamp string `json:"@timestamp"`
	Level     string `json:"lv"`
	Duration  string `json:"duration,omitempty"`
	Content   string `json:"ct"`
}

func outputJsonStyle(w WriterCloser, info, logLevel string) {
	logWrap := logPrometheusEntry{
		Timestamp: timex.Time().Format(timeFormat),
		Level:     logLevel,
		Content:   info,
	}
	if content, err := jsonx.Marshal(logWrap); err != nil {
		outputDirectString(w, err.Error())
	} else {
		outputDirectBytes(w, content)
	}
}
