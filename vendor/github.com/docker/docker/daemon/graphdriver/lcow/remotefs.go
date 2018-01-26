// +build windows

package lcow

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"

	"github.com/Microsoft/hcsshim"
	"github.com/Microsoft/opengcs/service/gcsutils/remotefs"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/sirupsen/logrus"
)

type lcowfs struct ***REMOVED***
	root        string
	d           *Driver
	mappedDisks []hcsshim.MappedVirtualDisk
	vmID        string
	currentSVM  *serviceVM
	sync.Mutex
***REMOVED***

var _ containerfs.ContainerFS = &lcowfs***REMOVED******REMOVED***

// ErrNotSupported is an error for unsupported operations in the remotefs
var ErrNotSupported = fmt.Errorf("not supported")

// Functions to implement the ContainerFS interface
func (l *lcowfs) Path() string ***REMOVED***
	return l.root
***REMOVED***

func (l *lcowfs) ResolveScopedPath(path string, rawPath bool) (string, error) ***REMOVED***
	logrus.Debugf("remotefs.resolvescopedpath inputs: %s %s ", path, l.root)

	arg1 := l.Join(l.root, path)
	if !rawPath ***REMOVED***
		// The l.Join("/", path) will make path an absolute path and then clean it
		// so if path = ../../X, it will become /X.
		arg1 = l.Join(l.root, l.Join("/", path))
	***REMOVED***
	arg2 := l.root

	output := &bytes.Buffer***REMOVED******REMOVED***
	if err := l.runRemoteFSProcess(nil, output, remotefs.ResolvePathCmd, arg1, arg2); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	logrus.Debugf("remotefs.resolvescopedpath success. Output: %s\n", output.String())
	return output.String(), nil
***REMOVED***

func (l *lcowfs) OS() string ***REMOVED***
	return "linux"
***REMOVED***

func (l *lcowfs) Architecture() string ***REMOVED***
	return runtime.GOARCH
***REMOVED***

// Other functions that are used by docker like the daemon Archiver/Extractor
func (l *lcowfs) ExtractArchive(src io.Reader, dst string, opts *archive.TarOptions) error ***REMOVED***
	logrus.Debugf("remotefs.ExtractArchve inputs: %s %+v", dst, opts)

	tarBuf := &bytes.Buffer***REMOVED******REMOVED***
	if err := remotefs.WriteTarOptions(tarBuf, opts); err != nil ***REMOVED***
		return fmt.Errorf("failed to marshall tar opts: %s", err)
	***REMOVED***

	input := io.MultiReader(tarBuf, src)
	if err := l.runRemoteFSProcess(input, nil, remotefs.ExtractArchiveCmd, dst); err != nil ***REMOVED***
		return fmt.Errorf("failed to extract archive to %s: %s", dst, err)
	***REMOVED***
	return nil
***REMOVED***

func (l *lcowfs) ArchivePath(src string, opts *archive.TarOptions) (io.ReadCloser, error) ***REMOVED***
	logrus.Debugf("remotefs.ArchivePath: %s %+v", src, opts)

	tarBuf := &bytes.Buffer***REMOVED******REMOVED***
	if err := remotefs.WriteTarOptions(tarBuf, opts); err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to marshall tar opts: %s", err)
	***REMOVED***

	r, w := io.Pipe()
	go func() ***REMOVED***
		defer w.Close()
		if err := l.runRemoteFSProcess(tarBuf, w, remotefs.ArchivePathCmd, src); err != nil ***REMOVED***
			logrus.Debugf("REMOTEFS: Failed to extract archive: %s %+v %s", src, opts, err)
		***REMOVED***
	***REMOVED***()
	return r, nil
***REMOVED***

// Helper functions
func (l *lcowfs) startVM() error ***REMOVED***
	l.Lock()
	defer l.Unlock()
	if l.currentSVM != nil ***REMOVED***
		return nil
	***REMOVED***

	svm, err := l.d.startServiceVMIfNotRunning(l.vmID, l.mappedDisks, fmt.Sprintf("lcowfs.startVM"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = svm.createUnionMount(l.root, l.mappedDisks...); err != nil ***REMOVED***
		return err
	***REMOVED***
	l.currentSVM = svm
	return nil
***REMOVED***

func (l *lcowfs) runRemoteFSProcess(stdin io.Reader, stdout io.Writer, args ...string) error ***REMOVED***
	if err := l.startVM(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Append remotefs prefix and setup as a command line string
	cmd := fmt.Sprintf("%s %s", remotefs.RemotefsCmd, strings.Join(args, " "))
	stderr := &bytes.Buffer***REMOVED******REMOVED***
	if err := l.currentSVM.runProcess(cmd, stdin, stdout, stderr); err != nil ***REMOVED***
		return err
	***REMOVED***

	eerr, err := remotefs.ReadError(stderr)
	if eerr != nil ***REMOVED***
		// Process returned an error so return that.
		return remotefs.ExportedToError(eerr)
	***REMOVED***
	return err
***REMOVED***
