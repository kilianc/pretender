package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"pretender/internal"
	"time"

	"github.com/lmittmann/tint"
)

const version = "v1.0.0"

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

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))

	err := internal.Listen(*port, *responseFileName)
	slog.Error("error starting server", "error", err)
	os.Exit(1)
}
