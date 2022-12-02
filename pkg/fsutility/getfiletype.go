package fsutility

import "os"

type pathType int

const (
	Regular pathType = 1 << iota
	Symlink
	Directory
	Notexisting
	Unknown
)

func (f pathType) String() string {
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

	panic("unknown pathType")
}

// GetPathType returns pathType. If permission denied occur then
// returns unknown.
func GetPathType(path string) pathType {
	pathInfo, err := os.Lstat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return Notexisting
		} else {
			return Unknown
		}
	}

	if pathInfo.Mode().IsRegular() {
		return Regular
	}

	if pathInfo.IsDir() {
		return Directory
	}

	if (pathInfo.Mode() & os.ModeSymlink) == os.ModeSymlink {
		return Symlink
	}

	return Unknown
}
