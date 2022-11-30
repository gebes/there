package header

// Request and Response Headers
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields
const (
	// CacheControl
	// Used to specify directives that must be obeyed by all caching mechanisms along the request-response chain.
	//
	//	Cache-Control: no-cache
	CacheControl = "Cache-Control"

	// Connection
	// Control options for the current connection and list of hop-by-hop request fields. Must not be used with HTTP/2.
	//
	//	Connection: keep-alive
	//	Connection: Upgrade
	Connection = "Connection"

	// ContentEncoding
	// The type of encoding used on the data. See HTTP compression.
	//
	//	Content-Encoding: gzip
	ContentEncoding = "Content-Encoding"

	// ContentLength
	// The length of the request body in octets (8-bit bytes).
	//
	//	Content-Length: 348
	ContentLength = "Content-Length"

	// ContentMd5
	// A Base64-encoded binary MD5 sum of the content of the request body.
	//
	//	Content-MD5: Q2hlY2sgSW50ZWdyaXR5IQ==
	ContentMd5 = "Content-MD5"

	// ContentType
	// The Media type of the body of the request (used with POST and PUT s).
	//
	//	Content-Type: application/x-www-form-urlencoded
	ContentType = "Content-Type"

	// Date
	// The date and time at which the message was originated (in "HTTP-date" bind as defined by RFC 7231 Date/Time Formats).
	//
	//	Date: Tue, 15 Nov 1994 08:12:31 GMT
	Date = "Date"

	// Pragma
	// Implementation-specific fields that may have various effects anywhere along the request-response chain.
	//
	//	Pragma: no-cache
	Pragma = "Pragma"

	// Trailer
	// The Trailer general field value indicates that the given set of header fields is present in the trailer of a message encoded with chunked transfer coding.
	//
	//	Trailer: Max-Forwards
	Trailer = "Trailer"

	// TransferEncoding
	// The form of encoding used to safely transfer the entity to the user. Currently defined methods are: chunked, compress, deflate, gzip, identity. Must not be used with HTTP/2.
	//
	//	Transfer-Encoding: chunked
	TransferEncoding = "Transfer-Encoding"

	// Upgrade
	// Ask the server to upgrade to another protocol. Must not be used in HTTP/2.
	//
	//	Upgrade: h2c, HTTPS/1.3, IRC/6.9, RTA/x11, websocket
	Upgrade = "Upgrade"

	// Via
	// Informs the server of proxies through which the request was sent.
	//
	//	Via: 1.0 fred, 1.1 example.com (Apache/1.1)
	Via = "Via"

	// Warning
	// A general warning about possible problems with the entity body.
	//
	//	Warning: 199 Miscellaneous warning
	Warning = "Warning"
)
