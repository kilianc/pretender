package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	fmt.Fprintf(os.Stderr, "\n")

	coverageSummary, err := os.ReadFile("cover.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[91m✘ error: reading cover.txt: %v\033[0m\n", err)
		os.Exit(1)
	}

	regexp := regexp.MustCompile(`total:.+(100.0%)`)

	match := regexp.FindStringSubmatch(string(coverageSummary))

	if len(match) < 2 {
		fmt.Fprintf(os.Stderr, "\033[91m✘ error: coverage not found in cover.txt\033[0m\n")
		os.Exit(1)
	}

	if match[1] != "100.0%" {
		fmt.Fprintf(os.Stderr, "\033[91m✘ error: coverage is below 100%%: %s\033[0m\n", match[1])
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "\033[92m✔ coverage is 100%%\033[0m\n")
}
