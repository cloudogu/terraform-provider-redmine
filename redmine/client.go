package redmine

import (
	rmapi "github.com/cloudogu/go-redmine"
)

type Client struct {
	config     Config
	redmineAPI *rmapi.Client
}

type Config struct {
	URL            string
	Username       string
	Password       string
	SkipCertVerify bool
	APIKey         string
}

func NewClient(config Config) *Client {
	redmineAPI := rmapi.NewClient(config.URL, config.APIKey)
	redmineAPI.Limit = -1
	redmineAPI.Offset = -1

	return &Client{config: config, redmineAPI: redmineAPI}
}
