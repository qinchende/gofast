package mid

import (
	"github.com/qinchende/gofast/fst"
	"net/http"

	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/sysx"
	"github.com/qinchende/gofast/skill/trace"
)

func Tracing(w *fst.GFResponse, r *http.Request) {
	carrier, err := trace.Extract(trace.HttpFormat, r.Header)
	// ErrInvalidCarrier means no trace id was set in http header
	if err != nil && err != trace.ErrInvalidCarrier {
		logx.Error(err)
	}

	ctx, span := trace.StartServerSpan(r.Context(), carrier, sysx.Hostname(), r.RequestURI)
	defer span.Finish()
	r = r.WithContext(ctx)

	w.NextFit(r)
}
