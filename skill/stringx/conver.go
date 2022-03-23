package stringx

import "bytes"

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
