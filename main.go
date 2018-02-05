package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
	"syscall"

	// Slack API
	"github.com/nlopes/slack"

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

var (
	cli      *client.Client
	cfgSlack Slack
)

// Slack is the configuration for writing feed data to slack channels
type Slack struct ***REMOVED***
	Channel string
	Token   string
***REMOVED***

// GetContainers returns all containers
// if all is false then only running containers are returned
func GetContainers(cli *client.Client, all bool) ([]types.Container, error) ***REMOVED***
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions***REMOVED***
		All: all,
	***REMOVED***)

	return containers, err
***REMOVED***

// HasContainerSprawl takes max amount of containers; if the total tops this
// then there is container sprawl
func HasContainerSprawl(cli *client.Client, sprawl_amount int) (bool, error) ***REMOVED***
	containers, err := GetContainers(cli, true)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if len(containers) >= sprawl_amount ***REMOVED***
		return true, nil
	***REMOVED***

	return false, nil
***REMOVED***

// GetImages returns all images on the host
func GetImages(cli *client.Client, all bool) ([]types.ImageSummary, error) ***REMOVED***
	images, err := cli.ImageList(context.Background(), types.ImageListOptions***REMOVED***
		All: all,
	***REMOVED***)

	return images, err
***REMOVED***

// HasImageSprawl takes max amount of images; if the total tops this
// then there is image sprawl
func HasImageSprawl(cli *client.Client, sprawl_amount int) (bool, error) ***REMOVED***
	images, err := GetImages(cli, true)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if len(images) >= sprawl_amount ***REMOVED***
		return true, nil
	***REMOVED***

	return false, nil
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
		if v[i] == a.Version ***REMOVED***
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

// GetFileStats takes a file/directory and returns its file info stat
func GetFileStats(fname string) (os.FileInfo, error) ***REMOVED***
	fd, err := os.Stat(fname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return fd, nil
***REMOVED***

// FileOwnedByRoot returns true if fname owned by root
func FileOwnedByRoot(fname string) (bool, error) ***REMOVED***
	fd, err := GetFileStats(fname)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if fd.Sys().(*syscall.Stat_t).Uid == 0 ***REMOVED***
		return true, nil
	***REMOVED***

	return false, nil
***REMOVED***

// GetHostname returns the systems hostname
func GetHostname() (string, error) ***REMOVED***
	return os.Hostname()
***REMOVED***

// GetIPs returns the local systems IP addresses
func GetIPs() ([]net.Addr, error) ***REMOVED***
	addrs, err := net.InterfaceAddrs()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return addrs, nil
***REMOVED***

// GetInstanceID returns the hosts AWS instance ID
func GetInstanceID() (string, error) ***REMOVED***
	doc, err := goquery.NewDocument("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return doc.Text(), nil
***REMOVED***

// GetECSAgentVersion returns the running ECS agent version
func GetECSAgentVersion() (string, error) ***REMOVED***
	doc, err := goquery.NewDocument("http://127.0.0.1:51678/v1/metadata")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	data := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	dec := json.NewDecoder(strings.NewReader(doc.Text()))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	v, err2 := jq.String("Version")
	if err2 != nil ***REMOVED***
		return "", err2
	***REMOVED***

	re := regexp.MustCompile("v[0-9]***REMOVED***1,2***REMOVED***.[0-9]***REMOVED***1,2***REMOVED***.[0-9]***REMOVED***1,2***REMOVED***?")
	m := re.FindStringSubmatch(v)

	if m != nil ***REMOVED***
		return m[0], nil
	***REMOVED***

	return v, nil
***REMOVED***

// GetECSClusterName returns the name of the current ECS cluster
func GetECSClusterName() (string, error) ***REMOVED***
	doc, err := goquery.NewDocument("http://127.0.0.1:51678/v1/metadata")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	data := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	dec := json.NewDecoder(strings.NewReader(doc.Text()))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	v, err2 := jq.String("Cluster")
	if err2 != nil ***REMOVED***
		return "", err2
	***REMOVED***

	return v, nil
***REMOVED***

// ToSlack writes the parsed feed data to a slack channel
func ToSlack(message string) ***REMOVED***
	api := slack.New(cfgSlack.Token)

	params := slack.PostMessageParameters***REMOVED******REMOVED***

	_, _, err := api.PostMessage(cfgSlack.Channel, message, params)

	if err != nil ***REMOVED***
		log.Error(err)
	***REMOVED***
***REMOVED***

func main() ***REMOVED***
	iz1, _ := GetHostname()
	fmt.Println(iz1)

	iz2, _ := GetIPs()
	fmt.Println(iz2)

	iz3, _ := GetInstanceID()
	fmt.Println(iz3)

	iz4, _ := GetECSAgentVersion()
	fmt.Printf("ECS version: %s\n", iz4)

	iz5, _ := GetECSClusterName()
	fmt.Printf("ECS Cluster: %s\n", iz5)

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

	zztx, _ := HasContainerSprawl(cli, 3)
	fmt.Printf("Container sprawl: %t\n", zztx)

	zztxs, _ := HasImageSprawl(cli, 2)
	fmt.Printf("Image sprawl: %t\n", zztxs)

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

	ff1, _ := FileOwnedByRoot("/var/lib/docker")
	fmt.Printf("/var/lib/docker owned by root: %t\n", ff1)
	ff2, _ := FileOwnedByRoot("/etc/docker")
	fmt.Printf("/var/lib/docker owned by root: %t\n", ff2)
	ff3, _ := FileOwnedByRoot("/etc/docker/daemon.json")
	fmt.Printf("/etc/docker/daemon.json owned by root: %t\n", ff3)
	ff4, _ := FileOwnedByRoot("/usr/bin/docker-containerd")
	fmt.Printf("/usr/bin/docker-containerd owned by root: %t\n", ff4)
	ff5, _ := FileOwnedByRoot("/usr/bin/docker-runc")
	fmt.Printf("/usr/bin/docker-runc owned by root: %t\n", ff5)
***REMOVED***
