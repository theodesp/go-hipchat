package hipchat

import "context"

const (
	listUsersRoute = "user"
)

type UsersService service

// UserListItem represents a HipChat User list item
type UserListItem struct {
	// The user Id.
	Id int64 `json:"id"`

	// User's @mention name.
	MentionName string `json:"mention_name"`

	// The display user name
	Name string `json:"name"`

	// An etag-like random version string.
	Version string `json:"version"`

	// URLs to retrieve user information
	Links *UserLinks `json:"links,omitempty"`
}

type UserLinks struct {
	// The link to use to retrieve the user information
	Self string `json:"self"`
}

type usersListResponse struct {
	Items []*UserListItem `json:"items,omitempty"`
}

type UserListOptions struct {
	ListOptions

	// Include active guest users in response
	IncludeGuests bool `url:"include-guests"`

	// Include deleted users in response
	IncludeDeleted bool `url:"include-deleted"`
}

func (s *UsersService) ListUsers(ctx context.Context, opt *UserListOptions) ([]*UserListItem, *PaginatedResponse, error) {
	opts, err := addUrlOptions(listUsersRoute, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.Get(opts)
	if err != nil {
		return nil, nil, err
	}

	var users *usersListResponse
	resp, err := s.client.Do(ctx, req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users.Items, resp, nil
}
