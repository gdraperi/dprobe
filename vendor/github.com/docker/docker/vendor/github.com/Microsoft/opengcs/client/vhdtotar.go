// +build windows

package client

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// VhdToTar does what is says - it exports a VHD in a specified
// folder (either a read-only layer.vhd, or a read-write sandbox.vhd) to a
// ReadCloser containing a tar-stream of the layers contents.
func (config *Config) VhdToTar(vhdFile string, uvmMountPath string, isSandbox bool, vhdSize int64) (io.ReadCloser, error) ***REMOVED***
	logrus.Debugf("opengcs: VhdToTar: %s isSandbox: %t", vhdFile, isSandbox)

	if config.Uvm == nil ***REMOVED***
		return nil, fmt.Errorf("cannot VhdToTar as no utility VM is in configuration")
	***REMOVED***

	defer config.DebugGCS()

	vhdHandle, err := os.Open(vhdFile)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("opengcs: VhdToTar: failed to open %s: %s", vhdFile, err)
	***REMOVED***
	defer vhdHandle.Close()
	logrus.Debugf("opengcs: VhdToTar: exporting %s, size %d, isSandbox %t", vhdHandle.Name(), vhdSize, isSandbox)

	// Different binary depending on whether a RO layer or a RW sandbox
	command := "vhd2tar"
	if isSandbox ***REMOVED***
		command = fmt.Sprintf("exportSandbox -path %s", uvmMountPath)
	***REMOVED***

	// Start the binary in the utility VM
	process, err := config.createUtilsProcess(command)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("opengcs: VhdToTar: %s: failed to create utils process %s: %s", vhdHandle.Name(), command, err)
	***REMOVED***

	if !isSandbox ***REMOVED***
		// Send the VHD contents to the utility VM processes stdin handle if not a sandbox
		logrus.Debugf("opengcs: VhdToTar: copying the layer VHD into the utility VM")
		if _, err = copyWithTimeout(process.Stdin, vhdHandle, vhdSize, config.UvmTimeoutSeconds, fmt.Sprintf("vhdtotarstream: sending %s to %s", vhdHandle.Name(), command)); err != nil ***REMOVED***
			process.Process.Close()
			return nil, fmt.Errorf("opengcs: VhdToTar: %s: failed to copyWithTimeout on the stdin pipe (to utility VM): %s", vhdHandle.Name(), err)
		***REMOVED***
	***REMOVED***

	// Start a goroutine which copies the stdout (ie the tar stream)
	reader, writer := io.Pipe()
	go func() ***REMOVED***
		defer writer.Close()
		defer process.Process.Close()
		logrus.Debugf("opengcs: VhdToTar: copying tar stream back from the utility VM")
		bytes, err := copyWithTimeout(writer, process.Stdout, vhdSize, config.UvmTimeoutSeconds, fmt.Sprintf("vhdtotarstream: copy tarstream from %s", command))
		if err != nil ***REMOVED***
			logrus.Errorf("opengcs: VhdToTar: %s:  copyWithTimeout on the stdout pipe (from utility VM) failed: %s", vhdHandle.Name(), err)
		***REMOVED***
		logrus.Debugf("opengcs: VhdToTar: copied %d bytes of the tarstream of %s from the utility VM", bytes, vhdHandle.Name())
	***REMOVED***()

	// Return the read-side of the pipe connected to the goroutine which is reading from the stdout of the process in the utility VM
	return reader, nil

***REMOVED***
