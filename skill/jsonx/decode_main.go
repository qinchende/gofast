package jsonx

import (
	"errors"
	"github.com/qinchende/gofast/cst"
)

var (
	sErr = errors.New("jsonx: json syntax error.")
)

type fastDecode struct {
	dst  cst.SuperKV
	src  string
	head int
	tail int
	//braces  bracesMark  // 大括号
	//squares squaresMark // 中括号
}

func (dd *fastDecode) init(dst cst.SuperKV, src string) error {
	//if err := dd.searchBrackets(); err != nil {
	//	return err
	//}
	dd.dst = dst
	dd.src = src
	dd.head = 0
	dd.tail = len(dd.src) - 1
	return nil
}

func (dd *fastDecode) warpError(err error) error {
	if err != nil {
		end := dd.head + 5
		if end > dd.tail {
			end = dd.tail
		}
		err = errors.New(err.Error() + " near: " + dd.src[dd.head:end])
	}
	return err
}
