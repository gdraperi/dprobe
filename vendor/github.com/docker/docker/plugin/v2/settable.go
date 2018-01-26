package v2

import (
	"errors"
	"fmt"
	"strings"
)

type settable struct ***REMOVED***
	name  string
	field string
	value string
***REMOVED***

var (
	allowedSettableFieldsEnv     = []string***REMOVED***"value"***REMOVED***
	allowedSettableFieldsArgs    = []string***REMOVED***"value"***REMOVED***
	allowedSettableFieldsDevices = []string***REMOVED***"path"***REMOVED***
	allowedSettableFieldsMounts  = []string***REMOVED***"source"***REMOVED***

	errMultipleFields = errors.New("multiple fields are settable, one must be specified")
	errInvalidFormat  = errors.New("invalid format, must be <name>[.<field>][=<value>]")
)

func newSettables(args []string) ([]settable, error) ***REMOVED***
	sets := make([]settable, 0, len(args))
	for _, arg := range args ***REMOVED***
		set, err := newSettable(arg)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		sets = append(sets, set)
	***REMOVED***
	return sets, nil
***REMOVED***

func newSettable(arg string) (settable, error) ***REMOVED***
	var set settable
	if i := strings.Index(arg, "="); i == 0 ***REMOVED***
		return set, errInvalidFormat
	***REMOVED*** else if i < 0 ***REMOVED***
		set.name = arg
	***REMOVED*** else ***REMOVED***
		set.name = arg[:i]
		set.value = arg[i+1:]
	***REMOVED***

	if i := strings.LastIndex(set.name, "."); i > 0 ***REMOVED***
		set.field = set.name[i+1:]
		set.name = arg[:i]
	***REMOVED***

	return set, nil
***REMOVED***

// prettyName return name.field if there is a field, otherwise name.
func (set *settable) prettyName() string ***REMOVED***
	if set.field != "" ***REMOVED***
		return fmt.Sprintf("%s.%s", set.name, set.field)
	***REMOVED***
	return set.name
***REMOVED***

func (set *settable) isSettable(allowedSettableFields []string, settable []string) (bool, error) ***REMOVED***
	if set.field == "" ***REMOVED***
		if len(settable) == 1 ***REMOVED***
			// if field is not specified and there only one settable, default to it.
			set.field = settable[0]
		***REMOVED*** else if len(settable) > 1 ***REMOVED***
			return false, errMultipleFields
		***REMOVED***
	***REMOVED***

	isAllowed := false
	for _, allowedSettableField := range allowedSettableFields ***REMOVED***
		if set.field == allowedSettableField ***REMOVED***
			isAllowed = true
			break
		***REMOVED***
	***REMOVED***

	if isAllowed ***REMOVED***
		for _, settableField := range settable ***REMOVED***
			if set.field == settableField ***REMOVED***
				return true, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return false, nil
***REMOVED***

func updateSettingsEnv(env *[]string, set *settable) ***REMOVED***
	for i, e := range *env ***REMOVED***
		if parts := strings.SplitN(e, "=", 2); parts[0] == set.name ***REMOVED***
			(*env)[i] = fmt.Sprintf("%s=%s", set.name, set.value)
			return
		***REMOVED***
	***REMOVED***

	*env = append(*env, fmt.Sprintf("%s=%s", set.name, set.value))
***REMOVED***
