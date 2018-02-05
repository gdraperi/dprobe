package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() ***REMOVED***
	api := slack.New("YOUR_TOKEN_HERE")
	channels, err := api.GetChannels(false)
	if err != nil ***REMOVED***
		fmt.Printf("%s\n", err)
		return
	***REMOVED***
	for _, channel := range channels ***REMOVED***
		fmt.Println(channel.Name)
		// channel is of type conversation & groupConversation
		// see all available methods in `conversation.go`
	***REMOVED***
***REMOVED***
