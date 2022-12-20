package realmain

import (
	"errors"
	"os"

	"github.com/backdround/deploy-configs/internal/config"
	"github.com/backdround/deploy-configs/internal/dataconverter"
	"github.com/backdround/deploy-configs/internal/deploy/commands"
	"github.com/backdround/deploy-configs/internal/deploy/links"
	"github.com/backdround/deploy-configs/internal/deploy/templates"
	"github.com/backdround/deploy-configs/internal/pathexpander"
	"github.com/backdround/deploy-configs/pkg/fsutility"
	"github.com/backdround/deploy-configs/pkg/logger"
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

func Main(l logger.Logger, cliArguments []string) int {
	// Gets config instance
	userInput := cliArguments[1:]
	if len(userInput) != 1 {
		l.Fail("expected config instance as argument")
		return 1
	}
	configInstance := userInput[0]

	// Gets cwd
	cwd, err := os.Getwd()
	if err != nil {
		l.Fail("unable to get current work directory:")
		l.Fail(err.Error())
		return 1
	}

	// Searches config path
	configPath, err := FindConfig(cwd, "deploy-configs.yml",
		"deploy-configs.yaml")
	if err != nil {
		l.Fail("error occurs while config searching:")
		l.Fail(err.Error())
		return 1
	}

	// Reads config yaml
	configData, err := os.ReadFile(configPath)
	if err != nil {
		l.Fail("unable to read config data:")
		l.Fail(err.Error())
		return 1
	}

	// Parse config data
	config, err := config.Get(configData, configInstance)
	if err != nil {
		l.Fail("fail to parse config data:")
		l.Fail(err.Error())
		return 1
	}

	// Restructures config to deploy data
	pathExpander := pathexpander.New(l, cwd)
	dataConverter := dataconverter.New(l, pathExpander)

	restructuredLinks, err := dataConverter.RestructureLinks(config.Links)
	if err != nil {
		l.Fail("invalid config links:")
		l.Fail(err.Error())
		return 1
	}

	restructuredTemplates, err := dataConverter.RestructureTemplates(
		config.Templates)
	if err != nil {
		l.Fail("invalid config templates:")
		l.Fail(err.Error())
		return 1
	}

	restructuredCommands, err := dataConverter.RestructureCommands(
		config.Commands)
	if err != nil {
		l.Fail("invalid config commands:")
		l.Fail(err.Error())
		return 1
	}

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

	return 0
}
