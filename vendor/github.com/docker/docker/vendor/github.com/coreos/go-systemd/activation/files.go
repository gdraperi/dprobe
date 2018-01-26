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

// Package activation implements primitives for systemd socket activation.
package activation

import (
	"os"
	"strconv"
	"syscall"
)

// based on: https://gist.github.com/alberts/4640792
const (
	listenFdsStart = 3
)

func Files(unsetEnv bool) []*os.File ***REMOVED***
	if unsetEnv ***REMOVED***
		defer os.Unsetenv("LISTEN_PID")
		defer os.Unsetenv("LISTEN_FDS")
	***REMOVED***

	pid, err := strconv.Atoi(os.Getenv("LISTEN_PID"))
	if err != nil || pid != os.Getpid() ***REMOVED***
		return nil
	***REMOVED***

	nfds, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
	if err != nil || nfds == 0 ***REMOVED***
		return nil
	***REMOVED***

	files := make([]*os.File, 0, nfds)
	for fd := listenFdsStart; fd < listenFdsStart+nfds; fd++ ***REMOVED***
		syscall.CloseOnExec(fd)
		files = append(files, os.NewFile(uintptr(fd), "LISTEN_FD_"+strconv.Itoa(fd)))
	***REMOVED***

	return files
***REMOVED***
