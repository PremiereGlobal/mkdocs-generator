package bitbucket

import (
  // "github.com/davecgh/go-spew/spew"
  "path/filepath"
)

func (b *BitbucketClient) RawByteSlice(repo Repo, path string) ([]byte, error) {

  data, err := b.get(filepath.Join("projects", repo.Project.Key, "repos", repo.Slug, "raw", path), path)
  if err != nil {
      return nil, err
  }

  return data, nil
}

func (b *BitbucketClient) RawByPath(path string, args string) ([]byte, error) {

  data, err := b.get(path, args)
  if err != nil {
      return nil, err
  }

  return data, nil
}
