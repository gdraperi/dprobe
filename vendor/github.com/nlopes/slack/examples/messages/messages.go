package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() ***REMOVED***
	api := slack.New("YOUR_TOKEN_HERE")
	params := slack.PostMessageParameters***REMOVED******REMOVED***
	attachment := slack.Attachment***REMOVED***
		Pretext: "some pretext",
		Text:    "some text",
		// Uncomment the following part to send a field too
		/*
			Fields: []slack.AttachmentField***REMOVED***
				slack.AttachmentField***REMOVED***
					Title: "a",
					Value: "no",
				***REMOVED***,
			***REMOVED***,
		*/
	***REMOVED***
	params.Attachments = []slack.Attachment***REMOVED***attachment***REMOVED***
	channelID, timestamp, err := api.PostMessage("CHANNEL_ID", "Some text", params)
	if err != nil ***REMOVED***
		fmt.Printf("%s\n", err)
		return
	***REMOVED***
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
***REMOVED***
