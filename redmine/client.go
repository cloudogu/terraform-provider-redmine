package redmine

import (
	"fmt"
	rmapi "github.com/cloudogu/go-redmine"
	"strconv"
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

func NewClient(config Config) (*Client, error) {
	redmineAPI, err := rmapi.NewClientBuilder().
		Endpoint(config.URL).
		AuthBasicAuth(config.Username, config.Password).
		SkipSSLVerify(config.SkipCertVerify).
		Build()
	if err != nil {
		return nil, err
	}

	return &Client{config: config, redmineAPI: redmineAPI}, nil
}

func verifyIDtoInt(id string) (int, error) {
	if id == "" || id == "0" {
		return 0, fmt.Errorf("invalid id '%s' found: must not be empty or 0", id)
	}

	idInt, _ := strconv.Atoi(id)
	if idInt < 0 {
		return 0, fmt.Errorf("invalid id '%d': must be strictly positive", idInt)
	}

	return idInt, nil
}
