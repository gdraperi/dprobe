package viper

import "github.com/spf13/pflag"

// FlagValueSet is an interface that users can implement
// to bind a set of flags to viper.
type FlagValueSet interface ***REMOVED***
	VisitAll(fn func(FlagValue))
***REMOVED***

// FlagValue is an interface that users can implement
// to bind different flags to viper.
type FlagValue interface ***REMOVED***
	HasChanged() bool
	Name() string
	ValueString() string
	ValueType() string
***REMOVED***

// pflagValueSet is a wrapper around *pflag.ValueSet
// that implements FlagValueSet.
type pflagValueSet struct ***REMOVED***
	flags *pflag.FlagSet
***REMOVED***

// VisitAll iterates over all *pflag.Flag inside the *pflag.FlagSet.
func (p pflagValueSet) VisitAll(fn func(flag FlagValue)) ***REMOVED***
	p.flags.VisitAll(func(flag *pflag.Flag) ***REMOVED***
		fn(pflagValue***REMOVED***flag***REMOVED***)
	***REMOVED***)
***REMOVED***

// pflagValue is a wrapper aroung *pflag.flag
// that implements FlagValue
type pflagValue struct ***REMOVED***
	flag *pflag.Flag
***REMOVED***

// HasChanges returns whether the flag has changes or not.
func (p pflagValue) HasChanged() bool ***REMOVED***
	return p.flag.Changed
***REMOVED***

// Name returns the name of the flag.
func (p pflagValue) Name() string ***REMOVED***
	return p.flag.Name
***REMOVED***

// ValueString returns the value of the flag as a string.
func (p pflagValue) ValueString() string ***REMOVED***
	return p.flag.Value.String()
***REMOVED***

// ValueType returns the type of the flag as a string.
func (p pflagValue) ValueType() string ***REMOVED***
	return p.flag.Value.Type()
***REMOVED***
