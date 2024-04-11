package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func Teardown() {
	osReadFile = os.ReadFile
	fmt.Println("teardown")
}

func Test_Listen(t *testing.T) {
	lines := []string{"hello", "world"}

	// mock os.ReadFile
	osReadFile = func(name string) ([]byte, error) {
		// append the file name to check it's being passed correctly to os.ReadFile
		lines = append(lines, name)
		return []byte(strings.Join(lines, "\n")), nil
	}

	go func() {
		err := Listen(8080, "responses.txt")
		t.Error(err)
	}()

	tests := []struct {
		want       string
		statusCode int
		err        error
	}{
		{"hello\n", http.StatusOK, nil},
		{"world\n", http.StatusOK, nil},
		{"responses.txt\n", http.StatusOK, nil},
		{"no responses left\n", http.StatusInternalServerError, nil},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			resp, err := http.Get("http://localhost:8080")
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != tt.statusCode {
				t.Errorf("got status code %d, want %d", resp.StatusCode, tt.statusCode)
			}

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if string(body) != tt.want {
				t.Errorf("got '%s', want '%s'", body, tt.want)
			}
		})
	}
}
