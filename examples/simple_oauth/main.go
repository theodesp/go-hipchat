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
	opts := &hipchat.RoomParticipantsOptions{}
	opts.StartIndex = 0
	opts.MaxResults = 100

	m, _, err := h.Rooms.SendRoomMessage(ctx, "TestRoom_1", "How are you?")
	fmt.Println(err)
	fmt.Println(m)
}
