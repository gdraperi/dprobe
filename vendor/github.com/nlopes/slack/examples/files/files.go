package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() ***REMOVED***
	api := slack.New("YOUR_TOKEN_HERE")
	params := slack.FileUploadParameters***REMOVED***
		Title: "Batman Example",
		//Filetype: "txt",
		File: "example.txt",
		//Content:  "Nan Nan Nan Nan Nan Nan Nan Nan Batman",
	***REMOVED***
	file, err := api.UploadFile(params)
	if err != nil ***REMOVED***
		fmt.Printf("%s\n", err)
		return
	***REMOVED***
	fmt.Printf("Name: %s, URL: %s\n", file.Name, file.URL)

	err = api.DeleteFile(file.ID)
	if err != nil ***REMOVED***
		fmt.Printf("%s\n", err)
		return
	***REMOVED***
	fmt.Printf("File %s deleted successfully.\n", file.Name)
***REMOVED***
