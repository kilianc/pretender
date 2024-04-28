package pretender

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var discardLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestNewHTTPMux(t *testing.T) {
	responseFileName := fmt.Sprintf("%s/response.json", t.TempDir())

	err := os.WriteFile(responseFileName, []byte(`[{"delay_ms":1}]`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should error out when file fails to load", func(t *testing.T) {
		server, rn, err := NewServeMux("no.json", discardLogger)
		if err == nil {
			t.Errorf("expected an error, got %v %d", server, rn)
		}
	})

	t.Run("should serve responses", func(t *testing.T) {
		server, rn, err := NewServeMux(responseFileName, discardLogger)
		if err != nil {
			t.Errorf("expected an error, got %v %d", server, rn)
		}

		{
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/path", nil)
			server.ServeHTTP(w, r)
			server.ServeHTTP(w, r)

			if w.Result().StatusCode != http.StatusOK {
				t.Errorf("expected status code 200, got %d", w.Result().StatusCode)
			}

			server.ServeHTTP(w, r)
		}
	})

	t.Run("should serve health checks", func(t *testing.T) {
		server, rn, err := NewServeMux(responseFileName, discardLogger, "/custom-health-check")
		if err != nil {
			t.Errorf("expected an error, got %v %d", server, rn)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/custom-health-check", nil)
		server.ServeHTTP(w, r)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", w.Result().StatusCode)
		}
	})
}

func TestNewServer(t *testing.T) {
	responseFileName := fmt.Sprintf("%s/response.txt", t.TempDir())

	err := os.WriteFile(responseFileName, []byte(`[{"delay_ms":1}]`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should error out when file fails to load", func(t *testing.T) {
		server, _, err := NewServer(8080, "no.json", discardLogger)
		if err == nil {
			t.Errorf("expected an error, got %v", server)
		}
	})

	t.Run("should return a server", func(t *testing.T) {
		server, _, err := NewServer(8080, responseFileName, discardLogger)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if server == nil {
			t.Errorf("expected a server, got %v", server)
		}
	})
}
