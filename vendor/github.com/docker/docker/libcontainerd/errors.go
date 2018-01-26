package libcontainerd

import (
	"errors"

	"github.com/docker/docker/errdefs"
)

func newNotFoundError(err string) error ***REMOVED*** return errdefs.NotFound(errors.New(err)) ***REMOVED***

func newInvalidParameterError(err string) error ***REMOVED*** return errdefs.InvalidParameter(errors.New(err)) ***REMOVED***

func newConflictError(err string) error ***REMOVED*** return errdefs.Conflict(errors.New(err)) ***REMOVED***
