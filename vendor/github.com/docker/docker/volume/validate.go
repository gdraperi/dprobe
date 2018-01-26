package volume

import (
	"fmt"

	"github.com/docker/docker/api/types/mount"
	"github.com/pkg/errors"
)

var errBindNotExist = errors.New("bind source path does not exist")

type errMountConfig struct ***REMOVED***
	mount *mount.Mount
	err   error
***REMOVED***

func (e *errMountConfig) Error() string ***REMOVED***
	return fmt.Sprintf("invalid mount config for type %q: %v", e.mount.Type, e.err.Error())
***REMOVED***

func errExtraField(name string) error ***REMOVED***
	return errors.Errorf("field %s must not be specified", name)
***REMOVED***
func errMissingField(name string) error ***REMOVED***
	return errors.Errorf("field %s must not be empty", name)
***REMOVED***
