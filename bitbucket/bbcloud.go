package bitbucket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type projectBBCloudResp struct {
	Page     int
	Size     int
	Pagelen  int
	Previous string
	Next     string
	Values   []*ProjectBBCloud
}

func (pl *projectBBCloudResp) GetValues() []BBProject {
	b := make([]BBProject, len(pl.Values))
	for i := range pl.Values {
		b[i] = pl.Values[i]
	}
	return b
}

func (b *BitbucketClient) listProjectsBBCloud() ([]BBProject, error) {
	projects := make([]*ProjectBBCloud, 0)
	path := fmt.Sprintf("/2.0/workspaces/%s/projects", b.Workspace)
	query := fmt.Sprintf("pagelen=%d", b.limit)
	for path != "" {
		projectList := &projectBBCloudResp{}
		body, err := b.get(path, query)
		if err != nil {
			return nil, err
		}
		path = ""
		query = ""

		err = json.Unmarshal(body, projectList)
		if err != nil {
			return nil, err
		}
		for _, v := range projectList.Values {
			v.client = b
			projects = append(projects, v)
		}
		if projectList.Next != "" {
			baseUrl, err := url.Parse(projectList.Next)
			if err != nil {
				return nil, err
			}
			path = baseUrl.Path
			query = baseUrl.RawQuery
		}
	}
	tmp := &projectBBCloudResp{Values: projects}
	return tmp.GetValues(), nil
}

type ProjectBBCloud struct {
	Key         string
	Id          string `json:"uuid"`
	Name        string
	Description string
	Private     bool   `json:"is_private"`
	ProjectType string `json:"type"`
	Links       map[string]map[string]string
	client      *BitbucketClient
}

func (p *ProjectBBCloud) GetKey() string {
	return p.Key
}
func (p *ProjectBBCloud) GetId() string {
	return p.Id
}
func (p *ProjectBBCloud) GetName() string {
	return p.Name
}
func (p *ProjectBBCloud) GetDescription() string {
	return p.Description
}
func (p *ProjectBBCloud) IsPublic() bool {
	return !p.Private
}
func (p *ProjectBBCloud) GetProjectType() string {
	return p.ProjectType
}

// MakePath creates the full relative path to the repo
func (p *ProjectBBCloud) ListReposURL() *url.URL {
	if value, ok := p.Links["repositories"]; ok {
		if href, ok := value["href"]; ok {
			baseUrl, err := url.Parse(href)
			if err != nil {
				return nil
			}
			return baseUrl
		}
	}
	return nil
}

func (p *ProjectBBCloud) GetBBClient() *BitbucketClient {
	return p.client
}

func (p *ProjectBBCloud) ListRepos() ([]BBRepo, error) {
	repos := make([]*RepoBBCloud, 0)
	lrurl := p.ListReposURL()
	query := lrurl.Query()
	query.Set("pagelen", fmt.Sprintf("%d", p.client.limit))
	lrurl.RawQuery = query.Encode()

	for lrurl != nil {
		lrepoList := &repoBBCloudResp{}
		body, err := p.client.getURL(lrurl)
		if err != nil {
			return nil, err
		}
		lrurl = nil
		query = nil

		err = json.Unmarshal(body, lrepoList)
		if err != nil {
			return nil, err
		}
		for _, v := range lrepoList.Values {
			v.Project = p
			repos = append(repos, v)
		}
		if lrepoList.Next != "" {
			baseUrl, err := url.Parse(lrepoList.Next)
			if err != nil {
				return nil, err
			}
			lrurl = baseUrl
		}
	}
	tmp := &repoBBCloudResp{Values: repos}
	return tmp.GetValues(), nil
}

type repoBBCloudResp struct {
	Page     int
	Size     int
	Pagelen  int
	Previous string
	Next     string
	Values   []*RepoBBCloud
}

func (rl *repoBBCloudResp) GetValues() []BBRepo {
	b := make([]BBRepo, len(rl.Values))
	for i := range rl.Values {
		b[i] = rl.Values[i]
	}
	return b
}

type RepoBBCloud struct {
	Slug          string
	Uuid          string
	Name          string
	ScmId         string
	State         string
	StatusMessage string
	Forkable      bool
	Project       *ProjectBBCloud `json:"-"`
	Public        bool
	Links         *RepoBBCloudLinks
}

type RepoBBCloudLinks struct {
	Source map[string]string
}

func (r *RepoBBCloud) GetSlug() string {
	return r.Slug
}
func (r *RepoBBCloud) GetId() string {
	return r.Uuid
}
func (r *RepoBBCloud) GetName() string {
	return r.Name
}
func (r *RepoBBCloud) GetScmId() string {
	return r.ScmId
}
func (r *RepoBBCloud) GetState() string {
	return r.State
}
func (r *RepoBBCloud) GetStatusMessage() string {
	return r.StatusMessage
}
func (r *RepoBBCloud) GetForkable() bool {
	return r.Forkable
}
func (r *RepoBBCloud) GetBBProject() BBProject {
	return r.Project
}
func (r *RepoBBCloud) IsPublic() bool {
	return r.Public
}
func (r *RepoBBCloud) GetFilesURL() *url.URL {
	if href, ok := r.Links.Source["href"]; ok {
		baseUrl, err := url.Parse(href)
		if err != nil {
			return nil
		}
		return baseUrl
	}
	return nil
}
func (r *RepoBBCloud) GetDir(path string) ([]BBFile, error) {
	files := make([]BBFile, 0)
	furl := r.GetFilesURL()
	if furl == nil {
		return nil, errors.New("Error creating URL to get files for Repo:" + r.GetName())
	}
	query := furl.Query()
	query.Set("pagelen", fmt.Sprintf("%d", r.GetBBProject().GetBBClient().limit))
	furl.RawQuery = query.Encode()
	furl.Path = filepath.Join(furl.Path, "/master/", path) + "/"

	for furl != nil {
		lfileList := &BBCloudFileList{}
		body, err := r.Project.client.getURL(furl)
		if err != nil {
			return nil, err
		}
		furl = nil
		query = nil

		err = json.Unmarshal(body, lfileList)
		if err != nil {
			return nil, err
		}
		for _, v := range lfileList.Values {
			v.repo = r
			v.basePath = path
			files = append(files, v)
		}
		if lfileList.Next != "" {
			baseUrl, err := url.Parse(lfileList.Next)
			if err != nil {
				return nil, err
			}
			furl = baseUrl
		}
	}
	return files, nil
}

func (r *RepoBBCloud) GetFile(path string) (BBFile, error) {
	fl, err := r.GetDir(filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	// fileName := filepath.Base(path)
	for _, f := range fl {
		if strings.HasSuffix(path, f.GetName()) {
			return f, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("File Not Found:%s", path))
}

type BBCloudFileList struct {
	Page     int
	Size     int
	Pagelen  int
	Previous string
	Next     string
	Values   []*BBCloudFile
}

type BBCloudFile struct {
	Name     string `json:"path"`
	Size     int
	Type     string
	Links    map[string]map[string]string
	basePath string
	repo     *RepoBBCloud
}

func (bbf *BBCloudFile) GetBasePath() string {
	return bbf.basePath
}
func (bbf *BBCloudFile) GetFullPath() string {
	return filepath.Join("/", bbf.Name)
}
func (bbf *BBCloudFile) GetName() string {
	return filepath.Base(bbf.Name)
}
func (bbf *BBCloudFile) GetFileType() string {
	if bbf.Type == "commit_file" {
		return "FILE"
	} else if bbf.Type == "commit_directory" {
		return "DICRECTORY"
	}
	return ""
}
func (bbf *BBCloudFile) GetSize() int {
	return bbf.Size
}
func (bbf *BBCloudFile) GetBBRepo() BBRepo {
	return bbf.repo
}
func (bbf *BBCloudFile) GetURL() (*url.URL, error) {
	var fdurl *url.URL
	var err error
	if href, ok := bbf.Links["self"]["href"]; ok {
		fdurl, err = url.Parse(href)
		if err != nil {
			return nil, err
		}
	}
	return fdurl, nil
}
func (bbf *BBCloudFile) GetData() ([]byte, error) {
	fdurl, err := bbf.GetURL()
	if err != nil {
		return nil, err
	}

	ba, err := bbf.GetBBRepo().GetBBProject().GetBBClient().getURL(fdurl)
	if err != nil {
		return nil, err
	}
	return ba, nil
}
