package hipchat

import (
	"github.com/google/go-querystring/query"
	"net/url"
	"reflect"
)

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addUrlOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func makeEmptyRoom() *Room  {
	roomListItem := &RoomListItem{ID: int64(1), IsArchived: false, Name: "", Privacy: "", Version: ""}
	room := &Room{
		roomListItem,
		"","","", false, "",
		false, "", "", nil, nil,
	}

	return room
}
