package cluster

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/ioutils"
)

func loadPersistentState(root string) (*nodeStartConfig, error) ***REMOVED***
	dt, err := ioutil.ReadFile(filepath.Join(root, stateFile))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// missing certificate means no actual state to restore from
	if _, err := os.Stat(filepath.Join(root, "certificates/swarm-node.crt")); err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			clearPersistentState(root)
		***REMOVED***
		return nil, err
	***REMOVED***
	var st nodeStartConfig
	if err := json.Unmarshal(dt, &st); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &st, nil
***REMOVED***

func savePersistentState(root string, config nodeStartConfig) error ***REMOVED***
	dt, err := json.Marshal(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutils.AtomicWriteFile(filepath.Join(root, stateFile), dt, 0600)
***REMOVED***

func clearPersistentState(root string) error ***REMOVED***
	// todo: backup this data instead of removing?
	// rather than delete the entire swarm directory, delete the contents in order to preserve the inode
	// (for example, allowing it to be bind-mounted)
	files, err := ioutil.ReadDir(root)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, f := range files ***REMOVED***
		if err := os.RemoveAll(filepath.Join(root, f.Name())); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func removingManagerCausesLossOfQuorum(reachable, unreachable int) bool ***REMOVED***
	return reachable-2 <= unreachable
***REMOVED***

func isLastManager(reachable, unreachable int) bool ***REMOVED***
	return reachable == 1 && unreachable == 0
***REMOVED***
