package bitbucket

import (
  "fmt"
  "errors"
  // "encoding/json"
  "io/ioutil"
  "net/http"
  // "net/url"
  "path/filepath"
  // "github.com/davecgh/go-spew/spew"
)

type PaginatedList struct {
  Size int
  Limit int
  IsLastPage bool
  Start int
  // Values interface{}
}

type Path struct {
  Components []string
  Parent string
  Name string
  Extension string
  ToString string
}

func (b *BitbucketClient) Get(path string, query string) ([]byte, error) {
  return b.get(path, query)
}

func (b *BitbucketClient) get(path string, query string) ([]byte, error) {

  u := b.BaseUrl
  u.Path = filepath.Join(b.BaseApiPath, path, "/")
  u.RawQuery = query

  // b.log.Debug("Fetching ", u.String())

  request := &http.Request{Method: "GET", URL: u, Header: http.Header{}}
  request.SetBasicAuth(b.Username, b.Password)

  // spew.Dump("xx"+ request.URL.String())
  // httpClient := &http.Client{}
  response, err := b.rawClient.Do(request)
  if err != nil {
    return nil, err
  }
  // spew.Dump(response.Request)
  defer response.Body.Close()

  if response.StatusCode == http.StatusOK {
    bodyBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    return bodyBytes, nil
  }

  // spew.Dump(response.Request, response.Header)
  return nil, errors.New(fmt.Sprintf("%d, %s", response.StatusCode, response.Request.URL))
}

// getList takes in an API path and list object and populates the list with
// the requested items
// func (b *BitbucketClient) getPaginatedList(path string, list interface{}) (error) {
//
//   responseList := &PaginatedList{Values: list}
//
//   body, err := b.get(path, fmt.Sprintf("limit=%d", b.limit))
//   if err != nil {
//     return err
//   }
//
//   err = json.Unmarshal(body, responseList)
//   if err != nil {
//     return err
//   }
//
//   return nil
// }
