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
	// Checks the target path
	targetType := fsutility.GetPathType(link.TargetPath)
	if targetType == fsutility.Notexisting {
		m.logFail(link, "target path isn't exist")
		return false
	}

	// Creates the link directory
	linkDirectory := path.Dir(link.LinkPath)
	err := fsutility.MakeDirectoryIfDoesntExist(linkDirectory)
	if err != nil {
		m.logFail(link, err.Error())
		return false
	}

	linkType := fsutility.GetPathType(link.LinkPath)

	// Checks that the link already points to target
	if linkType == fsutility.Symlink {
		skip := fsutility.IsLinkPointsToDestination(link.LinkPath,
			link.TargetPath)
		if skip {
			m.logSkip(link)
			return true
		}
	}

	// Checks the link to replace
	if linkType == fsutility.Symlink {
		err := os.Remove(link.LinkPath)
		if err != nil {
			message := "unable to replace link:\n  " + err.Error()
			m.logFail(link, message)
			return false
		}
	}

	// Creates the link
	linkType = fsutility.GetPathType(link.LinkPath)
	if linkType == fsutility.Notexisting {
		err = os.Symlink(link.TargetPath, link.LinkPath)
		if err != nil {
			message := "unable to create link:\n  " + err.Error()
			m.logFail(link, message)
			return false
		}

		m.logSuccess(link)
		return true
	}

	m.logFail(link, "link path is occupied")
	return false
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
