package main

import (
	"context"
	"fmt"
	"go-hipchat/hipchat"
	"golang.org/x/oauth2"
	"os"
	"strconv"
)

var authToken = os.Getenv("HIPCHAT_AUTH_TOKEN")

func main() {
	// Context for cancellation
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: authToken})
	tc := oauth2.NewClient(ctx, ts)

	h := hipchat.NewClient(tc)
	opts := &hipchat.RoomsListOptions{}
	opts.IncludeArchived = false
	opts.MaxResults = 1
	rooms, _, _ := h.Rooms.ListRooms(ctx, opts)
	for _, room := range rooms {
		item, _, _ := h.Rooms.Get(ctx, strconv.FormatInt(room.ID, 10))
		fmt.Println(item)
	}
}
