package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// docType specifies a type of a single document
type docType int

// Possible values of docType
const (
	markdownType docType = iota
	imageType
)

// document contains information about a single document
type document struct {

	// uid is the unique identifier for this document
	// this is constructed as follows
	// strings.ToLower(<project>/<repo>/<filepath>)
	uid string

	// docType describes the type of document (markdown, image, etc.)
	docType docType

	project  string
	repo     string
	filePath string
}

// NewDocument generates a new document object from its components
func NewDocument(project string, repo string, filePath string) *document {
	document := document{
		project:  strings.ToUpper(project),
		repo:     repo,
		filePath: filePath,
	}
	document.uidGen()

	return &document
}

// NewDocumentFromPath generates a new document object from a given scmPath
func NewDocumentFromPath(scmPath string) (*document, error) {

	var (
		project  string
		repo     string
		filePath string
	)

	// Break the path out into its components
	pathParts := strings.Split(scmPath, string(os.PathSeparator))

	// Ensure we have the right components to make up a document
	if pathParts[0] != "projects" || pathParts[2] != "repos" || !(pathParts[4] == "raw" || pathParts[4] == "browse") {
		return nil, errors.New(fmt.Sprintf("Reference to a bad Bitbucket file %s", scmPath))
	} else {
		project = pathParts[1]
		repo = pathParts[3]
		filePath = filepath.Join(pathParts[5:]...)
	}

	return NewDocument(project, repo, filePath), nil
}

func (d *document) uidGen() {
	slug := filepath.Join(d.project, d.repo, d.filePath)
	d.uid = strings.ToLower(slug)
}

func (d *document) scmRepoPath() string {
	return filepath.Join("projects", d.project, "repos", d.repo)
}

func (d *document) scmFilePath() string {
	return filepath.Join(d.scmRepoPath(), "raw", d.filePath)
}
