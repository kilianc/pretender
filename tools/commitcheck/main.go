package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	regexpStr = `^(feat|fix|chore)(\(.+\))?: [a-z].+`
	re        = regexp.MustCompile(regexpStr)
)

func main() {
	message := ""

	if len(os.Args) > 1 {
		message = os.Args[1]
	}

	if message == "" {
		commit, err := exec.Command("git", "log", "-1", "--pretty=%B").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[91m✘ error: git log -1 --pretty=%%B: %v\033[0m\n\n", err)
			os.Exit(1)
		}

		message = strings.Trim(string(commit), "\n")
	}

	if !isValid(message) {
		fmt.Fprintf(os.Stderr, "\033[91m✘ error: message %q does not match %s\033[0m\n\n", message, regexpStr)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "\033[92m✔ commit message %q is valid\033[0m\n", message)
}

func isValid(message string) bool {
	return re.Match([]byte(message))
}
