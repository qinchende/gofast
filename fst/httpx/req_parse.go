package httpx

import (
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/store/jde"
	"net/http"
	"strings"
)

func withJsonBody(r *http.Request) bool {
	return r.ContentLength > 0 && strings.Contains(r.Header.Get(cst.HeaderContentType), cst.MIMEAppJson)
}

func parseJsonResponse(resp *http.Response, err error) (cst.KV, error) {
	if resp == nil || err != nil {
		return nil, err
	}
	kv := cst.KV{}
	if err = jde.DecodeRequest(&kv, resp.Request); err != nil {
		return nil, err
	}
	return kv, err
}
