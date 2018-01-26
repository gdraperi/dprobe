package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/daemon/names"
)

var (
	validCheckpointNameChars   = names.RestrictedNameChars
	validCheckpointNamePattern = names.RestrictedNamePattern
)

// getCheckpointDir verifies checkpoint directory for create,remove, list options and checks if checkpoint already exists
func getCheckpointDir(checkDir, checkpointID, ctrName, ctrID, ctrCheckpointDir string, create bool) (string, error) ***REMOVED***
	var checkpointDir string
	var err2 error
	if checkDir != "" ***REMOVED***
		checkpointDir = checkDir
	***REMOVED*** else ***REMOVED***
		checkpointDir = ctrCheckpointDir
	***REMOVED***
	checkpointAbsDir := filepath.Join(checkpointDir, checkpointID)
	stat, err := os.Stat(checkpointAbsDir)
	if create ***REMOVED***
		switch ***REMOVED***
		case err == nil && stat.IsDir():
			err2 = fmt.Errorf("checkpoint with name %s already exists for container %s", checkpointID, ctrName)
		case err != nil && os.IsNotExist(err):
			err2 = os.MkdirAll(checkpointAbsDir, 0700)
		case err != nil:
			err2 = err
		case err == nil:
			err2 = fmt.Errorf("%s exists and is not a directory", checkpointAbsDir)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		switch ***REMOVED***
		case err != nil:
			err2 = fmt.Errorf("checkpoint %s does not exists for container %s", checkpointID, ctrName)
		case err == nil && stat.IsDir():
			err2 = nil
		case err == nil:
			err2 = fmt.Errorf("%s exists and is not a directory", checkpointAbsDir)
		***REMOVED***
	***REMOVED***
	return checkpointAbsDir, err2
***REMOVED***

// CheckpointCreate checkpoints the process running in a container with CRIU
func (daemon *Daemon) CheckpointCreate(name string, config types.CheckpointCreateOptions) error ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !container.IsRunning() ***REMOVED***
		return fmt.Errorf("Container %s not running", name)
	***REMOVED***

	if container.Config.Tty ***REMOVED***
		return fmt.Errorf("checkpoint not support on containers with tty")
	***REMOVED***

	if !validCheckpointNamePattern.MatchString(config.CheckpointID) ***REMOVED***
		return fmt.Errorf("Invalid checkpoint ID (%s), only %s are allowed", config.CheckpointID, validCheckpointNameChars)
	***REMOVED***

	checkpointDir, err := getCheckpointDir(config.CheckpointDir, config.CheckpointID, name, container.ID, container.CheckpointDir(), true)
	if err != nil ***REMOVED***
		return fmt.Errorf("cannot checkpoint container %s: %s", name, err)
	***REMOVED***

	err = daemon.containerd.CreateCheckpoint(context.Background(), container.ID, checkpointDir, config.Exit)
	if err != nil ***REMOVED***
		os.RemoveAll(checkpointDir)
		return fmt.Errorf("Cannot checkpoint container %s: %s", name, err)
	***REMOVED***

	daemon.LogContainerEvent(container, "checkpoint")

	return nil
***REMOVED***

// CheckpointDelete deletes the specified checkpoint
func (daemon *Daemon) CheckpointDelete(name string, config types.CheckpointDeleteOptions) error ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	checkpointDir, err := getCheckpointDir(config.CheckpointDir, config.CheckpointID, name, container.ID, container.CheckpointDir(), false)
	if err == nil ***REMOVED***
		return os.RemoveAll(filepath.Join(checkpointDir, config.CheckpointID))
	***REMOVED***
	return err
***REMOVED***

// CheckpointList lists all checkpoints of the specified container
func (daemon *Daemon) CheckpointList(name string, config types.CheckpointListOptions) ([]types.Checkpoint, error) ***REMOVED***
	var out []types.Checkpoint

	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	checkpointDir, err := getCheckpointDir(config.CheckpointDir, "", name, container.ID, container.CheckpointDir(), false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := os.MkdirAll(checkpointDir, 0755); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	dirs, err := ioutil.ReadDir(checkpointDir)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, d := range dirs ***REMOVED***
		if !d.IsDir() ***REMOVED***
			continue
		***REMOVED***
		path := filepath.Join(checkpointDir, d.Name(), "config.json")
		data, err := ioutil.ReadFile(path)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var cpt types.Checkpoint
		if err := json.Unmarshal(data, &cpt); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		out = append(out, cpt)
	***REMOVED***

	return out, nil
***REMOVED***
