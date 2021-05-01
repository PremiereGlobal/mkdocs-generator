package bitbucket

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
)

type BitbucketClientConfig struct {

	// Url is the full url to the Bitbucket instance, for example,
	// https://bitbucket.mysite.com
	Url string

	// Username is the username for authenticated instances of Bitbucket
	Username string
	Logger   *log.Logger
	// Password is the password for authenticated instances of Bitbucket
	Password string

	Workspace string
}

// BitbucketClient is the main
type BitbucketClient struct {
	rawClient   *http.Client
	config      *BitbucketClientConfig
	BaseUrl     *url.URL
	BaseApiPath string
	log         *log.Logger
	limit       int
	IsBBCloud   bool
	Workspace   string
}

func NewBitbucketClient(config *BitbucketClientConfig) (*BitbucketClient, error) {

	// We use Hashicorp's retry client to retry if we get errors
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.Logger = &WrappedLeveledLogger{logger: config.Logger}
	retryClient.RetryWaitMin = time.Duration(5 * time.Second)
	retryClient.RetryWaitMax = time.Duration(10 * time.Second)
	retryClient.CheckRetry = localRetryPolicy

	standardClient := retryClient.StandardClient()
	bbcloud := false
	baseapi := "/rest/api/1.0"
	limit := 1000
	if strings.Contains(config.Url, "bitbucket.org") {
		bbcloud = true
		baseapi = "/"
		limit = 100
	}
	baseUrl, err := url.Parse(config.Url)
	if err != nil {
		return nil, err
	}

	client := &BitbucketClient{
		rawClient:   standardClient,
		config:      config,
		BaseUrl:     baseUrl,
		IsBBCloud:   bbcloud,
		BaseApiPath: baseapi,
		limit:       limit,
		Workspace:   config.Workspace,
		log:         config.Logger,
	}

	return client, nil
}

var redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)
var schemeErrorRe = regexp.MustCompile(`unsupported protocol scheme`)

//Copied this code to add 429 retries since thats what BB.org gives
func localRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	if err != nil {
		if v, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to too many redirects.
			if redirectsErrorRe.MatchString(v.Error()) {
				return false, nil
			}

			// Don't retry if the error was due to an invalid protocol scheme.
			if schemeErrorRe.MatchString(v.Error()) {
				return false, nil
			}

			// Don't retry if the error was due to TLS cert verification failure.
			if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
				return false, nil
			}
		}

		// The error is likely recoverable so retry.
		return true, nil
	}
	if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != 501) || resp.StatusCode == 429 {
		return true, nil
	}

	return false, nil
}

type WrappedLeveledLogger struct {
	logger *log.Logger
}

func (wll *WrappedLeveledLogger) Error(s string, iface ...interface{}) {
	wll.logger.Error(toString(s, iface...))
}
func (wll *WrappedLeveledLogger) Info(s string, iface ...interface{}) {
	wll.logger.Info(toString(s, iface...))
}
func (wll *WrappedLeveledLogger) Debug(s string, iface ...interface{}) {
	wll.logger.Debug(toString(s, iface...))
}
func (wll *WrappedLeveledLogger) Warn(s string, iface ...interface{}) {
	wll.logger.Warn(toString(s, iface...))
}

func toString(s string, iface ...interface{}) string {
	sb := strings.Builder{}
	sb.WriteString(s)
	sb.WriteString(" ")
	for i := 0; i < len(iface); i++ {
		arg := iface[i]
		isString := arg != nil && reflect.TypeOf(arg).Kind() == reflect.String
		if isString {
			sb.WriteString(fmt.Sprintf("%s ", arg))
		} else {
			sb.WriteString(fmt.Sprintf("%v ", arg))
		}
	}
	return sb.String()
}
