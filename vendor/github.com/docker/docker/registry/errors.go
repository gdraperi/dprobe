package registry

import (
	"net/url"

	"github.com/docker/distribution/registry/api/errcode"
	"github.com/docker/docker/errdefs"
)

type notFoundError string

func (e notFoundError) Error() string ***REMOVED***
	return string(e)
***REMOVED***

func (notFoundError) NotFound() ***REMOVED******REMOVED***

func translateV2AuthError(err error) error ***REMOVED***
	switch e := err.(type) ***REMOVED***
	case *url.Error:
		switch e2 := e.Err.(type) ***REMOVED***
		case errcode.Error:
			switch e2.Code ***REMOVED***
			case errcode.ErrorCodeUnauthorized:
				return errdefs.Unauthorized(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return err
***REMOVED***
