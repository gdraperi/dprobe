//+build windows

package windows

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/Microsoft/go-winio"
	"github.com/Microsoft/go-winio/archive/tar"
	"github.com/Microsoft/go-winio/backuptar"
	"github.com/Microsoft/hcsshim"
	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/longpath"
	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/docker/pkg/system"
	units "github.com/docker/go-units"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

// filterDriver is an HCSShim driver type for the Windows Filter driver.
const filterDriver = 1

var (
	// mutatedFiles is a list of files that are mutated by the import process
	// and must be backed up and restored.
	mutatedFiles = map[string]string***REMOVED***
		"UtilityVM/Files/EFI/Microsoft/Boot/BCD":      "bcd.bak",
		"UtilityVM/Files/EFI/Microsoft/Boot/BCD.LOG":  "bcd.log.bak",
		"UtilityVM/Files/EFI/Microsoft/Boot/BCD.LOG1": "bcd.log1.bak",
		"UtilityVM/Files/EFI/Microsoft/Boot/BCD.LOG2": "bcd.log2.bak",
	***REMOVED***
	noreexec = false
)

// init registers the windows graph drivers to the register.
func init() ***REMOVED***
	graphdriver.Register("windowsfilter", InitFilter)
	// DOCKER_WINDOWSFILTER_NOREEXEC allows for inline processing which makes
	// debugging issues in the re-exec codepath significantly easier.
	if os.Getenv("DOCKER_WINDOWSFILTER_NOREEXEC") != "" ***REMOVED***
		logrus.Warnf("WindowsGraphDriver is set to not re-exec. This is intended for debugging purposes only.")
		noreexec = true
	***REMOVED*** else ***REMOVED***
		reexec.Register("docker-windows-write-layer", writeLayerReexec)
	***REMOVED***
***REMOVED***

type checker struct ***REMOVED***
***REMOVED***

func (c *checker) IsMounted(path string) bool ***REMOVED***
	return false
***REMOVED***

// Driver represents a windows graph driver.
type Driver struct ***REMOVED***
	// info stores the shim driver information
	info hcsshim.DriverInfo
	ctr  *graphdriver.RefCounter
	// it is safe for windows to use a cache here because it does not support
	// restoring containers when the daemon dies.
	cacheMu sync.Mutex
	cache   map[string]string
***REMOVED***

// InitFilter returns a new Windows storage filter driver.
func InitFilter(home string, options []string, uidMaps, gidMaps []idtools.IDMap) (graphdriver.Driver, error) ***REMOVED***
	logrus.Debugf("WindowsGraphDriver InitFilter at %s", home)

	fsType, err := getFileSystemType(string(home[0]))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if strings.ToLower(fsType) == "refs" ***REMOVED***
		return nil, fmt.Errorf("%s is on an ReFS volume - ReFS volumes are not supported", home)
	***REMOVED***

	if err := idtools.MkdirAllAndChown(home, 0700, idtools.IDPair***REMOVED***UID: 0, GID: 0***REMOVED***); err != nil ***REMOVED***
		return nil, fmt.Errorf("windowsfilter failed to create '%s': %v", home, err)
	***REMOVED***

	d := &Driver***REMOVED***
		info: hcsshim.DriverInfo***REMOVED***
			HomeDir: home,
			Flavour: filterDriver,
		***REMOVED***,
		cache: make(map[string]string),
		ctr:   graphdriver.NewRefCounter(&checker***REMOVED******REMOVED***),
	***REMOVED***
	return d, nil
***REMOVED***

// win32FromHresult is a helper function to get the win32 error code from an HRESULT
func win32FromHresult(hr uintptr) uintptr ***REMOVED***
	if hr&0x1fff0000 == 0x00070000 ***REMOVED***
		return hr & 0xffff
	***REMOVED***
	return hr
***REMOVED***

// getFileSystemType obtains the type of a file system through GetVolumeInformation
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa364993(v=vs.85).aspx
func getFileSystemType(drive string) (fsType string, hr error) ***REMOVED***
	var (
		modkernel32              = windows.NewLazySystemDLL("kernel32.dll")
		procGetVolumeInformation = modkernel32.NewProc("GetVolumeInformationW")
		buf                      = make([]uint16, 255)
		size                     = windows.MAX_PATH + 1
	)
	if len(drive) != 1 ***REMOVED***
		hr = errors.New("getFileSystemType must be called with a drive letter")
		return
	***REMOVED***
	drive += `:\`
	n := uintptr(unsafe.Pointer(nil))
	r0, _, _ := syscall.Syscall9(procGetVolumeInformation.Addr(), 8, uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(drive))), n, n, n, n, n, uintptr(unsafe.Pointer(&buf[0])), uintptr(size), 0)
	if int32(r0) < 0 ***REMOVED***
		hr = syscall.Errno(win32FromHresult(r0))
	***REMOVED***
	fsType = windows.UTF16ToString(buf)
	return
***REMOVED***

// String returns the string representation of a driver. This should match
// the name the graph driver has been registered with.
func (d *Driver) String() string ***REMOVED***
	return "windowsfilter"
***REMOVED***

// Status returns the status of the driver.
func (d *Driver) Status() [][2]string ***REMOVED***
	return [][2]string***REMOVED***
		***REMOVED***"Windows", ""***REMOVED***,
	***REMOVED***
***REMOVED***

// Exists returns true if the given id is registered with this driver.
func (d *Driver) Exists(id string) bool ***REMOVED***
	rID, err := d.resolveID(id)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	result, err := hcsshim.LayerExists(d.info, rID)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return result
***REMOVED***

// CreateReadWrite creates a layer that is writable for use as a container
// file system.
func (d *Driver) CreateReadWrite(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	if opts != nil ***REMOVED***
		return d.create(id, parent, opts.MountLabel, false, opts.StorageOpt)
	***REMOVED***
	return d.create(id, parent, "", false, nil)
***REMOVED***

// Create creates a new read-only layer with the given id.
func (d *Driver) Create(id, parent string, opts *graphdriver.CreateOpts) error ***REMOVED***
	if opts != nil ***REMOVED***
		return d.create(id, parent, opts.MountLabel, true, opts.StorageOpt)
	***REMOVED***
	return d.create(id, parent, "", true, nil)
***REMOVED***

func (d *Driver) create(id, parent, mountLabel string, readOnly bool, storageOpt map[string]string) error ***REMOVED***
	rPId, err := d.resolveID(parent)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	parentChain, err := d.getLayerChain(rPId)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var layerChain []string

	if rPId != "" ***REMOVED***
		parentPath, err := hcsshim.GetLayerMountPath(d.info, rPId)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, err := os.Stat(filepath.Join(parentPath, "Files")); err == nil ***REMOVED***
			// This is a legitimate parent layer (not the empty "-init" layer),
			// so include it in the layer chain.
			layerChain = []string***REMOVED***parentPath***REMOVED***
		***REMOVED***
	***REMOVED***

	layerChain = append(layerChain, parentChain...)

	if readOnly ***REMOVED***
		if err := hcsshim.CreateLayer(d.info, id, rPId); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var parentPath string
		if len(layerChain) != 0 ***REMOVED***
			parentPath = layerChain[0]
		***REMOVED***

		if err := hcsshim.CreateSandboxLayer(d.info, id, parentPath, layerChain); err != nil ***REMOVED***
			return err
		***REMOVED***

		storageOptions, err := parseStorageOpt(storageOpt)
		if err != nil ***REMOVED***
			return fmt.Errorf("Failed to parse storage options - %s", err)
		***REMOVED***

		if storageOptions.size != 0 ***REMOVED***
			if err := hcsshim.ExpandSandboxSize(d.info, id, storageOptions.size); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if _, err := os.Lstat(d.dir(parent)); err != nil ***REMOVED***
		if err2 := hcsshim.DestroyLayer(d.info, id); err2 != nil ***REMOVED***
			logrus.Warnf("Failed to DestroyLayer %s: %s", id, err2)
		***REMOVED***
		return fmt.Errorf("Cannot create layer with missing parent %s: %s", parent, err)
	***REMOVED***

	if err := d.setLayerChain(id, layerChain); err != nil ***REMOVED***
		if err2 := hcsshim.DestroyLayer(d.info, id); err2 != nil ***REMOVED***
			logrus.Warnf("Failed to DestroyLayer %s: %s", id, err2)
		***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// dir returns the absolute path to the layer.
func (d *Driver) dir(id string) string ***REMOVED***
	return filepath.Join(d.info.HomeDir, filepath.Base(id))
***REMOVED***

// Remove unmounts and removes the dir information.
func (d *Driver) Remove(id string) error ***REMOVED***
	rID, err := d.resolveID(id)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// This retry loop is due to a bug in Windows (Internal bug #9432268)
	// if GetContainers fails with ErrVmcomputeOperationInvalidState
	// it is a transient error. Retry until it succeeds.
	var computeSystems []hcsshim.ContainerProperties
	retryCount := 0
	osv := system.GetOSVersion()
	for ***REMOVED***
		// Get and terminate any template VMs that are currently using the layer.
		// Note: It is unfortunate that we end up in the graphdrivers Remove() call
		// for both containers and images, but the logic for template VMs is only
		// needed for images - specifically we are looking to see if a base layer
		// is in use by a template VM as a result of having started a Hyper-V
		// container at some point.
		//
		// We have a retry loop for ErrVmcomputeOperationInvalidState and
		// ErrVmcomputeOperationAccessIsDenied as there is a race condition
		// in RS1 and RS2 building during enumeration when a silo is going away
		// for example under it, in HCS. AccessIsDenied added to fix 30278.
		//
		// TODO @jhowardmsft - For RS3, we can remove the retries. Also consider
		// using platform APIs (if available) to get this more succinctly. Also
		// consider enhancing the Remove() interface to have context of why
		// the remove is being called - that could improve efficiency by not
		// enumerating compute systems during a remove of a container as it's
		// not required.
		computeSystems, err = hcsshim.GetContainers(hcsshim.ComputeSystemQuery***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			if (osv.Build < 15139) &&
				((err == hcsshim.ErrVmcomputeOperationInvalidState) || (err == hcsshim.ErrVmcomputeOperationAccessIsDenied)) ***REMOVED***
				if retryCount >= 500 ***REMOVED***
					break
				***REMOVED***
				retryCount++
				time.Sleep(10 * time.Millisecond)
				continue
			***REMOVED***
			return err
		***REMOVED***
		break
	***REMOVED***

	for _, computeSystem := range computeSystems ***REMOVED***
		if strings.Contains(computeSystem.RuntimeImagePath, id) && computeSystem.IsRuntimeTemplate ***REMOVED***
			container, err := hcsshim.OpenContainer(computeSystem.ID)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			defer container.Close()
			err = container.Terminate()
			if hcsshim.IsPending(err) ***REMOVED***
				err = container.Wait()
			***REMOVED*** else if hcsshim.IsAlreadyStopped(err) ***REMOVED***
				err = nil
			***REMOVED***

			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	layerPath := filepath.Join(d.info.HomeDir, rID)
	tmpID := fmt.Sprintf("%s-removing", rID)
	tmpLayerPath := filepath.Join(d.info.HomeDir, tmpID)
	if err := os.Rename(layerPath, tmpLayerPath); err != nil && !os.IsNotExist(err) ***REMOVED***
		return err
	***REMOVED***
	if err := hcsshim.DestroyLayer(d.info, tmpID); err != nil ***REMOVED***
		logrus.Errorf("Failed to DestroyLayer %s: %s", id, err)
	***REMOVED***

	return nil
***REMOVED***

// Get returns the rootfs path for the id. This will mount the dir at its given path.
func (d *Driver) Get(id, mountLabel string) (containerfs.ContainerFS, error) ***REMOVED***
	logrus.Debugf("WindowsGraphDriver Get() id %s mountLabel %s", id, mountLabel)
	var dir string

	rID, err := d.resolveID(id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if count := d.ctr.Increment(rID); count > 1 ***REMOVED***
		return containerfs.NewLocalContainerFS(d.cache[rID]), nil
	***REMOVED***

	// Getting the layer paths must be done outside of the lock.
	layerChain, err := d.getLayerChain(rID)
	if err != nil ***REMOVED***
		d.ctr.Decrement(rID)
		return nil, err
	***REMOVED***

	if err := hcsshim.ActivateLayer(d.info, rID); err != nil ***REMOVED***
		d.ctr.Decrement(rID)
		return nil, err
	***REMOVED***
	if err := hcsshim.PrepareLayer(d.info, rID, layerChain); err != nil ***REMOVED***
		d.ctr.Decrement(rID)
		if err2 := hcsshim.DeactivateLayer(d.info, rID); err2 != nil ***REMOVED***
			logrus.Warnf("Failed to Deactivate %s: %s", id, err)
		***REMOVED***
		return nil, err
	***REMOVED***

	mountPath, err := hcsshim.GetLayerMountPath(d.info, rID)
	if err != nil ***REMOVED***
		d.ctr.Decrement(rID)
		if err := hcsshim.UnprepareLayer(d.info, rID); err != nil ***REMOVED***
			logrus.Warnf("Failed to Unprepare %s: %s", id, err)
		***REMOVED***
		if err2 := hcsshim.DeactivateLayer(d.info, rID); err2 != nil ***REMOVED***
			logrus.Warnf("Failed to Deactivate %s: %s", id, err)
		***REMOVED***
		return nil, err
	***REMOVED***
	d.cacheMu.Lock()
	d.cache[rID] = mountPath
	d.cacheMu.Unlock()

	// If the layer has a mount path, use that. Otherwise, use the
	// folder path.
	if mountPath != "" ***REMOVED***
		dir = mountPath
	***REMOVED*** else ***REMOVED***
		dir = d.dir(id)
	***REMOVED***

	return containerfs.NewLocalContainerFS(dir), nil
***REMOVED***

// Put adds a new layer to the driver.
func (d *Driver) Put(id string) error ***REMOVED***
	logrus.Debugf("WindowsGraphDriver Put() id %s", id)

	rID, err := d.resolveID(id)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if count := d.ctr.Decrement(rID); count > 0 ***REMOVED***
		return nil
	***REMOVED***
	d.cacheMu.Lock()
	_, exists := d.cache[rID]
	delete(d.cache, rID)
	d.cacheMu.Unlock()

	// If the cache was not populated, then the layer was left unprepared and deactivated
	if !exists ***REMOVED***
		return nil
	***REMOVED***

	if err := hcsshim.UnprepareLayer(d.info, rID); err != nil ***REMOVED***
		return err
	***REMOVED***
	return hcsshim.DeactivateLayer(d.info, rID)
***REMOVED***

// Cleanup ensures the information the driver stores is properly removed.
// We use this opportunity to cleanup any -removing folders which may be
// still left if the daemon was killed while it was removing a layer.
func (d *Driver) Cleanup() error ***REMOVED***
	items, err := ioutil.ReadDir(d.info.HomeDir)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	// Note we don't return an error below - it's possible the files
	// are locked. However, next time around after the daemon exits,
	// we likely will be able to to cleanup successfully. Instead we log
	// warnings if there are errors.
	for _, item := range items ***REMOVED***
		if item.IsDir() && strings.HasSuffix(item.Name(), "-removing") ***REMOVED***
			if err := hcsshim.DestroyLayer(d.info, item.Name()); err != nil ***REMOVED***
				logrus.Warnf("Failed to cleanup %s: %s", item.Name(), err)
			***REMOVED*** else ***REMOVED***
				logrus.Infof("Cleaned up %s", item.Name())
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Diff produces an archive of the changes between the specified
// layer and its parent layer which may be "".
// The layer should be mounted when calling this function
func (d *Driver) Diff(id, parent string) (_ io.ReadCloser, err error) ***REMOVED***
	rID, err := d.resolveID(id)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	layerChain, err := d.getLayerChain(rID)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	// this is assuming that the layer is unmounted
	if err := hcsshim.UnprepareLayer(d.info, rID); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	prepare := func() ***REMOVED***
		if err := hcsshim.PrepareLayer(d.info, rID, layerChain); err != nil ***REMOVED***
			logrus.Warnf("Failed to Deactivate %s: %s", rID, err)
		***REMOVED***
	***REMOVED***

	arch, err := d.exportLayer(rID, layerChain)
	if err != nil ***REMOVED***
		prepare()
		return
	***REMOVED***
	return ioutils.NewReadCloserWrapper(arch, func() error ***REMOVED***
		err := arch.Close()
		prepare()
		return err
	***REMOVED***), nil
***REMOVED***

// Changes produces a list of changes between the specified layer
// and its parent layer. If parent is "", then all changes will be ADD changes.
// The layer should not be mounted when calling this function.
func (d *Driver) Changes(id, parent string) ([]archive.Change, error) ***REMOVED***
	rID, err := d.resolveID(id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	parentChain, err := d.getLayerChain(rID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := hcsshim.ActivateLayer(d.info, rID); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err2 := hcsshim.DeactivateLayer(d.info, rID); err2 != nil ***REMOVED***
			logrus.Errorf("changes() failed to DeactivateLayer %s %s: %s", id, rID, err2)
		***REMOVED***
	***REMOVED***()

	var changes []archive.Change
	err = winio.RunWithPrivilege(winio.SeBackupPrivilege, func() error ***REMOVED***
		r, err := hcsshim.NewLayerReader(d.info, id, parentChain)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer r.Close()

		for ***REMOVED***
			name, _, fileInfo, err := r.Next()
			if err == io.EOF ***REMOVED***
				return nil
			***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			name = filepath.ToSlash(name)
			if fileInfo == nil ***REMOVED***
				changes = append(changes, archive.Change***REMOVED***Path: name, Kind: archive.ChangeDelete***REMOVED***)
			***REMOVED*** else ***REMOVED***
				// Currently there is no way to tell between an add and a modify.
				changes = append(changes, archive.Change***REMOVED***Path: name, Kind: archive.ChangeModify***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return changes, nil
***REMOVED***

// ApplyDiff extracts the changeset from the given diff into the
// layer with the specified id and parent, returning the size of the
// new layer in bytes.
// The layer should not be mounted when calling this function
func (d *Driver) ApplyDiff(id, parent string, diff io.Reader) (int64, error) ***REMOVED***
	var layerChain []string
	if parent != "" ***REMOVED***
		rPId, err := d.resolveID(parent)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		parentChain, err := d.getLayerChain(rPId)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		parentPath, err := hcsshim.GetLayerMountPath(d.info, rPId)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		layerChain = append(layerChain, parentPath)
		layerChain = append(layerChain, parentChain...)
	***REMOVED***

	size, err := d.importLayer(id, diff, layerChain)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	if err = d.setLayerChain(id, layerChain); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return size, nil
***REMOVED***

// DiffSize calculates the changes between the specified layer
// and its parent and returns the size in bytes of the changes
// relative to its base filesystem directory.
func (d *Driver) DiffSize(id, parent string) (size int64, err error) ***REMOVED***
	rPId, err := d.resolveID(parent)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	changes, err := d.Changes(id, rPId)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	layerFs, err := d.Get(id, "")
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer d.Put(id)

	return archive.ChangesSize(layerFs.Path(), changes), nil
***REMOVED***

// GetMetadata returns custom driver information.
func (d *Driver) GetMetadata(id string) (map[string]string, error) ***REMOVED***
	m := make(map[string]string)
	m["dir"] = d.dir(id)
	return m, nil
***REMOVED***

func writeTarFromLayer(r hcsshim.LayerReader, w io.Writer) error ***REMOVED***
	t := tar.NewWriter(w)
	for ***REMOVED***
		name, size, fileInfo, err := r.Next()
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if fileInfo == nil ***REMOVED***
			// Write a whiteout file.
			hdr := &tar.Header***REMOVED***
				Name: filepath.ToSlash(filepath.Join(filepath.Dir(name), archive.WhiteoutPrefix+filepath.Base(name))),
			***REMOVED***
			err := t.WriteHeader(hdr)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = backuptar.WriteTarFileFromBackupStream(t, r, name, size, fileInfo)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return t.Close()
***REMOVED***

// exportLayer generates an archive from a layer based on the given ID.
func (d *Driver) exportLayer(id string, parentLayerPaths []string) (io.ReadCloser, error) ***REMOVED***
	archive, w := io.Pipe()
	go func() ***REMOVED***
		err := winio.RunWithPrivilege(winio.SeBackupPrivilege, func() error ***REMOVED***
			r, err := hcsshim.NewLayerReader(d.info, id, parentLayerPaths)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			err = writeTarFromLayer(r, w)
			cerr := r.Close()
			if err == nil ***REMOVED***
				err = cerr
			***REMOVED***
			return err
		***REMOVED***)
		w.CloseWithError(err)
	***REMOVED***()

	return archive, nil
***REMOVED***

// writeBackupStreamFromTarAndSaveMutatedFiles reads data from a tar stream and
// writes it to a backup stream, and also saves any files that will be mutated
// by the import layer process to a backup location.
func writeBackupStreamFromTarAndSaveMutatedFiles(buf *bufio.Writer, w io.Writer, t *tar.Reader, hdr *tar.Header, root string) (nextHdr *tar.Header, err error) ***REMOVED***
	var bcdBackup *os.File
	var bcdBackupWriter *winio.BackupFileWriter
	if backupPath, ok := mutatedFiles[hdr.Name]; ok ***REMOVED***
		bcdBackup, err = os.Create(filepath.Join(root, backupPath))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer func() ***REMOVED***
			cerr := bcdBackup.Close()
			if err == nil ***REMOVED***
				err = cerr
			***REMOVED***
		***REMOVED***()

		bcdBackupWriter = winio.NewBackupFileWriter(bcdBackup, false)
		defer func() ***REMOVED***
			cerr := bcdBackupWriter.Close()
			if err == nil ***REMOVED***
				err = cerr
			***REMOVED***
		***REMOVED***()

		buf.Reset(io.MultiWriter(w, bcdBackupWriter))
	***REMOVED*** else ***REMOVED***
		buf.Reset(w)
	***REMOVED***

	defer func() ***REMOVED***
		ferr := buf.Flush()
		if err == nil ***REMOVED***
			err = ferr
		***REMOVED***
	***REMOVED***()

	return backuptar.WriteBackupStreamFromTarFile(buf, t, hdr)
***REMOVED***

func writeLayerFromTar(r io.Reader, w hcsshim.LayerWriter, root string) (int64, error) ***REMOVED***
	t := tar.NewReader(r)
	hdr, err := t.Next()
	totalSize := int64(0)
	buf := bufio.NewWriter(nil)
	for err == nil ***REMOVED***
		base := path.Base(hdr.Name)
		if strings.HasPrefix(base, archive.WhiteoutPrefix) ***REMOVED***
			name := path.Join(path.Dir(hdr.Name), base[len(archive.WhiteoutPrefix):])
			err = w.Remove(filepath.FromSlash(name))
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			hdr, err = t.Next()
		***REMOVED*** else if hdr.Typeflag == tar.TypeLink ***REMOVED***
			err = w.AddLink(filepath.FromSlash(hdr.Name), filepath.FromSlash(hdr.Linkname))
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			hdr, err = t.Next()
		***REMOVED*** else ***REMOVED***
			var (
				name     string
				size     int64
				fileInfo *winio.FileBasicInfo
			)
			name, size, fileInfo, err = backuptar.FileInfoFromHeader(hdr)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			err = w.Add(filepath.FromSlash(name), fileInfo)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			hdr, err = writeBackupStreamFromTarAndSaveMutatedFiles(buf, w, t, hdr, root)
			totalSize += size
		***REMOVED***
	***REMOVED***
	if err != io.EOF ***REMOVED***
		return 0, err
	***REMOVED***
	return totalSize, nil
***REMOVED***

// importLayer adds a new layer to the tag and graph store based on the given data.
func (d *Driver) importLayer(id string, layerData io.Reader, parentLayerPaths []string) (size int64, err error) ***REMOVED***
	if !noreexec ***REMOVED***
		cmd := reexec.Command(append([]string***REMOVED***"docker-windows-write-layer", d.info.HomeDir, id***REMOVED***, parentLayerPaths...)...)
		output := bytes.NewBuffer(nil)
		cmd.Stdin = layerData
		cmd.Stdout = output
		cmd.Stderr = output

		if err = cmd.Start(); err != nil ***REMOVED***
			return
		***REMOVED***

		if err = cmd.Wait(); err != nil ***REMOVED***
			return 0, fmt.Errorf("re-exec error: %v: output: %s", err, output)
		***REMOVED***

		return strconv.ParseInt(output.String(), 10, 64)
	***REMOVED***
	return writeLayer(layerData, d.info.HomeDir, id, parentLayerPaths...)
***REMOVED***

// writeLayerReexec is the re-exec entry point for writing a layer from a tar file
func writeLayerReexec() ***REMOVED***
	size, err := writeLayer(os.Stdin, os.Args[1], os.Args[2], os.Args[3:]...)
	if err != nil ***REMOVED***
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	***REMOVED***
	fmt.Fprint(os.Stdout, size)
***REMOVED***

// writeLayer writes a layer from a tar file.
func writeLayer(layerData io.Reader, home string, id string, parentLayerPaths ...string) (int64, error) ***REMOVED***
	err := winio.EnableProcessPrivileges([]string***REMOVED***winio.SeBackupPrivilege, winio.SeRestorePrivilege***REMOVED***)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if noreexec ***REMOVED***
		defer func() ***REMOVED***
			if err := winio.DisableProcessPrivileges([]string***REMOVED***winio.SeBackupPrivilege, winio.SeRestorePrivilege***REMOVED***); err != nil ***REMOVED***
				// This should never happen, but just in case when in debugging mode.
				// See https://github.com/docker/docker/pull/28002#discussion_r86259241 for rationale.
				panic("Failed to disabled process privileges while in non re-exec mode")
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	info := hcsshim.DriverInfo***REMOVED***
		Flavour: filterDriver,
		HomeDir: home,
	***REMOVED***

	w, err := hcsshim.NewLayerWriter(info, id, parentLayerPaths)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	size, err := writeLayerFromTar(layerData, w, filepath.Join(home, id))
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	err = w.Close()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return size, nil
***REMOVED***

// resolveID computes the layerID information based on the given id.
func (d *Driver) resolveID(id string) (string, error) ***REMOVED***
	content, err := ioutil.ReadFile(filepath.Join(d.dir(id), "layerID"))
	if os.IsNotExist(err) ***REMOVED***
		return id, nil
	***REMOVED*** else if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return string(content), nil
***REMOVED***

// setID stores the layerId in disk.
func (d *Driver) setID(id, altID string) error ***REMOVED***
	return ioutil.WriteFile(filepath.Join(d.dir(id), "layerId"), []byte(altID), 0600)
***REMOVED***

// getLayerChain returns the layer chain information.
func (d *Driver) getLayerChain(id string) ([]string, error) ***REMOVED***
	jPath := filepath.Join(d.dir(id), "layerchain.json")
	content, err := ioutil.ReadFile(jPath)
	if os.IsNotExist(err) ***REMOVED***
		return nil, nil
	***REMOVED*** else if err != nil ***REMOVED***
		return nil, fmt.Errorf("Unable to read layerchain file - %s", err)
	***REMOVED***

	var layerChain []string
	err = json.Unmarshal(content, &layerChain)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Failed to unmarshall layerchain json - %s", err)
	***REMOVED***

	return layerChain, nil
***REMOVED***

// setLayerChain stores the layer chain information in disk.
func (d *Driver) setLayerChain(id string, chain []string) error ***REMOVED***
	content, err := json.Marshal(&chain)
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to marshall layerchain json - %s", err)
	***REMOVED***

	jPath := filepath.Join(d.dir(id), "layerchain.json")
	err = ioutil.WriteFile(jPath, content, 0600)
	if err != nil ***REMOVED***
		return fmt.Errorf("Unable to write layerchain file - %s", err)
	***REMOVED***

	return nil
***REMOVED***

type fileGetCloserWithBackupPrivileges struct ***REMOVED***
	path string
***REMOVED***

func (fg *fileGetCloserWithBackupPrivileges) Get(filename string) (io.ReadCloser, error) ***REMOVED***
	if backupPath, ok := mutatedFiles[filename]; ok ***REMOVED***
		return os.Open(filepath.Join(fg.path, backupPath))
	***REMOVED***

	var f *os.File
	// Open the file while holding the Windows backup privilege. This ensures that the
	// file can be opened even if the caller does not actually have access to it according
	// to the security descriptor. Also use sequential file access to avoid depleting the
	// standby list - Microsoft VSO Bug Tracker #9900466
	err := winio.RunWithPrivilege(winio.SeBackupPrivilege, func() error ***REMOVED***
		path := longpath.AddPrefix(filepath.Join(fg.path, filename))
		p, err := windows.UTF16FromString(path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		const fileFlagSequentialScan = 0x08000000 // FILE_FLAG_SEQUENTIAL_SCAN
		h, err := windows.CreateFile(&p[0], windows.GENERIC_READ, windows.FILE_SHARE_READ, nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS|fileFlagSequentialScan, 0)
		if err != nil ***REMOVED***
			return &os.PathError***REMOVED***Op: "open", Path: path, Err: err***REMOVED***
		***REMOVED***
		f = os.NewFile(uintptr(h), path)
		return nil
	***REMOVED***)
	return f, err
***REMOVED***

func (fg *fileGetCloserWithBackupPrivileges) Close() error ***REMOVED***
	return nil
***REMOVED***

// DiffGetter returns a FileGetCloser that can read files from the directory that
// contains files for the layer differences. Used for direct access for tar-split.
func (d *Driver) DiffGetter(id string) (graphdriver.FileGetCloser, error) ***REMOVED***
	id, err := d.resolveID(id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &fileGetCloserWithBackupPrivileges***REMOVED***d.dir(id)***REMOVED***, nil
***REMOVED***

type storageOptions struct ***REMOVED***
	size uint64
***REMOVED***

func parseStorageOpt(storageOpt map[string]string) (*storageOptions, error) ***REMOVED***
	options := storageOptions***REMOVED******REMOVED***

	// Read size to change the block device size per container.
	for key, val := range storageOpt ***REMOVED***
		key := strings.ToLower(key)
		switch key ***REMOVED***
		case "size":
			size, err := units.RAMInBytes(val)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			options.size = uint64(size)
		***REMOVED***
	***REMOVED***
	return &options, nil
***REMOVED***
