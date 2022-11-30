package header

// Request Headers
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields
const (
	// RequestAIm
	// Acceptable instance-manipulations for the request.
	//
	//	A-IM: feed
	RequestAIm = "A-IM"

	// RequestAccept
	// Media type(s) that is/are acceptable for the response. See Content negotiation.
	//
	//	Accept: text/html
	RequestAccept = "Accept"

	// RequestAcceptCharset
	// Character sets that are acceptable.
	//
	//	Accept-Charset: utf-8
	RequestAcceptCharset = "Accept-Charset"

	// RequestAcceptDatetime
	// Acceptable version in time.
	//
	//	Accept-Datetime: Thu, 31 May 2007 20:35:00 GMT
	RequestAcceptDatetime = "Accept-Datetime"

	// RequestAcceptEncoding
	// List of acceptable encodings. See HTTP compression.
	//
	//	Accept-Encoding: gzip, deflate
	RequestAcceptEncoding = "Accept-Encoding"

	// RequestAcceptLanguage
	// List of acceptable human languages for response. See Content negotiation.
	//
	//	Accept-Language: en-US
	RequestAcceptLanguage = "Accept-Language"

	// RequestAccessControlRequestMethod
	// Initiates a request for cross-origin resource sharing with Origin (below).
	//
	//	Access-Control-Request-Method: GET
	RequestAccessControlRequestMethod = "Access-Control-Request-Method"

	// RequestAccessControlRequests
	// Initiates a request for cross-origin resource sharing with Origin (below).
	//Access-Control-Request-Method: GET
	RequestAccessControlRequests = "Access-Control-Request-Headers"

	// RequestAuthorization
	// Authentication credentials for HTTP authentication.
	//
	//	Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==
	RequestAuthorization = "Authorization"

	// RequestCookie
	// An HTTP cookie previously sent by the server with Set-Cookie (below).
	//
	//	Cookie: $Version=1; Skin=new;
	RequestCookie = "Cookie"

	// RequestExpect
	// Indicates that particular server behaviors are required by the client.
	//
	//	Expect: 100-continue
	RequestExpect = "Expect"

	// RequestForwarded
	// Disclose original information of a client connecting to a web server through an HTTP proxy.
	//
	//	Forwarded: for=192.0.2.60;proto=http;by=203.0.113.43
	//	Forwarded: for=192.0.2.43, for=198.51.100.17
	RequestForwarded = "Forwarded"

	// RequestFrom
	// The email address of the user making the request.
	//
	//	From: user@example.com
	RequestFrom = "From"

	// RequestHost
	// The domain name of the server (for virtual hosting), and the TCP port number on which the server is listening. The port number may be omitted if the port is the standard port for the service requested. Mandatory since HTTP/1.1. If the request is generated directly in HTTP/2, it should not be used.
	//
	//	Host: en.wikipedia.org:8080
	//	Host: en.wikipedia.org
	RequestHost = "Host"

	// RequestHttp2Settings
	// A request that upgrades from HTTP/1.1 to HTTP/2 MUST include exactly one HTTP2-Setting header field. The HTTP2-Settings header field is a connection-specific header field that includes parameters that govern the HTTP/2 connection, provided in anticipation of the server accepting the request to upgrade.
	//
	//	HTTP2-Settings: token64
	RequestHttp2Settings = "HTTP2-Settings"

	// RequestIfMatch
	// Only perform the action if the client supplied entity matches the same entity on the server. This is mainly for methods like PUT to only update a resource if it has not been modified since the user last updated it.
	//
	//	If-Match: "737060cd8c284d8af7ad3082f209582d"
	RequestIfMatch = "If-Match"

	// RequestIfModifiedSince
	// Allows a 304 Not Modified to be returned if content is unchanged.
	//
	//	If-Modified-Since: Sat, 29 Oct 1994 19:43:31 GMT
	RequestIfModifiedSince = "If-Modified-Since"

	// RequestIfNoneMatch
	// Allows a 304 Not Modified to be returned if content is unchanged, see HTTP ETag.
	//
	//	If-None-Match: "737060cd8c284d8af7ad3082f209582d"
	RequestIfNoneMatch = "If-None-Match"

	// RequestIfRange
	// If the entity is unchanged, send me the part(s) that I am missing; otherwise, send me the entire new entity.
	//
	//	If-Range: "737060cd8c284d8af7ad3082f209582d"
	RequestIfRange = "If-Range"

	// RequestIfUnmodifiedSince
	// Only send the response if the entity has not been modified since a specific time.
	//
	//	If-Unmodified-Since: Sat, 29 Oct 1994 19:43:31 GMT
	RequestIfUnmodifiedSince = "If-Unmodified-Since"

	// RequestMaxForwards
	// Limit the number of times the message can be forwarded through proxies or gateways.
	//
	//	Max-Forwards: 10
	RequestMaxForwards = "Max-Forwards"

	// RequestOrigin
	// Initiates a request for cross-origin resource sharing (asks server for Access-Control-* response fields).
	//
	//	Origin: http://www.example-social-network.com
	RequestOrigin = "Origin"

	// RequestPrefer
	// Allows client to request that certain behaviors be employed by a server while processing a request.
	//
	//	Prefer: return=representation
	RequestPrefer = "Prefer"

	// RequestProxyAuthorization
	// Authorization credentials for connecting to a proxy.
	//
	//	Proxy-Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==
	RequestProxyAuthorization = "Proxy-Authorization"

	// RequestRange
	// Request only part of an entity.  ToBytes are numbered from 0.  See Byte serving.
	//
	//	Range: bytes=500-999
	RequestRange = "Range"

	// RequestReferer
	// This is the address of the previous web page from which a link to the currently requested page was followed. (The word "referrer" has been misspelled in the RFC as well as in most implementations to the point that it has become standard usage and is considered correct terminology)
	//
	//	Referer: http://en.wikipedia.org/wiki/Main_Page
	RequestReferer = "Referer"

	// RequestTe
	// The transfer encodings the user agent is willing to accept: the same values as for the response header field Transfer-Encoding can be used, plus the "trailers" value (related to the "chunked" transfer method) to notify the server it expects to receive additional fields in the trailer after the last, zero-sized, chunk. Only trailers is supported in HTTP/2.
	//
	//	TE: trailers, deflate
	RequestTe = "TE"

	// RequestUserAgent
	// The user agent string of the user agent.
	//
	//	User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:12.0) Gecko/20100101 Firefox/12.0
	RequestUserAgent = "User-Agent"
)
