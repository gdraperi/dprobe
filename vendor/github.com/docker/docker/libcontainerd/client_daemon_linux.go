package libcontainerd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/docker/docker/pkg/idtools"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
)

func summaryFromInterface(i interface***REMOVED******REMOVED***) (*Summary, error) ***REMOVED***
	return &Summary***REMOVED******REMOVED***, nil
***REMOVED***

func (c *client) UpdateResources(ctx context.Context, containerID string, resources *Resources) error ***REMOVED***
	p, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// go doesn't like the alias in 1.8, this means this need to be
	// platform specific
	return p.(containerd.Task).Update(ctx, containerd.WithResources((*specs.LinuxResources)(resources)))
***REMOVED***

func hostIDFromMap(id uint32, mp []specs.LinuxIDMapping) int ***REMOVED***
	for _, m := range mp ***REMOVED***
		if id >= m.ContainerID && id <= m.ContainerID+m.Size-1 ***REMOVED***
			return int(m.HostID + id - m.ContainerID)
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func getSpecUser(ociSpec *specs.Spec) (int, int) ***REMOVED***
	var (
		uid int
		gid int
	)

	for _, ns := range ociSpec.Linux.Namespaces ***REMOVED***
		if ns.Type == specs.UserNamespace ***REMOVED***
			uid = hostIDFromMap(0, ociSpec.Linux.UIDMappings)
			gid = hostIDFromMap(0, ociSpec.Linux.GIDMappings)
			break
		***REMOVED***
	***REMOVED***

	return uid, gid
***REMOVED***

func prepareBundleDir(bundleDir string, ociSpec *specs.Spec) (string, error) ***REMOVED***
	uid, gid := getSpecUser(ociSpec)
	if uid == 0 && gid == 0 ***REMOVED***
		return bundleDir, idtools.MkdirAllAndChownNew(bundleDir, 0755, idtools.IDPair***REMOVED***0, 0***REMOVED***)
	***REMOVED***

	p := string(filepath.Separator)
	components := strings.Split(bundleDir, string(filepath.Separator))
	for _, d := range components[1:] ***REMOVED***
		p = filepath.Join(p, d)
		fi, err := os.Stat(p)
		if err != nil && !os.IsNotExist(err) ***REMOVED***
			return "", err
		***REMOVED***
		if os.IsNotExist(err) || fi.Mode()&1 == 0 ***REMOVED***
			p = fmt.Sprintf("%s.%d.%d", p, uid, gid)
			if err := idtools.MkdirAndChown(p, 0700, idtools.IDPair***REMOVED***uid, gid***REMOVED***); err != nil && !os.IsExist(err) ***REMOVED***
				return "", err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return p, nil
***REMOVED***

func newFIFOSet(bundleDir, processID string, withStdin, withTerminal bool) *cio.FIFOSet ***REMOVED***
	config := cio.Config***REMOVED***
		Terminal: withTerminal,
		Stdout:   filepath.Join(bundleDir, processID+"-stdout"),
	***REMOVED***
	paths := []string***REMOVED***config.Stdout***REMOVED***

	if withStdin ***REMOVED***
		config.Stdin = filepath.Join(bundleDir, processID+"-stdin")
		paths = append(paths, config.Stdin)
	***REMOVED***
	if !withTerminal ***REMOVED***
		config.Stderr = filepath.Join(bundleDir, processID+"-stderr")
		paths = append(paths, config.Stderr)
	***REMOVED***
	closer := func() error ***REMOVED***
		for _, path := range paths ***REMOVED***
			if err := os.RemoveAll(path); err != nil ***REMOVED***
				logrus.Warnf("libcontainerd: failed to remove fifo %v: %v", path, err)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	return cio.NewFIFOSet(config, closer)
***REMOVED***
