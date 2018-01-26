// generated code - DO NOT EDIT

package volumedrivers

import (
	"errors"
	"time"

	"github.com/docker/docker/pkg/plugins"
	"github.com/docker/docker/volume"
)

const (
	longTimeout  = 2 * time.Minute
	shortTimeout = 1 * time.Minute
)

type client interface ***REMOVED***
	CallWithOptions(string, interface***REMOVED******REMOVED***, interface***REMOVED******REMOVED***, ...func(*plugins.RequestOpts)) error
***REMOVED***

type volumeDriverProxy struct ***REMOVED***
	client
***REMOVED***

type volumeDriverProxyCreateRequest struct ***REMOVED***
	Name string
	Opts map[string]string
***REMOVED***

type volumeDriverProxyCreateResponse struct ***REMOVED***
	Err string
***REMOVED***

func (pp *volumeDriverProxy) Create(name string, opts map[string]string) (err error) ***REMOVED***
	var (
		req volumeDriverProxyCreateRequest
		ret volumeDriverProxyCreateResponse
	)

	req.Name = name
	req.Opts = opts

	if err = pp.CallWithOptions("VolumeDriver.Create", req, &ret, plugins.WithRequestTimeout(longTimeout)); err != nil ***REMOVED***
		return
	***REMOVED***

	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***

	return
***REMOVED***

type volumeDriverProxyRemoveRequest struct ***REMOVED***
	Name string
***REMOVED***

type volumeDriverProxyRemoveResponse struct ***REMOVED***
	Err string
***REMOVED***

func (pp *volumeDriverProxy) Remove(name string) (err error) ***REMOVED***
	var (
		req volumeDriverProxyRemoveRequest
		ret volumeDriverProxyRemoveResponse
	)

	req.Name = name

	if err = pp.CallWithOptions("VolumeDriver.Remove", req, &ret, plugins.WithRequestTimeout(shortTimeout)); err != nil ***REMOVED***
		return
	***REMOVED***

	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***

	return
***REMOVED***

type volumeDriverProxyPathRequest struct ***REMOVED***
	Name string
***REMOVED***

type volumeDriverProxyPathResponse struct ***REMOVED***
	Mountpoint string
	Err        string
***REMOVED***

func (pp *volumeDriverProxy) Path(name string) (mountpoint string, err error) ***REMOVED***
	var (
		req volumeDriverProxyPathRequest
		ret volumeDriverProxyPathResponse
	)

	req.Name = name

	if err = pp.CallWithOptions("VolumeDriver.Path", req, &ret, plugins.WithRequestTimeout(shortTimeout)); err != nil ***REMOVED***
		return
	***REMOVED***

	mountpoint = ret.Mountpoint

	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***

	return
***REMOVED***

type volumeDriverProxyMountRequest struct ***REMOVED***
	Name string
	ID   string
***REMOVED***

type volumeDriverProxyMountResponse struct ***REMOVED***
	Mountpoint string
	Err        string
***REMOVED***

func (pp *volumeDriverProxy) Mount(name string, id string) (mountpoint string, err error) ***REMOVED***
	var (
		req volumeDriverProxyMountRequest
		ret volumeDriverProxyMountResponse
	)

	req.Name = name
	req.ID = id

	if err = pp.CallWithOptions("VolumeDriver.Mount", req, &ret, plugins.WithRequestTimeout(longTimeout)); err != nil ***REMOVED***
		return
	***REMOVED***

	mountpoint = ret.Mountpoint

	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***

	return
***REMOVED***

type volumeDriverProxyUnmountRequest struct ***REMOVED***
	Name string
	ID   string
***REMOVED***

type volumeDriverProxyUnmountResponse struct ***REMOVED***
	Err string
***REMOVED***

func (pp *volumeDriverProxy) Unmount(name string, id string) (err error) ***REMOVED***
	var (
		req volumeDriverProxyUnmountRequest
		ret volumeDriverProxyUnmountResponse
	)

	req.Name = name
	req.ID = id

	if err = pp.CallWithOptions("VolumeDriver.Unmount", req, &ret, plugins.WithRequestTimeout(shortTimeout)); err != nil ***REMOVED***
		return
	***REMOVED***

	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***

	return
***REMOVED***

type volumeDriverProxyListRequest struct ***REMOVED***
***REMOVED***

type volumeDriverProxyListResponse struct ***REMOVED***
	Volumes []*proxyVolume
	Err     string
***REMOVED***

func (pp *volumeDriverProxy) List() (volumes []*proxyVolume, err error) ***REMOVED***
	var (
		req volumeDriverProxyListRequest
		ret volumeDriverProxyListResponse
	)

	if err = pp.CallWithOptions("VolumeDriver.List", req, &ret, plugins.WithRequestTimeout(shortTimeout)); err != nil ***REMOVED***
		return
	***REMOVED***

	volumes = ret.Volumes

	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***

	return
***REMOVED***

type volumeDriverProxyGetRequest struct ***REMOVED***
	Name string
***REMOVED***

type volumeDriverProxyGetResponse struct ***REMOVED***
	Volume *proxyVolume
	Err    string
***REMOVED***

func (pp *volumeDriverProxy) Get(name string) (volume *proxyVolume, err error) ***REMOVED***
	var (
		req volumeDriverProxyGetRequest
		ret volumeDriverProxyGetResponse
	)

	req.Name = name

	if err = pp.CallWithOptions("VolumeDriver.Get", req, &ret, plugins.WithRequestTimeout(shortTimeout)); err != nil ***REMOVED***
		return
	***REMOVED***

	volume = ret.Volume

	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***

	return
***REMOVED***

type volumeDriverProxyCapabilitiesRequest struct ***REMOVED***
***REMOVED***

type volumeDriverProxyCapabilitiesResponse struct ***REMOVED***
	Capabilities volume.Capability
	Err          string
***REMOVED***

func (pp *volumeDriverProxy) Capabilities() (capabilities volume.Capability, err error) ***REMOVED***
	var (
		req volumeDriverProxyCapabilitiesRequest
		ret volumeDriverProxyCapabilitiesResponse
	)

	if err = pp.CallWithOptions("VolumeDriver.Capabilities", req, &ret, plugins.WithRequestTimeout(shortTimeout)); err != nil ***REMOVED***
		return
	***REMOVED***

	capabilities = ret.Capabilities

	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***

	return
***REMOVED***
