package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	// Docker API
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	// goquery (for HTTP requests)
	"github.com/PuerkitoBio/goquery"

	// Logging
	log "github.com/Sirupsen/logrus"

	// jsonq (easy json parsing)
	"github.com/jmoiron/jsonq"
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

// GetStableDockerCEVersions returns a list of the stable docker versions
func GetStableDockerCEVersions() ([]string, error) ***REMOVED***
	doc, err := goquery.NewDocument("https://docs.docker.com/release-notes/docker-ce/")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Find the review items
	var versions []string
	doc.Find("#my_toc > li:first-child > ul > li > a").Each(func(i int, s *goquery.Selection) ***REMOVED***
		// For each item found, get the release
		release := strings.Fields(s.Text())[0]

		versions = append(versions, release)
	***REMOVED***)

	return versions, nil
***REMOVED***

// GetDockerServerVersion returns the local docker server version
func GetDockerServerVersion(cli *client.Client) (types.Version, error) ***REMOVED***
	version, err := cli.ServerVersion(context.Background())

	return version, err
***REMOVED***

func HasStableDockerCEVersion() (bool, error) ***REMOVED***
	v, err1 := GetStableDockerCEVersions()
	if err1 != nil ***REMOVED***
		return false, err1
	***REMOVED***

	a, err2 := GetDockerServerVersion(cli)
	if err2 != nil ***REMOVED***
		return false, err2
	***REMOVED***

	for i := range v ***REMOVED***
		if v[i] == a.Components[0].Version ***REMOVED***
			return true, nil
		***REMOVED***
	***REMOVED***

	return false, fmt.Errorf("%s is not in the list of stable docker CE versions", a)
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

// HasHealthcheck returns true if there is a healthcheck set for the container
func HasHealthcheck(cli *client.Client, id string) (bool, error) ***REMOVED***
	c_insp, err := InspectContainer(cli, id)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if c_insp.Config.Healthcheck == nil ***REMOVED***
		return false, nil
	***REMOVED***

	return true, nil
***REMOVED***

// HasSharedMountPropagation returns true if any of the mount points on a
// container are shraed propagated
func HasSharedMountPropagation(cli *client.Client, id string) (bool, error) ***REMOVED***
	c_insp, err := InspectContainer(cli, id)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	var shared_prop_mounts bool
	for mount := range c_insp.Mounts ***REMOVED***
		if c_insp.Mounts[mount].Propagation == "shared" ***REMOVED***
			shared_prop_mounts = true
		***REMOVED***
	***REMOVED***

	return shared_prop_mounts, nil
***REMOVED***

// HasPrivilegedPorts returns true if the container is bound to a privileged port
func HasPrivilegedPorts(cli *client.Client, id string) (bool, error) ***REMOVED***
	c_insp, err := InspectContainer(cli, id)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	var priv_port bool
	for k := range c_insp.NetworkSettings.Ports ***REMOVED***
		if k.Int() <= 1024 ***REMOVED***
			priv_port = true
		***REMOVED***
	***REMOVED***

	return priv_port, nil
***REMOVED***

// HasUTSModeHost returns true if any containers UTSMode is "host"
func HasUTSModeHost(cli *client.Client, id string) (bool, error) ***REMOVED***
	c_insp, err := InspectContainer(cli, id)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if c_insp.HostConfig.UTSMode == "host" ***REMOVED***
		return true, nil
	***REMOVED***

	return false, nil
***REMOVED***

// HasIPCModeHost returns true if any containers IPCMode is "host"
func HasIPCModeHost(cli *client.Client, id string) (bool, error) ***REMOVED***
	c_insp, err := InspectContainer(cli, id)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if c_insp.HostConfig.IpcMode == "host" ***REMOVED***
		return true, nil
	***REMOVED***

	return false, nil
***REMOVED***

// HasProcessModeHost returns true if any containers Process Mode is "host"
func HasProcessModeHost(cli *client.Client, id string) (bool, error) ***REMOVED***
	c_insp, err := InspectContainer(cli, id)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if c_insp.HostConfig.PidMode == "host" ***REMOVED***
		return true, nil
	***REMOVED***

	return false, nil
***REMOVED***

// HasHostDevices returns true if a container has access to host devices
func HasHostDevices(cli *client.Client, id string) (bool, error) ***REMOVED***
	c_insp, err := InspectContainer(cli, id)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if len(c_insp.HostConfig.Devices) > 0 ***REMOVED***
		return true, nil
	***REMOVED***

	return false, nil
***REMOVED***

// GetServerInfo returns information about the server
func GetServerInfo(cli *client.Client) (types.Info, error) ***REMOVED***
	s_info, err := cli.Info(context.Background())
	if err != nil ***REMOVED***
		return s_info, err
	***REMOVED***

	return s_info, nil
***REMOVED***

// HasLiveRestore checks if the underlying docker server has --live-restore enabled
func HasLiveRestore(cli *client.Client) (bool, error) ***REMOVED***
	s_info, err := GetServerInfo(cli)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	return s_info.LiveRestoreEnabled, nil
***REMOVED***

// GetContainerStats returns a jq object that can be used to query container stats
func GetContainerStats(cli *client.Client, id string) (*jsonq.JsonQuery, error) ***REMOVED***
	c_stats, err := cli.ContainerStats(context.Background(), id, false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer c_stats.Body.Close()

	b, err2 := ioutil.ReadAll(c_stats.Body)
	if err2 != nil ***REMOVED***
		return nil, err2
	***REMOVED***

	data := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	return jq, nil
***REMOVED***

// HasMemoryLimit returns true if there is a memory limit on the container
func HasMemoryLimit(cli *client.Client, id string) (bool, error) ***REMOVED***
	jq, err1 := GetContainerStats(cli, id)
	if err1 != nil ***REMOVED***
		return false, err1
	***REMOVED***

	limit, err := jq.Int("memory_stats", "limit")
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if limit > 0 ***REMOVED***
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
	defer cli.Close()

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

		z, _ := HasMemoryLimit(cli, containers[c].ID)
		fmt.Printf("Memory limit: %t\n", z)

		y, _ := HasHealthcheck(cli, containers[c].ID)
		fmt.Printf("Health check: %t\n", y)

		x, _ := HasSharedMountPropagation(cli, containers[c].ID)
		fmt.Printf("Shared propagation: %t\n", x)

		xuu, _ := HasPrivilegedPorts(cli, containers[c].ID)
		fmt.Printf("Priv port: %t\n", xuu)

		xuuv, _ := HasUTSModeHost(cli, containers[c].ID)
		fmt.Printf("Exposed host UTS: %t\n", xuuv)

		xuuvx, _ := HasIPCModeHost(cli, containers[c].ID)
		fmt.Printf("Exposed host IPC: %t\n", xuuvx)

		aas, _ := HasProcessModeHost(cli, containers[c].ID)
		fmt.Printf("Exposed host Processes: %t\n", aas)

		abc1, _ := HasHostDevices(cli, containers[c].ID)
		fmt.Printf("Exposed host Devices: %t\n", abc1)
	***REMOVED***

	v, _ := GetStableDockerCEVersions()
	fmt.Println(v)
	a, _ := GetDockerServerVersion(cli)
	fmt.Printf("%+v\n", a)

	b, _ := HasStableDockerCEVersion()
	fmt.Println(b)

	HasLiveRestore(cli)
***REMOVED***
