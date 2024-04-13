package pretender

import (
	"io/fs"
	"os"
)

type osFileReader struct {
	fs.FS
}

func (fr osFileReader) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
