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
		{Id: int64(1), IsArchived: true, Name: "", Privacy: "", Version: ""},
		{Id: int64(2), IsArchived: false, Name: "", Privacy: "", Version: ""}}
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
	want.Id = int64(1)

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
	want.Id = int64(1)

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

	resp, err := suite.client.Rooms.SetRoomTopic(context.Background(), "1", input.Topic)
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, resp.StatusCode)
}

func (suite *HipChatClientTestSuite) TestRoomsService_GetRoomStatistics() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(getRoomStatisticsRoute, "1")
	route = fmt.Sprintf("/%s/%s", apiVersion2, route)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodGet)
		fmt.Fprint(w, `{"messages_sent":100}`)
	})

	st, _, err := suite.client.Rooms.GetRoomStatistics(context.Background(), "1")
	assert.Nil(err)

	want := &RoomStatistic{100, ""}
	assert.Equal(want, st)
}

func (suite *HipChatClientTestSuite) TestRoomsService_ShareLinkWithRoom() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(shareLinkWithRoomRoute, "1")
	route = fmt.Sprintf("/%s/%s", apiVersion2, route)

	input := shareLinkBody{"hello", "link"}

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodPost)

		link := shareLinkBody{"hello", "link"}
		json.NewDecoder(r.Body).Decode(link)
		assert.Equal(link, input)

		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := suite.client.Rooms.ShareLinkWithRoom(context.Background(), "1", input.Message, input.Link)
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, resp.StatusCode)
}

func (suite *HipChatClientTestSuite) TestRoomsService_GetRoomParticipants() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(getRoomParticipantsRoute, "1")
	route = fmt.Sprintf("/%s/%s", apiVersion2, route)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodGet)
		fmt.Fprint(w, `{"items":[{"id":1,"name":"Theo"},{"id":2,"name":"Alex"}]}`)
	})

	participants, _, err := suite.client.Rooms.GetRoomParticipants(context.Background(), "1", nil)
	assert.Nil(err)

	want := []*UserListItem{
		{Id: int64(1), Name: "Theo", Version: ""},
		{Id: int64(2), Name: "Alex", Version: ""}}
	assert.Equal(want, participants)
}

func (suite *HipChatClientTestSuite) TestRoomsService_ReplyToRoomMessage() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(replyToRoomMessageRoute, "1")
	route = fmt.Sprintf("/%s/%s", apiVersion2, route)

	input := replyToMessageBody{"1", "hello"}

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodPost)

		reply := replyToMessageBody{"1", "hello"}
		json.NewDecoder(r.Body).Decode(reply)
		assert.Equal(reply, input)

		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := suite.client.Rooms.ReplyToRoomMessage(context.Background(), "1", input.MessageId, input.Message)
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, resp.StatusCode)
}

func (suite *HipChatClientTestSuite) TestRoomsService_SendRoomMessage() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(sendRoomMessageRoute, "1")
	route = fmt.Sprintf("/%s/%s", apiVersion2, route)

	input := sendMessageBody{"hello"}

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodPost)

		reply := sendMessageBody{"hello"}
		json.NewDecoder(r.Body).Decode(reply)
		assert.Equal(reply, input)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id": "123"}`)
	})

	m, resp, err := suite.client.Rooms.SendRoomMessage(context.Background(), "1", input.Message)
	assert.Nil(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	assert.Equal("123", m.Id)
}

func (suite *HipChatClientTestSuite) TestRoomsService_GetRoomMembers() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(getRoomMembersRoute, "1")
	route = fmt.Sprintf("/%s/%s", apiVersion2, route)

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodGet)
		fmt.Fprint(w, `{"items":[{"id":1,"name":"Theo"},{"id":2,"name":"Alex"}]}`)
	})

	members, _, err := suite.client.Rooms.GetRoomMembers(context.Background(), "1", nil)
	assert.Nil(err)

	want := []*UserListItem{
		{Id: int64(1), Name: "Theo", Version: ""},
		{Id: int64(2), Name: "Alex", Version: ""}}
	assert.Equal(want, members)
}

func (suite *HipChatClientTestSuite) TestRoomsService_AddRoomMember() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(roomMemberRoute, "1")
	route = fmt.Sprintf("/%s/%s%s", apiVersion2, route, "theo")

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodPut)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := suite.client.Rooms.AddRoomMember(context.Background(), "1", "theo")
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, resp.StatusCode)
}

func (suite *HipChatClientTestSuite) TestRoomsService_RemoveRoomMember() {
	assert := assert.New(suite.T())
	route := fmt.Sprintf(roomMemberRoute, "1")
	route = fmt.Sprintf("/%s/%s%s", apiVersion2, route, "theo")

	suite.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := suite.client.Rooms.RemoveRoomMember(context.Background(), "1", "theo")
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, resp.StatusCode)
}

func (suite *HipChatClientTestSuite) TestRoomsService_EmptyRoomParams() {
	assert := assert.New(suite.T())
	_, _, err := suite.client.Rooms.GetRoom(context.Background(), "")
	assert.EqualError(err, emptyParam.Error())

	_, err = suite.client.Rooms.UpdateRoom(context.Background(), "", nil)
	assert.EqualError(err, emptyParam.Error())

	_, err = suite.client.Rooms.DeleteRoom(context.Background(), "")
	assert.EqualError(err, emptyParam.Error())

	_, err = suite.client.Rooms.SetRoomTopic(context.Background(), "", "")
	assert.EqualError(err, emptyParam.Error())

	_, _, err = suite.client.Rooms.GetRoomStatistics(context.Background(), "")
	assert.EqualError(err, emptyParam.Error())

	_, _, err = suite.client.Rooms.GetRoomParticipants(context.Background(), "", nil)
	assert.EqualError(err, emptyParam.Error())

	_, err = suite.client.Rooms.ReplyToRoomMessage(context.Background(),"", "", "")
	assert.EqualError(err, emptyParam.Error())

	_, _, err = suite.client.Rooms.SendRoomMessage(context.Background(),"", "")
	assert.EqualError(err, emptyParam.Error())

	_, _, err = suite.client.Rooms.GetRoomMembers(context.Background(),"", nil)
	assert.EqualError(err, emptyParam.Error())

	_, err = suite.client.Rooms.AddRoomMember(context.Background(),"", "")
	assert.EqualError(err, emptyParam.Error())

	_, err = suite.client.Rooms.RemoveRoomMember(context.Background(),"", "")
	assert.EqualError(err, emptyParam.Error())
}
