package fsutility

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

// FindEntryDescending searches an directory entry from the given
// topSearchPath and descending to root. It uses pathType as bitwise flags.
func FindEntryDescending(topSearchPath string, entryName string,
	types pathType) (
	desiredPath string, err error) {
	getDescendingParentDirectories := func(directory string) []string {
		parents := []string{directory}
		nextParent := path.Dir(directory)
		for directory != nextParent {
			parents = append(parents, nextParent)
			directory = nextParent
			nextParent = path.Dir(directory)
		}
		return parents
	}

	parentDirectories := getDescendingParentDirectories(topSearchPath)

	for _, currentDirectory := range parentDirectories {
		hypotheticalDesiredPath := path.Join(currentDirectory, entryName)
		pathType := GetPathType(hypotheticalDesiredPath)
		if pathType&types == pathType {
			return hypotheticalDesiredPath, nil
		}
	}

	return "", errors.New("Desired path isn't founded: " + entryName)
}

// GetFileHash calculates sha512 with file data.
// If file doesn't exist then it returns empty slice.
func GetFileHash(path string) []byte {
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
		pattern := "unable to create directory, because file exists: %q"
		return fmt.Errorf(pattern, directory)
	}

	return os.MkdirAll(directory, 0755)
}

func IsLinkPointsToDestination(linkPath string, destination string) bool {
	// Makes linkPath absolute
	if !path.IsAbs(linkPath) {
		wd, err := os.Getwd()
		if err != nil {
			return false
		}
		linkPath = path.Join(wd, linkPath)
	}

	makeAbsolute := func(baseDirectory string, p string) string {
		if path.IsAbs(p) {
			return p
		}
		p = path.Join(baseDirectory, p)
		p = path.Clean(p)

		return p
	}

	// Gets absolute destination
	linkDirectory := path.Dir(linkPath)
	destination = makeAbsolute(linkDirectory, destination)


	// Gets absolute link destination
	linkDestination, err := os.Readlink(linkPath)
	if err != nil {
		return false
	}
	linkDestination = makeAbsolute(linkDirectory, linkDestination)

	matched, err := path.Match(linkDestination, destination)

	if err != nil {
		return false
	}

	return matched
}
