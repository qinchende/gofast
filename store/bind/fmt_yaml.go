// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package bind

import (
	"encoding/json"
	"errors"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/iox"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/store/dts"
	"gopkg.in/yaml.v2"
	"io"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func BindYamlBytes(dst any, content []byte, model int8) error {
	return BindYamlBytesX(dst, content, dts.AsOptions(model))
}

func BindYamlBytesX(dst any, content []byte, opts *dts.BindOptions) error {
	var res any
	if err := UnmarshalYamlBytes(&res, content); err != nil {
		return err
	}

	if kvs, ok := res.(map[string]any); ok {
		return BindKVX(dst, cst.KV(kvs), opts)
	} else {
		return errors.New("only map[string]any type data supported")
	}
}

func BindYamlReader(dst any, reader io.Reader, model int8) error {
	return BindYamlReaderX(dst, reader, dts.AsOptions(model))
}

func BindYamlReaderX(dst any, reader io.Reader, opts *dts.BindOptions) error {
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

func cleanupMapValue(value any) any {
	switch val := value.(type) {
	case bool, string:
		return val
	case []any:
		res := make([]any, len(val))
		for i, v := range val {
			res[i] = cleanupMapValue(v)
		}
		return res
	case map[any]any:
		res := make(map[string]any)
		for k, v := range val {
			res[lang.ToString(k)] = cleanupMapValue(v)
		}
		return res
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return json.Number(lang.ToString(val))
	default:
		return lang.ToString(val)
	}
}
