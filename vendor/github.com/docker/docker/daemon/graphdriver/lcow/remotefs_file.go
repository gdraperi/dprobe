// +build windows

package lcow

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/Microsoft/hcsshim"
	"github.com/Microsoft/opengcs/service/gcsutils/remotefs"
	"github.com/containerd/continuity/driver"
)

type lcowfile struct ***REMOVED***
	process   hcsshim.Process
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	fs        *lcowfs
	guestPath string
***REMOVED***

func (l *lcowfs) Open(path string) (driver.File, error) ***REMOVED***
	return l.OpenFile(path, os.O_RDONLY, 0)
***REMOVED***

func (l *lcowfs) OpenFile(path string, flag int, perm os.FileMode) (_ driver.File, err error) ***REMOVED***
	flagStr := strconv.FormatInt(int64(flag), 10)
	permStr := strconv.FormatUint(uint64(perm), 8)

	commandLine := fmt.Sprintf("%s %s %s %s %s", remotefs.RemotefsCmd, remotefs.OpenFileCmd, path, flagStr, permStr)
	env := make(map[string]string)
	env["PATH"] = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:"
	processConfig := &hcsshim.ProcessConfig***REMOVED***
		EmulateConsole:    false,
		CreateStdInPipe:   true,
		CreateStdOutPipe:  true,
		CreateStdErrPipe:  true,
		CreateInUtilityVm: true,
		WorkingDirectory:  "/bin",
		Environment:       env,
		CommandLine:       commandLine,
	***REMOVED***

	process, err := l.currentSVM.config.Uvm.CreateProcess(processConfig)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to open file %s: %s", path, err)
	***REMOVED***

	stdin, stdout, stderr, err := process.Stdio()
	if err != nil ***REMOVED***
		process.Kill()
		process.Close()
		return nil, fmt.Errorf("failed to open file pipes %s: %s", path, err)
	***REMOVED***

	lf := &lcowfile***REMOVED***
		process:   process,
		stdin:     stdin,
		stdout:    stdout,
		stderr:    stderr,
		fs:        l,
		guestPath: path,
	***REMOVED***

	if _, err := lf.getResponse(); err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to open file %s: %s", path, err)
	***REMOVED***
	return lf, nil
***REMOVED***

func (l *lcowfile) Read(b []byte) (int, error) ***REMOVED***
	hdr := &remotefs.FileHeader***REMOVED***
		Cmd:  remotefs.Read,
		Size: uint64(len(b)),
	***REMOVED***

	if err := remotefs.WriteFileHeader(l.stdin, hdr, nil); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	buf, err := l.getResponse()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	n := copy(b, buf)
	return n, nil
***REMOVED***

func (l *lcowfile) Write(b []byte) (int, error) ***REMOVED***
	hdr := &remotefs.FileHeader***REMOVED***
		Cmd:  remotefs.Write,
		Size: uint64(len(b)),
	***REMOVED***

	if err := remotefs.WriteFileHeader(l.stdin, hdr, b); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	_, err := l.getResponse()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return len(b), nil
***REMOVED***

func (l *lcowfile) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	seekHdr := &remotefs.SeekHeader***REMOVED***
		Offset: offset,
		Whence: int32(whence),
	***REMOVED***

	buf := &bytes.Buffer***REMOVED******REMOVED***
	if err := binary.Write(buf, binary.BigEndian, seekHdr); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	hdr := &remotefs.FileHeader***REMOVED***
		Cmd:  remotefs.Write,
		Size: uint64(buf.Len()),
	***REMOVED***
	if err := remotefs.WriteFileHeader(l.stdin, hdr, buf.Bytes()); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	resBuf, err := l.getResponse()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	var res int64
	if err := binary.Read(bytes.NewBuffer(resBuf), binary.BigEndian, &res); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return res, nil
***REMOVED***

func (l *lcowfile) Close() error ***REMOVED***
	hdr := &remotefs.FileHeader***REMOVED***
		Cmd:  remotefs.Close,
		Size: 0,
	***REMOVED***

	if err := remotefs.WriteFileHeader(l.stdin, hdr, nil); err != nil ***REMOVED***
		return err
	***REMOVED***

	_, err := l.getResponse()
	return err
***REMOVED***

func (l *lcowfile) Readdir(n int) ([]os.FileInfo, error) ***REMOVED***
	nStr := strconv.FormatInt(int64(n), 10)

	// Unlike the other File functions, this one can just be run without maintaining state,
	// so just do the normal runRemoteFSProcess way.
	buf := &bytes.Buffer***REMOVED******REMOVED***
	if err := l.fs.runRemoteFSProcess(nil, buf, remotefs.ReadDirCmd, l.guestPath, nStr); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var info []remotefs.FileInfo
	if err := json.Unmarshal(buf.Bytes(), &info); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	osInfo := make([]os.FileInfo, len(info))
	for i := range info ***REMOVED***
		osInfo[i] = &info[i]
	***REMOVED***
	return osInfo, nil
***REMOVED***

func (l *lcowfile) getResponse() ([]byte, error) ***REMOVED***
	hdr, err := remotefs.ReadFileHeader(l.stdout)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if hdr.Cmd != remotefs.CmdOK ***REMOVED***
		// Something went wrong during the openfile in the server.
		// Parse stderr and return that as an error
		eerr, err := remotefs.ReadError(l.stderr)
		if eerr != nil ***REMOVED***
			return nil, remotefs.ExportedToError(eerr)
		***REMOVED***

		// Maybe the parsing went wrong?
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		// At this point, we know something went wrong in the remotefs program, but
		// we we don't know why.
		return nil, fmt.Errorf("unknown error")
	***REMOVED***

	// Successful command, we might have some data to read (for Read + Seek)
	buf := make([]byte, hdr.Size, hdr.Size)
	if _, err := io.ReadFull(l.stdout, buf); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return buf, nil
***REMOVED***
