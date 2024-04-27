package handlers

import (
	"io/fs"
	"os"
)

type osFileReader struct {
	fs.FS
}

func (fr osFileReader) ReadFile(name string) ([]byte, error) {
	//nolint
	return os.ReadFile(name)
}
