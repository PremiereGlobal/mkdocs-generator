package main

import (
	bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
	md "gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type fileTask struct {

	// Filetype should be "markdown" or "image"
	fileType string

	masterFilePath string
}

func (f fileTask) run(workerNum int) bool {

	// Decrement waitgroup counter when we're done
	defer wg.Done()

	log.Debug("Processing file task ", f.masterFilePath, " [file-worker:", workerNum, "]")

	// Create new Bitbucket client
	bb, err := bitbucket.NewBitbucketClient(config.bitbucketUrl, config.bitbucketUser, config.bitbucketPassword)
	if err != nil {
		log.Fatal(err)
	}

	// If the file is refrenced by the /browse/ path, we need to convert this
	// to "raw" for downloading
	downloadPath := f.masterFilePath
	parts := strings.Split(f.masterFilePath, string(os.PathSeparator))
	if len(parts) >= 6 && parts[4] == "browse" {
		buildPath := parts[0:4]
		buildPath = append(buildPath, "raw")
		buildPath = append(buildPath, parts[5:]...)
		downloadPath = filepath.Join(buildPath...)
	}

	// Download and save the file
	bodyBytes, err := bb.RawByPath(downloadPath, "at=refs/heads/master")
	if err != nil {
		log.Warn(err)
	} else {
		filename := filepath.Join(Args.GetString("build-dir"), f.masterFilePath)
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

	return false
}

type task interface {
	run(int) bool
}

// processMarkdownNode processes the markdown item. If we find a link or an image
// and it hasn't already been processed, add it to the file queue
func processMarkdownNode(node *md.Node, entering bool, sourceMasterFilePath string) md.WalkStatus {

	// Since this gets called twice, only execute on the entry event
	if entering == true {

		// We only care about links and images
		if node.Type == md.Link || node.Type == md.Image {

			// Parse the reference so we can get the parts we need
			// If it can't be parsed, just continue
			linkURL, err := url.Parse(string(node.LinkData.Destination))
			if err != nil {
				log.Warn("Unable to parse markdown reference ", string(node.LinkData.Destination))
				return md.GoToNext
			}

			// Determine what type of file we're dealing with and exit here if it's
			// not markdown or image
			fileType := ""
			if strings.ToLower(filepath.Ext(linkURL.Path)) == ".md" {
				fileType = "markdown"
			} else if node.Type == md.Image {
				fileType = "image"
			} else {
				return md.GoToNext
			}

			// Get our Bitbucket URL ready
			u, err := url.Parse(config.bitbucketUrl)
			if err != nil {
				log.Fatal(err)
			}

			// Continue nly if reference is a relative link or from the same Bitbucket host
			if (linkURL.Host == u.Host || (linkURL.Scheme == "" && linkURL.Host == "")) && linkURL.Path != "" {

				// If our path starts with a /, we don't need to add the project/repo info
				referenceMasterFilePath := ""
				if strings.HasPrefix(linkURL.Path, "/") {

					// Get rid of the leading slash
					referenceMasterFilePath = linkURL.Path[1:]

				} else {

					// Path is relative, use the masterFilePath directory to generate
					// the master file path for the reference file
					sourcePath := filepath.Dir(sourceMasterFilePath)
					referenceMasterFilePath = filepath.Join(sourcePath, linkURL.Path)
				}

				// Check if the file is on the master list, if not, add it
				if _, ok := masterFileList.Load(strings.ToLower(referenceMasterFilePath)); !ok {

					// Add the file to the master list, then add it to the queue
					masterFileList.Store(strings.ToLower(referenceMasterFilePath), true)
					task := fileTask{
						fileType:       fileType,
						masterFilePath: referenceMasterFilePath,
					}

					// Add the file to the master list so nothing else processes it
					masterFileList.Store(strings.ToLower(referenceMasterFilePath), true)

					// Add a count to the waitgroup and add the task to the queue
					wg.Add(1)
					fileChan <- task

				}
			}
		}
	}

	return md.GoToNext
}
