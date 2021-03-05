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

func NewClient(config Config) *Client {
	redmineAPI := rmapi.NewClient(config.URL, config.APIKey)
	redmineAPI.Limit = -1
	redmineAPI.Offset = -1

	return &Client{config: config, redmineAPI: redmineAPI}
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
