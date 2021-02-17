package redmine

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	config Config
}

type Config struct {
	URL            string
	Username       string
	Password       string
	SkipCertVerify bool
}

func NewClient(config Config) *Client {
	return &Client{config}
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.config.SkipCertVerify},
	}
	httpClient := &http.Client{Transport: tr}

	req.SetBasicAuth(c.config.Username, c.config.Password)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to do request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("response had statuscode: %d and body: %s", resp.StatusCode, body)
	}

	return body, err
}
