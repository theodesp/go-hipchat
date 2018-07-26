package hipchat

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
