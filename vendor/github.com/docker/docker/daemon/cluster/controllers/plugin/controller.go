package plugin

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/docker/distribution/reference"
	enginetypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm/runtime"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/plugin"
	"github.com/docker/docker/plugin/v2"
	"github.com/docker/swarmkit/api"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// Controller is the controller for the plugin backend.
// Plugins are managed as a singleton object with a desired state (different from containers).
// With the the plugin controller instead of having a strict create->start->stop->remove
// task lifecycle like containers, we manage the desired state of the plugin and let
// the plugin manager do what it already does and monitor the plugin.
// We'll also end up with many tasks all pointing to the same plugin ID.
//
// TODO(@cpuguy83): registry auth is intentionally not supported until we work out
// the right way to pass registry crednetials via secrets.
type Controller struct ***REMOVED***
	backend Backend
	spec    runtime.PluginSpec
	logger  *logrus.Entry

	pluginID  string
	serviceID string
	taskID    string

	// hook used to signal tests that `Wait()` is actually ready and waiting
	signalWaitReady func()
***REMOVED***

// Backend is the interface for interacting with the plugin manager
// Controller actions are passed to the configured backend to do the real work.
type Backend interface ***REMOVED***
	Disable(name string, config *enginetypes.PluginDisableConfig) error
	Enable(name string, config *enginetypes.PluginEnableConfig) error
	Remove(name string, config *enginetypes.PluginRmConfig) error
	Pull(ctx context.Context, ref reference.Named, name string, metaHeaders http.Header, authConfig *enginetypes.AuthConfig, privileges enginetypes.PluginPrivileges, outStream io.Writer, opts ...plugin.CreateOpt) error
	Upgrade(ctx context.Context, ref reference.Named, name string, metaHeaders http.Header, authConfig *enginetypes.AuthConfig, privileges enginetypes.PluginPrivileges, outStream io.Writer) error
	Get(name string) (*v2.Plugin, error)
	SubscribeEvents(buffer int, events ...plugin.Event) (eventCh <-chan interface***REMOVED******REMOVED***, cancel func())
***REMOVED***

// NewController returns a new cluster plugin controller
func NewController(backend Backend, t *api.Task) (*Controller, error) ***REMOVED***
	spec, err := readSpec(t)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Controller***REMOVED***
		backend:   backend,
		spec:      spec,
		serviceID: t.ServiceID,
		logger: logrus.WithFields(logrus.Fields***REMOVED***
			"controller": "plugin",
			"task":       t.ID,
			"plugin":     spec.Name,
		***REMOVED***)***REMOVED***, nil
***REMOVED***

func readSpec(t *api.Task) (runtime.PluginSpec, error) ***REMOVED***
	var cfg runtime.PluginSpec

	generic := t.Spec.GetGeneric()
	if err := proto.Unmarshal(generic.Payload.Value, &cfg); err != nil ***REMOVED***
		return cfg, errors.Wrap(err, "error reading plugin spec")
	***REMOVED***
	return cfg, nil
***REMOVED***

// Update is the update phase from swarmkit
func (p *Controller) Update(ctx context.Context, t *api.Task) error ***REMOVED***
	p.logger.Debug("Update")
	return nil
***REMOVED***

// Prepare is the prepare phase from swarmkit
func (p *Controller) Prepare(ctx context.Context) (err error) ***REMOVED***
	p.logger.Debug("Prepare")

	remote, err := reference.ParseNormalizedNamed(p.spec.Remote)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "error parsing remote reference %q", p.spec.Remote)
	***REMOVED***

	if p.spec.Name == "" ***REMOVED***
		p.spec.Name = remote.String()
	***REMOVED***

	var authConfig enginetypes.AuthConfig
	privs := convertPrivileges(p.spec.Privileges)

	pl, err := p.backend.Get(p.spec.Name)

	defer func() ***REMOVED***
		if pl != nil && err == nil ***REMOVED***
			pl.Acquire()
		***REMOVED***
	***REMOVED***()

	if err == nil && pl != nil ***REMOVED***
		if pl.SwarmServiceID != p.serviceID ***REMOVED***
			return errors.Errorf("plugin already exists: %s", p.spec.Name)
		***REMOVED***
		if pl.IsEnabled() ***REMOVED***
			if err := p.backend.Disable(pl.GetID(), &enginetypes.PluginDisableConfig***REMOVED***ForceDisable: true***REMOVED***); err != nil ***REMOVED***
				p.logger.WithError(err).Debug("could not disable plugin before running upgrade")
			***REMOVED***
		***REMOVED***
		p.pluginID = pl.GetID()
		return p.backend.Upgrade(ctx, remote, p.spec.Name, nil, &authConfig, privs, ioutil.Discard)
	***REMOVED***

	if err := p.backend.Pull(ctx, remote, p.spec.Name, nil, &authConfig, privs, ioutil.Discard, plugin.WithSwarmService(p.serviceID)); err != nil ***REMOVED***
		return err
	***REMOVED***
	pl, err = p.backend.Get(p.spec.Name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	p.pluginID = pl.GetID()

	return nil
***REMOVED***

// Start is the start phase from swarmkit
func (p *Controller) Start(ctx context.Context) error ***REMOVED***
	p.logger.Debug("Start")

	pl, err := p.backend.Get(p.pluginID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if p.spec.Disabled ***REMOVED***
		if pl.IsEnabled() ***REMOVED***
			return p.backend.Disable(p.pluginID, &enginetypes.PluginDisableConfig***REMOVED***ForceDisable: false***REMOVED***)
		***REMOVED***
		return nil
	***REMOVED***
	if !pl.IsEnabled() ***REMOVED***
		return p.backend.Enable(p.pluginID, &enginetypes.PluginEnableConfig***REMOVED***Timeout: 30***REMOVED***)
	***REMOVED***
	return nil
***REMOVED***

// Wait causes the task to wait until returned
func (p *Controller) Wait(ctx context.Context) error ***REMOVED***
	p.logger.Debug("Wait")

	pl, err := p.backend.Get(p.pluginID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	events, cancel := p.backend.SubscribeEvents(1, plugin.EventDisable***REMOVED***Plugin: pl.PluginObj***REMOVED***, plugin.EventRemove***REMOVED***Plugin: pl.PluginObj***REMOVED***, plugin.EventEnable***REMOVED***Plugin: pl.PluginObj***REMOVED***)
	defer cancel()

	if p.signalWaitReady != nil ***REMOVED***
		p.signalWaitReady()
	***REMOVED***

	if !p.spec.Disabled != pl.IsEnabled() ***REMOVED***
		return errors.New("mismatched plugin state")
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case e := <-events:
			p.logger.Debugf("got event %#T", e)

			switch e.(type) ***REMOVED***
			case plugin.EventEnable:
				if p.spec.Disabled ***REMOVED***
					return errors.New("plugin enabled")
				***REMOVED***
			case plugin.EventRemove:
				return errors.New("plugin removed")
			case plugin.EventDisable:
				if !p.spec.Disabled ***REMOVED***
					return errors.New("plugin disabled")
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func isNotFound(err error) bool ***REMOVED***
	return errdefs.IsNotFound(err)
***REMOVED***

// Shutdown is the shutdown phase from swarmkit
func (p *Controller) Shutdown(ctx context.Context) error ***REMOVED***
	p.logger.Debug("Shutdown")
	return nil
***REMOVED***

// Terminate is the terminate phase from swarmkit
func (p *Controller) Terminate(ctx context.Context) error ***REMOVED***
	p.logger.Debug("Terminate")
	return nil
***REMOVED***

// Remove is the remove phase from swarmkit
func (p *Controller) Remove(ctx context.Context) error ***REMOVED***
	p.logger.Debug("Remove")

	pl, err := p.backend.Get(p.pluginID)
	if err != nil ***REMOVED***
		if isNotFound(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	pl.Release()
	if pl.GetRefCount() > 0 ***REMOVED***
		p.logger.Debug("skipping remove due to ref count")
		return nil
	***REMOVED***

	// This may error because we have exactly 1 plugin, but potentially multiple
	// tasks which are calling remove.
	err = p.backend.Remove(p.pluginID, &enginetypes.PluginRmConfig***REMOVED***ForceRemove: true***REMOVED***)
	if isNotFound(err) ***REMOVED***
		return nil
	***REMOVED***
	return err
***REMOVED***

// Close is the close phase from swarmkit
func (p *Controller) Close() error ***REMOVED***
	p.logger.Debug("Close")
	return nil
***REMOVED***

func convertPrivileges(ls []*runtime.PluginPrivilege) enginetypes.PluginPrivileges ***REMOVED***
	var out enginetypes.PluginPrivileges
	for _, p := range ls ***REMOVED***
		pp := enginetypes.PluginPrivilege***REMOVED***
			Name:        p.Name,
			Description: p.Description,
			Value:       p.Value,
		***REMOVED***
		out = append(out, pp)
	***REMOVED***
	return out
***REMOVED***
