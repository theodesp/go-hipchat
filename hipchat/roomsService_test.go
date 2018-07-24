package hipchat

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func (suite *HipChatClientTestSuite) TestRoomsService_ListRooms() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf("/%s/%s", apiVersion2, listRoomsRoute)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodGet)
		fmt.Fprint(w, `{"items":[{"id":1,"is_archived": true},{"id":2,"is_archived": false}]}`)
	})

	rooms, _, err := suite.client.Rooms.ListRooms(context.Background(), nil)
	assert.Nil(err)

	want := []*RoomListItem{
		{ID: int64(1), IsArchived: true, Name: "", Privacy: "", Version: ""},
		{ID: int64(2), IsArchived: false, Name: "", Privacy: "", Version: ""}}
	assert.Equal(want, rooms)
}

func (suite *HipChatClientTestSuite) TestRoomsService_GetRoom() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(getRoomRoute, "1")
	route = fmt.Sprintf("/%s/%s", apiVersion2, route)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodGet)
		fmt.Fprint(w, `{"id":1,"is_archived": true, "name": "hello"}`)
	})

	room, _, err := suite.client.Rooms.Get(context.Background(), "1")
	assert.Nil(err)
	item := &RoomListItem{ID: int64(1), IsArchived: true, Name: "hello", Privacy: "", Version: ""}
	want := &Room{
		item,
		"","","", false, "",
		false, "", "", nil, nil,
	}
	assert.Equal(want, room)
}

func (suite *HipChatClientTestSuite) TestRoomsService_GetRoomEmptyId() {
	assert := assert.New(suite.T())
	_, _, err := suite.client.Rooms.Get(context.Background(), "")
	assert.EqualError(err, emptyParam.Error())
}
