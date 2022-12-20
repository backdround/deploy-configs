// links describes linkMaker which receives a bunch of links,
// creates these and logs all outcomes.
package links

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/backdround/deploy-configs/pkg/fsutility"
	"github.com/backdround/go-indent"
)

////////////////////////////////////////////////////////////
// linkAction

// linkAction describes what to do with a hypothetical link
type linkAction int

const (
	proceedNew linkAction = iota
	proceedRemove
	stopTargetDoesntExist
	stopLinkFileExists
	stopLinkPathExists
	skip
)

// linkDecisionMaker chooses what to do with link,
// based on the filesystem state
func linkDecisionMaker(link Link) linkAction {
	// Checks target path
	if fsutility.GetPathType(link.TargetPath) == fsutility.Notexisting {
		return stopTargetDoesntExist
	}

	// Checks link path
	linkType := fsutility.GetPathType(link.LinkPath)
	switch linkType {
	case fsutility.Notexisting:
		return proceedNew
	case fsutility.Regular:
		return stopLinkFileExists
	case fsutility.Directory, fsutility.Unknown:
		return stopLinkPathExists
	case fsutility.Symlink:
		if fsutility.IsLinkPointsToDestination(link.LinkPath, link.TargetPath) {
			return skip
		} else {
			return proceedRemove
		}
	}

	panic("unknown pathType")
}

////////////////////////////////////////////////////////////
// linkMaker

// linkMaker makes link and logs all outcomes.
type linkMaker struct {
	logger Logger
}

func NewLinkMaker(logger Logger) linkMaker {
	return linkMaker{
		logger: logger,
	}
}

func getDescription(link Link) string {
	return fmt.Sprintf("target: %q\nlink: %q",
		link.TargetPath, link.LinkPath)
}

func shift(message string, count int) string {
	return indent.Indent(message, "  ", count)
}

func (m linkMaker) logFail(link Link, reason string) {
	description := shift(getDescription(link), 1)
	errorMessage := shift("error: "+reason, 2)

	message := fmt.Sprintf("Unable to create %q link:\n%v\n%v",
		link.Name, description, errorMessage)
	m.logger.Fail(message)
}

func (m linkMaker) logSuccess(link Link) {
	message := fmt.Sprintf("Link %q created:\n%v", link.Name,
		shift(getDescription(link), 1))
	m.logger.Success(message)
}

func (m linkMaker) logSkip(link Link) {
	message := fmt.Sprintf("Link %q skipped", link.Name)
	m.logger.Log(message)
}

func (m linkMaker) makeLink(link Link) (success bool) {
	createLink := func() (success bool) {
		// Checks link directory
		linkDirectory := path.Dir(link.LinkPath)
		err := fsutility.MakeDirectoryIfDoesntExist(linkDirectory)
		if err != nil {
			m.logFail(link, err.Error())
			return false
		}

		// Creates link
		err = os.Symlink(link.TargetPath, link.LinkPath)
		if err != nil {
			m.logFail(link, err.Error())
			return false
		}

		m.logSuccess(link)
		return true
	}

	action := linkDecisionMaker(link)

	switch action {
	case proceedNew:
		return createLink()
	case proceedRemove:
		err := os.Remove(link.LinkPath)
		if err != nil {
			m.logFail(link, err.Error())
			return false
		}
		return createLink()
	case stopTargetDoesntExist:
		m.logFail(link, "Target file isn't exist")
		return false
	case stopLinkFileExists:
		m.logFail(link, "Link file already exists")
		return false
	case stopLinkPathExists:
		m.logFail(link, "Link path exists")
		return false
	case skip:
		m.logSkip(link)
		return true
	default:
		panic("unknown action")
	}
}

// CreateLinks creates links which are described in links parameter.
// If target is a directory it creates appropriate symlinks
// for all files in that directory
func (m linkMaker) CreateLinks(links []Link) (globalSuccess bool) {
	isDirectory := func(path string) bool {
		stat, err := os.Lstat(path)
		if err != nil {
			return false
		}
		return stat.IsDir()
	}

	globalSuccess = true
	for _, link := range links {
		if !isDirectory(link.TargetPath) {
			success := m.makeLink(link)
			globalSuccess = globalSuccess && success
			continue
		}

		// Reads all files in directory
		fileInfos, err := ioutil.ReadDir(link.TargetPath)
		if err != nil {
			m.logFail(link, err.Error())
			globalSuccess = false
			continue
		}

		// Makes link for every file in target directory.
		for _, fileInfo := range fileInfos {
			targetFileName := path.Base(fileInfo.Name())

			specificName := link.Name + "/" + targetFileName
			specificTargetFile := path.Join(link.TargetPath, targetFileName)
			specificLinkPath := path.Join(link.LinkPath, targetFileName)

			specificLink := Link{
				Name:       specificName,
				TargetPath: specificTargetFile,
				LinkPath:   specificLinkPath,
			}
			success := m.makeLink(specificLink)
			globalSuccess = globalSuccess && success
		}
	}

	return globalSuccess
}
