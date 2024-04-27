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

const version = "v1.6.1"

var isTTY = term.IsTerminal(int(os.Stdout.Fd()))

func main() {
	printVersion := flag.Bool("version", false, "print version and exit")
	responseFileName := flag.String("responses", "responses.json", "path to the file with responses")
	port := flag.Int("port", 8080, "port to listen")
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
	p("\033[32m•\033[0m starting server on port %d", *port)
	p("\033[32m•\033[0m using responses file: %s", *responseFileName)
	p("\033[32m•\033[0m press ctrl+c to stop")
	p("")

	logger := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			NoColor:    !isTTY,
		}),
	)

	server, rn, err := pretender.NewServer(*port, *responseFileName, logger, os.Getenv("PRETENDER_HEALTH_CHECK_PATH"))
	if err != nil {
		logger.Error("error loading responses file", "error", err)
		os.Exit(1)
	}

	logger.Info("loaded responses from file", "file", *responseFileName, "count", rn)

	go func() {
		err = server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error starting server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	signal := <-quit
	logger.Info("shutting down server", "signal", signal.String())

	err = server.Shutdown(context.Background())
	if err != nil {
		logger.Error("error shutting down server", "error", err)
		os.Exit(1)
	}
}

func p(format string, a ...any) {
	s := fmt.Sprintf(format, a...)

	if !isTTY {
		s = strings.ReplaceAll(s, "\033[32m", "")
		s = strings.ReplaceAll(s, "\033[0m", "")
	}

	fmt.Fprintln(os.Stderr, s)
}
