package main

import (
	"context"
	"go-hipchat/hipchat"
	"golang.org/x/oauth2"
	"os"
	"fmt"
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

	r, _ := h.Rooms.DeleteRoom(ctx, "TestRoom_1")

	fmt.Println(r.StatusCode)
}
