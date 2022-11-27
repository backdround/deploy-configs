package fsutility

import "os"

type fileType int

const (
	Regular fileType = iota
	Symlink
	Directory
	Notexisting
	Unknown
)

func (f fileType) String() string {
	switch f {
	case Regular:
		return "Reguilar"
	case Symlink:
		return "Symlink"
	case Directory:
		return "Directory"
	case Notexisting:
		return "Notexisting"
	case Unknown:
		return "Unknown"
	}

	panic("unknown fileType")
}

// GetFileType returns fileType. If permission denied occur then
// returns unknown.
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

	if fileInfo.IsDir() {
		return Directory
	}

	if (fileInfo.Mode() & os.ModeSymlink) == os.ModeSymlink {
		return Symlink
	}

	return Unknown
}
