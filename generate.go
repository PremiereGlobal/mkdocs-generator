package main

import (
	bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
	// "os"
	"sync"
)

// masterFileList holds all of the files that have been processed so we don't duplicate
// Map keys should be in the formats:
//   projects/<project>/repos/<repo>/raw/<filepath>
//   or
//   projects/<project>/repos/<repo>/browse/<filepath>
var masterFileList sync.Map

// This is our main waitgroup that counts items added/removed from the process
// queue.  When this gets to 0, we're done
var wg sync.WaitGroup

// Define our worker channels
// We create a channel for each type of workload because don't want to block.
// For example, files not being able to be processed because the work queue is
// full of repos
var projectChan chan task
var repoChan chan task
var fileChan chan task

// generateConfig is our configuration type
type generateConfig struct {
	bitbucketUrl      string
	bitbucketUser     string
	bitbucketPassword string
}

// config contains our configuration
var config generateConfig

// worker is the main worker function that processes all tasks
// This will be called in a goroutine
func worker(workerNum int, taskChan <-chan task) {
	for task := range taskChan {
		task.run(workerNum)
	}
}

func generate() {

	// Load up our config
	config = generateConfig{
		bitbucketUrl:      Args.GetString("bitbucket-url"),
		bitbucketUser:     Args.GetString("bitbucket-user"),
		bitbucketPassword: Args.GetString("bitbucket-password"),
	}

	// Ensure the build directory is good to go
	ensureBuildDir()

	// Create our channels that will buffer up to x tasks at a time
	// The buffer needs to be big enough so that one repo/file cannot fill it up
	projectChan = make(chan task, 1) // This can be anything
	repoChan = make(chan task, 20)   // This should be > the most repos in a project
	fileChan = make(chan task, 100)  // This should be > the most links in a md file

	// Start the workers
	projectworkerCount := 5
	for i := 0; i < projectworkerCount; i++ {
		go worker(i, projectChan)
	}
	repoWorkerCount := 20
	for i := 0; i < repoWorkerCount; i++ {
		go worker(i, repoChan)
	}
	fileWorkerCount := 40
	for i := 0; i < fileWorkerCount; i++ {
		go worker(i, fileChan)
	}

	// We add one to the waitgroup intitially because we want to make sure we block`
	// until we get through adding all the project tasks to the queue
	wg.Add(1)

	// Get the list of projects
	bb, err := bitbucket.NewBitbucketClient(config.bitbucketUrl, config.bitbucketUser, config.bitbucketPassword)
	projects, err := bb.ListProjects()
	if err != nil {
		log.Fatal("Unable to list projects: ", err)
	}

	// Loop through the projects and add a project task to the queue
	for _, p := range projects.Values {
		taskProject := p
		task := projectTask{project: &taskProject}

		// Add a count to the waitgroup and add the task to the queue
		wg.Add(1)
		projectChan <- task
	}

	// We're done adding all the projects, so remove our main blocker so that
	// the program can exit as soon as all the projects are done
	wg.Done()

	// Now wait for all the tasks to finish
	wg.Wait()
}

// ensureBuildDir ensures that the build directory exists, is a directory and
// is empty, creating it if need be
func ensureBuildDir() {
	buildDir := Args.GetString("build-dir")
	if ok, _ := PathExists(buildDir); ok {
		if ok, _ := IsDirectory(buildDir); !ok {
			log.Fatal("Build directory path exists and is not a directory")
		}
		if empty, _ := IsDirEmpty(buildDir); !empty {
			log.Fatal("Build directory exists and is not empty")
		}
	} else {
		log.Debug("Creating build directory ", buildDir)
		err := CreateDirIfNotExist(Args.GetString("build-dir"))
		if err != nil {
			log.Fatal("Unable to create build directory ", err)
		}
	}
}
