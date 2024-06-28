package pretender

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kilianc/pretender/internal/handlers"
)

// ErrorLoadingResponsesFile is the error returned when the responses file can't be loaded.
var ErrorLoadingResponsesFile = fmt.Errorf("loading responses file")

// NewHTTPHandler creates a new [http] handler function configured to serve the responses
// defined in the responseFileName. It also returns the number of responses loaded from the file.
// If the file can't be loaded, it returns an error.
// The healthCheckPath is an optional parameter to define a custom path for the health check endpoint.
// If not provided, the default path is "/healthz".
// The logger is used to log the requests and responses.
func NewHTTPHandler(
	responseFileName string,
	logger *slog.Logger,
	healthCheckPath ...string,
) (func(http.ResponseWriter, *http.Request), int, error) {
	hh := handlers.NewPretender(logger, healthCheckPath...)

	rn, err := hh.LoadResponsesFile(responseFileName)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %w", ErrorLoadingResponsesFile, err)
	}

	return hh.HandleFunc, rn, nil
}

// NewServeMux creates a new [http.NewServeMux] configured to serve the responses defined in the
// responseFileName. It also returns the number of responses loaded from the file.
// If the file can't be loaded, it returns an error.
// The healthCheckPath is an optional parameter to define a custom path for the health check endpoint.
// If not provided, the default path is "/healthz".
// The logger is used to log the requests and responses.
func NewServeMux(
	responseFileName string,
	logger *slog.Logger,
	healthCheckPath ...string,
) (*http.ServeMux, int, error) {
	handler, rn, err := NewHTTPHandler(responseFileName, logger, healthCheckPath...)
	if err != nil {
		return nil, 0, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	return mux, rn, nil
}

// NewServer creates a new [http.Server] with the given port configured to serve the responses
// defined in the responseFileName. It also returns the number of responses loaded from the file.
// If the file can't be loaded, it returns an error.
// The healthCheckPath is an optional parameter to define a custom path for the health check endpoint.
// If not provided, the default path is "/healthz".
// The logger is used to log the requests and responses.
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

	return &http.Server{Addr: fmt.Sprintf("%d", port), Handler: mux}, rn, nil
}
