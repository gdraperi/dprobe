package opts

import (
	"fmt"

	"github.com/docker/go-units"
)

// UlimitOpt defines a map of Ulimits
type UlimitOpt struct ***REMOVED***
	values *map[string]*units.Ulimit
***REMOVED***

// NewUlimitOpt creates a new UlimitOpt
func NewUlimitOpt(ref *map[string]*units.Ulimit) *UlimitOpt ***REMOVED***
	if ref == nil ***REMOVED***
		ref = &map[string]*units.Ulimit***REMOVED******REMOVED***
	***REMOVED***
	return &UlimitOpt***REMOVED***ref***REMOVED***
***REMOVED***

// Set validates a Ulimit and sets its name as a key in UlimitOpt
func (o *UlimitOpt) Set(val string) error ***REMOVED***
	l, err := units.ParseUlimit(val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	(*o.values)[l.Name] = l

	return nil
***REMOVED***

// String returns Ulimit values as a string.
func (o *UlimitOpt) String() string ***REMOVED***
	var out []string
	for _, v := range *o.values ***REMOVED***
		out = append(out, v.String())
	***REMOVED***

	return fmt.Sprintf("%v", out)
***REMOVED***

// GetList returns a slice of pointers to Ulimits.
func (o *UlimitOpt) GetList() []*units.Ulimit ***REMOVED***
	var ulimits []*units.Ulimit
	for _, v := range *o.values ***REMOVED***
		ulimits = append(ulimits, v)
	***REMOVED***

	return ulimits
***REMOVED***

// Type returns the option type
func (o *UlimitOpt) Type() string ***REMOVED***
	return "ulimit"
***REMOVED***

// NamedUlimitOpt defines a named map of Ulimits
type NamedUlimitOpt struct ***REMOVED***
	name string
	UlimitOpt
***REMOVED***

var _ NamedOption = &NamedUlimitOpt***REMOVED******REMOVED***

// NewNamedUlimitOpt creates a new NamedUlimitOpt
func NewNamedUlimitOpt(name string, ref *map[string]*units.Ulimit) *NamedUlimitOpt ***REMOVED***
	if ref == nil ***REMOVED***
		ref = &map[string]*units.Ulimit***REMOVED******REMOVED***
	***REMOVED***
	return &NamedUlimitOpt***REMOVED***
		name:      name,
		UlimitOpt: *NewUlimitOpt(ref),
	***REMOVED***
***REMOVED***

// Name returns the option name
func (o *NamedUlimitOpt) Name() string ***REMOVED***
	return o.name
***REMOVED***
