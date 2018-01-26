package labels

import (
	"github.com/containerd/containerd/errdefs"
	"github.com/pkg/errors"
)

const (
	maxSize = 4096
)

// Validate a label's key and value are under 4096 bytes
func Validate(k, v string) error ***REMOVED***
	if (len(k) + len(v)) > maxSize ***REMOVED***
		if len(k) > 10 ***REMOVED***
			k = k[:10]
		***REMOVED***
		return errors.Wrapf(errdefs.ErrInvalidArgument, "label key and value greater than maximum size (%d bytes), key: %s", maxSize, k)
	***REMOVED***
	return nil
***REMOVED***
