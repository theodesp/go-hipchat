package hipchat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/philippfranke/multipart-related/related"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
)

// A Client communicates with the HipChat API.
type Client struct {
	client *http.Client

	BaseUrl    *url.URL
	UserAgent  string
	common     service
	apiVersion string

	Rooms *RoomsService
	Users *UsersService
}

type service struct {
	client *Client
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, StartIndex is the starting index of the request.
	StartIndex int `url:"start-index,omitempty"`

	// For paginated result sets, StartIndex is the maximum number of results to include per page.
	MaxResults int `url:"max-results,omitempty"`
}

// NewClient returns a new HipChat API client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseUrl, _ := url.Parse(string(defaultBaseUrl + apiVersion2))

	c := &Client{client: httpClient, BaseUrl: baseUrl, UserAgent: userAgent}
	c.apiVersion = apiVersion2
	c.common.client = c

	// Services
	c.Rooms = (*RoomsService)(&c.common)
	c.Users = (*UsersService)(&c.common)

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	ref, err := c.BaseUrl.Parse(c.apiVersion + "/" + urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseUrl.ResolveReference(ref)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", contentTypeApplicationJson)
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// NewUploadRequest creates an upload request.
func (c *Client) NewUploadRequest(
	urlStr string,
	reader io.Reader,
	size int64,
	mediaType string,
	body interface{},
	fileName string) (*http.Request, error) {
	ref, err := c.BaseUrl.Parse(c.apiVersion + "/" + urlStr)
	if err != nil {
		return nil, err
	}

	if reader == nil {
		return nil, errors.New("upload_request: missing file reader parameter")
	}

	var buf *bytes.Buffer
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	var wBuf = new(bytes.Buffer)
	w := related.NewWriter(wBuf)
	_, err =
		w.CreateRoot("", "", nil)
	if err != nil {
		return nil, err
	}

	if body != nil {
		var header = make(textproto.MIMEHeader)
		header.Set("Content-Type", contentTypeApplicationJson)
		header.Set("Content-Disposition", contentDispositionMetadata)
		nextPart, err := w.CreatePart("", header)
		if err != nil {
			return nil, err
		}
		nextPart.Write(buf.Bytes())
	}

	var header = make(textproto.MIMEHeader)
	header.Set("Content-Type", mediaType)
	header.Set("Content-Disposition", fmt.Sprintf(contentDispositionFile, fileName))
	nextPart, err := w.CreatePart("", header)
	if err != nil {
		return nil, err
	}

	upload, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	nextPart.Write(upload)
	if err := w.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ref.String(), bytes.NewReader(wBuf.Bytes()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// PaginatedResponse is a HipChat API response. This wraps the standard http.Response
// returned from HipChat and provides convenient access to things like
// pagination links.
type PaginatedResponse struct {
	*http.Response
	MaxResults int
	StartIndex int
	Links      *PaginationLinks
}

// PaginationLinks consists of a list of links for the current, next and previous pagination
// results based on the current parameters.
type PaginationLinks struct {
	Next string
	Prev string
	Self string
}

type paginatedListResponse struct {
	MaxResults int `json:"maxResults,omitempty"`
	StartIndex int `json:"startIndex,omitempty"`
	Links      *struct {
		Next string `json:"next,omitempty"`
		Prev string `json:"prev,omitempty"`
		Self string `json:"self"`
	} `json:"links,omitempty"`
}

// newPaginatedResponse creates a new Response for the provided http.Response.
func newPaginatedResponse(r *http.Response) (*PaginatedResponse, []byte) {
	response := &PaginatedResponse{Response: r}

	var v paginatedListResponse
	rs, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(rs, &v)

	response.MaxResults = v.MaxResults
	response.StartIndex = v.StartIndex

	response.Links = (*PaginationLinks)(v.Links)
	return response, rs
}

// Convenient shorthand for GET requests
func (c *Client) Get(urlStr string) (*http.Request, error) {
	return c.NewRequest(http.MethodGet, urlStr, nil)
}

// Convenient shorthand for POST requests
func (c *Client) Post(urlStr string, body interface{}) (*http.Request, error) {
	return c.NewRequest(http.MethodPost, urlStr, body)
}

// Convenient shorthand for PUT requests
func (c *Client) Put(urlStr string, body interface{}) (*http.Request, error) {
	return c.NewRequest(http.MethodPut, urlStr, body)
}

// Convenient shorthand for DELETE requests
func (c *Client) Delete(urlStr string) (*http.Request, error) {
	return c.NewRequest(http.MethodDelete, urlStr, nil)
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.

// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*PaginatedResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		err := resp.Body.Close()
		if err != nil {
			log.Println("closing response body failed")
		}
	}()

	if resp.StatusCode >= 300 {
		rs, _ := ioutil.ReadAll(resp.Body)
		return &PaginatedResponse{Response: resp}, errors.New(string(rs))
	}

	// Populate pagination params
	response, body := newPaginatedResponse(resp)
	bodyReader := bytes.NewReader(body)

	if v != nil {
		if len(body) == 0 && resp.StatusCode == 204 {
			return response, nil
		} else if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, bodyReader)
			if err != nil {
				return nil, err
			}
		} else {
			decErr := json.NewDecoder(bodyReader).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

// Sets the HipChat API version. This defaults to v2
func (c *Client) SetApiVersion(apiVersion string) error {
	if strings.HasPrefix(apiVersion, "/") {
		return invalidSetApiVersion
	}

	baseUrl, _ := url.Parse(string(defaultBaseUrl + apiVersion))
	c.apiVersion = apiVersion
	c.BaseUrl = baseUrl

	return nil
}
