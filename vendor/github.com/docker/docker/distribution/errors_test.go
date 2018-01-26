package distribution

import (
	"errors"
	"strings"
	"syscall"
	"testing"

	"github.com/docker/distribution/registry/api/errcode"
	"github.com/docker/distribution/registry/api/v2"
	"github.com/docker/distribution/registry/client"
)

var alwaysContinue = []error***REMOVED***
	&client.UnexpectedHTTPResponseError***REMOVED******REMOVED***,

	// Some errcode.Errors that don't disprove the existence of a V1 image
	errcode.Error***REMOVED***Code: errcode.ErrorCodeUnauthorized***REMOVED***,
	errcode.Error***REMOVED***Code: v2.ErrorCodeManifestUnknown***REMOVED***,
	errcode.Error***REMOVED***Code: v2.ErrorCodeNameUnknown***REMOVED***,

	errors.New("some totally unexpected error"),
***REMOVED***

var continueFromMirrorEndpoint = []error***REMOVED***
	ImageConfigPullError***REMOVED******REMOVED***,

	// Some other errcode.Error that doesn't indicate we should search for a V1 image.
	errcode.Error***REMOVED***Code: errcode.ErrorCodeTooManyRequests***REMOVED***,
***REMOVED***

var neverContinue = []error***REMOVED***
	errors.New(strings.ToLower(syscall.ESRCH.Error())), // No such process
***REMOVED***

func TestContinueOnError_NonMirrorEndpoint(t *testing.T) ***REMOVED***
	for _, err := range alwaysContinue ***REMOVED***
		if !continueOnError(err, false) ***REMOVED***
			t.Errorf("Should continue from non-mirror endpoint: %T: '%s'", err, err.Error())
		***REMOVED***
	***REMOVED***

	for _, err := range continueFromMirrorEndpoint ***REMOVED***
		if continueOnError(err, false) ***REMOVED***
			t.Errorf("Should only continue from mirror endpoint: %T: '%s'", err, err.Error())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestContinueOnError_MirrorEndpoint(t *testing.T) ***REMOVED***
	errs := []error***REMOVED******REMOVED***
	errs = append(errs, alwaysContinue...)
	errs = append(errs, continueFromMirrorEndpoint...)
	for _, err := range errs ***REMOVED***
		if !continueOnError(err, true) ***REMOVED***
			t.Errorf("Should continue from mirror endpoint: %T: '%s'", err, err.Error())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestContinueOnError_NeverContinue(t *testing.T) ***REMOVED***
	for _, isMirrorEndpoint := range []bool***REMOVED***true, false***REMOVED*** ***REMOVED***
		for _, err := range neverContinue ***REMOVED***
			if continueOnError(err, isMirrorEndpoint) ***REMOVED***
				t.Errorf("Should never continue: %T: '%s'", err, err.Error())
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestContinueOnError_UnnestsErrors(t *testing.T) ***REMOVED***
	// ContinueOnError should evaluate nested errcode.Errors.

	// Assumes that v2.ErrorCodeNameUnknown is a continueable error code.
	err := errcode.Errors***REMOVED***errcode.Error***REMOVED***Code: v2.ErrorCodeNameUnknown***REMOVED******REMOVED***
	if !continueOnError(err, false) ***REMOVED***
		t.Fatal("ContinueOnError should unnest, base return value on errcode.Errors")
	***REMOVED***

	// Assumes that errcode.ErrorCodeTooManyRequests is not a V1-fallback indication
	err = errcode.Errors***REMOVED***errcode.Error***REMOVED***Code: errcode.ErrorCodeTooManyRequests***REMOVED******REMOVED***
	if continueOnError(err, false) ***REMOVED***
		t.Fatal("ContinueOnError should unnest, base return value on errcode.Errors")
	***REMOVED***
***REMOVED***
