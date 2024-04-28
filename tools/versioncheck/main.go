package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func main() {
	pattern := regexp.MustCompile(`v[0-9]\.[0-9]\.[0-9]`)
	versionSet := map[string][]string{}
	extensions := []string{".go", ".md"}
	skipList := []string{"CHANGELOG.md", "CHANGELOG.tpl.md"}

	{
		flag.Parse()
		versionTag := strings.Replace(flag.Arg(0), "refs/tags/", "", 1)

		if versionTag != "" {
			versionSet[versionTag] = append(versionSet[versionTag], ".git")
		}
	}

	{
		fmt.Fprintf(os.Stderr, " \033[38;5;245m• checking CHANGELOG.md\033[0m\n")

		content, err := os.ReadFile("CHANGELOG.md")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read file \"CHANGELOG.md\": %v", err)
			os.Exit(1)
		}

		latestVersion := pattern.FindString(string(content))
		versionSet[latestVersion] = append(versionSet[latestVersion], "CHANGELOG.md")
	}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !slices.Contains(extensions, filepath.Ext(path)) {
			return nil
		}

		if slices.Contains(skipList, filepath.Base(path)) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %q: %w", path, err)
		}

		versions := pattern.FindAllString(string(content), -1)

		if len(versions) == 0 {
			fmt.Fprintf(os.Stderr, " \033[38;5;245m• checking %s\033[0m\n", path)
		} else {
			fmt.Fprintf(os.Stderr, " \033[1;93m•\033[0m checking %q %q\033[0m\n", path, versions)
		}

		for _, version := range versions {
			versionSet[version] = append(versionSet[version], path)
		}

		return nil
	})
	if err != nil {
		fmt.Println("error: failed to walk dir tree:", err)
		os.Exit(1)
	}

	fmt.Println("")

	versions := keys(versionSet)

	if len(versions) > 1 {
		fmt.Println("\033[91m✘ error: multiple versions found:\033[0m")

		for version, path := range versionSet {
			fmt.Fprintf(os.Stderr, "\033[38;5;210m  %q: %q\033[0m\n", version, path)
		}

		fmt.Println("")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "\033[92m✔ All files have the same version: %q\033[0m\n", versions[0])
}

func keys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
