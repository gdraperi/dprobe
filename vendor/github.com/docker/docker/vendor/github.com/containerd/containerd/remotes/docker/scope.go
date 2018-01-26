package docker

import (
	"context"
	"net/url"
	"sort"
	"strings"

	"github.com/containerd/containerd/reference"
)

// repositoryScope returns a repository scope string such as "repository:foo/bar:pull"
// for "host/foo/bar:baz".
// When push is true, both pull and push are added to the scope.
func repositoryScope(refspec reference.Spec, push bool) (string, error) ***REMOVED***
	u, err := url.Parse("dummy://" + refspec.Locator)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	s := "repository:" + strings.TrimPrefix(u.Path, "/") + ":pull"
	if push ***REMOVED***
		s += ",push"
	***REMOVED***
	return s, nil
***REMOVED***

// tokenScopesKey is used for the key for context.WithValue().
// value: []string (e.g. ***REMOVED***"registry:foo/bar:pull"***REMOVED***)
type tokenScopesKey struct***REMOVED******REMOVED***

// contextWithRepositoryScope returns a context with tokenScopesKey***REMOVED******REMOVED*** and the repository scope value.
func contextWithRepositoryScope(ctx context.Context, refspec reference.Spec, push bool) (context.Context, error) ***REMOVED***
	s, err := repositoryScope(refspec, push)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return context.WithValue(ctx, tokenScopesKey***REMOVED******REMOVED***, []string***REMOVED***s***REMOVED***), nil
***REMOVED***

// getTokenScopes returns deduplicated and sorted scopes from ctx.Value(tokenScopesKey***REMOVED******REMOVED***) and params["scope"].
func getTokenScopes(ctx context.Context, params map[string]string) []string ***REMOVED***
	var scopes []string
	if x := ctx.Value(tokenScopesKey***REMOVED******REMOVED***); x != nil ***REMOVED***
		scopes = append(scopes, x.([]string)...)
	***REMOVED***
	if scope, ok := params["scope"]; ok ***REMOVED***
		for _, s := range scopes ***REMOVED***
			// Note: this comparison is unaware of the scope grammar (https://docs.docker.com/registry/spec/auth/scope/)
			// So, "repository:foo/bar:pull,push" != "repository:foo/bar:push,pull", although semantically they are equal.
			if s == scope ***REMOVED***
				// already appended
				goto Sort
			***REMOVED***
		***REMOVED***
		scopes = append(scopes, scope)
	***REMOVED***
Sort:
	sort.Strings(scopes)
	return scopes
***REMOVED***
