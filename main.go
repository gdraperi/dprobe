package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	log "github.com/Sirupsen/logrus"
)

var cli *client.Client

// GetContainers returns all containers
// if all is false then only running containers are returned
func GetContainers(cli *client.Client, all bool) ([]types.Container, error) ***REMOVED***
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions***REMOVED***
		All: all,
	***REMOVED***)

	return containers, err
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

	containers, err := GetContainers(cli, true)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	fmt.Printf("%+v\n", containers)
***REMOVED***
