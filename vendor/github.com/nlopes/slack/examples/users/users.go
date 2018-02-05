package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() ***REMOVED***
	api := slack.New("YOUR_TOKEN_HERE")
	user, err := api.GetUserInfo("U023BECGF")
	if err != nil ***REMOVED***
		fmt.Printf("%s\n", err)
		return
	***REMOVED***
	fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
***REMOVED***
