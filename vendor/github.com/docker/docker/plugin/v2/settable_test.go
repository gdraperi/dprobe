package v2

import (
	"reflect"
	"testing"
)

func TestNewSettable(t *testing.T) ***REMOVED***
	contexts := []struct ***REMOVED***
		arg   string
		name  string
		field string
		value string
		err   error
	***REMOVED******REMOVED***
		***REMOVED***"name=value", "name", "", "value", nil***REMOVED***,
		***REMOVED***"name", "name", "", "", nil***REMOVED***,
		***REMOVED***"name.field=value", "name", "field", "value", nil***REMOVED***,
		***REMOVED***"name.field", "name", "field", "", nil***REMOVED***,
		***REMOVED***"=value", "", "", "", errInvalidFormat***REMOVED***,
		***REMOVED***"=", "", "", "", errInvalidFormat***REMOVED***,
	***REMOVED***

	for _, c := range contexts ***REMOVED***
		s, err := newSettable(c.arg)
		if err != c.err ***REMOVED***
			t.Fatalf("expected error to be %v, got %v", c.err, err)
		***REMOVED***

		if s.name != c.name ***REMOVED***
			t.Fatalf("expected name to be %q, got %q", c.name, s.name)
		***REMOVED***

		if s.field != c.field ***REMOVED***
			t.Fatalf("expected field to be %q, got %q", c.field, s.field)
		***REMOVED***

		if s.value != c.value ***REMOVED***
			t.Fatalf("expected value to be %q, got %q", c.value, s.value)
		***REMOVED***

	***REMOVED***
***REMOVED***

func TestIsSettable(t *testing.T) ***REMOVED***
	contexts := []struct ***REMOVED***
		allowedSettableFields []string
		set                   settable
		settable              []string
		result                bool
		err                   error
	***REMOVED******REMOVED***
		***REMOVED***allowedSettableFieldsEnv, settable***REMOVED******REMOVED***, []string***REMOVED******REMOVED***, false, nil***REMOVED***,
		***REMOVED***allowedSettableFieldsEnv, settable***REMOVED***field: "value"***REMOVED***, []string***REMOVED******REMOVED***, false, nil***REMOVED***,
		***REMOVED***allowedSettableFieldsEnv, settable***REMOVED******REMOVED***, []string***REMOVED***"value"***REMOVED***, true, nil***REMOVED***,
		***REMOVED***allowedSettableFieldsEnv, settable***REMOVED***field: "value"***REMOVED***, []string***REMOVED***"value"***REMOVED***, true, nil***REMOVED***,
		***REMOVED***allowedSettableFieldsEnv, settable***REMOVED***field: "foo"***REMOVED***, []string***REMOVED***"value"***REMOVED***, false, nil***REMOVED***,
		***REMOVED***allowedSettableFieldsEnv, settable***REMOVED***field: "foo"***REMOVED***, []string***REMOVED***"foo"***REMOVED***, false, nil***REMOVED***,
		***REMOVED***allowedSettableFieldsEnv, settable***REMOVED******REMOVED***, []string***REMOVED***"value1", "value2"***REMOVED***, false, errMultipleFields***REMOVED***,
	***REMOVED***

	for _, c := range contexts ***REMOVED***
		if res, err := c.set.isSettable(c.allowedSettableFields, c.settable); res != c.result ***REMOVED***
			t.Fatalf("expected result to be %t, got %t", c.result, res)
		***REMOVED*** else if err != c.err ***REMOVED***
			t.Fatalf("expected error to be %v, got %v", c.err, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUpdateSettingsEnv(t *testing.T) ***REMOVED***
	contexts := []struct ***REMOVED***
		env    []string
		set    settable
		newEnv []string
	***REMOVED******REMOVED***
		***REMOVED***[]string***REMOVED******REMOVED***, settable***REMOVED***name: "DEBUG", value: "1"***REMOVED***, []string***REMOVED***"DEBUG=1"***REMOVED******REMOVED***,
		***REMOVED***[]string***REMOVED***"DEBUG=0"***REMOVED***, settable***REMOVED***name: "DEBUG", value: "1"***REMOVED***, []string***REMOVED***"DEBUG=1"***REMOVED******REMOVED***,
		***REMOVED***[]string***REMOVED***"FOO=0"***REMOVED***, settable***REMOVED***name: "DEBUG", value: "1"***REMOVED***, []string***REMOVED***"FOO=0", "DEBUG=1"***REMOVED******REMOVED***,
		***REMOVED***[]string***REMOVED***"FOO=0", "DEBUG=0"***REMOVED***, settable***REMOVED***name: "DEBUG", value: "1"***REMOVED***, []string***REMOVED***"FOO=0", "DEBUG=1"***REMOVED******REMOVED***,
		***REMOVED***[]string***REMOVED***"FOO=0", "DEBUG=0", "BAR=1"***REMOVED***, settable***REMOVED***name: "DEBUG", value: "1"***REMOVED***, []string***REMOVED***"FOO=0", "DEBUG=1", "BAR=1"***REMOVED******REMOVED***,
	***REMOVED***

	for _, c := range contexts ***REMOVED***
		updateSettingsEnv(&c.env, &c.set)

		if !reflect.DeepEqual(c.env, c.newEnv) ***REMOVED***
			t.Fatalf("expected env to be %q, got %q", c.newEnv, c.env)
		***REMOVED***
	***REMOVED***
***REMOVED***
