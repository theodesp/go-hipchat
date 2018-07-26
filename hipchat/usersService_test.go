package hipchat

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stretchr/testify/assert"
)

func (suite *HipChatClientTestSuite) TestUserService_ListRooms() {
	assertion := assert.New(suite.T())
	route := fmt.Sprintf("/%s/%s", apiVersion2, listUsersRoute)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(r.Method, http.MethodGet)
		fmt.Fprint(w, `{"items":[{"id":1,"mention_name":"john","name":"john doe"},{"id":2,"mention_name":"jane","name":"jane doe"}]}`)
	})

	users, _, err := suite.client.Users.ListUsers(context.Background(), nil)
	assertion.Nil(err)

	want := []*UserListItem{
		{Id: int64(1), MentionName: "john", Name: "john doe", Version: ""},
		{Id: int64(2), MentionName: "jane", Name: "jane doe", Version: ""}}
	assertion.Equal(want, users)
}
