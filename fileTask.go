package main

import (
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
	"unicode/utf8"

	bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
	md "gopkg.in/russross/blackfriday.v2"
)

// fileTask is a type of task that downloads and processes documents
type fileTask struct {
	file bitbucket.BBFile
	// document contains the document we want to process
	document *document

	// referencedBy is the document that referenced this file
	// Is only applicable for documents referenced by other documents
	referencedBy *document
}

// run describes how a fileTask should be processed
func (f fileTask) run(workerNum int, taskChan chan<- task) bool {

	log.Debugf("[worker:%03d] Processing file task %s:%s", workerNum, f.document.uid, f.file.GetName())

	// Download and save the file
	bodyBytes, err := f.file.GetData()
	if err != nil {
		log.Warnf("[worker:%03d] Error downloading file %s: %v", workerNum, f.document.scmFilePath(), err)
		return false
	}
	if f.document.docType == markdownType {
		if !utf8.Valid(bodyBytes) {
			//validate that the file is valid utf8 and/or ascii or mkdocs cant parse it
			log.Warnf("[worker:%03d] invalid utf8 on repo: %s", workerNum, f.document.uid)
			return false
		}
	}
	filename := filepath.Join(Args.GetString("build-dir"), "docs", f.document.scmFilePath())
	CreateFileIfNotExist(filename)
	ioutil.WriteFile(filename, bodyBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// If this file is markdown, parse it to find any more linked resources we
	// need to download
	if f.document.docType == markdownType {
		markdown := md.New(md.WithExtensions(md.CommonExtensions))
		parser := markdown.Parse(bodyBytes)
		parser.Walk(func(node *md.Node, entering bool) md.WalkStatus {
			return processMarkdownNode(node, entering, f.document, taskChan, workerNum)
		})
	}

	return false
}

// processMarkdownNode processes the markdown item. If we find a link or an image
// and it hasn't already been processed, add it to the file queue
func processMarkdownNode(node *md.Node, entering bool, sourceDocument *document, taskChan chan<- task, workerNum int) md.WalkStatus {

	// Since this gets called twice, only execute on the entry event
	if entering == true {
		bbrepo := sourceDocument.getBBFile().GetBBRepo()
		// We only care about links and images
		if node.Type == md.Link || node.Type == md.Image {

			// Parse the reference so we can get the parts we need
			// If it can't be parsed, just continue
			linkURL, err := url.Parse(string(node.LinkData.Destination))
			if err != nil {
				log.Warnf("[worker:%03d] Unable to parse markdown reference %s", workerNum, string(node.LinkData.Destination))
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
			u, err := url.Parse(sourceDocument.getBBFile().GetBBRepo().GetBBProject().GetBBClient().BaseUrl.String())
			if err != nil {
				log.Fatal(err)
			}

			log.Infof("[worker:%03d] Found link:%s in [%s]%s", workerNum, linkURL.String(), sourceDocument.getBBFile().GetBBRepo().GetName(), sourceDocument.getBBFile().GetFullPath())
			var newBBFile bitbucket.BBFile
			// Continue only if reference is a relative link or from the same Bitbucket host
			if (linkURL.Host == u.Host || (linkURL.Scheme == "" && linkURL.Host == "")) && linkURL.Path != "" {
				// If our path starts with a /, we don't need to add the project/repo info
				if strings.HasPrefix(linkURL.Path, "/") {
					// Get rid of the leading slash
					bbf, err := bbrepo.GetFile(linkURL.Path)
					if err != nil {
						log.Warnf("[worker:%03d] Found link:%s, But Got Error:%s", workerNum, linkURL.String(), err)
						return md.GoToNext
					}
					newBBFile = bbf
				} else {
					// Path is relative, use the masterFilePath directory to generate
					// the master file path for the reference file
					sourcePath := filepath.Dir(sourceDocument.getBBFile().GetBasePath())
					fp := filepath.Join(sourcePath, linkURL.Path)
					bbf, err := bbrepo.GetFile(fp)
					if err != nil {
						log.Warnf("[worker:%03d] Found2 link:%s, But Got Error:%s", workerNum, linkURL.String(), err)
						return md.GoToNext
					}
					newBBFile = bbf

				}
				if newBBFile == nil {
					return md.GoToNext
				}
				// Generate the document object for this file
				document := NewDocument(bbrepo.GetBBProject().GetKey(), bbrepo.GetSlug(), newBBFile.GetFullPath(), newBBFile)
				document.docType = docType

				// Check if the file is on the master list, if not, add it
				if _, ok := masterFileList.LoadOrStore(document.uid, document); !ok {
					// Create a task to process this file
					task := fileTask{document: document, referencedBy: sourceDocument, file: newBBFile}
					taskChan <- task
				} else {
					log.Infof("[worker:%03d] Skipped, already processed link:%s", workerNum, linkURL.String())
				}
			}
		}
	}

	return md.GoToNext
}
