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
	Body       json.RawMessage   `json:"body"`
	Headers    map[string]string `json:"headers"`
	DelayMs    uint              `json:"delay_ms"`
	Repeat     int               `json:"repeat"`
	count      int
}

type HTTPHandler struct {
	sync.Mutex
	index           int
	responses       []response
	fs              fs.FS
	logger          *slog.Logger
	healthCheckPath string
}

var errNoResponsesLeft = errors.New("no responses left")

var healthResponse = &response{
	StatusCode: 200,
	Body:       []byte("ok"),
	Headers:    map[string]string{},
	DelayMs:    0,
	Repeat:     1,
	count:      1,
}

func NewHTTPHandler(logger *slog.Logger, healthCheckPath ...string) *HTTPHandler {
	if len(healthCheckPath) == 0 || healthCheckPath[0] == "" {
		healthCheckPath = []string{"/healthz"}
	}

	return &HTTPHandler{
		logger:          logger,
		fs:              osFileReader{},
		healthCheckPath: healthCheckPath[0],
	}
}

func (hh *HTTPHandler) LoadResponsesFile(name string) (int, error) {
	content, err := fs.ReadFile(hh.fs, name)
	if err != nil {
		return 0, fmt.Errorf("failed to read responses file [%s]: %w", name, err)
	}

	//nolint:nestif
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

			if hh.responses[i].Repeat == 0 {
				hh.responses[i].Repeat = 1
			}

			// if the body is a string, remove the quotes
			if string(hh.responses[i].Body[0]) == `"` {
				hh.responses[i].Body = hh.responses[i].Body[1 : len(hh.responses[i].Body)-1]
			}
		}
	} else {
		lines := strings.Split(string(content), "\n")
		hh.responses = make([]response, len(lines))

		for i, line := range lines {
			hh.responses[i] = response{StatusCode: 200, Body: []byte(line), Repeat: 1}
		}
	}

	return len(hh.responses), nil
}

func (hh *HTTPHandler) getNextResponse(path string) (*response, error) {
	if path == hh.healthCheckPath {
		return healthResponse, nil
	}

	if hh.index >= len(hh.responses) {
		return &response{}, errNoResponsesLeft
	}

	response := &hh.responses[hh.index]
	response.count++

	if response.Repeat == response.count {
		hh.index++
	}

	return response, nil
}

func (hh *HTTPHandler) HandleFunc(w http.ResponseWriter, rq *http.Request) {
	hh.Lock()
	defer hh.Unlock()

	r, err := hh.getNextResponse(rq.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		hh.logger.Error("responding", "error", err)

		return
	}

	delay := time.Duration(r.DelayMs) * time.Millisecond
	if r.DelayMs > 0 {
		time.Sleep(delay)
	}

	for k, v := range r.Headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(int(r.StatusCode))

	_, err = fmt.Fprintf(w, "%s\n", r.Body)
	if err != nil {
		hh.logger.Error("responding", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	hh.logger.Info("responding",
		"status_code", r.StatusCode,
		"body", string(r.Body),
		"headers", r.Headers,
		"delay", delay,
		"repeat", strings.Replace(fmt.Sprintf("%d/%d", r.count, r.Repeat), "-1", "∞", -1),
	)
}
