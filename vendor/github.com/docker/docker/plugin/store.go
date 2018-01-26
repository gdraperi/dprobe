package plugin

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/pkg/plugins"
	"github.com/docker/docker/plugin/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/* allowV1PluginsFallback determines daemon's support for V1 plugins.
 * When the time comes to remove support for V1 plugins, flipping
 * this bool is all that will be needed.
 */
const allowV1PluginsFallback bool = true

/* defaultAPIVersion is the version of the plugin API for volume, network,
   IPAM and authz. This is a very stable API. When we update this API, then
   pluginType should include a version. e.g. "networkdriver/2.0".
*/
const defaultAPIVersion string = "1.0"

// GetV2Plugin retrieves a plugin by name, id or partial ID.
func (ps *Store) GetV2Plugin(refOrID string) (*v2.Plugin, error) ***REMOVED***
	ps.RLock()
	defer ps.RUnlock()

	id, err := ps.resolvePluginID(refOrID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p, idOk := ps.plugins[id]
	if !idOk ***REMOVED***
		return nil, errors.WithStack(errNotFound(id))
	***REMOVED***

	return p, nil
***REMOVED***

// validateName returns error if name is already reserved. always call with lock and full name
func (ps *Store) validateName(name string) error ***REMOVED***
	for _, p := range ps.plugins ***REMOVED***
		if p.Name() == name ***REMOVED***
			return alreadyExistsError(name)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// GetAll retrieves all plugins.
func (ps *Store) GetAll() map[string]*v2.Plugin ***REMOVED***
	ps.RLock()
	defer ps.RUnlock()
	return ps.plugins
***REMOVED***

// SetAll initialized plugins during daemon restore.
func (ps *Store) SetAll(plugins map[string]*v2.Plugin) ***REMOVED***
	ps.Lock()
	defer ps.Unlock()
	ps.plugins = plugins
***REMOVED***

func (ps *Store) getAllByCap(capability string) []plugingetter.CompatPlugin ***REMOVED***
	ps.RLock()
	defer ps.RUnlock()

	result := make([]plugingetter.CompatPlugin, 0, 1)
	for _, p := range ps.plugins ***REMOVED***
		if p.IsEnabled() ***REMOVED***
			if _, err := p.FilterByCap(capability); err == nil ***REMOVED***
				result = append(result, p)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***

// SetState sets the active state of the plugin and updates plugindb.
func (ps *Store) SetState(p *v2.Plugin, state bool) ***REMOVED***
	ps.Lock()
	defer ps.Unlock()

	p.PluginObj.Enabled = state
***REMOVED***

// Add adds a plugin to memory and plugindb.
// An error will be returned if there is a collision.
func (ps *Store) Add(p *v2.Plugin) error ***REMOVED***
	ps.Lock()
	defer ps.Unlock()

	if v, exist := ps.plugins[p.GetID()]; exist ***REMOVED***
		return fmt.Errorf("plugin %q has the same ID %s as %q", p.Name(), p.GetID(), v.Name())
	***REMOVED***
	ps.plugins[p.GetID()] = p
	return nil
***REMOVED***

// Remove removes a plugin from memory and plugindb.
func (ps *Store) Remove(p *v2.Plugin) ***REMOVED***
	ps.Lock()
	delete(ps.plugins, p.GetID())
	ps.Unlock()
***REMOVED***

// Get returns an enabled plugin matching the given name and capability.
func (ps *Store) Get(name, capability string, mode int) (plugingetter.CompatPlugin, error) ***REMOVED***
	// Lookup using new model.
	if ps != nil ***REMOVED***
		p, err := ps.GetV2Plugin(name)
		if err == nil ***REMOVED***
			if p.IsEnabled() ***REMOVED***
				fp, err := p.FilterByCap(capability)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				p.AddRefCount(mode)
				return fp, nil
			***REMOVED***

			// Plugin was found but it is disabled, so we should not fall back to legacy plugins
			// but we should error out right away
			return nil, errDisabled(name)
		***REMOVED***
		if _, ok := errors.Cause(err).(errNotFound); !ok ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if !allowV1PluginsFallback ***REMOVED***
		return nil, errNotFound(name)
	***REMOVED***

	p, err := plugins.Get(name, capability)
	if err == nil ***REMOVED***
		return p, nil
	***REMOVED***
	if errors.Cause(err) == plugins.ErrNotFound ***REMOVED***
		return nil, errNotFound(name)
	***REMOVED***
	return nil, errors.Wrap(errdefs.System(err), "legacy plugin")
***REMOVED***

// GetAllManagedPluginsByCap returns a list of managed plugins matching the given capability.
func (ps *Store) GetAllManagedPluginsByCap(capability string) []plugingetter.CompatPlugin ***REMOVED***
	return ps.getAllByCap(capability)
***REMOVED***

// GetAllByCap returns a list of enabled plugins matching the given capability.
func (ps *Store) GetAllByCap(capability string) ([]plugingetter.CompatPlugin, error) ***REMOVED***
	result := make([]plugingetter.CompatPlugin, 0, 1)

	/* Daemon start always calls plugin.Init thereby initializing a store.
	 * So store on experimental builds can never be nil, even while
	 * handling legacy plugins. However, there are legacy plugin unit
	 * tests where the volume subsystem directly talks with the plugin,
	 * bypassing the daemon. For such tests, this check is necessary.
	 */
	if ps != nil ***REMOVED***
		ps.RLock()
		result = ps.getAllByCap(capability)
		ps.RUnlock()
	***REMOVED***

	// Lookup with legacy model
	if allowV1PluginsFallback ***REMOVED***
		pl, err := plugins.GetAll(capability)
		if err != nil ***REMOVED***
			return nil, errors.Wrap(errdefs.System(err), "legacy plugin")
		***REMOVED***
		for _, p := range pl ***REMOVED***
			result = append(result, p)
		***REMOVED***
	***REMOVED***
	return result, nil
***REMOVED***

// Handle sets a callback for a given capability. It is only used by network
// and ipam drivers during plugin registration. The callback registers the
// driver with the subsystem (network, ipam).
func (ps *Store) Handle(capability string, callback func(string, *plugins.Client)) ***REMOVED***
	pluginType := fmt.Sprintf("docker.%s/%s", strings.ToLower(capability), defaultAPIVersion)

	// Register callback with new plugin model.
	ps.Lock()
	handlers, ok := ps.handlers[pluginType]
	if !ok ***REMOVED***
		handlers = []func(string, *plugins.Client)***REMOVED******REMOVED***
	***REMOVED***
	handlers = append(handlers, callback)
	ps.handlers[pluginType] = handlers
	ps.Unlock()

	// Register callback with legacy plugin model.
	if allowV1PluginsFallback ***REMOVED***
		plugins.Handle(capability, callback)
	***REMOVED***
***REMOVED***

// CallHandler calls the registered callback. It is invoked during plugin enable.
func (ps *Store) CallHandler(p *v2.Plugin) ***REMOVED***
	for _, typ := range p.GetTypes() ***REMOVED***
		for _, handler := range ps.handlers[typ.String()] ***REMOVED***
			handler(p.Name(), p.Client())
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ps *Store) resolvePluginID(idOrName string) (string, error) ***REMOVED***
	ps.RLock() // todo: fix
	defer ps.RUnlock()

	if validFullID.MatchString(idOrName) ***REMOVED***
		return idOrName, nil
	***REMOVED***

	ref, err := reference.ParseNormalizedNamed(idOrName)
	if err != nil ***REMOVED***
		return "", errors.WithStack(errNotFound(idOrName))
	***REMOVED***
	if _, ok := ref.(reference.Canonical); ok ***REMOVED***
		logrus.Warnf("canonical references cannot be resolved: %v", reference.FamiliarString(ref))
		return "", errors.WithStack(errNotFound(idOrName))
	***REMOVED***

	ref = reference.TagNameOnly(ref)

	for _, p := range ps.plugins ***REMOVED***
		if p.PluginObj.Name == reference.FamiliarString(ref) ***REMOVED***
			return p.PluginObj.ID, nil
		***REMOVED***
	***REMOVED***

	var found *v2.Plugin
	for id, p := range ps.plugins ***REMOVED*** // this can be optimized
		if strings.HasPrefix(id, idOrName) ***REMOVED***
			if found != nil ***REMOVED***
				return "", errors.WithStack(errAmbiguous(idOrName))
			***REMOVED***
			found = p
		***REMOVED***
	***REMOVED***
	if found == nil ***REMOVED***
		return "", errors.WithStack(errNotFound(idOrName))
	***REMOVED***
	return found.PluginObj.ID, nil
***REMOVED***
