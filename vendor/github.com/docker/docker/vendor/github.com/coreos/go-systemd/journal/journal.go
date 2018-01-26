// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package journal provides write bindings to the local systemd journal.
// It is implemented in pure Go and connects to the journal directly over its
// unix socket.
//
// To read from the journal, see the "sdjournal" package, which wraps the
// sd-journal a C API.
//
// http://www.freedesktop.org/software/systemd/man/systemd-journald.service.html
package journal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// Priority of a journal message
type Priority int

const (
	PriEmerg Priority = iota
	PriAlert
	PriCrit
	PriErr
	PriWarning
	PriNotice
	PriInfo
	PriDebug
)

var conn net.Conn

func init() ***REMOVED***
	var err error
	conn, err = net.Dial("unixgram", "/run/systemd/journal/socket")
	if err != nil ***REMOVED***
		conn = nil
	***REMOVED***
***REMOVED***

// Enabled returns true if the local systemd journal is available for logging
func Enabled() bool ***REMOVED***
	return conn != nil
***REMOVED***

// Send a message to the local systemd journal. vars is a map of journald
// fields to values.  Fields must be composed of uppercase letters, numbers,
// and underscores, but must not start with an underscore. Within these
// restrictions, any arbitrary field name may be used.  Some names have special
// significance: see the journalctl documentation
// (http://www.freedesktop.org/software/systemd/man/systemd.journal-fields.html)
// for more details.  vars may be nil.
func Send(message string, priority Priority, vars map[string]string) error ***REMOVED***
	if conn == nil ***REMOVED***
		return journalError("could not connect to journald socket")
	***REMOVED***

	data := new(bytes.Buffer)
	appendVariable(data, "PRIORITY", strconv.Itoa(int(priority)))
	appendVariable(data, "MESSAGE", message)
	for k, v := range vars ***REMOVED***
		appendVariable(data, k, v)
	***REMOVED***

	_, err := io.Copy(conn, data)
	if err != nil && isSocketSpaceError(err) ***REMOVED***
		file, err := tempFd()
		if err != nil ***REMOVED***
			return journalError(err.Error())
		***REMOVED***
		defer file.Close()
		_, err = io.Copy(file, data)
		if err != nil ***REMOVED***
			return journalError(err.Error())
		***REMOVED***

		rights := syscall.UnixRights(int(file.Fd()))

		/* this connection should always be a UnixConn, but better safe than sorry */
		unixConn, ok := conn.(*net.UnixConn)
		if !ok ***REMOVED***
			return journalError("can't send file through non-Unix connection")
		***REMOVED***
		unixConn.WriteMsgUnix([]byte***REMOVED******REMOVED***, rights, nil)
	***REMOVED*** else if err != nil ***REMOVED***
		return journalError(err.Error())
	***REMOVED***
	return nil
***REMOVED***

// Print prints a message to the local systemd journal using Send().
func Print(priority Priority, format string, a ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return Send(fmt.Sprintf(format, a...), priority, nil)
***REMOVED***

func appendVariable(w io.Writer, name, value string) ***REMOVED***
	if !validVarName(name) ***REMOVED***
		journalError("variable name contains invalid character, ignoring")
	***REMOVED***
	if strings.ContainsRune(value, '\n') ***REMOVED***
		/* When the value contains a newline, we write:
		 * - the variable name, followed by a newline
		 * - the size (in 64bit little endian format)
		 * - the data, followed by a newline
		 */
		fmt.Fprintln(w, name)
		binary.Write(w, binary.LittleEndian, uint64(len(value)))
		fmt.Fprintln(w, value)
	***REMOVED*** else ***REMOVED***
		/* just write the variable and value all on one line */
		fmt.Fprintf(w, "%s=%s\n", name, value)
	***REMOVED***
***REMOVED***

func validVarName(name string) bool ***REMOVED***
	/* The variable name must be in uppercase and consist only of characters,
	 * numbers and underscores, and may not begin with an underscore. (from the docs)
	 */

	valid := name[0] != '_'
	for _, c := range name ***REMOVED***
		valid = valid && ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || c == '_'
	***REMOVED***
	return valid
***REMOVED***

func isSocketSpaceError(err error) bool ***REMOVED***
	opErr, ok := err.(*net.OpError)
	if !ok ***REMOVED***
		return false
	***REMOVED***

	sysErr, ok := opErr.Err.(syscall.Errno)
	if !ok ***REMOVED***
		return false
	***REMOVED***

	return sysErr == syscall.EMSGSIZE || sysErr == syscall.ENOBUFS
***REMOVED***

func tempFd() (*os.File, error) ***REMOVED***
	file, err := ioutil.TempFile("/dev/shm/", "journal.XXXXX")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	syscall.Unlink(file.Name())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return file, nil
***REMOVED***

func journalError(s string) error ***REMOVED***
	s = "journal error: " + s
	fmt.Fprintln(os.Stderr, s)
	return errors.New(s)
***REMOVED***
