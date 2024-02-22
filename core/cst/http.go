package cst

const (
	CharsetUTF8 = "charset=utf-8"
	// PROPFIND Method can be used on collection and property resources.
	Propfind = "PROPFIND"
	// REPORT Method can be used to get information about a resource, see rfc 3253
	Report = "REPORT"
)

// Content-Type MIME of the most common data formats.
const (
	MIMEHtml          = "text/html"
	MIMEPlain         = "text/plain"
	MIMEXml           = "text/xml"
	MIMEAppJavascript = "application/javascript"
	MIMEAppXml        = "application/xml"
	MIMEAppJson       = "application/json"
	MIMEPostForm      = "application/x-www-form-urlencoded"
	MIMEMultiPostForm = "multipart/form-data"
	MIMEMultiMixed    = "multipart/mixed"
	MIMEProtoBuf      = "application/x-protobuf"
	MIMEMsgPack       = "application/msgpack"
	MIMEXMsgPack      = "application/x-msgpack"
	MIMEYaml          = "application/x-yaml"
)

// MIME types + CharsetUTF8
const (
	MIMEHtmlUTF8          = MIMEHtml + "; " + CharsetUTF8
	MIMEPlainUTF8         = MIMEPlain + "; " + CharsetUTF8
	MIMEXmlUTF8           = MIMEXml + "; " + CharsetUTF8
	MIMEAppJavascriptUTF8 = MIMEAppJavascript + "; " + CharsetUTF8
	MIMEAppXmlUTF8        = MIMEAppXml + "; " + CharsetUTF8
	MIMEAppJsonUTF8       = MIMEAppJson + "; " + CharsetUTF8
	MIMEPostFormUTF8      = MIMEPostForm + "; " + CharsetUTF8
	MIMEMultiPostFormUTF8 = MIMEMultiPostForm + "; " + CharsetUTF8
	MIMEMultiMixedUTF8    = MIMEMultiMixed + "; " + CharsetUTF8
)

// Headers
const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderXRequestedWith      = "X-Requested-With"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderXCSRFToken                      = "X-CSRF-Token"
	HeaderReferrerPolicy                  = "Referrer-Policy"
)
