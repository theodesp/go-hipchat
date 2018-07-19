package hipchat

import (
	"net/http"
	"net/url"
	"strings"
	"errors"
)

// A Client communicates with the HipChat API.
type Client struct {
	BaseUrl    *url.URL
	UserAgent  string
	common     service
	apiVersion string

	Room *RoomService
}

type service struct {
	client *http.Client
}

// NewClient returns a new HipChat API client
func NewClient(client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	baseUrl, _ := url.Parse(string(defaultBaseUrl + apiVersion2))

	c := &Client{BaseUrl: baseUrl, UserAgent: userAgent}
	c.apiVersion = apiVersion2
	c.common.client = client

	// Services
	c.Room = (*RoomService)(&c.common)

	return c
}

// Sets the HipChat API version. This defaults to v2
func (c *Client) SetApiVersion(apiVersion string) error {
	if strings.HasPrefix(apiVersion, "/") {
		return errors.New("set_api_version: apiVersion string parameter is prefixed with a forward slash (/)")
	}

	baseUrl, _ := url.Parse(string(defaultBaseUrl + apiVersion))
	c.apiVersion = apiVersion
	c.BaseUrl = baseUrl

	return nil
}
