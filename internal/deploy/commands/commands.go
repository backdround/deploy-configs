// commands describes commandExecuter which receives a bunch
// of commands, that create files from inputFiles, executes these
// and logs all outcomes.
package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"text/template"

	"github.com/backdround/deploy-configs/pkg/fsutility"
	"github.com/backdround/go-indent"
)

type commandExecuter struct {
	logger Logger
}

func NewCommandExecuter(logger Logger) *commandExecuter {
	return &commandExecuter{
		logger: logger,
	}
}

func getDescription(command Command) string {
	return fmt.Sprintf("input: %q\noutput: %q\ncommand: %q",
		command.InputPath, command.OutputPath, command.CommandTemplate)
}

func shift(message string, count int) string {
	return indent.Indent(message, "  ", count)
}

func (e commandExecuter) logFail(command Command, reason string) {
	description := shift(getDescription(command), 1)
	errorMessage := shift("error: "+reason, 2)

	message := fmt.Sprintf("Unable to execute %q command:\n%v\n%v\n",
		command.Name, description, errorMessage)
	e.logger.Fail(message)
}

func (e commandExecuter) logSuccess(command Command) {
	message := fmt.Sprintf("Command %q is executed:\n%v",
		command.Name, shift(getDescription(command), 1))
	e.logger.Success(message)
}

func (e commandExecuter) logSkip(command Command) {
	message := fmt.Sprintf("Command %q is skipped", command.Name)
	e.logger.Log(message)
}

// executeCommand expands command template, executes command,
// checks that the OutputPath is created and logs all outcomes.
func (e commandExecuter) executeCommand(c Command) (success bool) {
	// Checks that the input file exists
	inputPathType := fsutility.GetPathType(c.InputPath)
	if inputPathType == fsutility.Notexisting {
		e.logFail(c, "input file doesn't exist")
		return false
	}

	// Creates the output directory if it's needed
	outputDirectory := path.Dir(c.OutputPath)
	err := fsutility.MakeDirectoryIfDoesntExist(outputDirectory)
	if err != nil {
		e.logFail(c, err.Error())
		return false
	}

	// Saves a hash of the old output file (if it exists)
	oldOutputFileHash := fsutility.GetFileHash(c.OutputPath)

	// Removes the old output file if it exists
	outputPathType := fsutility.GetPathType(c.OutputPath)
	if outputPathType != fsutility.Notexisting {
		err := os.Remove(c.OutputPath)
		if err != nil {
			message := fmt.Sprintf("unable to replace output path:\n%v",
				shift(err.Error(), 1))
			e.logFail(c, message)
			return false
		}
	}

	// Gets command template
	commandTemplate := template.New(c.Name).Option("missingkey=error")
	commandTemplate, err = commandTemplate.Parse(c.CommandTemplate)
	if err != nil {
		e.logFail(c, err.Error())
		return false
	}

	// Gets expanded command
	expandData := map[string]string{
		"Input":  c.InputPath,
		"Output": c.OutputPath,
	}
	expandedCommand := bytes.NewBuffer([]byte{})
	err = commandTemplate.Execute(expandedCommand, expandData)
	if err != nil {
		e.logFail(c, err.Error())
		return false
	}

	// Executes the expanded command
	cmd := exec.Command("sh", "-c", expandedCommand.String())
	cmdOutput, err := cmd.Output()
	if err != nil {
		os.Remove(c.OutputPath)
		e.logFail(c, err.Error())
		return false
	}

	// Checks that the command created the output file
	outputPathType = fsutility.GetPathType(c.OutputPath)
	if outputPathType != fsutility.Regular {
		message := fmt.Sprintf("command didn't create file. output:\n%v",
			string(cmdOutput))
		e.logFail(c, message)
		return false
	}

	// Checks that output file is changed
	newOutputFileHash := fsutility.GetFileHash(c.OutputPath)
	if bytes.Equal(oldOutputFileHash, newOutputFileHash) {
		e.logSkip(c)
		return true
	}

	e.logSuccess(c)
	return true
}

// ExecuteCommands expands and executes given commands
func (e commandExecuter) ExecuteCommands(commands []Command) (success bool) {
	// Sorts commands by name
	sort.Slice(commands, func(i int, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	// Executes commands
	success = true
	for _, command := range commands {
		success = success && e.executeCommand(command)
	}
	return success
}
