// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fst

import (
	"encoding/xml"
	"gofast/skill"
	"net/http"
	"os"
	"path"
	"reflect"
	"runtime"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type onlyFilesFS struct {
	fs http.FileSystem
}

type neuteredReaddirFile struct {
	http.File
}

// Dir returns a http.Filesystem that can be used by http.FileServer(). It is used internally
// in router.Static().
// if listDirectory == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, listDirectory bool) http.FileSystem {
	fs := http.Dir(root)
	if listDirectory {
		return fs
	}
	return &onlyFilesFS{fs}
}

// Open conforms to http.Filesystem.
func (fs onlyFilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

// Readdir overrides the http.File default implementation.
func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// BindKey indicates a default bind key.
//const BindKey = "_gin-gonic/gin/bindkey"
//
//// Bind is a Kelper function for given interface object and returns a Gin middleware.
//func Bind(val interface{}) CHandler {
//	value := reflect.ValueOf(val)
//	if value.Kind() == reflect.Ptr {
//		panic(`Bind struct can not be a pointer. Example:
//	Use: gin.Bind(Struct{}) instead of gin.Bind(&Struct{})
//`)}
//	typ := value.Type()
//
//	return func(c *Context) {
//		obj := reflect.New(typ).Interface()
//		if c.Bind(obj) == nil {
//			c.Set(BindKey, obj)
//		}
//	}
//}
//
//// WrapF is a helper function for wrapping http.CHandler and returns a Gin middleware.
//func WrapF(f http.HandlerFunc) CHandler {
//	return func(c *Context) {
//		f(c.Reply, c.Request)
//	}
//}
//
//// WrapH is a helper function for wrapping http.Handler and returns a Gin middleware.
//func WrapH(h http.Handler) CHandler {
//	return func(c *Context) {
//		h.ServeHTTP(c.Reply, c.Request)
//	}
//}

// MarshalXML allows type KV to be used with xml.Marshal.
func (h KV) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func ifPanic(yn bool, text string) {
	if yn {
		panic(text)
	}
}

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

//
//func chooseData(custom, wildcard interface{}) interface{} {
//	if custom == nil {
//		if wildcard == nil {
//			panic("negotiation config is invalid")
//		}
//		return wildcard
//	}
//	return custom
//}
//
//func parseAccept(acceptHeader string) []string {
//	parts := strings.Split(acceptHeader, ",")
//	out := make([]string, 0, len(parts))
//	for _, part := range parts {
//		if part = strings.TrimSpace(strings.Split(part, ";")[0]); part != "" {
//			out = append(out, part)
//		}
//	}
//	return out
//}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			skill.DebugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		skill.DebugPrint("Environment variable PORT is undefined. Using port :8099 by default")
		return ":8099"
	case 1:
		return addr[0]
	default:
		panic("Too many parameters")
	}
}

func nameOfFunc(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
