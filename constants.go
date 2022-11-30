// All rights belong to the owners of the http package

package there

import (
	"strings"
)

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
	AllMethodsJoined = strings.Join(AllMethods, ",")
)

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
	ContentTypeApplicationMsgpack                        = "application/x-msgpack"
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

var fileContentTypes = map[string]string{
	"avi":   ContentTypeVideoXDashMsvideo,
	"bin":   ContentTypeApplicationOctetDashStream,
	"css":   ContentTypeTextCss,
	"csv":   ContentTypeTextCsv,
	"gif":   ContentTypeImageGif,
	"html":  ContentTypeTextHtml,
	"htm":   ContentTypeTextHtml,
	"ico":   ContentTypeImageXDashIcon,
	"jar":   ContentTypeApplicationJavaDashArchive,
	"jpeg":  ContentTypeImageJpeg,
	"jpg":   ContentTypeImageJpeg,
	"js":    ContentTypeTextJavascript,
	"json":  ContentTypeApplicationJson,
	"mjs":   ContentTypeTextJavascript,
	"mov":   ContentTypeVideoQuicktime,
	"mp3":   ContentTypeAudioMpeg,
	"mp4":   ContentTypeVideoMp4,
	"mpeg":  ContentTypeVideoMpeg,
	"ogx":   ContentTypeApplicationOgg,
	"png":   ContentTypeImagePng,
	"pdf":   ContentTypeApplicationPdf,
	"qt":    ContentTypeVideoQuicktime,
	"tif":   ContentTypeImageTiff,
	"tiff":  ContentTypeImageTiff,
	"txt":   ContentTypeTextPlain,
	"wav":   ContentTypeAudioXDashWav,
	"webm":  ContentTypeVideoWebm,
	"xhtml": ContentTypeApplicationXhtmlPlusXml,
	"xml":   ContentTypeTextXml,
	"zip":   ContentTypeApplicationZip,
}

// ContentType returns a content type header based on a given file-serving extension.
// The second returned var indicates, whether that content type was found in the list or not.
func ContentType(extension string) string {
	return fileContentTypes[extension]
}
