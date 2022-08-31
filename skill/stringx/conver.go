package stringx

import (
	"bytes"
	"github.com/qinchende/gofast/skill/lang"
)

func Camel2Snake(s string) string {
	newS := bytes.Buffer{}
	for i := 0; i < len(s); i++ {
		if s[i] >= 65 && s[i] <= 90 {
			if i > 0 && s[i-1] >= 97 && s[i-1] <= 122 {
				newS.WriteByte('_')
			}
			newS.WriteByte(s[i] + 32)
		} else {
			newS.WriteByte(s[i])
		}
	}
	return newS.String()
}

// ToString 获取变量的字符串值
// 浮点型 3.0 将会转换成字符串3 -> "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func ToString(v any) (s string) {
	s, _ = lang.ToString(v)
	return
}
