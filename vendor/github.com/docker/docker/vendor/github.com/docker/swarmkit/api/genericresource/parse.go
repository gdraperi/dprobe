package genericresource

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/swarmkit/api"
)

func newParseError(format string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return fmt.Errorf("could not parse GenericResource: "+format, args...)
***REMOVED***

// discreteResourceVal returns an int64 if the string is a discreteResource
// and an error if it isn't
func discreteResourceVal(res string) (int64, error) ***REMOVED***
	return strconv.ParseInt(res, 10, 64)
***REMOVED***

// allNamedResources returns true if the array of resources are all namedResources
// e.g: res = [red, orange, green]
func allNamedResources(res []string) bool ***REMOVED***
	for _, v := range res ***REMOVED***
		if _, err := discreteResourceVal(v); err == nil ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// ParseCmd parses the Generic Resource command line argument
// and returns a list of *api.GenericResource
func ParseCmd(cmd string) ([]*api.GenericResource, error) ***REMOVED***
	if strings.Contains(cmd, "\n") ***REMOVED***
		return nil, newParseError("unexpected '\\n' character")
	***REMOVED***

	r := csv.NewReader(strings.NewReader(cmd))
	records, err := r.ReadAll()

	if err != nil ***REMOVED***
		return nil, newParseError("%v", err)
	***REMOVED***

	if len(records) != 1 ***REMOVED***
		return nil, newParseError("found multiple records while parsing cmd %v", records)
	***REMOVED***

	return Parse(records[0])
***REMOVED***

// Parse parses a table of GenericResource resources
func Parse(cmds []string) ([]*api.GenericResource, error) ***REMOVED***
	tokens := make(map[string][]string)

	for _, term := range cmds ***REMOVED***
		kva := strings.Split(term, "=")
		if len(kva) != 2 ***REMOVED***
			return nil, newParseError("incorrect term %s, missing"+
				" '=' or malformed expression", term)
		***REMOVED***

		key := strings.TrimSpace(kva[0])
		val := strings.TrimSpace(kva[1])

		tokens[key] = append(tokens[key], val)
	***REMOVED***

	var rs []*api.GenericResource
	for k, v := range tokens ***REMOVED***
		if u, ok := isDiscreteResource(v); ok ***REMOVED***
			if u < 0 ***REMOVED***
				return nil, newParseError("cannot ask for"+
					" negative resource %s", k)
			***REMOVED***

			rs = append(rs, NewDiscrete(k, u))
			continue
		***REMOVED***

		if allNamedResources(v) ***REMOVED***
			rs = append(rs, NewSet(k, v...)...)
			continue
		***REMOVED***

		return nil, newParseError("mixed discrete and named resources"+
			" in expression '%s=%s'", k, v)
	***REMOVED***

	return rs, nil
***REMOVED***

// isDiscreteResource returns true if the array of resources is a
// Discrete Resource.
// e.g: res = [1]
func isDiscreteResource(values []string) (int64, bool) ***REMOVED***
	if len(values) != 1 ***REMOVED***
		return int64(0), false
	***REMOVED***

	u, err := discreteResourceVal(values[0])
	if err != nil ***REMOVED***
		return int64(0), false
	***REMOVED***

	return u, true

***REMOVED***
