package pretender

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

type HttpHandler struct {
	sync.Mutex
	index     int
	responses []string
	fs        fs.ReadFileFS
	logger    *slog.Logger
}

func NewHttpHandler(logger *slog.Logger) *HttpHandler {
	return &HttpHandler{
		logger: logger,
		fs:     osFileReader{},
	}
}

func (hh *HttpHandler) LoadResponsesFile(name string) error {
	content, err := hh.fs.ReadFile(name)
	if err != nil {
		return fmt.Errorf("failed to read responses file [%s]: %w", name, err)
	}

	hh.responses = strings.Split(string(content), "\n")

	return nil
}

func (hh *HttpHandler) getNextResponse() (string, error) {
	if hh.index >= len(hh.responses) {
		return "", fmt.Errorf("no responses left")
	}

	response := hh.responses[hh.index]
	hh.index++

	return response, nil
}

func (hh *HttpHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	hh.Lock()

	body, err := hh.getNextResponse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		hh.logger.Error("responding", "error", err)
		hh.Unlock()
		return
	}

	fmt.Fprintf(w, "%s\n", body)
	hh.logger.Info("responding", "response", body)
	hh.Unlock()
}
