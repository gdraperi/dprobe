package daemon

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types/backend"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/builder/dockerfile"
	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/system"
	"github.com/pkg/errors"
)

// merge merges two Config, the image container configuration (defaults values),
// and the user container configuration, either passed by the API or generated
// by the cli.
// It will mutate the specified user configuration (userConf) with the image
// configuration where the user configuration is incomplete.
func merge(userConf, imageConf *containertypes.Config) error ***REMOVED***
	if userConf.User == "" ***REMOVED***
		userConf.User = imageConf.User
	***REMOVED***
	if len(userConf.ExposedPorts) == 0 ***REMOVED***
		userConf.ExposedPorts = imageConf.ExposedPorts
	***REMOVED*** else if imageConf.ExposedPorts != nil ***REMOVED***
		for port := range imageConf.ExposedPorts ***REMOVED***
			if _, exists := userConf.ExposedPorts[port]; !exists ***REMOVED***
				userConf.ExposedPorts[port] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if len(userConf.Env) == 0 ***REMOVED***
		userConf.Env = imageConf.Env
	***REMOVED*** else ***REMOVED***
		for _, imageEnv := range imageConf.Env ***REMOVED***
			found := false
			imageEnvKey := strings.Split(imageEnv, "=")[0]
			for _, userEnv := range userConf.Env ***REMOVED***
				userEnvKey := strings.Split(userEnv, "=")[0]
				if runtime.GOOS == "windows" ***REMOVED***
					// Case insensitive environment variables on Windows
					imageEnvKey = strings.ToUpper(imageEnvKey)
					userEnvKey = strings.ToUpper(userEnvKey)
				***REMOVED***
				if imageEnvKey == userEnvKey ***REMOVED***
					found = true
					break
				***REMOVED***
			***REMOVED***
			if !found ***REMOVED***
				userConf.Env = append(userConf.Env, imageEnv)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if userConf.Labels == nil ***REMOVED***
		userConf.Labels = map[string]string***REMOVED******REMOVED***
	***REMOVED***
	for l, v := range imageConf.Labels ***REMOVED***
		if _, ok := userConf.Labels[l]; !ok ***REMOVED***
			userConf.Labels[l] = v
		***REMOVED***
	***REMOVED***

	if len(userConf.Entrypoint) == 0 ***REMOVED***
		if len(userConf.Cmd) == 0 ***REMOVED***
			userConf.Cmd = imageConf.Cmd
			userConf.ArgsEscaped = imageConf.ArgsEscaped
		***REMOVED***

		if userConf.Entrypoint == nil ***REMOVED***
			userConf.Entrypoint = imageConf.Entrypoint
		***REMOVED***
	***REMOVED***
	if imageConf.Healthcheck != nil ***REMOVED***
		if userConf.Healthcheck == nil ***REMOVED***
			userConf.Healthcheck = imageConf.Healthcheck
		***REMOVED*** else ***REMOVED***
			if len(userConf.Healthcheck.Test) == 0 ***REMOVED***
				userConf.Healthcheck.Test = imageConf.Healthcheck.Test
			***REMOVED***
			if userConf.Healthcheck.Interval == 0 ***REMOVED***
				userConf.Healthcheck.Interval = imageConf.Healthcheck.Interval
			***REMOVED***
			if userConf.Healthcheck.Timeout == 0 ***REMOVED***
				userConf.Healthcheck.Timeout = imageConf.Healthcheck.Timeout
			***REMOVED***
			if userConf.Healthcheck.StartPeriod == 0 ***REMOVED***
				userConf.Healthcheck.StartPeriod = imageConf.Healthcheck.StartPeriod
			***REMOVED***
			if userConf.Healthcheck.Retries == 0 ***REMOVED***
				userConf.Healthcheck.Retries = imageConf.Healthcheck.Retries
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if userConf.WorkingDir == "" ***REMOVED***
		userConf.WorkingDir = imageConf.WorkingDir
	***REMOVED***
	if len(userConf.Volumes) == 0 ***REMOVED***
		userConf.Volumes = imageConf.Volumes
	***REMOVED*** else ***REMOVED***
		for k, v := range imageConf.Volumes ***REMOVED***
			userConf.Volumes[k] = v
		***REMOVED***
	***REMOVED***

	if userConf.StopSignal == "" ***REMOVED***
		userConf.StopSignal = imageConf.StopSignal
	***REMOVED***
	return nil
***REMOVED***

// Commit creates a new filesystem image from the current state of a container.
// The image can optionally be tagged into a repository.
func (daemon *Daemon) Commit(name string, c *backend.ContainerCommitConfig) (string, error) ***REMOVED***
	start := time.Now()
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	// It is not possible to commit a running container on Windows
	if (runtime.GOOS == "windows") && container.IsRunning() ***REMOVED***
		return "", errors.Errorf("%+v does not support commit of a running container", runtime.GOOS)
	***REMOVED***

	if container.IsDead() ***REMOVED***
		err := fmt.Errorf("You cannot commit container %s which is Dead", container.ID)
		return "", errdefs.Conflict(err)
	***REMOVED***

	if container.IsRemovalInProgress() ***REMOVED***
		err := fmt.Errorf("You cannot commit container %s which is being removed", container.ID)
		return "", errdefs.Conflict(err)
	***REMOVED***

	if c.Pause && !container.IsPaused() ***REMOVED***
		daemon.containerPause(container)
		defer daemon.containerUnpause(container)
	***REMOVED***
	if !system.IsOSSupported(container.OS) ***REMOVED***
		return "", system.ErrNotSupportedOperatingSystem
	***REMOVED***

	if c.MergeConfigs && c.Config == nil ***REMOVED***
		c.Config = container.Config
	***REMOVED***

	newConfig, err := dockerfile.BuildFromConfig(c.Config, c.Changes, container.OS)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if c.MergeConfigs ***REMOVED***
		if err := merge(newConfig, container.Config); err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***

	rwTar, err := daemon.exportContainerRw(container)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer func() ***REMOVED***
		if rwTar != nil ***REMOVED***
			rwTar.Close()
		***REMOVED***
	***REMOVED***()

	var parent *image.Image
	if container.ImageID == "" ***REMOVED***
		parent = new(image.Image)
		parent.RootFS = image.NewRootFS()
	***REMOVED*** else ***REMOVED***
		parent, err = daemon.imageStore.Get(container.ImageID)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***

	l, err := daemon.layerStores[container.OS].Register(rwTar, parent.RootFS.ChainID())
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer layer.ReleaseAndLog(daemon.layerStores[container.OS], l)

	containerConfig := c.ContainerConfig
	if containerConfig == nil ***REMOVED***
		containerConfig = container.Config
	***REMOVED***
	cc := image.ChildConfig***REMOVED***
		ContainerID:     container.ID,
		Author:          c.Author,
		Comment:         c.Comment,
		ContainerConfig: containerConfig,
		Config:          newConfig,
		DiffID:          l.DiffID(),
	***REMOVED***
	config, err := json.Marshal(image.NewChildImage(parent, cc, container.OS))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	id, err := daemon.imageStore.Create(config)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if container.ImageID != "" ***REMOVED***
		if err := daemon.imageStore.SetParent(id, container.ImageID); err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***

	imageRef := ""
	if c.Repo != "" ***REMOVED***
		newTag, err := reference.ParseNormalizedNamed(c.Repo) // todo: should move this to API layer
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if !reference.IsNameOnly(newTag) ***REMOVED***
			return "", errors.Errorf("unexpected repository name: %s", c.Repo)
		***REMOVED***
		if c.Tag != "" ***REMOVED***
			if newTag, err = reference.WithTag(newTag, c.Tag); err != nil ***REMOVED***
				return "", err
			***REMOVED***
		***REMOVED***
		if err := daemon.TagImageWithReference(id, newTag); err != nil ***REMOVED***
			return "", err
		***REMOVED***
		imageRef = reference.FamiliarString(newTag)
	***REMOVED***

	attributes := map[string]string***REMOVED***
		"comment":  c.Comment,
		"imageID":  id.String(),
		"imageRef": imageRef,
	***REMOVED***
	daemon.LogContainerEventWithAttributes(container, "commit", attributes)
	containerActions.WithValues("commit").UpdateSince(start)
	return id.String(), nil
***REMOVED***

func (daemon *Daemon) exportContainerRw(container *container.Container) (arch io.ReadCloser, err error) ***REMOVED***
	// Note: Indexing by OS is safe as only called from `Commit` which has already performed validation
	rwlayer, err := daemon.layerStores[container.OS].GetRWLayer(container.ID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			daemon.layerStores[container.OS].ReleaseRWLayer(rwlayer)
		***REMOVED***
	***REMOVED***()

	// TODO: this mount call is not necessary as we assume that TarStream() should
	// mount the layer if needed. But the Diff() function for windows requests that
	// the layer should be mounted when calling it. So we reserve this mount call
	// until windows driver can implement Diff() interface correctly.
	_, err = rwlayer.Mount(container.GetMountLabel())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	archive, err := rwlayer.TarStream()
	if err != nil ***REMOVED***
		rwlayer.Unmount()
		return nil, err
	***REMOVED***
	return ioutils.NewReadCloserWrapper(archive, func() error ***REMOVED***
			archive.Close()
			err = rwlayer.Unmount()
			daemon.layerStores[container.OS].ReleaseRWLayer(rwlayer)
			return err
		***REMOVED***),
		nil
***REMOVED***
