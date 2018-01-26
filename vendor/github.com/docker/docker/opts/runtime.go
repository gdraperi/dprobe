package opts

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
)

// RuntimeOpt defines a map of Runtimes
type RuntimeOpt struct ***REMOVED***
	name             string
	stockRuntimeName string
	values           *map[string]types.Runtime
***REMOVED***

// NewNamedRuntimeOpt creates a new RuntimeOpt
func NewNamedRuntimeOpt(name string, ref *map[string]types.Runtime, stockRuntime string) *RuntimeOpt ***REMOVED***
	if ref == nil ***REMOVED***
		ref = &map[string]types.Runtime***REMOVED******REMOVED***
	***REMOVED***
	return &RuntimeOpt***REMOVED***name: name, values: ref, stockRuntimeName: stockRuntime***REMOVED***
***REMOVED***

// Name returns the name of the NamedListOpts in the configuration.
func (o *RuntimeOpt) Name() string ***REMOVED***
	return o.name
***REMOVED***

// Set validates and updates the list of Runtimes
func (o *RuntimeOpt) Set(val string) error ***REMOVED***
	parts := strings.SplitN(val, "=", 2)
	if len(parts) != 2 ***REMOVED***
		return fmt.Errorf("invalid runtime argument: %s", val)
	***REMOVED***

	parts[0] = strings.TrimSpace(parts[0])
	parts[1] = strings.TrimSpace(parts[1])
	if parts[0] == "" || parts[1] == "" ***REMOVED***
		return fmt.Errorf("invalid runtime argument: %s", val)
	***REMOVED***

	parts[0] = strings.ToLower(parts[0])
	if parts[0] == o.stockRuntimeName ***REMOVED***
		return fmt.Errorf("runtime name '%s' is reserved", o.stockRuntimeName)
	***REMOVED***

	if _, ok := (*o.values)[parts[0]]; ok ***REMOVED***
		return fmt.Errorf("runtime '%s' was already defined", parts[0])
	***REMOVED***

	(*o.values)[parts[0]] = types.Runtime***REMOVED***Path: parts[1]***REMOVED***

	return nil
***REMOVED***

// String returns Runtime values as a string.
func (o *RuntimeOpt) String() string ***REMOVED***
	var out []string
	for k := range *o.values ***REMOVED***
		out = append(out, k)
	***REMOVED***

	return fmt.Sprintf("%v", out)
***REMOVED***

// GetMap returns a map of Runtimes (name: path)
func (o *RuntimeOpt) GetMap() map[string]types.Runtime ***REMOVED***
	if o.values != nil ***REMOVED***
		return *o.values
	***REMOVED***

	return map[string]types.Runtime***REMOVED******REMOVED***
***REMOVED***

// Type returns the type of the option
func (o *RuntimeOpt) Type() string ***REMOVED***
	return "runtime"
***REMOVED***
