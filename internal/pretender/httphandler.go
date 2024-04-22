package pretender

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type response struct {
	StatusCode uint              `json:"status_code"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
	DelayMs    uint              `json:"delay_ms"`
}

type HTTPHandler struct {
	sync.Mutex
	index     int
	responses []response
	fs        fs.FS
	logger    *slog.Logger
}

var errNoResponsesLeft = errors.New("no responses left")

func NewHTTPHandler(logger *slog.Logger) *HTTPHandler {
	return &HTTPHandler{
		logger: logger,
		fs:     osFileReader{},
	}
}

func (hh *HTTPHandler) LoadResponsesFile(name string) (int, error) {
	content, err := fs.ReadFile(hh.fs, name)
	if err != nil {
		return 0, fmt.Errorf("failed to read responses file [%s]: %w", name, err)
	}

	if strings.HasSuffix(name, ".json") {
		hh.responses = []response{}

		err = json.Unmarshal(content, &hh.responses)
		if err != nil {
			return 0, fmt.Errorf("failed to unmarshal responses: %w", err)
		}

		for i := range hh.responses {
			if hh.responses[i].StatusCode == 0 {
				hh.responses[i].StatusCode = 200
			}
		}
	} else {
		lines := strings.Split(string(content), "\n")
		hh.responses = make([]response, len(lines))

		for i, line := range lines {
			hh.responses[i] = response{StatusCode: 200, Body: line}
		}
	}

	return len(hh.responses), nil
}

func (hh *HTTPHandler) getNextResponse() (response, error) {
	if hh.index >= len(hh.responses) {
		return response{}, errNoResponsesLeft
	}

	response := hh.responses[hh.index]
	hh.index++

	return response, nil
}

func (hh *HTTPHandler) HandleFunc(w http.ResponseWriter, _ *http.Request) {
	hh.Lock()
	defer hh.Unlock()

	res, err := hh.getNextResponse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		hh.logger.Error("responding", "error", err)

		return
	}

	delay := time.Duration(res.DelayMs) * time.Millisecond
	if res.DelayMs > 0 {
		time.Sleep(delay)
	}

	for k, v := range res.Headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(int(res.StatusCode))

	_, err = fmt.Fprintf(w, "%s\n", res.Body)
	if err != nil {
		hh.logger.Error("responding", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	hh.logger.Info("responding",
		"status_code", res.StatusCode,
		"body", res.Body,
		"headers", res.Headers,
		"delay", delay,
	)
}
