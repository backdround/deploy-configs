// Describes pathexpander that expands paths in given templates.
package pathexpander

import (
	"bytes"
	templatePackage "text/template"
)

type Logger interface {
	Log(message string)
	Warn(message string)
}

// pathexpander expands paths substitutions in given templates.
type pathexpander struct {
	data map[string]string
}

// New creates new pathexpander. searchGitFromDirectory is used to
// search project git root.
func New(l Logger, searchGitFromDirectory string) *pathexpander {
	log := func(message string) {
		l.Log("path-expander: " + message)
	}

	warn := func(message string) {
		l.Warn("path-expander: " + message)
	}

	expander := pathexpander{
		data: make(map[string]string),
	}

	// Adds "git-root" key
	gitRoot, err := getGitRoot(searchGitFromDirectory)
	if err == nil {
		expander.data["gitRoot"] = gitRoot
		log("gitRoot: " + gitRoot)
	} else {
		warn("Unable to get gitRoot")
	}

	// Adds "home" key
	homeDirectory, err := getHomeDirectory()
	if err == nil {
		expander.data["home"] = homeDirectory
		log("home: " + homeDirectory)
	} else {
		warn("Unable to get home")
	}

	return &expander
}

// Expand expands paths in template. It returns error if
// template is invalid or used keys that don't exist
func (expander pathexpander) Expand(template string) (string, error) {
	t := templatePackage.New("path-expander").Option("missingkey=error")
	t, err := t.Parse(template)
	if err != nil {
		return "", err
	}

	outputBuffer := bytes.NewBuffer([]byte{})
	err = t.Execute(outputBuffer, expander.data)
	if err != nil {
		return "", err
	}

	return outputBuffer.String(), nil
}
