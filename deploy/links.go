package deploy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func isLinkPointsToDestination(linkPath string, destination string) bool {
	linkDestination, err := os.Readlink(linkPath)
	if err != nil {
		return false
	}

	matched, err := path.Match(linkDestination, destination)

	if err != nil {
		return false
	}

	return matched
}

////////////////////////////////////////////////////////////
// linkAction
type linkAction int

const (
	proceedNew linkAction = iota
	proceedRemove
	stopTargetNotExisting
	stopLinkFileExist
	skip
)

// linkDecisionMaker chooses what to do with with link
func linkDecisionMaker(targetPath, linkPath string) linkAction {
	// Checks target path
	if getFileType(targetPath) == notexisting {
		return stopTargetNotExisting
	}

	// Checks link path
	linkType := getFileType(linkPath)
	switch linkType {
	case notexisting:
		return proceedNew
	case regular, unknown:
		return stopLinkFileExist
	case symlink:
		if isLinkPointsToDestination(linkPath, targetPath) {
			return skip
		} else {
			return proceedRemove
		}
	}

	panic("unknown fileType")
}

////////////////////////////////////////////////////////////
// linkMaker
type linkMaker struct {
	logger Logger
}

func NewLinkMaker(logger Logger) linkMaker {
	return linkMaker{
		logger: logger,
	}
}

func (m linkMaker) logFail(linkName string, link Link, reason string) {
	linkDescription := fmt.Sprintf("[%q, %q]", link.Target, link.LinkPath)
	message := fmt.Sprintf("Unable to create %q link:\n\t%v\n\t\t%v",
		linkName, linkDescription, reason)
	m.logger.Fail(message)
}

func (m linkMaker) logSuccess(linkName string, link Link) {
	linkDescription := fmt.Sprintf("[%q, %q]", link.Target, link.LinkPath)
	message := fmt.Sprintf("Link %q created: %v", linkName, linkDescription)
	m.logger.Success(message)
}

func (m linkMaker) logSkip(linkName string, link Link) {
	linkDescription := fmt.Sprintf("[%q, %q]", link.Target, link.LinkPath)
	message := fmt.Sprintf("Link %q skipped: %v", linkName, linkDescription)
	m.logger.Log(message)
}

func (m linkMaker) makeLink(linkName string, link Link) {
	createLink := func() {
		// Checks link directory
		linkDirectory := path.Dir(link.LinkPath)
		err := makeDirectoryIfDoesntExist(linkDirectory)
		if err != nil {
			m.logFail(linkName, link, err.Error())
		}

		// Creates link
		err = os.Symlink(link.Target, link.LinkPath)
		if err != nil {
			m.logFail(linkName, link, err.Error())
		} else {
			m.logSuccess(linkName, link)
		}
	}

	action := linkDecisionMaker(link.Target, link.LinkPath)

	switch action {
	case proceedNew:
		createLink()
	case proceedRemove:
		err := os.Remove(link.LinkPath)
		if err != nil {
			m.logFail(linkName, link, err.Error())
			break
		}
		createLink()
	case stopTargetNotExisting:
		m.logFail(linkName, link, "Target file isn't exist")
	case stopLinkFileExist:
		m.logFail(linkName, link, "Link file already exists")
	case skip:
		m.logSkip(linkName, link)
	}
}

// Links creates links. If target is a directory it creates
// appropriate symlinks for all files in that directory
func (m linkMaker) Links(links map[string]Link) {
	isDirectory := func(path string) bool {
		stat, err := os.Lstat(path)
		if err != nil {
			return false
		}
		return stat.IsDir()
	}

	for linkName, link := range links {
		if !isDirectory(link.Target) {
			m.makeLink(linkName, link)
			continue
		}

		// Reads all files in directory
		fileInfos, err := ioutil.ReadDir(link.Target)
		if err != nil {
			m.logFail(linkName, link, err.Error())
			continue
		}

		// Makes link for every file in target directory.
		for _, fileInfo := range fileInfos {
			targetFileName := path.Base(fileInfo.Name())

			specificName := linkName + "/" + targetFileName
			specificTargetFile := path.Join(link.Target, targetFileName)
			specificLinkPath := path.Join(link.LinkPath, targetFileName)

			specificLink := Link{
				Target:   specificTargetFile,
				LinkPath: specificLinkPath,
			}
			m.makeLink(specificName, specificLink)
		}
	}
}
