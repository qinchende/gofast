// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mapx

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func DecodeYamlBytes(dst any, content []byte, like int8) error {
	return DecodeYamlBytesX(dst, content, matchOptions(like))
}

func DecodeYamlBytesX(dst any, content []byte, opts *BindOptions) error {
	var o any
	if err := DecodeYaml(&o, content); err != nil {
		return err
	}

	if kv, ok := o.(map[string]any); ok {
		return BindKVX(dst, kv, opts)
	} else {
		return errors.New("only map-like configs supported")
	}
}

func DecodeYamlReader(dst any, reader io.Reader, like int8) error {
	return DecodeYamlReaderX(dst, reader, matchOptions(like))
}

func DecodeYamlReaderX(dst any, reader io.Reader, opts *BindOptions) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return DecodeYamlBytesX(dst, content, opts)
}

// yamlUnmarshal YAML to map[string]interface{} instead of map[interface{}]interface{}.
func DecodeYaml(out any, in []byte) error {
	var res any
	if err := yaml.Unmarshal(in, &res); err != nil {
		return err
	}

	*out.(*any) = cleanupMapValue(res)
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func cleanupInterfaceMap(in map[any]any) map[string]any {
	res := make(map[string]any)
	for k, v := range in {
		res[sdxAsString(k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceNumber(in any) json.Number {
	return json.Number(sdxAsString(in))
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
		return sdxAsString(v)
	}
}
