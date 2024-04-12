package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

var responsesIndex atomic.Uint32

func getNextResponse(s []string, ai *atomic.Uint32) (string, error) {
	i := int(ai.Load())

	if i >= len(s) {
		return "", fmt.Errorf("no responses left")
	}

	ai.Add(1)

	return s[i], nil
}

// exported for mocking
var osReadFile = os.ReadFile

func readFileLines(name string) ([]string, error) {
	content, err := osReadFile(name)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(content), "\n"), nil
}

func Listen(port int, responsesFileName string) error {
	responses, err := readFileLines(responsesFileName)
	if err != nil {
		return fmt.Errorf("failed to read responses file [%s]: %w", responsesFileName, err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := getNextResponse(responses, &responsesIndex)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.Error("responding", "error", err)
			return
		}

		fmt.Fprintf(w, "%s\n", body)
		slog.Info("responding", "response", body)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
