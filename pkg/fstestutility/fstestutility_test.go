package fstestutility_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"os"

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
