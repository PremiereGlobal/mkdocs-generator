package bitbucket

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

type BrowseList struct {
	Path     Path
	Revision string
	Children BrowsePaginatedList
}

type BrowsePaginatedList struct {
	paginatedList
	Values []File
}

type File struct {
	Path      Path
	ContentId string
	FileType  string `json:"type"`
	Size      int
}

func (b *BitbucketClient) Browse(repo *Repo, path string) (*BrowseList, error) {

	var browseList BrowseList

	body, err := b.get(filepath.Join(repo.MakePath(), "browse", path), fmt.Sprintf("limit=%d", b.limit))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &browseList)
	if err != nil {
		return nil, err
	}

	return &browseList, nil
}
