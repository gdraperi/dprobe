package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
	"golang.org/x/net/context"
)

var expectedNetworkInterfaceStats = strings.Split("rx_bytes rx_dropped rx_errors rx_packets tx_bytes tx_dropped tx_errors tx_packets", " ")

func (s *DockerSuite) TestAPIStatsNoStreamGetCpu(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "run", "-d", "busybox", "/bin/sh", "-c", "while true;usleep 100; do echo 'Hello'; done")

	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)
	resp, body, err := request.Get(fmt.Sprintf("/containers/%s/stats?stream=false", id))
	c.Assert(err, checker.IsNil)
	c.Assert(resp.StatusCode, checker.Equals, http.StatusOK)
	c.Assert(resp.Header.Get("Content-Type"), checker.Equals, "application/json")

	var v *types.Stats
	err = json.NewDecoder(body).Decode(&v)
	c.Assert(err, checker.IsNil)
	body.Close()

	var cpuPercent = 0.0

	if testEnv.OSType != "windows" ***REMOVED***
		cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta := float64(v.CPUStats.SystemUsage - v.PreCPUStats.SystemUsage)
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	***REMOVED*** else ***REMOVED***
		// Max number of 100ns intervals between the previous time read and now
		possIntervals := uint64(v.Read.Sub(v.PreRead).Nanoseconds()) // Start with number of ns intervals
		possIntervals /= 100                                         // Convert to number of 100ns intervals
		possIntervals *= uint64(v.NumProcs)                          // Multiple by the number of processors

		// Intervals used
		intervalsUsed := v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage

		// Percentage avoiding divide-by-zero
		if possIntervals > 0 ***REMOVED***
			cpuPercent = float64(intervalsUsed) / float64(possIntervals) * 100.0
		***REMOVED***
	***REMOVED***

	c.Assert(cpuPercent, check.Not(checker.Equals), 0.0, check.Commentf("docker stats with no-stream get cpu usage failed: was %v", cpuPercent))
***REMOVED***

func (s *DockerSuite) TestAPIStatsStoppedContainerInGoroutines(c *check.C) ***REMOVED***
	out, _ := dockerCmd(c, "run", "-d", "busybox", "/bin/sh", "-c", "echo 1")
	id := strings.TrimSpace(out)

	getGoRoutines := func() int ***REMOVED***
		_, body, err := request.Get(fmt.Sprintf("/info"))
		c.Assert(err, checker.IsNil)
		info := types.Info***REMOVED******REMOVED***
		err = json.NewDecoder(body).Decode(&info)
		c.Assert(err, checker.IsNil)
		body.Close()
		return info.NGoroutines
	***REMOVED***

	// When the HTTP connection is closed, the number of goroutines should not increase.
	routines := getGoRoutines()
	_, body, err := request.Get(fmt.Sprintf("/containers/%s/stats", id))
	c.Assert(err, checker.IsNil)
	body.Close()

	t := time.After(30 * time.Second)
	for ***REMOVED***
		select ***REMOVED***
		case <-t:
			c.Assert(getGoRoutines(), checker.LessOrEqualThan, routines)
			return
		default:
			if n := getGoRoutines(); n <= routines ***REMOVED***
				return
			***REMOVED***
			time.Sleep(200 * time.Millisecond)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestAPIStatsNetworkStats(c *check.C) ***REMOVED***
	testRequires(c, SameHostDaemon)

	out := runSleepingContainer(c)
	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)

	// Retrieve the container address
	net := "bridge"
	if testEnv.OSType == "windows" ***REMOVED***
		net = "nat"
	***REMOVED***
	contIP := findContainerIP(c, id, net)
	numPings := 1

	var preRxPackets uint64
	var preTxPackets uint64
	var postRxPackets uint64
	var postTxPackets uint64

	// Get the container networking stats before and after pinging the container
	nwStatsPre := getNetworkStats(c, id)
	for _, v := range nwStatsPre ***REMOVED***
		preRxPackets += v.RxPackets
		preTxPackets += v.TxPackets
	***REMOVED***

	countParam := "-c"
	if runtime.GOOS == "windows" ***REMOVED***
		countParam = "-n" // Ping count parameter is -n on Windows
	***REMOVED***
	pingout, err := exec.Command("ping", contIP, countParam, strconv.Itoa(numPings)).CombinedOutput()
	if err != nil && runtime.GOOS == "linux" ***REMOVED***
		// If it fails then try a work-around, but just for linux.
		// If this fails too then go back to the old error for reporting.
		//
		// The ping will sometimes fail due to an apparmor issue where it
		// denies access to the libc.so.6 shared library - running it
		// via /lib64/ld-linux-x86-64.so.2 seems to work around it.
		pingout2, err2 := exec.Command("/lib64/ld-linux-x86-64.so.2", "/bin/ping", contIP, "-c", strconv.Itoa(numPings)).CombinedOutput()
		if err2 == nil ***REMOVED***
			pingout = pingout2
			err = err2
		***REMOVED***
	***REMOVED***
	c.Assert(err, checker.IsNil)
	pingouts := string(pingout[:])
	nwStatsPost := getNetworkStats(c, id)
	for _, v := range nwStatsPost ***REMOVED***
		postRxPackets += v.RxPackets
		postTxPackets += v.TxPackets
	***REMOVED***

	// Verify the stats contain at least the expected number of packets
	// On Linux, account for ARP.
	expRxPkts := preRxPackets + uint64(numPings)
	expTxPkts := preTxPackets + uint64(numPings)
	if testEnv.OSType != "windows" ***REMOVED***
		expRxPkts++
		expTxPkts++
	***REMOVED***
	c.Assert(postTxPackets, checker.GreaterOrEqualThan, expTxPkts,
		check.Commentf("Reported less TxPackets than expected. Expected >= %d. Found %d. %s", expTxPkts, postTxPackets, pingouts))
	c.Assert(postRxPackets, checker.GreaterOrEqualThan, expRxPkts,
		check.Commentf("Reported less RxPackets than expected. Expected >= %d. Found %d. %s", expRxPkts, postRxPackets, pingouts))
***REMOVED***

func (s *DockerSuite) TestAPIStatsNetworkStatsVersioning(c *check.C) ***REMOVED***
	// Windows doesn't support API versions less than 1.25, so no point testing 1.17 .. 1.21
	testRequires(c, SameHostDaemon, DaemonIsLinux)

	out := runSleepingContainer(c)
	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)
	wg := sync.WaitGroup***REMOVED******REMOVED***

	for i := 17; i <= 21; i++ ***REMOVED***
		wg.Add(1)
		go func(i int) ***REMOVED***
			defer wg.Done()
			apiVersion := fmt.Sprintf("v1.%d", i)
			statsJSONBlob := getVersionedStats(c, id, apiVersion)
			if versions.LessThan(apiVersion, "v1.21") ***REMOVED***
				c.Assert(jsonBlobHasLTv121NetworkStats(statsJSONBlob), checker.Equals, true,
					check.Commentf("Stats JSON blob from API %s %#v does not look like a <v1.21 API stats structure", apiVersion, statsJSONBlob))
			***REMOVED*** else ***REMOVED***
				c.Assert(jsonBlobHasGTE121NetworkStats(statsJSONBlob), checker.Equals, true,
					check.Commentf("Stats JSON blob from API %s %#v does not look like a >=v1.21 API stats structure", apiVersion, statsJSONBlob))
			***REMOVED***
		***REMOVED***(i)
	***REMOVED***
	wg.Wait()
***REMOVED***

func getNetworkStats(c *check.C, id string) map[string]types.NetworkStats ***REMOVED***
	var st *types.StatsJSON

	_, body, err := request.Get(fmt.Sprintf("/containers/%s/stats?stream=false", id))
	c.Assert(err, checker.IsNil)

	err = json.NewDecoder(body).Decode(&st)
	c.Assert(err, checker.IsNil)
	body.Close()

	return st.Networks
***REMOVED***

// getVersionedStats returns stats result for the
// container with id using an API call with version apiVersion. Since the
// stats result type differs between API versions, we simply return
// map[string]interface***REMOVED******REMOVED***.
func getVersionedStats(c *check.C, id string, apiVersion string) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	stats := make(map[string]interface***REMOVED******REMOVED***)

	_, body, err := request.Get(fmt.Sprintf("/%s/containers/%s/stats?stream=false", apiVersion, id))
	c.Assert(err, checker.IsNil)
	defer body.Close()

	err = json.NewDecoder(body).Decode(&stats)
	c.Assert(err, checker.IsNil, check.Commentf("failed to decode stat: %s", err))

	return stats
***REMOVED***

func jsonBlobHasLTv121NetworkStats(blob map[string]interface***REMOVED******REMOVED***) bool ***REMOVED***
	networkStatsIntfc, ok := blob["network"]
	if !ok ***REMOVED***
		return false
	***REMOVED***
	networkStats, ok := networkStatsIntfc.(map[string]interface***REMOVED******REMOVED***)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	for _, expectedKey := range expectedNetworkInterfaceStats ***REMOVED***
		if _, ok := networkStats[expectedKey]; !ok ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func jsonBlobHasGTE121NetworkStats(blob map[string]interface***REMOVED******REMOVED***) bool ***REMOVED***
	networksStatsIntfc, ok := blob["networks"]
	if !ok ***REMOVED***
		return false
	***REMOVED***
	networksStats, ok := networksStatsIntfc.(map[string]interface***REMOVED******REMOVED***)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	for _, networkInterfaceStatsIntfc := range networksStats ***REMOVED***
		networkInterfaceStats, ok := networkInterfaceStatsIntfc.(map[string]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			return false
		***REMOVED***
		for _, expectedKey := range expectedNetworkInterfaceStats ***REMOVED***
			if _, ok := networkInterfaceStats[expectedKey]; !ok ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (s *DockerSuite) TestAPIStatsContainerNotFound(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	expected := "No such container: nonexistent"

	_, err = cli.ContainerStats(context.Background(), "nonexistent", true)
	c.Assert(err.Error(), checker.Contains, expected)
	_, err = cli.ContainerStats(context.Background(), "nonexistent", false)
	c.Assert(err.Error(), checker.Contains, expected)
***REMOVED***

func (s *DockerSuite) TestAPIStatsNoStreamConnectedContainers(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)

	out1 := runSleepingContainer(c)
	id1 := strings.TrimSpace(out1)
	c.Assert(waitRun(id1), checker.IsNil)

	out2 := runSleepingContainer(c, "--net", "container:"+id1)
	id2 := strings.TrimSpace(out2)
	c.Assert(waitRun(id2), checker.IsNil)

	ch := make(chan error, 1)
	go func() ***REMOVED***
		resp, body, err := request.Get(fmt.Sprintf("/containers/%s/stats?stream=false", id2))
		defer body.Close()
		if err != nil ***REMOVED***
			ch <- err
		***REMOVED***
		if resp.StatusCode != http.StatusOK ***REMOVED***
			ch <- fmt.Errorf("Invalid StatusCode %v", resp.StatusCode)
		***REMOVED***
		if resp.Header.Get("Content-Type") != "application/json" ***REMOVED***
			ch <- fmt.Errorf("Invalid 'Content-Type' %v", resp.Header.Get("Content-Type"))
		***REMOVED***
		var v *types.Stats
		if err := json.NewDecoder(body).Decode(&v); err != nil ***REMOVED***
			ch <- err
		***REMOVED***
		ch <- nil
	***REMOVED***()

	select ***REMOVED***
	case err := <-ch:
		c.Assert(err, checker.IsNil, check.Commentf("Error in stats Engine API: %v", err))
	case <-time.After(15 * time.Second):
		c.Fatalf("Stats did not return after timeout")
	***REMOVED***
***REMOVED***
