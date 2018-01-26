package fscache

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// NewNaiveCacheBackend is a basic backend implementation for fscache
func NewNaiveCacheBackend(root string) Backend ***REMOVED***
	return &naiveCacheBackend***REMOVED***root: root***REMOVED***
***REMOVED***

type naiveCacheBackend struct ***REMOVED***
	root string
***REMOVED***

func (tcb *naiveCacheBackend) Get(id string) (string, error) ***REMOVED***
	d := filepath.Join(tcb.root, id)
	if err := os.MkdirAll(d, 0700); err != nil ***REMOVED***
		return "", errors.Wrapf(err, "failed to create tmp dir for %s", d)
	***REMOVED***
	return d, nil
***REMOVED***
func (tcb *naiveCacheBackend) Remove(id string) error ***REMOVED***
	return errors.WithStack(os.RemoveAll(filepath.Join(tcb.root, id)))
***REMOVED***
