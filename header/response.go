package header

// Response Headers
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields
const (

	// ResponseAcceptCh
	// Requests HTTP Client Hints
	//
	//	Accept-CH: UA, Platform
	ResponseAcceptCh = "Accept-CH"

	// ResponseAccessControlAllowOrigin
	// Specifying which web sites can participate in cross-origin resource sharing
	//
	//	Access-Control-Allow-Origin: *
	ResponseAccessControlAllowOrigin = "Access-Control-Allow-Origin"

	// ResponseAccessControlAllowCredentials
	// Specifying which web sites can participate in cross-origin resource sharing
	//Access-Control-Allow-Origin: *
	ResponseAccessControlAllowCredentials = "Access-Control-Allow-Credentials"

	// ResponseAccessControlExposeHeaders
	// Specifying which web sites can participate in cross-origin resource sharing
	//Access-Control-Allow-Origin: *
	ResponseAccessControlExposeHeaders = "Access-Control-Expose-Headers"

	// ResponseAccessControlMaxAge
	// Specifying which web sites can participate in cross-origin resource sharing
	//Access-Control-Allow-Origin: *
	ResponseAccessControlMaxAge = "Access-Control-Max-Age"

	// ResponseAccessControlAllowMethods
	// Specifying which web sites can participate in cross-origin resource sharing
	//Access-Control-Allow-Origin: *
	ResponseAccessControlAllowMethods = "Access-Control-Allow-Methods"

	// ResponseAccessControlAllowHeaders
	// Specifying which web sites can participate in cross-origin resource sharing
	//Access-Control-Allow-Origin: *
	ResponseAccessControlAllowHeaders = "Access-Control-Allow-Headers"

	// ResponseAcceptPatch
	// Specifies which patch document formats this server supports
	//
	//	Accept-Patch: text/example;charset=utf-8
	ResponseAcceptPatch = "Accept-Patch"

	// ResponseAcceptRanges
	// What partial content range types this server supports via byte serving
	//
	//	Accept-Ranges: bytes
	ResponseAcceptRanges = "Accept-Ranges"

	// ResponseAge
	// The age the object has been in a proxy cache in seconds
	//
	//	Age: 12
	ResponseAge = "Age"

	// ResponseAllow
	// Valid methods for a specified resource. To be used for a 405 Method not allowed
	//
	//	Allow: GET, HEAD
	ResponseAllow = "Allow"

	// ResponseAltSvc
	// A server uses "Alt-Svc" header (meaning Alternative Services) to indicate that its resources can also be accessed at a different network location (host or port) or using a different protocol When using HTTP/2, servers should instead send an ALTSVC frame.
	//
	//	Alt-Svc: http/1.1="http2.example.com:8001"; ma=7200
	ResponseAltSvc = "Alt-Svc"

	// ResponseContentDisposition
	// An opportunity to raise a "File Download" dialogue box for a known MIME type with binary bind or suggest a filename for dynamic content. Quotes are necessary with special characters.
	//
	//	Content-Disposition: attachment; filename="fname.ext"
	ResponseContentDisposition = "Content-Disposition"

	// ResponseContentLanguage
	// The natural language or languages of the intended audience for the enclosed content
	//
	//	Content-Language: da
	ResponseContentLanguage = "Content-Language"

	// ResponseContentLocation
	// An alternate location for the returned data
	//
	//	Content-Location: /index.htm
	ResponseContentLocation = "Content-Location"

	// ResponseContentRange
	// Where in a full body message this partial message belongs
	//
	//	Content-Range: bytes 21010-47021/47022
	ResponseContentRange = "Content-Range"

	// ResponseDeltaBase
	// Specifies the delta-encoding entity tag of the response.
	//
	//	Delta-Base: "abc"
	ResponseDeltaBase = "Delta-Base"

	// ResponseEtag
	// An identifier for a specific version of a resource, often a message digest
	//
	//	ETag: "737060cd8c284d8af7ad3082f209582d"
	ResponseEtag = "ETag"

	// ResponseExpires
	// Gives the date/time after which the response is considered stale (in "HTTP-date" bind as defined by RFC 7231)
	//
	//	Expires: Thu, 01 Dec 1994 16:00:00 GMT
	ResponseExpires = "Expires"

	// ResponseIm
	// Instance-manipulations applied to the response.
	//
	//	IM: feed
	ResponseIm = "IM"

	// ResponseLastModified
	// The last modified date for the requested object (in "HTTP-date" bind as defined by RFC 7231)
	//
	//	Last-Modified: Tue, 15 Nov 1994 12:45:26 GMT
	ResponseLastModified = "Last-Modified"

	// ResponseLink
	// Used to express a typed relationship with another resource, where the relation type is defined by RFC 5988
	//
	//	Link: </feed>; rel="alternate"
	ResponseLink = "Link"

	// ResponseLocation
	// Used in redirection, or when a new resource has been created.
	// Example 1:
	//	Location: http://www.w3.org/pub/WWW/People.html Example 2:
	//	Location: /pub/WWW/People.html
	ResponseLocation = "Location"

	// ResponseP3p
	// This field is supposed to set P3P policy, in the form of P3P:CP="your_compact_policy". However, P3P did not take off, most browsers have never fully implemented it, a lot of websites set this field with fake policy text, that was enough to fool browsers the existence of P3P policy and grant permissions for third party cookies.
	//
	//	P3P: CP="This is not a
	//	P3P policy! See https://en.wikipedia.org/wiki/Special:CentralAutoLogin/
	//	P3P for more info."
	ResponseP3p = "P3P"

	// ResponsePreferenceApplied
	// Indicates which Prefer tokens were honored by the server and applied to the processing of the request.
	//
	//	Preference-Applied: return=representation
	ResponsePreferenceApplied = "Preference-Applied"

	// ResponseProxyAuthenticate
	// Request authentication to access the proxy.
	//
	//	Proxy-Authenticate: Basic
	ResponseProxyAuthenticate = "Proxy-Authenticate"

	// ResponsePublicKeyPins
	// HTTP Public Key Pinning, announces hash of website's authentic TLS certificate
	//
	//	Public-Key-Pins: max-age=2592000; pin-sha256="E9CZ9INDbd+2eRQozYqqbQ2yXLVKB9+xcprMF+44U1g=";
	ResponsePublicKeyPins = "Public-Key-Pins"

	// ResponseRetryAfter
	// If an entity is temporarily unavailable, this instructs the client to try again later. Value could be a specified period of time (in seconds) or a HTTP-date.
	// Example 1:
	//	Retry-After: 120 Example 2:
	//	Retry-After: Fri, 07 Nov 2014 23:59:59 GMT
	ResponseRetryAfter = "Retry-After"

	// ResponseServer
	// A name for the server
	//
	//	Server: Apache/2.4.1 (Unix)
	ResponseServer = "Server"

	// ResponseSetCookie
	// An HTTP cookie
	//
	//	Set-Cookie: UserID=JohnDoe; Max-Age=3600; Version=1
	ResponseSetCookie = "Set-Cookie"

	// ResponseStrictTransportSecurity
	// A HSTS Policy informing the HTTP client how long to cache the HTTPS only policy and whether this applies to subdomains.
	//
	//	Strict-Transport-Security: max-age=16070400; includeSubDomains
	ResponseStrictTransportSecurity = "Strict-Transport-Security"

	// ResponseTk
	// Tracking Status header, value suggested to be sent in response to a DNT(do-not-track), possible values: "!" — under construction "?" — dynamic "G" — gateway to multiple parties "N" — not tracking "T" — tracking "C" — tracking with consent "P" — tracking only if consented "D" — disregarding DNT "U" — updated
	//
	//	Tk: ?
	ResponseTk = "Tk"

	// ResponseVary
	// Tells downstream proxies how to match future request headers to decide whether the cached response can be used rather than requesting a fresh one from the origin server.
	// Example 1:
	//	Vary: * Example 2:
	//	Vary: Accept-Language
	ResponseVary = "Vary"

	// ResponseWwwAuthenticate
	// Indicates the authentication scheme that should be used to access the requested entity.
	//
	//	WWW-Authenticate: Basic
	ResponseWwwAuthenticate = "WWW-Authenticate"

	// ResponseXFrameOptions
	// Clickjacking protection: deny - no rendering within a frame, sameorigin - no rendering if origin mismatch, allow-from - allow from specified location, allowall - non-standard, allow from any location
	//
	//	X-Frame-Options: deny
	ResponseXFrameOptions = "X-Frame-Options"
)
