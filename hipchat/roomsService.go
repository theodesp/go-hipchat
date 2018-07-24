package hipchat

import (
	"context"
	"fmt"
)

const (
	listRoomsRoute = "room"
	getRoomRoute   = "room/%v"
)

// RoomsService handles communication with the room related
// methods of the HipChat API.
//
type RoomsService service

// RoomListItem represents a HipChat Room list item
type RoomListItem struct {
	ID         int64  `json:"id"`
	IsArchived bool   `json:"is_archived"`
	Name       string `json:"name"`
	Privacy    string `json:"privacy"`
	Version    string `json:"version"`

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
	*RoomListItem

	// XMPP/Jabber ID of the room.
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

	Owner *RoomOwner `json:"owner,omitempty"`

	// Statistics for this room.
	Statistics *struct{
		Links *struct{
			// The URL to use to retrieve room statistics
			Self string `json:"self"`
		} `json:"links,omitempty"`
	} `json:"statistics,omitempty"`
}

// RoomOwner represents a HipChat Room Owner
type RoomOwner struct {
	// User's @mention name
	MentionName string `json:"mention_name"`

	// An etag-like random version string.
	Version    string `json:"version"`

	// The user ID
	Id int64  `json:"id"`

	// URLs to retrieve user information
	Links *struct{
		// The link to use to retrieve the user information
		Self string `json:"self"`
	} `json:"links,omitempty"`

	// The display user name
	Name string `json:"name"`
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

// RoomsListResponse represents the response from the Rooms List request
type roomsListResponse struct {
	Items []*RoomListItem `json:"items,omitempty"`
}

// List non-archived rooms for this group.
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
func (s *RoomsService) Get(ctx context.Context, roomIdOrName string) (*Room, *PaginatedResponse, error) {
	var u string
	if roomIdOrName != "" {
		u = fmt.Sprintf(getRoomRoute, roomIdOrName)
	} else {
		return nil, nil, emptyParam
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
// Authentication required, with scope admin_room.
// Accessible by group clients, users.
func (s *RoomsService) Update(ctx context.Context, roomIdOrName string, room *Room) (*Room, *PaginatedResponse, error) {
	var u string
	if roomIdOrName != "" {
		u = fmt.Sprintf(getRoomRoute, roomIdOrName)
	} else {
		return nil, nil, emptyParam
	}

	req, err := s.client.Put(u, room)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(ctx, req, room)
	if err != nil {
		return nil, resp, err
	}

	return room, resp, nil
}

