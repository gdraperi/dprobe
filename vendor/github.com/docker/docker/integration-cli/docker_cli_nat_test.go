package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli"
	"github.com/go-check/check"
)

func startServerContainer(c *check.C, msg string, port int) string ***REMOVED***
	name := "server"
	cmd := []string***REMOVED***
		"run",
		"--name",
		name,
		"-d",
		"-p", fmt.Sprintf("%d:%d", port, port),
		"busybox",
		"sh", "-c", fmt.Sprintf("echo %q | nc -lp %d", msg, port),
	***REMOVED***
	cli.DockerCmd(c, cmd...)
	cli.WaitRun(c, name)
	return name
***REMOVED***

func getExternalAddress(c *check.C) net.IP ***REMOVED***
	iface, err := net.InterfaceByName("eth0")
	if err != nil ***REMOVED***
		c.Skip(fmt.Sprintf("Test not running with `make test`. Interface eth0 not found: %v", err))
	***REMOVED***

	ifaceAddrs, err := iface.Addrs()
	c.Assert(err, check.IsNil)
	c.Assert(ifaceAddrs, checker.Not(checker.HasLen), 0)

	ifaceIP, _, err := net.ParseCIDR(ifaceAddrs[0].String())
	c.Assert(err, check.IsNil)

	return ifaceIP
***REMOVED***

func (s *DockerSuite) TestNetworkNat(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, SameHostDaemon)
	msg := "it works"
	startServerContainer(c, msg, 8080)
	endpoint := getExternalAddress(c)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", endpoint.String(), 8080))
	c.Assert(err, check.IsNil)

	data, err := ioutil.ReadAll(conn)
	conn.Close()
	c.Assert(err, check.IsNil)

	final := strings.TrimRight(string(data), "\n")
	c.Assert(final, checker.Equals, msg)
***REMOVED***

func (s *DockerSuite) TestNetworkLocalhostTCPNat(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, SameHostDaemon)
	var (
		msg = "hi yall"
	)
	startServerContainer(c, msg, 8081)
	conn, err := net.Dial("tcp", "localhost:8081")
	c.Assert(err, check.IsNil)

	data, err := ioutil.ReadAll(conn)
	conn.Close()
	c.Assert(err, check.IsNil)

	final := strings.TrimRight(string(data), "\n")
	c.Assert(final, checker.Equals, msg)
***REMOVED***

func (s *DockerSuite) TestNetworkLoopbackNat(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, SameHostDaemon, NotUserNamespace)
	msg := "it works"
	startServerContainer(c, msg, 8080)
	endpoint := getExternalAddress(c)
	out, _ := dockerCmd(c, "run", "-t", "--net=container:server", "busybox",
		"sh", "-c", fmt.Sprintf("stty raw && nc -w 5 %s 8080", endpoint.String()))
	final := strings.TrimRight(string(out), "\n")
	c.Assert(final, checker.Equals, msg)
***REMOVED***
