package hipchat

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"encoding/json"
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
		fmt.Fprint(w, `{"id":1,"privacy":"public","name": "hello"}`)
	})

	room, _, err := suite.client.Rooms.GetRoom(context.Background(), "1")
	assert.Nil(err)

	want := NewRoom("hello")
	want.ID = int64(1)

	assert.Equal(want, room)
}

func (suite *HipChatClientTestSuite) TestRoomsService_UpdateRoom() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(getRoomRoute, "1")
	route = fmt.Sprintf("/%s/%s", apiVersion2, route)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodPut)
		w.WriteHeader(http.StatusNoContent)
	})

	room := NewRoom("")

	resp, err := suite.client.Rooms.UpdateRoom(context.Background(), "1", room)
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, resp.StatusCode)
}

func (suite *HipChatClientTestSuite) TestRoomsService_DeleteRoom() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(getRoomRoute, "1")
	route = fmt.Sprintf("/%s/%s", apiVersion2, route)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := suite.client.Rooms.DeleteRoom(context.Background(), "1")
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, resp.StatusCode)
}

func (suite *HipChatClientTestSuite) TestRoomsService_CreateRoom() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf("/%s/%s", apiVersion2, listRoomsRoute)

	room := NewRoom("hello")

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodPost)

		v := NewRoom("hello")
		json.NewDecoder(r.Body).Decode(v)
		assert.Equal(v, room)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":1}`)
	})

	v, resp, err := suite.client.Rooms.CreateRoom(context.Background(), room)
	want := NewRoom("hello")
	want.ID = int64(1)

	assert.Nil(err)
	assert.Equal(want, v)
	assert.Equal(http.StatusCreated, resp.StatusCode)
}

func (suite *HipChatClientTestSuite) TestRoomsService_SetRoomTopic() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(setRoomTopicRoute, "1")
	route = fmt.Sprintf("/%s/%s", apiVersion2, route)

	input := topicBody{"hello"}

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodPut)

		topic := topicBody{"hello"}
		json.NewDecoder(r.Body).Decode(topic)
		assert.Equal(topic, input)

		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := suite.client.Rooms.SetRoomTopic(context.Background(), "1", "hello")
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, resp.StatusCode)
}

func (suite *HipChatClientTestSuite) TestRoomsService_PassEmptyRoomId() {
	assert := assert.New(suite.T())
	_, _, err := suite.client.Rooms.GetRoom(context.Background(), "")
	assert.EqualError(err, emptyParam.Error())

	_, err = suite.client.Rooms.UpdateRoom(context.Background(), "", nil)
	assert.EqualError(err, emptyParam.Error())

	_, err = suite.client.Rooms.DeleteRoom(context.Background(), "")
	assert.EqualError(err, emptyParam.Error())

	_, err = suite.client.Rooms.SetRoomTopic(context.Background(), "", "")
	assert.EqualError(err, emptyParam.Error())
}
