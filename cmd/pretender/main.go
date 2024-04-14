package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kilianc/pretender/internal/pretender"
	"github.com/lmittmann/tint"
)

const version = "v1.0.2"

func main() {
	responseFileName := flag.String("responses", "responses.txt", "path to the file with responses")
	port := flag.Int("port", 8080, "port to listen")
	flag.Parse()

	fmt.Println("")
	fmt.Println("██████╗ ██████╗ ███████╗████████╗███████╗███╗   ██╗██████╗ ███████╗██████╗ ")
	fmt.Println("██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██╔════╝████╗  ██║██╔══██╗██╔════╝██╔══██╗")
	fmt.Println("██████╔╝██████╔╝█████╗     ██║   █████╗  ██╔██╗ ██║██║  ██║█████╗  ██████╔╝")
	fmt.Println("██╔═══╝ ██╔══██╗██╔══╝     ██║   ██╔══╝  ██║╚██╗██║██║  ██║██╔══╝  ██╔══██╗")
	fmt.Println("██║     ██║  ██║███████╗   ██║   ███████╗██║ ╚████║██████╔╝███████╗██║  ██║")
	fmt.Println("╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝", version)
	fmt.Println("")
	fmt.Printf("\033[32m•\033[0m starting server on port %d\n", *port)
	fmt.Printf("\033[32m•\033[0m using responses file: %s\n", *responseFileName)
	fmt.Printf("\033[32m•\033[0m press ctrl+c to stop\n")
	fmt.Println("")

	logger := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	)

	hh := pretender.NewHttpHandler(logger)

	err := hh.LoadResponsesFile(*responseFileName)
	if err != nil {
		logger.Error("error loading responses file", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", hh.HandleFunc)
	server := &http.Server{Addr: fmt.Sprintf(":%d", *port), Handler: mux}

	go func() {
		err = server.ListenAndServe()
		if err != http.ErrServerClosed {
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
