// +build linux freebsd

package daemon

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/links"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/runconfig"
	"github.com/docker/libnetwork"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func (daemon *Daemon) setupLinkedContainers(container *container.Container) ([]string, error) ***REMOVED***
	var env []string
	children := daemon.children(container)

	bridgeSettings := container.NetworkSettings.Networks[runconfig.DefaultDaemonNetworkMode().NetworkName()]
	if bridgeSettings == nil || bridgeSettings.EndpointSettings == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	for linkAlias, child := range children ***REMOVED***
		if !child.IsRunning() ***REMOVED***
			return nil, fmt.Errorf("Cannot link to a non running container: %s AS %s", child.Name, linkAlias)
		***REMOVED***

		childBridgeSettings := child.NetworkSettings.Networks[runconfig.DefaultDaemonNetworkMode().NetworkName()]
		if childBridgeSettings == nil || childBridgeSettings.EndpointSettings == nil ***REMOVED***
			return nil, fmt.Errorf("container %s not attached to default bridge network", child.ID)
		***REMOVED***

		link := links.NewLink(
			bridgeSettings.IPAddress,
			childBridgeSettings.IPAddress,
			linkAlias,
			child.Config.Env,
			child.Config.ExposedPorts,
		)

		env = append(env, link.ToEnv()...)
	***REMOVED***

	return env, nil
***REMOVED***

func (daemon *Daemon) getIpcContainer(id string) (*container.Container, error) ***REMOVED***
	errMsg := "can't join IPC of container " + id
	// Check the container exists
	container, err := daemon.GetContainer(id)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, errMsg)
	***REMOVED***
	// Check the container is running and not restarting
	if err := daemon.checkContainer(container, containerIsRunning, containerIsNotRestarting); err != nil ***REMOVED***
		return nil, errors.Wrap(err, errMsg)
	***REMOVED***
	// Check the container ipc is shareable
	if st, err := os.Stat(container.ShmPath); err != nil || !st.IsDir() ***REMOVED***
		if err == nil || os.IsNotExist(err) ***REMOVED***
			return nil, errors.New(errMsg + ": non-shareable IPC")
		***REMOVED***
		// stat() failed?
		return nil, errors.Wrap(err, errMsg+": unexpected error from stat "+container.ShmPath)
	***REMOVED***

	return container, nil
***REMOVED***

func (daemon *Daemon) getPidContainer(container *container.Container) (*container.Container, error) ***REMOVED***
	containerID := container.HostConfig.PidMode.Container()
	container, err := daemon.GetContainer(containerID)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "cannot join PID of a non running container: %s", containerID)
	***REMOVED***
	return container, daemon.checkContainer(container, containerIsRunning, containerIsNotRestarting)
***REMOVED***

func containerIsRunning(c *container.Container) error ***REMOVED***
	if !c.IsRunning() ***REMOVED***
		return errdefs.Conflict(errors.Errorf("container %s is not running", c.ID))
	***REMOVED***
	return nil
***REMOVED***

func containerIsNotRestarting(c *container.Container) error ***REMOVED***
	if c.IsRestarting() ***REMOVED***
		return errContainerIsRestarting(c.ID)
	***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) setupIpcDirs(c *container.Container) error ***REMOVED***
	ipcMode := c.HostConfig.IpcMode

	switch ***REMOVED***
	case ipcMode.IsContainer():
		ic, err := daemon.getIpcContainer(ipcMode.Container())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		c.ShmPath = ic.ShmPath

	case ipcMode.IsHost():
		if _, err := os.Stat("/dev/shm"); err != nil ***REMOVED***
			return fmt.Errorf("/dev/shm is not mounted, but must be for --ipc=host")
		***REMOVED***
		c.ShmPath = "/dev/shm"

	case ipcMode.IsPrivate(), ipcMode.IsNone():
		// c.ShmPath will/should not be used, so make it empty.
		// Container's /dev/shm mount comes from OCI spec.
		c.ShmPath = ""

	case ipcMode.IsEmpty():
		// A container was created by an older version of the daemon.
		// The default behavior used to be what is now called "shareable".
		fallthrough

	case ipcMode.IsShareable():
		rootIDs := daemon.idMappings.RootPair()
		if !c.HasMountFor("/dev/shm") ***REMOVED***
			shmPath, err := c.ShmResourcePath()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if err := idtools.MkdirAllAndChown(shmPath, 0700, rootIDs); err != nil ***REMOVED***
				return err
			***REMOVED***

			shmproperty := "mode=1777,size=" + strconv.FormatInt(c.HostConfig.ShmSize, 10)
			if err := unix.Mount("shm", shmPath, "tmpfs", uintptr(unix.MS_NOEXEC|unix.MS_NOSUID|unix.MS_NODEV), label.FormatMountLabel(shmproperty, c.GetMountLabel())); err != nil ***REMOVED***
				return fmt.Errorf("mounting shm tmpfs: %s", err)
			***REMOVED***
			if err := os.Chown(shmPath, rootIDs.UID, rootIDs.GID); err != nil ***REMOVED***
				return err
			***REMOVED***
			c.ShmPath = shmPath
		***REMOVED***

	default:
		return fmt.Errorf("invalid IPC mode: %v", ipcMode)
	***REMOVED***

	return nil
***REMOVED***

func (daemon *Daemon) setupSecretDir(c *container.Container) (setupErr error) ***REMOVED***
	if len(c.SecretReferences) == 0 ***REMOVED***
		return nil
	***REMOVED***

	localMountPath, err := c.SecretMountPath()
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error getting secrets mount dir")
	***REMOVED***
	logrus.Debugf("secrets: setting up secret dir: %s", localMountPath)

	// retrieve possible remapped range start for root UID, GID
	rootIDs := daemon.idMappings.RootPair()
	// create tmpfs
	if err := idtools.MkdirAllAndChown(localMountPath, 0700, rootIDs); err != nil ***REMOVED***
		return errors.Wrap(err, "error creating secret local mount path")
	***REMOVED***

	defer func() ***REMOVED***
		if setupErr != nil ***REMOVED***
			// cleanup
			_ = detachMounted(localMountPath)

			if err := os.RemoveAll(localMountPath); err != nil ***REMOVED***
				logrus.Errorf("error cleaning up secret mount: %s", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	tmpfsOwnership := fmt.Sprintf("uid=%d,gid=%d", rootIDs.UID, rootIDs.GID)
	if err := mount.Mount("tmpfs", localMountPath, "tmpfs", "nodev,nosuid,noexec,"+tmpfsOwnership); err != nil ***REMOVED***
		return errors.Wrap(err, "unable to setup secret mount")
	***REMOVED***

	if c.DependencyStore == nil ***REMOVED***
		return fmt.Errorf("secret store is not initialized")
	***REMOVED***

	for _, s := range c.SecretReferences ***REMOVED***
		// TODO (ehazlett): use type switch when more are supported
		if s.File == nil ***REMOVED***
			logrus.Error("secret target type is not a file target")
			continue
		***REMOVED***

		// secrets are created in the SecretMountPath on the host, at a
		// single level
		fPath, err := c.SecretFilePath(*s)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "error getting secret file path")
		***REMOVED***
		if err := idtools.MkdirAllAndChown(filepath.Dir(fPath), 0700, rootIDs); err != nil ***REMOVED***
			return errors.Wrap(err, "error creating secret mount path")
		***REMOVED***

		logrus.WithFields(logrus.Fields***REMOVED***
			"name": s.File.Name,
			"path": fPath,
		***REMOVED***).Debug("injecting secret")
		secret, err := c.DependencyStore.Secrets().Get(s.SecretID)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "unable to get secret from secret store")
		***REMOVED***
		if err := ioutil.WriteFile(fPath, secret.Spec.Data, s.File.Mode); err != nil ***REMOVED***
			return errors.Wrap(err, "error injecting secret")
		***REMOVED***

		uid, err := strconv.Atoi(s.File.UID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		gid, err := strconv.Atoi(s.File.GID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := os.Chown(fPath, rootIDs.UID+uid, rootIDs.GID+gid); err != nil ***REMOVED***
			return errors.Wrap(err, "error setting ownership for secret")
		***REMOVED***
	***REMOVED***

	label.Relabel(localMountPath, c.MountLabel, false)

	// remount secrets ro
	if err := mount.Mount("tmpfs", localMountPath, "tmpfs", "remount,ro,"+tmpfsOwnership); err != nil ***REMOVED***
		return errors.Wrap(err, "unable to remount secret dir as readonly")
	***REMOVED***

	return nil
***REMOVED***

func (daemon *Daemon) setupConfigDir(c *container.Container) (setupErr error) ***REMOVED***
	if len(c.ConfigReferences) == 0 ***REMOVED***
		return nil
	***REMOVED***

	localPath, err := c.ConfigsDirPath()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	logrus.Debugf("configs: setting up config dir: %s", localPath)

	// retrieve possible remapped range start for root UID, GID
	rootIDs := daemon.idMappings.RootPair()
	// create tmpfs
	if err := idtools.MkdirAllAndChown(localPath, 0700, rootIDs); err != nil ***REMOVED***
		return errors.Wrap(err, "error creating config dir")
	***REMOVED***

	defer func() ***REMOVED***
		if setupErr != nil ***REMOVED***
			if err := os.RemoveAll(localPath); err != nil ***REMOVED***
				logrus.Errorf("error cleaning up config dir: %s", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if c.DependencyStore == nil ***REMOVED***
		return fmt.Errorf("config store is not initialized")
	***REMOVED***

	for _, configRef := range c.ConfigReferences ***REMOVED***
		// TODO (ehazlett): use type switch when more are supported
		if configRef.File == nil ***REMOVED***
			logrus.Error("config target type is not a file target")
			continue
		***REMOVED***

		fPath, err := c.ConfigFilePath(*configRef)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		log := logrus.WithFields(logrus.Fields***REMOVED***"name": configRef.File.Name, "path": fPath***REMOVED***)

		if err := idtools.MkdirAllAndChown(filepath.Dir(fPath), 0700, rootIDs); err != nil ***REMOVED***
			return errors.Wrap(err, "error creating config path")
		***REMOVED***

		log.Debug("injecting config")
		config, err := c.DependencyStore.Configs().Get(configRef.ConfigID)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "unable to get config from config store")
		***REMOVED***
		if err := ioutil.WriteFile(fPath, config.Spec.Data, configRef.File.Mode); err != nil ***REMOVED***
			return errors.Wrap(err, "error injecting config")
		***REMOVED***

		uid, err := strconv.Atoi(configRef.File.UID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		gid, err := strconv.Atoi(configRef.File.GID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := os.Chown(fPath, rootIDs.UID+uid, rootIDs.GID+gid); err != nil ***REMOVED***
			return errors.Wrap(err, "error setting ownership for config")
		***REMOVED***

		label.Relabel(fPath, c.MountLabel, false)
	***REMOVED***

	return nil
***REMOVED***

func killProcessDirectly(cntr *container.Container) error ***REMOVED***
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Block until the container to stops or timeout.
	status := <-cntr.Wait(ctx, container.WaitConditionNotRunning)
	if status.Err() != nil ***REMOVED***
		// Ensure that we don't kill ourselves
		if pid := cntr.GetPID(); pid != 0 ***REMOVED***
			logrus.Infof("Container %s failed to exit within 10 seconds of kill - trying direct SIGKILL", stringid.TruncateID(cntr.ID))
			if err := unix.Kill(pid, 9); err != nil ***REMOVED***
				if err != unix.ESRCH ***REMOVED***
					return err
				***REMOVED***
				e := errNoSuchProcess***REMOVED***pid, 9***REMOVED***
				logrus.Debug(e)
				return e
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func detachMounted(path string) error ***REMOVED***
	return unix.Unmount(path, unix.MNT_DETACH)
***REMOVED***

func isLinkable(child *container.Container) bool ***REMOVED***
	// A container is linkable only if it belongs to the default network
	_, ok := child.NetworkSettings.Networks[runconfig.DefaultDaemonNetworkMode().NetworkName()]
	return ok
***REMOVED***

func enableIPOnPredefinedNetwork() bool ***REMOVED***
	return false
***REMOVED***

func (daemon *Daemon) isNetworkHotPluggable() bool ***REMOVED***
	return true
***REMOVED***

func setupPathsAndSandboxOptions(container *container.Container, sboxOptions *[]libnetwork.SandboxOption) error ***REMOVED***
	var err error

	container.HostsPath, err = container.GetRootResourcePath("hosts")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*sboxOptions = append(*sboxOptions, libnetwork.OptionHostsPath(container.HostsPath))

	container.ResolvConfPath, err = container.GetRootResourcePath("resolv.conf")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*sboxOptions = append(*sboxOptions, libnetwork.OptionResolvConfPath(container.ResolvConfPath))
	return nil
***REMOVED***

func (daemon *Daemon) initializeNetworkingPaths(container *container.Container, nc *container.Container) error ***REMOVED***
	container.HostnamePath = nc.HostnamePath
	container.HostsPath = nc.HostsPath
	container.ResolvConfPath = nc.ResolvConfPath
	return nil
***REMOVED***

func (daemon *Daemon) setupContainerMountsRoot(c *container.Container) error ***REMOVED***
	// get the root mount path so we can make it unbindable
	p, err := c.MountsResourcePath("")
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := idtools.MkdirAllAndChown(p, 0700, daemon.idMappings.RootPair()); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := mount.MakeUnbindable(p); err != nil ***REMOVED***
		// Setting unbindable is a precaution and is not neccessary for correct operation.
		// Do not error out if this fails.
		logrus.WithError(err).WithField("resource", p).WithField("container", c.ID).Warn("Error setting container resource mounts to unbindable, this may cause mount leakages, preventing removal of this container.")
	***REMOVED***
	return nil
***REMOVED***
