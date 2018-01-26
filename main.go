package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// GetContainers returns all containers
// if all is false then only running containers are returned
func GetContainers(all bool) ***REMOVED***

***REMOVED***

// InspectContainer returns information about the container back
// id is the id of the container
func InspectContainer(id string) ***REMOVED***

***REMOVED***

func main() ***REMOVED***
	cli, err := client.NewEnvClient()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	for _, container := range containers ***REMOVED***
		fmt.Println(container.ID)
	***REMOVED***
***REMOVED***
