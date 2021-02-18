package redmine

import (
	rmapi "github.com/mattn/go-redmine"
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
	redmindAPI := rmapi.NewClient(config.URL, config.APIKey)

	return &Client{config: config, redmineAPI: redmindAPI}
}
