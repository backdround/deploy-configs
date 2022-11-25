package deploy

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"os"
)

func TestGetFileType(t *testing.T) {
	t.Run("NotExistingType", func(t *testing.T) {
		notexistingPath := getNotExistingPath()
		resultType := getFileType(notexistingPath)
		require.Equal(t, notexisting.String(), resultType.String())
	})

	t.Run("RegularType", func(t *testing.T) {
		// Gets an unexisting filename
		file, err := os.CreateTemp("", "go_test.*.txt")
		assertNoError(err)
		defer os.Remove(file.Name())

		// Checks that file is not existing
		resultType := getFileType(file.Name())
		require.Equal(t, regular.String(), resultType.String())
	})

	t.Run("SymlinkType", func(t *testing.T) {
		t.Run("ValidSymlink", func(t *testing.T) {
			// Creates a valid link
			linkPath := getNotExistingPath()
			os.Symlink("/dev/null", linkPath)
			defer os.Remove(linkPath)

			// Checks that it's a link
			resultType := getFileType(linkPath)
			require.Equal(t, symlink.String(), resultType.String())
		})

		t.Run("BrokenSymlink", func(t *testing.T) {
			// Creates a broken link
			linkPath := getNotExistingPath()
			unexistingFilePath := getNotExistingPath()
			os.Symlink(unexistingFilePath, linkPath)
			defer os.Remove(linkPath)

			// Checks that it's a link
			resultType := getFileType(linkPath)
			require.Equal(t, symlink.String(), resultType.String())
		})
	})

	t.Run("UnknownType", func(t *testing.T) {
		t.Run("DirectoryIsUnknownType", func(t *testing.T) {
			resultType := getFileType(os.TempDir())
			require.Equal(t, unknown.String(), resultType.String())
		})

		t.Run("DeviceIsUnknownType", func(t *testing.T) {
			resultType := getFileType("/dev/null")
			require.Equal(t, unknown.String(), resultType.String())
		})
	})
}

func TestMakeDirectoryIfDoesntExist(t *testing.T) {
	t.Run("CreateNestedDirectory", func(t *testing.T) {
		// Makes directory path to create
		rootDirectory := getNotExistingPath()
		newDirectory := path.Join(rootDirectory, "level-two")
		defer os.RemoveAll(rootDirectory)

		// Executes the test
		err := makeDirectoryIfDoesntExist(newDirectory)
		require.NoError(t, err)

		// Asserts that the directory was created
		stat, err := os.Stat(newDirectory)
		require.NoError(t, err)
		require.True(t, stat.IsDir())
	})

	t.Run("FailOnExistingFile", func(t *testing.T) {
		// Makes directory path to create
		file, err := os.CreateTemp("", "go_test.*.txt")
		assertNoError(err)
		defer os.Remove(file.Name())

		dirictoryToCreate := file.Name()

		// Executes the test
		err = makeDirectoryIfDoesntExist(dirictoryToCreate)

		// Asserts
		require.Error(t, err)
	})

	t.Run("NoErrorOnExistingDirectory", func(t *testing.T) {
		// Makes directory path to create
		directory, err := os.MkdirTemp("", "go_test.*.d")
		assertNoError(err)
		defer os.Remove(directory)

		// Executes the test
		err = makeDirectoryIfDoesntExist(directory)

		// Asserts
		require.NoError(t, err)
	})
}
