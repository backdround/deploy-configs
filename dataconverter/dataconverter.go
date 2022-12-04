// dataconverter converts config data structures to deploy data structures.
// Also it expands paths by using pathexpander.PathExpander.
package dataconverter

import (
	"fmt"

	"github.com/backdround/deploy-configs/config"
	"github.com/backdround/deploy-configs/deploy/commands"
	"github.com/backdround/deploy-configs/deploy/links"
	"github.com/backdround/deploy-configs/deploy/templates"
	"github.com/backdround/deploy-configs/pathexpander"
)

type Logger interface {
	Log(message string)
	Fail(message string)
}

type dataConverter struct {
	logger       Logger
	pathExpander pathexpander.PathExpander
}

func New(logger Logger,
	pathExpander pathexpander.PathExpander) *dataConverter {
	return &dataConverter{
		logger:       logger,
		pathExpander: pathExpander,
	}
}

func (c dataConverter) pathExpand(unitName string, unitDescription string,
	templateToExpand string) (expandedTemplate string, err error) {
	expandedTemplate, err = c.pathExpander.Expand(templateToExpand)
	if err != nil {
		err = fmt.Errorf("unable to expand %v %q:\n%v\n\t%v",
			unitDescription, unitName, templateToExpand, err.Error())
	}
	return expandedTemplate, err
}

// RestructureLinks resturctures config links to deploy links
func (c dataConverter) RestructureLinks(
	configLinks map[string]config.Link) ([]links.Link, error) {

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
		expandedTemplate, err := c.pathExpand(link.Name, "link",
			link.TargetPath)

		if err != nil {
			return nil, err
		}
		newLinks[i].TargetPath = expandedTemplate

		expandedTemplate, err = c.pathExpand(link.Name, "link",
			link.LinkPath)
		if err != nil {
			return nil, err
		}
		newLinks[i].LinkPath = expandedTemplate
	}

	return newLinks, nil
}

// RestructureTemplates resturctures config templates to deploy templates
func (c dataConverter) RestructureTemplates(
	configTemplates map[string]config.Template) ([]templates.Template, error) {
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
		expandedTemplate, err := c.pathExpand(template.Name, "template",
			template.InputPath)
		if err != nil {
			return nil, err
		}
		newTemplates[i].InputPath = expandedTemplate

		expandedTemplate, err = c.pathExpand(template.Name, "template",
			template.OutputPath)
		if err != nil {
			return nil, err
		}
		newTemplates[i].OutputPath = expandedTemplate
	}

	return newTemplates, nil
}

// RestructureCommands resturctures config commands to deploy commands
func (c dataConverter) RestructureCommands(
	configCommands map[string]config.Command) ([]commands.Command, error) {
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
		expandedTemplate, err := c.pathExpand(command.Name, "command",
			command.InputPath)
		if err != nil {
			return nil, err
		}
		newCommands[i].InputPath = expandedTemplate

		expandedTemplate, err = c.pathExpand(command.Name, "command",
			command.OutputPath)
		if err != nil {
			return nil, err
		}
		newCommands[i].OutputPath = expandedTemplate
	}

	return newCommands, nil
}
