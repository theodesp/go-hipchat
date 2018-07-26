package hipchat

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type HipChatUtilsTestSuite struct {
	suite.Suite
}

func TestHipChatUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(HipChatUtilsTestSuite))
}

func (suite *HipChatUtilsTestSuite) TestAddUrlOptions() {
	assert := assert.New(suite.T())
	type opts struct {
		HasName bool `url:"has_name"`
	}
	testCases := []struct {
		name      string
		s         string
		opt       interface{}
		want      string
		wantError bool
	}{
		{"TestNilOpts", "123", nil, "123", false},
		{"TestInvalidUrl", "http://[::1]a", nil, "http://[::1]a", true},
		{"TestEmptyOpts", "http://[::1]:80", struct{}{}, "http://[::1]:80", false},
		{"TestOptsWithNoTags", "http://[::1]:80", struct{ name string }{"hello"}, "http://[::1]:80", false},
		{"TestOptsWithTags", "http://[::1]:80", opts{true}, "http://[::1]:80?has_name=true", false},
		{"TestInvalidOpts", "http://[::1]:80", "::", "http://[::1]:80", true},
	}
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			v, err := addUrlOptions(tc.s, tc.opt)
			assert.Equal(tc.want, v)
			if tc.wantError {
				assert.NotNil(err)
				return
			}
			assert.Nil(err)
		})

	}
}
