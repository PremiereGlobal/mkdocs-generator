package bitbucket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
)

type ProjectBBServerList struct {
	Size       int
	Limit      int
	IsLastPage bool
	Start      int
	Values     []*ProjectBBServer
}

func (pl *ProjectBBServerList) GetValues() []BBProject {
	b := make([]BBProject, len(pl.Values))
	for i := range pl.Values {
		b[i] = pl.Values[i]
	}
	return b
}

func (b *BitbucketClient) listProjectsBBServer() ([]BBProject, error) {

	projectList := &ProjectBBServerList{}

	body, err := b.get("projects", fmt.Sprintf("limit=%d", b.limit))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, projectList)
	if err != nil {
		return nil, err
	}
	for _, v := range projectList.Values {
		v.client = b
	}

	return projectList.GetValues(), nil
}

type ProjectBBServer struct {
	Key         string
	Id          int
	Name        string
	Description string
	Public      bool
	ProjectType string `json:"type"`
	client      *BitbucketClient
}

func (p *ProjectBBServer) GetKey() string {
	return p.Key
}
func (p *ProjectBBServer) GetId() string {
	return fmt.Sprintf("%d", p.Id)
}
func (p *ProjectBBServer) GetName() string {
	return p.Name
}
func (p *ProjectBBServer) GetDescription() string {
	return p.Description
}
func (p *ProjectBBServer) IsPublic() bool {
	return p.Public
}
func (p *ProjectBBServer) GetProjectType() string {
	return p.ProjectType
}
func (p *ProjectBBServer) ListReposURL() *url.URL {
	tmpURL := p.GetBBClient().BaseUrl.String() + "/" + filepath.Join(p.GetBBClient().BaseApiPath, "projects", p.Key, "repos") + fmt.Sprintf("?limit=%d", p.client.limit)
	rurl, err := url.Parse(tmpURL)
	if err != nil {
		return nil
	}
	return rurl
}
func (p *ProjectBBServer) ListRepos() ([]BBRepo, error) {
	repoList := &RepoBBServerList{}

	body, err := p.client.getURL(p.ListReposURL())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, repoList)
	if err != nil {
		return nil, err
	}
	for _, v := range repoList.Values {
		v.Project = p
	}

	return repoList.GetValues(), nil
}
func (p *ProjectBBServer) GetBBClient() *BitbucketClient {
	return p.client
}

type RepoBBServerList struct {
	Size       int
	Limit      int
	IsLastPage bool
	Start      int
	Values     []*RepoBBServer
}

func (rl *RepoBBServerList) GetValues() []BBRepo {
	b := make([]BBRepo, len(rl.Values))
	for i := range rl.Values {
		b[i] = rl.Values[i]
	}
	return b
}

type RepoBBServer struct {
	Slug          string
	Id            int
	Name          string
	ScmId         string
	State         string
	StatusMessage string
	Forkable      bool
	Project       *ProjectBBServer `json:"-"`
	Public        bool
	// Links         RepoLinks
}

func (r *RepoBBServer) GetSlug() string {
	return r.Slug
}
func (r *RepoBBServer) GetId() string {
	return fmt.Sprintf("%d", r.Id)
}
func (r *RepoBBServer) GetName() string {
	return r.Name
}
func (r *RepoBBServer) GetScmId() string {
	return r.ScmId
}
func (r *RepoBBServer) GetState() string {
	return r.State
}
func (r *RepoBBServer) GetStatusMessage() string {
	return r.StatusMessage
}
func (r *RepoBBServer) GetForkable() bool {
	return r.Forkable
}
func (r *RepoBBServer) GetBBProject() BBProject {
	return r.Project
}
func (r *RepoBBServer) IsPublic() bool {
	return r.Public
}
func (r *RepoBBServer) GetFilesURL() *url.URL {
	furl, err := url.Parse(r.GetBBProject().GetBBClient().BaseUrl.String() + "/" + filepath.Join(r.GetBBProject().GetBBClient().BaseApiPath, "projects", r.Project.GetKey(), "repos", r.Slug, "browse") + fmt.Sprintf("?limit=%d", r.Project.client.limit))
	if err != nil {
		return nil
	}
	return furl
}
func (r *RepoBBServer) GetDir(path string) ([]BBFile, error) {

	var browseList BBServerBrowseList
	furl := r.GetFilesURL()
	furl.Path = filepath.Join(furl.Path, path)
	body, err := r.Project.client.getURL(furl)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &browseList)
	if err != nil {
		return nil, err
	}
	fl := make([]BBFile, 0)
	for _, v := range browseList.Children.Values {
		v.basePath = path
		v.repo = r
		fl = append(fl, v)
	}
	return fl, nil
}
func (r *RepoBBServer) GetFile(path string) (BBFile, error) {
	fileName := filepath.Base(path)
	filePath := filepath.Dir(path)

	var browseList BBServerBrowseList
	furl := r.GetFilesURL()
	furl.Path = filepath.Join(furl.Path, filePath)
	body, err := r.Project.client.getURL(furl)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &browseList)
	if err != nil {
		return nil, err
	}
	for _, v := range browseList.Children.Values {
		if v.GetName() == fileName {
			v.basePath = filePath
			v.repo = r
			return v, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("File Not Found:%s", path))
}

type BBServerBrowseList struct {
	Path     Path
	Revision string
	Children BBServerBrowsePaginatedList
}

type BBServerBrowsePaginatedList struct {
	Size       int
	Limit      int
	IsLastPage bool
	Start      int
	Values     []*BBServerFile
}

type BBServerFile struct {
	Path      Path
	ContentId string
	FileType  string `json:"type"`
	Size      int
	basePath  string
	repo      BBRepo
}

type Path struct {
	Components []string
	Parent     string
	Name       string
	Extension  string
	ToString   string
}

func (bbf *BBServerFile) GetFullPath() string {
	return filepath.Join(bbf.basePath, bbf.Path.Name)
}

func (bbf *BBServerFile) GetBasePath() string {
	return bbf.basePath
}

func (bbf *BBServerFile) GetName() string {
	return bbf.Path.Name
}

func (bbf *BBServerFile) GetContentId() string {
	return bbf.ContentId
}

func (bbf *BBServerFile) GetFileType() string {
	return bbf.FileType
}

func (bbf *BBServerFile) GetSize() int {
	return bbf.Size
}
func (bbf *BBServerFile) GetBBRepo() BBRepo {
	return bbf.repo
}

func (bbf *BBServerFile) GetURL() (*url.URL, error) {
	r := bbf.GetBBRepo()
	furl, err := url.Parse(r.GetBBProject().GetBBClient().BaseUrl.String())
	if err != nil {
		return nil, err
	}
	//Must set to .Path here to get correct escaping
	furl.Path = "/" + filepath.Join(r.GetBBProject().GetBBClient().BaseApiPath, "projects", r.GetBBProject().GetKey(), "repos", r.GetSlug(), "raw", bbf.GetFullPath())
	return furl, nil
}

func (bbf *BBServerFile) GetData() ([]byte, error) {
	r := bbf.GetBBRepo()
	furl, err := bbf.GetURL()
	if err != nil {
		return nil, err
	}
	ba, err := r.GetBBProject().GetBBClient().getURL(furl)
	if err != nil {
		return nil, err
	}

	return ba, nil
}
