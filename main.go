package main

import (
	"errors"
	"os"

	"github.com/backdround/deploy-configs/config"
	"github.com/backdround/deploy-configs/config/validate"
	"github.com/backdround/deploy-configs/dataconverter"
	"github.com/backdround/deploy-configs/deploy/commands"
	"github.com/backdround/deploy-configs/deploy/links"
	"github.com/backdround/deploy-configs/deploy/templates"
	"github.com/backdround/deploy-configs/logger"
	"github.com/backdround/deploy-configs/pathexpander"
	"github.com/backdround/deploy-configs/pkg/fsutility"
)

func FindConfig(cwd string, names ...string) (configPath string, err error) {
	for _, name := range names {
		types := fsutility.Regular | fsutility.Symlink
		configPath, err = fsutility.FindEntryDescending(cwd, name, types)
		if err == nil {
			return configPath, nil
		}
	}

	return "", errors.New("unable to find config path")
}

func CheckFatalError(err error, l logger.Logger, message string) {
	if err != nil {
		if len(message) != 0 {
			l.Fail(message)
		}
		l.Fail(err.Error())
		os.Exit(1)
	}
}

func main() {
	l := logger.New()

	// Gets config instance
	userInput := os.Args[1:]
	if len(userInput) != 1 {
		l.Fail("expect config instance as argument")
		os.Exit(1)
	}
	configInstance := userInput[0]

	// Gets cwd
	cwd, err := os.Getwd()
	CheckFatalError(err, l, "unable to get current work directory:")

	// Searches config path
	configPath, err := FindConfig(cwd, "deploy-configs.yml",
		"deploy-configs.yaml")
	CheckFatalError(err, l, "")

	// Reads config yaml
	configData, err := os.ReadFile(configPath)
	CheckFatalError(err, l, "unable to get config data")

	// Validates config yaml
	err = validate.Validate(configData)
	CheckFatalError(err, l, "fail to validate config data")

	// Gets config data
	config, err := config.Get(configData, configInstance)
	CheckFatalError(err, l, "fail to get config data")

	// Restructures config to deploy data
	pathExpander := pathexpander.New(l, cwd)
	dataConverter := dataconverter.New(l, pathExpander)

	restructuredLinks, err := dataConverter.RestructureLinks(config.Links)
	CheckFatalError(err, l, "fail to convert config links to deploy links")

	restructuredTemplates, err := dataConverter.RestructureTemplates(
		config.Templates)
	CheckFatalError(err, l,
		"fail to convert config templates to deploy templates")

	restructuredCommands, err := dataConverter.RestructureCommands(
		config.Commands)
	CheckFatalError(err, l,
		"fail to convert config commands to deploy commands")

	// Deploys links
	linkMaker := links.NewLinkMaker(l)
	l.Title("Create links")
	linkMaker.CreateLinks(restructuredLinks)

	// Deploys templates
	templateMaker := templates.NewTemplateMaker(l)
	l.Title("Make templates")
	templateMaker.MakeTemplates(restructuredTemplates)

	// Deploys commands
	commandExecuter := commands.NewCommandExecuter(l)
	l.Title("Execute commands")
	commandExecuter.ExecuteCommands(restructuredCommands)
}
