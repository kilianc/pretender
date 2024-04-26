package pretender

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kilianc/pretender/internal/handlers"
)

var ErrorLoadingResponsesFile = fmt.Errorf("responses file")

func NewHTTPHandler(
	responseFileName string,
	logger *slog.Logger,
	healthCheckPath ...string,
) (*handlers.Pretender, int, error) {
	hh := handlers.NewPretender(logger, healthCheckPath...)

	rn, err := hh.LoadResponsesFile(responseFileName)
	if err != nil {
		return nil, 0, errors.Join(ErrorLoadingResponsesFile, err)
	}

	return hh, rn, nil
}

func NewServeMux(
	responseFileName string,
	logger *slog.Logger,
	healthCheckPath ...string,
) (*http.ServeMux, int, error) {
	hh, rn, err := NewHTTPHandler(responseFileName, logger, healthCheckPath...)
	if err != nil {
		return nil, 0, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", hh.HandleFunc)

	return mux, rn, nil
}

func NewServer(
	port int,
	responseFileName string,
	logger *slog.Logger,
	healthCheckPath ...string,
) (*http.Server, int, error) {
	mux, rn, err := NewServeMux(responseFileName, logger, healthCheckPath...)
	if err != nil {
		return nil, 0, err
	}

	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}, rn, nil
}
