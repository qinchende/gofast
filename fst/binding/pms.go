// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

type pmsBinding struct{}

func (pmsBinding) Name() string {
	return "pms"
}

// add by sdx on 20210305
func (pmsBinding) BindPms(values map[string]interface{}, obj interface{}) error {
	if err := mapPms(obj, values); err != nil {
		return err
	}

	// TODO: before validate.
	return validate(obj)
}
