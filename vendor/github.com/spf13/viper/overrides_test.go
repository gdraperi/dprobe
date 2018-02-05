package viper

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

type layer int

const (
	defaultLayer layer = iota + 1
	overrideLayer
)

func TestNestedOverrides(t *testing.T) ***REMOVED***
	assert := assert.New(t)
	var v *Viper

	// Case 0: value overridden by a value
	overrideDefault(assert, "tom", 10, "tom", 20) // "tom" is first given 10 as default value, then overridden by 20
	override(assert, "tom", 10, "tom", 20)        // "tom" is first given value 10, then overridden by 20
	overrideDefault(assert, "tom.age", 10, "tom.age", 20)
	override(assert, "tom.age", 10, "tom.age", 20)
	overrideDefault(assert, "sawyer.tom.age", 10, "sawyer.tom.age", 20)
	override(assert, "sawyer.tom.age", 10, "sawyer.tom.age", 20)

	// Case 1: key:value overridden by a value
	v = overrideDefault(assert, "tom.age", 10, "tom", "boy") // "tom.age" is first given 10 as default value, then "tom" is overridden by "boy"
	assert.Nil(v.Get("tom.age"))                             // "tom.age" should not exist anymore
	v = override(assert, "tom.age", 10, "tom", "boy")
	assert.Nil(v.Get("tom.age"))

	// Case 2: value overridden by a key:value
	overrideDefault(assert, "tom", "boy", "tom.age", 10) // "tom" is first given "boy" as default value, then "tom" is overridden by map***REMOVED***"age":10***REMOVED***
	override(assert, "tom.age", 10, "tom", "boy")

	// Case 3: key:value overridden by a key:value
	v = overrideDefault(assert, "tom.size", 4, "tom.age", 10)
	assert.Equal(4, v.Get("tom.size")) // value should still be reachable
	v = override(assert, "tom.size", 4, "tom.age", 10)
	assert.Equal(4, v.Get("tom.size"))
	deepCheckValue(assert, v, overrideLayer, []string***REMOVED***"tom", "size"***REMOVED***, 4)

	// Case 4: key:value overridden by a map
	v = overrideDefault(assert, "tom.size", 4, "tom", map[string]interface***REMOVED******REMOVED******REMOVED***"age": 10***REMOVED***) // "tom.size" is first given "4" as default value, then "tom" is overridden by map***REMOVED***"age":10***REMOVED***
	assert.Equal(4, v.Get("tom.size"))                                                   // "tom.size" should still be reachable
	assert.Equal(10, v.Get("tom.age"))                                                   // new value should be there
	deepCheckValue(assert, v, overrideLayer, []string***REMOVED***"tom", "age"***REMOVED***, 10)                 // new value should be there
	v = override(assert, "tom.size", 4, "tom", map[string]interface***REMOVED******REMOVED******REMOVED***"age": 10***REMOVED***)
	assert.Nil(v.Get("tom.size"))
	assert.Equal(10, v.Get("tom.age"))
	deepCheckValue(assert, v, overrideLayer, []string***REMOVED***"tom", "age"***REMOVED***, 10)

	// Case 5: array overridden by a value
	overrideDefault(assert, "tom", []int***REMOVED***10, 20***REMOVED***, "tom", 30)
	override(assert, "tom", []int***REMOVED***10, 20***REMOVED***, "tom", 30)
	overrideDefault(assert, "tom.age", []int***REMOVED***10, 20***REMOVED***, "tom.age", 30)
	override(assert, "tom.age", []int***REMOVED***10, 20***REMOVED***, "tom.age", 30)

	// Case 6: array overridden by an array
	overrideDefault(assert, "tom", []int***REMOVED***10, 20***REMOVED***, "tom", []int***REMOVED***30, 40***REMOVED***)
	override(assert, "tom", []int***REMOVED***10, 20***REMOVED***, "tom", []int***REMOVED***30, 40***REMOVED***)
	overrideDefault(assert, "tom.age", []int***REMOVED***10, 20***REMOVED***, "tom.age", []int***REMOVED***30, 40***REMOVED***)
	v = override(assert, "tom.age", []int***REMOVED***10, 20***REMOVED***, "tom.age", []int***REMOVED***30, 40***REMOVED***)
	// explicit array merge:
	s, ok := v.Get("tom.age").([]int)
	if assert.True(ok, "tom[\"age\"] is not a slice") ***REMOVED***
		v.Set("tom.age", append(s, []int***REMOVED***50, 60***REMOVED***...))
		assert.Equal([]int***REMOVED***30, 40, 50, 60***REMOVED***, v.Get("tom.age"))
		deepCheckValue(assert, v, overrideLayer, []string***REMOVED***"tom", "age"***REMOVED***, []int***REMOVED***30, 40, 50, 60***REMOVED***)
	***REMOVED***
***REMOVED***

func overrideDefault(assert *assert.Assertions, firstPath string, firstValue interface***REMOVED******REMOVED***, secondPath string, secondValue interface***REMOVED******REMOVED***) *Viper ***REMOVED***
	return overrideFromLayer(defaultLayer, assert, firstPath, firstValue, secondPath, secondValue)
***REMOVED***
func override(assert *assert.Assertions, firstPath string, firstValue interface***REMOVED******REMOVED***, secondPath string, secondValue interface***REMOVED******REMOVED***) *Viper ***REMOVED***
	return overrideFromLayer(overrideLayer, assert, firstPath, firstValue, secondPath, secondValue)
***REMOVED***

// overrideFromLayer performs the sequential override and low-level checks.
//
// First assignment is made on layer l for path firstPath with value firstValue,
// the second one on the override layer (i.e., with the Set() function)
// for path secondPath with value secondValue.
//
// firstPath and secondPath can include an arbitrary number of dots to indicate
// a nested element.
//
// After each assignment, the value is checked, retrieved both by its full path
// and by its key sequence (successive maps).
func overrideFromLayer(l layer, assert *assert.Assertions, firstPath string, firstValue interface***REMOVED******REMOVED***, secondPath string, secondValue interface***REMOVED******REMOVED***) *Viper ***REMOVED***
	v := New()
	firstKeys := strings.Split(firstPath, v.keyDelim)
	if assert == nil ||
		len(firstKeys) == 0 || len(firstKeys[0]) == 0 ***REMOVED***
		return v
	***REMOVED***

	// Set and check first value
	switch l ***REMOVED***
	case defaultLayer:
		v.SetDefault(firstPath, firstValue)
	case overrideLayer:
		v.Set(firstPath, firstValue)
	default:
		return v
	***REMOVED***
	assert.Equal(firstValue, v.Get(firstPath))
	deepCheckValue(assert, v, l, firstKeys, firstValue)

	// Override and check new value
	secondKeys := strings.Split(secondPath, v.keyDelim)
	if len(secondKeys) == 0 || len(secondKeys[0]) == 0 ***REMOVED***
		return v
	***REMOVED***
	v.Set(secondPath, secondValue)
	assert.Equal(secondValue, v.Get(secondPath))
	deepCheckValue(assert, v, overrideLayer, secondKeys, secondValue)

	return v
***REMOVED***

// deepCheckValue checks that all given keys correspond to a valid path in the
// configuration map of the given layer, and that the final value equals the one given
func deepCheckValue(assert *assert.Assertions, v *Viper, l layer, keys []string, value interface***REMOVED******REMOVED***) ***REMOVED***
	if assert == nil || v == nil ||
		len(keys) == 0 || len(keys[0]) == 0 ***REMOVED***
		return
	***REMOVED***

	// init
	var val interface***REMOVED******REMOVED***
	var ms string
	switch l ***REMOVED***
	case defaultLayer:
		val = v.defaults
		ms = "v.defaults"
	case overrideLayer:
		val = v.override
		ms = "v.override"
	***REMOVED***

	// loop through map
	var m map[string]interface***REMOVED******REMOVED***
	err := false
	for _, k := range keys ***REMOVED***
		if val == nil ***REMOVED***
			assert.Fail(fmt.Sprintf("%s is not a map[string]interface***REMOVED******REMOVED***", ms))
			return
		***REMOVED***

		// deep scan of the map to get the final value
		switch val.(type) ***REMOVED***
		case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
			m = cast.ToStringMap(val)
		case map[string]interface***REMOVED******REMOVED***:
			m = val.(map[string]interface***REMOVED******REMOVED***)
		default:
			assert.Fail(fmt.Sprintf("%s is not a map[string]interface***REMOVED******REMOVED***", ms))
			return
		***REMOVED***
		ms = ms + "[\"" + k + "\"]"
		val = m[k]
	***REMOVED***
	if !err ***REMOVED***
		assert.Equal(value, val)
	***REMOVED***
***REMOVED***
