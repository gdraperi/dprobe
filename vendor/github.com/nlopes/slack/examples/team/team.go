package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() ***REMOVED***
	api := slack.New("YOUR_TOKEN_HERE")
	//Example for single user
	billingActive, err := api.GetBillableInfo("U023BECGF")
	if err != nil ***REMOVED***
		fmt.Printf("%s\n", err)
		return
	***REMOVED***
	fmt.Printf("ID: U023BECGF, BillingActive: %v\n\n\n", billingActive["U023BECGF"])

	//Example for team
	billingActiveForTeam, _ := api.GetBillableInfoForTeam()
	for id, value := range billingActiveForTeam ***REMOVED***
		fmt.Printf("ID: %v, BillingActive: %v\n", id, value)
	***REMOVED***

***REMOVED***
