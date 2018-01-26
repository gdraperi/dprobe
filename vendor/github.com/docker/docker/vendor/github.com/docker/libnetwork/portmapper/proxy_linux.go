package portmapper

import (
	"net"
	"os/exec"
	"strconv"
	"syscall"
)

func newProxyCommand(proto string, hostIP net.IP, hostPort int, containerIP net.IP, containerPort int, proxyPath string) (userlandProxy, error) ***REMOVED***
	path := proxyPath
	if proxyPath == "" ***REMOVED***
		cmd, err := exec.LookPath(userlandProxyCommandName)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		path = cmd
	***REMOVED***

	args := []string***REMOVED***
		path,
		"-proto", proto,
		"-host-ip", hostIP.String(),
		"-host-port", strconv.Itoa(hostPort),
		"-container-ip", containerIP.String(),
		"-container-port", strconv.Itoa(containerPort),
	***REMOVED***

	return &proxyCommand***REMOVED***
		cmd: &exec.Cmd***REMOVED***
			Path: path,
			Args: args,
			SysProcAttr: &syscall.SysProcAttr***REMOVED***
				Pdeathsig: syscall.SIGTERM, // send a sigterm to the proxy if the daemon process dies
			***REMOVED***,
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***
