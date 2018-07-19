package hipchat

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/url"
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
	suite.sut = NewClient(nil)
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
	expectedBaseUrl, _ := url.Parse("https://api.hipchat.com/v3")

	err := suite.sut.SetApiVersion("v3")

	assert.Nil(err)
	assert.Equal(suite.sut.BaseUrl, expectedBaseUrl)

	err = suite.sut.SetApiVersion("/v3")
	assert.NotNil(err)
	assert.Equal(suite.sut.BaseUrl, expectedBaseUrl)

}
