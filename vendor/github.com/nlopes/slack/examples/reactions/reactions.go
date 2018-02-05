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

	var (
		postAsUserName  string
		postAsUserID    string
		postToUserName  string
		postToUserID    string
		postToChannelID string
	)

	// Find the user to post as.
	authTest, err := api.AuthTest()
	if err != nil ***REMOVED***
		fmt.Printf("Error getting channels: %s\n", err)
		return
	***REMOVED***

	// Post as the authenticated user.
	postAsUserName = authTest.User
	postAsUserID = authTest.UserID

	// Posting to DM with self causes a conversation with slackbot.
	postToUserName = authTest.User
	postToUserID = authTest.UserID

	// Find the channel.
	_, _, chanID, err := api.OpenIMChannel(postToUserID)
	if err != nil ***REMOVED***
		fmt.Printf("Error opening IM: %s\n", err)
		return
	***REMOVED***
	postToChannelID = chanID

	fmt.Printf("Posting as %s (%s) in DM with %s (%s), channel %s\n", postAsUserName, postAsUserID, postToUserName, postToUserID, postToChannelID)

	// Post a message.
	postParams := slack.PostMessageParameters***REMOVED******REMOVED***
	channelID, timestamp, err := api.PostMessage(postToChannelID, "Is this any good?", postParams)
	if err != nil ***REMOVED***
		fmt.Printf("Error posting message: %s\n", err)
		return
	***REMOVED***

	// Grab a reference to the message.
	msgRef := slack.NewRefToMessage(channelID, timestamp)

	// React with :+1:
	if err := api.AddReaction("+1", msgRef); err != nil ***REMOVED***
		fmt.Printf("Error adding reaction: %s\n", err)
		return
	***REMOVED***

	// React with :-1:
	if err := api.AddReaction("cry", msgRef); err != nil ***REMOVED***
		fmt.Printf("Error adding reaction: %s\n", err)
		return
	***REMOVED***

	// Get all reactions on the message.
	msgReactions, err := api.GetReactions(msgRef, slack.NewGetReactionsParameters())
	if err != nil ***REMOVED***
		fmt.Printf("Error getting reactions: %s\n", err)
		return
	***REMOVED***
	fmt.Printf("\n")
	fmt.Printf("%d reactions to message...\n", len(msgReactions))
	for _, r := range msgReactions ***REMOVED***
		fmt.Printf("  %d users say %s\n", r.Count, r.Name)
	***REMOVED***

	// List all of the users reactions.
	listReactions, _, err := api.ListReactions(slack.NewListReactionsParameters())
	if err != nil ***REMOVED***
		fmt.Printf("Error listing reactions: %s\n", err)
		return
	***REMOVED***
	fmt.Printf("\n")
	fmt.Printf("All reactions by %s...\n", authTest.User)
	for _, item := range listReactions ***REMOVED***
		fmt.Printf("%d on a %s...\n", len(item.Reactions), item.Type)
		for _, r := range item.Reactions ***REMOVED***
			fmt.Printf("  %s (along with %d others)\n", r.Name, r.Count-1)
		***REMOVED***
	***REMOVED***

	// Remove the :cry: reaction.
	err = api.RemoveReaction("cry", msgRef)
	if err != nil ***REMOVED***
		fmt.Printf("Error remove reaction: %s\n", err)
		return
	***REMOVED***

	// Get all reactions on the message.
	msgReactions, err = api.GetReactions(msgRef, slack.NewGetReactionsParameters())
	if err != nil ***REMOVED***
		fmt.Printf("Error getting reactions: %s\n", err)
		return
	***REMOVED***
	fmt.Printf("\n")
	fmt.Printf("%d reactions to message after removing cry...\n", len(msgReactions))
	for _, r := range msgReactions ***REMOVED***
		fmt.Printf("  %d users say %s\n", r.Count, r.Name)
	***REMOVED***
***REMOVED***
