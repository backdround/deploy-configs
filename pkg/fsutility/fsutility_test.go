package fsutility

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"os"
)

func TestGetNotExistingPath(t *testing.T) {
	path := GetNotExistingPath()
	_, err := os.Stat(path)
	require.True(t, os.IsNotExist(err))
}

func TestMakeDirectoryIfDoesntExist(t *testing.T) {
	t.Run("CreateNestedDirectory", func(t *testing.T) {
		// Makes directory path to create
		rootDirectory := GetNotExistingPath()
		newDirectory := path.Join(rootDirectory, "level-two")
		defer os.RemoveAll(rootDirectory)

		// Executes the test
		err := MakeDirectoryIfDoesntExist(newDirectory)
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
		err = MakeDirectoryIfDoesntExist(dirictoryToCreate)

		// Asserts
		require.Error(t, err)
	})

	t.Run("NoErrorOnExistingDirectory", func(t *testing.T) {
		// Makes directory path to create
		directory, err := os.MkdirTemp("", "go_test.*.d")
		assertNoError(err)
		defer os.Remove(directory)

		// Executes the test
		err = MakeDirectoryIfDoesntExist(directory)

		// Asserts
		require.NoError(t, err)
	})
}

func TestCreateTemporaryFiles(t *testing.T) {
	t.Run("CreateFiles", func (t *testing.T) {
		// Creates files
		file1 := "go_test.simplefile.*.txt"
		file2 := "go_test.simplefile.*.txt"
		cleanup := CreateTemporaryFiles(&file1, &file2)
		defer cleanup()

		// Check that the files exist
		fileInfo1, err := os.Stat(file1)
		require.NoError(t, err)
		require.True(t, fileInfo1.Mode().IsRegular())

		fileInfo2, err := os.Stat(file1)
		require.NoError(t, err)
		require.True(t, fileInfo2.Mode().IsRegular())
	})

	t.Run("CleanupDeletesFiles", func (t *testing.T) {
		// Creates a file
		file := "go_test.simplefile.*.txt"
		cleanup := CreateTemporaryFiles(&file)
		// Removes the file
		cleanup()

		// Checks that the file is removed
		_, err := os.Stat(file)
		require.True(t, os.IsNotExist(err))
	})
}

func TestIsLinkPointsToDestination(t *testing.T) {
	t.Run("LinkPointsToDestination", func(t *testing.T) {
		target := "/dev/null"
		link := GetNotExistingPath()
		err := os.Symlink(target, link)
		assertNoError(err)
		defer os.Remove(link)

		// Asserts
		require.True(t, IsLinkPointsToDestination(link, target))
	})

	t.Run("LinkDoesntPointToDestination", func(t *testing.T) {
		target := GetNotExistingPath()
		link := GetNotExistingPath()
		err := os.Symlink("/dev/null", link)
		assertNoError(err)
		defer os.Remove(link)

		// Asserts
		require.False(t, IsLinkPointsToDestination(link, target))
	})
}
