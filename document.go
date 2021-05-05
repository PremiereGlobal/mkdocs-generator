package main

import (
	"path/filepath"
	"strings"

	bitbucket "github.com/PremiereGlobal/mkdocs-generator/bitbucket"
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
	bbfile   bitbucket.BBFile
}

// NewDocument generates a new document object from its components
func NewDocument(project string, repo string, filePath string, bbfile bitbucket.BBFile) *document {
	document := document{
		project:  strings.ToUpper(project),
		repo:     repo,
		filePath: filePath,
		bbfile:   bbfile,
	}
	document.uidGen()

	return &document
}

func (d *document) getBBFile() bitbucket.BBFile {
	return d.bbfile
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
