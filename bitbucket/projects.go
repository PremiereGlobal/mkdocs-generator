package bitbucket

import (
  // "github.com/davecgh/go-spew/spew"
  "path/filepath"
  "encoding/json"
  "fmt"
)

type ProjectList struct {
  PaginatedList
  Values []Project
}

type Project struct {
  Key string
  Id int
  Name string
  Description string
  Public bool
  ProjectType string `json:"type"`
}

// MakePath creates the full relative path to the repo
func (p *Project) MakePath() string {
  return filepath.Join("projects", p.Key)
}

// ListProjects returns a list of all projects on the Bitbucket server
func (b *BitbucketClient) ListProjects() (*ProjectList, error) {

  projectList := ProjectList{}

  body, err := b.get("projects", fmt.Sprintf("limit=%d", b.limit))
  if err != nil {
    return nil, err
  }

  err = json.Unmarshal(body, &projectList)
  if err != nil {
    return nil, err
  }

  return &projectList, nil
}
