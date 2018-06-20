package hipchat

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseUrl = "https://api.hipchat.com"
	v2Part         = "v2"
	userAgent      = "go-hipchat"
)

// A Client communicates with the HipChat API.
type Client struct {
	BaseUrl   *url.URL
	UserAgent string
	common    service
	V2        *V2
}

type service struct {
	client *http.Client
}

// V2 is the baseline namespace for the  HipChat API v2 operations
type V2 struct {
	versionPart string
}

// NewClient returns a new HipChat API client
func NewCLient(client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	baseUrl, _ := url.Parse(defaultBaseUrl)

	c := &Client{BaseUrl: baseUrl, UserAgent: userAgent}
	c.common.client = client
	c.V2 = &V2{versionPart: v2Part}

	return c
}
