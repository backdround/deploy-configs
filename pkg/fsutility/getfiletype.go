package fsutility

import "os"

type fileType int

const (
	Regular fileType = iota
	Symlink
	Notexisting
	Unknown
)

func (f fileType) String() string {
	switch f {
	case Regular:
		return "reguilar"
	case Symlink:
		return "symlink"
	case Notexisting:
		return "notexisting"
	case Unknown:
		return "unknown"
	}

	panic("unknown fileType")
}

func GetFileType(path string) fileType {
	fileInfo, err := os.Lstat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return Notexisting
		} else {
			return Unknown
		}
	}

	if fileInfo.Mode().IsRegular() {
		return Regular
	}

	if (fileInfo.Mode() & os.ModeSymlink) == os.ModeSymlink {
		return Symlink
	}

	return Unknown
}
