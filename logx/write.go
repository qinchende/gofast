// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"fmt"
)

func writeBytesNow(text []byte) {
	_, _ = fmt.Fprint(DefaultWriter, text)
}

func writeStringNow(text string) {
	_, _ = fmt.Fprint(DefaultWriter, text)
}
