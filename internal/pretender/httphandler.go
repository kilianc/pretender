package pretender

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type HttpHandler struct {
	sync.Mutex
	index     int
	responses []response
	fs        fs.FS
	logger    *slog.Logger
}

func NewHttpHandler(logger *slog.Logger) *HttpHandler {
	return &HttpHandler{
		logger: logger,
		fs:     osFileReader{},
	}
}

func (hh *HttpHandler) LoadResponsesFile(name string) (int, error) {
	content, err := hh.fs.ReadFile(name)
	if err != nil {
		return 0, fmt.Errorf("failed to read responses file [%s]: %w", name, err)
	}

	if strings.HasSuffix(name, ".json") {
		hh.responses = []response{}

		err = json.Unmarshal(content, &hh.responses)
		if err != nil {
			return 0, fmt.Errorf("failed to unmarshal responses: %w", err)
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

func (hh *HttpHandler) getNextResponse() (response, error) {
	if hh.index >= len(hh.responses) {
		return response{}, fmt.Errorf("no responses left")
	}

	response := hh.responses[hh.index]
	hh.index++

	return response, nil
}

func (hh *HttpHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	hh.Lock()
	defer hh.Unlock()

	res, err := hh.getNextResponse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		hh.logger.Error("responding", "error", err)
		return
	}

	if res.Delay > 0 {
		time.Sleep(res.Delay)
	}

	for k, v := range res.Headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(int(res.StatusCode))
	fmt.Fprintf(w, "%s\n", res.Body)

	hh.logger.Info("responding",
		"status_code", res.StatusCode,
		"body", res.Body,
		"headers", res.Headers,
		"delay", res.Delay,
	)
}
