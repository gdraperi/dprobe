package viper

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestBindFlagValueSet(t *testing.T) ***REMOVED***
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)

	var testValues = map[string]*string***REMOVED***
		"host":     nil,
		"port":     nil,
		"endpoint": nil,
	***REMOVED***

	var mutatedTestValues = map[string]string***REMOVED***
		"host":     "localhost",
		"port":     "6060",
		"endpoint": "/public",
	***REMOVED***

	for name := range testValues ***REMOVED***
		testValues[name] = flagSet.String(name, "", "test")
	***REMOVED***

	flagValueSet := pflagValueSet***REMOVED***flagSet***REMOVED***

	err := BindFlagValues(flagValueSet)
	if err != nil ***REMOVED***
		t.Fatalf("error binding flag set, %v", err)
	***REMOVED***

	flagSet.VisitAll(func(flag *pflag.Flag) ***REMOVED***
		flag.Value.Set(mutatedTestValues[flag.Name])
		flag.Changed = true
	***REMOVED***)

	for name, expected := range mutatedTestValues ***REMOVED***
		assert.Equal(t, Get(name), expected)
	***REMOVED***
***REMOVED***

func TestBindFlagValue(t *testing.T) ***REMOVED***
	var testString = "testing"
	var testValue = newStringValue(testString, &testString)

	flag := &pflag.Flag***REMOVED***
		Name:    "testflag",
		Value:   testValue,
		Changed: false,
	***REMOVED***

	flagValue := pflagValue***REMOVED***flag***REMOVED***
	BindFlagValue("testvalue", flagValue)

	assert.Equal(t, testString, Get("testvalue"))

	flag.Value.Set("testing_mutate")
	flag.Changed = true //hack for pflag usage

	assert.Equal(t, "testing_mutate", Get("testvalue"))
***REMOVED***
