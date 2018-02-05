// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"fmt"
	"html/template"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToUintE(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect uint
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8.31), 8, false***REMOVED***,
		***REMOVED***float64(8.31), 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***nil, 0, false***REMOVED***,
		// errors
		***REMOVED***int(-8), 0, true***REMOVED***,
		***REMOVED***int8(-8), 0, true***REMOVED***,
		***REMOVED***int16(-8), 0, true***REMOVED***,
		***REMOVED***int32(-8), 0, true***REMOVED***,
		***REMOVED***int64(-8), 0, true***REMOVED***,
		***REMOVED***float32(-8.31), 0, true***REMOVED***,
		***REMOVED***float64(-8.31), 0, true***REMOVED***,
		***REMOVED***"-8", 0, true***REMOVED***,
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToUintE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test:
		v = ToUint(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToUint64E(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect uint64
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8.31), 8, false***REMOVED***,
		***REMOVED***float64(8.31), 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***nil, 0, false***REMOVED***,
		// errors
		***REMOVED***int(-8), 0, true***REMOVED***,
		***REMOVED***int8(-8), 0, true***REMOVED***,
		***REMOVED***int16(-8), 0, true***REMOVED***,
		***REMOVED***int32(-8), 0, true***REMOVED***,
		***REMOVED***int64(-8), 0, true***REMOVED***,
		***REMOVED***float32(-8.31), 0, true***REMOVED***,
		***REMOVED***float64(-8.31), 0, true***REMOVED***,
		***REMOVED***"-8", 0, true***REMOVED***,
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToUint64E(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test:
		v = ToUint64(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToUint32E(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect uint32
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8.31), 8, false***REMOVED***,
		***REMOVED***float64(8.31), 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***nil, 0, false***REMOVED***,
		***REMOVED***int(-8), 0, true***REMOVED***,
		***REMOVED***int8(-8), 0, true***REMOVED***,
		***REMOVED***int16(-8), 0, true***REMOVED***,
		***REMOVED***int32(-8), 0, true***REMOVED***,
		***REMOVED***int64(-8), 0, true***REMOVED***,
		***REMOVED***float32(-8.31), 0, true***REMOVED***,
		***REMOVED***float64(-8.31), 0, true***REMOVED***,
		***REMOVED***"-8", 0, true***REMOVED***,
		// errors
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToUint32E(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test:
		v = ToUint32(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToUint16E(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect uint16
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8.31), 8, false***REMOVED***,
		***REMOVED***float64(8.31), 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***nil, 0, false***REMOVED***,
		// errors
		***REMOVED***int(-8), 0, true***REMOVED***,
		***REMOVED***int8(-8), 0, true***REMOVED***,
		***REMOVED***int16(-8), 0, true***REMOVED***,
		***REMOVED***int32(-8), 0, true***REMOVED***,
		***REMOVED***int64(-8), 0, true***REMOVED***,
		***REMOVED***float32(-8.31), 0, true***REMOVED***,
		***REMOVED***float64(-8.31), 0, true***REMOVED***,
		***REMOVED***"-8", 0, true***REMOVED***,
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToUint16E(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToUint16(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToUint8E(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect uint8
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8.31), 8, false***REMOVED***,
		***REMOVED***float64(8.31), 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***nil, 0, false***REMOVED***,
		// errors
		***REMOVED***int(-8), 0, true***REMOVED***,
		***REMOVED***int8(-8), 0, true***REMOVED***,
		***REMOVED***int16(-8), 0, true***REMOVED***,
		***REMOVED***int32(-8), 0, true***REMOVED***,
		***REMOVED***int64(-8), 0, true***REMOVED***,
		***REMOVED***float32(-8.31), 0, true***REMOVED***,
		***REMOVED***float64(-8.31), 0, true***REMOVED***,
		***REMOVED***"-8", 0, true***REMOVED***,
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToUint8E(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToUint8(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToIntE(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect int
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8.31), 8, false***REMOVED***,
		***REMOVED***float64(8.31), 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***nil, 0, false***REMOVED***,
		// errors
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToIntE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToInt(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToInt64E(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect int64
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8.31), 8, false***REMOVED***,
		***REMOVED***float64(8.31), 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***nil, 0, false***REMOVED***,
		// errors
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToInt64E(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToInt64(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToInt32E(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect int32
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8.31), 8, false***REMOVED***,
		***REMOVED***float64(8.31), 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***nil, 0, false***REMOVED***,
		// errors
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToInt32E(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToInt32(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToInt16E(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect int16
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8.31), 8, false***REMOVED***,
		***REMOVED***float64(8.31), 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***nil, 0, false***REMOVED***,
		// errors
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToInt16E(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToInt16(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToInt8E(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect int8
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8.31), 8, false***REMOVED***,
		***REMOVED***float64(8.31), 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***nil, 0, false***REMOVED***,
		// errors
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToInt8E(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToInt8(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToFloat64E(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect float64
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8), 8, false***REMOVED***,
		***REMOVED***float64(8.31), 8.31, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		// errors
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToFloat64E(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToFloat64(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToFloat32E(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect float32
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), 8, false***REMOVED***,
		***REMOVED***int8(8), 8, false***REMOVED***,
		***REMOVED***int16(8), 8, false***REMOVED***,
		***REMOVED***int32(8), 8, false***REMOVED***,
		***REMOVED***int64(8), 8, false***REMOVED***,
		***REMOVED***uint(8), 8, false***REMOVED***,
		***REMOVED***uint8(8), 8, false***REMOVED***,
		***REMOVED***uint16(8), 8, false***REMOVED***,
		***REMOVED***uint32(8), 8, false***REMOVED***,
		***REMOVED***uint64(8), 8, false***REMOVED***,
		***REMOVED***float32(8.31), 8.31, false***REMOVED***,
		***REMOVED***float64(8.31), 8.31, false***REMOVED***,
		***REMOVED***"8", 8, false***REMOVED***,
		***REMOVED***true, 1, false***REMOVED***,
		***REMOVED***false, 0, false***REMOVED***,
		// errors
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToFloat32E(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToFloat32(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToStringE(t *testing.T) ***REMOVED***
	type Key struct ***REMOVED***
		k string
	***REMOVED***
	key := &Key***REMOVED***"foo"***REMOVED***

	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect string
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***int(8), "8", false***REMOVED***,
		***REMOVED***int8(8), "8", false***REMOVED***,
		***REMOVED***int16(8), "8", false***REMOVED***,
		***REMOVED***int32(8), "8", false***REMOVED***,
		***REMOVED***int64(8), "8", false***REMOVED***,
		***REMOVED***uint(8), "8", false***REMOVED***,
		***REMOVED***uint8(8), "8", false***REMOVED***,
		***REMOVED***uint16(8), "8", false***REMOVED***,
		***REMOVED***uint32(8), "8", false***REMOVED***,
		***REMOVED***uint64(8), "8", false***REMOVED***,
		***REMOVED***float32(8.31), "8.31", false***REMOVED***,
		***REMOVED***float64(8.31), "8.31", false***REMOVED***,
		***REMOVED***true, "true", false***REMOVED***,
		***REMOVED***false, "false", false***REMOVED***,
		***REMOVED***nil, "", false***REMOVED***,
		***REMOVED***[]byte("one time"), "one time", false***REMOVED***,
		***REMOVED***"one more time", "one more time", false***REMOVED***,
		***REMOVED***template.HTML("one time"), "one time", false***REMOVED***,
		***REMOVED***template.URL("http://somehost.foo"), "http://somehost.foo", false***REMOVED***,
		***REMOVED***template.JS("(1+2)"), "(1+2)", false***REMOVED***,
		***REMOVED***template.CSS("a"), "a", false***REMOVED***,
		***REMOVED***template.HTMLAttr("a"), "a", false***REMOVED***,
		// errors
		***REMOVED***testing.T***REMOVED******REMOVED***, "", true***REMOVED***,
		***REMOVED***key, "", true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToString(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

type foo struct ***REMOVED***
	val string
***REMOVED***

func (x foo) String() string ***REMOVED***
	return x.val
***REMOVED***

func TestStringerToString(t *testing.T) ***REMOVED***
	var x foo
	x.val = "bar"
	assert.Equal(t, "bar", ToString(x))
***REMOVED***

type fu struct ***REMOVED***
	val string
***REMOVED***

func (x fu) Error() string ***REMOVED***
	return x.val
***REMOVED***

func TestErrorToString(t *testing.T) ***REMOVED***
	var x fu
	x.val = "bar"
	assert.Equal(t, "bar", ToString(x))
***REMOVED***

func TestStringMapStringSliceE(t *testing.T) ***REMOVED***
	// ToStringMapString inputs/outputs
	var stringMapString = map[string]string***REMOVED***"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"***REMOVED***
	var stringMapInterface = map[string]interface***REMOVED******REMOVED******REMOVED***"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"***REMOVED***
	var interfaceMapString = map[interface***REMOVED******REMOVED***]string***REMOVED***"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"***REMOVED***
	var interfaceMapInterface = map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"***REMOVED***

	// ToStringMapStringSlice inputs/outputs
	var stringMapStringSlice = map[string][]string***REMOVED***"key 1": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 2": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 3": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED******REMOVED***
	var stringMapInterfaceSlice = map[string][]interface***REMOVED******REMOVED******REMOVED***"key 1": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 2": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 3": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED******REMOVED***
	var stringMapInterfaceInterfaceSlice = map[string]interface***REMOVED******REMOVED******REMOVED***"key 1": []interface***REMOVED******REMOVED******REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 2": []interface***REMOVED******REMOVED******REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 3": []interface***REMOVED******REMOVED******REMOVED***"value 1", "value 2", "value 3"***REMOVED******REMOVED***
	var stringMapStringSingleSliceFieldsResult = map[string][]string***REMOVED***"key 1": ***REMOVED***"value", "1"***REMOVED***, "key 2": ***REMOVED***"value", "2"***REMOVED***, "key 3": ***REMOVED***"value", "3"***REMOVED******REMOVED***
	var interfaceMapStringSlice = map[interface***REMOVED******REMOVED***][]string***REMOVED***"key 1": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 2": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 3": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED******REMOVED***
	var interfaceMapInterfaceSlice = map[interface***REMOVED******REMOVED***][]interface***REMOVED******REMOVED******REMOVED***"key 1": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 2": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 3": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED******REMOVED***

	var stringMapStringSliceMultiple = map[string][]string***REMOVED***"key 1": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 2": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED***, "key 3": ***REMOVED***"value 1", "value 2", "value 3"***REMOVED******REMOVED***
	var stringMapStringSliceSingle = map[string][]string***REMOVED***"key 1": ***REMOVED***"value 1"***REMOVED***, "key 2": ***REMOVED***"value 2"***REMOVED***, "key 3": ***REMOVED***"value 3"***REMOVED******REMOVED***

	var stringMapInterface1 = map[string]interface***REMOVED******REMOVED******REMOVED***"key 1": []string***REMOVED***"value 1"***REMOVED***, "key 2": []string***REMOVED***"value 2"***REMOVED******REMOVED***
	var stringMapInterfaceResult1 = map[string][]string***REMOVED***"key 1": ***REMOVED***"value 1"***REMOVED***, "key 2": ***REMOVED***"value 2"***REMOVED******REMOVED***

	type Key struct ***REMOVED***
		k string
	***REMOVED***

	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect map[string][]string
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***stringMapStringSlice, stringMapStringSlice, false***REMOVED***,
		***REMOVED***stringMapInterfaceSlice, stringMapStringSlice, false***REMOVED***,
		***REMOVED***stringMapInterfaceInterfaceSlice, stringMapStringSlice, false***REMOVED***,
		***REMOVED***stringMapStringSliceMultiple, stringMapStringSlice, false***REMOVED***,
		***REMOVED***stringMapStringSliceMultiple, stringMapStringSlice, false***REMOVED***,
		***REMOVED***stringMapString, stringMapStringSliceSingle, false***REMOVED***,
		***REMOVED***stringMapInterface, stringMapStringSliceSingle, false***REMOVED***,
		***REMOVED***stringMapInterface1, stringMapInterfaceResult1, false***REMOVED***,
		***REMOVED***interfaceMapStringSlice, stringMapStringSlice, false***REMOVED***,
		***REMOVED***interfaceMapInterfaceSlice, stringMapStringSlice, false***REMOVED***,
		***REMOVED***interfaceMapString, stringMapStringSingleSliceFieldsResult, false***REMOVED***,
		***REMOVED***interfaceMapInterface, stringMapStringSingleSliceFieldsResult, false***REMOVED***,
		// errors
		***REMOVED***nil, nil, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, nil, true***REMOVED***,
		***REMOVED***map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"foo": testing.T***REMOVED******REMOVED******REMOVED***, nil, true***REMOVED***,
		***REMOVED***map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***Key***REMOVED***"foo"***REMOVED***: "bar"***REMOVED***, nil, true***REMOVED***, // ToStringE(Key***REMOVED***"foo"***REMOVED***) should fail
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringMapStringSliceE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToStringMapStringSlice(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToStringMapE(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect map[string]interface***REMOVED******REMOVED***
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"tag": "tags", "group": "groups"***REMOVED***, map[string]interface***REMOVED******REMOVED******REMOVED***"tag": "tags", "group": "groups"***REMOVED***, false***REMOVED***,
		***REMOVED***map[string]interface***REMOVED******REMOVED******REMOVED***"tag": "tags", "group": "groups"***REMOVED***, map[string]interface***REMOVED******REMOVED******REMOVED***"tag": "tags", "group": "groups"***REMOVED***, false***REMOVED***,
		// errors
		***REMOVED***nil, nil, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, nil, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringMapE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToStringMap(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToStringMapBoolE(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect map[string]bool
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"v1": true, "v2": false***REMOVED***, map[string]bool***REMOVED***"v1": true, "v2": false***REMOVED***, false***REMOVED***,
		***REMOVED***map[string]interface***REMOVED******REMOVED******REMOVED***"v1": true, "v2": false***REMOVED***, map[string]bool***REMOVED***"v1": true, "v2": false***REMOVED***, false***REMOVED***,
		***REMOVED***map[string]bool***REMOVED***"v1": true, "v2": false***REMOVED***, map[string]bool***REMOVED***"v1": true, "v2": false***REMOVED***, false***REMOVED***,
		// errors
		***REMOVED***nil, nil, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, nil, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringMapBoolE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToStringMapBool(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToStringMapStringE(t *testing.T) ***REMOVED***
	var stringMapString = map[string]string***REMOVED***"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"***REMOVED***
	var stringMapInterface = map[string]interface***REMOVED******REMOVED******REMOVED***"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"***REMOVED***
	var interfaceMapString = map[interface***REMOVED******REMOVED***]string***REMOVED***"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"***REMOVED***
	var interfaceMapInterface = map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"***REMOVED***

	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect map[string]string
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***stringMapString, stringMapString, false***REMOVED***,
		***REMOVED***stringMapInterface, stringMapString, false***REMOVED***,
		***REMOVED***interfaceMapString, stringMapString, false***REMOVED***,
		***REMOVED***interfaceMapInterface, stringMapString, false***REMOVED***,
		// errors
		***REMOVED***nil, nil, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, nil, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringMapStringE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToStringMapString(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToBoolSliceE(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect []bool
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***[]bool***REMOVED***true, false, true***REMOVED***, []bool***REMOVED***true, false, true***REMOVED***, false***REMOVED***,
		***REMOVED***[]interface***REMOVED******REMOVED******REMOVED***true, false, true***REMOVED***, []bool***REMOVED***true, false, true***REMOVED***, false***REMOVED***,
		***REMOVED***[]int***REMOVED***1, 0, 1***REMOVED***, []bool***REMOVED***true, false, true***REMOVED***, false***REMOVED***,
		***REMOVED***[]string***REMOVED***"true", "false", "true"***REMOVED***, []bool***REMOVED***true, false, true***REMOVED***, false***REMOVED***,
		// errors
		***REMOVED***nil, nil, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, nil, true***REMOVED***,
		***REMOVED***[]string***REMOVED***"foo", "bar"***REMOVED***, nil, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToBoolSliceE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToBoolSlice(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToIntSliceE(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect []int
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***[]int***REMOVED***1, 3***REMOVED***, []int***REMOVED***1, 3***REMOVED***, false***REMOVED***,
		***REMOVED***[]interface***REMOVED******REMOVED******REMOVED***1.2, 3.2***REMOVED***, []int***REMOVED***1, 3***REMOVED***, false***REMOVED***,
		***REMOVED***[]string***REMOVED***"2", "3"***REMOVED***, []int***REMOVED***2, 3***REMOVED***, false***REMOVED***,
		***REMOVED***[2]string***REMOVED***"2", "3"***REMOVED***, []int***REMOVED***2, 3***REMOVED***, false***REMOVED***,
		// errors
		***REMOVED***nil, nil, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, nil, true***REMOVED***,
		***REMOVED***[]string***REMOVED***"foo", "bar"***REMOVED***, nil, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToIntSliceE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToIntSlice(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToSliceE(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect []interface***REMOVED******REMOVED***
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***[]interface***REMOVED******REMOVED******REMOVED***1, 3***REMOVED***, []interface***REMOVED******REMOVED******REMOVED***1, 3***REMOVED***, false***REMOVED***,
		***REMOVED***[]map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***"k1": 1***REMOVED***, ***REMOVED***"k2": 2***REMOVED******REMOVED***, []interface***REMOVED******REMOVED******REMOVED***map[string]interface***REMOVED******REMOVED******REMOVED***"k1": 1***REMOVED***, map[string]interface***REMOVED******REMOVED******REMOVED***"k2": 2***REMOVED******REMOVED***, false***REMOVED***,
		// errors
		***REMOVED***nil, nil, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, nil, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToSliceE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToSlice(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToStringSliceE(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect []string
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***[]string***REMOVED***"a", "b"***REMOVED***, []string***REMOVED***"a", "b"***REMOVED***, false***REMOVED***,
		***REMOVED***[]interface***REMOVED******REMOVED******REMOVED***1, 3***REMOVED***, []string***REMOVED***"1", "3"***REMOVED***, false***REMOVED***,
		***REMOVED***interface***REMOVED******REMOVED***(1), []string***REMOVED***"1"***REMOVED***, false***REMOVED***,
		// errors
		***REMOVED***nil, nil, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, nil, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringSliceE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToStringSlice(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToDurationSliceE(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect []time.Duration
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***[]string***REMOVED***"1s", "1m"***REMOVED***, []time.Duration***REMOVED***time.Second, time.Minute***REMOVED***, false***REMOVED***,
		***REMOVED***[]int***REMOVED***1, 2***REMOVED***, []time.Duration***REMOVED***1, 2***REMOVED***, false***REMOVED***,
		***REMOVED***[]interface***REMOVED******REMOVED******REMOVED***1, 3***REMOVED***, []time.Duration***REMOVED***1, 3***REMOVED***, false***REMOVED***,
		// errors
		***REMOVED***nil, nil, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, nil, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToDurationSliceE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToDurationSlice(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func TestToBoolE(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect bool
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***0, false, false***REMOVED***,
		***REMOVED***nil, false, false***REMOVED***,
		***REMOVED***"false", false, false***REMOVED***,
		***REMOVED***"FALSE", false, false***REMOVED***,
		***REMOVED***"False", false, false***REMOVED***,
		***REMOVED***"f", false, false***REMOVED***,
		***REMOVED***"F", false, false***REMOVED***,
		***REMOVED***false, false, false***REMOVED***,

		***REMOVED***"true", true, false***REMOVED***,
		***REMOVED***"TRUE", true, false***REMOVED***,
		***REMOVED***"True", true, false***REMOVED***,
		***REMOVED***"t", true, false***REMOVED***,
		***REMOVED***"T", true, false***REMOVED***,
		***REMOVED***1, true, false***REMOVED***,
		***REMOVED***true, true, false***REMOVED***,
		***REMOVED***-1, true, false***REMOVED***,

		// errors
		***REMOVED***"test", false, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, false, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToBoolE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToBool(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***

func BenchmarkTooBool(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		if !ToBool(true) ***REMOVED***
			b.Fatal("ToBool returned false")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIndirectPointers(t *testing.T) ***REMOVED***
	x := 13
	y := &x
	z := &y

	assert.Equal(t, ToInt(y), 13)
	assert.Equal(t, ToInt(z), 13)
***REMOVED***

func TestToTimeEE(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect time.Time
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***"2009-11-10 23:00:00 +0000 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,   // Time.String()
		***REMOVED***"Tue Nov 10 23:00:00 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,        // ANSIC
		***REMOVED***"Tue Nov 10 23:00:00 UTC 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,    // UnixDate
		***REMOVED***"Tue Nov 10 23:00:00 +0000 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,  // RubyDate
		***REMOVED***"10 Nov 09 23:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,             // RFC822
		***REMOVED***"10 Nov 09 23:00 +0000", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,           // RFC822Z
		***REMOVED***"Tuesday, 10-Nov-09 23:00:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***, // RFC850
		***REMOVED***"Tue, 10 Nov 2009 23:00:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,   // RFC1123
		***REMOVED***"Tue, 10 Nov 2009 23:00:00 +0000", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***, // RFC1123Z
		***REMOVED***"2009-11-10T23:00:00Z", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,            // RFC3339
		***REMOVED***"2009-11-10T23:00:00Z", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,            // RFC3339Nano
		***REMOVED***"11:00PM", time.Date(0, 1, 1, 23, 0, 0, 0, time.UTC), false***REMOVED***,                              // Kitchen
		***REMOVED***"Nov 10 23:00:00", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,                    // Stamp
		***REMOVED***"Nov 10 23:00:00.000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,                // StampMilli
		***REMOVED***"Nov 10 23:00:00.000000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,             // StampMicro
		***REMOVED***"Nov 10 23:00:00.000000000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false***REMOVED***,          // StampNano
		***REMOVED***"2016-03-06 15:28:01-00:00", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false***REMOVED***,        // RFC3339 without T
		***REMOVED***"2016-03-06 15:28:01", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false***REMOVED***,
		***REMOVED***"2016-03-06 15:28:01 -0000", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false***REMOVED***,
		***REMOVED***"2016-03-06 15:28:01 -00:00", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false***REMOVED***,
		***REMOVED***"2006-01-02", time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), false***REMOVED***,
		***REMOVED***"02 Jan 2006", time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), false***REMOVED***,
		***REMOVED***1472574600, time.Date(2016, 8, 30, 16, 30, 0, 0, time.UTC), false***REMOVED***,
		***REMOVED***int(1482597504), time.Date(2016, 12, 24, 16, 38, 24, 0, time.UTC), false***REMOVED***,
		***REMOVED***int64(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false***REMOVED***,
		***REMOVED***int32(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false***REMOVED***,
		***REMOVED***uint(1482597504), time.Date(2016, 12, 24, 16, 38, 24, 0, time.UTC), false***REMOVED***,
		***REMOVED***uint64(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false***REMOVED***,
		***REMOVED***uint32(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false***REMOVED***,
		***REMOVED***time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false***REMOVED***,
		// errors
		***REMOVED***"2006", time.Time***REMOVED******REMOVED***, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, time.Time***REMOVED******REMOVED***, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToTimeE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v.UTC(), errmsg)

		// Non-E test
		v = ToTime(test.input)
		assert.Equal(t, test.expect, v.UTC(), errmsg)
	***REMOVED***
***REMOVED***

func TestToDurationE(t *testing.T) ***REMOVED***
	var td time.Duration = 5

	tests := []struct ***REMOVED***
		input  interface***REMOVED******REMOVED***
		expect time.Duration
		iserr  bool
	***REMOVED******REMOVED***
		***REMOVED***time.Duration(5), td, false***REMOVED***,
		***REMOVED***int(5), td, false***REMOVED***,
		***REMOVED***int64(5), td, false***REMOVED***,
		***REMOVED***int32(5), td, false***REMOVED***,
		***REMOVED***int16(5), td, false***REMOVED***,
		***REMOVED***int8(5), td, false***REMOVED***,
		***REMOVED***uint(5), td, false***REMOVED***,
		***REMOVED***uint64(5), td, false***REMOVED***,
		***REMOVED***uint32(5), td, false***REMOVED***,
		***REMOVED***uint16(5), td, false***REMOVED***,
		***REMOVED***uint8(5), td, false***REMOVED***,
		***REMOVED***float64(5), td, false***REMOVED***,
		***REMOVED***float32(5), td, false***REMOVED***,
		***REMOVED***string("5"), td, false***REMOVED***,
		***REMOVED***string("5ns"), td, false***REMOVED***,
		***REMOVED***string("5us"), time.Microsecond * td, false***REMOVED***,
		***REMOVED***string("5µs"), time.Microsecond * td, false***REMOVED***,
		***REMOVED***string("5ms"), time.Millisecond * td, false***REMOVED***,
		***REMOVED***string("5s"), time.Second * td, false***REMOVED***,
		***REMOVED***string("5m"), time.Minute * td, false***REMOVED***,
		***REMOVED***string("5h"), time.Hour * td, false***REMOVED***,
		// errors
		***REMOVED***"test", 0, true***REMOVED***,
		***REMOVED***testing.T***REMOVED******REMOVED***, 0, true***REMOVED***,
	***REMOVED***

	for i, test := range tests ***REMOVED***
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToDurationE(test.input)
		if test.iserr ***REMOVED***
			assert.Error(t, err, errmsg)
			continue
		***REMOVED***

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToDuration(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	***REMOVED***
***REMOVED***
