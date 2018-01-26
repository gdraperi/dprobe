package server

import (
	"context"
	"os"

	"github.com/containerd/cgroups"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/sys"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// apply sets config settings on the server process
func apply(ctx context.Context, config *Config) error ***REMOVED***
	if !config.NoSubreaper ***REMOVED***
		log.G(ctx).Info("setting subreaper...")
		if err := sys.SetSubreaper(1); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if config.OOMScore != 0 ***REMOVED***
		log.G(ctx).Debugf("changing OOM score to %d", config.OOMScore)
		if err := sys.SetOOMScore(os.Getpid(), config.OOMScore); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed to change OOM score to %d", config.OOMScore)
		***REMOVED***
	***REMOVED***
	if config.Cgroup.Path != "" ***REMOVED***
		cg, err := cgroups.Load(cgroups.V1, cgroups.StaticPath(config.Cgroup.Path))
		if err != nil ***REMOVED***
			if err != cgroups.ErrCgroupDeleted ***REMOVED***
				return err
			***REMOVED***
			if cg, err = cgroups.New(cgroups.V1, cgroups.StaticPath(config.Cgroup.Path), &specs.LinuxResources***REMOVED******REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if err := cg.Add(cgroups.Process***REMOVED***
			Pid: os.Getpid(),
		***REMOVED***); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
