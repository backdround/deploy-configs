package fsutility

import (
	"os"
	"testing"

	"github.com/backdround/deploy-configs/pkg/fstestutility"
	"github.com/stretchr/testify/require"
)

func TestNotExistingType(t *testing.T) {
	notexistingPath := fstestutility.GetAvailableTempPath()
	resultType := GetPathType(notexistingPath)
	require.Equal(t, Notexisting.String(), resultType.String())
}

func TestRegularType(t *testing.T) {
	// Gets an unexisting filename
	file, err := os.CreateTemp("", "go_test.*.txt")
	fstestutility.AssertNoError(err)
	defer os.Remove(file.Name())

	// Checks that file is not existing
	resultType := GetPathType(file.Name())
	require.Equal(t, Regular.String(), resultType.String())
}

func TestSymlinkType(t *testing.T) {
	t.Run("ValidSymlink", func(t *testing.T) {
		// Creates a valid link
		linkPath := fstestutility.GetAvailableTempPath()
		os.Symlink("/dev/null", linkPath)
		defer os.Remove(linkPath)

		// Checks that it's a link
		resultType := GetPathType(linkPath)
		require.Equal(t, Symlink.String(), resultType.String())
	})

	t.Run("BrokenSymlink", func(t *testing.T) {
		// Creates a broken link
		linkPath := fstestutility.GetAvailableTempPath()
		unexistingFilePath := fstestutility.GetAvailableTempPath()
		os.Symlink(unexistingFilePath, linkPath)
		defer os.Remove(linkPath)

		// Checks that it's a link
		resultType := GetPathType(linkPath)
		require.Equal(t, Symlink.String(), resultType.String())
	})
}

func TestDirectoryType(t *testing.T) {
	resultType := GetPathType(os.TempDir())
	require.Equal(t, Directory.String(), resultType.String())
}

func TestUnknownType(t *testing.T) {
	resultType := GetPathType("/dev/null")
	require.Equal(t, Unknown.String(), resultType.String())
}
