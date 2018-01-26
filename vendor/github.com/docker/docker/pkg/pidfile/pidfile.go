// Package pidfile provides structure and helper functions to create and remove
// PID file. A PID file is usually a file used to store the process ID of a
// running process.
package pidfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/pkg/system"
)

// PIDFile is a file used to store the process ID of a running process.
type PIDFile struct ***REMOVED***
	path string
***REMOVED***

func checkPIDFileAlreadyExists(path string) error ***REMOVED***
	if pidByte, err := ioutil.ReadFile(path); err == nil ***REMOVED***
		pidString := strings.TrimSpace(string(pidByte))
		if pid, err := strconv.Atoi(pidString); err == nil ***REMOVED***
			if processExists(pid) ***REMOVED***
				return fmt.Errorf("pid file found, ensure docker is not running or delete %s", path)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// New creates a PIDfile using the specified path.
func New(path string) (*PIDFile, error) ***REMOVED***
	if err := checkPIDFileAlreadyExists(path); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Note MkdirAll returns nil if a directory already exists
	if err := system.MkdirAll(filepath.Dir(path), os.FileMode(0755), ""); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := ioutil.WriteFile(path, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &PIDFile***REMOVED***path: path***REMOVED***, nil
***REMOVED***

// Remove removes the PIDFile.
func (file PIDFile) Remove() error ***REMOVED***
	return os.Remove(file.path)
***REMOVED***
