package hipchat

import (
	"net/url"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type HipChatClientTestSuite struct {
	suite.Suite
	sut *Client
}

func TestHipChatClientTestSuite(t *testing.T) {
	suite.Run(t, new(HipChatClientTestSuite))
}

func (suite *HipChatClientTestSuite) SetupTest() {
	suite.sut = NewCLient(nil)
}

func (suite *HipChatClientTestSuite) TestNewClient() {
	assert := assert.New(suite.T())

	expectedBaseUrl, _ := url.Parse("https://api.hipchat.com/v2")
	expectedUserAgent := "go-hipchat"

	assert.Equal(suite.sut.BaseUrl, expectedBaseUrl)
	assert.Equal(suite.sut.UserAgent, expectedUserAgent)
	assert.NotNil(suite.sut.Room)
}

func (suite *HipChatClientTestSuite) TestClient_SetApiVersion() {
	assert := assert.New(suite.T())
	suite.sut.SetApiVersion(ApiVersion("v3"))

	expectedBaseUrl, _ := url.Parse("https://api.hipchat.com/v3")

	assert.Equal(suite.sut.BaseUrl, expectedBaseUrl)

}
