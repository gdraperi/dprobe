// Package reference provides a general type to represent any way of referencing images within the registry.
// Its main purpose is to abstract tags and digests (content-addressable hash).
//
// Grammar
//
// 	reference                       := name [ ":" tag ] [ "@" digest ]
//	name                            := [domain '/'] path-component ['/' path-component]*
//	domain                          := domain-component ['.' domain-component]* [':' port-number]
//	domain-component                := /([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])/
//	port-number                     := /[0-9]+/
//	path-component                  := alpha-numeric [separator alpha-numeric]*
// 	alpha-numeric                   := /[a-z0-9]+/
//	separator                       := /[_.]|__|[-]*/
//
//	tag                             := /[\w][\w.-]***REMOVED***0,127***REMOVED***/
//
//	digest                          := digest-algorithm ":" digest-hex
//	digest-algorithm                := digest-algorithm-component [ digest-algorithm-separator digest-algorithm-component ]*
//	digest-algorithm-separator      := /[+.-_]/
//	digest-algorithm-component      := /[A-Za-z][A-Za-z0-9]*/
//	digest-hex                      := /[0-9a-fA-F]***REMOVED***32,***REMOVED***/ ; At least 128 bit digest value
//
//	identifier                      := /[a-f0-9]***REMOVED***64***REMOVED***/
//	short-identifier                := /[a-f0-9]***REMOVED***6,64***REMOVED***/
package reference

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opencontainers/go-digest"
)

const (
	// NameTotalLengthMax is the maximum total number of characters in a repository name.
	NameTotalLengthMax = 255
)

var (
	// ErrReferenceInvalidFormat represents an error while trying to parse a string as a reference.
	ErrReferenceInvalidFormat = errors.New("invalid reference format")

	// ErrTagInvalidFormat represents an error while trying to parse a string as a tag.
	ErrTagInvalidFormat = errors.New("invalid tag format")

	// ErrDigestInvalidFormat represents an error while trying to parse a string as a tag.
	ErrDigestInvalidFormat = errors.New("invalid digest format")

	// ErrNameContainsUppercase is returned for invalid repository names that contain uppercase characters.
	ErrNameContainsUppercase = errors.New("repository name must be lowercase")

	// ErrNameEmpty is returned for empty, invalid repository names.
	ErrNameEmpty = errors.New("repository name must have at least one component")

	// ErrNameTooLong is returned when a repository name is longer than NameTotalLengthMax.
	ErrNameTooLong = fmt.Errorf("repository name must not be more than %v characters", NameTotalLengthMax)

	// ErrNameNotCanonical is returned when a name is not canonical.
	ErrNameNotCanonical = errors.New("repository name must be canonical")
)

// Reference is an opaque object reference identifier that may include
// modifiers such as a hostname, name, tag, and digest.
type Reference interface ***REMOVED***
	// String returns the full reference
	String() string
***REMOVED***

// Field provides a wrapper type for resolving correct reference types when
// working with encoding.
type Field struct ***REMOVED***
	reference Reference
***REMOVED***

// AsField wraps a reference in a Field for encoding.
func AsField(reference Reference) Field ***REMOVED***
	return Field***REMOVED***reference***REMOVED***
***REMOVED***

// Reference unwraps the reference type from the field to
// return the Reference object. This object should be
// of the appropriate type to further check for different
// reference types.
func (f Field) Reference() Reference ***REMOVED***
	return f.reference
***REMOVED***

// MarshalText serializes the field to byte text which
// is the string of the reference.
func (f Field) MarshalText() (p []byte, err error) ***REMOVED***
	return []byte(f.reference.String()), nil
***REMOVED***

// UnmarshalText parses text bytes by invoking the
// reference parser to ensure the appropriately
// typed reference object is wrapped by field.
func (f *Field) UnmarshalText(p []byte) error ***REMOVED***
	r, err := Parse(string(p))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	f.reference = r
	return nil
***REMOVED***

// Named is an object with a full name
type Named interface ***REMOVED***
	Reference
	Name() string
***REMOVED***

// Tagged is an object which has a tag
type Tagged interface ***REMOVED***
	Reference
	Tag() string
***REMOVED***

// NamedTagged is an object including a name and tag.
type NamedTagged interface ***REMOVED***
	Named
	Tag() string
***REMOVED***

// Digested is an object which has a digest
// in which it can be referenced by
type Digested interface ***REMOVED***
	Reference
	Digest() digest.Digest
***REMOVED***

// Canonical reference is an object with a fully unique
// name including a name with domain and digest
type Canonical interface ***REMOVED***
	Named
	Digest() digest.Digest
***REMOVED***

// namedRepository is a reference to a repository with a name.
// A namedRepository has both domain and path components.
type namedRepository interface ***REMOVED***
	Named
	Domain() string
	Path() string
***REMOVED***

// Domain returns the domain part of the Named reference
func Domain(named Named) string ***REMOVED***
	if r, ok := named.(namedRepository); ok ***REMOVED***
		return r.Domain()
	***REMOVED***
	domain, _ := splitDomain(named.Name())
	return domain
***REMOVED***

// Path returns the name without the domain part of the Named reference
func Path(named Named) (name string) ***REMOVED***
	if r, ok := named.(namedRepository); ok ***REMOVED***
		return r.Path()
	***REMOVED***
	_, path := splitDomain(named.Name())
	return path
***REMOVED***

func splitDomain(name string) (string, string) ***REMOVED***
	match := anchoredNameRegexp.FindStringSubmatch(name)
	if len(match) != 3 ***REMOVED***
		return "", name
	***REMOVED***
	return match[1], match[2]
***REMOVED***

// SplitHostname splits a named reference into a
// hostname and name string. If no valid hostname is
// found, the hostname is empty and the full value
// is returned as name
// DEPRECATED: Use Domain or Path
func SplitHostname(named Named) (string, string) ***REMOVED***
	if r, ok := named.(namedRepository); ok ***REMOVED***
		return r.Domain(), r.Path()
	***REMOVED***
	return splitDomain(named.Name())
***REMOVED***

// Parse parses s and returns a syntactically valid Reference.
// If an error was encountered it is returned, along with a nil Reference.
// NOTE: Parse will not handle short digests.
func Parse(s string) (Reference, error) ***REMOVED***
	matches := ReferenceRegexp.FindStringSubmatch(s)
	if matches == nil ***REMOVED***
		if s == "" ***REMOVED***
			return nil, ErrNameEmpty
		***REMOVED***
		if ReferenceRegexp.FindStringSubmatch(strings.ToLower(s)) != nil ***REMOVED***
			return nil, ErrNameContainsUppercase
		***REMOVED***
		return nil, ErrReferenceInvalidFormat
	***REMOVED***

	if len(matches[1]) > NameTotalLengthMax ***REMOVED***
		return nil, ErrNameTooLong
	***REMOVED***

	var repo repository

	nameMatch := anchoredNameRegexp.FindStringSubmatch(matches[1])
	if nameMatch != nil && len(nameMatch) == 3 ***REMOVED***
		repo.domain = nameMatch[1]
		repo.path = nameMatch[2]
	***REMOVED*** else ***REMOVED***
		repo.domain = ""
		repo.path = matches[1]
	***REMOVED***

	ref := reference***REMOVED***
		namedRepository: repo,
		tag:             matches[2],
	***REMOVED***
	if matches[3] != "" ***REMOVED***
		var err error
		ref.digest, err = digest.Parse(matches[3])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	r := getBestReferenceType(ref)
	if r == nil ***REMOVED***
		return nil, ErrNameEmpty
	***REMOVED***

	return r, nil
***REMOVED***

// ParseNamed parses s and returns a syntactically valid reference implementing
// the Named interface. The reference must have a name and be in the canonical
// form, otherwise an error is returned.
// If an error was encountered it is returned, along with a nil Reference.
// NOTE: ParseNamed will not handle short digests.
func ParseNamed(s string) (Named, error) ***REMOVED***
	named, err := ParseNormalizedNamed(s)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if named.String() != s ***REMOVED***
		return nil, ErrNameNotCanonical
	***REMOVED***
	return named, nil
***REMOVED***

// WithName returns a named object representing the given string. If the input
// is invalid ErrReferenceInvalidFormat will be returned.
func WithName(name string) (Named, error) ***REMOVED***
	if len(name) > NameTotalLengthMax ***REMOVED***
		return nil, ErrNameTooLong
	***REMOVED***

	match := anchoredNameRegexp.FindStringSubmatch(name)
	if match == nil || len(match) != 3 ***REMOVED***
		return nil, ErrReferenceInvalidFormat
	***REMOVED***
	return repository***REMOVED***
		domain: match[1],
		path:   match[2],
	***REMOVED***, nil
***REMOVED***

// WithTag combines the name from "name" and the tag from "tag" to form a
// reference incorporating both the name and the tag.
func WithTag(name Named, tag string) (NamedTagged, error) ***REMOVED***
	if !anchoredTagRegexp.MatchString(tag) ***REMOVED***
		return nil, ErrTagInvalidFormat
	***REMOVED***
	var repo repository
	if r, ok := name.(namedRepository); ok ***REMOVED***
		repo.domain = r.Domain()
		repo.path = r.Path()
	***REMOVED*** else ***REMOVED***
		repo.path = name.Name()
	***REMOVED***
	if canonical, ok := name.(Canonical); ok ***REMOVED***
		return reference***REMOVED***
			namedRepository: repo,
			tag:             tag,
			digest:          canonical.Digest(),
		***REMOVED***, nil
	***REMOVED***
	return taggedReference***REMOVED***
		namedRepository: repo,
		tag:             tag,
	***REMOVED***, nil
***REMOVED***

// WithDigest combines the name from "name" and the digest from "digest" to form
// a reference incorporating both the name and the digest.
func WithDigest(name Named, digest digest.Digest) (Canonical, error) ***REMOVED***
	if !anchoredDigestRegexp.MatchString(digest.String()) ***REMOVED***
		return nil, ErrDigestInvalidFormat
	***REMOVED***
	var repo repository
	if r, ok := name.(namedRepository); ok ***REMOVED***
		repo.domain = r.Domain()
		repo.path = r.Path()
	***REMOVED*** else ***REMOVED***
		repo.path = name.Name()
	***REMOVED***
	if tagged, ok := name.(Tagged); ok ***REMOVED***
		return reference***REMOVED***
			namedRepository: repo,
			tag:             tagged.Tag(),
			digest:          digest,
		***REMOVED***, nil
	***REMOVED***
	return canonicalReference***REMOVED***
		namedRepository: repo,
		digest:          digest,
	***REMOVED***, nil
***REMOVED***

// TrimNamed removes any tag or digest from the named reference.
func TrimNamed(ref Named) Named ***REMOVED***
	domain, path := SplitHostname(ref)
	return repository***REMOVED***
		domain: domain,
		path:   path,
	***REMOVED***
***REMOVED***

func getBestReferenceType(ref reference) Reference ***REMOVED***
	if ref.Name() == "" ***REMOVED***
		// Allow digest only references
		if ref.digest != "" ***REMOVED***
			return digestReference(ref.digest)
		***REMOVED***
		return nil
	***REMOVED***
	if ref.tag == "" ***REMOVED***
		if ref.digest != "" ***REMOVED***
			return canonicalReference***REMOVED***
				namedRepository: ref.namedRepository,
				digest:          ref.digest,
			***REMOVED***
		***REMOVED***
		return ref.namedRepository
	***REMOVED***
	if ref.digest == "" ***REMOVED***
		return taggedReference***REMOVED***
			namedRepository: ref.namedRepository,
			tag:             ref.tag,
		***REMOVED***
	***REMOVED***

	return ref
***REMOVED***

type reference struct ***REMOVED***
	namedRepository
	tag    string
	digest digest.Digest
***REMOVED***

func (r reference) String() string ***REMOVED***
	return r.Name() + ":" + r.tag + "@" + r.digest.String()
***REMOVED***

func (r reference) Tag() string ***REMOVED***
	return r.tag
***REMOVED***

func (r reference) Digest() digest.Digest ***REMOVED***
	return r.digest
***REMOVED***

type repository struct ***REMOVED***
	domain string
	path   string
***REMOVED***

func (r repository) String() string ***REMOVED***
	return r.Name()
***REMOVED***

func (r repository) Name() string ***REMOVED***
	if r.domain == "" ***REMOVED***
		return r.path
	***REMOVED***
	return r.domain + "/" + r.path
***REMOVED***

func (r repository) Domain() string ***REMOVED***
	return r.domain
***REMOVED***

func (r repository) Path() string ***REMOVED***
	return r.path
***REMOVED***

type digestReference digest.Digest

func (d digestReference) String() string ***REMOVED***
	return digest.Digest(d).String()
***REMOVED***

func (d digestReference) Digest() digest.Digest ***REMOVED***
	return digest.Digest(d)
***REMOVED***

type taggedReference struct ***REMOVED***
	namedRepository
	tag string
***REMOVED***

func (t taggedReference) String() string ***REMOVED***
	return t.Name() + ":" + t.tag
***REMOVED***

func (t taggedReference) Tag() string ***REMOVED***
	return t.tag
***REMOVED***

type canonicalReference struct ***REMOVED***
	namedRepository
	digest digest.Digest
***REMOVED***

func (c canonicalReference) String() string ***REMOVED***
	return c.Name() + "@" + c.digest.String()
***REMOVED***

func (c canonicalReference) Digest() digest.Digest ***REMOVED***
	return c.digest
***REMOVED***
