// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package dts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qinchende/gofast/skill/iox"
	"gopkg.in/yaml.v2"
	"io"
	"reflect"
	"strconv"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func BindYamlBytes(dst any, content []byte, like int8) error {
	return BindYamlBytesX(dst, content, AsOptions(like))
}

func BindYamlBytesX(dst any, content []byte, opts *BindOptions) error {
	var res any
	if err := UnmarshalYamlBytes(&res, content); err != nil {
		return err
	}

	if kvs, ok := res.(map[string]any); ok {
		return BindKVX(dst, kvs, opts)
	} else {
		return errors.New("only map-like configs supported")
	}
}

func BindYamlReader(dst any, reader io.Reader, like int8) error {
	return BindYamlReaderX(dst, reader, AsOptions(like))
}

func BindYamlReaderX(dst any, reader io.Reader, opts *BindOptions) error {
	content, err := iox.ReadAll(reader, 0)
	if err != nil {
		return err
	}
	return BindYamlBytesX(dst, content, opts)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// yamlUnmarshal YAML to map[string]interface{} instead of map[interface{}]interface{}.
func UnmarshalYamlBytes(dest any, content []byte) error {
	var res any
	if err := yaml.Unmarshal(content, &res); err != nil {
		return err
	}

	*dest.(*any) = cleanupMapValue(res)
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++
func cleanupInterfaceMap(in map[any]any) map[string]any {
	res := make(map[string]any)
	for k, v := range in {
		res[convToString(k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceNumber(in any) json.Number {
	return json.Number(convToString(in))
}

func cleanupInterfaceSlice(in []any) []any {
	res := make([]any, len(in))
	for i, v := range in {
		res[i] = cleanupMapValue(v)
	}
	return res
}

func cleanupMapValue(v any) any {
	switch v := v.(type) {
	case []any:
		return cleanupInterfaceSlice(v)
	case map[any]any:
		return cleanupInterfaceMap(v)
	case bool, string:
		return v
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return cleanupInterfaceNumber(v)
	default:
		return convToString(v)
	}
}

func convToString(src any) string {
	if src == nil {
		return ""
	}

	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	sv := reflect.ValueOf(src)
	switch sv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(sv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(sv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(sv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(sv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(sv.Bool())
	}
	//return fmt.Sprint("%v", src)
	return fmt.Sprint(src)
}
