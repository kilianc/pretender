package pretender

import (
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
	expected := []string{"hello", "world", "\n"}

	mfs := fstest.MapFS{
		"some/path/responses.txt": {
			Data: []byte(strings.Join(expected, "\n")),
			Mode: 0644,
		},
	}

	hh := HttpHandler{fs: mfs}
	err := hh.LoadResponsesFile("some/path/responses.txt")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	if reflect.DeepEqual(hh.responses, expected) {
		t.Errorf("got %v, expect %v", hh.responses, expected)
	}
}
