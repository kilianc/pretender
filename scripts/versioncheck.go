package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
)

func main() {
	pattern := regexp.MustCompile(`v[0-9]\.[0-9]\.[0-9]`)
	versionSet := map[string][]string{}
	extensions := []string{".go", ".md"}

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

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		versions := pattern.FindAllString(string(content), -1)
		for _, version := range versions {
			versionSet[version] = append(versionSet[version], path)
		}

		return nil
	})

	if err != nil {
		fmt.Println("error: failed to walk dir tree:", err)
		return
	}

	versions := keys(versionSet)

	if len(versions) > 1 {
		fmt.Println("error: multiple versions found")
		for version, path := range versionSet {
			fmt.Printf("  %s: %s\n", version, path)
		}
		os.Exit(1)
	}

	fmt.Println("All files have the same version:", versions[0])
}

func keys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
