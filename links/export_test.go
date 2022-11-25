package links

import (
	"github.com/stretchr/testify/mock"

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

// containsString returns a mock.matcher that match if argument contains
// a given string for mock.Mock.on function.
func containsString(str string) interface{} {
	return mock.MatchedBy(func(message string) bool {
		return strings.Contains(message, str)
	})
}
