// Package plugins provides structures and helper functions to manage Docker
// plugins.
//
// Docker discovers plugins by looking for them in the plugin directory whenever
// a user or container tries to use one by name. UNIX domain socket files must
// be located under /run/docker/plugins, whereas spec files can be located
// either under /etc/docker/plugins or /usr/lib/docker/plugins. This is handled
// by the Registry interface, which lets you list all plugins or get a plugin by
// its name if it exists.
//
// The plugins need to implement an HTTP server and bind this to the UNIX socket
// or the address specified in the spec files.
// A handshake is send at /Plugin.Activate, and plugins are expected to return
// a Manifest with a list of of Docker subsystems which this plugin implements.
//
// In order to use a plugins, you can use the ``Get`` with the name of the
// plugin and the subsystem it implements.
//
//	plugin, err := plugins.Get("example", "VolumeDriver")
//	if err != nil ***REMOVED***
//		return fmt.Errorf("Error looking up volume plugin example: %v", err)
//	***REMOVED***
package plugins

import (
	"errors"
	"sync"
	"time"

	"github.com/docker/go-connections/tlsconfig"
	"github.com/sirupsen/logrus"
)

var (
	// ErrNotImplements is returned if the plugin does not implement the requested driver.
	ErrNotImplements = errors.New("Plugin does not implement the requested driver")
)

type plugins struct ***REMOVED***
	sync.Mutex
	plugins map[string]*Plugin
***REMOVED***

type extpointHandlers struct ***REMOVED***
	sync.RWMutex
	extpointHandlers map[string][]func(string, *Client)
***REMOVED***

var (
	storage  = plugins***REMOVED***plugins: make(map[string]*Plugin)***REMOVED***
	handlers = extpointHandlers***REMOVED***extpointHandlers: make(map[string][]func(string, *Client))***REMOVED***
)

// Manifest lists what a plugin implements.
type Manifest struct ***REMOVED***
	// List of subsystem the plugin implements.
	Implements []string
***REMOVED***

// Plugin is the definition of a docker plugin.
type Plugin struct ***REMOVED***
	// Name of the plugin
	name string
	// Address of the plugin
	Addr string
	// TLS configuration of the plugin
	TLSConfig *tlsconfig.Options
	// Client attached to the plugin
	client *Client
	// Manifest of the plugin (see above)
	Manifest *Manifest `json:"-"`

	// wait for activation to finish
	activateWait *sync.Cond
	// error produced by activation
	activateErr error
	// keeps track of callback handlers run against this plugin
	handlersRun bool
***REMOVED***

// Name returns the name of the plugin.
func (p *Plugin) Name() string ***REMOVED***
	return p.name
***REMOVED***

// Client returns a ready-to-use plugin client that can be used to communicate with the plugin.
func (p *Plugin) Client() *Client ***REMOVED***
	return p.client
***REMOVED***

// IsV1 returns true for V1 plugins and false otherwise.
func (p *Plugin) IsV1() bool ***REMOVED***
	return true
***REMOVED***

// NewLocalPlugin creates a new local plugin.
func NewLocalPlugin(name, addr string) *Plugin ***REMOVED***
	return &Plugin***REMOVED***
		name: name,
		Addr: addr,
		// TODO: change to nil
		TLSConfig:    &tlsconfig.Options***REMOVED***InsecureSkipVerify: true***REMOVED***,
		activateWait: sync.NewCond(&sync.Mutex***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

func (p *Plugin) activate() error ***REMOVED***
	p.activateWait.L.Lock()

	if p.activated() ***REMOVED***
		p.runHandlers()
		p.activateWait.L.Unlock()
		return p.activateErr
	***REMOVED***

	p.activateErr = p.activateWithLock()

	p.runHandlers()
	p.activateWait.L.Unlock()
	p.activateWait.Broadcast()
	return p.activateErr
***REMOVED***

// runHandlers runs the registered handlers for the implemented plugin types
// This should only be run after activation, and while the activation lock is held.
func (p *Plugin) runHandlers() ***REMOVED***
	if !p.activated() ***REMOVED***
		return
	***REMOVED***

	handlers.RLock()
	if !p.handlersRun ***REMOVED***
		for _, iface := range p.Manifest.Implements ***REMOVED***
			hdlrs, handled := handlers.extpointHandlers[iface]
			if !handled ***REMOVED***
				continue
			***REMOVED***
			for _, handler := range hdlrs ***REMOVED***
				handler(p.name, p.client)
			***REMOVED***
		***REMOVED***
		p.handlersRun = true
	***REMOVED***
	handlers.RUnlock()

***REMOVED***

// activated returns if the plugin has already been activated.
// This should only be called with the activation lock held
func (p *Plugin) activated() bool ***REMOVED***
	return p.Manifest != nil
***REMOVED***

func (p *Plugin) activateWithLock() error ***REMOVED***
	c, err := NewClient(p.Addr, p.TLSConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	p.client = c

	m := new(Manifest)
	if err = p.client.Call("Plugin.Activate", nil, m); err != nil ***REMOVED***
		return err
	***REMOVED***

	p.Manifest = m
	return nil
***REMOVED***

func (p *Plugin) waitActive() error ***REMOVED***
	p.activateWait.L.Lock()
	for !p.activated() && p.activateErr == nil ***REMOVED***
		p.activateWait.Wait()
	***REMOVED***
	p.activateWait.L.Unlock()
	return p.activateErr
***REMOVED***

func (p *Plugin) implements(kind string) bool ***REMOVED***
	if p.Manifest == nil ***REMOVED***
		return false
	***REMOVED***
	for _, driver := range p.Manifest.Implements ***REMOVED***
		if driver == kind ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func load(name string) (*Plugin, error) ***REMOVED***
	return loadWithRetry(name, true)
***REMOVED***

func loadWithRetry(name string, retry bool) (*Plugin, error) ***REMOVED***
	registry := newLocalRegistry()
	start := time.Now()

	var retries int
	for ***REMOVED***
		pl, err := registry.Plugin(name)
		if err != nil ***REMOVED***
			if !retry ***REMOVED***
				return nil, err
			***REMOVED***

			timeOff := backoff(retries)
			if abort(start, timeOff) ***REMOVED***
				return nil, err
			***REMOVED***
			retries++
			logrus.Warnf("Unable to locate plugin: %s, retrying in %v", name, timeOff)
			time.Sleep(timeOff)
			continue
		***REMOVED***

		storage.Lock()
		if pl, exists := storage.plugins[name]; exists ***REMOVED***
			storage.Unlock()
			return pl, pl.activate()
		***REMOVED***
		storage.plugins[name] = pl
		storage.Unlock()

		err = pl.activate()

		if err != nil ***REMOVED***
			storage.Lock()
			delete(storage.plugins, name)
			storage.Unlock()
		***REMOVED***

		return pl, err
	***REMOVED***
***REMOVED***

func get(name string) (*Plugin, error) ***REMOVED***
	storage.Lock()
	pl, ok := storage.plugins[name]
	storage.Unlock()
	if ok ***REMOVED***
		return pl, pl.activate()
	***REMOVED***
	return load(name)
***REMOVED***

// Get returns the plugin given the specified name and requested implementation.
func Get(name, imp string) (*Plugin, error) ***REMOVED***
	pl, err := get(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := pl.waitActive(); err == nil && pl.implements(imp) ***REMOVED***
		logrus.Debugf("%s implements: %s", name, imp)
		return pl, nil
	***REMOVED***
	return nil, ErrNotImplements
***REMOVED***

// Handle adds the specified function to the extpointHandlers.
func Handle(iface string, fn func(string, *Client)) ***REMOVED***
	handlers.Lock()
	hdlrs, ok := handlers.extpointHandlers[iface]
	if !ok ***REMOVED***
		hdlrs = []func(string, *Client)***REMOVED******REMOVED***
	***REMOVED***

	hdlrs = append(hdlrs, fn)
	handlers.extpointHandlers[iface] = hdlrs

	storage.Lock()
	for _, p := range storage.plugins ***REMOVED***
		p.activateWait.L.Lock()
		if p.activated() && p.implements(iface) ***REMOVED***
			p.handlersRun = false
		***REMOVED***
		p.activateWait.L.Unlock()
	***REMOVED***
	storage.Unlock()

	handlers.Unlock()
***REMOVED***

// GetAll returns all the plugins for the specified implementation
func GetAll(imp string) ([]*Plugin, error) ***REMOVED***
	pluginNames, err := Scan()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	type plLoad struct ***REMOVED***
		pl  *Plugin
		err error
	***REMOVED***

	chPl := make(chan *plLoad, len(pluginNames))
	var wg sync.WaitGroup
	for _, name := range pluginNames ***REMOVED***
		storage.Lock()
		pl, ok := storage.plugins[name]
		storage.Unlock()
		if ok ***REMOVED***
			chPl <- &plLoad***REMOVED***pl, nil***REMOVED***
			continue
		***REMOVED***

		wg.Add(1)
		go func(name string) ***REMOVED***
			defer wg.Done()
			pl, err := loadWithRetry(name, false)
			chPl <- &plLoad***REMOVED***pl, err***REMOVED***
		***REMOVED***(name)
	***REMOVED***

	wg.Wait()
	close(chPl)

	var out []*Plugin
	for pl := range chPl ***REMOVED***
		if pl.err != nil ***REMOVED***
			logrus.Error(pl.err)
			continue
		***REMOVED***
		if err := pl.pl.waitActive(); err == nil && pl.pl.implements(imp) ***REMOVED***
			out = append(out, pl.pl)
		***REMOVED***
	***REMOVED***
	return out, nil
***REMOVED***
