package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() ***REMOVED***
	api := slack.New("YOUR_TOKEN_HERE")

	userID := "USER_ID"

	_, _, channelID, err := api.OpenIMChannel(userID)

	if err != nil ***REMOVED***
		fmt.Printf("%s\n", err)
	***REMOVED***

	api.PostMessage(channelID, "Hello World!", slack.PostMessageParameters***REMOVED******REMOVED***)
***REMOVED***
