package cgroups

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func NewCputset(root string) *cpusetController ***REMOVED***
	return &cpusetController***REMOVED***
		root: filepath.Join(root, string(Cpuset)),
	***REMOVED***
***REMOVED***

type cpusetController struct ***REMOVED***
	root string
***REMOVED***

func (c *cpusetController) Name() Name ***REMOVED***
	return Cpuset
***REMOVED***

func (c *cpusetController) Path(path string) string ***REMOVED***
	return filepath.Join(c.root, path)
***REMOVED***

func (c *cpusetController) Create(path string, resources *specs.LinuxResources) error ***REMOVED***
	if err := c.ensureParent(c.Path(path), c.root); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := os.MkdirAll(c.Path(path), defaultDirPerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := c.copyIfNeeded(c.Path(path), filepath.Dir(c.Path(path))); err != nil ***REMOVED***
		return err
	***REMOVED***
	if resources.CPU != nil ***REMOVED***
		for _, t := range []struct ***REMOVED***
			name  string
			value *string
		***REMOVED******REMOVED***
			***REMOVED***
				name:  "cpus",
				value: &resources.CPU.Cpus,
			***REMOVED***,
			***REMOVED***
				name:  "mems",
				value: &resources.CPU.Mems,
			***REMOVED***,
		***REMOVED*** ***REMOVED***
			if t.value != nil ***REMOVED***
				if err := ioutil.WriteFile(
					filepath.Join(c.Path(path), fmt.Sprintf("cpuset.%s", t.name)),
					[]byte(*t.value),
					defaultFilePerm,
				); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c *cpusetController) getValues(path string) (cpus []byte, mems []byte, err error) ***REMOVED***
	if cpus, err = ioutil.ReadFile(filepath.Join(path, "cpuset.cpus")); err != nil && !os.IsNotExist(err) ***REMOVED***
		return
	***REMOVED***
	if mems, err = ioutil.ReadFile(filepath.Join(path, "cpuset.mems")); err != nil && !os.IsNotExist(err) ***REMOVED***
		return
	***REMOVED***
	return cpus, mems, nil
***REMOVED***

// ensureParent makes sure that the parent directory of current is created
// and populated with the proper cpus and mems files copied from
// it's parent.
func (c *cpusetController) ensureParent(current, root string) error ***REMOVED***
	parent := filepath.Dir(current)
	if _, err := filepath.Rel(root, parent); err != nil ***REMOVED***
		return nil
	***REMOVED***
	// Avoid infinite recursion.
	if parent == current ***REMOVED***
		return fmt.Errorf("cpuset: cgroup parent path outside cgroup root")
	***REMOVED***
	if cleanPath(parent) != root ***REMOVED***
		if err := c.ensureParent(parent, root); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if err := os.MkdirAll(current, defaultDirPerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	return c.copyIfNeeded(current, parent)
***REMOVED***

// copyIfNeeded copies the cpuset.cpus and cpuset.mems from the parent
// directory to the current directory if the file's contents are 0
func (c *cpusetController) copyIfNeeded(current, parent string) error ***REMOVED***
	var (
		err                      error
		currentCpus, currentMems []byte
		parentCpus, parentMems   []byte
	)
	if currentCpus, currentMems, err = c.getValues(current); err != nil ***REMOVED***
		return err
	***REMOVED***
	if parentCpus, parentMems, err = c.getValues(parent); err != nil ***REMOVED***
		return err
	***REMOVED***
	if isEmpty(currentCpus) ***REMOVED***
		if err := ioutil.WriteFile(
			filepath.Join(current, "cpuset.cpus"),
			parentCpus,
			defaultFilePerm,
		); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if isEmpty(currentMems) ***REMOVED***
		if err := ioutil.WriteFile(
			filepath.Join(current, "cpuset.mems"),
			parentMems,
			defaultFilePerm,
		); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func isEmpty(b []byte) bool ***REMOVED***
	return len(bytes.Trim(b, "\n")) == 0
***REMOVED***
