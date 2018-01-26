// +build !windows

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"github.com/containerd/containerd/linux"
	"github.com/docker/docker/cmd/dockerd/hack"
	"github.com/docker/docker/daemon"
	"github.com/docker/docker/libcontainerd"
	"github.com/docker/libnetwork/portallocator"
	"golang.org/x/sys/unix"
)

const defaultDaemonConfigFile = "/etc/docker/daemon.json"

// setDefaultUmask sets the umask to 0022 to avoid problems
// caused by custom umask
func setDefaultUmask() error ***REMOVED***
	desiredUmask := 0022
	unix.Umask(desiredUmask)
	if umask := unix.Umask(desiredUmask); umask != desiredUmask ***REMOVED***
		return fmt.Errorf("failed to set umask: expected %#o, got %#o", desiredUmask, umask)
	***REMOVED***

	return nil
***REMOVED***

func getDaemonConfDir(_ string) string ***REMOVED***
	return "/etc/docker"
***REMOVED***

func (cli *DaemonCli) getPlatformRemoteOptions() ([]libcontainerd.RemoteOption, error) ***REMOVED***
	opts := []libcontainerd.RemoteOption***REMOVED***
		libcontainerd.WithOOMScore(cli.Config.OOMScoreAdjust),
		libcontainerd.WithPlugin("linux", &linux.Config***REMOVED***
			Shim:        daemon.DefaultShimBinary,
			Runtime:     daemon.DefaultRuntimeBinary,
			RuntimeRoot: filepath.Join(cli.Config.Root, "runc"),
			ShimDebug:   cli.Config.Debug,
		***REMOVED***),
	***REMOVED***
	if cli.Config.Debug ***REMOVED***
		opts = append(opts, libcontainerd.WithLogLevel("debug"))
	***REMOVED***
	if cli.Config.ContainerdAddr != "" ***REMOVED***
		opts = append(opts, libcontainerd.WithRemoteAddr(cli.Config.ContainerdAddr))
	***REMOVED*** else ***REMOVED***
		opts = append(opts, libcontainerd.WithStartDaemon(true))
	***REMOVED***

	return opts, nil
***REMOVED***

// setupConfigReloadTrap configures the USR2 signal to reload the configuration.
func (cli *DaemonCli) setupConfigReloadTrap() ***REMOVED***
	c := make(chan os.Signal, 1)
	signal.Notify(c, unix.SIGHUP)
	go func() ***REMOVED***
		for range c ***REMOVED***
			cli.reloadConfig()
		***REMOVED***
	***REMOVED***()
***REMOVED***

// getSwarmRunRoot gets the root directory for swarm to store runtime state
// For example, the control socket
func (cli *DaemonCli) getSwarmRunRoot() string ***REMOVED***
	return filepath.Join(cli.Config.ExecRoot, "swarm")
***REMOVED***

// allocateDaemonPort ensures that there are no containers
// that try to use any port allocated for the docker server.
func allocateDaemonPort(addr string) error ***REMOVED***
	host, port, err := net.SplitHostPort(addr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	intPort, err := strconv.Atoi(port)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var hostIPs []net.IP
	if parsedIP := net.ParseIP(host); parsedIP != nil ***REMOVED***
		hostIPs = append(hostIPs, parsedIP)
	***REMOVED*** else if hostIPs, err = net.LookupIP(host); err != nil ***REMOVED***
		return fmt.Errorf("failed to lookup %s address in host specification", host)
	***REMOVED***

	pa := portallocator.Get()
	for _, hostIP := range hostIPs ***REMOVED***
		if _, err := pa.RequestPort(hostIP, "tcp", intPort); err != nil ***REMOVED***
			return fmt.Errorf("failed to allocate daemon listening port %d (err: %v)", intPort, err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// notifyShutdown is called after the daemon shuts down but before the process exits.
func notifyShutdown(err error) ***REMOVED***
***REMOVED***

func wrapListeners(proto string, ls []net.Listener) []net.Listener ***REMOVED***
	switch proto ***REMOVED***
	case "unix":
		ls[0] = &hack.MalformedHostHeaderOverride***REMOVED***ls[0]***REMOVED***
	case "fd":
		for i := range ls ***REMOVED***
			ls[i] = &hack.MalformedHostHeaderOverride***REMOVED***ls[i]***REMOVED***
		***REMOVED***
	***REMOVED***
	return ls
***REMOVED***
