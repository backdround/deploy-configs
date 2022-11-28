package pathexpander

import (
	"os"
	"path"
	"testing"

	"github.com/backdround/deploy-configs/pkg/fstestutility"
	"github.com/stretchr/testify/require"
)

func TestGetGitRoot(t *testing.T) {
	t.Run("GitOnTop", func(t *testing.T) {
		// Creates a test root directory tree
		testRootDirectory, err := os.MkdirTemp("", "go_test_.*.d")
		fstestutility.AssertNoError(err)
		defer os.RemoveAll(testRootDirectory)

		// Creates a top level directory
		gitDirectory := path.Join(testRootDirectory, ".git")
		err = os.MkdirAll(gitDirectory, 0755)
		fstestutility.AssertNoError(err)

		// Executes the test
		resultDirectoryWithGit, err := getGitRoot(testRootDirectory)
		require.NoError(t, err)

		// Asserts expectations
		match, err := path.Match(path.Dir(gitDirectory), resultDirectoryWithGit)
		require.NoError(t, err)
		require.True(t, match)
	})

	t.Run("GitInMiddle", func(t *testing.T) {
		// Creates a test root directory tree
		testRootDirectory, err := os.MkdirTemp("", "go_test_.*.d")
		fstestutility.AssertNoError(err)
		defer os.RemoveAll(testRootDirectory)

		// Creates a top level directory
		topLevelDirectory := path.Join(testRootDirectory, "a/b/c")
		err = os.MkdirAll(topLevelDirectory, 0755)
		fstestutility.AssertNoError(err)

		// Creates the top level directory
		directoryWithGitDirectory := path.Join(testRootDirectory, "a/b")
		gitDirectory := path.Join(directoryWithGitDirectory, ".git")
		err = os.MkdirAll(gitDirectory, 0755)
		fstestutility.AssertNoError(err)

		// Executes the test
		resultDirectoryWithGit, err := getGitRoot(topLevelDirectory)
		require.NoError(t, err)

		// Asserts expectations
		match, err := path.Match(path.Dir(gitDirectory), resultDirectoryWithGit)
		require.NoError(t, err)
		require.True(t, match)
	})

	t.Run("GitDirectoryDoesntExist", func(t *testing.T) {
		// Creates a test root directory tree
		testRootDirectory, err := os.MkdirTemp("", "go_test_.*.d")
		fstestutility.AssertNoError(err)
		defer os.RemoveAll(testRootDirectory)

		// Creates the top level directory
		topLevelDirectory := path.Join(testRootDirectory, "a/b/c")
		err = os.MkdirAll(topLevelDirectory, 0755)
		fstestutility.AssertNoError(err)

		// Executes the test
		_, err = getGitRoot(topLevelDirectory)

		require.Error(t, err)
	})

	t.Run("InitialDirectoryDoesntExist", func(t *testing.T) {
		_, err := getGitRoot(fstestutility.GetAvailableTempPath())
		require.Error(t, err)
	})
}
