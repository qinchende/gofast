// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package drops

import (
	"github.com/qinchende/gofast/fst/render"
	"gopkg.in/yaml.v2"
	"net/http"
)

// YAML contains the given interface object.
type YAML struct {
	Data any
}

var yamlContentType = []string{"application/x-yaml; charset=utf-8"}

// Render (YAML) marshals the given interface object and writes data with custom ContentType.
func (r YAML) Write(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := yaml.Marshal(r.Data)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

// WriteContentType (YAML) writes YAML ContentType for response.
func (r YAML) WriteContentType(w http.ResponseWriter) {
	render.writeContentType(w, yamlContentType)
}
