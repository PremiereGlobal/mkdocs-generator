package main

import (
	bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
	"path/filepath"
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
			log.Debugf("Error browsing repo %s/%s/%s", r.repo.Project.Name, r.repo.Slug, path)
			continue
		}

		// Loop over the files
		for _, f := range browseList.Children.Values {

			// If the file is markdown, process it
			if f.Path.Extension == "md" {

				// Generate the document object for this file
				document := NewDocument(r.repo.Project.Key, r.repo.Slug, filepath.Join(browseList.Path.ToString, f.Path.ToString))
				document.docType = markdownType

				// Create a task to process this file
				task := fileTask{document: document}

				// Add the file to the master list so nothing else processes it
				masterFileList.Store(document.uid, document)

				// Add a count to the waitgroup and add the task to the queue
				wg.Add(1)
				taskChan <- task

			}
		}
	}

	return true
}
