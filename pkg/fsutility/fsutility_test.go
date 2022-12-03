package fsutility

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/backdround/deploy-configs/pkg/fstestutility"
	"github.com/stretchr/testify/require"
)

func TestFindEntryDescending(t *testing.T) {
	t.Run("PathOnTop", func(t *testing.T) {
		// Creates a test directory tree
		baseDirectory, cleanup := fstestutility.MakeTempDirectory("test_.*.d")
		defer cleanup()
		gitPath := fstestutility.MakeDirectory(baseDirectory, ".git")

		// Executes the test
		resultDesiredDirectory, err := FindEntryDescending(baseDirectory,
			".git", Directory)
		require.NoError(t, err)

		// Asserts expectations
		match, err := path.Match(gitPath, resultDesiredDirectory)
		require.NoError(t, err)
		require.True(t, match)
	})

	t.Run("PathInTheMiddle", func(t *testing.T) {
		// Creates a test directory tree
		baseDirectory, cleanup := fstestutility.MakeTempDirectory("test_.*.d")
		defer cleanup()
		topLevelPath := path.Join(baseDirectory, "/a/b/c")
		gitPath := fstestutility.MakeDirectory(baseDirectory, "/a", ".git")

		// Executes the test
		resultDesiredDirectory, err := FindEntryDescending(topLevelPath,
			".git", Directory)

		// Asserts expectations
		require.NoError(t, err)
		match, err := path.Match(resultDesiredDirectory, gitPath)
		require.NoError(t, err)
		require.True(t, match)
	})

	t.Run("IncorrectPathType", func(t *testing.T) {
		// Creates a test directory tree
		baseDirectory, cleanup := fstestutility.MakeTempDirectory("test_.*.d")
		defer cleanup()
		fstestutility.MakeDirectory(baseDirectory, ".git")

		// Executes the test
		_, err := FindEntryDescending(baseDirectory, ".git", Regular)
		require.Error(t, err)
	})

	t.Run("SeveralPathTypes", func(t *testing.T) {
		// Creates a test directory tree
		baseDirectory, cleanup := fstestutility.MakeTempDirectory("test_.*.d")
		defer cleanup()
		gitPath := fstestutility.MakeDirectory(baseDirectory, ".git")

		// Executes the test
		severalPathTypes := Regular | Directory
		resultDesiredPath, err := FindEntryDescending(baseDirectory, ".git",
			severalPathTypes)

		// Asserts expectations
		require.NoError(t, err)
		match, err := path.Match(resultDesiredPath, gitPath)
		require.NoError(t, err)
		require.True(t, match)
	})

	t.Run("PathDoesntExist", func(t *testing.T) {
		baseDirectory, cleanup := fstestutility.MakeTempDirectory("test_.*.d")
		defer cleanup()
		_, err := FindEntryDescending(baseDirectory, "someEntry", Regular)
		require.Error(t, err)
	})

	t.Run("InitialDirectoryDoesntExist", func(t *testing.T) {
		unexistingPath := fstestutility.GetAvailableTempPath()
		_, err := FindEntryDescending(unexistingPath, "someEntry", Regular)
		require.Error(t, err)
	})
}

func TestGetFileHash(t *testing.T) {
	data := "some data"
	path1, cleanup := fstestutility.CreateTemporaryFileWithData(data)
	defer cleanup()

	path2, cleanup := fstestutility.CreateTemporaryFileWithData(data)
	defer cleanup()

	hash1 := GetFileHash(path1)
	hash2 := GetFileHash(path2)

	require.True(t, bytes.Equal(hash1, hash2))
}

func TestMakeDirectoryIfDoesntExist(t *testing.T) {
	t.Run("CreateNestedDirectory", func(t *testing.T) {
		// Makes directory path to create
		rootDirectory := fstestutility.GetAvailableTempPath()
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
		fstestutility.AssertNoError(err)
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
		fstestutility.AssertNoError(err)
		defer os.Remove(directory)

		// Executes the test
		err = MakeDirectoryIfDoesntExist(directory)

		// Asserts
		require.NoError(t, err)
	})
}

func TestIsLinkPointsToDestination(t *testing.T) {
	t.Run("LinkPointsToDestination", func(t *testing.T) {
		target := "/dev/null"
		link := fstestutility.GetAvailableTempPath()
		err := os.Symlink(target, link)
		fstestutility.AssertNoError(err)
		defer os.Remove(link)

		// Asserts
		require.True(t, IsLinkPointsToDestination(link, target))
	})

	t.Run("LinkDoesntPointToDestination", func(t *testing.T) {
		target := fstestutility.GetAvailableTempPath()
		link := fstestutility.GetAvailableTempPath()
		err := os.Symlink("/dev/null", link)
		fstestutility.AssertNoError(err)
		defer os.Remove(link)

		// Asserts
		require.False(t, IsLinkPointsToDestination(link, target))
	})
}
