package fsutility

import (
	"fmt"
	"os"
	"path"
)

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

// GetAvailableTempPath returns path to available (nonexistent)
// temporary file
func GetAvailableTempPath() string {
	file, err := os.CreateTemp("", "go_test.*.txt")
	path := file.Name()
	assertNoError(err)
	assertNoError(os.Remove(path))

	return path
}

// MakeDirectoryIfDoesntExist creates directory if it doesn't exist.
// the error is return if unable to create directory.
func MakeDirectoryIfDoesntExist(directory string) error {
	stat, err := os.Stat(directory)
	if err == nil {
		if stat.IsDir() {
			return nil
		}
		pattern := "Unable to create directory, because file exists: %q"
		return fmt.Errorf(pattern, directory)
	}

	return os.MkdirAll(directory, 0755)
}

// CreateTemporaryFiles creates files by the patterns (* for random substitution).
// It changes the patterns parameters to the paths. It returns a cleanup function.
func CreateTemporaryFiles(patterns ...*string) (cleanup func()) {
	filesToRemove := []string{}

	for _, pattern := range patterns {
		file, err := os.CreateTemp("", *pattern)
		assertNoError(err)
		file.Close()
		filesToRemove = append(filesToRemove, file.Name())
		*pattern = file.Name()
	}

	removeAllFiles := func() {
		for _, path := range filesToRemove {
			err := os.Remove(path)
			assertNoError(err)
		}
	}

	return removeAllFiles
}

func IsLinkPointsToDestination(linkPath string, destination string) bool {
	linkDestination, err := os.Readlink(linkPath)
	if err != nil {
		return false
	}

	matched, err := path.Match(linkDestination, destination)

	if err != nil {
		return false
	}

	return matched
}

