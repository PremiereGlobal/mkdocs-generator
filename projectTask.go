package main

import (
	bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
)

type projectTask struct {
	project bitbucket.BBProject
}

func (p projectTask) run(workerNum int, taskChan chan<- task) bool {

	log.Infof("[worker:%03d] Processing project task %s", workerNum, p.project.GetKey())

	// Get the list of repos in this project
	repos, err := p.project.ListRepos()
	if err != nil {
		log.Fatalf("[worker:%03d] Unable to list repos for %s:%s", workerNum, p.project.GetKey(), err)
	}

	// Loop through the repos and add them to the queue
	for _, r := range repos {

		// Create the repo task
		task := repoTask{repo: r}

		taskChan <- task
	}
	return true
}
