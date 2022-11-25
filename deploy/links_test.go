package deploy

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"os"
	"path"
)

////////////////////////////////////////////////////////////
// isLinkPointsToDestination

func TestIsLinkPointsToDestination(t *testing.T) {
	t.Run("LinkPointsToDestination", func(t *testing.T) {
		target := "/dev/null"
		link := getNotExistingPath()
		err := os.Symlink(target, link)
		assertNoError(err)
		defer os.Remove(link)

		// Asserts
		require.True(t, isLinkPointsToDestination(link, target))
	})

	t.Run("LinkDoesntPointToDestination", func(t *testing.T) {
		target := getNotExistingPath()
		link := getNotExistingPath()
		err := os.Symlink("/dev/null", link)
		assertNoError(err)
		defer os.Remove(link)

		// Asserts
		require.False(t, isLinkPointsToDestination(link, target))
	})
}

////////////////////////////////////////////////////////////
// makeLink

func TestSuccessfulMakeLink(t *testing.T) {
	t.Run("LinksDoesntExists", func(t *testing.T) {
		// Creates a target file
		targetFile := "target.*.txt"
		linkPath := targetFile + ".link"
		cleanup := createTemporaryFiles(&targetFile)
		defer cleanup()
		defer os.Remove(linkPath)

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Success", containsString("target.")).Once()

		// Executes the test
		link := Link{
			Target:   targetFile,
			LinkPath: linkPath,
		}
		NewLinkMaker(loggerMock).makeLink("test-link", link)

		// Asserts the created symlink
		require.True(t, isLinkPointsToDestination(link.LinkPath, targetFile))
	})

	t.Run("LinkDirectoryDoesntExists", func(t *testing.T) {
		// Creates a target file
		targetFile := "target.*.txt"
		cleanup := createTemporaryFiles(&targetFile)
		defer cleanup()

		// Makes a link path inside a notexisting directory
		linkDirectory := getNotExistingPath()
		linkPath := path.Join(linkDirectory, "link")
		defer os.RemoveAll(linkDirectory)

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Success", containsString("target.")).Once()

		// Executes the test
		link := Link{
			Target:   targetFile,
			LinkPath: linkPath,
		}

		NewLinkMaker(loggerMock).makeLink("test-link", link)

		// Asserts the created symlink
		require.True(t, isLinkPointsToDestination(link.LinkPath, targetFile))
	})

	t.Run("LinksAreIncorrect", func(t *testing.T) {
		// Creates a target file
		targetFile := "target.*.txt"
		cleanup := createTemporaryFiles(&targetFile)
		defer cleanup()

		// Creates a link file
		linkPath := targetFile + ".link"
		err := os.Symlink("/dev/null", linkPath)
		require.NoError(t, err)
		defer os.Remove(linkPath)

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Success", containsString("target.")).Once()

		// Executes the test
		link := Link{
			Target:   targetFile,
			LinkPath: linkPath,
		}
		NewLinkMaker(loggerMock).makeLink("test-link", link)

		// Asserts the created symlink
		require.True(t, isLinkPointsToDestination(linkPath, targetFile))
	})
}

func TestFailedMakeLink(t *testing.T) {
	t.Run("LinksAreExistingFiles", func(t *testing.T) {
		// Creates test files
		targetFile := "target.*.txt"
		linkPath := targetFile + ".link"
		cleanup := createTemporaryFiles(&targetFile, &linkPath)
		defer cleanup()

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Fail", containsString("target.")).Once()

		// Executes the test
		link := Link{
			Target:   targetFile,
			LinkPath: linkPath,
		}
		NewLinkMaker(loggerMock).makeLink("test-link", link)

		// Asserts that the file on the link place wasn't deleted
		linkType := getFileType(link.LinkPath)
		require.Equal(t, regular.String(), linkType.String())
	})

	t.Run("TargetsDontExist", func(t *testing.T) {
		// Creates test paths
		targetFile := getNotExistingPath()
		linkPath := getNotExistingPath()

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Fail", mock.Anything).Once()

		// Executes the test
		link := Link{
			Target:   targetFile,
			LinkPath: linkPath,
		}
		NewLinkMaker(loggerMock).makeLink("test-link", link)

		// Asserts that files wasn't created
		linkType := getFileType(linkPath)
		require.Equal(t, notexisting.String(), linkType.String())
		targetType := getFileType(targetFile)
		require.Equal(t, notexisting.String(), targetType.String())
	})
}

func TestSkippedMakeLink(t *testing.T) {
	targetFile := "/dev/null"
	linkPath := getNotExistingPath()

	// Creates link
	err := os.Symlink(targetFile, linkPath)
	assertNoError(err)
	defer os.Remove(linkPath)

	// Sets up the mock
	loggerMock := new(LoggerMock)
	defer loggerMock.AssertExpectations(t)
	loggerMock.On("Log", mock.Anything).Once()

	// Executes the test
	link := Link{
		Target:   targetFile,
		LinkPath: linkPath,
	}
	NewLinkMaker(loggerMock).makeLink("test-link", link)

	// Asserts that the link exists
	require.True(t, isLinkPointsToDestination(linkPath, targetFile))
}

//////////////////////////////////////////////////////////
// Links

func TestLinks(t *testing.T) {
	t.Run("SeveralLinks", func(t *testing.T) {
		// Creates data
		targetFile := "/dev/null"
		link1Path := getNotExistingPath()
		link2Path := getNotExistingPath()
		defer os.Remove(link1Path)
		defer os.Remove(link2Path)

		links := map[string]Link{
			"link1": {
				Target:   targetFile,
				LinkPath: link1Path,
			},
			"link2": {
				Target:   targetFile,
				LinkPath: link2Path,
			},
		}

		// Executes the test
		NewLinkMaker(getLoggerDummy()).Links(links)

		// Asserts that the links are correct
		require.True(t, isLinkPointsToDestination(link1Path, targetFile))
		require.True(t, isLinkPointsToDestination(link2Path, targetFile))
	})

	t.Run("TargetIsADirectory", func(t *testing.T) {
		// Creates a target directory
		targetDirectory, err := os.MkdirTemp("", "target.*.d")
		assertNoError(err)
		defer os.Remove(targetDirectory)

		target1Path := path.Join(targetDirectory, "target1")
		_, err = os.Create(target1Path)
		assertNoError(err)
		defer os.Remove(target1Path)

		target2Path := path.Join(targetDirectory, "target2")
		_, err = os.Create(target2Path)
		assertNoError(err)
		defer os.Remove(target2Path)

		// Gets a link path
		linkPath := getNotExistingPath()
		defer os.RemoveAll(linkPath)

		// Makes data
		links := map[string]Link{
			"directory": {
				Target:   targetDirectory,
				LinkPath: linkPath,
			},
		}

		// Executes the test
		NewLinkMaker(getLoggerDummy()).Links(links)

		// Asserts that the links are valid
		expectedLink1Path := path.Join(linkPath, "target1")
		expectedLink2Path := path.Join(linkPath, "target2")
		require.True(t, isLinkPointsToDestination(expectedLink1Path, target1Path))
		require.True(t, isLinkPointsToDestination(expectedLink2Path, target2Path))
	})
}
