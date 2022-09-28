package httpx

//const (
//	ApplicationJson = "application/json"
//	ContentEncoding = "Content-Encoding"
//	ContentSecurity = "X-Content-Security"
//	ContentType     = "Content-Type"
//	KeyField        = "key"
//	SecretField     = "secret"
//	TypeField       = "type"
//	CryptionType    = 1
//)
//
//const (
//	CodeSignaturePass = iota
//	CodeSignatureInvalidHeader
//	CodeSignatureWrongTime
//	CodeSignatureInvalidToken
//)

import "errors"

const (
	pathKey   = "path"
	formKey   = "form"
	headerKey = "header"
	jsonKey   = "json"
	slash     = "/"
	colon     = ':'
)

// ErrGetWithBody indicates that GET request with body.
var ErrGetWithBody = errors.New("HTTP GET should not have body")
