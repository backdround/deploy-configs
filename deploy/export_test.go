package deploy

import (
	"github.com/stretchr/testify/mock"

	"os"
	"strings"
)

////////////////////////////////////////////////////////////
// LoggerMock
type LoggerMock struct {
	mock.Mock
}

func (l *LoggerMock) Success(message string) {
	l.Called(message)
}

func (l *LoggerMock) Fail(message string) {
	l.Called(message)
}

func (l *LoggerMock) Log(message string) {
	l.Called(message)
}

func getLoggerDummy() Logger {
	logger := &LoggerMock{}
	logger.On("Success", mock.Anything).Maybe()
	logger.On("Fail", mock.Anything).Maybe()
	logger.On("Log", mock.Anything).Maybe()
	return logger
}

////////////////////////////////////////////////////////////
// Utility functions

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func getNotExistingPath() string {
	file, err := os.CreateTemp("", "go_test.*.txt")
	path := file.Name()
	assertNoError(err)
	assertNoError(os.Remove(path))

	return path
}

// containsString returns a mock.matcher for mock.Mock.on function.
func containsString(str string) interface{} {
	return mock.MatchedBy(func(message string) bool {
		return strings.Contains(message, str)
	})
}

// createTemporaryFiles creates files by the patterns (* for random substitution).
// It changes the patterns parameters to the paths. It returns a cleanup function.
func createTemporaryFiles(patterns ...*string) (cleanup func()) {
	filesToRemove := []string{}

	for _, pattern := range patterns {
		file, err := os.CreateTemp("", *pattern)
		assertNoError(err)
		file.Close()
		filesToRemove = append(filesToRemove, file.Name())
		*pattern = file.Name()
	}

	removeAllFiles := func() {
		for _, path := range filesToRemove {
			err := os.Remove(path)
			assertNoError(err)
		}
	}

	return removeAllFiles
}
