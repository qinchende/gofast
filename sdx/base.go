// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

const (
	PmsToken = "tok"
)

var NeedPms = []string{PmsToken}

type UsePms struct {
	Tok string `v:"len=[0:128]"`
}
