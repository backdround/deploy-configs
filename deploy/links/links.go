// links describes linkMaker which receives a bunch of links,
// creates these and logs all outcomes.
package links

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/backdround/deploy-configs/pkg/fsutility"
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
	case fsutility.Regular, fsutility.Directory, fsutility.Unknown:
		return stopLinkFileExists
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

func (m linkMaker) logFail(link Link, reason string) {
	linkDescription := fmt.Sprintf("[%q, %q]", link.TargetPath, link.LinkPath)
	message := fmt.Sprintf("Unable to create %q link:\n\t%v\n\t\t%v",
		link.Name, linkDescription, reason)
	m.logger.Fail(message)
}

func (m linkMaker) logSuccess(link Link) {
	linkDescription := fmt.Sprintf("[%q, %q]", link.TargetPath, link.LinkPath)
	message := fmt.Sprintf("Link %q created:\n\t%v", link.Name, linkDescription)
	m.logger.Success(message)
}

func (m linkMaker) logSkip(link Link) {
	linkDescription := fmt.Sprintf("[%q, %q]", link.TargetPath, link.LinkPath)
	message := fmt.Sprintf("Link %q skipped:\n\t%v", link.Name, linkDescription)
	m.logger.Log(message)
}

func (m linkMaker) makeLink(link Link) {
	createLink := func() {
		// Checks link directory
		linkDirectory := path.Dir(link.LinkPath)
		err := fsutility.MakeDirectoryIfDoesntExist(linkDirectory)
		if err != nil {
			m.logFail(link, err.Error())
		}

		// Creates link
		err = os.Symlink(link.TargetPath, link.LinkPath)
		if err != nil {
			m.logFail(link, err.Error())
		} else {
			m.logSuccess(link)
		}
	}

	action := linkDecisionMaker(link)

	switch action {
	case proceedNew:
		createLink()
	case proceedRemove:
		err := os.Remove(link.LinkPath)
		if err != nil {
			m.logFail(link, err.Error())
			break
		}
		createLink()
	case stopTargetDoesntExist:
		m.logFail(link, "Target file isn't exist")
	case stopLinkFileExists:
		m.logFail(link, "Link file already exists")
	case skip:
		m.logSkip(link)
	}
}

// Links creates links which are described in links parameter.
// If target is a directory it creates appropriate symlinks
// for all files in that directory
func (m linkMaker) Links(links []Link) {
	isDirectory := func(path string) bool {
		stat, err := os.Lstat(path)
		if err != nil {
			return false
		}
		return stat.IsDir()
	}

	for _, link := range links {
		if !isDirectory(link.TargetPath) {
			m.makeLink(link)
			continue
		}

		// Reads all files in directory
		fileInfos, err := ioutil.ReadDir(link.TargetPath)
		if err != nil {
			m.logFail(link, err.Error())
			continue
		}

		// Makes link for every file in target directory.
		for _, fileInfo := range fileInfos {
			targetFileName := path.Base(fileInfo.Name())

			specificName := link.Name + "/" + targetFileName
			specificTargetFile := path.Join(link.TargetPath, targetFileName)
			specificLinkPath := path.Join(link.LinkPath, targetFileName)

			specificLink := Link{
				Name: specificName,
				TargetPath: specificTargetFile,
				LinkPath:   specificLinkPath,
			}
			m.makeLink(specificLink)
		}
	}
}