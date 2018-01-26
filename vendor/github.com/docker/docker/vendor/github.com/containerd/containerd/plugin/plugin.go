package plugin

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	// ErrNoType is returned when no type is specified
	ErrNoType = errors.New("plugin: no type")
	// ErrNoPluginID is returned when no id is specified
	ErrNoPluginID = errors.New("plugin: no id")

	// ErrSkipPlugin is used when a plugin is not initialized and should not be loaded,
	// this allows the plugin loader differentiate between a plugin which is configured
	// not to load and one that fails to load.
	ErrSkipPlugin = errors.New("skip plugin")

	// ErrInvalidRequires will be thrown if the requirements for a plugin are
	// defined in an invalid manner.
	ErrInvalidRequires = errors.New("invalid requires")
)

// IsSkipPlugin returns true if the error is skipping the plugin
func IsSkipPlugin(err error) bool ***REMOVED***
	if errors.Cause(err) == ErrSkipPlugin ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// Type is the type of the plugin
type Type string

func (t Type) String() string ***REMOVED*** return string(t) ***REMOVED***

const (
	// AllPlugins declares that the plugin should be initialized after all others.
	AllPlugins Type = "*"
	// RuntimePlugin implements a runtime
	RuntimePlugin Type = "io.containerd.runtime.v1"
	// GRPCPlugin implements a grpc service
	GRPCPlugin Type = "io.containerd.grpc.v1"
	// SnapshotPlugin implements a snapshotter
	SnapshotPlugin Type = "io.containerd.snapshotter.v1"
	// TaskMonitorPlugin implements a task monitor
	TaskMonitorPlugin Type = "io.containerd.monitor.v1"
	// DiffPlugin implements a differ
	DiffPlugin Type = "io.containerd.differ.v1"
	// MetadataPlugin implements a metadata store
	MetadataPlugin Type = "io.containerd.metadata.v1"
	// ContentPlugin implements a content store
	ContentPlugin Type = "io.containerd.content.v1"
	// GCPlugin implements garbage collection policy
	GCPlugin Type = "io.containerd.gc.v1"
)

// Registration contains information for registering a plugin
type Registration struct ***REMOVED***
	// Type of the plugin
	Type Type
	// ID of the plugin
	ID string
	// Config specific to the plugin
	Config interface***REMOVED******REMOVED***
	// Requires is a list of plugins that the registered plugin requires to be available
	Requires []Type

	// InitFn is called when initializing a plugin. The registration and
	// context are passed in. The init function may modify the registration to
	// add exports, capabilites and platform support declarations.
	InitFn func(*InitContext) (interface***REMOVED******REMOVED***, error)
***REMOVED***

// Init the registered plugin
func (r *Registration) Init(ic *InitContext) *Plugin ***REMOVED***
	p, err := r.InitFn(ic)
	return &Plugin***REMOVED***
		Registration: r,
		Config:       ic.Config,
		Meta:         ic.Meta,
		instance:     p,
		err:          err,
	***REMOVED***
***REMOVED***

// URI returns the full plugin URI
func (r *Registration) URI() string ***REMOVED***
	return fmt.Sprintf("%s.%s", r.Type, r.ID)
***REMOVED***

// Service allows GRPC services to be registered with the underlying server
type Service interface ***REMOVED***
	Register(*grpc.Server) error
***REMOVED***

var register = struct ***REMOVED***
	sync.RWMutex
	r []*Registration
***REMOVED******REMOVED******REMOVED***

// Load loads all plugins at the provided path into containerd
func Load(path string) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if v := recover(); v != nil ***REMOVED***
			rerr, ok := v.(error)
			if !ok ***REMOVED***
				rerr = fmt.Errorf("%s", v)
			***REMOVED***
			err = rerr
		***REMOVED***
	***REMOVED***()
	return loadPlugins(path)
***REMOVED***

// Register allows plugins to register
func Register(r *Registration) ***REMOVED***
	register.Lock()
	defer register.Unlock()
	if r.Type == "" ***REMOVED***
		panic(ErrNoType)
	***REMOVED***
	if r.ID == "" ***REMOVED***
		panic(ErrNoPluginID)
	***REMOVED***

	var last bool
	for _, requires := range r.Requires ***REMOVED***
		if requires == "*" ***REMOVED***
			last = true
		***REMOVED***
	***REMOVED***
	if last && len(r.Requires) != 1 ***REMOVED***
		panic(ErrInvalidRequires)
	***REMOVED***

	register.r = append(register.r, r)
***REMOVED***

// Graph returns an ordered list of registered plugins for initialization
func Graph() (ordered []*Registration) ***REMOVED***
	register.RLock()
	defer register.RUnlock()

	added := map[*Registration]bool***REMOVED******REMOVED***
	for _, r := range register.r ***REMOVED***

		children(r.ID, r.Requires, added, &ordered)
		if !added[r] ***REMOVED***
			ordered = append(ordered, r)
			added[r] = true
		***REMOVED***
	***REMOVED***
	return ordered
***REMOVED***

func children(id string, types []Type, added map[*Registration]bool, ordered *[]*Registration) ***REMOVED***
	for _, t := range types ***REMOVED***
		for _, r := range register.r ***REMOVED***
			if r.ID != id && (t == "*" || r.Type == t) ***REMOVED***
				children(r.ID, r.Requires, added, ordered)
				if !added[r] ***REMOVED***
					*ordered = append(*ordered, r)
					added[r] = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
