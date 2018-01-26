// Copyright 2014 Docker, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Code forked from Docker project
package daemon

import (
	"net"
	"os"
)

// SdNotify sends a message to the init daemon. It is common to ignore the error.
// If `unsetEnvironment` is true, the environment variable `NOTIFY_SOCKET`
// will be unconditionally unset.
//
// It returns one of the following:
// (false, nil) - notification not supported (i.e. NOTIFY_SOCKET is unset)
// (false, err) - notification supported, but failure happened (e.g. error connecting to NOTIFY_SOCKET or while sending data)
// (true, nil) - notification supported, data has been sent
func SdNotify(unsetEnvironment bool, state string) (sent bool, err error) ***REMOVED***
	socketAddr := &net.UnixAddr***REMOVED***
		Name: os.Getenv("NOTIFY_SOCKET"),
		Net:  "unixgram",
	***REMOVED***

	// NOTIFY_SOCKET not set
	if socketAddr.Name == "" ***REMOVED***
		return false, nil
	***REMOVED***

	if unsetEnvironment ***REMOVED***
		err = os.Unsetenv("NOTIFY_SOCKET")
	***REMOVED***
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
	// Error connecting to NOTIFY_SOCKET
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	defer conn.Close()

	_, err = conn.Write([]byte(state))
	// Error sending the message
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return true, nil
***REMOVED***
