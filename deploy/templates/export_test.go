package templates

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

////////////////////////////////////////////////////////////
// Utility functions

// containsString returns a mock.matcher that match if argument contains
// a given string for mock.Mock.on function.
func containsString(str string) interface{} {
	return mock.MatchedBy(func(message string) bool {
		return strings.Contains(message, str)
	})
}
