package hipchat

import (
	"context"
)

const (
	roomsListPrefix = "room"
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

	Links *struct {
		// The URL to use to retrieve the full room information
		Self string `json:"self"`

		// The URL to use to retrieve webhooks registered for this room
		Webhooks string `json:"webhooks"`

		// The URL to use to retrieve members for this room. Only available for private rooms.
		Members string `json:"members,omitempty"`

		//The URL to use to retrieve participants for this room
		Participants string `json:"participants"`
	} `json:"links,omitempty"`
}

// RoomsListOptions specifies the optional parameters to the
// RoomService.ListRooms
type RoomsListOptions struct {
	// Filter out private rooms
	IncludePrivate bool `url:"	include-private,omitempty"`

	// Filter rooms
	IncludeArchived bool `url:"	include-archived,omitempty"`
	ListOptions
}

// RoomsListResponse represents the response from the Rooms List request
type roomsListResponse struct {
	Items []*RoomListItem `json:"items,omitempty"`
}

// List non-archived rooms for this group.
// Authentication required, with scope view_group or view_room.
// Accessible by group clients, users.
func (s *RoomsService) List(ctx context.Context, opt *RoomsListOptions) ([]*RoomListItem, *PaginatedResponse, error) {
	opts, err := addUrlOptions(roomsListPrefix, opt)
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
