package main

import (
	bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
	"path/filepath"
  "io/ioutil"
  "net/url"
  "strings"
  "os"
  md "gopkg.in/russross/blackfriday.v2"
  // "github.com/davecgh/go-spew/spew"
)

type fileTask struct {

  // // The Bitbucket repo this file is part of
  // repo bitbucket.Repo
  //
  // // The basepath of the file, i.e. projects/key/repos/slug
  // repoPath string
  //
  // // The relative of the file, i.e. docs/dir/dir/test.md
  // filePath string

  // Path to download the file
	// fileName string

  // Path to write the file
  // writeFileName string

  // Filetype should be "markdown" or "image"
  fileType string

  //
  masterFilePath string
}

func (f fileTask) run(workerNum int) bool {

  // Decrement waitgroup counter when we're done
	defer wg.Done()

	log.Info("Processing file task ", f.masterFilePath, " [file-worker:", workerNum, "]")

  // // Check if the file is on the master list, if not, add it and process it
  // if _, ok := masterFileList.Load(f.writeFileName); !ok {
  //
  //   // Add the file to the master list, then add it to the queue
  //   // log.Debug("Adding file task for ", referencePath)
  //   masterFileList.Store(f.writeFileName, true)

  // Create new Bitbucket client
	bb, err := bitbucket.NewBitbucketClient(config.bitbucketUrl, config.bitbucketUser, config.bitbucketPassword)
	if err != nil {
		log.Fatal(err)
	}

  // in some cases if the file is refrenced by the /browse/ path, we need to convert this
  // to "raw" for downllading
  downloadPath := f.masterFilePath
  parts := strings.Split(f.masterFilePath, string(os.PathSeparator))
  // spew.Dump(parts)
 if len(parts) >= 6 && parts[4] == "browse" {
   // downloadPath := filepath.Join(path.Join(parts[0:3], "raw", parts[]...))
   buildPath := parts[0:4]// basePath := parts[0:3]
   buildPath = append(buildPath, "raw")// endPath := parts[5:]
   buildPath = append(buildPath, parts[5:]...)

   downloadPath = filepath.Join(buildPath...)
   // log.Warn("replacing browse path")
 }

  // Download and save the file
  bodyBytes, err := bb.RawByPath(downloadPath, "at=refs/heads/master")
	if err != nil {
		log.Warn(err)
	} else {
		filename := filepath.Join("test2/mkdocs", f.masterFilePath)
		CreateFileIfNotExist(filename)
		err = ioutil.WriteFile(filename, bodyBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

  // If this file is markdown, parse it to find any more linked resources we
  // need to download
  if f.fileType == "markdown" {
		markdown := md.New(md.WithExtensions(md.CommonExtensions))
		parser := markdown.Parse(bodyBytes)
		parser.Walk(func(node *md.Node, entering bool) md.WalkStatus {
			return processMarkdownNode(node, entering, f.masterFilePath)
		})
	}

  // } else {
  //   log.Debug("Skipping duplicate file ", f.writeFileName)
  // }

	return false
}

type task interface {
	run(int) bool
}

// processMarkdownNode processes the markdown item. If we find a link or an image
// and it hasn't already been processed, add it to the list
func processMarkdownNode(node *md.Node, entering bool, sourceMasterFilePath string) md.WalkStatus {

	// Since this gets called twice, only execute on the entry event
	if entering == true {

		// // If this is an image, add it to the queue (if it doesn't exist)
    // if  {
    //
    //   // Check if the file is on the master list, if not, add it
    //   if _, ok := masterFileList.Load(referencePath); !ok {
    //     // Add the file to the master list, then add it to the queue
    //     // log.Debug("Adding file task for ", referencePath)
    //     masterFileList.Store(referencePath, true)
    //     task := fileTask{
    //       fileName: referencePath,
    //       fileType: "markdown",
    //     }
    //
    //     // Add 1 to the wait group
    //     wg.Add(1)
    //
    //     // Add the task to the queue
    //     fileChan <- task
    //   }
    // }

    // We only care about links and images
		if node.Type == md.Link || node.Type == md.Image {

			// Parse the reference so we can get the parts we need
      // If it can't be parsed, just continue
      // spew.Dump(string(node.LinkData.Destination))
			linkURL, err := url.Parse(string(node.LinkData.Destination))
			if err != nil {
        log.Warn("Unable to parse markdown reference ", string(node.LinkData.Destination))
				return md.GoToNext
			}

      // Determine what type of file we're dealing with and exit here if it's
      // not markdown or image
      // spew.Dump(linkURL.Path)
      // spew.Dump(strings.ToLower(filepath.Ext(linkURL.Path)))
      fileType := ""
      if strings.ToLower(filepath.Ext(linkURL.Path)) == ".md" {
        fileType = "markdown"
      } else if node.Type == md.Image {
        fileType = "image"
      } else {
        return md.GoToNext
      }

      // spew.Dump(fileType)

      // Get our Bitbucket URL ready
      u, err := url.Parse(config.bitbucketUrl)
      if err != nil {
				log.Fatal(err)
			}

      // Continue nly if reference is a relative link or from the same Bitbucket host
			if (linkURL.Host == u.Host || (linkURL.Scheme == "" && linkURL.Host == "")) && linkURL.Path != "" {

				// If our path starts with a /, we don't need to add the project/repo info
				// referencePath is the full path of the file, from the root of the site
        // Also, paths need to be lowercased in case there is variation however
        // the filenames themselves must remain in their original case
				// referencePath := ""
        // referenceRepoPath := ""
        // referenceFilePath := ""
        referenceMasterFilePath := ""
        // log.Debug("Found link to ", linkURL.Path)
				if strings.HasPrefix(linkURL.Path, "/") {

          // Split the link up to make sure we've got the right parts so we can
          // construct the new master file path
          // parts := filepath.SplitList(linkURL.Path)
          // if len(parts) >= 6 && parts[0] == "projects" && parts[2] == "repos" && (parts[4] == "browse" || parts[4] == "raw") {
          //   filepath.Join(strings.ToLower(filepath.Join(parts[0], parts[1], parts[2], parts[3], parts[4]))
          // }

          // Get rid of the leading slash
          // referencePath = filepath.Join(linkURL.Path[1:])
					// referenceRepoPath = strings.ToLower(filepath.Dir(referencePath))
          // referenceFilePath = filepath.Base(referencePath))
          referenceMasterFilePath = linkURL.Path[1:]
          // log.Debug("Reference is absolute; master file path: ", referenceMasterFilePath)
				} else {
          // Path is relative, use the masterFilePath directory to generate
          // the master file path for the reference file
          sourcePath := filepath.Dir(sourceMasterFilePath)
          referenceMasterFilePath = filepath.Join(sourcePath, linkURL.Path)
          // referencePath = filepath.Join(repoPath, "raw", filepath.Dir(filePath), linkURL.Path)
          // referenceRepoPath =
          // referenceFilePath = filepath.Base(referencePath))

          // log.Debug("Reference is relative; master file path: ", referenceMasterFilePath)
				}
				// referencePath = strings.ToLower(referencePath)



  			// Check if the file is on the master list, if not, add it
        // masterName := filepath.Join(repoPath, "raw", filePath)
				if _, ok := masterFileList.Load(strings.ToLower(referenceMasterFilePath)); !ok {
					// Add the file to the master list, then add it to the queue
					// log.Debug("Adding file task for ", referenceMasterFilePath, " based on markdown reference")
					masterFileList.Store(strings.ToLower(referenceMasterFilePath), true)
					task := fileTask{
            fileType: fileType,
            masterFilePath: referenceMasterFilePath,
          }

          // Add the file to the master list so nothing else processes it
          masterFileList.Store(strings.ToLower(referenceMasterFilePath), true)

          // Add 1 to the wait group
          wg.Add(1)
        //
          // Add the task to the queue
					fileChan <- task
				} else {
          // log.Debug("Reference has already been processed ", referenceMasterFilePath)
        }

			} else {
        // log.Debug("Reference is not in Bitbucket ", string(node.LinkData.Destination))
      }
		}
	}

	return md.GoToNext
}
