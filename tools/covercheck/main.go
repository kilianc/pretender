package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	fmt.Fprintf(os.Stderr, "\n")

	coverageSummary, err := os.ReadFile("cover.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[91m✘ error: reading cover.txt: %v\033[0m\n\n", err)
		os.Exit(1)
	}

	regexp := regexp.MustCompile(`total:.+\t+(\d?\d?\d.\d)%`)

	match := regexp.FindStringSubmatch(string(coverageSummary))

	if len(match) < 2 {
		fmt.Fprintf(os.Stderr, "\033[91m✘ error: coverage not found in cover.txt\033[0m\n\n")
		os.Exit(1)
	}

	coverage, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[91m✘ error: failed to convert coverage to float: %v\033[0m\n\n", err)
		os.Exit(1)
	}

	if coverage != 100.0 {
		fmt.Fprintf(os.Stderr, "\033[91m✘ error: coverage is below 100%%: %.1f%%\033[0m\n\n", coverage)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "\033[92m✔ coverage is 100%%\033[0m\n")
}
