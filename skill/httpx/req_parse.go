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

//func Parse(r *http.Request, v any) error {
//	if err := ParsePath(r, v); err != nil {
//		return err
//	}
//
//	if err := ParseForm(r, v); err != nil {
//		return err
//	}
//
//	return ParseJsonBody(r, v)
//}

//// Parses the form request.
//func ParseForm(r *http.Request, v any) error {
//	if err := r.ParseForm(); err != nil {
//		return err
//	}
//
//	if err := r.ParseMultipartForm(maxMemory); err != nil {
//		if err != http.ErrNotMultipart {
//			return err
//		}
//	}
//
//	params := make(map[string]any, len(r.Form))
//	for name := range r.Form {
//		formValue := r.Form.Get(name)
//		if len(formValue) > 0 {
//			params[name] = formValue
//		}
//	}
//
//	return mapx.ApplyKVOfData(v, params)
//}

//func ParseHeader(headerValue string) map[string]string {
//	ret := make(map[string]string)
//	fields := strings.Split(headerValue, separator)
//
//	for _, field := range fields {
//		field = strings.TrimSpace(field)
//		if len(field) == 0 {
//			continue
//		}
//
//		kv := strings.SplitN(field, "=", tokensInAttribute)
//		if len(kv) != tokensInAttribute {
//			continue
//		}
//
//		ret[kv[0]] = kv[1]
//	}
//
//	return ret
//}

//// Parses the post request which contains json in body.
//func ParseJsonBody(r *http.Request, v any) error {
//	var reader io.Reader
//	if withJsonBody(r) {
//		reader = io.LimitReader(r.Body, maxBodyLen)
//	} else {
//		reader = strings.NewReader(emptyJson)
//	}
//
//	return mapx.DecodeJsonReaderOfData(v, reader)
//}

//// Parses the symbols reside in url path.
//// Like http://localhost/bag/:name
//func ParsePath(r *http.Request, v any) error {
//	vars := Vars(r)
//	kv := make(map[string]any, len(vars))
//	for k, v := range vars {
//		kv[k] = v
//	}
//
//	return mapx.ApplyKVOfData(v, kv)
//}
