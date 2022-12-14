// fstestutility describes functions that are useful for tests
package fstestutility

import (
	"os"
	"path"
)

func AssertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

// CreateTemporaryFiles creates files by the patterns (* for random substitution).
// It changes the patterns parameters to the paths. It returns a cleanup function.
func CreateTemporaryFiles(patterns ...*string) (cleanup func()) {
	filesToRemove := []string{}

	for _, pattern := range patterns {
		file, err := os.CreateTemp("", *pattern)
		AssertNoError(err)
		file.Close()
		filesToRemove = append(filesToRemove, file.Name())
		*pattern = file.Name()
	}

	removeAllFiles := func() {
		for _, path := range filesToRemove {
			err := os.Remove(path)
			AssertNoError(err)
		}
	}

	return removeAllFiles
}

// CreateTemporaryFileWithData Creates a file with given data.
// It returns path to file and a cleanup function.
func CreateTemporaryFileWithData(data string) (path string, cleanup func()) {
	path = GetAvailableTempPath()
	err := os.WriteFile(path, []byte(data), 0644)
	AssertNoError(err)

	return path, func() {
		os.Remove(path)
	}
}

// MakeTempDirectory creates temporary directory by pattern and
// returns cleanup function.
func MakeTempDirectory(pattern string) (directoryPath string, cleanup func()) {
	p, err := os.MkdirTemp("", pattern)
	AssertNoError(err)
	cleanup = func() { os.Remove(p) }
	return p, cleanup
}

// MakeDirectory creates directory by path parts.
func MakeDirectory(pathParts ...string) string {
	resultPath := path.Join(pathParts...)
	err := os.MkdirAll(resultPath, 0755)
	AssertNoError(err)
	return resultPath
}
