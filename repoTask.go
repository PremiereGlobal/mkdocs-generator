package main

import (
	"path/filepath"

	bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
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

	// Get the list of files at this path
	// There's a good chance that there's not a docs folder so don't do anything here...
	browseList, err := bb.Browse(&r.repo, "/")
	if err != nil {
		log.Debugf("Error browsing repo %s/%s/%s: %v", r.repo.Project.Name, r.repo.Slug, "/", err)
		return true
	}

	// Loop over the files
	r.checkFiles(browseList)
	for _, f := range browseList.Children.Values {
		if f.FileType == "DIRECTORY" && f.Path.Name == "docs" {
			browseList, err := bb.Browse(&r.repo, "/docs")
			if err != nil {
				log.Debugf("Error browsing repo %s/%s/%s: %v", r.repo.Project.Name, r.repo.Slug, "/docs", err)
				return true
			}
			r.checkFiles(browseList)
			return true
		}
	}

	return true
}

func (r repoTask) checkFiles(bl *bitbucket.BrowseList) {
	for _, f := range bl.Children.Values {

		// If the file is markdown, process it
		if f.Path.Extension == "md" {

			// Generate the document object for this file
			document := NewDocument(r.repo.Project.Key, r.repo.Slug, filepath.Join(bl.Path.ToString, f.Path.ToString))
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
