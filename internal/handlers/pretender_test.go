package handlers

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
	hh := Pretender{
		responses: []response{
			{
				StatusCode: http.StatusOK,
				Body:       []byte("hello"),
				Headers:    map[string]string{"content-type": "my/type1"},
				Repeat:     1,
			},
			{
				StatusCode: http.StatusOK,
				Body:       []byte("world"),
				Headers:    map[string]string{"content-type": "my/type2"},
				Repeat:     1,
			},
			{
				StatusCode: http.StatusOK,
				Body:       []byte("twice"),
				Repeat:     2,
			},
		},
		logger:          slog.New(slog.NewTextHandler(io.Discard, nil)),
		healthCheckPath: "/healthz",
	}

	tests := []struct {
		path        string
		statusCode  int
		body        string
		contentType string
	}{
		{
			path:        "/",
			statusCode:  http.StatusOK,
			body:        "hello\n",
			contentType: "my/type1",
		},
		{
			path:        "/",
			statusCode:  http.StatusOK,
			body:        "world\n",
			contentType: "my/type2",
		},
		{
			path:       "/healthz",
			statusCode: http.StatusOK,
			body:       "ok\n",
		},
		{
			path:       "/",
			statusCode: http.StatusOK,
			body:       "twice\n",
		},
		{
			path:       "/",
			statusCode: http.StatusOK,
			body:       "twice\n",
		},
		{
			path:        "/",
			statusCode:  http.StatusInternalServerError,
			body:        "no responses left\n",
			contentType: "text/plain; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(" %v %v %v", tt.path, tt.statusCode, tt.body), func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.path, nil)
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

func Test_HandleFuncNegativeRepeat(t *testing.T) {
	hh := Pretender{
		responses: []response{
			{
				StatusCode: http.StatusOK,
				Body:       []byte("hello"),
				Headers:    map[string]string{"content-type": "my/type1"},
				Repeat:     -1,
			},
		},
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	t.Run("should repeat forever", func(t *testing.T) {
		for range 100 {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			hh.HandleFunc(w, r)

			if w.Body.String() != "hello\n" {
				t.Fatalf("got %q, expect %q", w.Body, "hello\n")
			}
		}
	})
}

func Test_LoadResponsesFile(t *testing.T) {
	mfs := fstest.MapFS{
		"some/path/valid.json": {Data: []byte(`[
			{"body":"hello","headers":{"content-type":"text/plain"},"delay_ms":1000,"repeat":-1},
			{"status_code":404,"body":"world","headers":{"content-type":"text/plain"},"repeat":0},
			{"status_code":404,"body":{"hello":"world"},"headers":{"content-type":"application/json"},"repeat":5}
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
					Repeat:     -1,
				},
				{
					StatusCode: http.StatusNotFound,
					Body:       []byte("world"),
					Headers:    map[string]string{"content-type": "text/plain"},
					DelayMs:    0,
					Repeat:     1,
				},
				{
					StatusCode: http.StatusNotFound,
					Body:       []byte(`{"hello":"world"}`),
					Headers:    map[string]string{"content-type": "application/json"},
					DelayMs:    0,
					Repeat:     5,
				},
			},
		},
		{
			"some/path/plain.text",
			"",
			[]response{
				{StatusCode: http.StatusOK, Body: []byte("hello"), Repeat: 1},
				{StatusCode: http.StatusOK, Body: []byte("world"), Repeat: 1},
				{StatusCode: http.StatusOK, Body: []byte(""), Repeat: 1},
			},
		},
		{"some/path/invalid.json", "failed to unmarshal responses", []response{}},
		{"some/path/doesnt.exist", "failed to read responses file", []response{}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.path), func(t *testing.T) {
			hh := Pretender{
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
