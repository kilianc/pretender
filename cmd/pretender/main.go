package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kilianc/pretender/pkg/pretender"
	"github.com/lmittmann/tint"
	"golang.org/x/term"
)

const version = "v1.8.0"

var (
	isTTY            = term.IsTerminal(int(os.Stdout.Fd()))
	printVersion     = flag.Bool("version", false, "print version and exit")
	responseFileName = flag.String("responses", "responses.json", "path to the file with responses")
	host             = flag.String("host", "127.0.0.1", "host to bind the server to")
	port             = flag.Int("port", 8080, "port to listen")
	noColor          = flag.Bool("no-color", false, "disable color output")
)

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	p("")
	p("██████╗ ██████╗ ███████╗████████╗███████╗███╗   ██╗██████╗ ███████╗██████╗")
	p("██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██╔════╝████╗  ██║██╔══██╗██╔════╝██╔══██╗")
	p("██████╔╝██████╔╝█████╗     ██║   █████╗  ██╔██╗ ██║██║  ██║█████╗  ██████╔╝")
	p("██╔═══╝ ██╔══██╗██╔══╝     ██║   ██╔══╝  ██║╚██╗██║██║  ██║██╔══╝  ██╔══██╗")
	p("██║     ██║  ██║███████╗   ██║   ███████╗██║ ╚████║██████╔╝███████╗██║  ██║")
	p("╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝ %s", version)
	p("")
	p("\033[32m•\033[0m starting server on %s:%d", *host, *port)
	p("\033[32m•\033[0m using responses file: %s", *responseFileName)
	p("\033[32m•\033[0m press ctrl+c to stop")
	p("")

	logger := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			NoColor:    !isTTY || *noColor,
		}),
	)

	healthCheckPath := os.Getenv("PRETENDER_HEALTH_CHECK_PATH")

	server, rn, err := pretender.NewServer(*host, *port, *responseFileName, logger, healthCheckPath)
	if err != nil {
		logger.Error("creating server", "error", err)
		os.Exit(1)
	}

	logger.Info("loaded responses from file", "file", *responseFileName, "count", rn)

	go func() {
		err = server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("starting server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	signal := <-quit
	logger.Info("shutting down server", "signal", signal.String())

	err = server.Shutdown(context.Background())
	if err != nil {
		logger.Error("shutting down server", "error", err)
		os.Exit(1)
	}
}

func p(format string, a ...any) {
	s := fmt.Sprintf(format, a...)

	if !isTTY || *noColor {
		s = strings.ReplaceAll(s, "\033[32m", "")
		s = strings.ReplaceAll(s, "\033[0m", "")
	}

	fmt.Fprintln(os.Stderr, s)
}
