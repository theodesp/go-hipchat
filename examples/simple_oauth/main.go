package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"os"
)

var authToken = os.Getenv("HIPCHAT_AUTH_TOKEN")
var baseHost = "https://api.hipchat.com/v2/"

func main() {
	// Context for cancellation
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: authToken})
	tc := oauth2.NewClient(ctx, ts)

	resp, _ := tc.Get(baseHost + "room")
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
	}
}
