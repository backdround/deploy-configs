package links

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"

	"os"
	"path"

	"github.com/backdround/deploy-configs/pkg/fstestutility"
	"github.com/backdround/deploy-configs/pkg/fsutility"
)

////////////////////////////////////////////////////////////
// makeLink

func TestSuccessfulMakeLink(t *testing.T) {
	t.Run("LinksDoesntExist", func(t *testing.T) {
		// Creates a target file
		targetFile := "target.*.txt"
		linkPath := targetFile + ".link"
		cleanup := fstestutility.CreateTemporaryFiles(&targetFile)
		defer cleanup()
		defer os.Remove(linkPath)

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Success", containsString("target.")).Once()

		// Executes the test
		link := Link{
			Name:       "test-link",
			TargetPath: targetFile,
			LinkPath:   linkPath,
		}
		NewLinkMaker(loggerMock).makeLink(link)

		// Asserts the created symlink
		require.True(t, fsutility.IsLinkPointsToDestination(link.LinkPath,
			targetFile))
	})

	t.Run("LinkDirectoryDoesntExist", func(t *testing.T) {
		// Creates a target file
		targetFile := "target.*.txt"
		cleanup := fstestutility.CreateTemporaryFiles(&targetFile)
		defer cleanup()

		// Makes a link path inside a notexisting directory
		linkDirectory := fstestutility.GetAvailableTempPath()
		linkPath := path.Join(linkDirectory, "link")
		defer os.RemoveAll(linkDirectory)

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Success", containsString("target.")).Once()

		// Executes the test
		link := Link{
			Name:       "test-link",
			TargetPath: targetFile,
			LinkPath:   linkPath,
		}

		NewLinkMaker(loggerMock).makeLink(link)

		// Asserts the created symlink
		require.True(t, fsutility.IsLinkPointsToDestination(link.LinkPath,
			targetFile))
	})

	t.Run("LinksAreIncorrect", func(t *testing.T) {
		// Creates a target file
		targetFile := "target.*.txt"
		cleanup := fstestutility.CreateTemporaryFiles(&targetFile)
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
			Name:       "test-link",
			TargetPath: targetFile,
			LinkPath:   linkPath,
		}
		NewLinkMaker(loggerMock).makeLink(link)

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
		cleanup := fstestutility.CreateTemporaryFiles(&targetFile, &linkPath)
		defer cleanup()

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Fail", containsString("target.")).Once()

		// Executes the test
		link := Link{
			Name:       "test-link",
			TargetPath: targetFile,
			LinkPath:   linkPath,
		}
		NewLinkMaker(loggerMock).makeLink(link)

		// Asserts that the file on the link place wasn't deleted
		linkType := fsutility.GetPathType(link.LinkPath)
		require.Equal(t, fsutility.Regular.String(), linkType.String())
	})

	t.Run("LinkPathIsADirectory", func(t *testing.T) {
		// Creates target file
		targetFile := "target.*.txt"
		cleanup := fstestutility.CreateTemporaryFiles(&targetFile)
		defer cleanup()

		// Creates directory by link path
		linkPath := fstestutility.GetAvailableTempPath()
		fstestutility.AssertNoError(os.Mkdir(linkPath, 0755))
		defer os.Remove(linkPath)

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Fail", containsString("target.")).Once()

		// Executes the test
		link := Link{
			Name:       "test-link",
			TargetPath: targetFile,
			LinkPath:   linkPath,
		}
		NewLinkMaker(loggerMock).makeLink(link)

		// Asserts that the file on the link place wasn't deleted
		linkType := fsutility.GetPathType(link.LinkPath)
		require.Equal(t, fsutility.Directory.String(), linkType.String())
	})

	t.Run("TargetPathsDontExist", func(t *testing.T) {
		// Creates test paths
		targetFile := fstestutility.GetAvailableTempPath()
		linkPath := fstestutility.GetAvailableTempPath()

		// Sets up the mock
		loggerMock := new(LoggerMock)
		defer loggerMock.AssertExpectations(t)
		loggerMock.On("Fail", mock.Anything).Once()

		// Executes the test
		link := Link{
			Name:       "test-link",
			TargetPath: targetFile,
			LinkPath:   linkPath,
		}
		NewLinkMaker(loggerMock).makeLink(link)

		// Asserts that files wasn't created
		linkType := fsutility.GetPathType(linkPath)
		require.Equal(t, fsutility.Notexisting.String(), linkType.String())
		targetType := fsutility.GetPathType(targetFile)
		require.Equal(t, fsutility.Notexisting.String(), targetType.String())
	})
}

func TestSkippedMakeLink(t *testing.T) {
	targetFile := "/dev/null"
	linkPath := fstestutility.GetAvailableTempPath()

	// Creates link
	err := os.Symlink(targetFile, linkPath)
	fstestutility.AssertNoError(err)
	defer os.Remove(linkPath)

	// Sets up the mock
	loggerMock := new(LoggerMock)
	defer loggerMock.AssertExpectations(t)
	loggerMock.On("Log", mock.Anything).Once()

	// Executes the test
	link := Link{
		Name:       "test-link",
		TargetPath: targetFile,
		LinkPath:   linkPath,
	}
	NewLinkMaker(loggerMock).makeLink(link)

	// Asserts that the link exists
	require.True(t, fsutility.IsLinkPointsToDestination(linkPath, targetFile))
}

//////////////////////////////////////////////////////////
// Links

func TestLinks(t *testing.T) {
	t.Run("SeveralLinks", func(t *testing.T) {
		// Creates data
		targetFile := "/dev/null"
		link1Path := fstestutility.GetAvailableTempPath()
		link2Path := fstestutility.GetAvailableTempPath()
		defer os.Remove(link1Path)
		defer os.Remove(link2Path)

		links := []Link{{
			Name:       "link1",
			TargetPath: targetFile,
			LinkPath:   link1Path,
		}, {
			Name:       "link1",
			TargetPath: targetFile,
			LinkPath:   link2Path,
		},
		}

		// Executes the test
		NewLinkMaker(getLoggerDummy()).CreateLinks(links)

		// Asserts that the links are correct
		require.True(t, fsutility.IsLinkPointsToDestination(link1Path,
			targetFile))
		require.True(t, fsutility.IsLinkPointsToDestination(link2Path,
			targetFile))
	})

	t.Run("TargetPathIsADirectory", func(t *testing.T) {
		// Creates a target directory
		targetDirectory, err := os.MkdirTemp("", "target.*.d")
		fstestutility.AssertNoError(err)
		defer os.Remove(targetDirectory)

		target1Path := path.Join(targetDirectory, "target1")
		_, err = os.Create(target1Path)
		fstestutility.AssertNoError(err)
		defer os.Remove(target1Path)

		target2Path := path.Join(targetDirectory, "target2")
		_, err = os.Create(target2Path)
		fstestutility.AssertNoError(err)
		defer os.Remove(target2Path)

		// Gets a link path
		linkPath := fstestutility.GetAvailableTempPath()
		defer os.RemoveAll(linkPath)

		// Makes data
		links := []Link{{
			Name:       "directory",
			TargetPath: targetDirectory,
			LinkPath:   linkPath,
		},
		}

		// Executes the test
		NewLinkMaker(getLoggerDummy()).CreateLinks(links)

		// Asserts that the links are valid
		expectedLink1Path := path.Join(linkPath, "target1")
		expectedLink2Path := path.Join(linkPath, "target2")
		require.True(t,
			fsutility.IsLinkPointsToDestination(expectedLink1Path, target1Path))
		require.True(t,
			fsutility.IsLinkPointsToDestination(expectedLink2Path, target2Path))
	})
}
