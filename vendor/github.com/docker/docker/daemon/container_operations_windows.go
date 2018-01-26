package daemon

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/libnetwork"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (daemon *Daemon) setupLinkedContainers(container *container.Container) ([]string, error) ***REMOVED***
	return nil, nil
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

	// create local config root
	if err := system.MkdirAllWithACL(localPath, 0, system.SddlAdministratorsLocalSystem); err != nil ***REMOVED***
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

		log.Debug("injecting config")
		config, err := c.DependencyStore.Configs().Get(configRef.ConfigID)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "unable to get config from config store")
		***REMOVED***
		if err := ioutil.WriteFile(fPath, config.Spec.Data, configRef.File.Mode); err != nil ***REMOVED***
			return errors.Wrap(err, "error injecting config")
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// getSize returns real size & virtual size
func (daemon *Daemon) getSize(containerID string) (int64, int64) ***REMOVED***
	// TODO Windows
	return 0, 0
***REMOVED***

func (daemon *Daemon) setupIpcDirs(container *container.Container) error ***REMOVED***
	return nil
***REMOVED***

// TODO Windows: Fix Post-TP5. This is a hack to allow docker cp to work
// against containers which have volumes. You will still be able to cp
// to somewhere on the container drive, but not to any mounted volumes
// inside the container. Without this fix, docker cp is broken to any
// container which has a volume, regardless of where the file is inside the
// container.
func (daemon *Daemon) mountVolumes(container *container.Container) error ***REMOVED***
	return nil
***REMOVED***

func detachMounted(path string) error ***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) setupSecretDir(c *container.Container) (setupErr error) ***REMOVED***
	if len(c.SecretReferences) == 0 ***REMOVED***
		return nil
	***REMOVED***

	localMountPath, err := c.SecretMountPath()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	logrus.Debugf("secrets: setting up secret dir: %s", localMountPath)

	// create local secret root
	if err := system.MkdirAllWithACL(localMountPath, 0, system.SddlAdministratorsLocalSystem); err != nil ***REMOVED***
		return errors.Wrap(err, "error creating secret local directory")
	***REMOVED***

	defer func() ***REMOVED***
		if setupErr != nil ***REMOVED***
			if err := os.RemoveAll(localMountPath); err != nil ***REMOVED***
				logrus.Errorf("error cleaning up secret mount: %s", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

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
			return err
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
	***REMOVED***

	return nil
***REMOVED***

func killProcessDirectly(container *container.Container) error ***REMOVED***
	return nil
***REMOVED***

func isLinkable(child *container.Container) bool ***REMOVED***
	return false
***REMOVED***

func enableIPOnPredefinedNetwork() bool ***REMOVED***
	return true
***REMOVED***

func (daemon *Daemon) isNetworkHotPluggable() bool ***REMOVED***
	return false
***REMOVED***

func setupPathsAndSandboxOptions(container *container.Container, sboxOptions *[]libnetwork.SandboxOption) error ***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) initializeNetworkingPaths(container *container.Container, nc *container.Container) error ***REMOVED***

	if nc.HostConfig.Isolation.IsHyperV() ***REMOVED***
		return fmt.Errorf("sharing of hyperv containers network is not supported")
	***REMOVED***

	container.NetworkSharedContainerID = nc.ID

	if nc.NetworkSettings != nil ***REMOVED***
		for n := range nc.NetworkSettings.Networks ***REMOVED***
			sn, err := daemon.FindNetwork(n)
			if err != nil ***REMOVED***
				continue
			***REMOVED***

			ep, err := nc.GetEndpointInNetwork(sn)
			if err != nil ***REMOVED***
				continue
			***REMOVED***

			data, err := ep.DriverInfo()
			if err != nil ***REMOVED***
				continue
			***REMOVED***

			if data["GW_INFO"] != nil ***REMOVED***
				gwInfo := data["GW_INFO"].(map[string]interface***REMOVED******REMOVED***)
				if gwInfo["hnsid"] != nil ***REMOVED***
					container.SharedEndpointList = append(container.SharedEndpointList, gwInfo["hnsid"].(string))
				***REMOVED***
			***REMOVED***

			if data["hnsid"] != nil ***REMOVED***
				container.SharedEndpointList = append(container.SharedEndpointList, data["hnsid"].(string))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
