package main

import (
	"os"
	"testing"
)

func Test(_ *testing.T) {
	os.Chdir("../..")
	main()
}
