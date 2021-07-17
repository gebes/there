// All rights belong to the owners of the http package

package there

import (
	"errors"
	"strings"
)

var ErrorNotRunning = errors.New("there: Server is not running")

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

var (
	AllMethods       = []string{MethodGet, MethodHead, MethodPost, MethodPut, MethodPatch, MethodDelete, MethodConnect, MethodOptions, MethodTrace}
	AllMethodsString = strings.Join(AllMethods, ",")
)

// HTTP status codes as registered with IANA.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
const (
	StatusContinue           = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1
	StatusEarlyHints         = 103 // RFC 8297

	StatusOK                   = 200 // RFC 7231, 6.3.1
	StatusCreated              = 201 // RFC 7231, 6.3.2
	StatusAccepted             = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
	StatusNoContent            = 204 // RFC 7231, 6.3.5
	StatusResetContent         = 205 // RFC 7231, 6.3.6
	StatusPartialContent       = 206 // RFC 7233, 4.1
	StatusMultiStatus          = 207 // RFC 4918, 11.1
	StatusAlreadyReported      = 208 // RFC 5842, 7.1
	StatusIMUsed               = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  = 301 // RFC 7231, 6.4.2
	StatusFound             = 302 // RFC 7231, 6.4.3
	StatusSeeOther          = 303 // RFC 7231, 6.4.4
	StatusNotModified       = 304 // RFC 7232, 4.1
	StatusUseProxy          = 305 // RFC 7231, 6.4.5
	_                       = 306 // RFC 7231, 6.4.6 (Unused)
	StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect = 308 // RFC 7538, 3

	StatusBadRequest                   = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 = 401 // RFC 7235, 3.1
	StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
	StatusForbidden                    = 403 // RFC 7231, 6.5.3
	StatusNotFound                     = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
	StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
	StatusConflict                     = 409 // RFC 7231, 6.5.8
	StatusGone                         = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
	StatusExpectationFailed            = 417 // RFC 7231, 6.5.14
	StatusTeapot                       = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest           = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity          = 422 // RFC 4918, 11.2
	StatusLocked                       = 423 // RFC 4918, 11.3
	StatusFailedDependency             = 424 // RFC 4918, 11.4
	StatusTooEarly                     = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         = 428 // RFC 6585, 3
	StatusTooManyRequests              = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

	StatusInternalServerError           = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5
	StatusLoopDetected                  = 508 // RFC 5842, 7.2
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)

var statusText = map[int]string{
	StatusContinue:           "Continue",
	StatusSwitchingProtocols: "Switching Protocols",
	StatusProcessing:         "Processing",
	StatusEarlyHints:         "Early Hints",

	StatusOK:                   "OK",
	StatusCreated:              "Created",
	StatusAccepted:             "Accepted",
	StatusNonAuthoritativeInfo: "Non-Authoritative Information",
	StatusNoContent:            "No Content",
	StatusResetContent:         "Reset Content",
	StatusPartialContent:       "Partial Content",
	StatusMultiStatus:          "Multi-Status",
	StatusAlreadyReported:      "Already Reported",
	StatusIMUsed:               "IM Used",

	StatusMultipleChoices:   "Multiple Choices",
	StatusMovedPermanently:  "Moved Permanently",
	StatusFound:             "Found",
	StatusSeeOther:          "See Other",
	StatusNotModified:       "Not Modified",
	StatusUseProxy:          "Use Proxy",
	StatusTemporaryRedirect: "Temporary Redirect",
	StatusPermanentRedirect: "Permanent Redirect",

	StatusBadRequest:                   "Bad Request",
	StatusUnauthorized:                 "Unauthorized",
	StatusPaymentRequired:              "Payment Required",
	StatusForbidden:                    "Forbidden",
	StatusNotFound:                     "Not Found",
	StatusMethodNotAllowed:             "Method Not Allowed",
	StatusNotAcceptable:                "Not Acceptable",
	StatusProxyAuthRequired:            "Proxy Authentication Required",
	StatusRequestTimeout:               "Request Timeout",
	StatusConflict:                     "Conflict",
	StatusGone:                         "Gone",
	StatusLengthRequired:               "Length Required",
	StatusPreconditionFailed:           "Precondition Failed",
	StatusRequestEntityTooLarge:        "Request Entity Too Large",
	StatusRequestURITooLong:            "Request URI Too Long",
	StatusUnsupportedMediaType:         "Unsupported Media Type",
	StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	StatusExpectationFailed:            "Expectation Failed",
	StatusTeapot:                       "I'm a teapot",
	StatusMisdirectedRequest:           "Misdirected Request",
	StatusUnprocessableEntity:          "Unprocessable Entity",
	StatusLocked:                       "Locked",
	StatusFailedDependency:             "Failed Dependency",
	StatusTooEarly:                     "Too Early",
	StatusUpgradeRequired:              "Upgrade Required",
	StatusPreconditionRequired:         "Precondition Required",
	StatusTooManyRequests:              "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	StatusInternalServerError:           "Internal Server Error",
	StatusNotImplemented:                "Not Implemented",
	StatusBadGateway:                    "Bad Gateway",
	StatusServiceUnavailable:            "Service Unavailable",
	StatusGatewayTimeout:                "Gateway Timeout",
	StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
	StatusVariantAlsoNegotiates:         "Variant Also Negotiates",
	StatusInsufficientStorage:           "Insufficient Storage",
	StatusLoopDetected:                  "Loop Detected",
	StatusNotExtended:                   "Not Extended",
	StatusNetworkAuthenticationRequired: "Network Authentication Required",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}

const (
	ContentTypeApplicationJavaDashArchive                = "application/java-archive"
	ContentTypeApplicationEdiDashX12                     = "application/EDI-X12"
	ContentTypeApplicationEdifact                        = "application/EDIFACT"
	ContentTypeApplicationJavascript                     = "application/javascript"
	ContentTypeApplicationOctetDashStream                = "application/octet-stream"
	ContentTypeApplicationOgg                            = "application/ogg"
	ContentTypeApplicationPdf                            = "application/pdf"
	ContentTypeApplicationXhtmlPlusXml                   = "application/xhtml+xml"
	ContentTypeApplicationXDashShockwaveDashFlash        = "application/x-shockwave-flash"
	ContentTypeApplicationJson                           = "application/json"
	ContentTypeApplicationLdPlusJson                     = "application/ld+json"
	ContentTypeApplicationXml                            = "application/xml"
	ContentTypeApplicationZip                            = "application/zip"
	ContentTypeApplicationXDashWwwDashFormDashUrlencoded = "application/x-www-form-urlencoded"
	ContentTypeAudioMpeg                                 = "audio/mpeg"
	ContentTypeAudioXDashMsDashWma                       = "audio/x-ms-wma"
	ContentTypeAudioVndDotRnDashRealaudio                = "audio/vnd.rn-realaudio"
	ContentTypeAudioXDashWav                             = "audio/x-wav"
	ContentTypeImageGif                                  = "image/gif"
	ContentTypeImageJpeg                                 = "image/jpeg"
	ContentTypeImagePng                                  = "image/png"
	ContentTypeImageTiff                                 = "image/tiff"
	ContentTypeImageVndDotMicrosoftDotIcon               = "image/vnd.microsoft.icon"
	ContentTypeImageXDashIcon                            = "image/x-icon"
	ContentTypeImageVndDotDjvu                           = "image/vnd.djvu"
	ContentTypeImageSvgPlusXml                           = "image/svg+xml"
	ContentTypeMultipartMixed                            = "multipart/mixed"
	ContentTypeMultipartAlternative                      = "multipart/alternative"
	ContentTypeMultipartRelated                          = "multipart/related"
	ContentTypeMultipartFormDashData                     = "multipart/form-data"
	ContentTypeTextCss                                   = "text/css"
	ContentTypeTextCsv                                   = "text/csv"
	ContentTypeTextHtml                                  = "text/html"
	ContentTypeTextJavascript                            = "text/javascript"
	ContentTypeTextPlain                                 = "text/plain"
	ContentTypeTextXml                                   = "text/xml"
	ContentTypeVideoMpeg                                 = "video/mpeg"
	ContentTypeVideoMp4                                  = "video/mp4"
	ContentTypeVideoQuicktime                            = "video/quicktime"
	ContentTypeVideoXDashMsDashWmv                       = "video/x-ms-wmv"
	ContentTypeVideoXDashMsvideo                         = "video/x-msvideo"
	ContentTypeVideoXDashFlv                             = "video/x-flv"
	ContentTypeVideoWebm                                 = "video/webm"
)

// Request Headers
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields
const (
	// RequestHeaderAIm
	// Acceptable instance-manipulations for the request.
	//
	//	A-IM: feed
	RequestHeaderAIm="A-IM"

	// RequestHeaderAccept
	// Media type(s) that is/are acceptable for the response. See Content negotiation.
	//
	//	Accept: text/html
	RequestHeaderAccept="Accept"

	// RequestHeaderAcceptCharset
	// Character sets that are acceptable.
	//
	//	Accept-Charset: utf-8
	RequestHeaderAcceptCharset="Accept-Charset"

	// RequestHeaderAcceptDatetime
	// Acceptable version in time.
	//
	//	Accept-Datetime: Thu, 31 May 2007 20:35:00 GMT
	RequestHeaderAcceptDatetime="Accept-Datetime"

	// RequestHeaderAcceptEncoding
	// List of acceptable encodings. See HTTP compression.
	//
	//	Accept-Encoding: gzip, deflate
	RequestHeaderAcceptEncoding="Accept-Encoding"

	// RequestHeaderAcceptLanguage
	// List of acceptable human languages for response. See Content negotiation.
	//
	//	Accept-Language: en-US
	RequestHeaderAcceptLanguage="Accept-Language"

	// RequestHeaderAccessControlRequestMethod
	// Initiates a request for cross-origin resource sharing with Origin (below).
	//
	//	Access-Control-Request-Method: GET
	RequestHeaderAccessControlRequestMethod="Access-Control-Request-Method"

	// RequestHeaderAccessControlRequestHeaders
	// Initiates a request for cross-origin resource sharing with Origin (below).
	//Access-Control-Request-Method: GET
	RequestHeaderAccessControlRequestHeaders="Access-Control-Request-Headers"

	// RequestHeaderAuthorization
	// Authentication credentials for HTTP authentication.
	//
	//	Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==
	RequestHeaderAuthorization="Authorization"

	// RequestHeaderCacheControl
	// Used to specify directives that must be obeyed by all caching mechanisms along the request-response chain.
	//
	//	Cache-Control: no-cache
	RequestHeaderCacheControl="Cache-Control"

	// RequestHeaderConnection
	// Control options for the current connection and list of hop-by-hop request fields. Must not be used with HTTP/2.
	//
	//	Connection: keep-alive
	//	Connection: Upgrade
	RequestHeaderConnection="Connection"

	// RequestHeaderContentEncoding
	// The type of encoding used on the data. See HTTP compression.
	//
	//	Content-Encoding: gzip
	RequestHeaderContentEncoding="Content-Encoding"

	// RequestHeaderContentLength
	// The length of the request body in octets (8-bit bytes).
	//
	//	Content-Length: 348
	RequestHeaderContentLength="Content-Length"

	// RequestHeaderContentMd5
	// A Base64-encoded binary MD5 sum of the content of the request body.
	//
	//	Content-MD5: Q2hlY2sgSW50ZWdyaXR5IQ==
	RequestHeaderContentMd5="Content-MD5"

	// RequestHeaderContentType
	// The Media type of the body of the request (used with POST and PUT requests).
	//
	//	Content-Type: application/x-www-form-urlencoded
	RequestHeaderContentType="Content-Type"

	// RequestHeaderCookie
	// An HTTP cookie previously sent by the server with Set-Cookie (below).
	//
	//	Cookie: $Version=1; Skin=new;
	RequestHeaderCookie="Cookie"

	// RequestHeaderDate
	// The date and time at which the message was originated (in "HTTP-date" format as defined by RFC 7231 Date/Time Formats).
	//
	//	Date: Tue, 15 Nov 1994 08:12:31 GMT
	RequestHeaderDate="Date"

	// RequestHeaderExpect
	// Indicates that particular server behaviors are required by the client.
	//
	//	Expect: 100-continue
	RequestHeaderExpect="Expect"

	// RequestHeaderForwarded
	// Disclose original information of a client connecting to a web server through an HTTP proxy.
	//
	//	Forwarded: for=192.0.2.60;proto=http;by=203.0.113.43
	//	Forwarded: for=192.0.2.43, for=198.51.100.17
	RequestHeaderForwarded="Forwarded"

	// RequestHeaderFrom
	// The email address of the user making the request.
	//
	//	From: user@example.com
	RequestHeaderFrom="From"

	// RequestHeaderHost
	// The domain name of the server (for virtual hosting), and the TCP port number on which the server is listening. The port number may be omitted if the port is the standard port for the service requested. Mandatory since HTTP/1.1. If the request is generated directly in HTTP/2, it should not be used.
	//
	//	Host: en.wikipedia.org:8080
	//	Host: en.wikipedia.org
	RequestHeaderHost="Host"

	// RequestHeaderHttp2Settings
	// A request that upgrades from HTTP/1.1 to HTTP/2 MUST include exactly one HTTP2-Setting header field. The HTTP2-Settings header field is a connection-specific header field that includes parameters that govern the HTTP/2 connection, provided in anticipation of the server accepting the request to upgrade.
	//
	//	HTTP2-Settings: token64
	RequestHeaderHttp2Settings="HTTP2-Settings"

	// RequestHeaderIfMatch
	// Only perform the action if the client supplied entity matches the same entity on the server. This is mainly for methods like PUT to only update a resource if it has not been modified since the user last updated it.
	//
	//	If-Match: "737060cd8c284d8af7ad3082f209582d"
	RequestHeaderIfMatch="If-Match"

	// RequestHeaderIfModifiedSince
	// Allows a 304 Not Modified to be returned if content is unchanged.
	//
	//	If-Modified-Since: Sat, 29 Oct 1994 19:43:31 GMT
	RequestHeaderIfModifiedSince="If-Modified-Since"

	// RequestHeaderIfNoneMatch
	// Allows a 304 Not Modified to be returned if content is unchanged, see HTTP ETag.
	//
	//	If-None-Match: "737060cd8c284d8af7ad3082f209582d"
	RequestHeaderIfNoneMatch="If-None-Match"

	// RequestHeaderIfRange
	// If the entity is unchanged, send me the part(s) that I am missing; otherwise, send me the entire new entity.
	//
	//	If-Range: "737060cd8c284d8af7ad3082f209582d"
	RequestHeaderIfRange="If-Range"

	// RequestHeaderIfUnmodifiedSince
	// Only send the response if the entity has not been modified since a specific time.
	//
	//	If-Unmodified-Since: Sat, 29 Oct 1994 19:43:31 GMT
	RequestHeaderIfUnmodifiedSince="If-Unmodified-Since"

	// RequestHeaderMaxForwards
	// Limit the number of times the message can be forwarded through proxies or gateways.
	//
	//	Max-Forwards: 10
	RequestHeaderMaxForwards="Max-Forwards"

	// RequestHeaderOrigin
	// Initiates a request for cross-origin resource sharing (asks server for Access-Control-* response fields).
	//
	//	Origin: http://www.example-social-network.com
	RequestHeaderOrigin="Origin"

	// RequestHeaderPragma
	// Implementation-specific fields that may have various effects anywhere along the request-response chain.
	//
	//	Pragma: no-cache
	RequestHeaderPragma="Pragma"

	// RequestHeaderPrefer
	// Allows client to request that certain behaviors be employed by a server while processing a request.
	//
	//	Prefer: return=representation
	RequestHeaderPrefer="Prefer"

	// RequestHeaderProxyAuthorization
	// Authorization credentials for connecting to a proxy.
	//
	//	Proxy-Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==
	RequestHeaderProxyAuthorization="Proxy-Authorization"

	// RequestHeaderRange
	// Request only part of an entity.  Bytes are numbered from 0.  See Byte serving.
	//
	//	Range: bytes=500-999
	RequestHeaderRange="Range"

	// RequestHeaderReferer
	// This is the address of the previous web page from which a link to the currently requested page was followed. (The word "referrer" has been misspelled in the RFC as well as in most implementations to the point that it has become standard usage and is considered correct terminology)
	//
	//	Referer: http://en.wikipedia.org/wiki/Main_Page
	RequestHeaderReferer ="Referer"

	// RequestHeaderTe
	// The transfer encodings the user agent is willing to accept: the same values as for the response header field Transfer-Encoding can be used, plus the "trailers" value (related to the "chunked" transfer method) to notify the server it expects to receive additional fields in the trailer after the last, zero-sized, chunk. Only trailers is supported in HTTP/2.
	//
	//	TE: trailers, deflate
	RequestHeaderTe="TE"

	// RequestHeaderTrailer
	// The Trailer general field value indicates that the given set of header fields is present in the trailer of a message encoded with chunked transfer coding.
	//
	//	Trailer: Max-Forwards
	RequestHeaderTrailer="Trailer"

	// RequestHeaderTransferEncoding
	// The form of encoding used to safely transfer the entity to the user. Currently defined methods are: chunked, compress, deflate, gzip, identity. Must not be used with HTTP/2.
	//
	//	Transfer-Encoding: chunked
	RequestHeaderTransferEncoding="Transfer-Encoding"

	// RequestHeaderUserAgent
	// The user agent string of the user agent.
	//
	//	User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:12.0) Gecko/20100101 Firefox/12.0
	RequestHeaderUserAgent="User-Agent"

	// RequestHeaderUpgrade
	// Ask the server to upgrade to another protocol. Must not be used in HTTP/2.
	//
	//	Upgrade: h2c, HTTPS/1.3, IRC/6.9, RTA/x11, websocket
	RequestHeaderUpgrade="Upgrade"

	// RequestHeaderVia
	// Informs the server of proxies through which the request was sent.
	//
	//	Via: 1.0 fred, 1.1 example.com (Apache/1.1)
	RequestHeaderVia="Via"

	// RequestHeaderWarning
	// A general warning about possible problems with the entity body.
	//
	//	Warning: 199 Miscellaneous warning
	RequestHeaderWarning="Warning"
)

// Response Headers
const (

	// ResponseHeaderAcceptCh
	// Requests HTTP Client Hints
	//
	//	Accept-CH: UA, Platform
	ResponseHeaderAcceptCh="Accept-CH"

	// ResponseHeaderAccessControlAllowOrigin
	// Specifying which web sites can participate in cross-origin resource sharing
	//
	//	Access-Control-Allow-Origin: *
	ResponseHeaderAccessControlAllowOrigin="Access-Control-Allow-Origin"

	// ResponseHeaderAccessControlAllowCredentials
	// Specifying which web sites can participate in cross-origin resource sharing
	//Access-Control-Allow-Origin: *
	ResponseHeaderAccessControlAllowCredentials="Access-Control-Allow-Credentials"

	// ResponseHeaderAccessControlExposeHeaders
	// Specifying which web sites can participate in cross-origin resource sharing
	//Access-Control-Allow-Origin: *
	ResponseHeaderAccessControlExposeHeaders="Access-Control-Expose-Headers"

	// ResponseHeaderAccessControlMaxAge
	// Specifying which web sites can participate in cross-origin resource sharing
	//Access-Control-Allow-Origin: *
	ResponseHeaderAccessControlMaxAge="Access-Control-Max-Age"

	// ResponseHeaderAccessControlAllowMethods
	// Specifying which web sites can participate in cross-origin resource sharing
	//Access-Control-Allow-Origin: *
	ResponseHeaderAccessControlAllowMethods="Access-Control-Allow-Methods"

	// ResponseHeaderAccessControlAllowHeaders
	// Specifying which web sites can participate in cross-origin resource sharing
	//Access-Control-Allow-Origin: *
	ResponseHeaderAccessControlAllowHeaders="Access-Control-Allow-Headers"

	// ResponseHeaderAcceptPatch
	// Specifies which patch document formats this server supports
	//
	//	Accept-Patch: text/example;charset=utf-8
	ResponseHeaderAcceptPatch="Accept-Patch"

	// ResponseHeaderAcceptRanges
	// What partial content range types this server supports via byte serving
	//
	//	Accept-Ranges: bytes
	ResponseHeaderAcceptRanges="Accept-Ranges"

	// ResponseHeaderAge
	// The age the object has been in a proxy cache in seconds
	//
	//	Age: 12
	ResponseHeaderAge="Age"

	// ResponseHeaderAllow
	// Valid methods for a specified resource. To be used for a 405 Method not allowed
	//
	//	Allow: GET, HEAD
	ResponseHeaderAllow="Allow"

	// ResponseHeaderAltSvc
	// A server uses "Alt-Svc" header (meaning Alternative Services) to indicate that its resources can also be accessed at a different network location (host or port) or using a different protocol When using HTTP/2, servers should instead send an ALTSVC frame.
	//
	//	Alt-Svc: http/1.1="http2.example.com:8001"; ma=7200
	ResponseHeaderAltSvc="Alt-Svc"

	// ResponseHeaderCacheControl
	// Tells all caching mechanisms from server to client whether they may cache this object. It is measured in seconds
	//
	//	Cache-Control: max-age=3600
	ResponseHeaderCacheControl="Cache-Control"

	// ResponseHeaderConnection
	// Control options for the current connection and list of hop-by-hop response fields. Must not be used with HTTP/2.
	//
	//	Connection: close
	ResponseHeaderConnection="Connection"

	// ResponseHeaderContentDisposition
	// An opportunity to raise a "File Download" dialogue box for a known MIME type with binary format or suggest a filename for dynamic content. Quotes are necessary with special characters.
	//
	//	Content-Disposition: attachment; filename="fname.ext"
	ResponseHeaderContentDisposition="Content-Disposition"

	// ResponseHeaderContentEncoding
	// The type of encoding used on the data. See HTTP compression.
	//
	//	Content-Encoding: gzip
	ResponseHeaderContentEncoding="Content-Encoding"

	// ResponseHeaderContentLanguage
	// The natural language or languages of the intended audience for the enclosed content
	//
	//	Content-Language: da
	ResponseHeaderContentLanguage="Content-Language"

	// ResponseHeaderContentLength
	// The length of the response body in octets (8-bit bytes)
	//
	//	Content-Length: 348
	ResponseHeaderContentLength="Content-Length"

	// ResponseHeaderContentLocation
	// An alternate location for the returned data
	//
	//	Content-Location: /index.htm
	ResponseHeaderContentLocation="Content-Location"

	// ResponseHeaderContentMd5
	// A Base64-encoded binary MD5 sum of the content of the response
	//
	//	Content-MD5: Q2hlY2sgSW50ZWdyaXR5IQ==
	ResponseHeaderContentMd5="Content-MD5"

	// ResponseHeaderContentRange
	// Where in a full body message this partial message belongs
	//
	//	Content-Range: bytes 21010-47021/47022
	ResponseHeaderContentRange="Content-Range"

	// ResponseHeaderContentType
	// The MIME type of this content
	//
	//	Content-Type: text/html; charset=utf-8
	ResponseHeaderContentType="Content-Type"

	// ResponseHeaderDate
	// The date and time that the message was sent (in "HTTP-date" format as defined by RFC 7231)
	//
	//	Date: Tue, 15 Nov 1994 08:12:31 GMT
	ResponseHeaderDate="Date"

	// ResponseHeaderDeltaBase
	// Specifies the delta-encoding entity tag of the response.
	//
	//	Delta-Base: "abc"
	ResponseHeaderDeltaBase="Delta-Base"

	// ResponseHeaderEtag
	// An identifier for a specific version of a resource, often a message digest
	//
	//	ETag: "737060cd8c284d8af7ad3082f209582d"
	ResponseHeaderEtag="ETag"

	// ResponseHeaderExpires
	// Gives the date/time after which the response is considered stale (in "HTTP-date" format as defined by RFC 7231)
	//
	//	Expires: Thu, 01 Dec 1994 16:00:00 GMT
	ResponseHeaderExpires="Expires"

	// ResponseHeaderIm
	// Instance-manipulations applied to the response.
	//
	//	IM: feed
	ResponseHeaderIm="IM"

	// ResponseHeaderLastModified
	// The last modified date for the requested object (in "HTTP-date" format as defined by RFC 7231)
	//
	//	Last-Modified: Tue, 15 Nov 1994 12:45:26 GMT
	ResponseHeaderLastModified="Last-Modified"

	// ResponseHeaderLink
	// Used to express a typed relationship with another resource, where the relation type is defined by RFC 5988
	//
	//	Link: </feed>; rel="alternate"
	ResponseHeaderLink="Link"

	// ResponseHeaderLocation
	// Used in redirection, or when a new resource has been created.
	//Example 1:
	//	Location: http://www.w3.org/pub/WWW/People.html Example 2:
	//	Location: /pub/WWW/People.html
	ResponseHeaderLocation="Location"

	// ResponseHeaderP3p
	// This field is supposed to set P3P policy, in the form of P3P:CP="your_compact_policy". However, P3P did not take off, most browsers have never fully implemented it, a lot of websites set this field with fake policy text, that was enough to fool browsers the existence of P3P policy and grant permissions for third party cookies.
	//
	//	P3P: CP="This is not a
	//	P3P policy! See https://en.wikipedia.org/wiki/Special:CentralAutoLogin/
	//	P3P for more info."
	ResponseHeaderP3p="P3P"

	// ResponseHeaderPragma
	// Implementation-specific fields that may have various effects anywhere along the request-response chain.
	//
	//	Pragma: no-cache
	ResponseHeaderPragma="Pragma"

	// ResponseHeaderPreferenceApplied
	// Indicates which Prefer tokens were honored by the server and applied to the processing of the request.
	//
	//	Preference-Applied: return=representation
	ResponseHeaderPreferenceApplied="Preference-Applied"

	// ResponseHeaderProxyAuthenticate
	// Request authentication to access the proxy.
	//
	//	Proxy-Authenticate: Basic
	ResponseHeaderProxyAuthenticate="Proxy-Authenticate"

	// ResponseHeaderPublicKeyPins
	// HTTP Public Key Pinning, announces hash of website's authentic TLS certificate
	//
	//	Public-Key-Pins: max-age=2592000; pin-sha256="E9CZ9INDbd+2eRQozYqqbQ2yXLVKB9+xcprMF+44U1g=";
	ResponseHeaderPublicKeyPins="Public-Key-Pins"

	// ResponseHeaderRetryAfter
	// If an entity is temporarily unavailable, this instructs the client to try again later. Value could be a specified period of time (in seconds) or a HTTP-date.
	//Example 1:
	//	Retry-After: 120 Example 2:
	//	Retry-After: Fri, 07 Nov 2014 23:59:59 GMT
	ResponseHeaderRetryAfter="Retry-After"

	// ResponseHeaderServer
	// A name for the server
	//
	//	Server: Apache/2.4.1 (Unix)
	ResponseHeaderServer="Server"

	// ResponseHeaderSetCookie
	// An HTTP cookie
	//
	//	Set-Cookie: UserID=JohnDoe; Max-Age=3600; Version=1
	ResponseHeaderSetCookie="Set-Cookie"

	// ResponseHeaderStrictTransportSecurity
	// A HSTS Policy informing the HTTP client how long to cache the HTTPS only policy and whether this applies to subdomains.
	//
	//	Strict-Transport-Security: max-age=16070400; includeSubDomains
	ResponseHeaderStrictTransportSecurity="Strict-Transport-Security"

	// ResponseHeaderTrailer
	// The Trailer general field value indicates that the given set of header fields is present in the trailer of a message encoded with chunked transfer coding.
	//
	//	Trailer: Max-Forwards
	ResponseHeaderTrailer="Trailer"

	// ResponseHeaderTransferEncoding
	// The form of encoding used to safely transfer the entity to the user. Currently defined methods are: chunked, compress, deflate, gzip, identity. Must not be used with HTTP/2.
	//
	//	Transfer-Encoding: chunked
	ResponseHeaderTransferEncoding="Transfer-Encoding"

	// ResponseHeaderTk
	// Tracking Status header, value suggested to be sent in response to a DNT(do-not-track), possible values: "!" — under construction "?" — dynamic "G" — gateway to multiple parties "N" — not tracking "T" — tracking "C" — tracking with consent "P" — tracking only if consented "D" — disregarding DNT "U" — updated
	//
	//	Tk: ?
	ResponseHeaderTk="Tk"

	// ResponseHeaderUpgrade
	// Ask the client to upgrade to another protocol. Must not be used in HTTP/2
	//
	//	Upgrade: h2c, HTTPS/1.3, IRC/6.9, RTA/x11, websocket
	ResponseHeaderUpgrade="Upgrade"

	// ResponseHeaderVary
	// Tells downstream proxies how to match future request headers to decide whether the cached response can be used rather than requesting a fresh one from the origin server.
	//Example 1:
	//	Vary: * Example 2:
	//	Vary: Accept-Language
	ResponseHeaderVary="Vary"

	// ResponseHeaderVia
	// Informs the client of proxies through which the response was sent.
	//
	//	Via: 1.0 fred, 1.1 example.com (Apache/1.1)
	ResponseHeaderVia="Via"

	// ResponseHeaderWarning
	// A general warning about possible problems with the entity body.
	//
	//	Warning: 199 Miscellaneous warning
	ResponseHeaderWarning="Warning"

	// ResponseHeaderWwwAuthenticate
	// Indicates the authentication scheme that should be used to access the requested entity.
	//
	//	WWW-Authenticate: Basic
	ResponseHeaderWwwAuthenticate="WWW-Authenticate"

	// ResponseHeaderXFrameOptions
	// Clickjacking protection: deny - no rendering within a frame, sameorigin - no rendering if origin mismatch, allow-from - allow from specified location, allowall - non-standard, allow from any location
	//
	//	X-Frame-Options: deny
	ResponseHeaderXFrameOptions="X-Frame-Options"
)
