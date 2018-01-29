package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var cli *client.Client

// GetContainers returns all containers
// if all is false then only running containers are returned
func GetContainers(cli *client.Client, all bool) ***REMOVED***
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions***REMOVED***
		All: all,
	***REMOVED***)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	for _, container := range containers ***REMOVED***
		fmt.Println(container.ID)
		fmt.Printf("%+v\n", container)
	***REMOVED***
***REMOVED***

// GetImages returns all images on the host
func GetImages() ***REMOVED***

***REMOVED***

// InspectContainer returns information about the container back
// id is the id of the container
func InspectContainer(id string) ***REMOVED***

***REMOVED***

func main() ***REMOVED***
	var err error

	cli, err = client.NewEnvClient()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	GetContainers(cli, true)
***REMOVED***
