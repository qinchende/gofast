package mapping

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// modify by sdx on 2021.08.30 (change json to cnf)
// To make .json & .yaml consistent, we just use [cnf] as the tag key.
const yamlTagKey = "cnf"

var (
	ErrUnsupportedType = errors.New("only map-like configs are suported")

	yamlUnmarshaler = NewUnmarshaler(yamlTagKey)
)

func UnmarshalYamlBytes(content []byte, v any) error {
	return unmarshalYamlBytes(content, v, yamlUnmarshaler)
}

func UnmarshalYamlReader(reader io.Reader, v any) error {
	return unmarshalYamlReader(reader, v, yamlUnmarshaler)
}

func unmarshalYamlBytes(content []byte, v any, unmarshaler *Unmarshaler) error {
	var o any
	if err := yamlUnmarshal(content, &o); err != nil {
		return err
	}

	if m, ok := o.(map[string]any); ok {
		return unmarshaler.Unmarshal(m, v)
	} else {
		return ErrUnsupportedType
	}
}

func unmarshalYamlReader(reader io.Reader, v any, unmarshaler *Unmarshaler) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	return unmarshalYamlBytes(content, v, unmarshaler)
}

// yamlUnmarshal YAML to map[string]interface{} instead of map[interface{}]interface{}.
func yamlUnmarshal(in []byte, out any) error {
	var res any
	if err := yaml.Unmarshal(in, &res); err != nil {
		return err
	}

	*out.(*any) = cleanupMapValue(res)
	return nil
}

func cleanupInterfaceMap(in map[any]any) map[string]any {
	res := make(map[string]any)
	for k, v := range in {
		res[Repr(k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceNumber(in any) json.Number {
	return json.Number(Repr(in))
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
		return Repr(v)
	}
}
