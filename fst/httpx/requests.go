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

func Do(req *http.Request) (*http.Response, error) {
	return myClient.Do(req)
}

func DoGetKV(req *http.Request) (cst.KV, error) {
	return parseJsonResponse(myClient.Do(req))
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
		return ClientByPet(pet).Do(req)
	}
}

func DoRequestGetKV(pet *RequestPet) (cst.KV, error) {
	return DoRequestGetKVCtx(context.Background(), pet)
}

func DoRequestGetKVCtx(ctx context.Context, pet *RequestPet) (cst.KV, error) {
	// TODO：根据不同的 content-type 解析数据到 cst.KV 形式
	return parseJsonResponse(DoRequestCtx(ctx, pet))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func checkRequestPet(pet *RequestPet) {
	if pet.Method == "" {
		if pet.BodyArgs != nil {
			pet.Method = http.MethodPost
		} else {
			pet.Method = http.MethodGet
		}
	}
}

func fillQueryArgs(u *url.URL, args cst.WebKV) {
	if args == nil {
		return
	}

	query := u.Query()
	for k, v := range args {
		if v == "" {
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

	if r.Header.Get(cst.HeaderContentType) == "" {
		ctValue := cst.MIMEPlainUTF8
		switch {
		case pet.Method == http.MethodGet:
			ctValue = cst.MIMEPlainUTF8
		case pet.BodyFormat == FormatJson:
			ctValue = cst.MIMEAppJsonUTF8
		case pet.BodyFormat == FormatUrlEncoding:
			ctValue = cst.MIMEPostFormUTF8
		case pet.BodyFormat == FormatXml:
			ctValue = cst.MIMEAppXmlUTF8
		}
		r.Header.Set(cst.HeaderContentType, ctValue)
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
			if v == "" {
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
