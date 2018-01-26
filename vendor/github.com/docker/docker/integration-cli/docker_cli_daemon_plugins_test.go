// +build linux

package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/pkg/mount"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
	"golang.org/x/sys/unix"
)

// TestDaemonRestartWithPluginEnabled tests state restore for an enabled plugin
func (s *DockerDaemonSuite) TestDaemonRestartWithPluginEnabled(c *check.C) ***REMOVED***
	testRequires(c, IsAmd64, Network)

	s.d.Start(c)

	if out, err := s.d.Cmd("plugin", "install", "--grant-all-permissions", pName); err != nil ***REMOVED***
		c.Fatalf("Could not install plugin: %v %s", err, out)
	***REMOVED***

	defer func() ***REMOVED***
		if out, err := s.d.Cmd("plugin", "disable", pName); err != nil ***REMOVED***
			c.Fatalf("Could not disable plugin: %v %s", err, out)
		***REMOVED***
		if out, err := s.d.Cmd("plugin", "remove", pName); err != nil ***REMOVED***
			c.Fatalf("Could not remove plugin: %v %s", err, out)
		***REMOVED***
	***REMOVED***()

	s.d.Restart(c)

	out, err := s.d.Cmd("plugin", "ls")
	if err != nil ***REMOVED***
		c.Fatalf("Could not list plugins: %v %s", err, out)
	***REMOVED***
	c.Assert(out, checker.Contains, pName)
	c.Assert(out, checker.Contains, "true")
***REMOVED***

// TestDaemonRestartWithPluginDisabled tests state restore for a disabled plugin
func (s *DockerDaemonSuite) TestDaemonRestartWithPluginDisabled(c *check.C) ***REMOVED***
	testRequires(c, IsAmd64, Network)

	s.d.Start(c)

	if out, err := s.d.Cmd("plugin", "install", "--grant-all-permissions", pName, "--disable"); err != nil ***REMOVED***
		c.Fatalf("Could not install plugin: %v %s", err, out)
	***REMOVED***

	defer func() ***REMOVED***
		if out, err := s.d.Cmd("plugin", "remove", pName); err != nil ***REMOVED***
			c.Fatalf("Could not remove plugin: %v %s", err, out)
		***REMOVED***
	***REMOVED***()

	s.d.Restart(c)

	out, err := s.d.Cmd("plugin", "ls")
	if err != nil ***REMOVED***
		c.Fatalf("Could not list plugins: %v %s", err, out)
	***REMOVED***
	c.Assert(out, checker.Contains, pName)
	c.Assert(out, checker.Contains, "false")
***REMOVED***

// TestDaemonKillLiveRestoreWithPlugins SIGKILLs daemon started with --live-restore.
// Plugins should continue to run.
func (s *DockerDaemonSuite) TestDaemonKillLiveRestoreWithPlugins(c *check.C) ***REMOVED***
	testRequires(c, IsAmd64, Network)

	s.d.Start(c, "--live-restore")
	if out, err := s.d.Cmd("plugin", "install", "--grant-all-permissions", pName); err != nil ***REMOVED***
		c.Fatalf("Could not install plugin: %v %s", err, out)
	***REMOVED***
	defer func() ***REMOVED***
		s.d.Restart(c, "--live-restore")
		if out, err := s.d.Cmd("plugin", "disable", pName); err != nil ***REMOVED***
			c.Fatalf("Could not disable plugin: %v %s", err, out)
		***REMOVED***
		if out, err := s.d.Cmd("plugin", "remove", pName); err != nil ***REMOVED***
			c.Fatalf("Could not remove plugin: %v %s", err, out)
		***REMOVED***
	***REMOVED***()

	if err := s.d.Kill(); err != nil ***REMOVED***
		c.Fatalf("Could not kill daemon: %v", err)
	***REMOVED***

	icmd.RunCommand("pgrep", "-f", pluginProcessName).Assert(c, icmd.Success)
***REMOVED***

// TestDaemonShutdownLiveRestoreWithPlugins SIGTERMs daemon started with --live-restore.
// Plugins should continue to run.
func (s *DockerDaemonSuite) TestDaemonShutdownLiveRestoreWithPlugins(c *check.C) ***REMOVED***
	testRequires(c, IsAmd64, Network)

	s.d.Start(c, "--live-restore")
	if out, err := s.d.Cmd("plugin", "install", "--grant-all-permissions", pName); err != nil ***REMOVED***
		c.Fatalf("Could not install plugin: %v %s", err, out)
	***REMOVED***
	defer func() ***REMOVED***
		s.d.Restart(c, "--live-restore")
		if out, err := s.d.Cmd("plugin", "disable", pName); err != nil ***REMOVED***
			c.Fatalf("Could not disable plugin: %v %s", err, out)
		***REMOVED***
		if out, err := s.d.Cmd("plugin", "remove", pName); err != nil ***REMOVED***
			c.Fatalf("Could not remove plugin: %v %s", err, out)
		***REMOVED***
	***REMOVED***()

	if err := s.d.Interrupt(); err != nil ***REMOVED***
		c.Fatalf("Could not kill daemon: %v", err)
	***REMOVED***

	icmd.RunCommand("pgrep", "-f", pluginProcessName).Assert(c, icmd.Success)
***REMOVED***

// TestDaemonShutdownWithPlugins shuts down running plugins.
func (s *DockerDaemonSuite) TestDaemonShutdownWithPlugins(c *check.C) ***REMOVED***
	testRequires(c, IsAmd64, Network, SameHostDaemon)

	s.d.Start(c)
	if out, err := s.d.Cmd("plugin", "install", "--grant-all-permissions", pName); err != nil ***REMOVED***
		c.Fatalf("Could not install plugin: %v %s", err, out)
	***REMOVED***

	defer func() ***REMOVED***
		s.d.Restart(c)
		if out, err := s.d.Cmd("plugin", "disable", pName); err != nil ***REMOVED***
			c.Fatalf("Could not disable plugin: %v %s", err, out)
		***REMOVED***
		if out, err := s.d.Cmd("plugin", "remove", pName); err != nil ***REMOVED***
			c.Fatalf("Could not remove plugin: %v %s", err, out)
		***REMOVED***
	***REMOVED***()

	if err := s.d.Interrupt(); err != nil ***REMOVED***
		c.Fatalf("Could not kill daemon: %v", err)
	***REMOVED***

	for ***REMOVED***
		if err := unix.Kill(s.d.Pid(), 0); err == unix.ESRCH ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	icmd.RunCommand("pgrep", "-f", pluginProcessName).Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
		Error:    "exit status 1",
	***REMOVED***)

	s.d.Start(c)
	icmd.RunCommand("pgrep", "-f", pluginProcessName).Assert(c, icmd.Success)
***REMOVED***

// TestDaemonKillWithPlugins leaves plugins running.
func (s *DockerDaemonSuite) TestDaemonKillWithPlugins(c *check.C) ***REMOVED***
	testRequires(c, IsAmd64, Network, SameHostDaemon)

	s.d.Start(c)
	if out, err := s.d.Cmd("plugin", "install", "--grant-all-permissions", pName); err != nil ***REMOVED***
		c.Fatalf("Could not install plugin: %v %s", err, out)
	***REMOVED***

	defer func() ***REMOVED***
		s.d.Restart(c)
		if out, err := s.d.Cmd("plugin", "disable", pName); err != nil ***REMOVED***
			c.Fatalf("Could not disable plugin: %v %s", err, out)
		***REMOVED***
		if out, err := s.d.Cmd("plugin", "remove", pName); err != nil ***REMOVED***
			c.Fatalf("Could not remove plugin: %v %s", err, out)
		***REMOVED***
	***REMOVED***()

	if err := s.d.Kill(); err != nil ***REMOVED***
		c.Fatalf("Could not kill daemon: %v", err)
	***REMOVED***

	// assert that plugins are running.
	icmd.RunCommand("pgrep", "-f", pluginProcessName).Assert(c, icmd.Success)
***REMOVED***

// TestVolumePlugin tests volume creation using a plugin.
func (s *DockerDaemonSuite) TestVolumePlugin(c *check.C) ***REMOVED***
	testRequires(c, IsAmd64, Network)

	volName := "plugin-volume"
	destDir := "/tmp/data/"
	destFile := "foo"

	s.d.Start(c)
	out, err := s.d.Cmd("plugin", "install", pName, "--grant-all-permissions")
	if err != nil ***REMOVED***
		c.Fatalf("Could not install plugin: %v %s", err, out)
	***REMOVED***
	pluginID, err := s.d.Cmd("plugin", "inspect", "-f", "***REMOVED******REMOVED***.Id***REMOVED******REMOVED***", pName)
	pluginID = strings.TrimSpace(pluginID)
	if err != nil ***REMOVED***
		c.Fatalf("Could not retrieve plugin ID: %v %s", err, pluginID)
	***REMOVED***
	mountpointPrefix := filepath.Join(s.d.RootDir(), "plugins", pluginID, "rootfs")
	defer func() ***REMOVED***
		if out, err := s.d.Cmd("plugin", "disable", pName); err != nil ***REMOVED***
			c.Fatalf("Could not disable plugin: %v %s", err, out)
		***REMOVED***

		if out, err := s.d.Cmd("plugin", "remove", pName); err != nil ***REMOVED***
			c.Fatalf("Could not remove plugin: %v %s", err, out)
		***REMOVED***

		exists, err := existsMountpointWithPrefix(mountpointPrefix)
		c.Assert(err, checker.IsNil)
		c.Assert(exists, checker.Equals, false)

	***REMOVED***()

	out, err = s.d.Cmd("volume", "create", "-d", pName, volName)
	if err != nil ***REMOVED***
		c.Fatalf("Could not create volume: %v %s", err, out)
	***REMOVED***
	defer func() ***REMOVED***
		if out, err := s.d.Cmd("volume", "remove", volName); err != nil ***REMOVED***
			c.Fatalf("Could not remove volume: %v %s", err, out)
		***REMOVED***
	***REMOVED***()

	out, err = s.d.Cmd("volume", "ls")
	if err != nil ***REMOVED***
		c.Fatalf("Could not list volume: %v %s", err, out)
	***REMOVED***
	c.Assert(out, checker.Contains, volName)
	c.Assert(out, checker.Contains, pName)

	mountPoint, err := s.d.Cmd("volume", "inspect", volName, "--format", "***REMOVED******REMOVED***.Mountpoint***REMOVED******REMOVED***")
	if err != nil ***REMOVED***
		c.Fatalf("Could not inspect volume: %v %s", err, mountPoint)
	***REMOVED***
	mountPoint = strings.TrimSpace(mountPoint)

	out, err = s.d.Cmd("run", "--rm", "-v", volName+":"+destDir, "busybox", "touch", destDir+destFile)
	c.Assert(err, checker.IsNil, check.Commentf(out))
	path := filepath.Join(s.d.RootDir(), "plugins", pluginID, "rootfs", mountPoint, destFile)
	_, err = os.Lstat(path)
	c.Assert(err, checker.IsNil)

	exists, err := existsMountpointWithPrefix(mountpointPrefix)
	c.Assert(err, checker.IsNil)
	c.Assert(exists, checker.Equals, true)
***REMOVED***

func (s *DockerDaemonSuite) TestGraphdriverPlugin(c *check.C) ***REMOVED***
	testRequires(c, Network, IsAmd64, DaemonIsLinux, overlay2Supported, ExperimentalDaemon)

	s.d.Start(c)

	// install the plugin
	plugin := "cpuguy83/docker-overlay2-graphdriver-plugin"
	out, err := s.d.Cmd("plugin", "install", "--grant-all-permissions", plugin)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	// restart the daemon with the plugin set as the storage driver
	s.d.Restart(c, "-s", plugin, "--storage-opt", "overlay2.override_kernel_check=1")

	// run a container
	out, err = s.d.Cmd("run", "--rm", "busybox", "true") // this will pull busybox using the plugin
	c.Assert(err, checker.IsNil, check.Commentf(out))
***REMOVED***

func (s *DockerDaemonSuite) TestPluginVolumeRemoveOnRestart(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux, Network, IsAmd64)

	s.d.Start(c, "--live-restore=true")

	out, err := s.d.Cmd("plugin", "install", "--grant-all-permissions", pName)
	c.Assert(err, checker.IsNil, check.Commentf(out))
	c.Assert(strings.TrimSpace(out), checker.Contains, pName)

	out, err = s.d.Cmd("volume", "create", "--driver", pName, "test")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	s.d.Restart(c, "--live-restore=true")

	out, err = s.d.Cmd("plugin", "disable", pName)
	c.Assert(err, checker.NotNil, check.Commentf(out))
	c.Assert(out, checker.Contains, "in use")

	out, err = s.d.Cmd("volume", "rm", "test")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	out, err = s.d.Cmd("plugin", "disable", pName)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	out, err = s.d.Cmd("plugin", "rm", pName)
	c.Assert(err, checker.IsNil, check.Commentf(out))
***REMOVED***

func existsMountpointWithPrefix(mountpointPrefix string) (bool, error) ***REMOVED***
	mounts, err := mount.GetMounts()
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	for _, mnt := range mounts ***REMOVED***
		if strings.HasPrefix(mnt.Mountpoint, mountpointPrefix) ***REMOVED***
			return true, nil
		***REMOVED***
	***REMOVED***
	return false, nil
***REMOVED***

func (s *DockerDaemonSuite) TestPluginListFilterEnabled(c *check.C) ***REMOVED***
	testRequires(c, IsAmd64, Network)

	s.d.Start(c)

	out, err := s.d.Cmd("plugin", "install", "--grant-all-permissions", pNameWithTag, "--disable")
	c.Assert(err, check.IsNil, check.Commentf(out))

	defer func() ***REMOVED***
		if out, err := s.d.Cmd("plugin", "remove", pNameWithTag); err != nil ***REMOVED***
			c.Fatalf("Could not remove plugin: %v %s", err, out)
		***REMOVED***
	***REMOVED***()

	out, err = s.d.Cmd("plugin", "ls", "--filter", "enabled=true")
	c.Assert(err, checker.IsNil)
	c.Assert(out, checker.Not(checker.Contains), pName)

	out, err = s.d.Cmd("plugin", "ls", "--filter", "enabled=false")
	c.Assert(err, checker.IsNil)
	c.Assert(out, checker.Contains, pName)
	c.Assert(out, checker.Contains, "false")

	out, err = s.d.Cmd("plugin", "ls")
	c.Assert(err, checker.IsNil)
	c.Assert(out, checker.Contains, pName)
***REMOVED***

func (s *DockerDaemonSuite) TestPluginListFilterCapability(c *check.C) ***REMOVED***
	testRequires(c, IsAmd64, Network)

	s.d.Start(c)

	out, err := s.d.Cmd("plugin", "install", "--grant-all-permissions", pNameWithTag, "--disable")
	c.Assert(err, check.IsNil, check.Commentf(out))

	defer func() ***REMOVED***
		if out, err := s.d.Cmd("plugin", "remove", pNameWithTag); err != nil ***REMOVED***
			c.Fatalf("Could not remove plugin: %v %s", err, out)
		***REMOVED***
	***REMOVED***()

	out, err = s.d.Cmd("plugin", "ls", "--filter", "capability=volumedriver")
	c.Assert(err, checker.IsNil)
	c.Assert(out, checker.Contains, pName)

	out, err = s.d.Cmd("plugin", "ls", "--filter", "capability=authz")
	c.Assert(err, checker.IsNil)
	c.Assert(out, checker.Not(checker.Contains), pName)

	out, err = s.d.Cmd("plugin", "ls")
	c.Assert(err, checker.IsNil)
	c.Assert(out, checker.Contains, pName)
***REMOVED***
