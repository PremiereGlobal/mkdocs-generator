package main

import (
	bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
)

type repoTask struct {
	repo bitbucket.BBRepo
}

func (r repoTask) run(workerNum int, taskChan chan<- task) bool {
	log.Infof("[worker:%03d] Processing repo task %s", workerNum, r.repo.GetFilesURL())

	// Get the list of files at this path
	// There's a good chance that there's not a docs folder so don't do anything here...
	browseList, err := r.repo.GetDir("/")
	if err != nil {
		log.Debugf("[worker:%03d] Error browsing repo %s/%s/%s: %v", workerNum, r.repo.GetBBProject().GetName(), r.repo.GetSlug(), "/", err)
		return true
	}

	// Loop over the files
	r.checkFiles(browseList, taskChan, workerNum)
	for _, f := range browseList {
		if f.GetFileType() == "DIRECTORY" && f.GetName() == "docs" {
			browseList, err := r.repo.GetDir("/docs")
			if err != nil {
				log.Debugf("[worker:%03d] Error browsing repo %s/%s/%s: %v", workerNum, r.repo.GetBBProject().GetName(), r.repo.GetSlug(), "/docs", err)
				return true
			}
			r.checkFiles(browseList, taskChan, workerNum)
			return true
		}
	}

	return true
}

func (r repoTask) checkFiles(bbfl []bitbucket.BBFile, taskChan chan<- task, workerNum int) {
	for _, f := range bbfl {

		// If the file is markdown, process it
		fname := f.GetName()
		if len(fname) > 4 && fname[len(fname)-3:] == ".md" {
			log.Debugf("[worker:%03d] Added File:%s", workerNum, fname)
			// Generate the document object for this file
			document := NewDocument(r.repo.GetBBProject().GetKey(), r.repo.GetSlug(), f.GetFullPath(), f)
			document.docType = markdownType

			//See if we can store the document if we can, because its not there, we will add a task (on store ok==false)
			if _, ok := masterFileList.LoadOrStore(document.uid, document); !ok {
				// Create a task to process this file
				task := fileTask{document: document, file: f}
				taskChan <- task
			}
		}
	}

}
