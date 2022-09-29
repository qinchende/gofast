package httpx

import (
	"bytes"
	"context"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/lang"
	"io"
	"net/http"
	"net/url"
)

func ParseJsonResponse(resp *http.Response, err error) (cst.KV, error) {
	if resp == nil || err != nil {
		return nil, err
	}
	kv := cst.KV{}
	if err = jsonx.UnmarshalFromReader(&kv, resp.Body); err != nil {
		return nil, err
	}
	return kv, err
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func DoGetKV(req *http.Request) (cst.KV, error) {
	return ParseJsonResponse(http.DefaultClient.Do(req))
}

func NewRequest(pet *RequestPet) (*http.Request, error) {
	return NewRequestCtx(context.Background(), pet)
}

func NewRequestCtx(ctx context.Context, pet *RequestPet) (*http.Request, error) {
	return buildRequest(ctx, pet)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func DoRequest(pet *RequestPet) (*http.Response, error) {
	return DoRequestCtx(context.Background(), pet)
}

func DoRequestCtx(ctx context.Context, pet *RequestPet) (*http.Response, error) {
	if req, err := buildRequest(ctx, pet); err != nil {
		return nil, err
	} else {
		return http.DefaultClient.Do(req)
	}
}

func DoRequestGetKV(pet *RequestPet) (cst.KV, error) {
	return DoRequestGetKVCtx(context.Background(), pet)
}

func DoRequestGetKVCtx(ctx context.Context, pet *RequestPet) (cst.KV, error) {
	// TODO：根据不同的 content-type 解析数据到 cst.KV 形式
	return ParseJsonResponse(DoRequestCtx(ctx, pet))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func checkRequestPet(pet *RequestPet) {
	if pet.Method == "" {
		pet.Method = http.MethodPost
	}
}

func fillQueryArgs(u *url.URL, args cst.KV) {
	if args == nil {
		return
	}

	query := u.Query()
	for k, v := range args {
		if v == nil {
			continue
		}
		query.Add(k, lang.ToString(v))
	}
	u.RawQuery = query.Encode()
}

func fillHeader(r *http.Request, pet *RequestPet) {
	if pet.Headers != nil {
		for k, v := range pet.Headers {
			r.Header.Add(k, v)
		}
	}

	switch pet.BodyFormat {
	case FormatJson:
		r.Header.Set(cst.HeaderContentType, cst.MIMEAppJsonUTF8)
	case FormatUrlEncoding:
		r.Header.Set(cst.HeaderContentType, cst.MIMEPostFormUTF8)
	case FormatXml:
		r.Header.Set(cst.HeaderContentType, cst.MIMEAppXmlUTF8)
	default:
		r.Header.Set(cst.HeaderContentType, cst.MIMEPlainUTF8)
	}
}

func buildBody(pet *RequestPet) io.Reader {
	if pet.BodyArgs == nil {
		return nil
	}

	var buf bytes.Buffer
	switch pet.BodyFormat {
	case FormatJson:
		_ = jsonx.NewEncoder(&buf).Encode(pet.BodyArgs)
	case FormatUrlEncoding:
		data := url.Values{}
		for k, v := range pet.BodyArgs {
			if v == nil {
				continue
			}
			data.Add(k, lang.ToString(v))
		}
		buf.WriteString(data.Encode())
	case FormatXml:
	default:

	}
	return &buf
}

func buildRequest(ctx context.Context, pet *RequestPet) (*http.Request, error) {
	checkRequestPet(pet)

	req, err := http.NewRequestWithContext(ctx, pet.Method, pet.Url, buildBody(pet))
	fillQueryArgs(req.URL, pet.QueryArgs)
	fillHeader(req, pet)

	return req, err
}
