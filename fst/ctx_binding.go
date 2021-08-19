// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/fst/binding"
	"net/http"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// GoFast框架自定义的绑定方法，按照GoFast的模式，以前Gin的绑定方式很多都要失效了。

// add by sdx on 20210305
// 就当 c.Pms (c.ReqRaw.Form) 中的是 JSON 对象，我们需要用这个数据源绑定任意的对象
func (c *Context) BindPms(obj interface{}) error {
	// add preBind events by sdx on 2021.03.18
	//c.execPreBindHandlers()

	return binding.Pms.BindPms(c.Pms, obj)
}

///************************************/
///******* binding and validate *******/
///************************************/
//
//// Bind checks the Content-Type to select a binding gftApp automatically,
//// Depending the "Content-Type" header different bindings are used:
////     "application/json" --> JSON binding
////     "application/xml"  --> XML binding
//// otherwise --> returns an error.
//// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
//// It decodes the json payload into the struct specified as a pointer.
//// It writes a 400 error and sets Content-Type header "text/plain" in the response if input is not valid.
//func (c *Context) Bind(obj interface{}) error {
//	b := binding.Default(c.ReqRaw.Method, c.ContentType())
//	return c.MustBindWith(obj, b)
//}

// BindJSON is a shortcut for c.MustBindWith(obj, binding.JSON).
func (c *Context) BindJSON(obj interface{}) error {
	return c.MustBindWith(obj, binding.JSON)
}

// BindXML is a shortcut for c.MustBindWith(obj, binding.BindXML).
func (c *Context) BindXML(obj interface{}) error {
	return c.MustBindWith(obj, binding.XML)
}

//// BindQuery is a shortcut for c.MustBindWith(obj, binding.Query).
//func (c *Context) BindQuery(obj interface{}) error {
//	return c.MustBindWith(obj, binding.Query)
//}
//
//// BindYAML is a shortcut for c.MustBindWith(obj, binding.YAML).
//func (c *Context) BindYAML(obj interface{}) error {
//	return c.MustBindWith(obj, binding.YAML)
//}
//
//// BindHeader is a shortcut for c.MustBindWith(obj, binding.Header).
//func (c *Context) BindHeader(obj interface{}) error {
//	return c.MustBindWith(obj, binding.Header)
//}
//
//// BindUri binds the passed struct pointer using binding.Uri.
//// It will abort the request with HTTP 400 if any error occurs.
//func (c *Context) BindUri(obj interface{}) error {
//	if err := c.ShouldBindUri(obj); err != nil {
//		c.AbortWithError(http.StatusBadRequest, err).SetType(ErrorTypeBind) // nolint: errcheck
//		return err
//	}
//	return nil
//}

// MustBindWith binds the passed struct pointer using the specified binding gftApp.
// It will abort the request with HTTP 400 if any error occurs.
// See the binding package.
func (c *Context) MustBindWith(obj interface{}, b binding.Binding) error {
	if err := c.ShouldBindWith(obj, b); err != nil {
		c.AbortWithError(http.StatusBadRequest, err).SetType(ErrorTypeBind) // nolint: errcheck
		return err
	}
	return nil
}

//// ShouldBind checks the Content-Type to select a binding gftApp automatically,
//// Depending the "Content-Type" header different bindings are used:
////     "application/json" --> JSON binding
////     "application/xml"  --> XML binding
//// otherwise --> returns an error
//// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
//// It decodes the json payload into the struct specified as a pointer.
//// Like c.Bind() but this method does not set the response status code to 400 and abort if the json is not valid.
//func (c *Context) ShouldBind(obj interface{}) error {
//	b := binding.Default(c.ReqRaw.Method, c.ContentType())
//	return c.ShouldBindWith(obj, b)
//}
//
//// ShouldBindJSON is a shortcut for c.ShouldBindWith(obj, binding.JSON).
//func (c *Context) ShouldBindJSON(obj interface{}) error {
//	return c.ShouldBindWith(obj, binding.JSON)
//}
//
//// ShouldBindXML is a shortcut for c.ShouldBindWith(obj, binding.XML).
//func (c *Context) ShouldBindXML(obj interface{}) error {
//	return c.ShouldBindWith(obj, binding.XML)
//}
//
//// ShouldBindQuery is a shortcut for c.ShouldBindWith(obj, binding.Query).
//func (c *Context) ShouldBindQuery(obj interface{}) error {
//	return c.ShouldBindWith(obj, binding.Query)
//}
//
//// ShouldBindYAML is a shortcut for c.ShouldBindWith(obj, binding.YAML).
//func (c *Context) ShouldBindYAML(obj interface{}) error {
//	return c.ShouldBindWith(obj, binding.YAML)
//}
//
//// ShouldBindHeader is a shortcut for c.ShouldBindWith(obj, binding.Header).
//func (c *Context) ShouldBindHeader(obj interface{}) error {
//	return c.ShouldBindWith(obj, binding.Header)
//}
//
//// ShouldBindUri binds the passed struct pointer using the specified binding gftApp.
//func (c *Context) ShouldBindUri(obj interface{}) error {
//	m := make(map[string][]string)
//	for _, v := range *c.match.params {
//		m[v.Key] = []string{v.Value}
//	}
//	return binding.Uri.BindUri(m, obj)
//}

// ShouldBindWith binds the passed struct pointer using the specified binding gftApp.
// See the binding package.
func (c *Context) ShouldBindWith(obj interface{}, b binding.Binding) error {
	// add preBind events by sdx on 2021.03.18
	//c.execPreBindHandlers()

	return b.Bind(c.ReqRaw, obj)
}

//// ShouldBindBodyWith is similar with ShouldBindWith, but it stores the request
//// body into the context, and reuse when it is called again.
////
//// NOTE: This method reads the body before binding. So you should use
//// ShouldBindWith for better performance if you need to call only once.
//func (c *Context) ShouldBindBodyWith(obj interface{}, bb binding.BindingBody) (err error) {
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
//	return bb.BindBody(body, obj)
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
//
//func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
//	// But also call it here, in case we're dealing with an array of BindUnmarshalers
//	if ok, err := unmarshalField(valueKind, val, structField); ok {
//		return err
//	}
//
//	switch valueKind {
//	case reflect.Ptr:
//		return setWithProperType(structField.Elem().Kind(), val, structField.Elem())
//	case reflect.Int:
//		return setIntField(val, 0, structField)
//	case reflect.Int8:
//		return setIntField(val, 8, structField)
//	case reflect.Int16:
//		return setIntField(val, 16, structField)
//	case reflect.Int32:
//		return setIntField(val, 32, structField)
//	case reflect.Int64:
//		return setIntField(val, 64, structField)
//	case reflect.Uint:
//		return setUintField(val, 0, structField)
//	case reflect.Uint8:
//		return setUintField(val, 8, structField)
//	case reflect.Uint16:
//		return setUintField(val, 16, structField)
//	case reflect.Uint32:
//		return setUintField(val, 32, structField)
//	case reflect.Uint64:
//		return setUintField(val, 64, structField)
//	case reflect.Bool:
//		return setBoolField(val, structField)
//	case reflect.Float32:
//		return setFloatField(val, 32, structField)
//	case reflect.Float64:
//		return setFloatField(val, 64, structField)
//	case reflect.String:
//		structField.SetString(val)
//	default:
//		return errors.New("unknown type")
//	}
//	return nil
//}
//
//func unmarshalField(valueKind reflect.Kind, val string, field reflect.Value) (bool, error) {
//	switch valueKind {
//	case reflect.Ptr:
//		return unmarshalFieldPtr(val, field)
//	default:
//		return unmarshalFieldNonPtr(val, field)
//	}
//}
//
//func unmarshalFieldNonPtr(value string, field reflect.Value) (bool, error) {
//	fieldIValue := field.Addr().Interface()
//	if unmarshaler, ok := fieldIValue.(BindUnmarshaler); ok {
//		return true, unmarshaler.UnmarshalParam(value)
//	}
//	if unmarshaler, ok := fieldIValue.(encoding.TextUnmarshaler); ok {
//		return true, unmarshaler.UnmarshalText([]byte(value))
//	}
//
//	return false, nil
//}
//
//func unmarshalFieldPtr(value string, field reflect.Value) (bool, error) {
//	if field.IsNil() {
//		// Initialize the pointer to a nil value
//		field.Set(reflect.New(field.Type().Elem()))
//	}
//	return unmarshalFieldNonPtr(value, field.Elem())
//}
//
//func setIntField(value string, bitSize int, field reflect.Value) error {
//	if value == "" {
//		value = "0"
//	}
//	intVal, err := strconv.ParseInt(value, 10, bitSize)
//	if err == nil {
//		field.SetInt(intVal)
//	}
//	return err
//}
//
//func setUintField(value string, bitSize int, field reflect.Value) error {
//	if value == "" {
//		value = "0"
//	}
//	uintVal, err := strconv.ParseUint(value, 10, bitSize)
//	if err == nil {
//		field.SetUint(uintVal)
//	}
//	return err
//}
//
//func setBoolField(value string, field reflect.Value) error {
//	if value == "" {
//		value = "false"
//	}
//	boolVal, err := strconv.ParseBool(value)
//	if err == nil {
//		field.SetBool(boolVal)
//	}
//	return err
//}
//
//func setFloatField(value string, bitSize int, field reflect.Value) error {
//	if value == "" {
//		value = "0.0"
//	}
//	floatVal, err := strconv.ParseFloat(value, bitSize)
//	if err == nil {
//		field.SetFloat(floatVal)
//	}
//	return err
//}
//
