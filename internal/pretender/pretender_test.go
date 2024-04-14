package pretender

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

func Test_HandleFunc(t *testing.T) {
	tests := []struct {
		expect     string
		statusCode int
		err        error
	}{
		{"hello\n", http.StatusOK, nil},
		{"world\n", http.StatusOK, nil},
		{"no responses left\n", http.StatusInternalServerError, nil},
	}

	hh := HttpHandler{
		responses: []string{"hello", "world"},
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		hh.HandleFunc(w, r)

		if w.Body.String() != tt.expect {
			t.Errorf("got %q, expect %q", w.Body, tt.expect)
		}

		if w.Result().StatusCode != tt.statusCode {
			t.Errorf("got %d, expect %d", w.Result().StatusCode, tt.statusCode)
		}
	}

}

func Test_loadResponsesFile(t *testing.T) {
	expected := []string{"hello", "world", ""}

	mfs := fstest.MapFS{
		"some/path/responses.txt": {
			Data: []byte(strings.Join(expected, "\n")),
		},
	}

	tests := []struct {
		path      string
		errPrefix string
	}{
		{"some/path/responses.txt", ""},
		{"some/path/not-exists.txt", "failed to read responses file"},
	}

	for _, tt := range tests {
		hh := HttpHandler{
			logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
			fs:     mfs,
		}

		err := hh.LoadResponsesFile(tt.path)

		// check if error is nil when expected prefix is empty
		if tt.errPrefix == "" && err != nil {
			t.Errorf("got \"%v\", expect nil", err)
		}

		// check if error message starts with expected prefix
		if !strings.HasPrefix(fmt.Errorf("%w", err).Error(), tt.errPrefix) {
			t.Errorf("got \"%v\", expect \"%v*\"", err, tt.errPrefix)
		}

		// check if responses in the file are equal to expected
		if err == nil && !reflect.DeepEqual(hh.responses, expected) {
			t.Errorf("got \"%v\", expect \"%v\"", hh.responses, expected)
		}
	}
}
