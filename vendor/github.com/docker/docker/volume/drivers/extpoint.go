//go:generate pluginrpc-gen -i $GOFILE -o proxy.go -type volumeDriver -name VolumeDriver

package volumedrivers

import (
	"fmt"
	"sort"
	"sync"

	"github.com/docker/docker/pkg/locker"
	getter "github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/volume"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// currently created by hand. generation tool would generate this like:
// $ extpoint-gen Driver > volume/extpoint.go

var drivers = &driverExtpoint***REMOVED***
	extensions: make(map[string]volume.Driver),
	driverLock: &locker.Locker***REMOVED******REMOVED***,
***REMOVED***

const extName = "VolumeDriver"

// NewVolumeDriver returns a driver has the given name mapped on the given client.
func NewVolumeDriver(name string, baseHostPath string, c client) volume.Driver ***REMOVED***
	proxy := &volumeDriverProxy***REMOVED***c***REMOVED***
	return &volumeDriverAdapter***REMOVED***name: name, baseHostPath: baseHostPath, proxy: proxy***REMOVED***
***REMOVED***

// volumeDriver defines the available functions that volume plugins must implement.
// This interface is only defined to generate the proxy objects.
// It's not intended to be public or reused.
// nolint: deadcode
type volumeDriver interface ***REMOVED***
	// Create a volume with the given name
	Create(name string, opts map[string]string) (err error)
	// Remove the volume with the given name
	Remove(name string) (err error)
	// Get the mountpoint of the given volume
	Path(name string) (mountpoint string, err error)
	// Mount the given volume and return the mountpoint
	Mount(name, id string) (mountpoint string, err error)
	// Unmount the given volume
	Unmount(name, id string) (err error)
	// List lists all the volumes known to the driver
	List() (volumes []*proxyVolume, err error)
	// Get retrieves the volume with the requested name
	Get(name string) (volume *proxyVolume, err error)
	// Capabilities gets the list of capabilities of the driver
	Capabilities() (capabilities volume.Capability, err error)
***REMOVED***

type driverExtpoint struct ***REMOVED***
	extensions map[string]volume.Driver
	sync.Mutex
	driverLock   *locker.Locker
	plugingetter getter.PluginGetter
***REMOVED***

// RegisterPluginGetter sets the plugingetter
func RegisterPluginGetter(plugingetter getter.PluginGetter) ***REMOVED***
	drivers.plugingetter = plugingetter
***REMOVED***

// Register associates the given driver to the given name, checking if
// the name is already associated
func Register(extension volume.Driver, name string) bool ***REMOVED***
	if name == "" ***REMOVED***
		return false
	***REMOVED***

	drivers.Lock()
	defer drivers.Unlock()

	_, exists := drivers.extensions[name]
	if exists ***REMOVED***
		return false
	***REMOVED***

	if err := validateDriver(extension); err != nil ***REMOVED***
		return false
	***REMOVED***

	drivers.extensions[name] = extension

	return true
***REMOVED***

// Unregister dissociates the name from its driver, if the association exists.
func Unregister(name string) bool ***REMOVED***
	drivers.Lock()
	defer drivers.Unlock()

	_, exists := drivers.extensions[name]
	if !exists ***REMOVED***
		return false
	***REMOVED***
	delete(drivers.extensions, name)
	return true
***REMOVED***

type driverNotFoundError string

func (e driverNotFoundError) Error() string ***REMOVED***
	return "volume driver not found: " + string(e)
***REMOVED***

func (driverNotFoundError) NotFound() ***REMOVED******REMOVED***

// lookup returns the driver associated with the given name. If a
// driver with the given name has not been registered it checks if
// there is a VolumeDriver plugin available with the given name.
func lookup(name string, mode int) (volume.Driver, error) ***REMOVED***
	drivers.driverLock.Lock(name)
	defer drivers.driverLock.Unlock(name)

	drivers.Lock()
	ext, ok := drivers.extensions[name]
	drivers.Unlock()
	if ok ***REMOVED***
		return ext, nil
	***REMOVED***
	if drivers.plugingetter != nil ***REMOVED***
		p, err := drivers.plugingetter.Get(name, extName, mode)
		if err != nil ***REMOVED***
			return nil, errors.Wrap(err, "error looking up volume plugin "+name)
		***REMOVED***

		d := NewVolumeDriver(p.Name(), p.BasePath(), p.Client())
		if err := validateDriver(d); err != nil ***REMOVED***
			if mode > 0 ***REMOVED***
				// Undo any reference count changes from the initial `Get`
				if _, err := drivers.plugingetter.Get(name, extName, mode*-1); err != nil ***REMOVED***
					logrus.WithError(err).WithField("action", "validate-driver").WithField("plugin", name).Error("error releasing reference to plugin")
				***REMOVED***
			***REMOVED***
			return nil, err
		***REMOVED***

		if p.IsV1() ***REMOVED***
			drivers.Lock()
			drivers.extensions[name] = d
			drivers.Unlock()
		***REMOVED***
		return d, nil
	***REMOVED***
	return nil, driverNotFoundError(name)
***REMOVED***

func validateDriver(vd volume.Driver) error ***REMOVED***
	scope := vd.Scope()
	if scope != volume.LocalScope && scope != volume.GlobalScope ***REMOVED***
		return fmt.Errorf("Driver %q provided an invalid capability scope: %s", vd.Name(), scope)
	***REMOVED***
	return nil
***REMOVED***

// GetDriver returns a volume driver by its name.
// If the driver is empty, it looks for the local driver.
func GetDriver(name string) (volume.Driver, error) ***REMOVED***
	if name == "" ***REMOVED***
		name = volume.DefaultDriverName
	***REMOVED***
	return lookup(name, getter.Lookup)
***REMOVED***

// CreateDriver returns a volume driver by its name and increments RefCount.
// If the driver is empty, it looks for the local driver.
func CreateDriver(name string) (volume.Driver, error) ***REMOVED***
	if name == "" ***REMOVED***
		name = volume.DefaultDriverName
	***REMOVED***
	return lookup(name, getter.Acquire)
***REMOVED***

// ReleaseDriver returns a volume driver by its name and decrements RefCount..
// If the driver is empty, it looks for the local driver.
func ReleaseDriver(name string) (volume.Driver, error) ***REMOVED***
	if name == "" ***REMOVED***
		name = volume.DefaultDriverName
	***REMOVED***
	return lookup(name, getter.Release)
***REMOVED***

// GetDriverList returns list of volume drivers registered.
// If no driver is registered, empty string list will be returned.
func GetDriverList() []string ***REMOVED***
	var driverList []string
	drivers.Lock()
	for driverName := range drivers.extensions ***REMOVED***
		driverList = append(driverList, driverName)
	***REMOVED***
	drivers.Unlock()
	sort.Strings(driverList)
	return driverList
***REMOVED***

// GetAllDrivers lists all the registered drivers
func GetAllDrivers() ([]volume.Driver, error) ***REMOVED***
	var plugins []getter.CompatPlugin
	if drivers.plugingetter != nil ***REMOVED***
		var err error
		plugins, err = drivers.plugingetter.GetAllByCap(extName)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("error listing plugins: %v", err)
		***REMOVED***
	***REMOVED***
	var ds []volume.Driver

	drivers.Lock()
	defer drivers.Unlock()

	for _, d := range drivers.extensions ***REMOVED***
		ds = append(ds, d)
	***REMOVED***

	for _, p := range plugins ***REMOVED***
		name := p.Name()

		if _, ok := drivers.extensions[name]; ok ***REMOVED***
			continue
		***REMOVED***

		ext := NewVolumeDriver(name, p.BasePath(), p.Client())
		if p.IsV1() ***REMOVED***
			drivers.extensions[name] = ext
		***REMOVED***
		ds = append(ds, ext)
	***REMOVED***
	return ds, nil
***REMOVED***
