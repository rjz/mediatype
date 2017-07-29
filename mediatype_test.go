package mediatype_test

import (
	"github.com/rjz/mediatype"
	"net/http"
	"strings"
	"testing"
)

func TestPrefersInvalid(t *testing.T) {
	tests := []string{
		"text/foo/bar",
		"text/plain;q=asdklj",
		"text/plain;q=14",
	}

	for _, test := range tests {
		mediaRange := test
		r, _ := http.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("Accept", mediaRange)

		if _, err := mediatype.Prefers(r); err == nil {
			t.Errorf("Expected error from invalid Accept header '%s'", mediaRange)
		}
	}
}

func TestPrefers(t *testing.T) {
	tests := [...][2]string{
		{"text/json;q=0.8, text/plain;q=0.5", "text/json,text/plain"},
		{"text/json;q=0.8, text/plain", "text/plain,text/json"},
		{"text/json, text/plain", "text/plain,text/json"},
	}

	for _, test := range tests {
		mediaRange := test[0]
		expected := test[1]

		r, _ := http.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("Accept", mediaRange)

		prefs, _ := mediatype.Prefers(r)
		if strings.Join(prefs, ",") != expected {
			t.Errorf("Expected '%s' to equal '%s'", strings.Join(prefs, ","), expected)
		}
	}
}

func TestAcceptsInvalid(t *testing.T) {
	mediaRange := "text/foo/bar"
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	r.Header.Set("Accept", mediaRange)

	if _, err := mediatype.Accepts(r, "application/json"); err == nil {
		t.Errorf("Expected error from invalid Accept header '%s'", mediaRange)
	}
}

func TestAccepts(t *testing.T) {
	tests := [...][2]string{
		{"text/json, text/plain", "text/json"},
		{"text/json, application/octet-stream, text/plain;q=0.8", "text/json"},
		{"*/*", "text/json"},
		{"text/*", "text/json"},
	}

	for _, test := range tests {
		mediaRange := test[0]
		mimetype := test[1]
		r, _ := http.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("Accept", mediaRange)

		if ok, _ := mediatype.Accepts(r, mimetype); !ok {
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
		r, _ := http.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("Accept", mediaRange)

		if ok, _ := mediatype.Accepts(r, mimetype); ok {
			t.Errorf("Expected '%s' to reject '%s'", mediaRange, mimetype)
		}
	}
}

func TestHasContentType(t *testing.T) {
	tests := [...][2]string{
		{"application/json", "application/json"},
		{"application/octet-stream", mediatype.DefaultMimeType},
		{"", mediatype.DefaultMimeType},
	}

	for _, test := range tests {
		contentType := test[0]
		mimetype := test[1]
		r, _ := http.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("Content-type", contentType)

		if !mediatype.HasContentType(r, mimetype) {
			t.Errorf("Expected '%s' to have content-type '%s'", contentType, mimetype)
		}
	}
}
