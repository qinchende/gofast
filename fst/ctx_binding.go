// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/fst/bind"
	"github.com/qinchende/gofast/skill/mapx"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// GoFast框架自定义的绑定方法，按照GoFast的模式，以前Gin的绑定方式很多都要失效了。

// add by sdx on 20210305
// 就当 c.Pms (c.ReqRaw.Form) 中的是 JSON 对象，我们需要用这个数据源绑定任意的对象
func (c *Context) BindPms(dst any) error {
	// add preBind events by sdx on 2021.03.18
	// c.execPreBindHandlers()
	//return bind.Pms.BindPms(dst, c.Pms)
	return mapx.ApplyKV(dst, c.Pms, mapx.DataOptions())
}

func (c *Context) BindJSON(dst any) error {
	return c.ShouldBindWith(dst, bind.JSON)
}

func (c *Context) BindXML(dst any) error {
	return c.ShouldBindWith(dst, bind.XML)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (c *Context) ShouldBindWith(dst any, b bind.Binding) error {
	// add preBind events by sdx on 2021.03.18
	//c.execPreBindHandlers()
	return b.Bind(c.ReqRaw, dst)
}

// MustBindWith binds the passed struct pointer using the specified binding format.
// It will abort the request with HTTP 400 if any error occurs.
// See the binding package.
func (c *Context) MustBindWith(dst any, b bind.Binding) error {
	if err := c.ShouldBindWith(dst, b); err != nil {
		//c.AbortWithError(http.StatusBadRequest, err).SetType(ErrorTypeBind) // nolint: errcheck
		return err
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 以下来自gin的代码，暂时不需要了
///************************************/
///******* binding and validate *******/
///************************************/
//
//// Bind checks the Content-Type to select a binding myApp automatically,
//// Depending the "Content-Type" header different bindings are used:
////     "application/json" --> JSON binding
////     "application/xml"  --> XML binding
//// otherwise --> returns an error.
//// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
//// It decodes the json payload into the struct specified as a pointer.
//// It writes a 400 error and sets Content-Type header "text/plain" in the response if input is not valid.
//func (c *Context) Bind(dst interface{}) error {
//	b := binding.Default(c.ReqRaw.Method, c.ContentType())
//	return c.MustBindWith(dst, b)
//}

//// BindQuery is a shortcut for c.MustBindWith(dst, binding.Query).
//func (c *Context) BindQuery(dst interface{}) error {
//	return c.MustBindWith(dst, binding.Query)
//}
//
//// BindYAML is a shortcut for c.MustBindWith(dst, binding.YAML).
//func (c *Context) BindYAML(dst interface{}) error {
//	return c.MustBindWith(dst, binding.YAML)
//}
//
//// BindHeader is a shortcut for c.MustBindWith(dst, binding.Header).
//func (c *Context) BindHeader(dst interface{}) error {
//	return c.MustBindWith(dst, binding.Header)
//}
//
//// BindUri binds the passed struct pointer using binding.Uri.
//// It will abort the request with HTTP 400 if any error occurs.
//func (c *Context) BindUri(dst interface{}) error {
//	if err := c.ShouldBindUri(dst); err != nil {
//		c.AbortWithError(http.StatusBadRequest, err).SetType(ErrorTypeBind) // nolint: errcheck
//		return err
//	}
//	return nil
//}

//// ShouldBind checks the Content-Type to select a binding myApp automatically,
//// Depending the "Content-Type" header different bindings are used:
////     "application/json" --> JSON binding
////     "application/xml"  --> XML binding
//// otherwise --> returns an error
//// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
//// It decodes the json payload into the struct specified as a pointer.
//// Like c.Bind() but this method does not set the response status code to 400 and abort if the json is not valid.
//func (c *Context) ShouldBind(dst interface{}) error {
//	b := binding.Default(c.ReqRaw.Method, c.ContentType())
//	return c.ShouldBindWith(dst, b)
//}
//
//// ShouldBindJSON is a shortcut for c.ShouldBindWith(dst, binding.JSON).
//func (c *Context) ShouldBindJSON(dst interface{}) error {
//	return c.ShouldBindWith(dst, binding.JSON)
//}
//
//// ShouldBindXML is a shortcut for c.ShouldBindWith(dst, binding.XML).
//func (c *Context) ShouldBindXML(dst interface{}) error {
//	return c.ShouldBindWith(dst, binding.XML)
//}
//
//// ShouldBindQuery is a shortcut for c.ShouldBindWith(dst, binding.Query).
//func (c *Context) ShouldBindQuery(dst interface{}) error {
//	return c.ShouldBindWith(dst, binding.Query)
//}
//
//// ShouldBindYAML is a shortcut for c.ShouldBindWith(dst, binding.YAML).
//func (c *Context) ShouldBindYAML(dst interface{}) error {
//	return c.ShouldBindWith(dst, binding.YAML)
//}
//
//// ShouldBindHeader is a shortcut for c.ShouldBindWith(dst, binding.Header).
//func (c *Context) ShouldBindHeader(dst interface{}) error {
//	return c.ShouldBindWith(dst, binding.Header)
//}
//
//// ShouldBindUri binds the passed struct pointer using the specified binding myApp.
//func (c *Context) ShouldBindUri(dst interface{}) error {
//	m := make(map[string][]string)
//	for _, v := range *c.match.params {
//		m[v.Key] = []string{v.Value}
//	}
//	return binding.Uri.BindUri(m, dst)
//}

//// ShouldBindBodyWith is similar with ShouldBindWith, but it stores the request
//// body into the context, and reuse when it is called again.
////
//// NOTE: This method reads the body before binding. So you should use
//// ShouldBindWith for better performance if you need to call only once.
//func (c *Context) ShouldBindBodyWith(dst interface{}, bb binding.BindingBody) (err error) {
//	var body []byte
//	if cb, ok := c.Get(BodyBytesKey); ok {
//		if cbb, ok := cb.([]byte); ok {
//			body = cbb
//		}
//	}
//	if body == nil {
//		body, err = ioutil.ReadAll(c.ReqRaw.Body)
//		if err != nil {
//			return err
//		}
//		c.Set(BodyBytesKey, body)
//	}
//	return bb.BindBody(body, dst)
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//
//// Bind implements the `Binder#Bind` function.
//func (b *DefaultBinder) Bind(i interface{}, c Context) (err error) {
//	req := c.Request()
//
//	names := c.ParamNames()
//	values := c.ParamValues()
//	params := map[string][]string{}
//	for i, name := range names {
//		params[name] = []string{values[i]}
//	}
//	if err := b.bindData(i, params, "param"); err != nil {
//		return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
//	}
//	if err = b.bindData(i, c.QueryParams(), "query"); err != nil {
//		return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
//	}
//	if req.ContentLength == 0 {
//		return
//	}
//	ctype := req.Header.Get(HeaderContentType)
//	switch {
//	case strings.HasPrefix(ctype, MIMEApplicationJSON):
//		if err = json.NewDecoder(req.Body).Decode(i); err != nil {
//			if ute, ok := err.(*json.UnmarshalTypeError); ok {
//				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
//			} else if se, ok := err.(*json.SyntaxError); ok {
//				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
//			}
//			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
//		}
//	case strings.HasPrefix(ctype, MIMEApplicationXML), strings.HasPrefix(ctype, MIMETextXML):
//		if err = xml.NewDecoder(req.Body).Decode(i); err != nil {
//			if ute, ok := err.(*xml.UnsupportedTypeError); ok {
//				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unsupported type error: type=%v, error=%v", ute.Type, ute.Error())).SetInternal(err)
//			} else if se, ok := err.(*xml.SyntaxError); ok {
//				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: line=%v, error=%v", se.Line, se.Error())).SetInternal(err)
//			}
//			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
//		}
//	case strings.HasPrefix(ctype, MIMEApplicationForm), strings.HasPrefix(ctype, MIMEMultipartForm):
//		params, err := c.FormParams()
//		if err != nil {
//			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
//		}
//		if err = b.bindData(i, params, "form"); err != nil {
//			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
//		}
//	default:
//		return ErrUnsupportedMediaType
//	}
//	return
//}
//
//func (b *DefaultBinder) bindData(ptr interface{}, data map[string][]string, tag string) error {
//	if ptr == nil || len(data) == 0 {
//		return nil
//	}
//	typ := reflect.TypeOf(ptr).Elem()
//	val := reflect.ValueOf(ptr).Elem()
//
//	// Map
//	if typ.Kind() == reflect.Map {
//		for k, v := range data {
//			val.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v[0]))
//		}
//		return nil
//	}
//
//	// !struct
//	if typ.Kind() != reflect.Struct {
//		return errors.New("binding element must be a struct")
//	}
//
//	for i := 0; i < typ.NumField(); i++ {
//		typeField := typ.Field(i)
//		structField := val.Field(i)
//		if !structField.CanSet() {
//			continue
//		}
//		structFieldKind := structField.Kind()
//		inputFieldName := typeField.Tag.Get(tag)
//
//		if inputFieldName == "" {
//			inputFieldName = typeField.Name
//			// If tag is nil, we inspect if the field is a struct.
//			if _, ok := structField.Addr().Interface().(BindUnmarshaler); !ok && structFieldKind == reflect.Struct {
//				if err := b.bindData(structField.Addr().Interface(), data, tag); err != nil {
//					return err
//				}
//				continue
//			}
//		}
//
//		inputValue, exists := data[inputFieldName]
//		if !exists {
//			// Go json.Unmarshal supports case insensitive binding.  However the
//			// url params are bound case sensitive which is inconsistent.  To
//			// fix this we must check all of the map values in a
//			// case-insensitive search.
//			for k, v := range data {
//				if strings.EqualFold(k, inputFieldName) {
//					inputValue = v
//					exists = true
//					break
//				}
//			}
//		}
//
//		if !exists {
//			continue
//		}
//
//		// Call this first, in case we're dealing with an alias to an array type
//		if ok, err := unmarshalField(typeField.Type.Kind(), inputValue[0], structField); ok {
//			if err != nil {
//				return err
//			}
//			continue
//		}
//
//		numElems := len(inputValue)
//		if structFieldKind == reflect.Slice && numElems > 0 {
//			sliceOf := structField.Type().Elem().Kind()
//			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
//			for j := 0; j < numElems; j++ {
//				if err := setWithProperType(sliceOf, inputValue[j], slice.Index(j)); err != nil {
//					return err
//				}
//			}
//			val.Field(i).Set(slice)
//		} else if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
//			return err
//
//		}
//	}
//	return nil
//}
