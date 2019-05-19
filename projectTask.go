package main

import (
  bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
)

type projectTask struct {
  project *bitbucket.Project
}

func (p projectTask) run(workerNum int) bool {

  // Decrement waitgroup counter when we're done
  defer wg.Done()

  log.Info("Processing project task ", p.project.MakePath(), " [project-worker:", workerNum, "]")
  bb, err := bitbucket.NewBitbucketClient(config.bitbucketUrl, config.bitbucketUser, config.bitbucketPassword)
  repos, err := bb.ListRepos(p.project)
  if err != nil {
    log.Fatal("Unable to list repos for ", p.project.MakePath(), ": ", err)
  }
  for _, r := range repos.Values {
    // if r.Slug == "doc-test" {
      // Create the repo task
      task := repoTask{repo: r}

      // Add a count to the waitgroup and add the task to the queue
      wg.Add(1)
      repoChan <- task
    // }
  }
  return true
}
