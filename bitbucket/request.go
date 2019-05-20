package bitbucket

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type paginatedList struct {
	Size       int
	Limit      int
	IsLastPage bool
	Start      int
}

// Path element for
type Path struct {
	Components []string
	Parent     string
	Name       string
	Extension  string
	ToString   string
}

func (b *BitbucketClient) Get(path string, query string) ([]byte, error) {
	return b.get(path, query)
}

func (b *BitbucketClient) get(path string, query string) ([]byte, error) {

	u := b.BaseUrl
	u.Path = filepath.Join(b.BaseApiPath, path, "/")
	u.RawQuery = query

	request := &http.Request{Method: "GET", URL: u, Header: http.Header{}}
	request.SetBasicAuth(b.config.Username, b.config.Password)

	response, err := b.rawClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return bodyBytes, nil
	}

	return nil, errors.New(fmt.Sprintf("%d, %s", response.StatusCode, response.Request.URL))
}
