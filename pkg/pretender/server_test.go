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

func Test_NewHTTPMux(t *testing.T) {
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	responseFileName := fmt.Sprintf("%s/response.txt", t.TempDir())
	os.WriteFile(responseFileName, []byte(`hi!\n`), 0o644)

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

			if w.Result().StatusCode != http.StatusOK {
				t.Errorf("expected status code 200, got %d", w.Result().StatusCode)
			}
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
