package pretender

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

type HttpHandler struct {
	index     atomic.Uint32
	responses []string
	fs        fs.FS
	logger    *slog.Logger
}

func NewHttpHandler(logger *slog.Logger) *HttpHandler {
	return &HttpHandler{
		logger: logger,
	}
}

func (hh *HttpHandler) LoadResponsesFile(name string) error {
	// this is needed to allow relative paths as input from the CLI
	absoluteFilePath, err := filepath.Abs(name)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for responses file [%s]: %w", name, err)
	}

	// httpHandler.fs is only set in tests, we default to os filesystem
	if hh.fs == nil {
		hh.fs = os.DirFS(filepath.Dir(absoluteFilePath))
	}

	file, err := hh.fs.Open(filepath.Base(absoluteFilePath))
	if err != nil {
		return fmt.Errorf("failed to open responses file [%s]: %w", name, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read responses file [%s]: %w", name, err)
	}

	hh.responses = strings.Split(string(content), "\n")

	return nil
}

func (hh *HttpHandler) getNextResponse() (string, error) {
	i := int(hh.index.Load())

	if i >= len(hh.responses) {
		return "", fmt.Errorf("no responses left")
	}

	hh.index.Add(1)

	return hh.responses[i], nil
}

func (hh *HttpHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	body, err := hh.getNextResponse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		hh.logger.Error("responding", "error", err)
		return
	}

	fmt.Fprintf(w, "%s\n", body)
	hh.logger.Info("responding", "response", body)
}
