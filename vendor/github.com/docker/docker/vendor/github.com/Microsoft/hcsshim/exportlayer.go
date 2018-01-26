package hcsshim

import (
	"io"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/Microsoft/go-winio"
	"github.com/sirupsen/logrus"
)

// ExportLayer will create a folder at exportFolderPath and fill that folder with
// the transport format version of the layer identified by layerId. This transport
// format includes any metadata required for later importing the layer (using
// ImportLayer), and requires the full list of parent layer paths in order to
// perform the export.
func ExportLayer(info DriverInfo, layerId string, exportFolderPath string, parentLayerPaths []string) error ***REMOVED***
	title := "hcsshim::ExportLayer "
	logrus.Debugf(title+"flavour %d layerId %s folder %s", info.Flavour, layerId, exportFolderPath)

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

	err = exportLayer(&infop, layerId, exportFolderPath, layers)
	if err != nil ***REMOVED***
		err = makeErrorf(err, title, "layerId=%s flavour=%d folder=%s", layerId, info.Flavour, exportFolderPath)
		logrus.Error(err)
		return err
	***REMOVED***

	logrus.Debugf(title+"succeeded flavour=%d layerId=%s folder=%s", info.Flavour, layerId, exportFolderPath)
	return nil
***REMOVED***

type LayerReader interface ***REMOVED***
	Next() (string, int64, *winio.FileBasicInfo, error)
	Read(b []byte) (int, error)
	Close() error
***REMOVED***

// FilterLayerReader provides an interface for extracting the contents of an on-disk layer.
type FilterLayerReader struct ***REMOVED***
	context uintptr
***REMOVED***

// Next reads the next available file from a layer, ensuring that parent directories are always read
// before child files and directories.
//
// Next returns the file's relative path, size, and basic file metadata. Read() should be used to
// extract a Win32 backup stream with the remainder of the metadata and the data.
func (r *FilterLayerReader) Next() (string, int64, *winio.FileBasicInfo, error) ***REMOVED***
	var fileNamep *uint16
	fileInfo := &winio.FileBasicInfo***REMOVED******REMOVED***
	var deleted uint32
	var fileSize int64
	err := exportLayerNext(r.context, &fileNamep, fileInfo, &fileSize, &deleted)
	if err != nil ***REMOVED***
		if err == syscall.ERROR_NO_MORE_FILES ***REMOVED***
			err = io.EOF
		***REMOVED*** else ***REMOVED***
			err = makeError(err, "ExportLayerNext", "")
		***REMOVED***
		return "", 0, nil, err
	***REMOVED***
	fileName := convertAndFreeCoTaskMemString(fileNamep)
	if deleted != 0 ***REMOVED***
		fileInfo = nil
	***REMOVED***
	if fileName[0] == '\\' ***REMOVED***
		fileName = fileName[1:]
	***REMOVED***
	return fileName, fileSize, fileInfo, nil
***REMOVED***

// Read reads from the current file's Win32 backup stream.
func (r *FilterLayerReader) Read(b []byte) (int, error) ***REMOVED***
	var bytesRead uint32
	err := exportLayerRead(r.context, b, &bytesRead)
	if err != nil ***REMOVED***
		return 0, makeError(err, "ExportLayerRead", "")
	***REMOVED***
	if bytesRead == 0 ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	return int(bytesRead), nil
***REMOVED***

// Close frees resources associated with the layer reader. It will return an
// error if there was an error while reading the layer or of the layer was not
// completely read.
func (r *FilterLayerReader) Close() (err error) ***REMOVED***
	if r.context != 0 ***REMOVED***
		err = exportLayerEnd(r.context)
		if err != nil ***REMOVED***
			err = makeError(err, "ExportLayerEnd", "")
		***REMOVED***
		r.context = 0
	***REMOVED***
	return
***REMOVED***

// NewLayerReader returns a new layer reader for reading the contents of an on-disk layer.
// The caller must have taken the SeBackupPrivilege privilege
// to call this and any methods on the resulting LayerReader.
func NewLayerReader(info DriverInfo, layerID string, parentLayerPaths []string) (LayerReader, error) ***REMOVED***
	if procExportLayerBegin.Find() != nil ***REMOVED***
		// The new layer reader is not available on this Windows build. Fall back to the
		// legacy export code path.
		path, err := ioutil.TempDir("", "hcs")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		err = ExportLayer(info, layerID, path, parentLayerPaths)
		if err != nil ***REMOVED***
			os.RemoveAll(path)
			return nil, err
		***REMOVED***
		return &legacyLayerReaderWrapper***REMOVED***newLegacyLayerReader(path)***REMOVED***, nil
	***REMOVED***

	layers, err := layerPathsToDescriptors(parentLayerPaths)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	infop, err := convertDriverInfo(info)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	r := &FilterLayerReader***REMOVED******REMOVED***
	err = exportLayerBegin(&infop, layerID, layers, &r.context)
	if err != nil ***REMOVED***
		return nil, makeError(err, "ExportLayerBegin", "")
	***REMOVED***
	return r, err
***REMOVED***

type legacyLayerReaderWrapper struct ***REMOVED***
	*legacyLayerReader
***REMOVED***

func (r *legacyLayerReaderWrapper) Close() error ***REMOVED***
	err := r.legacyLayerReader.Close()
	os.RemoveAll(r.root)
	return err
***REMOVED***
