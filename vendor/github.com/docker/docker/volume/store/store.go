package store

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/boltdb/bolt"
	"github.com/docker/docker/pkg/locker"
	"github.com/docker/docker/volume"
	"github.com/docker/docker/volume/drivers"
	"github.com/sirupsen/logrus"
)

const (
	volumeDataDir = "volumes"
)

type volumeWrapper struct ***REMOVED***
	volume.Volume
	labels  map[string]string
	scope   string
	options map[string]string
***REMOVED***

func (v volumeWrapper) Options() map[string]string ***REMOVED***
	options := map[string]string***REMOVED******REMOVED***
	for key, value := range v.options ***REMOVED***
		options[key] = value
	***REMOVED***
	return options
***REMOVED***

func (v volumeWrapper) Labels() map[string]string ***REMOVED***
	return v.labels
***REMOVED***

func (v volumeWrapper) Scope() string ***REMOVED***
	return v.scope
***REMOVED***

func (v volumeWrapper) CachedPath() string ***REMOVED***
	if vv, ok := v.Volume.(interface ***REMOVED***
		CachedPath() string
	***REMOVED***); ok ***REMOVED***
		return vv.CachedPath()
	***REMOVED***
	return v.Volume.Path()
***REMOVED***

// New initializes a VolumeStore to keep
// reference counting of volumes in the system.
func New(rootPath string) (*VolumeStore, error) ***REMOVED***
	vs := &VolumeStore***REMOVED***
		locks:   &locker.Locker***REMOVED******REMOVED***,
		names:   make(map[string]volume.Volume),
		refs:    make(map[string]map[string]struct***REMOVED******REMOVED***),
		labels:  make(map[string]map[string]string),
		options: make(map[string]map[string]string),
	***REMOVED***

	if rootPath != "" ***REMOVED***
		// initialize metadata store
		volPath := filepath.Join(rootPath, volumeDataDir)
		if err := os.MkdirAll(volPath, 750); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		dbPath := filepath.Join(volPath, "metadata.db")

		var err error
		vs.db, err = bolt.Open(dbPath, 0600, &bolt.Options***REMOVED***Timeout: 1 * time.Second***REMOVED***)
		if err != nil ***REMOVED***
			return nil, errors.Wrap(err, "error while opening volume store metadata database")
		***REMOVED***

		// initialize volumes bucket
		if err := vs.db.Update(func(tx *bolt.Tx) error ***REMOVED***
			if _, err := tx.CreateBucketIfNotExists(volumeBucketName); err != nil ***REMOVED***
				return errors.Wrap(err, "error while setting up volume store metadata database")
			***REMOVED***
			return nil
		***REMOVED***); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	vs.restore()

	return vs, nil
***REMOVED***

func (s *VolumeStore) getNamed(name string) (volume.Volume, bool) ***REMOVED***
	s.globalLock.RLock()
	v, exists := s.names[name]
	s.globalLock.RUnlock()
	return v, exists
***REMOVED***

func (s *VolumeStore) setNamed(v volume.Volume, ref string) ***REMOVED***
	name := v.Name()

	s.globalLock.Lock()
	s.names[name] = v
	if len(ref) > 0 ***REMOVED***
		if s.refs[name] == nil ***REMOVED***
			s.refs[name] = make(map[string]struct***REMOVED******REMOVED***)
		***REMOVED***
		s.refs[name][ref] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	s.globalLock.Unlock()
***REMOVED***

// hasRef returns true if the given name has at least one ref.
// Callers of this function are expected to hold the name lock.
func (s *VolumeStore) hasRef(name string) bool ***REMOVED***
	s.globalLock.RLock()
	l := len(s.refs[name])
	s.globalLock.RUnlock()
	return l > 0
***REMOVED***

// getRefs gets the list of refs for a given name
// Callers of this function are expected to hold the name lock.
func (s *VolumeStore) getRefs(name string) []string ***REMOVED***
	s.globalLock.RLock()
	defer s.globalLock.RUnlock()

	refs := make([]string, 0, len(s.refs[name]))
	for r := range s.refs[name] ***REMOVED***
		refs = append(refs, r)
	***REMOVED***

	return refs
***REMOVED***

// Purge allows the cleanup of internal data on docker in case
// the internal data is out of sync with volumes driver plugins.
func (s *VolumeStore) Purge(name string) ***REMOVED***
	s.globalLock.Lock()
	v, exists := s.names[name]
	if exists ***REMOVED***
		driverName := v.DriverName()
		if _, err := volumedrivers.ReleaseDriver(driverName); err != nil ***REMOVED***
			logrus.WithError(err).WithField("driver", driverName).Error("Error releasing reference to volume driver")
		***REMOVED***
	***REMOVED***
	if err := s.removeMeta(name); err != nil ***REMOVED***
		logrus.Errorf("Error removing volume metadata for volume %q: %v", name, err)
	***REMOVED***
	delete(s.names, name)
	delete(s.refs, name)
	delete(s.labels, name)
	delete(s.options, name)
	s.globalLock.Unlock()
***REMOVED***

// VolumeStore is a struct that stores the list of volumes available and keeps track of their usage counts
type VolumeStore struct ***REMOVED***
	// locks ensures that only one action is being performed on a particular volume at a time without locking the entire store
	// since actions on volumes can be quite slow, this ensures the store is free to handle requests for other volumes.
	locks *locker.Locker
	// globalLock is used to protect access to mutable structures used by the store object
	globalLock sync.RWMutex
	// names stores the volume name -> volume relationship.
	// This is used for making lookups faster so we don't have to probe all drivers
	names map[string]volume.Volume
	// refs stores the volume name and the list of things referencing it
	refs map[string]map[string]struct***REMOVED******REMOVED***
	// labels stores volume labels for each volume
	labels map[string]map[string]string
	// options stores volume options for each volume
	options map[string]map[string]string
	db      *bolt.DB
***REMOVED***

// List proxies to all registered volume drivers to get the full list of volumes
// If a driver returns a volume that has name which conflicts with another volume from a different driver,
// the first volume is chosen and the conflicting volume is dropped.
func (s *VolumeStore) List() ([]volume.Volume, []string, error) ***REMOVED***
	vols, warnings, err := s.list()
	if err != nil ***REMOVED***
		return nil, nil, &OpErr***REMOVED***Err: err, Op: "list"***REMOVED***
	***REMOVED***
	var out []volume.Volume

	for _, v := range vols ***REMOVED***
		name := normalizeVolumeName(v.Name())

		s.locks.Lock(name)
		storedV, exists := s.getNamed(name)
		// Note: it's not safe to populate the cache here because the volume may have been
		// deleted before we acquire a lock on its name
		if exists && storedV.DriverName() != v.DriverName() ***REMOVED***
			logrus.Warnf("Volume name %s already exists for driver %s, not including volume returned by %s", v.Name(), storedV.DriverName(), v.DriverName())
			s.locks.Unlock(v.Name())
			continue
		***REMOVED***

		out = append(out, v)
		s.locks.Unlock(v.Name())
	***REMOVED***
	return out, warnings, nil
***REMOVED***

// list goes through each volume driver and asks for its list of volumes.
func (s *VolumeStore) list() ([]volume.Volume, []string, error) ***REMOVED***
	var (
		ls       []volume.Volume
		warnings []string
	)

	drivers, err := volumedrivers.GetAllDrivers()
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	type vols struct ***REMOVED***
		vols       []volume.Volume
		err        error
		driverName string
	***REMOVED***
	chVols := make(chan vols, len(drivers))

	for _, vd := range drivers ***REMOVED***
		go func(d volume.Driver) ***REMOVED***
			vs, err := d.List()
			if err != nil ***REMOVED***
				chVols <- vols***REMOVED***driverName: d.Name(), err: &OpErr***REMOVED***Err: err, Name: d.Name(), Op: "list"***REMOVED******REMOVED***
				return
			***REMOVED***
			for i, v := range vs ***REMOVED***
				s.globalLock.RLock()
				vs[i] = volumeWrapper***REMOVED***v, s.labels[v.Name()], d.Scope(), s.options[v.Name()]***REMOVED***
				s.globalLock.RUnlock()
			***REMOVED***

			chVols <- vols***REMOVED***vols: vs***REMOVED***
		***REMOVED***(vd)
	***REMOVED***

	badDrivers := make(map[string]struct***REMOVED******REMOVED***)
	for i := 0; i < len(drivers); i++ ***REMOVED***
		vs := <-chVols

		if vs.err != nil ***REMOVED***
			warnings = append(warnings, vs.err.Error())
			badDrivers[vs.driverName] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			logrus.Warn(vs.err)
		***REMOVED***
		ls = append(ls, vs.vols...)
	***REMOVED***

	if len(badDrivers) > 0 ***REMOVED***
		s.globalLock.RLock()
		for _, v := range s.names ***REMOVED***
			if _, exists := badDrivers[v.DriverName()]; exists ***REMOVED***
				ls = append(ls, v)
			***REMOVED***
		***REMOVED***
		s.globalLock.RUnlock()
	***REMOVED***
	return ls, warnings, nil
***REMOVED***

// CreateWithRef creates a volume with the given name and driver and stores the ref
// This ensures there's no race between creating a volume and then storing a reference.
func (s *VolumeStore) CreateWithRef(name, driverName, ref string, opts, labels map[string]string) (volume.Volume, error) ***REMOVED***
	name = normalizeVolumeName(name)
	s.locks.Lock(name)
	defer s.locks.Unlock(name)

	v, err := s.create(name, driverName, opts, labels)
	if err != nil ***REMOVED***
		if _, ok := err.(*OpErr); ok ***REMOVED***
			return nil, err
		***REMOVED***
		return nil, &OpErr***REMOVED***Err: err, Name: name, Op: "create"***REMOVED***
	***REMOVED***

	s.setNamed(v, ref)
	return v, nil
***REMOVED***

// Create creates a volume with the given name and driver.
// This is just like CreateWithRef() except we don't store a reference while holding the lock.
func (s *VolumeStore) Create(name, driverName string, opts, labels map[string]string) (volume.Volume, error) ***REMOVED***
	return s.CreateWithRef(name, driverName, "", opts, labels)
***REMOVED***

// checkConflict checks the local cache for name collisions with the passed in name,
// for existing volumes with the same name but in a different driver.
// This is used by `Create` as a best effort to prevent name collisions for volumes.
// If a matching volume is found that is not a conflict that is returned so the caller
// does not need to perform an additional lookup.
// When no matching volume is found, both returns will be nil
//
// Note: This does not probe all the drivers for name collisions because v1 plugins
// are very slow, particularly if the plugin is down, and cause other issues,
// particularly around locking the store.
// TODO(cpuguy83): With v2 plugins this shouldn't be a problem. Could also potentially
// use a connect timeout for this kind of check to ensure we aren't blocking for a
// long time.
func (s *VolumeStore) checkConflict(name, driverName string) (volume.Volume, error) ***REMOVED***
	// check the local cache
	v, _ := s.getNamed(name)
	if v == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	vDriverName := v.DriverName()
	var conflict bool
	if driverName != "" ***REMOVED***
		// Retrieve canonical driver name to avoid inconsistencies (for example
		// "plugin" vs. "plugin:latest")
		vd, err := volumedrivers.GetDriver(driverName)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if vDriverName != vd.Name() ***REMOVED***
			conflict = true
		***REMOVED***
	***REMOVED***

	// let's check if the found volume ref
	// is stale by checking with the driver if it still exists
	exists, err := volumeExists(v)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(errNameConflict, "found reference to volume '%s' in driver '%s', but got an error while checking the driver: %v", name, vDriverName, err)
	***REMOVED***

	if exists ***REMOVED***
		if conflict ***REMOVED***
			return nil, errors.Wrapf(errNameConflict, "driver '%s' already has volume '%s'", vDriverName, name)
		***REMOVED***
		return v, nil
	***REMOVED***

	if s.hasRef(v.Name()) ***REMOVED***
		// Containers are referencing this volume but it doesn't seem to exist anywhere.
		// Return a conflict error here, the user can fix this with `docker volume rm -f`
		return nil, errors.Wrapf(errNameConflict, "found references to volume '%s' in driver '%s' but the volume was not found in the driver -- you may need to remove containers referencing this volume or force remove the volume to re-create it", name, vDriverName)
	***REMOVED***

	// doesn't exist, so purge it from the cache
	s.Purge(name)
	return nil, nil
***REMOVED***

// volumeExists returns if the volume is still present in the driver.
// An error is returned if there was an issue communicating with the driver.
func volumeExists(v volume.Volume) (bool, error) ***REMOVED***
	exists, err := lookupVolume(v.DriverName(), v.Name())
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return exists != nil, nil
***REMOVED***

// create asks the given driver to create a volume with the name/opts.
// If a volume with the name is already known, it will ask the stored driver for the volume.
// If the passed in driver name does not match the driver name which is stored
//  for the given volume name, an error is returned after checking if the reference is stale.
// If the reference is stale, it will be purged and this create can continue.
// It is expected that callers of this function hold any necessary locks.
func (s *VolumeStore) create(name, driverName string, opts, labels map[string]string) (volume.Volume, error) ***REMOVED***
	// Validate the name in a platform-specific manner

	// volume name validation is specific to the host os and not on container image
	// windows/lcow should have an equivalent volumename validation logic so we create a parser for current host OS
	parser := volume.NewParser(runtime.GOOS)
	err := parser.ValidateVolumeName(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	v, err := s.checkConflict(name, driverName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if v != nil ***REMOVED***
		// there is an existing volume, if we already have this stored locally, return it.
		// TODO: there could be some inconsistent details such as labels here
		if vv, _ := s.getNamed(v.Name()); vv != nil ***REMOVED***
			return vv, nil
		***REMOVED***
	***REMOVED***

	// Since there isn't a specified driver name, let's see if any of the existing drivers have this volume name
	if driverName == "" ***REMOVED***
		v, _ = s.getVolume(name)
		if v != nil ***REMOVED***
			return v, nil
		***REMOVED***
	***REMOVED***

	vd, err := volumedrivers.CreateDriver(driverName)
	if err != nil ***REMOVED***
		return nil, &OpErr***REMOVED***Op: "create", Name: name, Err: err***REMOVED***
	***REMOVED***

	logrus.Debugf("Registering new volume reference: driver %q, name %q", vd.Name(), name)
	if v, _ = vd.Get(name); v == nil ***REMOVED***
		v, err = vd.Create(name, opts)
		if err != nil ***REMOVED***
			if _, err := volumedrivers.ReleaseDriver(driverName); err != nil ***REMOVED***
				logrus.WithError(err).WithField("driver", driverName).Error("Error releasing reference to volume driver")
			***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	s.globalLock.Lock()
	s.labels[name] = labels
	s.options[name] = opts
	s.refs[name] = make(map[string]struct***REMOVED******REMOVED***)
	s.globalLock.Unlock()

	metadata := volumeMetadata***REMOVED***
		Name:    name,
		Driver:  vd.Name(),
		Labels:  labels,
		Options: opts,
	***REMOVED***

	if err := s.setMeta(name, metadata); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return volumeWrapper***REMOVED***v, labels, vd.Scope(), opts***REMOVED***, nil
***REMOVED***

// GetWithRef gets a volume with the given name from the passed in driver and stores the ref
// This is just like Get(), but we store the reference while holding the lock.
// This makes sure there are no races between checking for the existence of a volume and adding a reference for it
func (s *VolumeStore) GetWithRef(name, driverName, ref string) (volume.Volume, error) ***REMOVED***
	name = normalizeVolumeName(name)
	s.locks.Lock(name)
	defer s.locks.Unlock(name)

	vd, err := volumedrivers.GetDriver(driverName)
	if err != nil ***REMOVED***
		return nil, &OpErr***REMOVED***Err: err, Name: name, Op: "get"***REMOVED***
	***REMOVED***

	v, err := vd.Get(name)
	if err != nil ***REMOVED***
		return nil, &OpErr***REMOVED***Err: err, Name: name, Op: "get"***REMOVED***
	***REMOVED***

	s.setNamed(v, ref)

	s.globalLock.RLock()
	defer s.globalLock.RUnlock()
	return volumeWrapper***REMOVED***v, s.labels[name], vd.Scope(), s.options[name]***REMOVED***, nil
***REMOVED***

// Get looks if a volume with the given name exists and returns it if so
func (s *VolumeStore) Get(name string) (volume.Volume, error) ***REMOVED***
	name = normalizeVolumeName(name)
	s.locks.Lock(name)
	defer s.locks.Unlock(name)

	v, err := s.getVolume(name)
	if err != nil ***REMOVED***
		return nil, &OpErr***REMOVED***Err: err, Name: name, Op: "get"***REMOVED***
	***REMOVED***
	s.setNamed(v, "")
	return v, nil
***REMOVED***

// getVolume requests the volume, if the driver info is stored it just accesses that driver,
// if the driver is unknown it probes all drivers until it finds the first volume with that name.
// it is expected that callers of this function hold any necessary locks
func (s *VolumeStore) getVolume(name string) (volume.Volume, error) ***REMOVED***
	var meta volumeMetadata
	meta, err := s.getMeta(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	driverName := meta.Driver
	if driverName == "" ***REMOVED***
		s.globalLock.RLock()
		v, exists := s.names[name]
		s.globalLock.RUnlock()
		if exists ***REMOVED***
			meta.Driver = v.DriverName()
			if err := s.setMeta(name, meta); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if meta.Driver != "" ***REMOVED***
		vol, err := lookupVolume(meta.Driver, name)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if vol == nil ***REMOVED***
			s.Purge(name)
			return nil, errNoSuchVolume
		***REMOVED***

		var scope string
		vd, err := volumedrivers.GetDriver(meta.Driver)
		if err == nil ***REMOVED***
			scope = vd.Scope()
		***REMOVED***
		return volumeWrapper***REMOVED***vol, meta.Labels, scope, meta.Options***REMOVED***, nil
	***REMOVED***

	logrus.Debugf("Probing all drivers for volume with name: %s", name)
	drivers, err := volumedrivers.GetAllDrivers()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, d := range drivers ***REMOVED***
		v, err := d.Get(name)
		if err != nil || v == nil ***REMOVED***
			continue
		***REMOVED***
		meta.Driver = v.DriverName()
		if err := s.setMeta(name, meta); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return volumeWrapper***REMOVED***v, meta.Labels, d.Scope(), meta.Options***REMOVED***, nil
	***REMOVED***
	return nil, errNoSuchVolume
***REMOVED***

// lookupVolume gets the specified volume from the specified driver.
// This will only return errors related to communications with the driver.
// If the driver returns an error that is not communication related the
//   error is logged but not returned.
// If the volume is not found it will return `nil, nil``
func lookupVolume(driverName, volumeName string) (volume.Volume, error) ***REMOVED***
	vd, err := volumedrivers.GetDriver(driverName)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "error while checking if volume %q exists in driver %q", volumeName, driverName)
	***REMOVED***
	v, err := vd.Get(volumeName)
	if err != nil ***REMOVED***
		err = errors.Cause(err)
		if _, ok := err.(net.Error); ok ***REMOVED***
			if v != nil ***REMOVED***
				volumeName = v.Name()
				driverName = v.DriverName()
			***REMOVED***
			return nil, errors.Wrapf(err, "error while checking if volume %q exists in driver %q", volumeName, driverName)
		***REMOVED***

		// At this point, the error could be anything from the driver, such as "no such volume"
		// Let's not check an error here, and instead check if the driver returned a volume
		logrus.WithError(err).WithField("driver", driverName).WithField("volume", volumeName).Warnf("Error while looking up volume")
	***REMOVED***
	return v, nil
***REMOVED***

// Remove removes the requested volume. A volume is not removed if it has any refs
func (s *VolumeStore) Remove(v volume.Volume) error ***REMOVED***
	name := normalizeVolumeName(v.Name())
	s.locks.Lock(name)
	defer s.locks.Unlock(name)

	if s.hasRef(name) ***REMOVED***
		return &OpErr***REMOVED***Err: errVolumeInUse, Name: v.Name(), Op: "remove", Refs: s.getRefs(name)***REMOVED***
	***REMOVED***

	vd, err := volumedrivers.GetDriver(v.DriverName())
	if err != nil ***REMOVED***
		return &OpErr***REMOVED***Err: err, Name: v.DriverName(), Op: "remove"***REMOVED***
	***REMOVED***

	logrus.Debugf("Removing volume reference: driver %s, name %s", v.DriverName(), name)
	vol := unwrapVolume(v)
	if err := vd.Remove(vol); err != nil ***REMOVED***
		return &OpErr***REMOVED***Err: err, Name: name, Op: "remove"***REMOVED***
	***REMOVED***

	s.Purge(name)
	return nil
***REMOVED***

// Dereference removes the specified reference to the volume
func (s *VolumeStore) Dereference(v volume.Volume, ref string) ***REMOVED***
	name := v.Name()

	s.locks.Lock(name)
	defer s.locks.Unlock(name)

	s.globalLock.Lock()
	defer s.globalLock.Unlock()

	if s.refs[name] != nil ***REMOVED***
		delete(s.refs[name], ref)
	***REMOVED***
***REMOVED***

// Refs gets the current list of refs for the given volume
func (s *VolumeStore) Refs(v volume.Volume) []string ***REMOVED***
	name := v.Name()

	s.locks.Lock(name)
	defer s.locks.Unlock(name)

	return s.getRefs(name)
***REMOVED***

// FilterByDriver returns the available volumes filtered by driver name
func (s *VolumeStore) FilterByDriver(name string) ([]volume.Volume, error) ***REMOVED***
	vd, err := volumedrivers.GetDriver(name)
	if err != nil ***REMOVED***
		return nil, &OpErr***REMOVED***Err: err, Name: name, Op: "list"***REMOVED***
	***REMOVED***
	ls, err := vd.List()
	if err != nil ***REMOVED***
		return nil, &OpErr***REMOVED***Err: err, Name: name, Op: "list"***REMOVED***
	***REMOVED***
	for i, v := range ls ***REMOVED***
		options := map[string]string***REMOVED******REMOVED***
		s.globalLock.RLock()
		for key, value := range s.options[v.Name()] ***REMOVED***
			options[key] = value
		***REMOVED***
		ls[i] = volumeWrapper***REMOVED***v, s.labels[v.Name()], vd.Scope(), options***REMOVED***
		s.globalLock.RUnlock()
	***REMOVED***
	return ls, nil
***REMOVED***

// FilterByUsed returns the available volumes filtered by if they are in use or not.
// `used=true` returns only volumes that are being used, while `used=false` returns
// only volumes that are not being used.
func (s *VolumeStore) FilterByUsed(vols []volume.Volume, used bool) []volume.Volume ***REMOVED***
	return s.filter(vols, func(v volume.Volume) bool ***REMOVED***
		s.locks.Lock(v.Name())
		hasRef := s.hasRef(v.Name())
		s.locks.Unlock(v.Name())
		return used == hasRef
	***REMOVED***)
***REMOVED***

// filterFunc defines a function to allow filter volumes in the store
type filterFunc func(vol volume.Volume) bool

// filter returns the available volumes filtered by a filterFunc function
func (s *VolumeStore) filter(vols []volume.Volume, f filterFunc) []volume.Volume ***REMOVED***
	var ls []volume.Volume
	for _, v := range vols ***REMOVED***
		if f(v) ***REMOVED***
			ls = append(ls, v)
		***REMOVED***
	***REMOVED***
	return ls
***REMOVED***

func unwrapVolume(v volume.Volume) volume.Volume ***REMOVED***
	if vol, ok := v.(volumeWrapper); ok ***REMOVED***
		return vol.Volume
	***REMOVED***

	return v
***REMOVED***

// Shutdown releases all resources used by the volume store
// It does not make any changes to volumes, drivers, etc.
func (s *VolumeStore) Shutdown() error ***REMOVED***
	return s.db.Close()
***REMOVED***
