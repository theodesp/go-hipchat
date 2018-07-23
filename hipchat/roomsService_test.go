package hipchat

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func (suite *HipChatClientTestSuite) TestRoomsService_List() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf("/%s/%s", apiVersion2, roomsListPrefix)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodGet)
		fmt.Fprint(w, `{"items":[{"id":1,"is_archived": true},{"id":2,"is_archived": false}]}`)
	})

	rooms, _, err := suite.client.Rooms.List(context.Background(), nil)
	assert.Nil(err)

	want := []*RoomListItem{
		{ID: int64(1), IsArchived: true, Name: "", Privacy: "", Version: ""},
		{ID: int64(2), IsArchived: false, Name: "", Privacy: "", Version: ""}}
	assert.Equal(want, rooms)
}
