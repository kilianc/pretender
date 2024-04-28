package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Chdir("../..")

	t.Run("from args", func(_ *testing.T) {
		os.Args = []string{"", "fix: test"}

		main()
	})

	t.Run("from git", func(_ *testing.T) {
		os.Args = []string{}

		main()
	})
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		message string
		want    bool
	}{
		{"fix: test", true},
		{"feat: test", true},
		{"chore: test", true},
		{"fix(test): test", true},
		{"feat(test): test", true},
		{"chore(test): test", true},
		{"fix(test): Test", false},
		{"fix(test):  test", false},
		{"fix(test): 0 aaa", false},
		{"fix something", false},
		{"fixes: something", false},
	}

	for _, test := range tests {
		t.Run(test.message, func(t *testing.T) {
			if isValid(test.message) != test.want {
				t.Errorf("isValid(%q) = %v; want %v", test.message, !test.want, test.want)
			}
		})
	}
}
