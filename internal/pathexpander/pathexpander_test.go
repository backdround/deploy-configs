package pathexpander

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/backdround/deploy-configs/pkg/fstestutility"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type LoggerMock struct {
	mock.Mock
}

func (l *LoggerMock) Log(message string) {
	l.Called(message)
}

func (l *LoggerMock) Warn(message string) {
	l.Called(message)
}

func getLoggerDummy() *LoggerMock {
	logger := &LoggerMock{}
	logger.On("Log", mock.Anything).Maybe()
	logger.On("Warn", mock.Anything).Maybe()
	return logger
}

// containsString returns a mock.matcher that match if argument contains
// a given string for mock.Mock.on function.
func containsString(str string) interface{} {
	return mock.MatchedBy(func(message string) bool {
		return strings.Contains(message, str)
	})
}

func TestExpandWithoutGitRoot(t *testing.T) {
	// Makes templates to expand
	path1 := "{{.GitRoot}}/configs/file1"
	path2 := "{{.Home}}/.config/file1"

	// Creates a logger
	logger := &LoggerMock{}
	logger.On("Warn", containsString("GitRoot")).Once()
	logger.On("Log", containsString("Home")).Once()
	defer logger.AssertExpectations(t)

	// Executes the test
	expander := New(logger, os.TempDir())
	_, err1 := expander.Expand(path1)
	_, err2 := expander.Expand(path2)

	// Asserts expansions
	require.Error(t, err1)
	require.NoError(t, err2)
}

func TestExpandInvalidTemplate(t *testing.T) {
	// Executes the test
	expander := New(getLoggerDummy(), os.TempDir())
	_, err1 := expander.Expand("{{Home}}")

	// Asserts expansions
	require.Error(t, err1)
}

func TestExpandWithGitRoot(t *testing.T) {
	// Creates a test root directory tree
	testRootDirectory, err := os.MkdirTemp("", "go_test_.*.d")
	fstestutility.AssertNoError(err)
	defer os.RemoveAll(testRootDirectory)

	// Creates a git directory
	err = os.MkdirAll(path.Join(testRootDirectory, ".git"), 0755)
	fstestutility.AssertNoError(err)

	// Makes templates to expand
	path1 := "{{.GitRoot}}/configs/file1"
	path2 := "{{.Home}}/.config/file1"

	// Creates a logger
	logger := &LoggerMock{}
	logger.On("Log", containsString("GitRoot")).Once()
	logger.On("Log", containsString("Home")).Once()
	defer logger.AssertExpectations(t)

	// Executes the test
	expander := New(logger, testRootDirectory)
	path1, err1 := expander.Expand(path1)
	_, err2 := expander.Expand(path2)

	// Asserts expansions
	require.NoError(t, err1)
	require.Equal(t, testRootDirectory+"/configs/file1", path1)
	require.NoError(t, err2)
}
