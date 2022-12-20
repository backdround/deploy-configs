// Package testcase used for integration tests
package testcase

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/backdround/go-fstree/v2"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"

	"github.com/backdround/deploy-configs/internal/realmain"
)

type TestCase struct {
	returnCode    int
	fakeLogger    *FakeLogger
	testDirectory string
}

func RunCase(t *testing.T, fileTreeYaml string, arguments ...string) TestCase {
	c := TestCase{}
	c.fakeLogger = &FakeLogger{}
	c.prepareTestEnvirenment(t, fileTreeYaml)

	c.returnCode = realmain.Main(c.fakeLogger, arguments)

	return c
}

////////////////////////////////////////////////////////////
// Public fucntions

func (c *TestCase) RequireFileTree(t *testing.T, fileTreeYaml string) {
	t.Helper()

	fileTreeYaml = c.prepareYaml(fileTreeYaml)
	difference, err := fstree.CheckOverOSFS(c.testDirectory, fileTreeYaml)
	if err != nil {
		panic(err)
	}
	require.Nil(t, difference)
}

func (c *TestCase) RequireReturnCode(t *testing.T, returnCode int) {
	t.Helper()
	require.Equal(t, returnCode, c.returnCode)
}

func (c *TestCase) RequireFailMessage(t *testing.T, message string) {
	t.Helper()
	message = c.prepareOutput(message)
	c.fakeLogger.RequireFailContains(t, message)
}

func (c *TestCase) RequireSuccessMessage(t *testing.T, message string) {
	t.Helper()
	message = c.prepareOutput(message)
	c.fakeLogger.RequireSuccessContains(t, message)
}

func (c *TestCase) RequireLogMessage(t *testing.T, message string) {
	t.Helper()
	message = c.prepareOutput(message)
	c.fakeLogger.RequireLogContains(t, message)
}

////////////////////////////////////////////////////////////
// Private fucntions

func (c *TestCase) prepareYaml(yaml string) string {
	yaml = dedent.Dedent(yaml)
	yaml = strings.ReplaceAll(yaml, "\t", "  ")
	return yaml
}

func (c *TestCase) prepareOutput(output string) string {
	// Converts code pseudo yaml to real yaml
	output = dedent.Dedent(output)
	output = strings.ReplaceAll(output, "\t", "  ")

	// Removes first and last empty line
	outsideEmptyLines := regexp.MustCompile("^\n|\n$")
	output = outsideEmptyLines.ReplaceAllString(output, "")

	// Substitutes root
	output = strings.ReplaceAll(output, "{Root}", c.testDirectory)

	return output
}

func (c *TestCase) prepareTestEnvirenment(t *testing.T, fileTreeYaml string) {
	// Saves current work directory
	oldWorkDirectory, err := os.Getwd()
	assertNoError(err)

	// Gets test directory
	c.testDirectory, err = os.MkdirTemp("", "go-test-deploy-configs-*.d")
	assertNoError(err)

	// Cd to test directory
	err = os.Chdir(c.testDirectory)
	assertNoError(err)

	// Creates filetree structure
	fileTreeYaml = c.prepareYaml(fileTreeYaml)
	err = fstree.MakeOverOSFS(c.testDirectory, fileTreeYaml)
	assertNoError(err)

	// Restores work directory and removes test directory after test complition
	t.Cleanup(func() {
		// Cd to old work directory
		err := os.Chdir(oldWorkDirectory)
		assertNoError(err)

		// Removes test directory
		err = os.RemoveAll(c.testDirectory)
		assertNoError(err)
	})
}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}
