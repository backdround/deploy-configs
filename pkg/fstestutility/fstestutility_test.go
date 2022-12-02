package fstestutility_test

import (
	"os"

	"testing"
	"github.com/stretchr/testify/require"

	"github.com/backdround/deploy-configs/pkg/fstestutility"
)

func TestCreateTemporaryFiles(t *testing.T) {
	t.Run("CreateFiles", func (t *testing.T) {
		// Creates files
		file1 := "go_test.simplefile.*.txt"
		file2 := "go_test.simplefile.*.txt"
		cleanup := fstestutility.CreateTemporaryFiles(&file1, &file2)
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
		cleanup := fstestutility.CreateTemporaryFiles(&file)
		// Removes the file
		cleanup()

		// Checks that the file is removed
		_, err := os.Stat(file)
		require.True(t, os.IsNotExist(err))
	})
}

func TestCreateTemporaryFileWithData(t *testing.T) {
	data := "some data"
	path, cleanup := fstestutility.CreateTemporaryFileWithData(data)

	// Asserts that the file is created
	file, err := os.Stat(path)
	require.NoError(t, err)
	require.True(t, file.Mode().IsRegular())

	// Asserts that the cleanup removes file
	cleanup()
	_, err = os.Stat(path)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestMakeTempDirectory(t *testing.T) {
	path, cleanup := fstestutility.MakeTempDirectory("")
	directoryInfo, err := os.Stat(path)

	// Asserts created directory
	require.NoError(t, err)
	require.True(t, directoryInfo.IsDir())

	// Asserts directory deletion
	cleanup()
	_, err = os.Stat(path)
	require.True(t, os.IsNotExist(err))
}

func TestMakeDirectory(t *testing.T) {
	directoryBase := fstestutility.GetAvailableTempPath()
	path := fstestutility.MakeDirectory(directoryBase, "some", "temp", "dir")

	// Asserts created directory
	directoryInfo, err := os.Stat(path)
	require.NoError(t, err)
	require.True(t, directoryInfo.IsDir())

	os.RemoveAll(directoryBase)
}
