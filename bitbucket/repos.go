package bitbucket

import (
  // "github.com/davecgh/go-spew/spew"
  "path/filepath"
  "encoding/json"
  "fmt"
)

type RepoList struct {
  PaginatedList
  Values []Repo
}

type Repo struct {
  Slug string
  Id int
  Name string
  ScmId string
  State string
  StatusMessage string
  Forkable bool
  Project Project
  Public bool
  Links RepoLinks
}

type RepoLinks struct {
  Clone []RepoClone
  Self []RepoSelf
}

type RepoClone struct {
  Href string
  Name string
}

type RepoSelf struct {
  Href string
}

// MakePath creates the full relative path to the repo
func (r *Repo) MakePath() string {
  return filepath.Join("projects", r.Project.Key, "repos", r.Slug)
}

func (b *BitbucketClient) ListRepos(project *Project) (*RepoList, error) {

  repoList := RepoList{}

  body, err := b.get(filepath.Join(project.MakePath(), "repos"), fmt.Sprintf("limit=%d", b.limit))
  if err != nil {
    return nil, err
  }

  err = json.Unmarshal(body, &repoList)
  if err != nil {
    return nil, err
  }

  return &repoList, nil
}
