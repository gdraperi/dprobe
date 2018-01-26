package vfs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/daemon/graphdriver/quota"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/system"
	units "github.com/docker/go-units"
	"github.com/opencontainers/selinux/go-selinux/label"
)

var (
	// CopyDir defines the copy method to use.
	CopyDir = dirCopy
)

func init() ***REMOVED***
	graphdriver.Register("vfs", Init)
***REMOVED***

// Init returns a new VFS driver.
// This sets the home directory for the driver and returns NaiveDiffDriver.
func Init(home string, options []string, uidMaps, gidMaps []idtools.IDMap) (graphdriver.Driver, error) ***REMOVED***
	d := &Driver***REMOVED***
		home:       home,
		idMappings: idtools.NewIDMappingsFromMaps(uidMaps, gidMaps),
	***REMOVED***
	rootIDs := d.idMappings.RootPair()
	if err := idtools.MkdirAllAndChown(home, 0700, rootIDs); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	setupDriverQuota(d)

	return graphdriver.NewNaiveDiffDriver(d, uidMaps, gidMaps), nil
***REMOVED***

// Driver holds information about the driver, home directory of the driver.
// Driver implements graphdriver.ProtoDriver. It uses only basic vfs operations.
// In order to support layering, files are copied from the parent layer into the new layer. There is no copy-on-write support.
// Driver must be wrapped in NaiveDiffDriver to be used as a graphdriver.Driver
type Driver struct ***REMOVED***
	driverQuota
	home       string
	idMappings *idtools.IDMappings
***REMOVED***

func (d *Driver) String() string ***REMOVED***
	return "vfs"
***REMOVED***

// Status is used for implementing the graphdriver.ProtoDriver interface. VFS does not currently have any status information.
func (d *Driver) Status() [][2]string ***REMOVED***
	return nil
***REMOVED***

// GetMetadata is used for implementing the graphdriver.ProtoDriver interface. VFS does not currently have any meta data.
func (d *Driver) GetMetadata(id string) (map[string]string, error) ***REMOVED***
	return nil, nil
***REMOVED***

// Cleanup is used to implement graphdriver.ProtoDriver. There is no cleanup required for this driver.
func (d *Driver) Cleanup() error ***REMOVED***
	return nil
***REMOVED***

// CreateReadWrite creates a layer that is writable for use as a container
// file system.
func (d *Driver) CreateReadWrite(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	var err error
	var size int64

	if opts != nil ***REMOVED***
		for key, val := range opts.StorageOpt ***REMOVED***
			switch key ***REMOVED***
			case "size":
				if !d.quotaSupported() ***REMOVED***
					return quota.ErrQuotaNotSupported
				***REMOVED***
				if size, err = units.RAMInBytes(val); err != nil ***REMOVED***
					return err
				***REMOVED***
			default:
				return fmt.Errorf("Storage opt %s not supported", key)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return d.create(id, parent, uint64(size))
***REMOVED***

// Create prepares the filesystem for the VFS driver and copies the directory for the given id under the parent.
func (d *Driver) Create(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	if opts != nil && len(opts.StorageOpt) != 0 ***REMOVED***
		return fmt.Errorf("--storage-opt is not supported for vfs on read-only layers")
	***REMOVED***

	return d.create(id, parent, 0)
***REMOVED***

func (d *Driver) create(id, parent string, size uint64) error ***REMOVED***
	dir := d.dir(id)
	rootIDs := d.idMappings.RootPair()
	if err := idtools.MkdirAllAndChown(filepath.Dir(dir), 0700, rootIDs); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := idtools.MkdirAndChown(dir, 0755, rootIDs); err != nil ***REMOVED***
		return err
	***REMOVED***

	if size != 0 ***REMOVED***
		if err := d.setupQuota(dir, size); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	labelOpts := []string***REMOVED***"level:s0"***REMOVED***
	if _, mountLabel, err := label.InitLabels(labelOpts); err == nil ***REMOVED***
		label.SetFileLabel(dir, mountLabel)
	***REMOVED***
	if parent == "" ***REMOVED***
		return nil
	***REMOVED***
	parentDir, err := d.Get(parent, "")
	if err != nil ***REMOVED***
		return fmt.Errorf("%s: %s", parent, err)
	***REMOVED***
	return CopyDir(parentDir.Path(), dir)
***REMOVED***

func (d *Driver) dir(id string) string ***REMOVED***
	return filepath.Join(d.home, "dir", filepath.Base(id))
***REMOVED***

// Remove deletes the content from the directory for a given id.
func (d *Driver) Remove(id string) error ***REMOVED***
	return system.EnsureRemoveAll(d.dir(id))
***REMOVED***

// Get returns the directory for the given id.
func (d *Driver) Get(id, mountLabel string) (containerfs.ContainerFS, error) ***REMOVED***
	dir := d.dir(id)
	if st, err := os.Stat(dir); err != nil ***REMOVED***
		return nil, err
	***REMOVED*** else if !st.IsDir() ***REMOVED***
		return nil, fmt.Errorf("%s: not a directory", dir)
	***REMOVED***
	return containerfs.NewLocalContainerFS(dir), nil
***REMOVED***

// Put is a noop for vfs that return nil for the error, since this driver has no runtime resources to clean up.
func (d *Driver) Put(id string) error ***REMOVED***
	// The vfs driver has no runtime resources (e.g. mounts)
	// to clean up, so we don't need anything here
	return nil
***REMOVED***

// Exists checks to see if the directory exists for the given id.
func (d *Driver) Exists(id string) bool ***REMOVED***
	_, err := os.Stat(d.dir(id))
	return err == nil
***REMOVED***
