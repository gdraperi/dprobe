package testutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/pkg/plugins"
	"github.com/docker/docker/volume"
)

// NoopVolume is a volume that doesn't perform any operation
type NoopVolume struct***REMOVED******REMOVED***

// Name is the name of the volume
func (NoopVolume) Name() string ***REMOVED*** return "noop" ***REMOVED***

// DriverName is the name of the driver
func (NoopVolume) DriverName() string ***REMOVED*** return "noop" ***REMOVED***

// Path is the filesystem path to the volume
func (NoopVolume) Path() string ***REMOVED*** return "noop" ***REMOVED***

// Mount mounts the volume in the container
func (NoopVolume) Mount(_ string) (string, error) ***REMOVED*** return "noop", nil ***REMOVED***

// Unmount unmounts the volume from the container
func (NoopVolume) Unmount(_ string) error ***REMOVED*** return nil ***REMOVED***

// Status provides low-level details about the volume
func (NoopVolume) Status() map[string]interface***REMOVED******REMOVED*** ***REMOVED*** return nil ***REMOVED***

// CreatedAt provides the time the volume (directory) was created at
func (NoopVolume) CreatedAt() (time.Time, error) ***REMOVED*** return time.Now(), nil ***REMOVED***

// FakeVolume is a fake volume with a random name
type FakeVolume struct ***REMOVED***
	name       string
	driverName string
***REMOVED***

// NewFakeVolume creates a new fake volume for testing
func NewFakeVolume(name string, driverName string) volume.Volume ***REMOVED***
	return FakeVolume***REMOVED***name: name, driverName: driverName***REMOVED***
***REMOVED***

// Name is the name of the volume
func (f FakeVolume) Name() string ***REMOVED*** return f.name ***REMOVED***

// DriverName is the name of the driver
func (f FakeVolume) DriverName() string ***REMOVED*** return f.driverName ***REMOVED***

// Path is the filesystem path to the volume
func (FakeVolume) Path() string ***REMOVED*** return "fake" ***REMOVED***

// Mount mounts the volume in the container
func (FakeVolume) Mount(_ string) (string, error) ***REMOVED*** return "fake", nil ***REMOVED***

// Unmount unmounts the volume from the container
func (FakeVolume) Unmount(_ string) error ***REMOVED*** return nil ***REMOVED***

// Status provides low-level details about the volume
func (FakeVolume) Status() map[string]interface***REMOVED******REMOVED*** ***REMOVED*** return nil ***REMOVED***

// CreatedAt provides the time the volume (directory) was created at
func (FakeVolume) CreatedAt() (time.Time, error) ***REMOVED*** return time.Now(), nil ***REMOVED***

// FakeDriver is a driver that generates fake volumes
type FakeDriver struct ***REMOVED***
	name string
	vols map[string]volume.Volume
***REMOVED***

// NewFakeDriver creates a new FakeDriver with the specified name
func NewFakeDriver(name string) volume.Driver ***REMOVED***
	return &FakeDriver***REMOVED***
		name: name,
		vols: make(map[string]volume.Volume),
	***REMOVED***
***REMOVED***

// Name is the name of the driver
func (d *FakeDriver) Name() string ***REMOVED*** return d.name ***REMOVED***

// Create initializes a fake volume.
// It returns an error if the options include an "error" key with a message
func (d *FakeDriver) Create(name string, opts map[string]string) (volume.Volume, error) ***REMOVED***
	if opts != nil && opts["error"] != "" ***REMOVED***
		return nil, fmt.Errorf(opts["error"])
	***REMOVED***
	v := NewFakeVolume(name, d.name)
	d.vols[name] = v
	return v, nil
***REMOVED***

// Remove deletes a volume.
func (d *FakeDriver) Remove(v volume.Volume) error ***REMOVED***
	if _, exists := d.vols[v.Name()]; !exists ***REMOVED***
		return fmt.Errorf("no such volume")
	***REMOVED***
	delete(d.vols, v.Name())
	return nil
***REMOVED***

// List lists the volumes
func (d *FakeDriver) List() ([]volume.Volume, error) ***REMOVED***
	var vols []volume.Volume
	for _, v := range d.vols ***REMOVED***
		vols = append(vols, v)
	***REMOVED***
	return vols, nil
***REMOVED***

// Get gets the volume
func (d *FakeDriver) Get(name string) (volume.Volume, error) ***REMOVED***
	if v, exists := d.vols[name]; exists ***REMOVED***
		return v, nil
	***REMOVED***
	return nil, fmt.Errorf("no such volume")
***REMOVED***

// Scope returns the local scope
func (*FakeDriver) Scope() string ***REMOVED***
	return "local"
***REMOVED***

type fakePlugin struct ***REMOVED***
	client *plugins.Client
	name   string
	refs   int
***REMOVED***

// MakeFakePlugin creates a fake plugin from the passed in driver
// Note: currently only "Create" is implemented because that's all that's needed
// so far. If you need it to test something else, add it here, but probably you
// shouldn't need to use this except for very specific cases with v2 plugin handling.
func MakeFakePlugin(d volume.Driver, l net.Listener) (plugingetter.CompatPlugin, error) ***REMOVED***
	c, err := plugins.NewClient(l.Addr().Network()+"://"+l.Addr().String(), nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	mux := http.NewServeMux()

	mux.HandleFunc("/VolumeDriver.Create", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		createReq := struct ***REMOVED***
			Name string
			Opts map[string]string
		***REMOVED******REMOVED******REMOVED***
		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil ***REMOVED***
			fmt.Fprintf(w, `***REMOVED***"Err": "%s"***REMOVED***`, err.Error())
			return
		***REMOVED***
		_, err := d.Create(createReq.Name, createReq.Opts)
		if err != nil ***REMOVED***
			fmt.Fprintf(w, `***REMOVED***"Err": "%s"***REMOVED***`, err.Error())
			return
		***REMOVED***
		w.Write([]byte("***REMOVED******REMOVED***"))
	***REMOVED***)

	go http.Serve(l, mux)
	return &fakePlugin***REMOVED***client: c, name: d.Name()***REMOVED***, nil
***REMOVED***

func (p *fakePlugin) Client() *plugins.Client ***REMOVED***
	return p.client
***REMOVED***

func (p *fakePlugin) Name() string ***REMOVED***
	return p.name
***REMOVED***

func (p *fakePlugin) IsV1() bool ***REMOVED***
	return false
***REMOVED***

func (p *fakePlugin) BasePath() string ***REMOVED***
	return ""
***REMOVED***

type fakePluginGetter struct ***REMOVED***
	plugins map[string]plugingetter.CompatPlugin
***REMOVED***

// NewFakePluginGetter returns a plugin getter for fake plugins
func NewFakePluginGetter(pls ...plugingetter.CompatPlugin) plugingetter.PluginGetter ***REMOVED***
	idx := make(map[string]plugingetter.CompatPlugin, len(pls))
	for _, p := range pls ***REMOVED***
		idx[p.Name()] = p
	***REMOVED***
	return &fakePluginGetter***REMOVED***plugins: idx***REMOVED***
***REMOVED***

// This ignores the second argument since we only care about volume drivers here,
// there shouldn't be any other kind of plugin in here
func (g *fakePluginGetter) Get(name, _ string, mode int) (plugingetter.CompatPlugin, error) ***REMOVED***
	p, ok := g.plugins[name]
	if !ok ***REMOVED***
		return nil, errors.New("not found")
	***REMOVED***
	p.(*fakePlugin).refs += mode
	return p, nil
***REMOVED***

func (g *fakePluginGetter) GetAllByCap(capability string) ([]plugingetter.CompatPlugin, error) ***REMOVED***
	panic("GetAllByCap shouldn't be called")
***REMOVED***

func (g *fakePluginGetter) GetAllManagedPluginsByCap(capability string) []plugingetter.CompatPlugin ***REMOVED***
	panic("GetAllManagedPluginsByCap should not be called")
***REMOVED***

func (g *fakePluginGetter) Handle(capability string, callback func(string, *plugins.Client)) ***REMOVED***
	panic("Handle should not be called")
***REMOVED***

// FakeRefs checks ref count on a fake plugin.
func FakeRefs(p plugingetter.CompatPlugin) int ***REMOVED***
	// this should panic if something other than a `*fakePlugin` is passed in
	return p.(*fakePlugin).refs
***REMOVED***
