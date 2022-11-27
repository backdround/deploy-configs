package fsutility

import (
	"crypto/sha512"
	"fmt"
	"io"
	"os"
	"path"
	"bytes"
)

// GetFileHash calculates sha512 with file data.
// If file doesn't exist then it returns empty slice.
func GetFileHash (path string) []byte {
	// Opens file
	file, err := os.Open(path)
	if err != nil {
		return []byte{}
	}
	defer file.Close()

	// Calculates hash
	hash := sha512.New()
	if _, err := io.Copy(hash, file); err != nil {
		return []byte{}
	}

	return hash.Sum(nil)
}

// GetHash calculates sha512 from given data
func GetHash(data []byte) []byte {
	dataReader := bytes.NewReader(data)

	// Calculates hash
	hash := sha512.New()
	if _, err := io.Copy(hash, dataReader); err != nil {
		return []byte{}
	}

	return hash.Sum(nil)
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

