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
func GetImages(cli *client.Client, all bool) ([]types.ImageSummary, error) ***REMOVED***
	images, err := cli.ImageList(context.Background(), types.ImageListOptions***REMOVED***
		All: all,
	***REMOVED***)

	return images, err
***REMOVED***

// InspectContainer returns information about the container back
// id is the id of the container
func InspectContainer(cli *client.Client, id string) (types.ContainerJSON, error) ***REMOVED***
	inspection, err := cli.ContainerInspect(context.Background(), id)

	return inspection, err
***REMOVED***

// HasPrivilegedExecution returns true/false if the container has
// privileged execution
func HasPrivilegedExecution(cli *client.Client, id string) (bool, error) ***REMOVED***
	c_insp, err := InspectContainer(cli, id)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	return c_insp.HostConfig.Privileged, nil
***REMOVED***

// HasExtendedCapabilities returns true/false if the container has extended capabilities
func HasExtendedCapabilities(cli *client.Client, id string) (bool, error) ***REMOVED***
	c_insp, err := InspectContainer(cli, id)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if len(c_insp.HostConfig.CapAdd) > 0 ***REMOVED***
		return true, nil
	***REMOVED***

	return false, nil
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

	images, err := GetImages(cli, true)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	fmt.Printf("%+v\n", images)

	for c := range containers ***REMOVED***
		t, _ := HasPrivilegedExecution(cli, containers[c].ID)
		fmt.Println(t)
		s, _ := HasExtendedCapabilities(cli, containers[c].ID)
		fmt.Println(s)
	***REMOVED***
***REMOVED***
