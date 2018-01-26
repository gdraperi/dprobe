package cgroups

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

func NewFreezer(root string) *freezerController ***REMOVED***
	return &freezerController***REMOVED***
		root: filepath.Join(root, string(Freezer)),
	***REMOVED***
***REMOVED***

type freezerController struct ***REMOVED***
	root string
***REMOVED***

func (f *freezerController) Name() Name ***REMOVED***
	return Freezer
***REMOVED***

func (f *freezerController) Path(path string) string ***REMOVED***
	return filepath.Join(f.root, path)
***REMOVED***

func (f *freezerController) Freeze(path string) error ***REMOVED***
	return f.waitState(path, Frozen)
***REMOVED***

func (f *freezerController) Thaw(path string) error ***REMOVED***
	return f.waitState(path, Thawed)
***REMOVED***

func (f *freezerController) changeState(path string, state State) error ***REMOVED***
	return ioutil.WriteFile(
		filepath.Join(f.root, path, "freezer.state"),
		[]byte(strings.ToUpper(string(state))),
		defaultFilePerm,
	)
***REMOVED***

func (f *freezerController) state(path string) (State, error) ***REMOVED***
	current, err := ioutil.ReadFile(filepath.Join(f.root, path, "freezer.state"))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return State(strings.ToLower(strings.TrimSpace(string(current)))), nil
***REMOVED***

func (f *freezerController) waitState(path string, state State) error ***REMOVED***
	for ***REMOVED***
		if err := f.changeState(path, state); err != nil ***REMOVED***
			return err
		***REMOVED***
		current, err := f.state(path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if current == state ***REMOVED***
			return nil
		***REMOVED***
		time.Sleep(1 * time.Millisecond)
	***REMOVED***
***REMOVED***
