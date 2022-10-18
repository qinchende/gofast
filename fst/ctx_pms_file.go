// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"io"
	"mime/multipart"
	"os"
)

// FormFile returns the first file for the provided form key.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.ReqRaw.MultipartForm == nil {
		if err := c.ReqRaw.ParseMultipartForm(c.myApp.WebConfig.MaxMultipartBytes); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.ReqRaw.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

// MultipartForm is the parsed multipart form, including file uploads.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.ReqRaw.ParseMultipartForm(c.myApp.WebConfig.MaxMultipartBytes)
	return c.ReqRaw.MultipartForm, err
}

// SaveUploadedFile uploads the form file to specific dst.
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
