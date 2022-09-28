package httpx

import (
	"bytes"
	"context"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/lang"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type (
	Client interface {
		Do() (*http.Response, error)
	}

	superClient struct {
		req *http.Request
	}
)

func (client *superClient) Do() (*http.Response, error) {
	return http.DefaultClient.Do(client.req)
}

func NewRequest(method, url string, data any) (Client, error) {
	return NewRequestCtx(context.Background(), method, url, data)
}
func NewRequestCtx(ctx context.Context, method, url string, data any) (Client, error) {
	r, err := buildRequest(ctx, method, url, data)
	return &superClient{req: r}, err
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func DoRequest(method, url string, data any) (*http.Response, error) {
	return DoRequestCtx(context.Background(), method, url, data)
}

func DoRequestCtx(ctx context.Context, method, url string, data any) (*http.Response, error) {
	if req, err := buildRequest(ctx, method, url, data); err != nil {
		return nil, err
	} else {
		return request(req)
	}
}

//
//// Do sends an HTTP request with the given arguments and returns an HTTP response.
//// data is automatically marshal into a *httpRequest, typically it's defined in an API file.
//func Do(ctx context.Context, method, url string, data any) (*http.Response, error) {
//	req, err := buildRequest(ctx, method, url, data)
//	if err != nil {
//		return nil, err
//	}
//
//	return DoRequest(req)
//}

//// DoRequest sends an HTTP request and returns an HTTP response.
//func DoRequest(r *http.Request) (*http.Response, error) {
//	return request(r, defaultClient{})
//}

func buildFormQuery(u *url.URL, val map[string]any) string {
	query := u.Query()
	for k, v := range val {
		query.Add(k, fmt.Sprint(v))
	}

	return query.Encode()
}

func buildRequest(ctx context.Context, method, rUrl string, data any) (*http.Request, error) {
	u, err := url.Parse(rUrl)
	if err != nil {
		return nil, err
	}

	var val map[string]map[string]any
	if data != nil {
		//val, err = mapping.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	if err := fillPath(u, val[pathKey]); err != nil {
		return nil, err
	}

	var reader io.Reader
	jsonVars, hasJsonBody := val[jsonKey]
	if hasJsonBody {
		if method == http.MethodGet {
			return nil, ErrGetWithBody
		}

		var buf bytes.Buffer
		enc := jsonx.NewEncoder(&buf)
		if err := enc.Encode(jsonVars); err != nil {
			return nil, err
		}

		reader = &buf
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = buildFormQuery(u, val[formKey])
	fillHeader(req, val[headerKey])
	if hasJsonBody {
		req.Header.Set(cst.HeaderContentType, cst.MIMEAppJsonUTF8)
	}

	return req, nil
}

func fillHeader(r *http.Request, val map[string]any) {
	for k, v := range val {
		r.Header.Add(k, fmt.Sprint(v))
	}
}

func fillPath(u *url.URL, val map[string]any) error {
	used := make(map[string]lang.PlaceholderType)
	fields := strings.Split(u.Path, slash)

	for i := range fields {
		field := fields[i]
		if len(field) > 0 && field[0] == colon {
			name := field[1:]
			ival, ok := val[name]
			if !ok {
				return fmt.Errorf("missing path variable %q", name)
			}
			value := fmt.Sprint(ival)
			if len(value) == 0 {
				return fmt.Errorf("empty path variable %q", name)
			}
			fields[i] = value
			used[name] = lang.Placeholder
		}
	}

	if len(val) != len(used) {
		for key := range used {
			delete(val, key)
		}

		var unused []string
		for key := range val {
			unused = append(unused, key)
		}

		return fmt.Errorf("more path variables are provided: %q", strings.Join(unused, ", "))
	}

	u.Path = strings.Join(fields, slash)
	return nil
}

func request(r *http.Request) (*http.Response, error) {
	//tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	//propagator := otel.GetTextMapPropagator()

	//spanName := r.URL.Path
	//ctx, span := tracer.Start(
	//	r.Context(),
	//	spanName,
	//	oteltrace.WithSpanKind(oteltrace.SpanKindClient),
	//	oteltrace.WithAttributes(semconv.HTTPClientAttributesFromHTTPRequest(r)...),
	//)
	//defer span.End()

	//respHandlers := make([]internal.ResponseHandler, len(interceptors))
	//for i, interceptor := range interceptors {
	//	var h internal.ResponseHandler
	//	r, h = interceptor(r)
	//	respHandlers[i] = h
	//}

	//clientTrace := httptrace.ContextClientTrace(ctx)
	//if clientTrace != nil {
	//	ctx = httptrace.WithClientTrace(ctx, clientTrace)
	//}

	//r = r.WithContext(ctx)
	//propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))

	resp, err := http.DefaultClient.Do(r)
	//for i := len(respHandlers) - 1; i >= 0; i-- {
	//	respHandlers[i](resp, err)
	//}

	//if err != nil {
	//	span.RecordError(err)
	//	span.SetStatus(codes.Error, err.Error())
	//	return resp, err
	//}

	//span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(resp.StatusCode)...)
	//span.SetStatus(semconv.SpanStatusFromHTTPStatusCode(resp.StatusCode))

	return resp, err
}
