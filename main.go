package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"
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

	// CLI/config parsing
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command***REMOVED***
	Use:   "dprobe",
	Short: "dprobe",
	Long:  `dprobe`,
	Run: func(cmd *cobra.Command, args []string) ***REMOVED***

	***REMOVED***,
***REMOVED***

var (
	cfgDebug           bool
	cli                *client.Client
	cfgSlack           Slack
	cfgOutput          string
	cfgImageSprawl     uint32
	cfgContainerSprawl uint32
	message            string
)

// Slack is the configuration for writing feed data to slack channels
type Slack struct ***REMOVED***
	Channel string
	Token   string
***REMOVED***

func setFlags() ***REMOVED***
	rootCmd.PersistentFlags().StringVarP(&cfgOutput, "output", "o", "stdout", "Sets the output method (slack, or stdout)")
	rootCmd.PersistentFlags().Uint32VarP(&cfgImageSprawl, "isprawl", "i", 100, "Sets the minimum amount of images on a host to trip the image sprawl flag")
	rootCmd.PersistentFlags().Uint32VarP(&cfgContainerSprawl, "csprawl", "c", 100, "Sets the minimum amount of containers on a host to trip the container sprawl flag")
***REMOVED***

// PreInit initializes initializes cobra
func PreInit() ***REMOVED***
	setFlags()

	helpCmd := rootCmd.HelpFunc()

	var helpFlag bool

	newHelpCmd := func(c *cobra.Command, args []string) ***REMOVED***
		helpFlag = true
		helpCmd(c, args)
	***REMOVED***
	rootCmd.SetHelpFunc(newHelpCmd)

	err := rootCmd.Execute()

	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	if helpFlag ***REMOVED***
		os.Exit(0)
	***REMOVED***

	if cfgImageSprawl <= 3 ***REMOVED***
		log.Infof("%d is pretty low for image sprawl setting; setting to 10 instead...", cfgImageSprawl)
		cfgImageSprawl = 10
	***REMOVED***

	if cfgContainerSprawl <= 3 ***REMOVED***
		log.Infof("%d is pretty low for container sprawl setting; setting to 10 instead...", cfgContainerSprawl)
		cfgContainerSprawl = 10
	***REMOVED***
***REMOVED***

// InitViper initializes viper (configuration file) and links cobra and viper
func InitViper() ***REMOVED***
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("csprawl", rootCmd.PersistentFlags().Lookup("csprawl"))
	viper.BindPFlag("isprawl", rootCmd.PersistentFlags().Lookup("isprawl"))

	viper.SetConfigName("dprobe")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", "dprobe"))
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	err3 := viper.UnmarshalKey("slack", &cfgSlack)

	if err3 != nil ***REMOVED***
		log.Fatal(err3)
	***REMOVED***

	if cfgDebug ***REMOVED***
		log.Println(viper.AllSettings())
	***REMOVED***
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
func GetIPs() ([]string, error) ***REMOVED***
	addrs, err := net.InterfaceAddrs()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var strAddrs []string
	strAddrs = append(strAddrs, "IPs:")
	for i := range addrs ***REMOVED***
		strAddrs = append(strAddrs, addrs[i].String())
	***REMOVED***

	return strAddrs, nil
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
func ToSlack(data string) ***REMOVED***
	api := slack.New(cfgSlack.Token)

	params := slack.PostMessageParameters***REMOVED******REMOVED***

	_, _, err := api.PostMessage(cfgSlack.Channel, data, params)

	if err != nil ***REMOVED***
		log.Error(err)
	***REMOVED***
***REMOVED***

// MakeOutput formats output based on the type of output
func MakeOutput(data ...string) error ***REMOVED***
	var line string

	for m := range data ***REMOVED***
		if m != 0 ***REMOVED***
			line = fmt.Sprintf("%s %s", line, data[m])
		***REMOVED*** else ***REMOVED***
			line = fmt.Sprintf("%s%s", line, data[m])
		***REMOVED***
	***REMOVED***

	line = fmt.Sprintf("%s\n", line)
	message = message + line

	return nil
***REMOVED***

// SendOutput used to send the output to the defined location
func SendOutput(output string) error ***REMOVED***
	if output == "stdout" ***REMOVED***
		fmt.Println(message)
	***REMOVED*** else if output == "slack" ***REMOVED***

	***REMOVED*** else ***REMOVED***
		return fmt.Errorf("Invalid output format")
	***REMOVED***

	return nil
***REMOVED***

func main() ***REMOVED***
	PreInit()
	InitViper()

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

	MakeOutput("Host Information")
	iz1, err1 := GetHostname()
	if err1 != nil ***REMOVED***
		log.Fatal(err1)
	***REMOVED***
	MakeOutput("Hostname:", iz1)

	iz2, err2 := GetIPs()
	if err2 != nil ***REMOVED***
		log.Fatal(err2)
	***REMOVED***
	MakeOutput(iz2...)

	iz3, err3 := GetInstanceID()
	if err3 != nil ***REMOVED***
		log.Error(err3)
	***REMOVED***
	if len(iz3) > 0 ***REMOVED***
		MakeOutput("Instance ID:", iz3)
	***REMOVED***

	iz4, err4 := GetECSAgentVersion()
	if err4 != nil ***REMOVED***
		log.Error(err4)
	***REMOVED***
	if len(iz4) > 0 ***REMOVED***
		MakeOutput("ECS version:", iz4)
	***REMOVED***

	iz5, err5 := GetECSClusterName()
	if err5 != nil ***REMOVED***
		log.Error(err5)
	***REMOVED***
	if len(iz5) > 0 ***REMOVED***
		MakeOutput("ECS Cluster:", iz5)
	***REMOVED***

	MakeOutput("\n")
	MakeOutput("Docker Host Audit")
	iz6, err6 := HasContainerSprawl(cli, int(cfgContainerSprawl))
	if err6 != nil ***REMOVED***
		log.Error(err6)
	***REMOVED***
	strconv.FormatBool(iz6)
	MakeOutput("Container sprawl:", strconv.FormatBool(iz6))

	iz7, err7 := HasImageSprawl(cli, int(cfgImageSprawl))
	if err7 != nil ***REMOVED***
		log.Error(err7)
	***REMOVED***
	strconv.FormatBool(iz7)
	MakeOutput("Image sprawl:", strconv.FormatBool(iz7))

	iz17, err17 := HasStableDockerCEVersion()
	if err17 != nil ***REMOVED***
		log.Error(err17)
	***REMOVED***
	MakeOutput("Stable docker version:", strconv.FormatBool(iz17))

	iz18, err18 := HasLiveRestore(cli)
	if err18 != nil ***REMOVED***
		log.Error(err18)
	***REMOVED***
	MakeOutput("Live Restore:", strconv.FormatBool(iz18))

	iz19, err19 := FileOwnedByRoot("/var/lib/docker")
	if err19 != nil ***REMOVED***
		log.Error(err19)
	***REMOVED***
	MakeOutput("/var/lib/docker owned by root:", strconv.FormatBool(iz19))

	iz20, err20 := FileOwnedByRoot("/etc/docker")
	if err20 != nil ***REMOVED***
		log.Error(err20)
	***REMOVED***
	MakeOutput("/etc/docker owned by root:", strconv.FormatBool(iz20))

	iz21, err21 := FileOwnedByRoot("/etc/docker/daemon.json")
	if err21 != nil ***REMOVED***
		log.Error(err21)
	***REMOVED***
	MakeOutput("/etc/docker/daemon.json owned by root:", strconv.FormatBool(iz21))

	iz22, err22 := FileOwnedByRoot("/usr/bin/docker-containerd")
	if err22 != nil ***REMOVED***
		log.Error(err22)
	***REMOVED***
	MakeOutput("/usr/bin/docker-containerd owned by root:", strconv.FormatBool(iz22))

	iz23, err23 := FileOwnedByRoot("/usr/bin/docker-runc")
	if err23 != nil ***REMOVED***
		log.Error(err23)
	***REMOVED***
	MakeOutput("/usr/bin/docker-runc owned by root:", strconv.FormatBool(iz23))

	for c := range containers ***REMOVED***
		img := fmt.Sprintf("(%s)", containers[c].Image)
		MakeOutput("\n")
		MakeOutput("Container:", containers[c].ID, img)

		iz8, err8 := HasPrivilegedExecution(cli, containers[c].ID)
		if err8 != nil ***REMOVED***
			log.Error(err8)
		***REMOVED***
		MakeOutput("Privileged Execution:", strconv.FormatBool(iz8))

		iz9, err9 := HasExtendedCapabilities(cli, containers[c].ID)
		if err9 != nil ***REMOVED***
			log.Error(err9)
		***REMOVED***
		MakeOutput("Extended Capabilities:", strconv.FormatBool(iz9))

		iz10, err10 := HasHealthcheck(cli, containers[c].ID)
		if err10 != nil ***REMOVED***
			log.Error(err10)
		***REMOVED***
		MakeOutput("Memory limit:", strconv.FormatBool(iz10))

		iz11, err11 := HasSharedMountPropagation(cli, containers[c].ID)
		if err11 != nil ***REMOVED***
			log.Error(err11)
		***REMOVED***
		MakeOutput("Shared Propagation:", strconv.FormatBool(iz11))

		iz12, err12 := HasPrivilegedPorts(cli, containers[c].ID)
		if err12 != nil ***REMOVED***
			log.Error(err12)
		***REMOVED***
		MakeOutput("Privileged Ports:", strconv.FormatBool(iz12))

		iz13, err13 := HasUTSModeHost(cli, containers[c].ID)
		if err13 != nil ***REMOVED***
			log.Error(err13)
		***REMOVED***
		MakeOutput("UTS Mode Host:", strconv.FormatBool(iz13))

		iz14, err14 := HasIPCModeHost(cli, containers[c].ID)
		if err14 != nil ***REMOVED***
			log.Error(err14)
		***REMOVED***
		MakeOutput("IPC Mode Host:", strconv.FormatBool(iz14))

		iz15, err15 := HasProcessModeHost(cli, containers[c].ID)
		if err15 != nil ***REMOVED***
			log.Error(err15)
		***REMOVED***
		MakeOutput("Process Mode Host:", strconv.FormatBool(iz15))

		iz16, err16 := HasHostDevices(cli, containers[c].ID)
		if err16 != nil ***REMOVED***
			log.Error(err16)
		***REMOVED***
		MakeOutput("Has Host Devices:", strconv.FormatBool(iz16))
	***REMOVED***

	SendOutput(cfgOutput)
***REMOVED***
