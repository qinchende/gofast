package jsonx

import (
	"fmt"
)

func formatError(v string, err error) error {
	return fmt.Errorf("string: `%s`, error: `%s`", v, err.Error())
}
