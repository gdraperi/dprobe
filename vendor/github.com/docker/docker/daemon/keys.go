// +build linux

package daemon

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	rootKeyFile   = "/proc/sys/kernel/keys/root_maxkeys"
	rootBytesFile = "/proc/sys/kernel/keys/root_maxbytes"
	rootKeyLimit  = 1000000
	// it is standard configuration to allocate 25 bytes per key
	rootKeyByteMultiplier = 25
)

// ModifyRootKeyLimit checks to see if the root key limit is set to
// at least 1000000 and changes it to that limit along with the maxbytes
// allocated to the keys at a 25 to 1 multiplier.
func ModifyRootKeyLimit() error ***REMOVED***
	value, err := readRootKeyLimit(rootKeyFile)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if value < rootKeyLimit ***REMOVED***
		return setRootKeyLimit(rootKeyLimit)
	***REMOVED***
	return nil
***REMOVED***

func setRootKeyLimit(limit int) error ***REMOVED***
	keys, err := os.OpenFile(rootKeyFile, os.O_WRONLY, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer keys.Close()
	if _, err := fmt.Fprintf(keys, "%d", limit); err != nil ***REMOVED***
		return err
	***REMOVED***
	bytes, err := os.OpenFile(rootBytesFile, os.O_WRONLY, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer bytes.Close()
	_, err = fmt.Fprintf(bytes, "%d", limit*rootKeyByteMultiplier)
	return err
***REMOVED***

func readRootKeyLimit(path string) (int, error) ***REMOVED***
	data, err := ioutil.ReadFile(path)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	return strconv.Atoi(strings.Trim(string(data), "\n"))
***REMOVED***
