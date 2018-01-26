package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

// New returns a new control via the cgroup cgroups interface
func New(hierarchy Hierarchy, path Path, resources *specs.LinuxResources) (Cgroup, error) ***REMOVED***
	subsystems, err := hierarchy()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, s := range subsystems ***REMOVED***
		if err := initializeSubsystem(s, path, resources); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return &cgroup***REMOVED***
		path:       path,
		subsystems: subsystems,
	***REMOVED***, nil
***REMOVED***

// Load will load an existing cgroup and allow it to be controlled
func Load(hierarchy Hierarchy, path Path) (Cgroup, error) ***REMOVED***
	subsystems, err := hierarchy()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// check the the subsystems still exist
	for _, s := range pathers(subsystems) ***REMOVED***
		p, err := path(s.Name())
		if err != nil ***REMOVED***
			if os.IsNotExist(errors.Cause(err)) ***REMOVED***
				return nil, ErrCgroupDeleted
			***REMOVED***
			return nil, err
		***REMOVED***
		if _, err := os.Lstat(s.Path(p)); err != nil ***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				return nil, ErrCgroupDeleted
			***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return &cgroup***REMOVED***
		path:       path,
		subsystems: subsystems,
	***REMOVED***, nil
***REMOVED***

type cgroup struct ***REMOVED***
	path Path

	subsystems []Subsystem
	mu         sync.Mutex
	err        error
***REMOVED***

// New returns a new sub cgroup
func (c *cgroup) New(name string, resources *specs.LinuxResources) (Cgroup, error) ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil ***REMOVED***
		return nil, c.err
	***REMOVED***
	path := subPath(c.path, name)
	for _, s := range c.subsystems ***REMOVED***
		if err := initializeSubsystem(s, path, resources); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return &cgroup***REMOVED***
		path:       path,
		subsystems: c.subsystems,
	***REMOVED***, nil
***REMOVED***

// Subsystems returns all the subsystems that are currently being
// consumed by the group
func (c *cgroup) Subsystems() []Subsystem ***REMOVED***
	return c.subsystems
***REMOVED***

// Add moves the provided process into the new cgroup
func (c *cgroup) Add(process Process) error ***REMOVED***
	if process.Pid <= 0 ***REMOVED***
		return ErrInvalidPid
	***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil ***REMOVED***
		return c.err
	***REMOVED***
	return c.add(process)
***REMOVED***

func (c *cgroup) add(process Process) error ***REMOVED***
	for _, s := range pathers(c.subsystems) ***REMOVED***
		p, err := c.path(s.Name())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := ioutil.WriteFile(
			filepath.Join(s.Path(p), cgroupProcs),
			[]byte(strconv.Itoa(process.Pid)),
			defaultFilePerm,
		); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Delete will remove the control group from each of the subsystems registered
func (c *cgroup) Delete() error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil ***REMOVED***
		return c.err
	***REMOVED***
	var errors []string
	for _, s := range c.subsystems ***REMOVED***
		if d, ok := s.(deleter); ok ***REMOVED***
			sp, err := c.path(s.Name())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := d.Delete(sp); err != nil ***REMOVED***
				errors = append(errors, string(s.Name()))
			***REMOVED***
			continue
		***REMOVED***
		if p, ok := s.(pather); ok ***REMOVED***
			sp, err := c.path(s.Name())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			path := p.Path(sp)
			if err := remove(path); err != nil ***REMOVED***
				errors = append(errors, path)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(errors) > 0 ***REMOVED***
		return fmt.Errorf("cgroups: unable to remove paths %s", strings.Join(errors, ", "))
	***REMOVED***
	c.err = ErrCgroupDeleted
	return nil
***REMOVED***

// Stat returns the current metrics for the cgroup
func (c *cgroup) Stat(handlers ...ErrorHandler) (*Metrics, error) ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil ***REMOVED***
		return nil, c.err
	***REMOVED***
	if len(handlers) == 0 ***REMOVED***
		handlers = append(handlers, errPassthrough)
	***REMOVED***
	var (
		stats = &Metrics***REMOVED***
			CPU: &CPUStat***REMOVED***
				Throttling: &Throttle***REMOVED******REMOVED***,
				Usage:      &CPUUsage***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***
		wg   = &sync.WaitGroup***REMOVED******REMOVED***
		errs = make(chan error, len(c.subsystems))
	)
	for _, s := range c.subsystems ***REMOVED***
		if ss, ok := s.(stater); ok ***REMOVED***
			sp, err := c.path(s.Name())
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			wg.Add(1)
			go func() ***REMOVED***
				defer wg.Done()
				if err := ss.Stat(sp, stats); err != nil ***REMOVED***
					for _, eh := range handlers ***REMOVED***
						if herr := eh(err); herr != nil ***REMOVED***
							errs <- herr
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***()
		***REMOVED***
	***REMOVED***
	wg.Wait()
	close(errs)
	for err := range errs ***REMOVED***
		return nil, err
	***REMOVED***
	return stats, nil
***REMOVED***

// Update updates the cgroup with the new resource values provided
//
// Be prepared to handle EBUSY when trying to update a cgroup with
// live processes and other operations like Stats being performed at the
// same time
func (c *cgroup) Update(resources *specs.LinuxResources) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil ***REMOVED***
		return c.err
	***REMOVED***
	for _, s := range c.subsystems ***REMOVED***
		if u, ok := s.(updater); ok ***REMOVED***
			sp, err := c.path(s.Name())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := u.Update(sp, resources); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Processes returns the processes running inside the cgroup along
// with the subsystem used, pid, and path
func (c *cgroup) Processes(subsystem Name, recursive bool) ([]Process, error) ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil ***REMOVED***
		return nil, c.err
	***REMOVED***
	return c.processes(subsystem, recursive)
***REMOVED***

func (c *cgroup) processes(subsystem Name, recursive bool) ([]Process, error) ***REMOVED***
	s := c.getSubsystem(subsystem)
	sp, err := c.path(subsystem)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	path := s.(pather).Path(sp)
	var processes []Process
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if !recursive && info.IsDir() ***REMOVED***
			if p == path ***REMOVED***
				return nil
			***REMOVED***
			return filepath.SkipDir
		***REMOVED***
		dir, name := filepath.Split(p)
		if name != cgroupProcs ***REMOVED***
			return nil
		***REMOVED***
		procs, err := readPids(dir, subsystem)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		processes = append(processes, procs...)
		return nil
	***REMOVED***)
	return processes, err
***REMOVED***

// Freeze freezes the entire cgroup and all the processes inside it
func (c *cgroup) Freeze() error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil ***REMOVED***
		return c.err
	***REMOVED***
	s := c.getSubsystem(Freezer)
	if s == nil ***REMOVED***
		return ErrFreezerNotSupported
	***REMOVED***
	sp, err := c.path(Freezer)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.(*freezerController).Freeze(sp)
***REMOVED***

// Thaw thaws out the cgroup and all the processes inside it
func (c *cgroup) Thaw() error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil ***REMOVED***
		return c.err
	***REMOVED***
	s := c.getSubsystem(Freezer)
	if s == nil ***REMOVED***
		return ErrFreezerNotSupported
	***REMOVED***
	sp, err := c.path(Freezer)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.(*freezerController).Thaw(sp)
***REMOVED***

// OOMEventFD returns the memory cgroup's out of memory event fd that triggers
// when processes inside the cgroup receive an oom event. Returns
// ErrMemoryNotSupported if memory cgroups is not supported.
func (c *cgroup) OOMEventFD() (uintptr, error) ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil ***REMOVED***
		return 0, c.err
	***REMOVED***
	s := c.getSubsystem(Memory)
	if s == nil ***REMOVED***
		return 0, ErrMemoryNotSupported
	***REMOVED***
	sp, err := c.path(Memory)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return s.(*memoryController).OOMEventFD(sp)
***REMOVED***

// State returns the state of the cgroup and its processes
func (c *cgroup) State() State ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checkExists()
	if c.err != nil && c.err == ErrCgroupDeleted ***REMOVED***
		return Deleted
	***REMOVED***
	s := c.getSubsystem(Freezer)
	if s == nil ***REMOVED***
		return Thawed
	***REMOVED***
	sp, err := c.path(Freezer)
	if err != nil ***REMOVED***
		return Unknown
	***REMOVED***
	state, err := s.(*freezerController).state(sp)
	if err != nil ***REMOVED***
		return Unknown
	***REMOVED***
	return state
***REMOVED***

// MoveTo does a recursive move subsystem by subsystem of all the processes
// inside the group
func (c *cgroup) MoveTo(destination Cgroup) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil ***REMOVED***
		return c.err
	***REMOVED***
	for _, s := range c.subsystems ***REMOVED***
		processes, err := c.processes(s.Name(), true)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, p := range processes ***REMOVED***
			if err := destination.Add(p); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c *cgroup) getSubsystem(n Name) Subsystem ***REMOVED***
	for _, s := range c.subsystems ***REMOVED***
		if s.Name() == n ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c *cgroup) checkExists() ***REMOVED***
	for _, s := range pathers(c.subsystems) ***REMOVED***
		p, err := c.path(s.Name())
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if _, err := os.Lstat(s.Path(p)); err != nil ***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				c.err = ErrCgroupDeleted
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
