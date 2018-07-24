package main

import (
	"context"
	"fmt"
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
	rooms, _, _ := h.Rooms.List(ctx, opts)
	for _, room := range rooms {
		fmt.Println(room.Name)
	}
}
