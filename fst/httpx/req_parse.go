package httpx

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
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
	if err = jsonx.UnmarshalFromReader(&kv, resp.Body); err != nil {
		return nil, err
	}
	return kv, err
}
