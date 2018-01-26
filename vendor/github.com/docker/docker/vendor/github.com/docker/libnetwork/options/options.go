// Package options provides a way to pass unstructured sets of options to a
// component expecting a strongly-typed configuration structure.
package options

import (
	"fmt"
	"reflect"
)

// NoSuchFieldError is the error returned when the generic parameters hold a
// value for a field absent from the destination structure.
type NoSuchFieldError struct ***REMOVED***
	Field string
	Type  string
***REMOVED***

func (e NoSuchFieldError) Error() string ***REMOVED***
	return fmt.Sprintf("no field %q in type %q", e.Field, e.Type)
***REMOVED***

// CannotSetFieldError is the error returned when the generic parameters hold a
// value for a field that cannot be set in the destination structure.
type CannotSetFieldError struct ***REMOVED***
	Field string
	Type  string
***REMOVED***

func (e CannotSetFieldError) Error() string ***REMOVED***
	return fmt.Sprintf("cannot set field %q of type %q", e.Field, e.Type)
***REMOVED***

// TypeMismatchError is the error returned when the type of the generic value
// for a field mismatches the type of the destination structure.
type TypeMismatchError struct ***REMOVED***
	Field      string
	ExpectType string
	ActualType string
***REMOVED***

func (e TypeMismatchError) Error() string ***REMOVED***
	return fmt.Sprintf("type mismatch, field %s require type %v, actual type %v", e.Field, e.ExpectType, e.ActualType)
***REMOVED***

// Generic is a basic type to store arbitrary settings.
type Generic map[string]interface***REMOVED******REMOVED***

// NewGeneric returns a new Generic instance.
func NewGeneric() Generic ***REMOVED***
	return make(Generic)
***REMOVED***

// GenerateFromModel takes the generic options, and tries to build a new
// instance of the model's type by matching keys from the generic options to
// fields in the model.
//
// The return value is of the same type than the model (including a potential
// pointer qualifier).
func GenerateFromModel(options Generic, model interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	modType := reflect.TypeOf(model)

	// If the model is of pointer type, we need to dereference for New.
	resType := reflect.TypeOf(model)
	if modType.Kind() == reflect.Ptr ***REMOVED***
		resType = resType.Elem()
	***REMOVED***

	// Populate the result structure with the generic layout content.
	res := reflect.New(resType)
	for name, value := range options ***REMOVED***
		field := res.Elem().FieldByName(name)
		if !field.IsValid() ***REMOVED***
			return nil, NoSuchFieldError***REMOVED***name, resType.String()***REMOVED***
		***REMOVED***
		if !field.CanSet() ***REMOVED***
			return nil, CannotSetFieldError***REMOVED***name, resType.String()***REMOVED***
		***REMOVED***
		if reflect.TypeOf(value) != field.Type() ***REMOVED***
			return nil, TypeMismatchError***REMOVED***name, field.Type().String(), reflect.TypeOf(value).String()***REMOVED***
		***REMOVED***
		field.Set(reflect.ValueOf(value))
	***REMOVED***

	// If the model is not of pointer type, return content of the result.
	if modType.Kind() == reflect.Ptr ***REMOVED***
		return res.Interface(), nil
	***REMOVED***
	return res.Elem().Interface(), nil
***REMOVED***
