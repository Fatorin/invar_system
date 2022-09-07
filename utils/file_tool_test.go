package utils

import (
	"os"
	"testing"
)

func TestFilePathGenerator(t *testing.T) {
	path := FilePathGenerator("test.png", "../static/public/")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Log("Success")
		t.Log("Path=", path[1:])
	} else {
		t.Log("False")
	}
}
