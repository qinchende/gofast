package fst

import (
	"fmt"
	"gofast/fst/binding"
)

type I interface{}
type KV map[string]I
type CHandler func(*Context)
type CHandlers []CHandler

type FHandler func(*Faster)
type FHandlers []FHandler

const (
	routePathMaxLen    uint8 = 255
	routeMaxHandlers   uint8 = 255
	defMultipartMemory int   = 32 << 20 // 32 MB
)

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = binding.MIMEJSON
	MIMEHTML              = binding.MIMEHTML
	MIMEXML               = binding.MIMEXML
	MIMEXML2              = binding.MIMEXML2
	MIMEPlain             = binding.MIMEPlain
	MIMEPOSTForm          = binding.MIMEPOSTForm
	MIMEMultipartPOSTForm = binding.MIMEMultipartPOSTForm
	MIMEYAML              = binding.MIMEYAML
)

var (
	spf            = fmt.Sprintf
	mimePlain      = []string{MIMEPlain}
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
)
