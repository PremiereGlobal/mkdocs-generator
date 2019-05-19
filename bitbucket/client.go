package bitbucket

import (
  "net/http"
  "net/url"
  log "github.com/sirupsen/logrus"
)

type BitbucketClient struct {
  rawClient *http.Client
  BaseUrl *url.URL
  BaseApiPath string
  Username string
  Password string
  log *log.Logger
  limit int
}

func NewBitbucketClient(baseUrl string, username string, password string) (*BitbucketClient, error) {

  u, err := url.Parse(baseUrl)
  if err != nil {
    return nil, err
  }

  client := &BitbucketClient{
    rawClient: &http.Client{},
    BaseUrl: u,
    BaseApiPath: "/rest/api/1.0",
    Username: username,
    Password: password,
    limit: 1000,
  }

  return client, nil
}

func (b *BitbucketClient) SetLogger(log *log.Logger) {
  b.log = log
}
