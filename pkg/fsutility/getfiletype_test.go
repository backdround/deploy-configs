package fsutility

import (
	"testing"

	"github.com/stretchr/testify/require"

	"os"
)

func TestNotExistingType(t *testing.T) {
	notexistingPath := GetAvailableTempPath()
	resultType := GetFileType(notexistingPath)
	require.Equal(t, Notexisting.String(), resultType.String())
}

func TestRegularType(t *testing.T) {
	// Gets an unexisting filename
	file, err := os.CreateTemp("", "go_test.*.txt")
	assertNoError(err)
	defer os.Remove(file.Name())

	// Checks that file is not existing
	resultType := GetFileType(file.Name())
	require.Equal(t, Regular.String(), resultType.String())
}

func TestSymlinkType(t *testing.T) {
	t.Run("ValidSymlink", func(t *testing.T) {
		// Creates a valid link
		linkPath := GetAvailableTempPath()
		os.Symlink("/dev/null", linkPath)
		defer os.Remove(linkPath)

		// Checks that it's a link
		resultType := GetFileType(linkPath)
		require.Equal(t, Symlink.String(), resultType.String())
	})

	t.Run("BrokenSymlink", func(t *testing.T) {
		// Creates a broken link
		linkPath := GetAvailableTempPath()
		unexistingFilePath := GetAvailableTempPath()
		os.Symlink(unexistingFilePath, linkPath)
		defer os.Remove(linkPath)

		// Checks that it's a link
		resultType := GetFileType(linkPath)
		require.Equal(t, Symlink.String(), resultType.String())
	})
}

func TestDirectoryType (t *testing.T) {
	resultType := GetFileType(os.TempDir())
	require.Equal(t, Directory.String(), resultType.String())
}

func TestUnknownType (t *testing.T) {
	resultType := GetFileType("/dev/null")
	require.Equal(t, Unknown.String(), resultType.String())
}
