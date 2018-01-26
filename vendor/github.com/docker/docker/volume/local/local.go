// Package local provides the default implementation for volumes. It
// is used to mount data volume containers and directories local to
// the host server.
package local

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/docker/docker/daemon/names"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/volume"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// VolumeDataPathName is the name of the directory where the volume data is stored.
// It uses a very distinctive name to avoid collisions migrating data between
// Docker versions.
const (
	VolumeDataPathName = "_data"
	volumesPathName    = "volumes"
)

var (
	// ErrNotFound is the typed error returned when the requested volume name can't be found
	ErrNotFound = fmt.Errorf("volume not found")
	// volumeNameRegex ensures the name assigned for the volume is valid.
	// This name is used to create the bind directory, so we need to avoid characters that
	// would make the path to escape the root directory.
	volumeNameRegex = names.RestrictedNamePattern
)

type activeMount struct ***REMOVED***
	count   uint64
	mounted bool
***REMOVED***

// New instantiates a new Root instance with the provided scope. Scope
// is the base path that the Root instance uses to store its
// volumes. The base path is created here if it does not exist.
func New(scope string, rootIDs idtools.IDPair) (*Root, error) ***REMOVED***
	rootDirectory := filepath.Join(scope, volumesPathName)

	if err := idtools.MkdirAllAndChown(rootDirectory, 0700, rootIDs); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	r := &Root***REMOVED***
		scope:   scope,
		path:    rootDirectory,
		volumes: make(map[string]*localVolume),
		rootIDs: rootIDs,
	***REMOVED***

	dirs, err := ioutil.ReadDir(rootDirectory)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	mountInfos, err := mount.GetMounts()
	if err != nil ***REMOVED***
		logrus.Debugf("error looking up mounts for local volume cleanup: %v", err)
	***REMOVED***

	for _, d := range dirs ***REMOVED***
		if !d.IsDir() ***REMOVED***
			continue
		***REMOVED***

		name := filepath.Base(d.Name())
		v := &localVolume***REMOVED***
			driverName: r.Name(),
			name:       name,
			path:       r.DataPath(name),
		***REMOVED***
		r.volumes[name] = v
		optsFilePath := filepath.Join(rootDirectory, name, "opts.json")
		if b, err := ioutil.ReadFile(optsFilePath); err == nil ***REMOVED***
			opts := optsConfig***REMOVED******REMOVED***
			if err := json.Unmarshal(b, &opts); err != nil ***REMOVED***
				return nil, errors.Wrapf(err, "error while unmarshaling volume options for volume: %s", name)
			***REMOVED***
			// Make sure this isn't an empty optsConfig.
			// This could be empty due to buggy behavior in older versions of Docker.
			if !reflect.DeepEqual(opts, optsConfig***REMOVED******REMOVED***) ***REMOVED***
				v.opts = &opts
			***REMOVED***

			// unmount anything that may still be mounted (for example, from an unclean shutdown)
			for _, info := range mountInfos ***REMOVED***
				if info.Mountpoint == v.path ***REMOVED***
					mount.Unmount(v.path)
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return r, nil
***REMOVED***

// Root implements the Driver interface for the volume package and
// manages the creation/removal of volumes. It uses only standard vfs
// commands to create/remove dirs within its provided scope.
type Root struct ***REMOVED***
	m       sync.Mutex
	scope   string
	path    string
	volumes map[string]*localVolume
	rootIDs idtools.IDPair
***REMOVED***

// List lists all the volumes
func (r *Root) List() ([]volume.Volume, error) ***REMOVED***
	var ls []volume.Volume
	r.m.Lock()
	for _, v := range r.volumes ***REMOVED***
		ls = append(ls, v)
	***REMOVED***
	r.m.Unlock()
	return ls, nil
***REMOVED***

// DataPath returns the constructed path of this volume.
func (r *Root) DataPath(volumeName string) string ***REMOVED***
	return filepath.Join(r.path, volumeName, VolumeDataPathName)
***REMOVED***

// Name returns the name of Root, defined in the volume package in the DefaultDriverName constant.
func (r *Root) Name() string ***REMOVED***
	return volume.DefaultDriverName
***REMOVED***

// Create creates a new volume.Volume with the provided name, creating
// the underlying directory tree required for this volume in the
// process.
func (r *Root) Create(name string, opts map[string]string) (volume.Volume, error) ***REMOVED***
	if err := r.validateName(name); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	r.m.Lock()
	defer r.m.Unlock()

	v, exists := r.volumes[name]
	if exists ***REMOVED***
		return v, nil
	***REMOVED***

	path := r.DataPath(name)
	if err := idtools.MkdirAllAndChown(path, 0755, r.rootIDs); err != nil ***REMOVED***
		return nil, errors.Wrapf(errdefs.System(err), "error while creating volume path '%s'", path)
	***REMOVED***

	var err error
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			os.RemoveAll(filepath.Dir(path))
		***REMOVED***
	***REMOVED***()

	v = &localVolume***REMOVED***
		driverName: r.Name(),
		name:       name,
		path:       path,
	***REMOVED***

	if len(opts) != 0 ***REMOVED***
		if err = setOpts(v, opts); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var b []byte
		b, err = json.Marshal(v.opts)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if err = ioutil.WriteFile(filepath.Join(filepath.Dir(path), "opts.json"), b, 600); err != nil ***REMOVED***
			return nil, errdefs.System(errors.Wrap(err, "error while persisting volume options"))
		***REMOVED***
	***REMOVED***

	r.volumes[name] = v
	return v, nil
***REMOVED***

// Remove removes the specified volume and all underlying data. If the
// given volume does not belong to this driver and an error is
// returned. The volume is reference counted, if all references are
// not released then the volume is not removed.
func (r *Root) Remove(v volume.Volume) error ***REMOVED***
	r.m.Lock()
	defer r.m.Unlock()

	lv, ok := v.(*localVolume)
	if !ok ***REMOVED***
		return errdefs.System(errors.Errorf("unknown volume type %T", v))
	***REMOVED***

	if lv.active.count > 0 ***REMOVED***
		return errdefs.System(errors.Errorf("volume has active mounts"))
	***REMOVED***

	if err := lv.unmount(); err != nil ***REMOVED***
		return err
	***REMOVED***

	realPath, err := filepath.EvalSymlinks(lv.path)
	if err != nil ***REMOVED***
		if !os.IsNotExist(err) ***REMOVED***
			return err
		***REMOVED***
		realPath = filepath.Dir(lv.path)
	***REMOVED***

	if !r.scopedPath(realPath) ***REMOVED***
		return errdefs.System(errors.Errorf("Unable to remove a directory of out the Docker root %s: %s", r.scope, realPath))
	***REMOVED***

	if err := removePath(realPath); err != nil ***REMOVED***
		return err
	***REMOVED***

	delete(r.volumes, lv.name)
	return removePath(filepath.Dir(lv.path))
***REMOVED***

func removePath(path string) error ***REMOVED***
	if err := os.RemoveAll(path); err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return nil
		***REMOVED***
		return errdefs.System(errors.Wrapf(err, "error removing volume path '%s'", path))
	***REMOVED***
	return nil
***REMOVED***

// Get looks up the volume for the given name and returns it if found
func (r *Root) Get(name string) (volume.Volume, error) ***REMOVED***
	r.m.Lock()
	v, exists := r.volumes[name]
	r.m.Unlock()
	if !exists ***REMOVED***
		return nil, ErrNotFound
	***REMOVED***
	return v, nil
***REMOVED***

// Scope returns the local volume scope
func (r *Root) Scope() string ***REMOVED***
	return volume.LocalScope
***REMOVED***

type validationError string

func (e validationError) Error() string ***REMOVED***
	return string(e)
***REMOVED***

func (e validationError) InvalidParameter() ***REMOVED******REMOVED***

func (r *Root) validateName(name string) error ***REMOVED***
	if len(name) == 1 ***REMOVED***
		return validationError("volume name is too short, names should be at least two alphanumeric characters")
	***REMOVED***
	if !volumeNameRegex.MatchString(name) ***REMOVED***
		return validationError(fmt.Sprintf("%q includes invalid characters for a local volume name, only %q are allowed. If you intended to pass a host directory, use absolute path", name, names.RestrictedNameChars))
	***REMOVED***
	return nil
***REMOVED***

// localVolume implements the Volume interface from the volume package and
// represents the volumes created by Root.
type localVolume struct ***REMOVED***
	m sync.Mutex
	// unique name of the volume
	name string
	// path is the path on the host where the data lives
	path string
	// driverName is the name of the driver that created the volume.
	driverName string
	// opts is the parsed list of options used to create the volume
	opts *optsConfig
	// active refcounts the active mounts
	active activeMount
***REMOVED***

// Name returns the name of the given Volume.
func (v *localVolume) Name() string ***REMOVED***
	return v.name
***REMOVED***

// DriverName returns the driver that created the given Volume.
func (v *localVolume) DriverName() string ***REMOVED***
	return v.driverName
***REMOVED***

// Path returns the data location.
func (v *localVolume) Path() string ***REMOVED***
	return v.path
***REMOVED***

// CachedPath returns the data location
func (v *localVolume) CachedPath() string ***REMOVED***
	return v.path
***REMOVED***

// Mount implements the localVolume interface, returning the data location.
// If there are any provided mount options, the resources will be mounted at this point
func (v *localVolume) Mount(id string) (string, error) ***REMOVED***
	v.m.Lock()
	defer v.m.Unlock()
	if v.opts != nil ***REMOVED***
		if !v.active.mounted ***REMOVED***
			if err := v.mount(); err != nil ***REMOVED***
				return "", errdefs.System(err)
			***REMOVED***
			v.active.mounted = true
		***REMOVED***
		v.active.count++
	***REMOVED***
	return v.path, nil
***REMOVED***

// Unmount dereferences the id, and if it is the last reference will unmount any resources
// that were previously mounted.
func (v *localVolume) Unmount(id string) error ***REMOVED***
	v.m.Lock()
	defer v.m.Unlock()

	// Always decrement the count, even if the unmount fails
	// Essentially docker doesn't care if this fails, it will send an error, but
	// ultimately there's nothing that can be done. If we don't decrement the count
	// this volume can never be removed until a daemon restart occurs.
	if v.opts != nil ***REMOVED***
		v.active.count--
	***REMOVED***

	if v.active.count > 0 ***REMOVED***
		return nil
	***REMOVED***

	return v.unmount()
***REMOVED***

func (v *localVolume) unmount() error ***REMOVED***
	if v.opts != nil ***REMOVED***
		if err := mount.Unmount(v.path); err != nil ***REMOVED***
			if mounted, mErr := mount.Mounted(v.path); mounted || mErr != nil ***REMOVED***
				return errdefs.System(errors.Wrapf(err, "error while unmounting volume path '%s'", v.path))
			***REMOVED***
		***REMOVED***
		v.active.mounted = false
	***REMOVED***
	return nil
***REMOVED***

func validateOpts(opts map[string]string) error ***REMOVED***
	for opt := range opts ***REMOVED***
		if !validOpts[opt] ***REMOVED***
			return validationError(fmt.Sprintf("invalid option key: %q", opt))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (v *localVolume) Status() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***

// getAddress finds out address/hostname from options
func getAddress(opts string) string ***REMOVED***
	optsList := strings.Split(opts, ",")
	for i := 0; i < len(optsList); i++ ***REMOVED***
		if strings.HasPrefix(optsList[i], "addr=") ***REMOVED***
			addr := (strings.SplitN(optsList[i], "=", 2)[1])
			return addr
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***
