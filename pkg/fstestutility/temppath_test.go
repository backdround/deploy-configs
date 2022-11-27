package fstestutility

import (
	"testing"

	"os"
)

func TestGetAvailableTempPath(t *testing.T) {
	path := GetAvailableTempPath()
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		t.Error("Path must be unoccupied")
	}
}
