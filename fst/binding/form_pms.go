// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "net/url"

// add by sdx on 20210305
func (jsonBinding) BindPms(values url.Values, obj interface{}) error {
	if err := mapForm(obj, values); err != nil {
		return err
	}

	// TODO: before validate.
	return validate(obj)
}
