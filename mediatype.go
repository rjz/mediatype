package mediatype

import (
	"mime"
	"net/http"
	"strings"
)

func hasMediaRange(mediaRange, mimetype string) bool {
	for _, v := range strings.Split(mediaRange, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}

		if t == mimetype || t == "*/*" || (strings.HasSuffix(t, "/*") && strings.HasPrefix(mimetype, t[0:len(t)-2])) {
			return true
		}
	}
	return false
}

// Determine whether the server can respond with a mime-type allowed by the
// request's `accept` header
func Accepts(r *http.Request, mimetype string) bool {
	return hasMediaRange(r.Header.Get("Accept"), mimetype)
}

// Determine whether the request `content-type` includes a server-acceptable
// mime-type
//
// Failure should yield an HTTP 415 (`http.StatusUnsupportedMediaType`)
func HasContentType(r *http.Request, mimetype string) bool {
	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		return mimetype == "application/octet-stream"
	}
	return hasMediaRange(contentType, mimetype)
}
