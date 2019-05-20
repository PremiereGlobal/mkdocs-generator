package main

import (
	bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
	"path/filepath"
	"strings"
)

type repoTask struct {
	repo bitbucket.Repo
}

func (r repoTask) run(workerNum int) bool {

	// Decrement waitgroup counter when we're done
	defer wg.Done()

	log.Info("Processing repo task ", r.repo.MakePath(), " [worker:", workerNum, "]")

  // Create new Bitbucket client
	bb := NewBitbucketClient()

	// Paths to browse for markdown
	browsePaths := []string{"/", "/docs"}
	for _, path := range browsePaths {

		// Get the list of files at this path
		// There's a good chance that there's not a docs folder so don't do anything here...
		browseList, err := bb.Browse(&r.repo, path)
		if err != nil {
			continue
		}

		// Loop over the files
		for _, f := range browseList.Children.Values {

			// If the file is markdown, process it
			if strings.ToLower(f.Path.Extension) == "md" {

				// Generate the master file path for this file
				masterFilePath := filepath.Join(r.repo.MakePath(), "raw", filepath.Join(browseList.Path.ToString, f.Path.ToString))

				// Create a task to process this file
				task := fileTask{
					fileType:       "markdown",
					masterFilePath: masterFilePath,
				}

				// Add the file to the master list so nothing else processes it
				masterFileList.Store(strings.ToLower(masterFilePath), true)

				// Add a count to the waitgroup and add the task to the queue
				wg.Add(1)
				taskChan <- task

			}
		}
	}

	return true
}
