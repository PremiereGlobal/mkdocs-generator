package bitbucket

// ListProjects returns a list of all projects on the Bitbucket server
func (b *BitbucketClient) ListProjects() ([]BBProject, error) {
	if b.IsBBCloud {
		return b.listProjectsBBCloud()
	}
	return b.listProjectsBBServer()
}
