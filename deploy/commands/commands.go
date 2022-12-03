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
	"text/template"

	"github.com/backdround/deploy-configs/pkg/fsutility"
)

type commandExecuter struct {
	logger Logger
}

func NewCommandExecuter(logger Logger) *commandExecuter {
	return &commandExecuter{
		logger: logger,
	}
}

func (e commandExecuter) getDescription(command Command) string {
	return fmt.Sprintf("command %q [%q, %q]",
		command.CommandTemplate, command.InputPath, command.OutputPath)
}

func (e commandExecuter) logFail(command Command, reason string) {
	message := fmt.Sprintf("Unable to expand %q link:\n\t%v\n\t\t%v",
		command.Name, e.getDescription(command), reason)
	e.logger.Fail(message)
}

func (e commandExecuter) logSuccess(command Command) {
	message := fmt.Sprintf("Command %q is executed:\n\t%v",
		command.Name, e.getDescription(command))
	e.logger.Success(message)
}

func (e commandExecuter) logSkip(command Command) {
	message := fmt.Sprintf("Command %q is skipped:\n\t%v",
		command.Name, e.getDescription(command))
	e.logger.Log(message)
}

// executeCommand expands command template, executes command,
// checks that the OutputPath is created and logs all outcomes.
func (e commandExecuter) executeCommand(c Command) {
	// Checks that the input file exists
	inputPathType := fsutility.GetPathType(c.InputPath)
	if inputPathType == fsutility.Notexisting {
		e.logFail(c, "input file doesn't exist")
		return
	}

	// Saves a hash of the old output file (if it exists)
	oldOutputFileHash := fsutility.GetFileHash(c.OutputPath)

	// Removes the old output file if it exists
	outputPathType := fsutility.GetPathType(c.OutputPath)
	if outputPathType == fsutility.Regular {
		err := os.Remove(c.OutputPath)
		if err != nil {
			message := fmt.Sprintf("unable to remove output file:\n\t%v",
				err.Error())
			e.logFail(c, message)
			return
		}
	}

	// Creates the output directory if it's needed
	outputDirectory := path.Dir(c.OutputPath)
	err := fsutility.MakeDirectoryIfDoesntExist(outputDirectory)
	if err != nil {
		e.logFail(c, err.Error())
		return
	}

	// Gets command template
	commandTemplate := template.New(c.Name).Option("missingkey=error")
	commandTemplate, err = commandTemplate.Parse(c.CommandTemplate)
	if err != nil {
		e.logFail(c, err.Error())
		return
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
		return
	}

	// Executes the expanded command
	cmd := exec.Command("sh", "-c", expandedCommand.String())
	cmdOutput, err := cmd.Output()
	if err != nil {
		os.Remove(c.OutputPath)
		e.logFail(c, err.Error())
		return
	}

	// Checks that the command created the output file
	outputPathType = fsutility.GetPathType(c.OutputPath)
	if outputPathType != fsutility.Regular {
		message := fmt.Sprintf("command didn't create file. output:\n%v",
			cmdOutput)
		e.logFail(c, message)
		return
	}

	// Checks that output file is changed
	newOutputFileHash := fsutility.GetFileHash(c.OutputPath)
	if bytes.Equal(oldOutputFileHash, newOutputFileHash) {
		e.logSkip(c)
		return
	}

	e.logSuccess(c)
}

// ExecuteCommands expands and executes given commands
func (e commandExecuter) ExecuteCommands(commands []Command) {
	for _, command := range commands {
		e.executeCommand(command)
	}
}
