package main

import (
	md "gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
)

// fileTask is a type of task that downloads and processes documents
type fileTask struct {

	// document contains the document we want to process
	document *document

	// referencedBy is the document that referenced this file
	// Is only applicable for documents referenced by other documents
	referencedBy *document
}

// run describes how a fileTask should be processed
func (f fileTask) run(workerNum int) bool {

	// Decrement waitgroup counter when we're done
	defer wg.Done()

	log.Debug("Processing file task ", f.document.uid, " [worker:", workerNum, "]")

	// Create new Bitbucket client
	bb := NewBitbucketClient()

	// Download and save the file
	bodyBytes, err := bb.RawByPath(f.document.scmFilePath(), "at=refs/heads/master")
	if err != nil {
		log.Warnf("Error downloading file %s: %v", f.document.scmFilePath(), err)
	} else {
		filename := filepath.Join(Args.GetString("build-dir"), "docs", f.document.scmFilePath())
		CreateFileIfNotExist(filename)
		err = ioutil.WriteFile(filename, bodyBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// If this file is markdown, parse it to find any more linked resources we
	// need to download
	if f.document.docType == markdownType {
		markdown := md.New(md.WithExtensions(md.CommonExtensions))
		parser := markdown.Parse(bodyBytes)
		parser.Walk(func(node *md.Node, entering bool) md.WalkStatus {
			return processMarkdownNode(node, entering, f.document)
		})
	}

	return false
}

// processMarkdownNode processes the markdown item. If we find a link or an image
// and it hasn't already been processed, add it to the file queue
func processMarkdownNode(node *md.Node, entering bool, sourceDocument *document) md.WalkStatus {

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
			var docType docType
			if filepath.Ext(linkURL.Path) == ".md" {
				docType = markdownType
			} else if node.Type == md.Image {
				docType = imageType
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
					sourcePath := filepath.Dir(sourceDocument.scmFilePath())
					referenceMasterFilePath = filepath.Join(sourcePath, linkURL.Path)

				}

				// Generate the document object for this file
				document, err := NewDocumentFromPath(referenceMasterFilePath)
				if err != nil {
					log.Warn("Bad file reference in ", sourceDocument.uid, ": ", err)
					return md.GoToNext
				}
				document.docType = docType

				// Check if the file is on the master list, if not, add it
				if _, ok := masterFileList.Load(document.uid); !ok {

					// Create a task to process this file
					task := fileTask{document: document, referencedBy: sourceDocument}

					// Add the file to the master list so nothing else processes it
					masterFileList.Store(document.uid, document)

					// Add a count to the waitgroup and add the task to the queue
					wg.Add(1)
					taskChan <- task

				}
			}
		}
	}

	return md.GoToNext
}
