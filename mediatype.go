package mediatype

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var ErrInvalid = errors.New("invalid media-range")

// DefaultMimeType is applied to requests without a Content-type header
const DefaultMimeType = "application/octet-stream"

// IsAccepted validates the mimetype against the specified pattern
func IsAccepted(mimetype, pattern string) bool {
	if pattern == "*/*" {
		return true
	}

	delimIndex := strings.Index(mimetype, "/")
	if strings.EqualFold(pattern[0:delimIndex], mimetype[0:delimIndex]) {
		if pattern[delimIndex:] == "/*" {
			return true
		}
		return strings.EqualFold(pattern[delimIndex:], mimetype[delimIndex:])
	}
	return false
}

// Prefers returns a list of accepted mime-types ordered by quality
func Prefers(request *http.Request) ([]string, error) {
	mediaRange := request.Header.Get("Accept")
	acceptedTypes := make(map[string]string)

	var qualities []string
	for i, v := range strings.Split(mediaRange, ",") {
		t, details, err := mime.ParseMediaType(v)
		if err != nil {
			return nil, ErrInvalid
		}

		quality := float64(1)
		if q, ok := details["q"]; ok {
			qf, err := strconv.ParseFloat(q, 32)
			if err != nil || qf < 0 || qf > 1 {
				return nil, ErrInvalid
			}

			quality = qf
		}
		qualStr := fmt.Sprintf("%0.2f#%d", quality, i)
		acceptedTypes[qualStr] = t
		qualities = append(qualities, qualStr)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(qualities)))
	fmt.Println(qualities)

	var prefs []string
	for _, k := range qualities {
		prefs = append(prefs, acceptedTypes[k])
	}

	return prefs, nil
}

// Accepts verifies that a request will accept the specified mimetype
func Accepts(request *http.Request, mimetype string) (bool, error) {
	mediaRange := request.Header.Get("Accept")
	for _, v := range strings.Split(mediaRange, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			return false, ErrInvalid
		}

		if IsAccepted(mimetype, t) {
			return true, nil
		}
	}
	return false, nil
}

// HasContentType verifies that the request matches the specified mimetype
func HasContentType(request *http.Request, mimetype string) bool {
	contentType := request.Header.Get("Content-type")
	if contentType == "" {
		return mimetype == DefaultMimeType
	}
	return strings.EqualFold(contentType, mimetype)
}
