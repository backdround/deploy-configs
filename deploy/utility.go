package deploy

import (
	"fmt"
	"os"
)

type fileType int

const (
	regular fileType = iota
	symlink
	notexisting
	unknown
)

func (f fileType) String() string {
	switch f {
	case regular:
		return "reguilar"
	case symlink:
		return "symlink"
	case notexisting:
		return "notexisting"
	case unknown:
		return "unknown"
	}

	panic("unknown fileType")
}

func getFileType(path string) fileType {
	fileInfo, err := os.Lstat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return notexisting
		} else {
			return unknown
		}
	}

	if fileInfo.Mode().IsRegular() {
		return regular
	}

	if (fileInfo.Mode() & os.ModeSymlink) == os.ModeSymlink {
		return symlink
	}

	return unknown
}

func makeDirectoryIfDoesntExist(directory string) error {
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
