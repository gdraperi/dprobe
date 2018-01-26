// +build linux

package linux

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/containerd/containerd/events/exchange"
	"github.com/containerd/containerd/linux/runctypes"
	"github.com/containerd/containerd/linux/shim"
	"github.com/containerd/containerd/linux/shim/client"
	"github.com/pkg/errors"
)

// loadBundle loads an existing bundle from disk
func loadBundle(id, path, workdir string) *bundle ***REMOVED***
	return &bundle***REMOVED***
		id:      id,
		path:    path,
		workDir: workdir,
	***REMOVED***
***REMOVED***

// newBundle creates a new bundle on disk at the provided path for the given id
func newBundle(id, path, workDir string, spec []byte) (b *bundle, err error) ***REMOVED***
	if err := os.MkdirAll(path, 0711); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	path = filepath.Join(path, id)
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			os.RemoveAll(path)
		***REMOVED***
	***REMOVED***()
	workDir = filepath.Join(workDir, id)
	if err := os.MkdirAll(workDir, 0711); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			os.RemoveAll(workDir)
		***REMOVED***
	***REMOVED***()

	if err := os.Mkdir(path, 0711); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := os.Mkdir(filepath.Join(path, "rootfs"), 0711); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	f, err := os.Create(filepath.Join(path, configFilename))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()
	_, err = io.Copy(f, bytes.NewReader(spec))
	return &bundle***REMOVED***
		id:      id,
		path:    path,
		workDir: workDir,
	***REMOVED***, err
***REMOVED***

type bundle struct ***REMOVED***
	id      string
	path    string
	workDir string
***REMOVED***

// ShimOpt specifies shim options for initialization and connection
type ShimOpt func(*bundle, string, *runctypes.RuncOptions) (shim.Config, client.Opt)

// ShimRemote is a ShimOpt for connecting and starting a remote shim
func ShimRemote(shimBinary, daemonAddress, cgroup string, debug bool, exitHandler func()) ShimOpt ***REMOVED***
	return func(b *bundle, ns string, ropts *runctypes.RuncOptions) (shim.Config, client.Opt) ***REMOVED***
		return b.shimConfig(ns, ropts),
			client.WithStart(shimBinary, b.shimAddress(ns), daemonAddress, cgroup, debug, exitHandler)
	***REMOVED***
***REMOVED***

// ShimLocal is a ShimOpt for using an in process shim implementation
func ShimLocal(exchange *exchange.Exchange) ShimOpt ***REMOVED***
	return func(b *bundle, ns string, ropts *runctypes.RuncOptions) (shim.Config, client.Opt) ***REMOVED***
		return b.shimConfig(ns, ropts), client.WithLocal(exchange)
	***REMOVED***
***REMOVED***

// ShimConnect is a ShimOpt for connecting to an existing remote shim
func ShimConnect() ShimOpt ***REMOVED***
	return func(b *bundle, ns string, ropts *runctypes.RuncOptions) (shim.Config, client.Opt) ***REMOVED***
		return b.shimConfig(ns, ropts), client.WithConnect(b.shimAddress(ns))
	***REMOVED***
***REMOVED***

// NewShimClient connects to the shim managing the bundle and tasks creating it if needed
func (b *bundle) NewShimClient(ctx context.Context, namespace string, getClientOpts ShimOpt, runcOpts *runctypes.RuncOptions) (*client.Client, error) ***REMOVED***
	cfg, opt := getClientOpts(b, namespace, runcOpts)
	return client.New(ctx, cfg, opt)
***REMOVED***

// Delete deletes the bundle from disk
func (b *bundle) Delete() error ***REMOVED***
	err := os.RemoveAll(b.path)
	if err == nil ***REMOVED***
		return os.RemoveAll(b.workDir)
	***REMOVED***
	// error removing the bundle path; still attempt removing work dir
	err2 := os.RemoveAll(b.workDir)
	if err2 == nil ***REMOVED***
		return err
	***REMOVED***
	return errors.Wrapf(err, "Failed to remove both bundle and workdir locations: %v", err2)
***REMOVED***

func (b *bundle) shimAddress(namespace string) string ***REMOVED***
	return filepath.Join(string(filepath.Separator), "containerd-shim", namespace, b.id, "shim.sock")
***REMOVED***

func (b *bundle) shimConfig(namespace string, runcOptions *runctypes.RuncOptions) shim.Config ***REMOVED***
	var (
		criuPath      string
		runtimeRoot   string
		systemdCgroup bool
	)
	if runcOptions != nil ***REMOVED***
		criuPath = runcOptions.CriuPath
		systemdCgroup = runcOptions.SystemdCgroup
		runtimeRoot = runcOptions.RuntimeRoot
	***REMOVED***
	return shim.Config***REMOVED***
		Path:          b.path,
		WorkDir:       b.workDir,
		Namespace:     namespace,
		Criu:          criuPath,
		RuntimeRoot:   runtimeRoot,
		SystemdCgroup: systemdCgroup,
	***REMOVED***
***REMOVED***
