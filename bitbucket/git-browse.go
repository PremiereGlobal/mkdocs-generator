package bitbucket

import (
  // "github.com/davecgh/go-spew/spew"
  "path/filepath"
    "encoding/json"
    "fmt"
    // "reflect"
)

type BrowseList struct {
  Path Path
  Revision string
  Children BrowsePaginatedList
}

type BrowsePaginatedList struct {
  PaginatedList
  Values []File
}

type File struct {
  Path Path
  ContentId string
  FileType string `json:"type"`
  Size int
}

func (b *BitbucketClient) Browse(repo *Repo, path string) (*BrowseList, error) {

  var browseList BrowseList
  // body, err := b.get(filepath.Join(repo.MakePath(), "browse"), fmt.Sprintf("limit=%d", b.limit))

  // var fileList []File
  // spew.Dump(filepath.Join(repo.MakePath(), "browse", path))
  // browseList.Children = &PaginatedList{Values: fileList}
  body, err := b.get(filepath.Join(repo.MakePath(), "browse", path), fmt.Sprintf("limit=%d", b.limit))
  if err != nil {
    return nil, err
  }

  err = json.Unmarshal(body, &browseList)
  if err != nil {
    return nil, err
  }
  // browseList, err := b.getBrowseList(repo.MakePath(), &fileList)
  // if err != nil {
  //   return nil, err
  // }
  //
  // fmt.Println(reflect.TypeOf(browseList.Children.Values).String())

  // test := browseList.Children.Values


  // browseList.Children.Values = browseList.Children.Values.([]File)

  return &browseList, nil
}
//
// func (b *BitbucketClient) getBrowseList(path string, list interface{}) (*BrowseList, error) {
//
//   browseList := BrowseList{Children: PaginatedList{Values: &[]File{}}}
//   // var fileList []File
//
//   // paginatedList := PaginatedList{Values: &fileList}
//   // browseList.Children = paginatedList
//
//
//
//   err = json.Unmarshal(body, &browseList)
//   if err != nil {
//     return nil, err
//   }
//
//   spew.Dump(browseList.Children.Values)
//   spew.Dump("A")
//
//   x := *browseList.Children.Values.(*[]File)
//
// //   values := browseList.Children.Values.(*[]File)
// //   browseList.Children.Values = values
// //
// fmt.Println(reflect.TypeOf(x[0]).String())
//
//   // spew.Dump(browseList)
//
//   return &browseList, nil
// }


// func (b *BitbucketClient) BrowsePath(repo Repo, path string) ([]File, error) {
//
//   fileList := []File{}
//
//   err := b.browse(filepath.Join("projects", repo.Project.Key, "repos", repo.Slug, "browse", path), &fileList)
//   if err != nil {
//     return nil, err
//   }
//
//   return fileList, nil
// }

//
// func (b *BitbucketClient) browse(path string, list interface{}) (error) {
//
//   responseList := &BrowseList{}
//   responseList.Children.Values = list
//
//   // Setting this to a high value statically for now
//   // Should probably get pagination working...
//   limit := 1000
//
//   body, err := b.get(path, fmt.Sprintf("limit=%d", limit))
//   if err != nil {
//     // b.log.Debug(body)
//     return err
//   }
//
//   err = json.Unmarshal(body, responseList)
//   if err != nil {
//     // b.log.Info(string(body))
//     // b.log.Fatal(path)
//     // b.log.Debug(string(body))
//     return err
//   }
//
//   return nil
// }
