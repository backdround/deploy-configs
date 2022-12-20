// templates describes templateMaker which receives a bunch
// of templates, expands these and logs all outcomes.
package templates

import (
	"bytes"
	"fmt"
	"os"
	"path"
	templatePackage "text/template"

	"github.com/backdround/deploy-configs/pkg/fsutility"
	"github.com/backdround/go-indent"
)

type templateMaker struct {
	logger Logger
}

func NewTemplateMaker(logger Logger) templateMaker {
	return templateMaker{
		logger: logger,
	}
}

func getDescription(template Template) string {
	return fmt.Sprintf("input: %q\noutput: %q\ndata: %q",
		template.InputPath, template.OutputPath, template.Data)
}

func shift(message string, count int) string {
	return indent.Indent(message, "  ", count)
}

func (m templateMaker) logFail(template Template, reason string) {
	description := shift(getDescription(template), 1)
	errorMessage := shift("error: " + reason, 2)

	message := fmt.Sprintf("Unable to expand %q template:\n%v\n%v",
		template.Name, description, errorMessage)
	m.logger.Fail(message)
}

func (m templateMaker) logSuccess(template Template) {
	message := fmt.Sprintf("Template %q expanded:\n%v",
		template.Name,
		shift(getDescription(template), 1))
	m.logger.Success(message)
}

func (m templateMaker) logSkip(template Template) {
	message := fmt.Sprintf("Template %q skipped", template.Name)
	m.logger.Log(message)
}

func (m templateMaker) makeTemplate(t Template) (success bool) {
	// Gets expanded data
	template, err := templatePackage.ParseFiles(t.InputPath)
	if err != nil {
		m.logFail(t, err.Error())
		return false
	}

	outputBuffer := bytes.NewBuffer([]byte{})
	err = template.Option("missingkey=error").Execute(outputBuffer, t.Data)
	if err != nil {
		m.logFail(t, err.Error())
		return false
	}

	// Checks if the output file is already expanded
	oldOutputFileHash := fsutility.GetFileHash(t.OutputPath)
	newOutputFileHash := fsutility.GetHash(outputBuffer.Bytes())
	if bytes.Equal(oldOutputFileHash, newOutputFileHash) {
		m.logSkip(t)
		return true
	}

	// Creates tha output file directory
	err = fsutility.MakeDirectoryIfDoesntExist(path.Dir(t.OutputPath))
	if err != nil {
		m.logFail(t, err.Error())
		return false
	}

	// Creates the expanded file
	err = os.WriteFile(t.OutputPath, outputBuffer.Bytes(), 0644)
	if err != nil {
		m.logFail(t, err.Error())
		return false
	}

	m.logSuccess(t)
	return true
}

// MakeTemplates expands the given templates.
func (m templateMaker) MakeTemplates(templates []Template) (success bool) {
	success = true
	for _, template := range templates {
		success = success && m.makeTemplate(template)
	}
	return success
}
