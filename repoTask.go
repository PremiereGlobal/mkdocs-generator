package main

import (
  "strings"
  "path/filepath"
  bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
  // "github.com/davecgh/go-spew/spew"
)

type repoTask struct {
	repo bitbucket.Repo
}

func (r repoTask) run(workerNum int) bool {

  // Decrement waitgroup counter when we're done
  defer wg.Done()

  log.Info("Processing repo task ", r.repo.MakePath(), " [repo-worker:", workerNum, "]")

  // Create new connection to Bitbucket
  bb, err := bitbucket.NewBitbucketClient(config.bitbucketUrl, config.bitbucketUser, config.bitbucketPassword)
  if err != nil {
		log.Fatal(err)
	}

  // Paths to browse for markdown
  browsePaths := []string{"/", "/docs"}

  for _, path := range browsePaths {

    // spew.Dump(path)

    // Get the list of files at this path
    // There's a good chance that there's not a docs folder so don't do anything here...
  	browseList, err := bb.Browse(&r.repo, path)
  	if err != nil {
  		// log.Warn("Error fetching files at ", path, " for ", r.repo.MakePath(), " ", err)
      continue
  	}

    // spew.Dump(browseList.Children.Values[0])
    // spew.Dump(browseList)
    // Loop over the files
    // values := browseList.Children.Values.([]bitbucket.File)
  	for _, f := range browseList.Children.Values {
      // file := f.(bitbucket.File)
      // spew.Dump(filepath.Join(r.repo.MakePath(), browseList.Path.ToString, f.Path.ToString))
  		if strings.ToLower(f.Path.Extension) == "md" {
    //
        // Generate the master file path for this file
        masterFilePath := filepath.Join(r.repo.MakePath(), "raw", filepath.Join(browseList.Path.ToString, f.Path.ToString))

        // Create a task to process this file
        task := fileTask{
          // repoPath: r.repo.MakePath(),
          // filePath: filepath.Join(browseList.Path.ToString, f.Path.ToString),
          // fileName: filepath.Join(r.repo.MakePath(), "raw", browseList.Path.ToString, f.Path.ToString),
          // writeFileName: filepath.Join(r.repo.Project.Key, r.repo.Slug, browseList.Path.ToString, f.Path.ToString),
          fileType: "markdown",
          masterFilePath: masterFilePath,
        }

        // Add the file to the master list so nothing else processes it
        masterFileList.Store(strings.ToLower(masterFilePath), true)
    //     //
        // Add a count to the waitgroup and add the task to the queue
        wg.Add(1)
        fileChan <- task
    //
    //
  	// 		// data, err := bb.RawByteSlice(r.repo, f.Path.Name)
  	// 		// if err != nil {
  	// 		// 	log.Fatal("Error getting markdown file ", f.Path.Name, " in ", r.repo.Project.Key, "/", r.repo.Slug, " ", err)
  	// 		// }
    //     //
  	// 		// markdown := md.New(md.WithExtensions(md.CommonExtensions))
  	// 		// parser := markdown.Parse(data)
  	// 		// parser.Walk(func(node *md.Node, entering bool) md.WalkStatus {
  	// 		// 	return processMarkdownNode(node, entering, &r.repo, &f)
  	// 		// })
  		}
  	}

  }


  //
  // // Get the list of files in the "docs" directory of this project
  // docfiles, err := bb.Browse(&r.repo, "/docs")
	// if err != nil {
  //   // There's a good chance that there's not a docs folder so don't do anything here...
	// }
  //
  // // Combine the list of markdown from the root and docs directories
  // files = append(files, docfiles...)



  return true
}
