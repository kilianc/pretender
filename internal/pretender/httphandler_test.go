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
		statusCode  int
		body        string
		contentType string
	}{
		{
			statusCode:  http.StatusOK,
			body:        "hello\n",
			contentType: "my/type1",
		},
		{
			statusCode:  http.StatusOK,
			body:        "world\n",
			contentType: "my/type2",
		},
		{
			statusCode:  http.StatusInternalServerError,
			body:        "no responses left\n",
			contentType: "text/plain; charset=utf-8",
		},
	}

	hh := HTTPHandler{
		responses: []response{
			{
				StatusCode: http.StatusOK,
				Body:       []byte("hello"),
				Headers:    map[string]string{"content-type": "my/type1"},
			},
			{
				StatusCode: http.StatusOK,
				Body:       []byte("world"),
				Headers:    map[string]string{"content-type": "my/type2"},
			},
		},
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v %q", tt.statusCode, tt.body), func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			hh.HandleFunc(w, r)

			if w.Body.String() != tt.body {
				t.Errorf("got %q, expect %q", w.Body, tt.body)
			}

			if w.Result().StatusCode != tt.statusCode {
				t.Errorf("got %d, expect %d", w.Result().StatusCode, tt.statusCode)
			}

			if w.Result().Header.Get("content-type") != tt.contentType {
				t.Errorf("got %q, expect %q", w.Result().Header.Get("content-type"), tt.contentType)
			}
		})
	}
}

func Test_LoadResponsesFile(t *testing.T) {
	mfs := fstest.MapFS{
		"some/path/valid.json": {Data: []byte(`[
			{"body":"hello","headers":{"content-type":"text/plain"},"delay_ms":1000},
			{"status_code":404,"body":"world","headers":{"content-type":"text/plain"}},
			{"status_code":404,"body":{"hello":"world"},"headers":{"content-type":"application/json"}}
		]`)},
		"some/path/plain.text":   {Data: []byte("hello\nworld\n")},
		"some/path/invalid.json": {Data: []byte("invalid json")},
	}

	tests := []struct {
		path      string
		errPrefix string
		expected  []response
	}{
		{
			"some/path/valid.json",
			"",
			[]response{
				{
					StatusCode: http.StatusOK,
					Body:       []byte("hello"),
					Headers:    map[string]string{"content-type": "text/plain"},
					DelayMs:    1000,
				},
				{
					StatusCode: http.StatusNotFound,
					Body:       []byte("world"),
					Headers:    map[string]string{"content-type": "text/plain"},
					DelayMs:    0,
				},
				{
					StatusCode: http.StatusNotFound,
					Body:       []byte(`{"hello":"world"}`),
					Headers:    map[string]string{"content-type": "application/json"},
					DelayMs:    0,
				},
			},
		},
		{
			"some/path/plain.text",
			"",
			[]response{
				{StatusCode: http.StatusOK, Body: []byte("hello")},
				{StatusCode: http.StatusOK, Body: []byte("world")},
				{StatusCode: http.StatusOK, Body: []byte("")},
			},
		},
		{"some/path/invalid.json", "failed to unmarshal responses", []response{}},
		{"some/path/doesnt.exist", "failed to read responses file", []response{}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.path), func(t *testing.T) {
			hh := HTTPHandler{
				logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
				fs:     mfs,
			}

			_, err := hh.LoadResponsesFile(tt.path)

			// check if error is nil when expected prefix is empty
			if tt.errPrefix == "" && err != nil {
				t.Errorf("got \"%v\", expect nil", err)
			}

			// check if error message starts with expected prefix
			if !strings.HasPrefix(fmt.Errorf("%w", err).Error(), tt.errPrefix) {
				t.Errorf("got \"%v\", expect \"%v*\"", err, tt.errPrefix)
			}

			// check if responses in the file are equal to expected
			if err == nil && !reflect.DeepEqual(hh.responses, tt.expected) {
				t.Errorf("got \"%v\", expect \"%v\"", hh.responses, tt.expected)
			}
		})
	}
}
