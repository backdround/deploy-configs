package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/backdround/deploy-configs/config"
	"github.com/backdround/deploy-configs/config/validate"
	"github.com/backdround/deploy-configs/deploy/commands"
	"github.com/backdround/deploy-configs/deploy/links"
	"github.com/backdround/deploy-configs/deploy/templates"
	"github.com/backdround/deploy-configs/logger"
	"github.com/backdround/deploy-configs/pathexpander"
	"github.com/backdround/deploy-configs/pkg/fsutility"
)

func FindConfig(cwd string, names ...string) (configPath string, err error) {
	for _, name := range names {
		types := fsutility.Regular & fsutility.Symlink
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

type MustPathExpander = func(unitName string, unitDescription string,
	templateToExpand string) string

// RestructureLinks resturctures config links to deploy links
func RestructureLinks(mustPathExpand MustPathExpander,
	configLinks map[string]config.Link) []links.Link {
	// Restructures config links to deploy links
	newLinks := []links.Link{}
	for linkName, link := range configLinks {
		newStructuredLink := links.Link{
			Name:       linkName,
			TargetPath: link.TargetPath,
			LinkPath:   link.LinkPath,
		}
		newLinks = append(newLinks, newStructuredLink)
	}

	// Expands links paths
	for i, link := range newLinks {
		newLinks[i].TargetPath = mustPathExpand(link.Name, "link",
			link.TargetPath)
		newLinks[i].LinkPath = mustPathExpand(link.Name, "link",
			link.LinkPath)
	}

	return newLinks
}

// RestructureTemplates resturctures config templates to deploy templates
func RestructureTemplates(mustPathExpand MustPathExpander,
	configTemplates map[string]config.Template) []templates.Template {
	// Restructures config templates to deploy templates
	newTemplates := []templates.Template{}
	for templateName, template := range configTemplates {
		newStructuredTemplate := templates.Template{
			Name:       templateName,
			InputPath:  template.InputPath,
			OutputPath: template.OutputPath,
			Data:       template.Data,
		}
		newTemplates = append(newTemplates, newStructuredTemplate)
	}

	// Expands templates paths
	for i, template := range newTemplates {
		newTemplates[i].InputPath = mustPathExpand(template.Name, "template",
			template.InputPath)
		newTemplates[i].OutputPath = mustPathExpand(template.Name, "template",
			template.OutputPath)
	}

	return newTemplates
}

// RestructureCommands resturctures config commands to deploy commands
func RestructureCommands(mustPathExpand MustPathExpander,
	configCommands map[string]config.Command) []commands.Command {
	// Restructures config commands to deploy commands
	newCommands := []commands.Command{}
	for commandName, command := range configCommands {
		newStructuredCommand := commands.Command{
			Name:            commandName,
			InputPath:       command.InputPath,
			OutputPath:      command.OutputPath,
			CommandTemplate: command.Command,
		}
		newCommands = append(newCommands, newStructuredCommand)
	}

	// Expands commands paths
	for i, command := range newCommands {
		newCommands[i].InputPath = mustPathExpand(command.Name, "command",
			command.InputPath)
		newCommands[i].OutputPath = mustPathExpand(command.Name, "command",
			command.OutputPath)
	}

	return newCommands
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
	configPath, err := FindConfig(cwd, "deploy-config.yml", "deploy-config.yaml")
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

	// Creates auxiliaries to expand
	pathExpander := pathexpander.New(l, cwd)
	mustExpand := func(unitName string, unitDescription string,
		templateToExpand string) string {
		expandedTemplate, err := pathExpander.Expand(templateToExpand)
		if err != nil {
			message := fmt.Sprintf("unable to expand %v %q:\n%v",
				unitDescription, unitName, templateToExpand)
			CheckFatalError(err, l, message)
		}
		return expandedTemplate
	}

	// Restructures config to deploy data
	restructuredLinks := RestructureLinks(mustExpand, config.Links)
	restructuredTemplates := RestructureTemplates(mustExpand, config.Templates)
	restructuredCommands := RestructureCommands(mustExpand, config.Commands)

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
