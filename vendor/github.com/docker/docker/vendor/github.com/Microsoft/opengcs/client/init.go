// +build windows

package client

import (
	"os"
	"strconv"
)

var (
	logDataFromUVM int64
)

func init() ***REMOVED***
	bytes := os.Getenv("OPENGCS_LOG_DATA_FROM_UVM")
	if len(bytes) == 0 ***REMOVED***
		return
	***REMOVED***
	u, err := strconv.ParseUint(bytes, 10, 32)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	logDataFromUVM = int64(u)
***REMOVED***
