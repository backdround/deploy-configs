package fstestutility

import (
	"os"
	"sync"
)

// tempPathSearcher searches and gives available, unoccupied paths
// to temporary files
type tempPathSearcher struct {
	issuedPaths map[string]bool
	mutex sync.Mutex
}

func (searcher *tempPathSearcher) getUnoccupiedPath() string {
	path := ""

	searcher.mutex.Lock()
	defer searcher.mutex.Unlock()

	for true {
		file, err := os.CreateTemp("", "go_test.*.txt")
		AssertNoError(err)
		AssertNoError(os.Remove(file.Name()))

		path = file.Name()
		if searcher.issuedPaths[path] {
			continue
		} else {
			searcher.issuedPaths[path] = true
			break
		}
	}

	return path
}

var tempPathSearcherInstance = tempPathSearcher{
	issuedPaths: make(map[string]bool),
	mutex: sync.Mutex{},
}

// GetAvailableTempPath returns path to available (nonexistent)
// temporary file.
func GetAvailableTempPath() string {
	return tempPathSearcherInstance.getUnoccupiedPath()
}
