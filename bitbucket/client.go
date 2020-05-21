package bitbucket

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

type BitbucketClientConfig struct {

	// Url is the full url to the Bitbucket instance, for example,
	// https://bitbucket.mysite.com
	Url string

	// Username is the username for authenticated instances of Bitbucket
	Username string

	// Password is the password for authenticated instances of Bitbucket
	Password string
}

// BitbucketClient is the main
type BitbucketClient struct {
	rawClient   *http.Client
	config      *BitbucketClientConfig
	BaseUrl     *url.URL
	BaseApiPath string
	log         *log.Logger
	limit       int
}

func NewBitbucketClient(config *BitbucketClientConfig) (*BitbucketClient, error) {

	// We use Hashicorp's retry client to retry if we get errors
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	standardClient := retryClient.StandardClient()

	baseUrl, err := url.Parse(config.Url)
	if err != nil {
		return nil, err
	}

	client := &BitbucketClient{
		rawClient:   standardClient,
		config:      config,
		BaseUrl:     baseUrl,
		BaseApiPath: "/rest/api/1.0",
		limit:       1000,
	}

	return client, nil
}

func (b *BitbucketClient) SetLogger(log *log.Logger) {
	b.log = log
}
