// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"io"
	"mime/multipart"
	"os"
)

// 返回所有上传文件的Form表单结构
// MultipartForm is the parsed multipart form, including file uploads.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Req.Raw.ParseMultipartForm(c.myApp.WebConfig.MaxMultipartBytes)
	return c.Req.Raw.MultipartForm, err
}

// 查找一个上传的文件
// FormFile returns the first file for the provided form key.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Req.Raw.MultipartForm == nil {
		if err := c.Req.Raw.ParseMultipartForm(c.myApp.WebConfig.MaxMultipartBytes); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.Req.Raw.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

// 指定文件，临时保存上传的文件流
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
