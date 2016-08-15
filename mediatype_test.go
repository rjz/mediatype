package mediatype_test

import (
	"github.com/rjz/mediatype"
	"net/http"
	"testing"
)

func req(mediaRange string) *http.Request {
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	r.Header.Set("Accept", mediaRange)
	r.Header.Set("Content-type", mediaRange)
	return r
}

func TestAccepts(t *testing.T) {
	tests := [...][2]string{
		{"text/json, text/plain", "text/json"},
		{"*/*", "text/json"},
		{"text/*", "text/json"},
	}

	for _, test := range tests {
		mediaRange := test[0]
		mimetype := test[1]

		if !mediatype.Accepts(req(mediaRange), mimetype) {
			t.Errorf("Expected '%s' to accept '%s'", mediaRange, mimetype)
		}
	}
}

func TestNotAccepts(t *testing.T) {
	tests := [...][2]string{
		{"text/plain", "text/json"},
		{"application/*", "text/json"},
	}

	for _, test := range tests {
		mediaRange := test[0]
		mimetype := test[1]
		if mediatype.Accepts(req(mediaRange), mimetype) {
			t.Errorf("Expected '%s' to reject '%s'", mediaRange, mimetype)
		}
	}
}

func TestContentType(t *testing.T) {
	tests := [...][2]string{
		{"text/json, text/plain", "text/json"},
		{"", "application/octet-stream"},
	}

	for _, test := range tests {
		mediaRange := test[0]
		mimetype := test[1]
		if !mediatype.HasContentType(req(mediaRange), mimetype) {
			t.Errorf("Expected '%s' to have content-type '%s'", mediaRange, mimetype)
		}
	}
}
