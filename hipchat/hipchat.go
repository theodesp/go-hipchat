package hipchat

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseUrl = "https://api.hipchat.com"
	userAgent      = "go-hipchat"
)

// A Client communicates with the HipChat API.
type Client struct {
	BaseUrl   *url.URL
	UserAgent string
	common    service
	apiVersion string

	Room *RoomService
}

type service struct {
	client *http.Client
}

// NewClient returns a new HipChat API client
func NewCLient(client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	baseUrl, _ := url.Parse(string(defaultBaseUrl + "/" + ApiVersions.V2))

	c := &Client{BaseUrl: baseUrl, UserAgent: userAgent}
	c.apiVersion = ApiVersions.V2
	c.common.client = client

	// Services
	c.Room = (*RoomService)(&c.common)

	return c
}

// Sets the HipChat API version. This defaults to v2
func (c *Client)SetApiVersion(apiVersion string)  {
	baseUrl, _ := url.Parse(string(defaultBaseUrl + "/" + apiVersion))
	c.BaseUrl = baseUrl
}
