// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"fmt"
)

// 向对应的文件（描述符）写入日志记录
func logBytes(buf []byte) {
	_, _ = fmt.Fprint(DefaultWriter, buf)
}

func logString(text string) {
	_, _ = fmt.Fprint(DefaultWriter, text)
}
