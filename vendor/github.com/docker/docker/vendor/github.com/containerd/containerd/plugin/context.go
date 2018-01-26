package plugin

import (
	"context"
	"path/filepath"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/events/exchange"
	"github.com/containerd/containerd/log"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// InitContext is used for plugin inititalization
type InitContext struct ***REMOVED***
	Context context.Context
	Root    string
	State   string
	Config  interface***REMOVED******REMOVED***
	Address string
	Events  *exchange.Exchange

	Meta *Meta // plugins can fill in metadata at init.

	plugins *Set
***REMOVED***

// NewContext returns a new plugin InitContext
func NewContext(ctx context.Context, r *Registration, plugins *Set, root, state string) *InitContext ***REMOVED***
	return &InitContext***REMOVED***
		Context: log.WithModule(ctx, r.URI()),
		Root:    filepath.Join(root, r.URI()),
		State:   filepath.Join(state, r.URI()),
		Meta: &Meta***REMOVED***
			Exports: map[string]string***REMOVED******REMOVED***,
		***REMOVED***,
		plugins: plugins,
	***REMOVED***
***REMOVED***

// Get returns the first plugin by its type
func (i *InitContext) Get(t Type) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return i.plugins.Get(t)
***REMOVED***

// Meta contains information gathered from the registration and initialization
// process.
type Meta struct ***REMOVED***
	Platforms    []ocispec.Platform // platforms supported by plugin
	Exports      map[string]string  // values exported by plugin
	Capabilities []string           // feature switches for plugin
***REMOVED***

// Plugin represents an initialized plugin, used with an init context.
type Plugin struct ***REMOVED***
	Registration *Registration // registration, as initialized
	Config       interface***REMOVED******REMOVED***   // config, as initialized
	Meta         *Meta

	instance interface***REMOVED******REMOVED***
	err      error // will be set if there was an error initializing the plugin
***REMOVED***

// Err returns the errors during initialization.
// returns nil if not error was encountered
func (p *Plugin) Err() error ***REMOVED***
	return p.err
***REMOVED***

// Instance returns the instance and any initialization error of the plugin
func (p *Plugin) Instance() (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return p.instance, p.err
***REMOVED***

// Set defines a plugin collection, used with InitContext.
//
// This maintains ordering and unique indexing over the set.
//
// After iteratively instantiating plugins, this set should represent, the
// ordered, initialization set of plugins for a containerd instance.
type Set struct ***REMOVED***
	ordered     []*Plugin // order of initialization
	byTypeAndID map[Type]map[string]*Plugin
***REMOVED***

// NewPluginSet returns an initialized plugin set
func NewPluginSet() *Set ***REMOVED***
	return &Set***REMOVED***
		byTypeAndID: make(map[Type]map[string]*Plugin),
	***REMOVED***
***REMOVED***

// Add a plugin to the set
func (ps *Set) Add(p *Plugin) error ***REMOVED***
	if byID, typeok := ps.byTypeAndID[p.Registration.Type]; !typeok ***REMOVED***
		ps.byTypeAndID[p.Registration.Type] = map[string]*Plugin***REMOVED***
			p.Registration.ID: p,
		***REMOVED***
	***REMOVED*** else if _, idok := byID[p.Registration.ID]; !idok ***REMOVED***
		byID[p.Registration.ID] = p
	***REMOVED*** else ***REMOVED***
		return errors.Wrapf(errdefs.ErrAlreadyExists, "plugin %v already initialized", p.Registration.URI())
	***REMOVED***

	ps.ordered = append(ps.ordered, p)
	return nil
***REMOVED***

// Get returns the first plugin by its type
func (ps *Set) Get(t Type) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	for _, v := range ps.byTypeAndID[t] ***REMOVED***
		return v.Instance()
	***REMOVED***
	return nil, errors.Wrapf(errdefs.ErrNotFound, "no plugins registered for %s", t)
***REMOVED***

// GetAll plugins in the set
func (i *InitContext) GetAll() []*Plugin ***REMOVED***
	return i.plugins.ordered
***REMOVED***

// GetByType returns all plugins with the specific type.
func (i *InitContext) GetByType(t Type) (map[string]*Plugin, error) ***REMOVED***
	p, ok := i.plugins.byTypeAndID[t]
	if !ok ***REMOVED***
		return nil, errors.Wrapf(errdefs.ErrNotFound, "no plugins registered for %s", t)
	***REMOVED***

	return p, nil
***REMOVED***
