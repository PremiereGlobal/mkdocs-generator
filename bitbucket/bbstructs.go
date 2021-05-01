package bitbucket

import "net/url"

type BBProject interface {
	GetKey() string
	GetId() string
	GetName() string
	GetDescription() string
	IsPublic() bool
	GetProjectType() string
	ListReposURL() *url.URL
	ListRepos() ([]BBRepo, error)
	GetBBClient() *BitbucketClient
}

type BBRepo interface {
	GetSlug() string
	GetId() string
	GetName() string
	GetScmId() string
	GetState() string
	GetStatusMessage() string
	GetForkable() bool
	GetBBProject() BBProject
	IsPublic() bool
	GetFilesURL() *url.URL
	GetDir(string) ([]BBFile, error)
	GetFile(string) (BBFile, error)
}

type BBFile interface {
	GetBasePath() string
	GetFullPath() string
	GetName() string
	GetFileType() string
	GetSize() int
	GetBBRepo() BBRepo
	GetURL() (*url.URL, error)
	GetData() ([]byte, error)
}
