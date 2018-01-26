package errdefs

import (
	"errors"
	"testing"
)

var errTest = errors.New("this is a test")

type causal interface ***REMOVED***
	Cause() error
***REMOVED***

func TestNotFound(t *testing.T) ***REMOVED***
	e := NotFound(errTest)
	if !IsNotFound(e) ***REMOVED***
		t.Fatalf("expected not found error, got: %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***

func TestConflict(t *testing.T) ***REMOVED***
	e := Conflict(errTest)
	if !IsConflict(e) ***REMOVED***
		t.Fatalf("expected conflcit error, got: %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***

func TestForbidden(t *testing.T) ***REMOVED***
	e := Forbidden(errTest)
	if !IsForbidden(e) ***REMOVED***
		t.Fatalf("expected forbidden error, got: %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***

func TestInvalidParameter(t *testing.T) ***REMOVED***
	e := InvalidParameter(errTest)
	if !IsInvalidParameter(e) ***REMOVED***
		t.Fatalf("expected invalid argument error, got %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***

func TestNotImplemented(t *testing.T) ***REMOVED***
	e := NotImplemented(errTest)
	if !IsNotImplemented(e) ***REMOVED***
		t.Fatalf("expected not implemented error, got %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***

func TestNotModified(t *testing.T) ***REMOVED***
	e := NotModified(errTest)
	if !IsNotModified(e) ***REMOVED***
		t.Fatalf("expected not modified error, got %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***

func TestAlreadyExists(t *testing.T) ***REMOVED***
	e := AlreadyExists(errTest)
	if !IsAlreadyExists(e) ***REMOVED***
		t.Fatalf("expected already exists error, got %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***

func TestUnauthorized(t *testing.T) ***REMOVED***
	e := Unauthorized(errTest)
	if !IsUnauthorized(e) ***REMOVED***
		t.Fatalf("expected unauthorized error, got %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***

func TestUnknown(t *testing.T) ***REMOVED***
	e := Unknown(errTest)
	if !IsUnknown(e) ***REMOVED***
		t.Fatalf("expected unknown error, got %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***

func TestCancelled(t *testing.T) ***REMOVED***
	e := Cancelled(errTest)
	if !IsCancelled(e) ***REMOVED***
		t.Fatalf("expected canclled error, got %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***

func TestDeadline(t *testing.T) ***REMOVED***
	e := Deadline(errTest)
	if !IsDeadline(e) ***REMOVED***
		t.Fatalf("expected deadline error, got %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***

func TestIsDataLoss(t *testing.T) ***REMOVED***
	e := DataLoss(errTest)
	if !IsDataLoss(e) ***REMOVED***
		t.Fatalf("expected data loss error, got %T", e)
	***REMOVED***
	if cause := e.(causal).Cause(); cause != errTest ***REMOVED***
		t.Fatalf("causual should be errTest, got: %v", cause)
	***REMOVED***
***REMOVED***
