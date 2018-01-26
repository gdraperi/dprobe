// +build !windows

package libcontainerd

import "fmt"

// WithRemoteAddr sets the external containerd socket to connect to.
func WithRemoteAddr(addr string) RemoteOption ***REMOVED***
	return rpcAddr(addr)
***REMOVED***

type rpcAddr string

func (a rpcAddr) Apply(r Remote) error ***REMOVED***
	if remote, ok := r.(*remote); ok ***REMOVED***
		remote.GRPC.Address = string(a)
		return nil
	***REMOVED***
	return fmt.Errorf("WithRemoteAddr option not supported for this remote")
***REMOVED***

// WithRemoteAddrUser sets the uid and gid to create the RPC address with
func WithRemoteAddrUser(uid, gid int) RemoteOption ***REMOVED***
	return rpcUser***REMOVED***uid, gid***REMOVED***
***REMOVED***

type rpcUser struct ***REMOVED***
	uid int
	gid int
***REMOVED***

func (u rpcUser) Apply(r Remote) error ***REMOVED***
	if remote, ok := r.(*remote); ok ***REMOVED***
		remote.GRPC.UID = u.uid
		remote.GRPC.GID = u.gid
		return nil
	***REMOVED***
	return fmt.Errorf("WithRemoteAddr option not supported for this remote")
***REMOVED***

// WithStartDaemon defines if libcontainerd should also run containerd daemon.
func WithStartDaemon(start bool) RemoteOption ***REMOVED***
	return startDaemon(start)
***REMOVED***

type startDaemon bool

func (s startDaemon) Apply(r Remote) error ***REMOVED***
	if remote, ok := r.(*remote); ok ***REMOVED***
		remote.startDaemon = bool(s)
		return nil
	***REMOVED***
	return fmt.Errorf("WithStartDaemon option not supported for this remote")
***REMOVED***

// WithLogLevel defines which log level to starts containerd with.
// This only makes sense if WithStartDaemon() was set to true.
func WithLogLevel(lvl string) RemoteOption ***REMOVED***
	return logLevel(lvl)
***REMOVED***

type logLevel string

func (l logLevel) Apply(r Remote) error ***REMOVED***
	if remote, ok := r.(*remote); ok ***REMOVED***
		remote.Debug.Level = string(l)
		return nil
	***REMOVED***
	return fmt.Errorf("WithDebugLog option not supported for this remote")
***REMOVED***

// WithDebugAddress defines at which location the debug GRPC connection
// should be made
func WithDebugAddress(addr string) RemoteOption ***REMOVED***
	return debugAddress(addr)
***REMOVED***

type debugAddress string

func (d debugAddress) Apply(r Remote) error ***REMOVED***
	if remote, ok := r.(*remote); ok ***REMOVED***
		remote.Debug.Address = string(d)
		return nil
	***REMOVED***
	return fmt.Errorf("WithDebugAddress option not supported for this remote")
***REMOVED***

// WithMetricsAddress defines at which location the debug GRPC connection
// should be made
func WithMetricsAddress(addr string) RemoteOption ***REMOVED***
	return metricsAddress(addr)
***REMOVED***

type metricsAddress string

func (m metricsAddress) Apply(r Remote) error ***REMOVED***
	if remote, ok := r.(*remote); ok ***REMOVED***
		remote.Metrics.Address = string(m)
		return nil
	***REMOVED***
	return fmt.Errorf("WithMetricsAddress option not supported for this remote")
***REMOVED***

// WithSnapshotter defines snapshotter driver should be used
func WithSnapshotter(name string) RemoteOption ***REMOVED***
	return snapshotter(name)
***REMOVED***

type snapshotter string

func (s snapshotter) Apply(r Remote) error ***REMOVED***
	if remote, ok := r.(*remote); ok ***REMOVED***
		remote.snapshotter = string(s)
		return nil
	***REMOVED***
	return fmt.Errorf("WithSnapshotter option not supported for this remote")
***REMOVED***

// WithPlugin allow configuring a containerd plugin
// configuration values passed needs to be quoted if quotes are needed in
// the toml format.
func WithPlugin(name string, conf interface***REMOVED******REMOVED***) RemoteOption ***REMOVED***
	return pluginConf***REMOVED***
		name: name,
		conf: conf,
	***REMOVED***
***REMOVED***

type pluginConf struct ***REMOVED***
	// Name is the name of the plugin
	name string
	conf interface***REMOVED******REMOVED***
***REMOVED***

func (p pluginConf) Apply(r Remote) error ***REMOVED***
	if remote, ok := r.(*remote); ok ***REMOVED***
		remote.pluginConfs.Plugins[p.name] = p.conf
		return nil
	***REMOVED***
	return fmt.Errorf("WithPlugin option not supported for this remote")
***REMOVED***
