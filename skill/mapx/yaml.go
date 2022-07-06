package mapx

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func DecodeYamlReader(dst any, reader io.Reader) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	return DecodeYamlBytes(dst, content)
}

func DecodeYamlBytes(dst any, content []byte) error {
	var o any
	if err := DecodeYaml(&o, content); err != nil {
		return err
	}

	if m, ok := o.(map[string]any); ok {
		return ApplyKVByNameWithDef(dst, m)
	} else {
		return errNotKVType
	}
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
