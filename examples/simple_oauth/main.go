package main

import (
	"context"
	"go-hipchat/hipchat"
	"golang.org/x/oauth2"
	"os"
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
	opts.StartIndex = 10
	opts.MaxResults = 1

	h.Rooms.SetRoomTopic(ctx, "TestRoom_1", "This is a topic")
}
