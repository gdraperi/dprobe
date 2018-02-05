package main

import (
	"flag"
	"fmt"

	"github.com/nlopes/slack"
)

func main() ***REMOVED***
	var (
		apiToken string
		debug    bool
	)

	flag.StringVar(&apiToken, "token", "YOUR_TOKEN_HERE", "Your Slack API Token")
	flag.BoolVar(&debug, "debug", false, "Show JSON output")
	flag.Parse()

	api := slack.New(apiToken)
	if debug ***REMOVED***
		api.SetDebug(true)
	***REMOVED***

	// Get all stars for the usr.
	params := slack.NewStarsParameters()
	starredItems, _, err := api.GetStarred(params)
	if err != nil ***REMOVED***
		fmt.Printf("Error getting stars: %s\n", err)
		return
	***REMOVED***
	for _, s := range starredItems ***REMOVED***
		var desc string
		switch s.Type ***REMOVED***
		case slack.TYPE_MESSAGE:
			desc = s.Message.Text
		case slack.TYPE_FILE:
			desc = s.File.Name
		case slack.TYPE_FILE_COMMENT:
			desc = s.File.Name + " - " + s.Comment.Comment
		case slack.TYPE_CHANNEL, slack.TYPE_IM, slack.TYPE_GROUP:
			desc = s.Channel
		***REMOVED***
		fmt.Printf("Starred %s: %s\n", s.Type, desc)
	***REMOVED***
***REMOVED***
