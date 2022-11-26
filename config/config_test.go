package config

import (
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"

	"fmt"
	"strings"
)

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

func assertNoTab(data string) {
	if strings.Contains(data, "\t") {
		message := fmt.Sprintf("data contains tab: \n%v", data)
		panic(message)
	}
}

func TestNonExistenInstanceConfig(t *testing.T) {
	data := dedent.Dedent(`
	  instances:
	    instance1:
	      commands:
	    instance2:
	      links:
	`)
	assertNoTab(data)

	config, err := Get([]byte(data), "not-existen")
	require.Nil(t, config)
	require.Error(t, err)
}

func TestChooseInstanceConfig(t *testing.T) {
	data := dedent.Dedent(`
	  instances:
	    instance1:
	      commands:
	    instance2:
	      links:
	`)
	assertNoTab(data)

	config, err := Get([]byte(data), "instance1")
	require.NotNil(t, config)
	require.NoError(t, err)
}

func TestLinksConfig(t *testing.T) {
	data := dedent.Dedent(`
	  instances:
	    instance1:
	      links:
	        link1: ["./file1.txt", "./link1"]
	        link2: ["./file2.txt", "./link2"]
	      commands:
	      templates:
	`)
	assertNoTab(data)

	config, err := Get([]byte(data), "instance1")
	require.NotNil(t, config)
	require.NoError(t, err)

	require.Contains(t, config.Links, "link1")
	link1 := config.Links["link1"]
	require.Equal(t, "./file1.txt", link1.TargetPath)
	require.Equal(t, "./link1", link1.LinkPath)

	require.Contains(t, config.Links, "link2")
	link2 := config.Links["link2"]
	require.Equal(t, "./file2.txt", link2.TargetPath)
	require.Equal(t, "./link2", link2.LinkPath)
}

func TestCommandsConfig(t *testing.T) {
	data := dedent.Dedent(`
	  instances:
	    instance1:
	      commands:
	        command1:
	          input: "./file.txt"
	          output: "~/file.txt"
	          command: "seq 3"
	      links:
	      templates:
	`)
	assertNoTab(data)

	config, err := Get([]byte(data), "instance1")
	require.NotNil(t, config)
	require.NoError(t, err)

	require.Contains(t, config.Commands, "command1")

	command := config.Commands["command1"]
	require.Equal(t, "./file.txt", command.InputPath)
	require.Equal(t, "~/file.txt", command.OutputPath)
	require.Equal(t, "seq 3", command.Command)
}

func TestTemplatesConfig(t *testing.T) {
	data := dedent.Dedent(`
	  instances:
	    instance1:
	      templates:
	        template1:
	          input: "./file.txt"
	          output: "~/file.txt"
	          data:
	            variable1: "value1"
	            variable2: "value2"
	`)
	assertNoTab(data)

	config, err := Get([]byte(data), "instance1")
	require.NotNil(t, config)
	require.NoError(t, err)

	require.Contains(t, config.Templates, "template1")

	template := config.Templates["template1"]
	require.Equal(t, "./file.txt", template.InputPath)
	require.Equal(t, "~/file.txt", template.OutputPath)

	require.Contains(t, template.Data, "variable1")
	require.Contains(t, template.Data, "variable2")
	templateData := template.Data.(map[string]interface{})
	require.Equal(t, "value1", templateData["variable1"].(string))
	require.Equal(t, "value2", templateData["variable2"].(string))
}
