package pathexpander

import (
	"errors"
	"os"
	"os/user"
	"path"

	"github.com/backdround/deploy-configs/pkg/fsutility"
)

// getGitRoot searches a git directory in parent directories
// descending up to the root.
func getGitRoot(initialDirectoryToSearch string) (string, error) {
	// Checks that the input directory is a directory
	if fsutility.GetPathType(initialDirectoryToSearch) != fsutility.Directory {
		return "", errors.New("initial directory to search isn't a directory")
	}

	gitPath, err := fsutility.FindEntryDescending(initialDirectoryToSearch,
		".git", fsutility.Directory)
	if err != nil {
		return "", errors.New("git directory wasn't founded")
	}
	gitRoot := path.Dir(gitPath)

	return gitRoot, nil
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
