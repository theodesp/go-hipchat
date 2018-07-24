package hipchat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type HipChatClientTestSuite struct {
	suite.Suite

	// client is the GitHub client being tested.
	client *Client

	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
}

func TestHipChatClientTestSuite(t *testing.T) {
	suite.Run(t, new(HipChatClientTestSuite))
}

func (suite *HipChatClientTestSuite) SetupTest() {
	suite.client = NewClient(nil)
	suite.mux = http.NewServeMux()
	suite.server = httptest.NewServer(suite.mux)

	// hipchat client configured to use test server
	url, _ := url.Parse(suite.server.URL)
	suite.client.BaseUrl = url
}

func (suite *HipChatClientTestSuite) TearDownTest() {
	suite.server.Close()
}

func (suite *HipChatClientTestSuite) TestNewClient() {
	assert := assert.New(suite.T())

	assert.NotNil(suite.client.Rooms)
}

func (suite *HipChatClientTestSuite) TestClient_SetApiVersion() {
	assert := assert.New(suite.T())
	expectedBaseUrl, _ := url.Parse("https://api.hipchat.com/v3")

	err := suite.client.SetApiVersion("v3")

	assert.Nil(err)
	assert.Equal(suite.client.BaseUrl, expectedBaseUrl)

	err = suite.client.SetApiVersion("/v3")
	assert.NotNil(err)
	assert.Equal(suite.client.BaseUrl, expectedBaseUrl)

}

func (suite *HipChatClientTestSuite) TestClient_TestNewRequest() {
	assert := assert.New(suite.T())

	inURL, outURL := "bar", suite.server.URL+"/v2/bar"
	inBody, outBody := &RoomListItem{Id: 1}, `{"id":1,"is_archived":false,"name":"","privacy":"","version":""}`+"\n"

	req, err := suite.client.NewRequest("GET", inURL, inBody)
	assert.Nil(err)

	assert.Equal(req.URL.String(), outURL)

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	assert.Equal(outBody, string(body))

	// test that default user-agent is attached to the request
	assert.Equal(req.Header.Get("User-Agent"), suite.client.UserAgent)
}

func (suite *HipChatClientTestSuite) TestNewRequest_invalidJSON() {
	assert := assert.New(suite.T())
	type T struct {
		A map[interface{}]interface{}
	}
	_, err := suite.client.NewRequest("GET", "/", &T{})
	assert.NotNil(err)
	assert.IsType(&json.UnsupportedTypeError{}, err)
}

func (suite *HipChatClientTestSuite) TestResponse_paginationValues() {
	assert := assert.New(suite.T())

	r := http.Response{
		Body: ioutil.NopCloser(
			bytes.NewReader([]byte(`{"maxResults":10,"startIndex":100,"links":{"next":"123","prev":"123","self":"123"}}`))),
	}

	response, _ := newPaginatedResponse(&r)
	assert.Equal("123", response.Links.Next)
	assert.Equal("123", response.Links.Prev)
	assert.Equal("123", response.Links.Self)
	assert.Equal(10, response.MaxResults)
	assert.Equal(100, response.StartIndex)
}

func (suite *HipChatClientTestSuite) TestResponse_paginationValuesInvalid() {
	assert := assert.New(suite.T())

	r := http.Response{
		Body: ioutil.NopCloser(
			bytes.NewReader([]byte(`{"maxResults":"10","startIndex":"abc","links":{"next":1,"previus":"123"}}`))),
	}

	response, _ := newPaginatedResponse(&r)
	assert.Equal("", response.Links.Next)
	assert.Equal("", response.Links.Prev)
	assert.Equal("", response.Links.Self)
	assert.Equal(0, response.MaxResults)
	assert.Equal(0, response.StartIndex)
}

func (suite *HipChatClientTestSuite) TestDo() {
	assert := assert.New(suite.T())
	type foo struct {
		A string
	}

	suite.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := suite.client.Get(".")
	body := new(foo)
	suite.client.Do(context.Background(), req, body)

	want := &foo{"a"}
	assert.Equal(want, body)
}

func (suite *HipChatClientTestSuite) TestDo_httpError() {
	assert := assert.New(suite.T())

	suite.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := suite.client.Get(".")
	resp, _ := suite.client.Do(context.Background(), req, nil)

	assert.Equal(400, resp.StatusCode)
}
