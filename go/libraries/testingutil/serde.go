package testingutil

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func ReadToObject(t *testing.T, reader io.Reader, object interface{}) {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(object)
	if err != nil {
		t.Fatalf("failed to decode input, %v", err)
	}
}

func ReadFile(t *testing.T, path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to read file from path %v, %v", path, err)
	}
	return file
}
