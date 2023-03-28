// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
)

//// 返回所有上传文件的Form表单结构
//// MultipartForm is the parsed multipart form, including file uploads.
//func (c *Context) MultipartForm() (*multipart.Form, error) {
//	err := httpx.ParseMultipartForm(c.Pms2, c.Req.Raw, c.myApp.WebConfig.MaxMultipartBytes)
//	return c.Req.Raw.MultipartForm, err
//}

// 查找一个上传的文件
// FormFile returns the first file for the provided form key.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	mForm := c.Req.Raw.MultipartForm
	if mForm != nil && mForm.File != nil {
		if fhs := mForm.File[name]; len(fhs) > 0 {
			return fhs[0], nil
		}
	}
	return nil, errors.New("http: no such file")
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
