package volumedrivers

import (
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/volume"
	"github.com/sirupsen/logrus"
)

var (
	errNoSuchVolume = errors.New("no such volume")
)

type volumeDriverAdapter struct ***REMOVED***
	name         string
	baseHostPath string
	capabilities *volume.Capability
	proxy        *volumeDriverProxy
***REMOVED***

func (a *volumeDriverAdapter) Name() string ***REMOVED***
	return a.name
***REMOVED***

func (a *volumeDriverAdapter) Create(name string, opts map[string]string) (volume.Volume, error) ***REMOVED***
	if err := a.proxy.Create(name, opts); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &volumeAdapter***REMOVED***
		proxy:        a.proxy,
		name:         name,
		driverName:   a.name,
		baseHostPath: a.baseHostPath,
	***REMOVED***, nil
***REMOVED***

func (a *volumeDriverAdapter) Remove(v volume.Volume) error ***REMOVED***
	return a.proxy.Remove(v.Name())
***REMOVED***

func hostPath(baseHostPath, path string) string ***REMOVED***
	if baseHostPath != "" ***REMOVED***
		path = filepath.Join(baseHostPath, path)
	***REMOVED***
	return path
***REMOVED***

func (a *volumeDriverAdapter) List() ([]volume.Volume, error) ***REMOVED***
	ls, err := a.proxy.List()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var out []volume.Volume
	for _, vp := range ls ***REMOVED***
		out = append(out, &volumeAdapter***REMOVED***
			proxy:        a.proxy,
			name:         vp.Name,
			baseHostPath: a.baseHostPath,
			driverName:   a.name,
			eMount:       hostPath(a.baseHostPath, vp.Mountpoint),
		***REMOVED***)
	***REMOVED***
	return out, nil
***REMOVED***

func (a *volumeDriverAdapter) Get(name string) (volume.Volume, error) ***REMOVED***
	v, err := a.proxy.Get(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// plugin may have returned no volume and no error
	if v == nil ***REMOVED***
		return nil, errNoSuchVolume
	***REMOVED***

	return &volumeAdapter***REMOVED***
		proxy:        a.proxy,
		name:         v.Name,
		driverName:   a.Name(),
		eMount:       v.Mountpoint,
		createdAt:    v.CreatedAt,
		status:       v.Status,
		baseHostPath: a.baseHostPath,
	***REMOVED***, nil
***REMOVED***

func (a *volumeDriverAdapter) Scope() string ***REMOVED***
	cap := a.getCapabilities()
	return cap.Scope
***REMOVED***

func (a *volumeDriverAdapter) getCapabilities() volume.Capability ***REMOVED***
	if a.capabilities != nil ***REMOVED***
		return *a.capabilities
	***REMOVED***
	cap, err := a.proxy.Capabilities()
	if err != nil ***REMOVED***
		// `GetCapabilities` is a not a required endpoint.
		// On error assume it's a local-only driver
		logrus.Warnf("Volume driver %s returned an error while trying to query its capabilities, using default capabilities: %v", a.name, err)
		return volume.Capability***REMOVED***Scope: volume.LocalScope***REMOVED***
	***REMOVED***

	// don't spam the warn log below just because the plugin didn't provide a scope
	if len(cap.Scope) == 0 ***REMOVED***
		cap.Scope = volume.LocalScope
	***REMOVED***

	cap.Scope = strings.ToLower(cap.Scope)
	if cap.Scope != volume.LocalScope && cap.Scope != volume.GlobalScope ***REMOVED***
		logrus.Warnf("Volume driver %q returned an invalid scope: %q", a.Name(), cap.Scope)
		cap.Scope = volume.LocalScope
	***REMOVED***

	a.capabilities = &cap
	return cap
***REMOVED***

type volumeAdapter struct ***REMOVED***
	proxy        *volumeDriverProxy
	name         string
	baseHostPath string
	driverName   string
	eMount       string    // ephemeral host volume path
	createdAt    time.Time // time the directory was created
	status       map[string]interface***REMOVED******REMOVED***
***REMOVED***

type proxyVolume struct ***REMOVED***
	Name       string
	Mountpoint string
	CreatedAt  time.Time
	Status     map[string]interface***REMOVED******REMOVED***
***REMOVED***

func (a *volumeAdapter) Name() string ***REMOVED***
	return a.name
***REMOVED***

func (a *volumeAdapter) DriverName() string ***REMOVED***
	return a.driverName
***REMOVED***

func (a *volumeAdapter) Path() string ***REMOVED***
	if len(a.eMount) == 0 ***REMOVED***
		mountpoint, _ := a.proxy.Path(a.name)
		a.eMount = hostPath(a.baseHostPath, mountpoint)
	***REMOVED***
	return a.eMount
***REMOVED***

func (a *volumeAdapter) CachedPath() string ***REMOVED***
	return a.eMount
***REMOVED***

func (a *volumeAdapter) Mount(id string) (string, error) ***REMOVED***
	mountpoint, err := a.proxy.Mount(a.name, id)
	a.eMount = hostPath(a.baseHostPath, mountpoint)
	return a.eMount, err
***REMOVED***

func (a *volumeAdapter) Unmount(id string) error ***REMOVED***
	err := a.proxy.Unmount(a.name, id)
	if err == nil ***REMOVED***
		a.eMount = ""
	***REMOVED***
	return err
***REMOVED***

func (a *volumeAdapter) CreatedAt() (time.Time, error) ***REMOVED***
	return a.createdAt, nil
***REMOVED***
func (a *volumeAdapter) Status() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	out := make(map[string]interface***REMOVED******REMOVED***, len(a.status))
	for k, v := range a.status ***REMOVED***
		out[k] = v
	***REMOVED***
	return out
***REMOVED***
