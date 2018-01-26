package distribution

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opencontainers/go-digest"
)

// ErrAccessDenied is returned when an access to a requested resource is
// denied.
var ErrAccessDenied = errors.New("access denied")

// ErrManifestNotModified is returned when a conditional manifest GetByTag
// returns nil due to the client indicating it has the latest version
var ErrManifestNotModified = errors.New("manifest not modified")

// ErrUnsupported is returned when an unimplemented or unsupported action is
// performed
var ErrUnsupported = errors.New("operation unsupported")

// ErrTagUnknown is returned if the given tag is not known by the tag service
type ErrTagUnknown struct ***REMOVED***
	Tag string
***REMOVED***

func (err ErrTagUnknown) Error() string ***REMOVED***
	return fmt.Sprintf("unknown tag=%s", err.Tag)
***REMOVED***

// ErrRepositoryUnknown is returned if the named repository is not known by
// the registry.
type ErrRepositoryUnknown struct ***REMOVED***
	Name string
***REMOVED***

func (err ErrRepositoryUnknown) Error() string ***REMOVED***
	return fmt.Sprintf("unknown repository name=%s", err.Name)
***REMOVED***

// ErrRepositoryNameInvalid should be used to denote an invalid repository
// name. Reason may set, indicating the cause of invalidity.
type ErrRepositoryNameInvalid struct ***REMOVED***
	Name   string
	Reason error
***REMOVED***

func (err ErrRepositoryNameInvalid) Error() string ***REMOVED***
	return fmt.Sprintf("repository name %q invalid: %v", err.Name, err.Reason)
***REMOVED***

// ErrManifestUnknown is returned if the manifest is not known by the
// registry.
type ErrManifestUnknown struct ***REMOVED***
	Name string
	Tag  string
***REMOVED***

func (err ErrManifestUnknown) Error() string ***REMOVED***
	return fmt.Sprintf("unknown manifest name=%s tag=%s", err.Name, err.Tag)
***REMOVED***

// ErrManifestUnknownRevision is returned when a manifest cannot be found by
// revision within a repository.
type ErrManifestUnknownRevision struct ***REMOVED***
	Name     string
	Revision digest.Digest
***REMOVED***

func (err ErrManifestUnknownRevision) Error() string ***REMOVED***
	return fmt.Sprintf("unknown manifest name=%s revision=%s", err.Name, err.Revision)
***REMOVED***

// ErrManifestUnverified is returned when the registry is unable to verify
// the manifest.
type ErrManifestUnverified struct***REMOVED******REMOVED***

func (ErrManifestUnverified) Error() string ***REMOVED***
	return "unverified manifest"
***REMOVED***

// ErrManifestVerification provides a type to collect errors encountered
// during manifest verification. Currently, it accepts errors of all types,
// but it may be narrowed to those involving manifest verification.
type ErrManifestVerification []error

func (errs ErrManifestVerification) Error() string ***REMOVED***
	var parts []string
	for _, err := range errs ***REMOVED***
		parts = append(parts, err.Error())
	***REMOVED***

	return fmt.Sprintf("errors verifying manifest: %v", strings.Join(parts, ","))
***REMOVED***

// ErrManifestBlobUnknown returned when a referenced blob cannot be found.
type ErrManifestBlobUnknown struct ***REMOVED***
	Digest digest.Digest
***REMOVED***

func (err ErrManifestBlobUnknown) Error() string ***REMOVED***
	return fmt.Sprintf("unknown blob %v on manifest", err.Digest)
***REMOVED***

// ErrManifestNameInvalid should be used to denote an invalid manifest
// name. Reason may set, indicating the cause of invalidity.
type ErrManifestNameInvalid struct ***REMOVED***
	Name   string
	Reason error
***REMOVED***

func (err ErrManifestNameInvalid) Error() string ***REMOVED***
	return fmt.Sprintf("manifest name %q invalid: %v", err.Name, err.Reason)
***REMOVED***
