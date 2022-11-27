package fstestutility

import (
	"os"
)

func assertNoError(err error) {
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
