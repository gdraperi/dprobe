package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestAPINetworkGetDefaults(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	// By default docker daemon creates 3 networks. check if they are present
	defaults := []string***REMOVED***"bridge", "host", "none"***REMOVED***
	for _, nn := range defaults ***REMOVED***
		c.Assert(isNetworkAvailable(c, nn), checker.Equals, true)
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestAPINetworkCreateDelete(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	// Create a network
	name := "testnetwork"
	config := types.NetworkCreateRequest***REMOVED***
		Name: name,
		NetworkCreate: types.NetworkCreate***REMOVED***
			CheckDuplicate: true,
		***REMOVED***,
	***REMOVED***
	id := createNetwork(c, config, http.StatusCreated)
	c.Assert(isNetworkAvailable(c, name), checker.Equals, true)

	// delete the network and make sure it is deleted
	deleteNetwork(c, id, true)
	c.Assert(isNetworkAvailable(c, name), checker.Equals, false)
***REMOVED***

func (s *DockerSuite) TestAPINetworkCreateCheckDuplicate(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	name := "testcheckduplicate"
	configOnCheck := types.NetworkCreateRequest***REMOVED***
		Name: name,
		NetworkCreate: types.NetworkCreate***REMOVED***
			CheckDuplicate: true,
		***REMOVED***,
	***REMOVED***
	configNotCheck := types.NetworkCreateRequest***REMOVED***
		Name: name,
		NetworkCreate: types.NetworkCreate***REMOVED***
			CheckDuplicate: false,
		***REMOVED***,
	***REMOVED***

	// Creating a new network first
	createNetwork(c, configOnCheck, http.StatusCreated)
	c.Assert(isNetworkAvailable(c, name), checker.Equals, true)

	// Creating another network with same name and CheckDuplicate must fail
	createNetwork(c, configOnCheck, http.StatusConflict)

	// Creating another network with same name and not CheckDuplicate must succeed
	createNetwork(c, configNotCheck, http.StatusCreated)
***REMOVED***

func (s *DockerSuite) TestAPINetworkFilter(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	nr := getNetworkResource(c, getNetworkIDByName(c, "bridge"))
	c.Assert(nr.Name, checker.Equals, "bridge")
***REMOVED***

func (s *DockerSuite) TestAPINetworkInspectBridge(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	// Inspect default bridge network
	nr := getNetworkResource(c, "bridge")
	c.Assert(nr.Name, checker.Equals, "bridge")

	// run a container and attach it to the default bridge network
	out, _ := dockerCmd(c, "run", "-d", "--name", "test", "busybox", "top")
	containerID := strings.TrimSpace(out)
	containerIP := findContainerIP(c, "test", "bridge")

	// inspect default bridge network again and make sure the container is connected
	nr = getNetworkResource(c, nr.ID)
	c.Assert(nr.Driver, checker.Equals, "bridge")
	c.Assert(nr.Scope, checker.Equals, "local")
	c.Assert(nr.Internal, checker.Equals, false)
	c.Assert(nr.EnableIPv6, checker.Equals, false)
	c.Assert(nr.IPAM.Driver, checker.Equals, "default")
	c.Assert(nr.Containers[containerID], checker.NotNil)

	ip, _, err := net.ParseCIDR(nr.Containers[containerID].IPv4Address)
	c.Assert(err, checker.IsNil)
	c.Assert(ip.String(), checker.Equals, containerIP)
***REMOVED***

func (s *DockerSuite) TestAPINetworkInspectUserDefinedNetwork(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	// IPAM configuration inspect
	ipam := &network.IPAM***REMOVED***
		Driver: "default",
		Config: []network.IPAMConfig***REMOVED******REMOVED***Subnet: "172.28.0.0/16", IPRange: "172.28.5.0/24", Gateway: "172.28.5.254"***REMOVED******REMOVED***,
	***REMOVED***
	config := types.NetworkCreateRequest***REMOVED***
		Name: "br0",
		NetworkCreate: types.NetworkCreate***REMOVED***
			Driver:  "bridge",
			IPAM:    ipam,
			Options: map[string]string***REMOVED***"foo": "bar", "opts": "dopts"***REMOVED***,
		***REMOVED***,
	***REMOVED***
	id0 := createNetwork(c, config, http.StatusCreated)
	c.Assert(isNetworkAvailable(c, "br0"), checker.Equals, true)

	nr := getNetworkResource(c, id0)
	c.Assert(len(nr.IPAM.Config), checker.Equals, 1)
	c.Assert(nr.IPAM.Config[0].Subnet, checker.Equals, "172.28.0.0/16")
	c.Assert(nr.IPAM.Config[0].IPRange, checker.Equals, "172.28.5.0/24")
	c.Assert(nr.IPAM.Config[0].Gateway, checker.Equals, "172.28.5.254")
	c.Assert(nr.Options["foo"], checker.Equals, "bar")
	c.Assert(nr.Options["opts"], checker.Equals, "dopts")

	// delete the network and make sure it is deleted
	deleteNetwork(c, id0, true)
	c.Assert(isNetworkAvailable(c, "br0"), checker.Equals, false)
***REMOVED***

func (s *DockerSuite) TestAPINetworkConnectDisconnect(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	// Create test network
	name := "testnetwork"
	config := types.NetworkCreateRequest***REMOVED***
		Name: name,
	***REMOVED***
	id := createNetwork(c, config, http.StatusCreated)
	nr := getNetworkResource(c, id)
	c.Assert(nr.Name, checker.Equals, name)
	c.Assert(nr.ID, checker.Equals, id)
	c.Assert(len(nr.Containers), checker.Equals, 0)

	// run a container
	out, _ := dockerCmd(c, "run", "-d", "--name", "test", "busybox", "top")
	containerID := strings.TrimSpace(out)

	// connect the container to the test network
	connectNetwork(c, nr.ID, containerID)

	// inspect the network to make sure container is connected
	nr = getNetworkResource(c, nr.ID)
	c.Assert(len(nr.Containers), checker.Equals, 1)
	c.Assert(nr.Containers[containerID], checker.NotNil)

	// check if container IP matches network inspect
	ip, _, err := net.ParseCIDR(nr.Containers[containerID].IPv4Address)
	c.Assert(err, checker.IsNil)
	containerIP := findContainerIP(c, "test", "testnetwork")
	c.Assert(ip.String(), checker.Equals, containerIP)

	// disconnect container from the network
	disconnectNetwork(c, nr.ID, containerID)
	nr = getNetworkResource(c, nr.ID)
	c.Assert(nr.Name, checker.Equals, name)
	c.Assert(len(nr.Containers), checker.Equals, 0)

	// delete the network
	deleteNetwork(c, nr.ID, true)
***REMOVED***

func (s *DockerSuite) TestAPINetworkIPAMMultipleBridgeNetworks(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	// test0 bridge network
	ipam0 := &network.IPAM***REMOVED***
		Driver: "default",
		Config: []network.IPAMConfig***REMOVED******REMOVED***Subnet: "192.178.0.0/16", IPRange: "192.178.128.0/17", Gateway: "192.178.138.100"***REMOVED******REMOVED***,
	***REMOVED***
	config0 := types.NetworkCreateRequest***REMOVED***
		Name: "test0",
		NetworkCreate: types.NetworkCreate***REMOVED***
			Driver: "bridge",
			IPAM:   ipam0,
		***REMOVED***,
	***REMOVED***
	id0 := createNetwork(c, config0, http.StatusCreated)
	c.Assert(isNetworkAvailable(c, "test0"), checker.Equals, true)

	ipam1 := &network.IPAM***REMOVED***
		Driver: "default",
		Config: []network.IPAMConfig***REMOVED******REMOVED***Subnet: "192.178.128.0/17", Gateway: "192.178.128.1"***REMOVED******REMOVED***,
	***REMOVED***
	// test1 bridge network overlaps with test0
	config1 := types.NetworkCreateRequest***REMOVED***
		Name: "test1",
		NetworkCreate: types.NetworkCreate***REMOVED***
			Driver: "bridge",
			IPAM:   ipam1,
		***REMOVED***,
	***REMOVED***
	createNetwork(c, config1, http.StatusForbidden)
	c.Assert(isNetworkAvailable(c, "test1"), checker.Equals, false)

	ipam2 := &network.IPAM***REMOVED***
		Driver: "default",
		Config: []network.IPAMConfig***REMOVED******REMOVED***Subnet: "192.169.0.0/16", Gateway: "192.169.100.100"***REMOVED******REMOVED***,
	***REMOVED***
	// test2 bridge network does not overlap
	config2 := types.NetworkCreateRequest***REMOVED***
		Name: "test2",
		NetworkCreate: types.NetworkCreate***REMOVED***
			Driver: "bridge",
			IPAM:   ipam2,
		***REMOVED***,
	***REMOVED***
	createNetwork(c, config2, http.StatusCreated)
	c.Assert(isNetworkAvailable(c, "test2"), checker.Equals, true)

	// remove test0 and retry to create test1
	deleteNetwork(c, id0, true)
	createNetwork(c, config1, http.StatusCreated)
	c.Assert(isNetworkAvailable(c, "test1"), checker.Equals, true)

	// for networks w/o ipam specified, docker will choose proper non-overlapping subnets
	createNetwork(c, types.NetworkCreateRequest***REMOVED***Name: "test3"***REMOVED***, http.StatusCreated)
	c.Assert(isNetworkAvailable(c, "test3"), checker.Equals, true)
	createNetwork(c, types.NetworkCreateRequest***REMOVED***Name: "test4"***REMOVED***, http.StatusCreated)
	c.Assert(isNetworkAvailable(c, "test4"), checker.Equals, true)
	createNetwork(c, types.NetworkCreateRequest***REMOVED***Name: "test5"***REMOVED***, http.StatusCreated)
	c.Assert(isNetworkAvailable(c, "test5"), checker.Equals, true)

	for i := 1; i < 6; i++ ***REMOVED***
		deleteNetwork(c, fmt.Sprintf("test%d", i), true)
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestAPICreateDeletePredefinedNetworks(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	createDeletePredefinedNetwork(c, "bridge")
	createDeletePredefinedNetwork(c, "none")
	createDeletePredefinedNetwork(c, "host")
***REMOVED***

func createDeletePredefinedNetwork(c *check.C, name string) ***REMOVED***
	// Create pre-defined network
	config := types.NetworkCreateRequest***REMOVED***
		Name: name,
		NetworkCreate: types.NetworkCreate***REMOVED***
			CheckDuplicate: true,
		***REMOVED***,
	***REMOVED***
	createNetwork(c, config, http.StatusForbidden)
	deleteNetwork(c, name, false)
***REMOVED***

func isNetworkAvailable(c *check.C, name string) bool ***REMOVED***
	resp, body, err := request.Get("/networks")
	c.Assert(err, checker.IsNil)
	defer resp.Body.Close()
	c.Assert(resp.StatusCode, checker.Equals, http.StatusOK)

	nJSON := []types.NetworkResource***REMOVED******REMOVED***
	err = json.NewDecoder(body).Decode(&nJSON)
	c.Assert(err, checker.IsNil)

	for _, n := range nJSON ***REMOVED***
		if n.Name == name ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func getNetworkIDByName(c *check.C, name string) string ***REMOVED***
	var (
		v          = url.Values***REMOVED******REMOVED***
		filterArgs = filters.NewArgs()
	)
	filterArgs.Add("name", name)
	filterJSON, err := filters.ToJSON(filterArgs)
	c.Assert(err, checker.IsNil)
	v.Set("filters", filterJSON)

	resp, body, err := request.Get("/networks?" + v.Encode())
	c.Assert(resp.StatusCode, checker.Equals, http.StatusOK)
	c.Assert(err, checker.IsNil)

	nJSON := []types.NetworkResource***REMOVED******REMOVED***
	err = json.NewDecoder(body).Decode(&nJSON)
	c.Assert(err, checker.IsNil)
	var res string
	for _, n := range nJSON ***REMOVED***
		// Find exact match
		if n.Name == name ***REMOVED***
			res = n.ID
		***REMOVED***
	***REMOVED***
	c.Assert(res, checker.Not(checker.Equals), "")

	return res
***REMOVED***

func getNetworkResource(c *check.C, id string) *types.NetworkResource ***REMOVED***
	_, obj, err := request.Get("/networks/" + id)
	c.Assert(err, checker.IsNil)

	nr := types.NetworkResource***REMOVED******REMOVED***
	err = json.NewDecoder(obj).Decode(&nr)
	c.Assert(err, checker.IsNil)

	return &nr
***REMOVED***

func createNetwork(c *check.C, config types.NetworkCreateRequest, expectedStatusCode int) string ***REMOVED***
	resp, body, err := request.Post("/networks/create", request.JSONBody(config))
	c.Assert(err, checker.IsNil)
	defer resp.Body.Close()

	c.Assert(resp.StatusCode, checker.Equals, expectedStatusCode)

	if expectedStatusCode == http.StatusCreated ***REMOVED***
		var nr types.NetworkCreateResponse
		err = json.NewDecoder(body).Decode(&nr)
		c.Assert(err, checker.IsNil)

		return nr.ID
	***REMOVED***
	return ""
***REMOVED***

func connectNetwork(c *check.C, nid, cid string) ***REMOVED***
	config := types.NetworkConnect***REMOVED***
		Container: cid,
	***REMOVED***

	resp, _, err := request.Post("/networks/"+nid+"/connect", request.JSONBody(config))
	c.Assert(resp.StatusCode, checker.Equals, http.StatusOK)
	c.Assert(err, checker.IsNil)
***REMOVED***

func disconnectNetwork(c *check.C, nid, cid string) ***REMOVED***
	config := types.NetworkConnect***REMOVED***
		Container: cid,
	***REMOVED***

	resp, _, err := request.Post("/networks/"+nid+"/disconnect", request.JSONBody(config))
	c.Assert(resp.StatusCode, checker.Equals, http.StatusOK)
	c.Assert(err, checker.IsNil)
***REMOVED***

func deleteNetwork(c *check.C, id string, shouldSucceed bool) ***REMOVED***
	resp, _, err := request.Delete("/networks/" + id)
	c.Assert(err, checker.IsNil)
	defer resp.Body.Close()
	if !shouldSucceed ***REMOVED***
		c.Assert(resp.StatusCode, checker.Not(checker.Equals), http.StatusOK)
		return
	***REMOVED***
	c.Assert(resp.StatusCode, checker.Equals, http.StatusNoContent)
***REMOVED***
