package hipchat

import (
	"context"
	"fmt"
)

const (
	// Public Room access
	RoomPrivacyPublic = "public"

	// Private Room access
	RoomPrivacyPrivate = "private"


	listRoomsRoute = "room"
	getRoomRoute   = "room/%v"
	setRoomTopicRoute = "room/%v/topic"
	getRoomStatisticsRoute = "room/%v/statistics"
	shareLinkWithRoomRoute = "room/%v/share/link"
	getRoomParticipantsRoute = "room/%v/participant"
	replyToRoomMessageRoute = "room/%v/reply"
	sendRoomMessageRoute = "room/%v/message"
	getRoomMembersRoute = "room/%v/member"
)

// RoomsService handles communication with the room related
// methods of the HipChat API.
//
type RoomsService service

// RoomListItem represents a HipChat Room list item
type RoomListItem struct {
	// Id of the room.
	Id         int64  `json:"id"`

	// Whether or not this room is archived.
	IsArchived bool   `json:"is_archived"`

	// Name of the room.
	Name       string `json:"name"`

	// Privacy setting. Valid values: public, private.
	Privacy    string `json:"privacy"`

	// An etag-like random version string.
	Version    string `json:"version"`

	// URLs to retrieve room information
	Links *RoomLinks `json:"links,omitempty"`
}

type RoomLinks struct {
	// The URL to use to retrieve the full room information
	Self string `json:"self"`

	// The URL to use to retrieve webhooks registered for this room
	Webhooks string `json:"webhooks"`

	// The URL to use to retrieve members for this room. Only available for private rooms.
	Members string `json:"members,omitempty"`

	//The URL to use to retrieve participants for this room
	Participants string `json:"participants"`
}

// Room represents a HipChat Room
type Room struct {
	RoomListItem

	// XMPP/Jabber Id of the room.
	XmppJid string `json:"xmpp_jid"`

	// Time the room was created in ISO 8601 format UTC.
	Created string `json:"created"`

	// Privacy setting
	// Valid values: public, private.
	Privacy string `json:"privacy"`

	// Whether or not guests can access this room.
	IsGuestAccessible bool `json:"is_guest_accessible"`

	// URL to rooms's avatar. 125px on the longest side.
	AvatarUrl string `json:"avatar_url"`

	// Whether the room is visible to delegate admins, may be null to use the group default.
	// May be null.
	DelegateAdminVisibility bool `json:"delegate_admin_visibility"`

	// Current topic.
	Topic string `json:"topic"`

	// URL for guest access, if enabled.
	// May be null.
	GuestAccessUrl string `json:"guest_access_url"`

	Owner *UserListItem `json:"owner,omitempty"`

	// Statistics for this room.
	Statistics *struct{
		Links *struct{
			// The URL to use to retrieve room statistics
			Self string `json:"self"`
		} `json:"links,omitempty"`
	} `json:"statistics,omitempty"`
}

// RoomStatistic represents a HipChat Room Statistic
type RoomStatistic struct {
	// The number of messages sent in this room for its entire history.
	MessagesSent int64 `json:"messages_sent"`

	// Time of last activity (sent message) in the room in UNIX time (UTC).
	// May be null in rare cases when the time is unknown.
	LastActive string `json:"last_active"`
}

// RoomMessage represents a HipChat Room Message
type RoomMessage struct {
	// The unique identifier of the sent message.
	Id string `json:"id"`

	// The UTC timestamp representing when the message was processed.
	Timestamp string `json:"timestamp"`
}

// RoomsListOptions specifies the optional parameters to the
// RoomService.ListRooms
type RoomsListOptions struct {
	// Filter out private rooms
	IncludePrivate bool `url:"include-private,omitempty"`

	// Filter rooms
	IncludeArchived bool `url:"include-archived,omitempty"`
	ListOptions
}

// RoomParticipantsOptions specifies the optional parameters to the
// RoomService.GetRoomParticipants
type RoomParticipantsOptions struct {
	// Filter users by status (boolean). Only valid for private rooms.
	//
	// Defaults to 'false'.
	IncludeOffline bool `url:"include-offline,omitempty"`
	ListOptions
}

// List non-archived rooms for this group.
//
// Authentication required, with scope view_group or view_room.
// Accessible by group clients, users.
func (s *RoomsService) ListRooms(ctx context.Context, opt *RoomsListOptions) ([]*RoomListItem, *PaginatedResponse, error) {
	opts, err := addUrlOptions(listRoomsRoute, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.Get(opts)
	if err != nil {
		return nil, nil, err
	}

	var rooms *roomsListResponse
	resp, err := s.client.Do(ctx, req, &rooms)
	if err != nil {
		return nil, resp, err
	}

	return rooms.Items, resp, nil
}


// Get room details.
//
// Authentication required, with scope view_group or view_room.
// Accessible by group clients, room clients, users.
func (s *RoomsService) GetRoom(ctx context.Context, roomIdOrName string) (*Room, *PaginatedResponse, error) {
	var u, err = getRoomResourcePath(roomIdOrName, getRoomRoute)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.Get(u)
	if err != nil {
		return nil, nil, err
	}

	app := new(Room)
	resp, err := s.client.Do(ctx, req, app)
	if err != nil {
		return nil, resp, err
	}

	return app, resp, nil
}

// Updates a room.
//
// Authentication required, with scope admin_room.
// Accessible by group clients, users.
func (s *RoomsService) UpdateRoom(ctx context.Context, roomIdOrName string, room *Room) (*PaginatedResponse, error) {
	var u, err = getRoomResourcePath(roomIdOrName, getRoomRoute)
	if err != nil {
		return nil, err
	}

	req, err := s.client.Put(u, room)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, room)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Deletes a room and kicks the current participants.
//
// Authentication required, with scope manage_rooms.
// Accessible by group clients, users.
func (s *RoomsService) DeleteRoom(ctx context.Context, roomIdOrName string) (*PaginatedResponse, error) {
	var u, err = getRoomResourcePath(roomIdOrName, getRoomRoute)
	if err != nil {
		return nil, err
	}

	req, err := s.client.Delete(u)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Creates a new room.
//
// Authentication required, with scope manage_rooms.
// Accessible by group clients, users.
func (s *RoomsService) CreateRoom(ctx context.Context, room *Room) (*Room, *PaginatedResponse, error) {
	req, err := s.client.Post(listRoomsRoute, room)
	if err != nil {
		return nil, nil, err
	}

	r := new(Room)
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	room.Id = r.Id
	room.Links = r.Links

	return room, resp, nil
}

// Set a room's topic. Useful for displaying statistics, important links, server status, you name it!
//
// Authentication required, with scope admin_room.
// Accessible by group clients, room clients, users.
func (s *RoomsService) SetRoomTopic(ctx context.Context, roomIdOrName string, topic string) (*PaginatedResponse, error) {
	var u, err = getRoomResourcePath(roomIdOrName, setRoomTopicRoute)
	if err != nil {
		return nil, err
	}

	req, err := s.client.Put(u, topicBody{topic})
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Fetch statistics for this room.
//
// Authentication required, with scope view_group or view_room.
// Accessible by group clients, room clients, users.
func (s *RoomsService) GetRoomStatistics(ctx context.Context, roomIdOrName string) (*RoomStatistic, *PaginatedResponse, error) {
	var u, err = getRoomResourcePath(roomIdOrName, getRoomStatisticsRoute)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.Get(u)
	if err != nil {
		return nil, nil, err
	}

	st := new(RoomStatistic)
	resp, err := s.client.Do(ctx, req, st)
	if err != nil {
		return nil, resp, err
	}

	return st, resp, nil
}

// Share a link with the room.
//
// Authentication required, with scope send_message.
// Accessible by users.
func (s *RoomsService) ShareLinkWithRoom(ctx context.Context, roomIdOrName string, message string, link string) (*PaginatedResponse, error) {
	var u, err = getRoomResourcePath(roomIdOrName, shareLinkWithRoomRoute)
	if err != nil {
		return nil, err
	}

	req, err := s.client.Post(u, shareLinkBody{message, link})
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Gets all participants in this room.
//
// Authentication required, with scope view_room.
// Accessible by group clients, room clients, users.
func (s *RoomsService) GetRoomParticipants(ctx context.Context, roomIdOrName string, opt *RoomParticipantsOptions) ([]*UserListItem, *PaginatedResponse, error) {
	var u, err = getRoomResourcePath(roomIdOrName, getRoomParticipantsRoute)
	if err != nil {
		return nil, nil, err
	}

	opts, err := addUrlOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.Get(opts)
	if err != nil {
		return nil, nil, err
	}

	var participants *usersListResponse
	resp, err := s.client.Do(ctx, req, &participants)
	if err != nil {
		return nil, resp, err
	}

	return participants.Items, resp, nil
}

// Reply to a message in a room.
//
// Authentication required, with scope send_message.
// Accessible by users.
func (s *RoomsService) ReplyToRoomMessage(ctx context.Context, roomIdOrName string, messageId string, message string) (*PaginatedResponse, error) {
	var u, err = getRoomResourcePath(roomIdOrName, replyToRoomMessageRoute)
	if err != nil {
		return nil, err
	}

	req, err := s.client.Post(u, replyToMessageBody{messageId, message})
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Send a message to a room.
//
// Authentication required, with scope send_message.
// Accessible by users.
func (s *RoomsService) SendRoomMessage(ctx context.Context, roomIdOrName string, message string) (*RoomMessage, *PaginatedResponse, error) {
	var u, err = getRoomResourcePath(roomIdOrName, sendRoomMessageRoute)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.Post(u, sendMessageBody{message})
	if err != nil {
		return nil, nil, err
	}

	m := new(RoomMessage)
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// Gets all members for this private room.
//
// Authentication required, with scope view_room.
// Accessible by group clients, room clients, users.
func (s *RoomsService) GetRoomMembers(ctx context.Context, roomIdOrName string, opt *ListOptions) ([]*UserListItem, *PaginatedResponse, error) {
	var u, err = getRoomResourcePath(roomIdOrName, getRoomMembersRoute)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.Get(u)
	if err != nil {
		return nil, nil, err
	}

	var members *usersListResponse
	resp, err := s.client.Do(ctx, req, &members)
	if err != nil {
		return nil, resp, err
	}

	return members.Items, resp, nil
}


// Creates a new Room Object
func NewRoom(name string) *Room {
	r:= &Room{}
	r.RoomListItem = RoomListItem{}
	r.Name = name

	// Room access defaults to 'public'.
	r.Privacy = RoomPrivacyPublic

	return r
}

func getRoomResourcePath(roomIdOrName string, route string) (string, error)  {
	if roomIdOrName != "" {
		return fmt.Sprintf(route, roomIdOrName), nil
	} else {
		return "", emptyParam
	}
}

type roomsListResponse struct {
	Items []*RoomListItem `json:"items,omitempty"`
}

type topicBody struct {
	Topic string `json:"topic"`
}

type sendMessageBody struct {
	Message string `json:"message"`
}

type shareLinkBody struct {
	Message string `json:"message"`
	Link string `json:"link"`
}

type replyToMessageBody struct {
	MessageId string `json:"parentMessageId"`
	Message string `json:"message"`
}
