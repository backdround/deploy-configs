package dataconverter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/backdround/deploy-configs/internal/config"
	"github.com/stretchr/testify/require"
)

// //////////////////////////////////////////////////////////
// fakeLogger
type fakeLogger struct{}

func (l fakeLogger) Log(message string)  {}
func (l fakeLogger) Fail(message string) {}

// //////////////////////////////////////////////////////////
// lenExpander
type lenExpander struct{}

// Expand replaces given template to its len
func (e lenExpander) Expand(template string) (string, error) {
	return fmt.Sprint(len(template)), nil
}

// //////////////////////////////////////////////////////////
// errorExpander
type errorExpander struct{}

func (e errorExpander) Expand(template string) (string, error) {
	return "", errors.New("something goes wrong")
}

////////////////////////////////////////////////////////////
// link converting tests

func TestSuccessfulLinkConverting(t *testing.T) {
	// Creates data to convert
	configLinks := map[string]config.Link{
		"l1": {
			TargetPath: "ab",
			LinkPath:   "abcd",
		},
	}

	// Makes conversion
	dataConverter := New(fakeLogger{}, lenExpander{})
	deployLinks, err := dataConverter.RestructureLinks(configLinks)

	// Asserts converted data
	require.NoError(t, err)
	require.Len(t, deployLinks, 1)
	require.Equal(t, "l1", deployLinks[0].Name)
	require.Equal(t, "2", deployLinks[0].TargetPath)
	require.Equal(t, "4", deployLinks[0].LinkPath)
}

func TestFailedLinkConverting(t *testing.T) {
	// Creates data to convert
	configLinks := map[string]config.Link{
		"l1": {
			TargetPath: "ab",
			LinkPath:   "abcd",
		},
	}

	// Fails conversion
	dataConverter := New(fakeLogger{}, errorExpander{})
	deployLinks, err := dataConverter.RestructureLinks(configLinks)

	// Asserts fail
	require.Error(t, err)
	require.Len(t, deployLinks, 0)
}

////////////////////////////////////////////////////////////
// template converting tests

func TestSuccessfulTemplateConverting(t *testing.T) {
	// Creates data to convert
	configTemplates := map[string]config.Template{
		"t1": {
			InputPath:  "ab",
			OutputPath: "abcd",
			Data:       "some data",
		},
	}

	// Makes conversion
	dataConverter := New(fakeLogger{}, lenExpander{})
	deployTemplates, err := dataConverter.RestructureTemplates(configTemplates)

	// Asserts converted data
	require.NoError(t, err)
	require.Len(t, deployTemplates, 1)
	require.Equal(t, "t1", deployTemplates[0].Name)
	require.Equal(t, "2", deployTemplates[0].InputPath)
	require.Equal(t, "4", deployTemplates[0].OutputPath)
	require.Equal(t, "some data", deployTemplates[0].Data)
}

func TestFailedTemplateConverting(t *testing.T) {
	// Creates data to convert
	configTemplates := map[string]config.Template{
		"t1": {
			InputPath:  "ab",
			OutputPath: "abcd",
			Data:       "some data",
		},
	}

	// Fails conversion
	dataConverter := New(fakeLogger{}, errorExpander{})
	deployTemplates, err := dataConverter.RestructureTemplates(configTemplates)

	// Asserts fail
	require.Error(t, err)
	require.Len(t, deployTemplates, 0)
}

////////////////////////////////////////////////////////////
// commands converting tests

func TestSuccessfulCommandsConverting(t *testing.T) {
	// Creates data to convert
	configCommands := map[string]config.Command{
		"c1": {
			InputPath:  "ab",
			OutputPath: "abcd",
			Command:    "do something",
		},
	}

	// Makes conversion
	dataConverter := New(fakeLogger{}, lenExpander{})
	deployCommands, err := dataConverter.RestructureCommands(configCommands)

	// Asserts converted data
	require.NoError(t, err)
	require.Len(t, deployCommands, 1)
	require.Equal(t, "c1", deployCommands[0].Name)
	require.Equal(t, "2", deployCommands[0].InputPath)
	require.Equal(t, "4", deployCommands[0].OutputPath)
	require.Equal(t, "do something", deployCommands[0].CommandTemplate)
}

func TestFailedCommandsConverting(t *testing.T) {
	// Creates data to convert
	configCommands := map[string]config.Command{
		"c1": {
			InputPath:  "ab",
			OutputPath: "abcd",
			Command:    "do something",
		},
	}

	// Fails conversion
	dataConverter := New(fakeLogger{}, errorExpander{})
	deployCommands, err := dataConverter.RestructureCommands(configCommands)

	// Asserts fail
	require.Error(t, err)
	require.Len(t, deployCommands, 0)
}
