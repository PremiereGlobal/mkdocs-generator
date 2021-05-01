package bitbucket

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
)

func (b *BitbucketClient) get(path string, query string) ([]byte, error) {
	u, _ := url.Parse(b.BaseUrl.String())
	u.Path = filepath.Join(b.BaseApiPath, path, "/")
	u.RawQuery = query
	return b.getURL(u)
}

func (b *BitbucketClient) getURL(gurl *url.URL) ([]byte, error) {
	request := &http.Request{Method: "GET", URL: gurl, Header: http.Header{}}
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
