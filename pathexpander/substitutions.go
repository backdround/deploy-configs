package pathexpander

import (
	"errors"
	"os"
	"os/user"
	"path"
)

// getGitRoot searches a git directory in parent directories
// descending up to the root.
func getGitRoot(initialDirectoryToSearch string) (string, error) {
	isADirectory := func (path string) bool {
		fileInfo, err := os.Stat(path)
		return err == nil && fileInfo.IsDir()
	}

	getDescendingParentDirectories := func(directory string) []string {
		parents := []string{directory}
		nextParent := path.Dir(directory)
		for directory != nextParent {
			parents = append(parents, nextParent)
			directory = nextParent
			nextParent = path.Dir(directory)
		}
		return parents
	}

	// Checks that the input directory is a directory
	if !isADirectory(initialDirectoryToSearch) {
		return "", errors.New("initial directory to search isn't a directory")
	}

	parentDirectories := getDescendingParentDirectories(initialDirectoryToSearch)

	// Searches a git directory descending to root
	for _, currentDirectory := range  parentDirectories {
		hypotheticalGitDirectory := path.Join(currentDirectory, ".git")
		if isADirectory(hypotheticalGitDirectory) {
			return currentDirectory, nil
		}
	}

	return "", errors.New("Git directory wasn't founded")
}

// Searches the home directory
func getHomeDirectory() (string, error) {
	homedir, err := os.UserHomeDir()
	if err == nil {
		return homedir, nil
	}

	user, err := user.Current()
	if err == nil {
		return user.HomeDir, nil
	}

	return "", err
}
