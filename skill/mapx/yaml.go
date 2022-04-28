package mapx

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decodeYamlReader(dst interface{}, reader io.Reader) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	return decodeYamlBytes(dst, content)
}

func decodeYamlBytes(dst interface{}, content []byte) error {
	var o interface{}
	if err := decodeYaml(&o, content); err != nil {
		return err
	}

	if m, ok := o.(map[string]interface{}); ok {
		return ApplyKVByNameWithDef(dst, m)
	} else {
		return errNotKVType
	}
}

// yamlUnmarshal YAML to map[string]interface{} instead of map[interface{}]interface{}.
func decodeYaml(out interface{}, in []byte) error {
	var res interface{}
	if err := yaml.Unmarshal(in, &res); err != nil {
		return err
	}

	*out.(*interface{}) = cleanupMapValue(res)
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[Repr(k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceNumber(in interface{}) json.Number {
	return json.Number(Repr(in))
}

func cleanupInterfaceSlice(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = cleanupMapValue(v)
	}
	return res
}

func cleanupMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanupInterfaceSlice(v)
	case map[interface{}]interface{}:
		return cleanupInterfaceMap(v)
	case bool, string:
		return v
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return cleanupInterfaceNumber(v)
	default:
		return Repr(v)
	}
}
