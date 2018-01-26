package reference

import (
	"errors"
	"fmt"
	"strings"

	"github.com/docker/distribution/digestset"
	"github.com/opencontainers/go-digest"
)

var (
	legacyDefaultDomain = "index.docker.io"
	defaultDomain       = "docker.io"
	officialRepoName    = "library"
	defaultTag          = "latest"
)

// normalizedNamed represents a name which has been
// normalized and has a familiar form. A familiar name
// is what is used in Docker UI. An example normalized
// name is "docker.io/library/ubuntu" and corresponding
// familiar name of "ubuntu".
type normalizedNamed interface ***REMOVED***
	Named
	Familiar() Named
***REMOVED***

// ParseNormalizedNamed parses a string into a named reference
// transforming a familiar name from Docker UI to a fully
// qualified reference. If the value may be an identifier
// use ParseAnyReference.
func ParseNormalizedNamed(s string) (Named, error) ***REMOVED***
	if ok := anchoredIdentifierRegexp.MatchString(s); ok ***REMOVED***
		return nil, fmt.Errorf("invalid repository name (%s), cannot specify 64-byte hexadecimal strings", s)
	***REMOVED***
	domain, remainder := splitDockerDomain(s)
	var remoteName string
	if tagSep := strings.IndexRune(remainder, ':'); tagSep > -1 ***REMOVED***
		remoteName = remainder[:tagSep]
	***REMOVED*** else ***REMOVED***
		remoteName = remainder
	***REMOVED***
	if strings.ToLower(remoteName) != remoteName ***REMOVED***
		return nil, errors.New("invalid reference format: repository name must be lowercase")
	***REMOVED***

	ref, err := Parse(domain + "/" + remainder)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	named, isNamed := ref.(Named)
	if !isNamed ***REMOVED***
		return nil, fmt.Errorf("reference %s has no name", ref.String())
	***REMOVED***
	return named, nil
***REMOVED***

// splitDockerDomain splits a repository name to domain and remotename string.
// If no valid domain is found, the default domain is used. Repository name
// needs to be already validated before.
func splitDockerDomain(name string) (domain, remainder string) ***REMOVED***
	i := strings.IndexRune(name, '/')
	if i == -1 || (!strings.ContainsAny(name[:i], ".:") && name[:i] != "localhost") ***REMOVED***
		domain, remainder = defaultDomain, name
	***REMOVED*** else ***REMOVED***
		domain, remainder = name[:i], name[i+1:]
	***REMOVED***
	if domain == legacyDefaultDomain ***REMOVED***
		domain = defaultDomain
	***REMOVED***
	if domain == defaultDomain && !strings.ContainsRune(remainder, '/') ***REMOVED***
		remainder = officialRepoName + "/" + remainder
	***REMOVED***
	return
***REMOVED***

// familiarizeName returns a shortened version of the name familiar
// to to the Docker UI. Familiar names have the default domain
// "docker.io" and "library/" repository prefix removed.
// For example, "docker.io/library/redis" will have the familiar
// name "redis" and "docker.io/dmcgowan/myapp" will be "dmcgowan/myapp".
// Returns a familiarized named only reference.
func familiarizeName(named namedRepository) repository ***REMOVED***
	repo := repository***REMOVED***
		domain: named.Domain(),
		path:   named.Path(),
	***REMOVED***

	if repo.domain == defaultDomain ***REMOVED***
		repo.domain = ""
		// Handle official repositories which have the pattern "library/<official repo name>"
		if split := strings.Split(repo.path, "/"); len(split) == 2 && split[0] == officialRepoName ***REMOVED***
			repo.path = split[1]
		***REMOVED***
	***REMOVED***
	return repo
***REMOVED***

func (r reference) Familiar() Named ***REMOVED***
	return reference***REMOVED***
		namedRepository: familiarizeName(r.namedRepository),
		tag:             r.tag,
		digest:          r.digest,
	***REMOVED***
***REMOVED***

func (r repository) Familiar() Named ***REMOVED***
	return familiarizeName(r)
***REMOVED***

func (t taggedReference) Familiar() Named ***REMOVED***
	return taggedReference***REMOVED***
		namedRepository: familiarizeName(t.namedRepository),
		tag:             t.tag,
	***REMOVED***
***REMOVED***

func (c canonicalReference) Familiar() Named ***REMOVED***
	return canonicalReference***REMOVED***
		namedRepository: familiarizeName(c.namedRepository),
		digest:          c.digest,
	***REMOVED***
***REMOVED***

// TagNameOnly adds the default tag "latest" to a reference if it only has
// a repo name.
func TagNameOnly(ref Named) Named ***REMOVED***
	if IsNameOnly(ref) ***REMOVED***
		namedTagged, err := WithTag(ref, defaultTag)
		if err != nil ***REMOVED***
			// Default tag must be valid, to create a NamedTagged
			// type with non-validated input the WithTag function
			// should be used instead
			panic(err)
		***REMOVED***
		return namedTagged
	***REMOVED***
	return ref
***REMOVED***

// ParseAnyReference parses a reference string as a possible identifier,
// full digest, or familiar name.
func ParseAnyReference(ref string) (Reference, error) ***REMOVED***
	if ok := anchoredIdentifierRegexp.MatchString(ref); ok ***REMOVED***
		return digestReference("sha256:" + ref), nil
	***REMOVED***
	if dgst, err := digest.Parse(ref); err == nil ***REMOVED***
		return digestReference(dgst), nil
	***REMOVED***

	return ParseNormalizedNamed(ref)
***REMOVED***

// ParseAnyReferenceWithSet parses a reference string as a possible short
// identifier to be matched in a digest set, a full digest, or familiar name.
func ParseAnyReferenceWithSet(ref string, ds *digestset.Set) (Reference, error) ***REMOVED***
	if ok := anchoredShortIdentifierRegexp.MatchString(ref); ok ***REMOVED***
		dgst, err := ds.Lookup(ref)
		if err == nil ***REMOVED***
			return digestReference(dgst), nil
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if dgst, err := digest.Parse(ref); err == nil ***REMOVED***
			return digestReference(dgst), nil
		***REMOVED***
	***REMOVED***

	return ParseNormalizedNamed(ref)
***REMOVED***
