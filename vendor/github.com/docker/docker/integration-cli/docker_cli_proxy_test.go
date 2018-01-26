package main

import (
	"net"
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

func (s *DockerSuite) TestCLIProxyDisableProxyUnixSock(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, SameHostDaemon)

	icmd.RunCmd(icmd.Cmd***REMOVED***
		Command: []string***REMOVED***dockerBinary, "info"***REMOVED***,
		Env:     appendBaseEnv(false, "HTTP_PROXY=http://127.0.0.1:9999"),
	***REMOVED***).Assert(c, icmd.Success)
***REMOVED***

// Can't use localhost here since go has a special case to not use proxy if connecting to localhost
// See https://golang.org/pkg/net/http/#ProxyFromEnvironment
func (s *DockerDaemonSuite) TestCLIProxyProxyTCPSock(c *check.C) ***REMOVED***
	testRequires(c, SameHostDaemon)
	// get the IP to use to connect since we can't use localhost
	addrs, err := net.InterfaceAddrs()
	c.Assert(err, checker.IsNil)
	var ip string
	for _, addr := range addrs ***REMOVED***
		sAddr := addr.String()
		if !strings.Contains(sAddr, "127.0.0.1") ***REMOVED***
			addrArr := strings.Split(sAddr, "/")
			ip = addrArr[0]
			break
		***REMOVED***
	***REMOVED***

	c.Assert(ip, checker.Not(checker.Equals), "")

	s.d.Start(c, "-H", "tcp://"+ip+":2375")

	icmd.RunCmd(icmd.Cmd***REMOVED***
		Command: []string***REMOVED***dockerBinary, "info"***REMOVED***,
		Env:     []string***REMOVED***"DOCKER_HOST=tcp://" + ip + ":2375", "HTTP_PROXY=127.0.0.1:9999"***REMOVED***,
	***REMOVED***).Assert(c, icmd.Expected***REMOVED***Error: "exit status 1", ExitCode: 1***REMOVED***)
	// Test with no_proxy
	icmd.RunCmd(icmd.Cmd***REMOVED***
		Command: []string***REMOVED***dockerBinary, "info"***REMOVED***,
		Env:     []string***REMOVED***"DOCKER_HOST=tcp://" + ip + ":2375", "HTTP_PROXY=127.0.0.1:9999", "NO_PROXY=" + ip***REMOVED***,
	***REMOVED***).Assert(c, icmd.Success)
***REMOVED***
