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

