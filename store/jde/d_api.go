package jde

import (
	"github.com/qinchende/gofast/skill/iox"
	"github.com/qinchende/gofast/skill/lang"
	"io"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decodeFromReader(dst any, reader io.Reader, ctSize int64) error {
	// 一次性读取完成，或者遇到EOF标记或者其它错误
	if ctSize > maxJsonLength {
		ctSize = maxJsonLength
	}
	bytes, err1 := iox.ReadAll(reader, ctSize)
	if err1 != nil {
		return err1
	}
	return decodeFromString(dst, lang.BTS(bytes))
}

func decodeFromString(dst any, source string) error {
	if len(source) > maxJsonLength {
		return errJsonTooLarge
	}

	dd := fastDecode{}
	if err := dd.init(dst, source); err != nil {
		return err
	}
	return dd.warpError(dd.parseJson())
}
