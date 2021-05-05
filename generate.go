package main

import (
	"path/filepath"
	"sync"

	bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
)

// masterFileList holds all of the files that have been processed so we don't duplicate
// Map keys should be in the formats:
//   projects/<project>/repos/<repo>/raw/<filepath>
//   or
//   projects/<project>/repos/<repo>/browse/<filepath>
var masterFileList sync.Map

func generate() {

	// Ensure the build directory is good to go
	ensureBuildDir()

	bbConfig := bitbucket.BitbucketClientConfig{
		Url:       Args.GetString("bitbucket-url"),
		Username:  Args.GetString("bitbucket-user"),
		Password:  Args.GetString("bitbucket-password"),
		Workspace: Args.GetString("bitbucket-workspace"),
		Logger:    log,
	}

	bb, err := bitbucket.NewBitbucketClient(&bbConfig)
	if err != nil {
		log.Fatal("Unable to create Bitbucket client ", err)
	}

	workerCount := Args.GetInt("workers")
	if workerCount <= 0 {
		workerCount = 1
	}

	taskChan, _, wg := NewTaskQueue(workerCount)

	// Get the list of projects
	projects, err := bb.ListProjects()
	if err != nil {
		log.Fatal("Unable to list projects: ", err)
	}

	// Loop through the projects and add a project task to the queue
	for _, p := range projects {
		taskProject := p
		task := projectTask{project: taskProject}
		taskChan <- task
	}

	// Now wait for all the tasks to finish
	wg.Wait()

	// If user provided mkdocs directory
	docsDir := Args.GetString("docs-dir")
	if docsDir != "" {
		makeNav(docsDir)
	}
}

// ensureBuildDir ensures that the build directory exists, is a directory and
// is empty, creating it if need be
func ensureBuildDir() {
	buildDir := filepath.Join(Args.GetString("build-dir"), "docs")
	if ok, _ := PathExists(buildDir); ok {
		if ok, _ := IsDirectory(buildDir); !ok {
			log.Fatal("Build directory path exists and is not a directory")
		}
		if empty, _ := IsDirEmpty(buildDir); !empty {
			log.Fatal("Build directory exists and is not empty")
		}
	} else {
		log.Debug("Creating build directory ", buildDir)
		err := CreateDirIfNotExist(buildDir)
		if err != nil {
			log.Fatal("Unable to create build directory ", err)
		}
	}
}
