// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bind

import "github.com/qinchende/gofast/cst"

type pmsBinding struct{}

func (pmsBinding) Name() string {
	return "pms"
}

// add by sdx on 20210305
func (pmsBinding) BindPms(dest interface{}, values cst.KV) error {
	if err := mapPms(dest, values); err != nil {
		return err
	}

	// TODO: before validate.
	return validate(dest)
}
