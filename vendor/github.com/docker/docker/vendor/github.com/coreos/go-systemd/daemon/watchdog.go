// Copyright 2016 CoreOS, Inc.
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

package daemon

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// SdWatchdogEnabled return watchdog information for a service.
// Process should send daemon.SdNotify("WATCHDOG=1") every time / 2.
// If `unsetEnvironment` is true, the environment variables `WATCHDOG_USEC`
// and `WATCHDOG_PID` will be unconditionally unset.
//
// It returns one of the following:
// (0, nil) - watchdog isn't enabled or we aren't the watched PID.
// (0, err) - an error happened (e.g. error converting time).
// (time, nil) - watchdog is enabled and we can send ping.
//   time is delay before inactive service will be killed.
func SdWatchdogEnabled(unsetEnvironment bool) (time.Duration, error) ***REMOVED***
	wusec := os.Getenv("WATCHDOG_USEC")
	wpid := os.Getenv("WATCHDOG_PID")
	if unsetEnvironment ***REMOVED***
		wusecErr := os.Unsetenv("WATCHDOG_USEC")
		wpidErr := os.Unsetenv("WATCHDOG_PID")
		if wusecErr != nil ***REMOVED***
			return 0, wusecErr
		***REMOVED***
		if wpidErr != nil ***REMOVED***
			return 0, wpidErr
		***REMOVED***
	***REMOVED***

	if wusec == "" ***REMOVED***
		return 0, nil
	***REMOVED***
	s, err := strconv.Atoi(wusec)
	if err != nil ***REMOVED***
		return 0, fmt.Errorf("error converting WATCHDOG_USEC: %s", err)
	***REMOVED***
	if s <= 0 ***REMOVED***
		return 0, fmt.Errorf("error WATCHDOG_USEC must be a positive number")
	***REMOVED***
	interval := time.Duration(s) * time.Microsecond

	if wpid == "" ***REMOVED***
		return interval, nil
	***REMOVED***
	p, err := strconv.Atoi(wpid)
	if err != nil ***REMOVED***
		return 0, fmt.Errorf("error converting WATCHDOG_PID: %s", err)
	***REMOVED***
	if os.Getpid() != p ***REMOVED***
		return 0, nil
	***REMOVED***

	return interval, nil
***REMOVED***
