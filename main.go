package main

import (

  "os"
  // "strings"
  "sync"
  "net/url"
  // "path/filepath"

	"github.com/PremiereGlobal/mkdocs-generator/cmd"
	// "github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	// "sync"
  bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"

)

type mkDocsGenConfig struct {
  bitbucketUrl string
  bitbucketUser string
  bitbucketPassword string
  baseUrl url.URL
}

var config mkDocsGenConfig
// var bb *bitbucket.BitbucketClient
var log *logrus.Logger

// This holds all of the files that have been processed so we don't duplicate
// Keys should be in the format projects/<project>/repos/<repo>/raw/<filepath>
var masterFileList sync.Map

// This is our main waitgroup that lets us know when we are done.
var wg sync.WaitGroup

// Define our worker channels
// We create a channel for each type of workload because don't want to block,
// for example, repos from being processed if all the workers are used by projects
var projectChan chan task
var repoChan chan task
var fileChan chan task

func main() {

	log = logrus.New()
	// log.SetLevel(logrus.DebugLevel)
  log.SetLevel(logrus.WarnLevel)

	cmd.Init()
	cmd.Execute()

  config = mkDocsGenConfig{
    bitbucketUrl: cmd.Lookup("bitbucket-url"),
    bitbucketUser: cmd.Lookup("bitbucket-user"),
    bitbucketPassword: cmd.Lookup("bitbucket-password"),
  }

  // Clean up the build directory
  os.RemoveAll("test2")

  // Create our channels that will buffer up to x tasks at a time
  // The buffer needs to be big enough so that one repo/file cannot fill it up
  projectChan = make(chan task, 500)
  repoChan = make(chan task, 500)
  fileChan = make(chan task, 500)

  // Start the workers
  workerCount := 20
  for i:=0; i<workerCount; i++ {
      go worker(i, projectChan)
      go worker(i, repoChan)
      go worker(i, fileChan)
  }

  // go worker(1, fileChan)
  // go worker(2, fileChan)

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
		// if p.Key == "SRE" {
      taskProject := p
      task := projectTask{project: &taskProject}

      // Add a count to the waitgroup and add the task to the queue
      wg.Add(1)
      projectChan <- task
		// }
	}

  // We're done adding all the projects, so remove our main blocker so that
  // the program can exit as soon as all the projects are done
  wg.Done()

  // Now wait for all the tasks to finish
  wg.Wait()


  // masterFileList.Range(func(key, value interface{}) bool {
  //   spew.Dump(key)
  //   return true
  // })

	log.Debug("done")
}

// func processRepo(key interface{}, repo interface{}) bool {
//
//   bb, err := bitbucket.NewBitbucketClient(config.bitbucketUrl, config.bitbucketUser, config.bitbucketPassword)
//   if err != nil {
// 		log.Fatal(err)
// 	}
//
//   r := repo.(bitbucket.Repo)
// 	files, err := bb.BrowsePath(r, "/")
// 	if err != nil {
// 		log.Fatal("Error fetching files at / for ", r.Project.Key, "/", r.Slug, " ", err)
// 	}
//
// 	for _, f := range files {
// 		if f.Path.Extension == "md" {
// 			data, err := bb.RawByteSlice(r, f.Path.Name)
// 			if err != nil {
// 				log.Fatal("Error getting markdown file ", f.Path.Name, " in ", r.Project.Key, "/", r.Slug, " ", err)
// 			}
//
// 			// parserOpts := md.Option{}
// 			markdown := md.New(md.WithExtensions(md.CommonExtensions))
// 			parser := markdown.Parse(data)
// 			parser.Walk(func(node *md.Node, entering bool) md.WalkStatus {
// 				return processMarkdownNode(node, entering, &r, &f)
// 			})
// 			// spew.Dump(output)
// 		}
// 	}
//
// 	return true
//
// }

// worker is the main worker function that processes all tasks
func worker(workerNum int, taskChan <-chan task) {
	for task := range taskChan {
		task.run(workerNum)
	}
}
