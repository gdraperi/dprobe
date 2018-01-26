package hcsshim

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Microsoft/go-winio"
	"github.com/sirupsen/logrus"
)

// ImportLayer will take the contents of the folder at importFolderPath and import
// that into a layer with the id layerId.  Note that in order to correctly populate
// the layer and interperet the transport format, all parent layers must already
// be present on the system at the paths provided in parentLayerPaths.
func ImportLayer(info DriverInfo, layerID string, importFolderPath string, parentLayerPaths []string) error ***REMOVED***
	title := "hcsshim::ImportLayer "
	logrus.Debugf(title+"flavour %d layerId %s folder %s", info.Flavour, layerID, importFolderPath)

	// Generate layer descriptors
	layers, err := layerPathsToDescriptors(parentLayerPaths)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Convert info to API calling convention
	infop, err := convertDriverInfo(info)
	if err != nil ***REMOVED***
		logrus.Error(err)
		return err
	***REMOVED***

	err = importLayer(&infop, layerID, importFolderPath, layers)
	if err != nil ***REMOVED***
		err = makeErrorf(err, title, "layerId=%s flavour=%d folder=%s", layerID, info.Flavour, importFolderPath)
		logrus.Error(err)
		return err
	***REMOVED***

	logrus.Debugf(title+"succeeded flavour=%d layerId=%s folder=%s", info.Flavour, layerID, importFolderPath)
	return nil
***REMOVED***

// LayerWriter is an interface that supports writing a new container image layer.
type LayerWriter interface ***REMOVED***
	// Add adds a file to the layer with given metadata.
	Add(name string, fileInfo *winio.FileBasicInfo) error
	// AddLink adds a hard link to the layer. The target must already have been added.
	AddLink(name string, target string) error
	// Remove removes a file that was present in a parent layer from the layer.
	Remove(name string) error
	// Write writes data to the current file. The data must be in the format of a Win32
	// backup stream.
	Write(b []byte) (int, error)
	// Close finishes the layer writing process and releases any resources.
	Close() error
***REMOVED***

// FilterLayerWriter provides an interface to write the contents of a layer to the file system.
type FilterLayerWriter struct ***REMOVED***
	context uintptr
***REMOVED***

// Add adds a file or directory to the layer. The file's parent directory must have already been added.
//
// name contains the file's relative path. fileInfo contains file times and file attributes; the rest
// of the file metadata and the file data must be written as a Win32 backup stream to the Write() method.
// winio.BackupStreamWriter can be used to facilitate this.
func (w *FilterLayerWriter) Add(name string, fileInfo *winio.FileBasicInfo) error ***REMOVED***
	if name[0] != '\\' ***REMOVED***
		name = `\` + name
	***REMOVED***
	err := importLayerNext(w.context, name, fileInfo)
	if err != nil ***REMOVED***
		return makeError(err, "ImportLayerNext", "")
	***REMOVED***
	return nil
***REMOVED***

// AddLink adds a hard link to the layer. The target of the link must have already been added.
func (w *FilterLayerWriter) AddLink(name string, target string) error ***REMOVED***
	return errors.New("hard links not yet supported")
***REMOVED***

// Remove removes a file from the layer. The file must have been present in the parent layer.
//
// name contains the file's relative path.
func (w *FilterLayerWriter) Remove(name string) error ***REMOVED***
	if name[0] != '\\' ***REMOVED***
		name = `\` + name
	***REMOVED***
	err := importLayerNext(w.context, name, nil)
	if err != nil ***REMOVED***
		return makeError(err, "ImportLayerNext", "")
	***REMOVED***
	return nil
***REMOVED***

// Write writes more backup stream data to the current file.
func (w *FilterLayerWriter) Write(b []byte) (int, error) ***REMOVED***
	err := importLayerWrite(w.context, b)
	if err != nil ***REMOVED***
		err = makeError(err, "ImportLayerWrite", "")
		return 0, err
	***REMOVED***
	return len(b), err
***REMOVED***

// Close completes the layer write operation. The error must be checked to ensure that the
// operation was successful.
func (w *FilterLayerWriter) Close() (err error) ***REMOVED***
	if w.context != 0 ***REMOVED***
		err = importLayerEnd(w.context)
		if err != nil ***REMOVED***
			err = makeError(err, "ImportLayerEnd", "")
		***REMOVED***
		w.context = 0
	***REMOVED***
	return
***REMOVED***

type legacyLayerWriterWrapper struct ***REMOVED***
	*legacyLayerWriter
	info             DriverInfo
	layerID          string
	path             string
	parentLayerPaths []string
***REMOVED***

func (r *legacyLayerWriterWrapper) Close() error ***REMOVED***
	defer os.RemoveAll(r.root)
	err := r.legacyLayerWriter.Close()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Use the original path here because ImportLayer does not support long paths for the source in TP5.
	// But do use a long path for the destination to work around another bug with directories
	// with MAX_PATH - 12 < length < MAX_PATH.
	info := r.info
	fullPath, err := makeLongAbsPath(filepath.Join(info.HomeDir, r.layerID))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	info.HomeDir = ""
	if err = ImportLayer(info, fullPath, r.path, r.parentLayerPaths); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Add any hard links that were collected.
	for _, lnk := range r.PendingLinks ***REMOVED***
		if err = os.Remove(lnk.Path); err != nil && !os.IsNotExist(err) ***REMOVED***
			return err
		***REMOVED***
		if err = os.Link(lnk.Target, lnk.Path); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	// Prepare the utility VM for use if one is present in the layer.
	if r.HasUtilityVM ***REMOVED***
		err = ProcessUtilityVMImage(filepath.Join(fullPath, "UtilityVM"))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// NewLayerWriter returns a new layer writer for creating a layer on disk.
// The caller must have taken the SeBackupPrivilege and SeRestorePrivilege privileges
// to call this and any methods on the resulting LayerWriter.
func NewLayerWriter(info DriverInfo, layerID string, parentLayerPaths []string) (LayerWriter, error) ***REMOVED***
	if len(parentLayerPaths) == 0 ***REMOVED***
		// This is a base layer. It gets imported differently.
		return &baseLayerWriter***REMOVED***
			root: filepath.Join(info.HomeDir, layerID),
		***REMOVED***, nil
	***REMOVED***

	if procImportLayerBegin.Find() != nil ***REMOVED***
		// The new layer reader is not available on this Windows build. Fall back to the
		// legacy export code path.
		path, err := ioutil.TempDir("", "hcs")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return &legacyLayerWriterWrapper***REMOVED***
			legacyLayerWriter: newLegacyLayerWriter(path, parentLayerPaths, filepath.Join(info.HomeDir, layerID)),
			info:              info,
			layerID:           layerID,
			path:              path,
			parentLayerPaths:  parentLayerPaths,
		***REMOVED***, nil
	***REMOVED***
	layers, err := layerPathsToDescriptors(parentLayerPaths)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	infop, err := convertDriverInfo(info)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	w := &FilterLayerWriter***REMOVED******REMOVED***
	err = importLayerBegin(&infop, layerID, layers, &w.context)
	if err != nil ***REMOVED***
		return nil, makeError(err, "ImportLayerStart", "")
	***REMOVED***
	return w, nil
***REMOVED***
