// Copyright 2017, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package cmp

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/google/go-cmp/cmp/internal/function"
)

// Option configures for specific behavior of Equal and Diff. In particular,
// the fundamental Option functions (Ignore, Transformer, and Comparer),
// configure how equality is determined.
//
// The fundamental options may be composed with filters (FilterPath and
// FilterValues) to control the scope over which they are applied.
//
// The cmp/cmpopts package provides helper functions for creating options that
// may be used with Equal and Diff.
type Option interface ***REMOVED***
	// filter applies all filters and returns the option that remains.
	// Each option may only read s.curPath and call s.callTTBFunc.
	//
	// An Options is returned only if multiple comparers or transformers
	// can apply simultaneously and will only contain values of those types
	// or sub-Options containing values of those types.
	filter(s *state, vx, vy reflect.Value, t reflect.Type) applicableOption
***REMOVED***

// applicableOption represents the following types:
//	Fundamental: ignore | invalid | *comparer | *transformer
//	Grouping:    Options
type applicableOption interface ***REMOVED***
	Option

	// apply executes the option and reports whether the option was applied.
	// Each option may mutate s.
	apply(s *state, vx, vy reflect.Value) bool
***REMOVED***

// coreOption represents the following types:
//	Fundamental: ignore | invalid | *comparer | *transformer
//	Filters:     *pathFilter | *valuesFilter
type coreOption interface ***REMOVED***
	Option
	isCore()
***REMOVED***

type core struct***REMOVED******REMOVED***

func (core) isCore() ***REMOVED******REMOVED***

// Options is a list of Option values that also satisfies the Option interface.
// Helper comparison packages may return an Options value when packing multiple
// Option values into a single Option. When this package processes an Options,
// it will be implicitly expanded into a flat list.
//
// Applying a filter on an Options is equivalent to applying that same filter
// on all individual options held within.
type Options []Option

func (opts Options) filter(s *state, vx, vy reflect.Value, t reflect.Type) (out applicableOption) ***REMOVED***
	for _, opt := range opts ***REMOVED***
		switch opt := opt.filter(s, vx, vy, t); opt.(type) ***REMOVED***
		case ignore:
			return ignore***REMOVED******REMOVED*** // Only ignore can short-circuit evaluation
		case invalid:
			out = invalid***REMOVED******REMOVED*** // Takes precedence over comparer or transformer
		case *comparer, *transformer, Options:
			switch out.(type) ***REMOVED***
			case nil:
				out = opt
			case invalid:
				// Keep invalid
			case *comparer, *transformer, Options:
				out = Options***REMOVED***out, opt***REMOVED*** // Conflicting comparers or transformers
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return out
***REMOVED***

func (opts Options) apply(s *state, _, _ reflect.Value) bool ***REMOVED***
	const warning = "ambiguous set of applicable options"
	const help = "consider using filters to ensure at most one Comparer or Transformer may apply"
	var ss []string
	for _, opt := range flattenOptions(nil, opts) ***REMOVED***
		ss = append(ss, fmt.Sprint(opt))
	***REMOVED***
	set := strings.Join(ss, "\n\t")
	panic(fmt.Sprintf("%s at %#v:\n\t%s\n%s", warning, s.curPath, set, help))
***REMOVED***

func (opts Options) String() string ***REMOVED***
	var ss []string
	for _, opt := range opts ***REMOVED***
		ss = append(ss, fmt.Sprint(opt))
	***REMOVED***
	return fmt.Sprintf("Options***REMOVED***%s***REMOVED***", strings.Join(ss, ", "))
***REMOVED***

// FilterPath returns a new Option where opt is only evaluated if filter f
// returns true for the current Path in the value tree.
//
// The option passed in may be an Ignore, Transformer, Comparer, Options, or
// a previously filtered Option.
func FilterPath(f func(Path) bool, opt Option) Option ***REMOVED***
	if f == nil ***REMOVED***
		panic("invalid path filter function")
	***REMOVED***
	if opt := normalizeOption(opt); opt != nil ***REMOVED***
		return &pathFilter***REMOVED***fnc: f, opt: opt***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type pathFilter struct ***REMOVED***
	core
	fnc func(Path) bool
	opt Option
***REMOVED***

func (f pathFilter) filter(s *state, vx, vy reflect.Value, t reflect.Type) applicableOption ***REMOVED***
	if f.fnc(s.curPath) ***REMOVED***
		return f.opt.filter(s, vx, vy, t)
	***REMOVED***
	return nil
***REMOVED***

func (f pathFilter) String() string ***REMOVED***
	fn := getFuncName(reflect.ValueOf(f.fnc).Pointer())
	return fmt.Sprintf("FilterPath(%s, %v)", fn, f.opt)
***REMOVED***

// FilterValues returns a new Option where opt is only evaluated if filter f,
// which is a function of the form "func(T, T) bool", returns true for the
// current pair of values being compared. If the type of the values is not
// assignable to T, then this filter implicitly returns false.
//
// The filter function must be
// symmetric (i.e., agnostic to the order of the inputs) and
// deterministic (i.e., produces the same result when given the same inputs).
// If T is an interface, it is possible that f is called with two values with
// different concrete types that both implement T.
//
// The option passed in may be an Ignore, Transformer, Comparer, Options, or
// a previously filtered Option.
func FilterValues(f interface***REMOVED******REMOVED***, opt Option) Option ***REMOVED***
	v := reflect.ValueOf(f)
	if !function.IsType(v.Type(), function.ValueFilter) || v.IsNil() ***REMOVED***
		panic(fmt.Sprintf("invalid values filter function: %T", f))
	***REMOVED***
	if opt := normalizeOption(opt); opt != nil ***REMOVED***
		vf := &valuesFilter***REMOVED***fnc: v, opt: opt***REMOVED***
		if ti := v.Type().In(0); ti.Kind() != reflect.Interface || ti.NumMethod() > 0 ***REMOVED***
			vf.typ = ti
		***REMOVED***
		return vf
	***REMOVED***
	return nil
***REMOVED***

type valuesFilter struct ***REMOVED***
	core
	typ reflect.Type  // T
	fnc reflect.Value // func(T, T) bool
	opt Option
***REMOVED***

func (f valuesFilter) filter(s *state, vx, vy reflect.Value, t reflect.Type) applicableOption ***REMOVED***
	if !vx.IsValid() || !vy.IsValid() ***REMOVED***
		return invalid***REMOVED******REMOVED***
	***REMOVED***
	if (f.typ == nil || t.AssignableTo(f.typ)) && s.callTTBFunc(f.fnc, vx, vy) ***REMOVED***
		return f.opt.filter(s, vx, vy, t)
	***REMOVED***
	return nil
***REMOVED***

func (f valuesFilter) String() string ***REMOVED***
	fn := getFuncName(f.fnc.Pointer())
	return fmt.Sprintf("FilterValues(%s, %v)", fn, f.opt)
***REMOVED***

// Ignore is an Option that causes all comparisons to be ignored.
// This value is intended to be combined with FilterPath or FilterValues.
// It is an error to pass an unfiltered Ignore option to Equal.
func Ignore() Option ***REMOVED*** return ignore***REMOVED******REMOVED*** ***REMOVED***

type ignore struct***REMOVED*** core ***REMOVED***

func (ignore) isFiltered() bool                                                     ***REMOVED*** return false ***REMOVED***
func (ignore) filter(_ *state, _, _ reflect.Value, _ reflect.Type) applicableOption ***REMOVED*** return ignore***REMOVED******REMOVED*** ***REMOVED***
func (ignore) apply(_ *state, _, _ reflect.Value) bool                              ***REMOVED*** return true ***REMOVED***
func (ignore) String() string                                                       ***REMOVED*** return "Ignore()" ***REMOVED***

// invalid is a sentinel Option type to indicate that some options could not
// be evaluated due to unexported fields.
type invalid struct***REMOVED*** core ***REMOVED***

func (invalid) filter(_ *state, _, _ reflect.Value, _ reflect.Type) applicableOption ***REMOVED*** return invalid***REMOVED******REMOVED*** ***REMOVED***
func (invalid) apply(s *state, _, _ reflect.Value) bool ***REMOVED***
	const help = "consider using AllowUnexported or cmpopts.IgnoreUnexported"
	panic(fmt.Sprintf("cannot handle unexported field: %#v\n%s", s.curPath, help))
***REMOVED***

// Transformer returns an Option that applies a transformation function that
// converts values of a certain type into that of another.
//
// The transformer f must be a function "func(T) R" that converts values of
// type T to those of type R and is implicitly filtered to input values
// assignable to T. The transformer must not mutate T in any way.
// If T and R are the same type, an additional filter must be applied to
// act as the base case to prevent an infinite recursion applying the same
// transform to itself (see the SortedSlice example).
//
// The name is a user provided label that is used as the Transform.Name in the
// transformation PathStep. If empty, an arbitrary name is used.
func Transformer(name string, f interface***REMOVED******REMOVED***) Option ***REMOVED***
	v := reflect.ValueOf(f)
	if !function.IsType(v.Type(), function.Transformer) || v.IsNil() ***REMOVED***
		panic(fmt.Sprintf("invalid transformer function: %T", f))
	***REMOVED***
	if name == "" ***REMOVED***
		name = "λ" // Lambda-symbol as place-holder for anonymous transformer
	***REMOVED***
	if !isValid(name) ***REMOVED***
		panic(fmt.Sprintf("invalid name: %q", name))
	***REMOVED***
	tr := &transformer***REMOVED***name: name, fnc: reflect.ValueOf(f)***REMOVED***
	if ti := v.Type().In(0); ti.Kind() != reflect.Interface || ti.NumMethod() > 0 ***REMOVED***
		tr.typ = ti
	***REMOVED***
	return tr
***REMOVED***

type transformer struct ***REMOVED***
	core
	name string
	typ  reflect.Type  // T
	fnc  reflect.Value // func(T) R
***REMOVED***

func (tr *transformer) isFiltered() bool ***REMOVED*** return tr.typ != nil ***REMOVED***

func (tr *transformer) filter(_ *state, _, _ reflect.Value, t reflect.Type) applicableOption ***REMOVED***
	if tr.typ == nil || t.AssignableTo(tr.typ) ***REMOVED***
		return tr
	***REMOVED***
	return nil
***REMOVED***

func (tr *transformer) apply(s *state, vx, vy reflect.Value) bool ***REMOVED***
	// Update path before calling the Transformer so that dynamic checks
	// will use the updated path.
	s.curPath.push(&transform***REMOVED***pathStep***REMOVED***tr.fnc.Type().Out(0)***REMOVED***, tr***REMOVED***)
	defer s.curPath.pop()

	vx = s.callTRFunc(tr.fnc, vx)
	vy = s.callTRFunc(tr.fnc, vy)
	s.compareAny(vx, vy)
	return true
***REMOVED***

func (tr transformer) String() string ***REMOVED***
	return fmt.Sprintf("Transformer(%s, %s)", tr.name, getFuncName(tr.fnc.Pointer()))
***REMOVED***

// Comparer returns an Option that determines whether two values are equal
// to each other.
//
// The comparer f must be a function "func(T, T) bool" and is implicitly
// filtered to input values assignable to T. If T is an interface, it is
// possible that f is called with two values of different concrete types that
// both implement T.
//
// The equality function must be:
//	• Symmetric: equal(x, y) == equal(y, x)
//	• Deterministic: equal(x, y) == equal(x, y)
//	• Pure: equal(x, y) does not modify x or y
func Comparer(f interface***REMOVED******REMOVED***) Option ***REMOVED***
	v := reflect.ValueOf(f)
	if !function.IsType(v.Type(), function.Equal) || v.IsNil() ***REMOVED***
		panic(fmt.Sprintf("invalid comparer function: %T", f))
	***REMOVED***
	cm := &comparer***REMOVED***fnc: v***REMOVED***
	if ti := v.Type().In(0); ti.Kind() != reflect.Interface || ti.NumMethod() > 0 ***REMOVED***
		cm.typ = ti
	***REMOVED***
	return cm
***REMOVED***

type comparer struct ***REMOVED***
	core
	typ reflect.Type  // T
	fnc reflect.Value // func(T, T) bool
***REMOVED***

func (cm *comparer) isFiltered() bool ***REMOVED*** return cm.typ != nil ***REMOVED***

func (cm *comparer) filter(_ *state, _, _ reflect.Value, t reflect.Type) applicableOption ***REMOVED***
	if cm.typ == nil || t.AssignableTo(cm.typ) ***REMOVED***
		return cm
	***REMOVED***
	return nil
***REMOVED***

func (cm *comparer) apply(s *state, vx, vy reflect.Value) bool ***REMOVED***
	eq := s.callTTBFunc(cm.fnc, vx, vy)
	s.report(eq, vx, vy)
	return true
***REMOVED***

func (cm comparer) String() string ***REMOVED***
	return fmt.Sprintf("Comparer(%s)", getFuncName(cm.fnc.Pointer()))
***REMOVED***

// AllowUnexported returns an Option that forcibly allows operations on
// unexported fields in certain structs, which are specified by passing in a
// value of each struct type.
//
// Users of this option must understand that comparing on unexported fields
// from external packages is not safe since changes in the internal
// implementation of some external package may cause the result of Equal
// to unexpectedly change. However, it may be valid to use this option on types
// defined in an internal package where the semantic meaning of an unexported
// field is in the control of the user.
//
// For some cases, a custom Comparer should be used instead that defines
// equality as a function of the public API of a type rather than the underlying
// unexported implementation.
//
// For example, the reflect.Type documentation defines equality to be determined
// by the == operator on the interface (essentially performing a shallow pointer
// comparison) and most attempts to compare *regexp.Regexp types are interested
// in only checking that the regular expression strings are equal.
// Both of these are accomplished using Comparers:
//
//	Comparer(func(x, y reflect.Type) bool ***REMOVED*** return x == y ***REMOVED***)
//	Comparer(func(x, y *regexp.Regexp) bool ***REMOVED*** return x.String() == y.String() ***REMOVED***)
//
// In other cases, the cmpopts.IgnoreUnexported option can be used to ignore
// all unexported fields on specified struct types.
func AllowUnexported(types ...interface***REMOVED******REMOVED***) Option ***REMOVED***
	if !supportAllowUnexported ***REMOVED***
		panic("AllowUnexported is not supported on App Engine Classic or GopherJS")
	***REMOVED***
	m := make(map[reflect.Type]bool)
	for _, typ := range types ***REMOVED***
		t := reflect.TypeOf(typ)
		if t.Kind() != reflect.Struct ***REMOVED***
			panic(fmt.Sprintf("invalid struct type: %T", typ))
		***REMOVED***
		m[t] = true
	***REMOVED***
	return visibleStructs(m)
***REMOVED***

type visibleStructs map[reflect.Type]bool

func (visibleStructs) filter(_ *state, _, _ reflect.Value, _ reflect.Type) applicableOption ***REMOVED***
	panic("not implemented")
***REMOVED***

// reporter is an Option that configures how differences are reported.
type reporter interface ***REMOVED***
	// TODO: Not exported yet.
	//
	// Perhaps add PushStep and PopStep and change Report to only accept
	// a PathStep instead of the full-path? Adding a PushStep and PopStep makes
	// it clear that we are traversing the value tree in a depth-first-search
	// manner, which has an effect on how values are printed.

	Option

	// Report is called for every comparison made and will be provided with
	// the two values being compared, the equality result, and the
	// current path in the value tree. It is possible for x or y to be an
	// invalid reflect.Value if one of the values is non-existent;
	// which is possible with maps and slices.
	Report(x, y reflect.Value, eq bool, p Path)
***REMOVED***

// normalizeOption normalizes the input options such that all Options groups
// are flattened and groups with a single element are reduced to that element.
// Only coreOptions and Options containing coreOptions are allowed.
func normalizeOption(src Option) Option ***REMOVED***
	switch opts := flattenOptions(nil, Options***REMOVED***src***REMOVED***); len(opts) ***REMOVED***
	case 0:
		return nil
	case 1:
		return opts[0]
	default:
		return opts
	***REMOVED***
***REMOVED***

// flattenOptions copies all options in src to dst as a flat list.
// Only coreOptions and Options containing coreOptions are allowed.
func flattenOptions(dst, src Options) Options ***REMOVED***
	for _, opt := range src ***REMOVED***
		switch opt := opt.(type) ***REMOVED***
		case nil:
			continue
		case Options:
			dst = flattenOptions(dst, opt)
		case coreOption:
			dst = append(dst, opt)
		default:
			panic(fmt.Sprintf("invalid option type: %T", opt))
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// getFuncName returns a short function name from the pointer.
// The string parsing logic works up until Go1.9.
func getFuncName(p uintptr) string ***REMOVED***
	fnc := runtime.FuncForPC(p)
	if fnc == nil ***REMOVED***
		return "<unknown>"
	***REMOVED***
	name := fnc.Name() // E.g., "long/path/name/mypkg.(mytype).(long/path/name/mypkg.myfunc)-fm"
	if strings.HasSuffix(name, ")-fm") || strings.HasSuffix(name, ")·fm") ***REMOVED***
		// Strip the package name from method name.
		name = strings.TrimSuffix(name, ")-fm")
		name = strings.TrimSuffix(name, ")·fm")
		if i := strings.LastIndexByte(name, '('); i >= 0 ***REMOVED***
			methodName := name[i+1:] // E.g., "long/path/name/mypkg.myfunc"
			if j := strings.LastIndexByte(methodName, '.'); j >= 0 ***REMOVED***
				methodName = methodName[j+1:] // E.g., "myfunc"
			***REMOVED***
			name = name[:i] + methodName // E.g., "long/path/name/mypkg.(mytype)." + "myfunc"
		***REMOVED***
	***REMOVED***
	if i := strings.LastIndexByte(name, '/'); i >= 0 ***REMOVED***
		// Strip the package name.
		name = name[i+1:] // E.g., "mypkg.(mytype).myfunc"
	***REMOVED***
	return name
***REMOVED***
