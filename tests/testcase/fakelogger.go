package testcase

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type FakeLogger struct {
	titles    []string
	successes []string
	warns     []string
	fails     []string
	logs      []string
}

////////////////////////////////////////////////////////////
// Implement logger

func (l *FakeLogger) Title(title string) {
	l.titles = append(l.titles, title)
}

func (l *FakeLogger) Success(message string) {
	l.successes = append(l.successes, message)
}

func (l *FakeLogger) Warn(message string) {
	l.warns = append(l.warns, message)
}

func (l *FakeLogger) Fail(message string) {
	l.fails = append(l.fails, message)
}

func (l *FakeLogger) Log(message string) {
	l.logs = append(l.logs, message)
}

////////////////////////////////////////////////////////////
// Test utility members

func (l *FakeLogger) requireContains(t *testing.T, messages []string,
	substring string) {
	t.Helper()

	for _, message := range messages {
		if strings.Contains(message, substring) {
			return
		}
	}

	errMessage := fmt.Sprintf("%#v doesn't fuzzy contain %#v", messages,
		substring)
	require.FailNowf(t, errMessage, "")
}

func (l *FakeLogger) RequireSuccessContains(t *testing.T, message string) {
	t.Helper()
	l.requireContains(t, l.successes, message)
}

func (l *FakeLogger) RequireWarnContains(t *testing.T, message string) {
	t.Helper()
	l.requireContains(t, l.warns, message)
}

func (l *FakeLogger) RequireFailContains(t *testing.T, message string) {
	t.Helper()
	l.requireContains(t, l.fails, message)
}

func (l *FakeLogger) RequireLogContains(t *testing.T, message string) {
	t.Helper()
	l.requireContains(t, l.logs, message)
}
