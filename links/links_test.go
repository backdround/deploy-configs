package links

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"os"
	"path"

	"github.com/backdround/deploy-configs/pkg/fsutility"
)

////////////////////////////////////////////////////////////
// makeLink

func TestSuccessfulMakeLink(t *testing.T) {
	t.Run("LinksDoesntExist", func(t *testing.T) {
		// Creates a target file
		targetFile := "target.*.txt"
		linkPath := targetFile + ".link"
		cleanup := fsutility.CreateTemporaryFiles(&targetFile)
		defer cleanup()
		defer os.Remove(linkPath)

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Success", containsString("target.")).Once()

		// Executes the test
		link := Link{
			TargetPath: targetFile,
			LinkPath:   linkPath,
		}
		NewLinkMaker(loggerMock).makeLink("test-link", link)

		// Asserts the created symlink
		require.True(t, fsutility.IsLinkPointsToDestination(link.LinkPath,
			targetFile))
	})

	t.Run("LinkDirectoryDoesntExist", func(t *testing.T) {
		// Creates a target file
		targetFile := "target.*.txt"
		cleanup := fsutility.CreateTemporaryFiles(&targetFile)
		defer cleanup()

		// Makes a link path inside a notexisting directory
		linkDirectory := fsutility.GetNotExistingPath()
		linkPath := path.Join(linkDirectory, "link")
		defer os.RemoveAll(linkDirectory)

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Success", containsString("target.")).Once()

		// Executes the test
		link := Link{
			TargetPath: targetFile,
			LinkPath:   linkPath,
		}

		NewLinkMaker(loggerMock).makeLink("test-link", link)

		// Asserts the created symlink
		require.True(t, fsutility.IsLinkPointsToDestination(link.LinkPath,
			targetFile))
	})

	t.Run("LinksAreIncorrect", func(t *testing.T) {
		// Creates a target file
		targetFile := "target.*.txt"
		cleanup := fsutility.CreateTemporaryFiles(&targetFile)
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
			TargetPath: targetFile,
			LinkPath:   linkPath,
		}
		NewLinkMaker(loggerMock).makeLink("test-link", link)

		// Asserts the created symlink
		require.True(t, fsutility.IsLinkPointsToDestination(linkPath,
			targetFile))
	})
}

func TestFailedMakeLink(t *testing.T) {
	t.Run("LinksAreExistingFiles", func(t *testing.T) {
		// Creates test files
		targetFile := "target.*.txt"
		linkPath := targetFile + ".link"
		cleanup := fsutility.CreateTemporaryFiles(&targetFile, &linkPath)
		defer cleanup()

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Fail", containsString("target.")).Once()

		// Executes the test
		link := Link{
			TargetPath: targetFile,
			LinkPath:   linkPath,
		}
		NewLinkMaker(loggerMock).makeLink("test-link", link)

		// Asserts that the file on the link place wasn't deleted
		linkType := fsutility.GetFileType(link.LinkPath)
		require.Equal(t, fsutility.Regular.String(), linkType.String())
	})

	t.Run("TargetPathsDontExist", func(t *testing.T) {
		// Creates test paths
		targetFile := fsutility.GetNotExistingPath()
		linkPath := fsutility.GetNotExistingPath()

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Fail", mock.Anything).Once()

		// Executes the test
		link := Link{
			TargetPath: targetFile,
			LinkPath:   linkPath,
		}
		NewLinkMaker(loggerMock).makeLink("test-link", link)

		// Asserts that files wasn't created
		linkType := fsutility.GetFileType(linkPath)
		require.Equal(t, fsutility.Notexisting.String(), linkType.String())
		targetType := fsutility.GetFileType(targetFile)
		require.Equal(t, fsutility.Notexisting.String(), targetType.String())
	})
}

func TestSkippedMakeLink(t *testing.T) {
	targetFile := "/dev/null"
	linkPath := fsutility.GetNotExistingPath()

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
		TargetPath: targetFile,
		LinkPath:   linkPath,
	}
	NewLinkMaker(loggerMock).makeLink("test-link", link)

	// Asserts that the link exists
	require.True(t, fsutility.IsLinkPointsToDestination(linkPath, targetFile))
}

//////////////////////////////////////////////////////////
// Links

func TestLinks(t *testing.T) {
	t.Run("SeveralLinks", func(t *testing.T) {
		// Creates data
		targetFile := "/dev/null"
		link1Path := fsutility.GetNotExistingPath()
		link2Path := fsutility.GetNotExistingPath()
		defer os.Remove(link1Path)
		defer os.Remove(link2Path)

		links := map[string]Link{
			"link1": {
				TargetPath: targetFile,
				LinkPath:   link1Path,
			},
			"link2": {
				TargetPath: targetFile,
				LinkPath:   link2Path,
			},
		}

		// Executes the test
		NewLinkMaker(getLoggerDummy()).Links(links)

		// Asserts that the links are correct
		require.True(t, fsutility.IsLinkPointsToDestination(link1Path,
			targetFile))
		require.True(t, fsutility.IsLinkPointsToDestination(link2Path,
			targetFile))
	})

	t.Run("TargetPathIsADirectory", func(t *testing.T) {
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
		linkPath := fsutility.GetNotExistingPath()
		defer os.RemoveAll(linkPath)

		// Makes data
		links := map[string]Link{
			"directory": {
				TargetPath: targetDirectory,
				LinkPath:   linkPath,
			},
		}

		// Executes the test
		NewLinkMaker(getLoggerDummy()).Links(links)

		// Asserts that the links are valid
		expectedLink1Path := path.Join(linkPath, "target1")
		expectedLink2Path := path.Join(linkPath, "target2")
		require.True(t,
			fsutility.IsLinkPointsToDestination(expectedLink1Path, target1Path))
		require.True(t,
			fsutility.IsLinkPointsToDestination(expectedLink2Path, target2Path))
	})
}
