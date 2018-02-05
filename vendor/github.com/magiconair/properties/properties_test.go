// Copyright 2017 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

var verbose = flag.Bool("verbose", false, "Verbose output")

func init() ***REMOVED***
	ErrorHandler = PanicHandler
***REMOVED***

// ----------------------------------------------------------------------------

// define test cases in the form of
// ***REMOVED***"input", "key1", "value1", "key2", "value2", ...***REMOVED***
var complexTests = [][]string***REMOVED***
	// whitespace prefix
	***REMOVED***" key=value", "key", "value"***REMOVED***,     // SPACE prefix
	***REMOVED***"\fkey=value", "key", "value"***REMOVED***,    // FF prefix
	***REMOVED***"\tkey=value", "key", "value"***REMOVED***,    // TAB prefix
	***REMOVED***" \f\tkey=value", "key", "value"***REMOVED***, // mix prefix

	// multiple keys
	***REMOVED***"key1=value1\nkey2=value2\n", "key1", "value1", "key2", "value2"***REMOVED***,
	***REMOVED***"key1=value1\rkey2=value2\r", "key1", "value1", "key2", "value2"***REMOVED***,
	***REMOVED***"key1=value1\r\nkey2=value2\r\n", "key1", "value1", "key2", "value2"***REMOVED***,

	// blank lines
	***REMOVED***"\nkey=value\n", "key", "value"***REMOVED***,
	***REMOVED***"\rkey=value\r", "key", "value"***REMOVED***,
	***REMOVED***"\r\nkey=value\r\n", "key", "value"***REMOVED***,
	***REMOVED***"\nkey=value\n \nkey2=value2", "key", "value", "key2", "value2"***REMOVED***,
	***REMOVED***"\nkey=value\n\t\nkey2=value2", "key", "value", "key2", "value2"***REMOVED***,

	// escaped chars in key
	***REMOVED***"k\\ ey = value", "k ey", "value"***REMOVED***,
	***REMOVED***"k\\:ey = value", "k:ey", "value"***REMOVED***,
	***REMOVED***"k\\=ey = value", "k=ey", "value"***REMOVED***,
	***REMOVED***"k\\fey = value", "k\fey", "value"***REMOVED***,
	***REMOVED***"k\\ney = value", "k\ney", "value"***REMOVED***,
	***REMOVED***"k\\rey = value", "k\rey", "value"***REMOVED***,
	***REMOVED***"k\\tey = value", "k\tey", "value"***REMOVED***,

	// escaped chars in value
	***REMOVED***"key = v\\ alue", "key", "v alue"***REMOVED***,
	***REMOVED***"key = v\\:alue", "key", "v:alue"***REMOVED***,
	***REMOVED***"key = v\\=alue", "key", "v=alue"***REMOVED***,
	***REMOVED***"key = v\\falue", "key", "v\falue"***REMOVED***,
	***REMOVED***"key = v\\nalue", "key", "v\nalue"***REMOVED***,
	***REMOVED***"key = v\\ralue", "key", "v\ralue"***REMOVED***,
	***REMOVED***"key = v\\talue", "key", "v\talue"***REMOVED***,

	// silently dropped escape character
	***REMOVED***"k\\zey = value", "kzey", "value"***REMOVED***,
	***REMOVED***"key = v\\zalue", "key", "vzalue"***REMOVED***,

	// unicode literals
	***REMOVED***"key\\u2318 = value", "key⌘", "value"***REMOVED***,
	***REMOVED***"k\\u2318ey = value", "k⌘ey", "value"***REMOVED***,
	***REMOVED***"key = value\\u2318", "key", "value⌘"***REMOVED***,
	***REMOVED***"key = valu\\u2318e", "key", "valu⌘e"***REMOVED***,

	// multiline values
	***REMOVED***"key = valueA,\\\n    valueB", "key", "valueA,valueB"***REMOVED***,   // SPACE indent
	***REMOVED***"key = valueA,\\\n\f\f\fvalueB", "key", "valueA,valueB"***REMOVED***, // FF indent
	***REMOVED***"key = valueA,\\\n\t\t\tvalueB", "key", "valueA,valueB"***REMOVED***, // TAB indent
	***REMOVED***"key = valueA,\\\n \f\tvalueB", "key", "valueA,valueB"***REMOVED***,  // mix indent

	// comments
	***REMOVED***"# this is a comment\n! and so is this\nkey1=value1\nkey#2=value#2\n\nkey!3=value!3\n# and another one\n! and the final one", "key1", "value1", "key#2", "value#2", "key!3", "value!3"***REMOVED***,

	// expansion tests
	***REMOVED***"key=value\nkey2=$***REMOVED***key***REMOVED***", "key", "value", "key2", "value"***REMOVED***,
	***REMOVED***"key=value\nkey2=aa$***REMOVED***key***REMOVED***", "key", "value", "key2", "aavalue"***REMOVED***,
	***REMOVED***"key=value\nkey2=$***REMOVED***key***REMOVED***bb", "key", "value", "key2", "valuebb"***REMOVED***,
	***REMOVED***"key=value\nkey2=aa$***REMOVED***key***REMOVED***bb", "key", "value", "key2", "aavaluebb"***REMOVED***,
	***REMOVED***"key=value\nkey2=$***REMOVED***key***REMOVED***\nkey3=$***REMOVED***key2***REMOVED***", "key", "value", "key2", "value", "key3", "value"***REMOVED***,
	***REMOVED***"key=$***REMOVED***USER***REMOVED***", "key", os.Getenv("USER")***REMOVED***,
	***REMOVED***"key=$***REMOVED***USER***REMOVED***\nUSER=value", "key", "value", "USER", "value"***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var commentTests = []struct ***REMOVED***
	input, key, value string
	comments          []string
***REMOVED******REMOVED***
	***REMOVED***"key=value", "key", "value", nil***REMOVED***,
	***REMOVED***"#\nkey=value", "key", "value", []string***REMOVED***""***REMOVED******REMOVED***,
	***REMOVED***"#comment\nkey=value", "key", "value", []string***REMOVED***"comment"***REMOVED******REMOVED***,
	***REMOVED***"# comment\nkey=value", "key", "value", []string***REMOVED***"comment"***REMOVED******REMOVED***,
	***REMOVED***"#  comment\nkey=value", "key", "value", []string***REMOVED***"comment"***REMOVED******REMOVED***,
	***REMOVED***"# comment\n\nkey=value", "key", "value", []string***REMOVED***"comment"***REMOVED******REMOVED***,
	***REMOVED***"# comment1\n# comment2\nkey=value", "key", "value", []string***REMOVED***"comment1", "comment2"***REMOVED******REMOVED***,
	***REMOVED***"# comment1\n\n# comment2\n\nkey=value", "key", "value", []string***REMOVED***"comment1", "comment2"***REMOVED******REMOVED***,
	***REMOVED***"!comment\nkey=value", "key", "value", []string***REMOVED***"comment"***REMOVED******REMOVED***,
	***REMOVED***"! comment\nkey=value", "key", "value", []string***REMOVED***"comment"***REMOVED******REMOVED***,
	***REMOVED***"!  comment\nkey=value", "key", "value", []string***REMOVED***"comment"***REMOVED******REMOVED***,
	***REMOVED***"! comment\n\nkey=value", "key", "value", []string***REMOVED***"comment"***REMOVED******REMOVED***,
	***REMOVED***"! comment1\n! comment2\nkey=value", "key", "value", []string***REMOVED***"comment1", "comment2"***REMOVED******REMOVED***,
	***REMOVED***"! comment1\n\n! comment2\n\nkey=value", "key", "value", []string***REMOVED***"comment1", "comment2"***REMOVED******REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var errorTests = []struct ***REMOVED***
	input, msg string
***REMOVED******REMOVED***
	// unicode literals
	***REMOVED***"key\\u1 = value", "invalid unicode literal"***REMOVED***,
	***REMOVED***"key\\u12 = value", "invalid unicode literal"***REMOVED***,
	***REMOVED***"key\\u123 = value", "invalid unicode literal"***REMOVED***,
	***REMOVED***"key\\u123g = value", "invalid unicode literal"***REMOVED***,
	***REMOVED***"key\\u123", "invalid unicode literal"***REMOVED***,

	// circular references
	***REMOVED***"key=$***REMOVED***key***REMOVED***", "circular reference"***REMOVED***,
	***REMOVED***"key1=$***REMOVED***key2***REMOVED***\nkey2=$***REMOVED***key1***REMOVED***", "circular reference"***REMOVED***,

	// malformed expressions
	***REMOVED***"key=$***REMOVED***ke", "malformed expression"***REMOVED***,
	***REMOVED***"key=valu$***REMOVED***ke", "malformed expression"***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var writeTests = []struct ***REMOVED***
	input, output, encoding string
***REMOVED******REMOVED***
	// ISO-8859-1 tests
	***REMOVED***"key = value", "key = value\n", "ISO-8859-1"***REMOVED***,
	***REMOVED***"key = value \\\n   continued", "key = value continued\n", "ISO-8859-1"***REMOVED***,
	***REMOVED***"key⌘ = value", "key\\u2318 = value\n", "ISO-8859-1"***REMOVED***,
	***REMOVED***"ke\\ \\:y = value", "ke\\ \\:y = value\n", "ISO-8859-1"***REMOVED***,

	// UTF-8 tests
	***REMOVED***"key = value", "key = value\n", "UTF-8"***REMOVED***,
	***REMOVED***"key = value \\\n   continued", "key = value continued\n", "UTF-8"***REMOVED***,
	***REMOVED***"key⌘ = value⌘", "key⌘ = value⌘\n", "UTF-8"***REMOVED***,
	***REMOVED***"ke\\ \\:y = value", "ke\\ \\:y = value\n", "UTF-8"***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var writeCommentTests = []struct ***REMOVED***
	input, output, encoding string
***REMOVED******REMOVED***
	// ISO-8859-1 tests
	***REMOVED***"key = value", "key = value\n", "ISO-8859-1"***REMOVED***,
	***REMOVED***"#\nkey = value", "key = value\n", "ISO-8859-1"***REMOVED***,
	***REMOVED***"#\n#\n#\nkey = value", "key = value\n", "ISO-8859-1"***REMOVED***,
	***REMOVED***"# comment\nkey = value", "# comment\nkey = value\n", "ISO-8859-1"***REMOVED***,
	***REMOVED***"\n# comment\nkey = value", "# comment\nkey = value\n", "ISO-8859-1"***REMOVED***,
	***REMOVED***"# comment\n\nkey = value", "# comment\nkey = value\n", "ISO-8859-1"***REMOVED***,
	***REMOVED***"# comment1\n# comment2\nkey = value", "# comment1\n# comment2\nkey = value\n", "ISO-8859-1"***REMOVED***,
	***REMOVED***"#comment1\nkey1 = value1\n#comment2\nkey2 = value2", "# comment1\nkey1 = value1\n\n# comment2\nkey2 = value2\n", "ISO-8859-1"***REMOVED***,

	// UTF-8 tests
	***REMOVED***"key = value", "key = value\n", "UTF-8"***REMOVED***,
	***REMOVED***"# comment⌘\nkey = value⌘", "# comment⌘\nkey = value⌘\n", "UTF-8"***REMOVED***,
	***REMOVED***"\n# comment⌘\nkey = value⌘", "# comment⌘\nkey = value⌘\n", "UTF-8"***REMOVED***,
	***REMOVED***"# comment⌘\n\nkey = value⌘", "# comment⌘\nkey = value⌘\n", "UTF-8"***REMOVED***,
	***REMOVED***"# comment1⌘\n# comment2⌘\nkey = value⌘", "# comment1⌘\n# comment2⌘\nkey = value⌘\n", "UTF-8"***REMOVED***,
	***REMOVED***"#comment1⌘\nkey1 = value1⌘\n#comment2⌘\nkey2 = value2⌘", "# comment1⌘\nkey1 = value1⌘\n\n# comment2⌘\nkey2 = value2⌘\n", "UTF-8"***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var boolTests = []struct ***REMOVED***
	input, key string
	def, value bool
***REMOVED******REMOVED***
	// valid values for TRUE
	***REMOVED***"key = 1", "key", false, true***REMOVED***,
	***REMOVED***"key = on", "key", false, true***REMOVED***,
	***REMOVED***"key = On", "key", false, true***REMOVED***,
	***REMOVED***"key = ON", "key", false, true***REMOVED***,
	***REMOVED***"key = true", "key", false, true***REMOVED***,
	***REMOVED***"key = True", "key", false, true***REMOVED***,
	***REMOVED***"key = TRUE", "key", false, true***REMOVED***,
	***REMOVED***"key = yes", "key", false, true***REMOVED***,
	***REMOVED***"key = Yes", "key", false, true***REMOVED***,
	***REMOVED***"key = YES", "key", false, true***REMOVED***,

	// valid values for FALSE (all other)
	***REMOVED***"key = 0", "key", true, false***REMOVED***,
	***REMOVED***"key = off", "key", true, false***REMOVED***,
	***REMOVED***"key = false", "key", true, false***REMOVED***,
	***REMOVED***"key = no", "key", true, false***REMOVED***,

	// non existent key
	***REMOVED***"key = true", "key2", false, false***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var durationTests = []struct ***REMOVED***
	input, key string
	def, value time.Duration
***REMOVED******REMOVED***
	// valid values
	***REMOVED***"key = 1", "key", 999, 1***REMOVED***,
	***REMOVED***"key = 0", "key", 999, 0***REMOVED***,
	***REMOVED***"key = -1", "key", 999, -1***REMOVED***,
	***REMOVED***"key = 0123", "key", 999, 123***REMOVED***,

	// invalid values
	***REMOVED***"key = 0xff", "key", 999, 999***REMOVED***,
	***REMOVED***"key = 1.0", "key", 999, 999***REMOVED***,
	***REMOVED***"key = a", "key", 999, 999***REMOVED***,

	// non existent key
	***REMOVED***"key = 1", "key2", 999, 999***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var parsedDurationTests = []struct ***REMOVED***
	input, key string
	def, value time.Duration
***REMOVED******REMOVED***
	// valid values
	***REMOVED***"key = -1ns", "key", 999, -1 * time.Nanosecond***REMOVED***,
	***REMOVED***"key = 300ms", "key", 999, 300 * time.Millisecond***REMOVED***,
	***REMOVED***"key = 5s", "key", 999, 5 * time.Second***REMOVED***,
	***REMOVED***"key = 3h", "key", 999, 3 * time.Hour***REMOVED***,
	***REMOVED***"key = 2h45m", "key", 999, 2*time.Hour + 45*time.Minute***REMOVED***,

	// invalid values
	***REMOVED***"key = 0xff", "key", 999, 999***REMOVED***,
	***REMOVED***"key = 1.0", "key", 999, 999***REMOVED***,
	***REMOVED***"key = a", "key", 999, 999***REMOVED***,
	***REMOVED***"key = 1", "key", 999, 999***REMOVED***,
	***REMOVED***"key = 0", "key", 999, 0***REMOVED***,

	// non existent key
	***REMOVED***"key = 1", "key2", 999, 999***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var floatTests = []struct ***REMOVED***
	input, key string
	def, value float64
***REMOVED******REMOVED***
	// valid values
	***REMOVED***"key = 1.0", "key", 999, 1.0***REMOVED***,
	***REMOVED***"key = 0.0", "key", 999, 0.0***REMOVED***,
	***REMOVED***"key = -1.0", "key", 999, -1.0***REMOVED***,
	***REMOVED***"key = 1", "key", 999, 1***REMOVED***,
	***REMOVED***"key = 0", "key", 999, 0***REMOVED***,
	***REMOVED***"key = -1", "key", 999, -1***REMOVED***,
	***REMOVED***"key = 0123", "key", 999, 123***REMOVED***,

	// invalid values
	***REMOVED***"key = 0xff", "key", 999, 999***REMOVED***,
	***REMOVED***"key = a", "key", 999, 999***REMOVED***,

	// non existent key
	***REMOVED***"key = 1", "key2", 999, 999***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var int64Tests = []struct ***REMOVED***
	input, key string
	def, value int64
***REMOVED******REMOVED***
	// valid values
	***REMOVED***"key = 1", "key", 999, 1***REMOVED***,
	***REMOVED***"key = 0", "key", 999, 0***REMOVED***,
	***REMOVED***"key = -1", "key", 999, -1***REMOVED***,
	***REMOVED***"key = 0123", "key", 999, 123***REMOVED***,

	// invalid values
	***REMOVED***"key = 0xff", "key", 999, 999***REMOVED***,
	***REMOVED***"key = 1.0", "key", 999, 999***REMOVED***,
	***REMOVED***"key = a", "key", 999, 999***REMOVED***,

	// non existent key
	***REMOVED***"key = 1", "key2", 999, 999***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var uint64Tests = []struct ***REMOVED***
	input, key string
	def, value uint64
***REMOVED******REMOVED***
	// valid values
	***REMOVED***"key = 1", "key", 999, 1***REMOVED***,
	***REMOVED***"key = 0", "key", 999, 0***REMOVED***,
	***REMOVED***"key = 0123", "key", 999, 123***REMOVED***,

	// invalid values
	***REMOVED***"key = -1", "key", 999, 999***REMOVED***,
	***REMOVED***"key = 0xff", "key", 999, 999***REMOVED***,
	***REMOVED***"key = 1.0", "key", 999, 999***REMOVED***,
	***REMOVED***"key = a", "key", 999, 999***REMOVED***,

	// non existent key
	***REMOVED***"key = 1", "key2", 999, 999***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var stringTests = []struct ***REMOVED***
	input, key string
	def, value string
***REMOVED******REMOVED***
	// valid values
	***REMOVED***"key = abc", "key", "def", "abc"***REMOVED***,

	// non existent key
	***REMOVED***"key = abc", "key2", "def", "def"***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var keysTests = []struct ***REMOVED***
	input string
	keys  []string
***REMOVED******REMOVED***
	***REMOVED***"", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"key = abc", []string***REMOVED***"key"***REMOVED******REMOVED***,
	***REMOVED***"key = abc\nkey2=def", []string***REMOVED***"key", "key2"***REMOVED******REMOVED***,
	***REMOVED***"key2 = abc\nkey=def", []string***REMOVED***"key2", "key"***REMOVED******REMOVED***,
	***REMOVED***"key = abc\nkey=def", []string***REMOVED***"key"***REMOVED******REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var filterTests = []struct ***REMOVED***
	input   string
	pattern string
	keys    []string
	err     string
***REMOVED******REMOVED***
	***REMOVED***"", "", []string***REMOVED******REMOVED***, ""***REMOVED***,
	***REMOVED***"", "abc", []string***REMOVED******REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value", "", []string***REMOVED***"key"***REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value", "key=", []string***REMOVED******REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "", []string***REMOVED***"foo", "key"***REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "f", []string***REMOVED***"foo"***REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "fo", []string***REMOVED***"foo"***REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "foo", []string***REMOVED***"foo"***REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "fooo", []string***REMOVED******REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value\nkey2=value2\nfoo=bar", "ey", []string***REMOVED***"key", "key2"***REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value\nkey2=value2\nfoo=bar", "key", []string***REMOVED***"key", "key2"***REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value\nkey2=value2\nfoo=bar", "^key", []string***REMOVED***"key", "key2"***REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value\nkey2=value2\nfoo=bar", "^(key|foo)", []string***REMOVED***"foo", "key", "key2"***REMOVED***, ""***REMOVED***,
	***REMOVED***"key=value\nkey2=value2\nfoo=bar", "[ abc", nil, "error parsing regexp.*"***REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var filterPrefixTests = []struct ***REMOVED***
	input  string
	prefix string
	keys   []string
***REMOVED******REMOVED***
	***REMOVED***"", "", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"", "abc", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"key=value", "", []string***REMOVED***"key"***REMOVED******REMOVED***,
	***REMOVED***"key=value", "key=", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "", []string***REMOVED***"foo", "key"***REMOVED******REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "f", []string***REMOVED***"foo"***REMOVED******REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "fo", []string***REMOVED***"foo"***REMOVED******REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "foo", []string***REMOVED***"foo"***REMOVED******REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "fooo", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"key=value\nkey2=value2\nfoo=bar", "key", []string***REMOVED***"key", "key2"***REMOVED******REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var filterStripPrefixTests = []struct ***REMOVED***
	input  string
	prefix string
	keys   []string
***REMOVED******REMOVED***
	***REMOVED***"", "", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"", "abc", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"key=value", "", []string***REMOVED***"key"***REMOVED******REMOVED***,
	***REMOVED***"key=value", "key=", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "", []string***REMOVED***"foo", "key"***REMOVED******REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "f", []string***REMOVED***"foo"***REMOVED******REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "fo", []string***REMOVED***"foo"***REMOVED******REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "foo", []string***REMOVED***"foo"***REMOVED******REMOVED***,
	***REMOVED***"key=value\nfoo=bar", "fooo", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"key=value\nkey2=value2\nfoo=bar", "key", []string***REMOVED***"key", "key2"***REMOVED******REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

var setTests = []struct ***REMOVED***
	input      string
	key, value string
	prev       string
	ok         bool
	err        string
	keys       []string
***REMOVED******REMOVED***
	***REMOVED***"", "", "", "", false, "", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"", "key", "value", "", false, "", []string***REMOVED***"key"***REMOVED******REMOVED***,
	***REMOVED***"key=value", "key2", "value2", "", false, "", []string***REMOVED***"key", "key2"***REMOVED******REMOVED***,
	***REMOVED***"key=value", "abc", "value3", "", false, "", []string***REMOVED***"key", "abc"***REMOVED******REMOVED***,
	***REMOVED***"key=value", "key", "value3", "value", true, "", []string***REMOVED***"key"***REMOVED******REMOVED***,
***REMOVED***

// ----------------------------------------------------------------------------

// TestBasic tests basic single key/value combinations with all possible
// whitespace, delimiter and newline permutations.
func TestBasic(t *testing.T) ***REMOVED***
	testWhitespaceAndDelimiterCombinations(t, "key", "")
	testWhitespaceAndDelimiterCombinations(t, "key", "value")
	testWhitespaceAndDelimiterCombinations(t, "key", "value   ")
***REMOVED***

func TestComplex(t *testing.T) ***REMOVED***
	for _, test := range complexTests ***REMOVED***
		testKeyValue(t, test[0], test[1:]...)
	***REMOVED***
***REMOVED***

func TestErrors(t *testing.T) ***REMOVED***
	for _, test := range errorTests ***REMOVED***
		_, err := Load([]byte(test.input), ISO_8859_1)
		assert.Equal(t, err != nil, true, "want error")
		assert.Equal(t, strings.Contains(err.Error(), test.msg), true)
	***REMOVED***
***REMOVED***

func TestDisableExpansion(t *testing.T) ***REMOVED***
	input := "key=value\nkey2=$***REMOVED***key***REMOVED***"
	p := mustParse(t, input)
	p.DisableExpansion = true
	assert.Equal(t, p.MustGet("key"), "value")
	assert.Equal(t, p.MustGet("key2"), "$***REMOVED***key***REMOVED***")

	// with expansion disabled we can introduce circular references
	p.MustSet("keyA", "$***REMOVED***keyB***REMOVED***")
	p.MustSet("keyB", "$***REMOVED***keyA***REMOVED***")
	assert.Equal(t, p.MustGet("keyA"), "$***REMOVED***keyB***REMOVED***")
	assert.Equal(t, p.MustGet("keyB"), "$***REMOVED***keyA***REMOVED***")
***REMOVED***

func TestDisableExpansionStillUpdatesKeys(t *testing.T) ***REMOVED***
	p := NewProperties()
	p.MustSet("p1", "a")
	assert.Equal(t, p.Keys(), []string***REMOVED***"p1"***REMOVED***)
	assert.Equal(t, p.String(), "p1 = a\n")

	p.DisableExpansion = true
	p.MustSet("p2", "b")

	assert.Equal(t, p.Keys(), []string***REMOVED***"p1", "p2"***REMOVED***)
	assert.Equal(t, p.String(), "p1 = a\np2 = b\n")
***REMOVED***

func TestMustGet(t *testing.T) ***REMOVED***
	input := "key = value\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGet("key"), "value")
	assert.Panic(t, func() ***REMOVED*** p.MustGet("invalid") ***REMOVED***, "unknown property: invalid")
***REMOVED***

func TestGetBool(t *testing.T) ***REMOVED***
	for _, test := range boolTests ***REMOVED***
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetBool(test.key, test.def), test.value)
	***REMOVED***
***REMOVED***

func TestMustGetBool(t *testing.T) ***REMOVED***
	input := "key = true\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetBool("key"), true)
	assert.Panic(t, func() ***REMOVED*** p.MustGetBool("invalid") ***REMOVED***, "unknown property: invalid")
***REMOVED***

func TestGetDuration(t *testing.T) ***REMOVED***
	for _, test := range durationTests ***REMOVED***
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetDuration(test.key, test.def), test.value)
	***REMOVED***
***REMOVED***

func TestMustGetDuration(t *testing.T) ***REMOVED***
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetDuration("key"), time.Duration(123))
	assert.Panic(t, func() ***REMOVED*** p.MustGetDuration("key2") ***REMOVED***, "strconv.ParseInt: parsing.*")
	assert.Panic(t, func() ***REMOVED*** p.MustGetDuration("invalid") ***REMOVED***, "unknown property: invalid")
***REMOVED***

func TestGetParsedDuration(t *testing.T) ***REMOVED***
	for _, test := range parsedDurationTests ***REMOVED***
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetParsedDuration(test.key, test.def), test.value)
	***REMOVED***
***REMOVED***

func TestMustGetParsedDuration(t *testing.T) ***REMOVED***
	input := "key = 123ms\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetParsedDuration("key"), 123*time.Millisecond)
	assert.Panic(t, func() ***REMOVED*** p.MustGetParsedDuration("key2") ***REMOVED***, "time: invalid duration ghi")
	assert.Panic(t, func() ***REMOVED*** p.MustGetParsedDuration("invalid") ***REMOVED***, "unknown property: invalid")
***REMOVED***

func TestGetFloat64(t *testing.T) ***REMOVED***
	for _, test := range floatTests ***REMOVED***
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetFloat64(test.key, test.def), test.value)
	***REMOVED***
***REMOVED***

func TestMustGetFloat64(t *testing.T) ***REMOVED***
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetFloat64("key"), float64(123))
	assert.Panic(t, func() ***REMOVED*** p.MustGetFloat64("key2") ***REMOVED***, "strconv.ParseFloat: parsing.*")
	assert.Panic(t, func() ***REMOVED*** p.MustGetFloat64("invalid") ***REMOVED***, "unknown property: invalid")
***REMOVED***

func TestGetInt(t *testing.T) ***REMOVED***
	for _, test := range int64Tests ***REMOVED***
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetInt(test.key, int(test.def)), int(test.value))
	***REMOVED***
***REMOVED***

func TestMustGetInt(t *testing.T) ***REMOVED***
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetInt("key"), int(123))
	assert.Panic(t, func() ***REMOVED*** p.MustGetInt("key2") ***REMOVED***, "strconv.ParseInt: parsing.*")
	assert.Panic(t, func() ***REMOVED*** p.MustGetInt("invalid") ***REMOVED***, "unknown property: invalid")
***REMOVED***

func TestGetInt64(t *testing.T) ***REMOVED***
	for _, test := range int64Tests ***REMOVED***
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetInt64(test.key, test.def), test.value)
	***REMOVED***
***REMOVED***

func TestMustGetInt64(t *testing.T) ***REMOVED***
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetInt64("key"), int64(123))
	assert.Panic(t, func() ***REMOVED*** p.MustGetInt64("key2") ***REMOVED***, "strconv.ParseInt: parsing.*")
	assert.Panic(t, func() ***REMOVED*** p.MustGetInt64("invalid") ***REMOVED***, "unknown property: invalid")
***REMOVED***

func TestGetUint(t *testing.T) ***REMOVED***
	for _, test := range uint64Tests ***REMOVED***
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetUint(test.key, uint(test.def)), uint(test.value))
	***REMOVED***
***REMOVED***

func TestMustGetUint(t *testing.T) ***REMOVED***
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetUint("key"), uint(123))
	assert.Panic(t, func() ***REMOVED*** p.MustGetUint64("key2") ***REMOVED***, "strconv.ParseUint: parsing.*")
	assert.Panic(t, func() ***REMOVED*** p.MustGetUint64("invalid") ***REMOVED***, "unknown property: invalid")
***REMOVED***

func TestGetUint64(t *testing.T) ***REMOVED***
	for _, test := range uint64Tests ***REMOVED***
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetUint64(test.key, test.def), test.value)
	***REMOVED***
***REMOVED***

func TestMustGetUint64(t *testing.T) ***REMOVED***
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetUint64("key"), uint64(123))
	assert.Panic(t, func() ***REMOVED*** p.MustGetUint64("key2") ***REMOVED***, "strconv.ParseUint: parsing.*")
	assert.Panic(t, func() ***REMOVED*** p.MustGetUint64("invalid") ***REMOVED***, "unknown property: invalid")
***REMOVED***

func TestGetString(t *testing.T) ***REMOVED***
	for _, test := range stringTests ***REMOVED***
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetString(test.key, test.def), test.value)
	***REMOVED***
***REMOVED***

func TestMustGetString(t *testing.T) ***REMOVED***
	input := `key = value`
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetString("key"), "value")
	assert.Panic(t, func() ***REMOVED*** p.MustGetString("invalid") ***REMOVED***, "unknown property: invalid")
***REMOVED***

func TestComment(t *testing.T) ***REMOVED***
	for _, test := range commentTests ***REMOVED***
		p := mustParse(t, test.input)
		assert.Equal(t, p.MustGetString(test.key), test.value)
		assert.Equal(t, p.GetComments(test.key), test.comments)
		if test.comments != nil ***REMOVED***
			assert.Equal(t, p.GetComment(test.key), test.comments[len(test.comments)-1])
		***REMOVED*** else ***REMOVED***
			assert.Equal(t, p.GetComment(test.key), "")
		***REMOVED***

		// test setting comments
		if len(test.comments) > 0 ***REMOVED***
			// set single comment
			p.ClearComments()
			assert.Equal(t, len(p.c), 0)
			p.SetComment(test.key, test.comments[0])
			assert.Equal(t, p.GetComment(test.key), test.comments[0])

			// set multiple comments
			p.ClearComments()
			assert.Equal(t, len(p.c), 0)
			p.SetComments(test.key, test.comments)
			assert.Equal(t, p.GetComments(test.key), test.comments)

			// clear comments for a key
			p.SetComments(test.key, nil)
			assert.Equal(t, p.GetComment(test.key), "")
			assert.Equal(t, p.GetComments(test.key), ([]string)(nil))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFilter(t *testing.T) ***REMOVED***
	for _, test := range filterTests ***REMOVED***
		p := mustParse(t, test.input)
		pp, err := p.Filter(test.pattern)
		if err != nil ***REMOVED***
			assert.Matches(t, err.Error(), test.err)
			continue
		***REMOVED***
		assert.Equal(t, pp != nil, true, "want properties")
		assert.Equal(t, pp.Len(), len(test.keys))
		for _, key := range test.keys ***REMOVED***
			v1, ok1 := p.Get(key)
			v2, ok2 := pp.Get(key)
			assert.Equal(t, ok1, true)
			assert.Equal(t, ok2, true)
			assert.Equal(t, v1, v2)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFilterPrefix(t *testing.T) ***REMOVED***
	for _, test := range filterPrefixTests ***REMOVED***
		p := mustParse(t, test.input)
		pp := p.FilterPrefix(test.prefix)
		assert.Equal(t, pp != nil, true, "want properties")
		assert.Equal(t, pp.Len(), len(test.keys))
		for _, key := range test.keys ***REMOVED***
			v1, ok1 := p.Get(key)
			v2, ok2 := pp.Get(key)
			assert.Equal(t, ok1, true)
			assert.Equal(t, ok2, true)
			assert.Equal(t, v1, v2)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFilterStripPrefix(t *testing.T) ***REMOVED***
	for _, test := range filterStripPrefixTests ***REMOVED***
		p := mustParse(t, test.input)
		pp := p.FilterPrefix(test.prefix)
		assert.Equal(t, pp != nil, true, "want properties")
		assert.Equal(t, pp.Len(), len(test.keys))
		for _, key := range test.keys ***REMOVED***
			v1, ok1 := p.Get(key)
			v2, ok2 := pp.Get(key)
			assert.Equal(t, ok1, true)
			assert.Equal(t, ok2, true)
			assert.Equal(t, v1, v2)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestKeys(t *testing.T) ***REMOVED***
	for _, test := range keysTests ***REMOVED***
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), len(test.keys))
		assert.Equal(t, len(p.Keys()), len(test.keys))
		assert.Equal(t, p.Keys(), test.keys)
	***REMOVED***
***REMOVED***

func TestSet(t *testing.T) ***REMOVED***
	for _, test := range setTests ***REMOVED***
		p := mustParse(t, test.input)
		prev, ok, err := p.Set(test.key, test.value)
		if test.err != "" ***REMOVED***
			assert.Matches(t, err.Error(), test.err)
			continue
		***REMOVED***

		assert.Equal(t, err, nil)
		assert.Equal(t, ok, test.ok)
		if ok ***REMOVED***
			assert.Equal(t, prev, test.prev)
		***REMOVED***
		assert.Equal(t, p.Keys(), test.keys)
	***REMOVED***
***REMOVED***

func TestSetValue(t *testing.T) ***REMOVED***
	tests := []interface***REMOVED******REMOVED******REMOVED***
		true, false,
		int8(123), int16(123), int32(123), int64(123), int(123),
		uint8(123), uint16(123), uint32(123), uint64(123), uint(123),
		float32(1.23), float64(1.23),
		"abc",
	***REMOVED***

	for _, v := range tests ***REMOVED***
		p := NewProperties()
		err := p.SetValue("x", v)
		assert.Equal(t, err, nil)
		assert.Equal(t, p.GetString("x", ""), fmt.Sprintf("%v", v))
	***REMOVED***
***REMOVED***

func TestMustSet(t *testing.T) ***REMOVED***
	input := "key=$***REMOVED***key***REMOVED***"
	p := mustParse(t, input)
	assert.Panic(t, func() ***REMOVED*** p.MustSet("key", "$***REMOVED***key***REMOVED***") ***REMOVED***, "circular reference .*")
***REMOVED***

func TestWrite(t *testing.T) ***REMOVED***
	for _, test := range writeTests ***REMOVED***
		p, err := parse(test.input)

		buf := new(bytes.Buffer)
		var n int
		switch test.encoding ***REMOVED***
		case "UTF-8":
			n, err = p.Write(buf, UTF8)
		case "ISO-8859-1":
			n, err = p.Write(buf, ISO_8859_1)
		***REMOVED***
		assert.Equal(t, err, nil)
		s := string(buf.Bytes())
		assert.Equal(t, n, len(test.output), fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
		assert.Equal(t, s, test.output, fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
	***REMOVED***
***REMOVED***

func TestWriteComment(t *testing.T) ***REMOVED***
	for _, test := range writeCommentTests ***REMOVED***
		p, err := parse(test.input)

		buf := new(bytes.Buffer)
		var n int
		switch test.encoding ***REMOVED***
		case "UTF-8":
			n, err = p.WriteComment(buf, "# ", UTF8)
		case "ISO-8859-1":
			n, err = p.WriteComment(buf, "# ", ISO_8859_1)
		***REMOVED***
		assert.Equal(t, err, nil)
		s := string(buf.Bytes())
		assert.Equal(t, n, len(test.output), fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
		assert.Equal(t, s, test.output, fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
	***REMOVED***
***REMOVED***

func TestCustomExpansionExpression(t *testing.T) ***REMOVED***
	testKeyValuePrePostfix(t, "*[", "]*", "key=value\nkey2=*[key]*", "key", "value", "key2", "value")
***REMOVED***

func TestPanicOn32BitIntOverflow(t *testing.T) ***REMOVED***
	is32Bit = true
	var min, max int64 = math.MinInt32 - 1, math.MaxInt32 + 1
	input := fmt.Sprintf("min=%d\nmax=%d", min, max)
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetInt64("min"), min)
	assert.Equal(t, p.MustGetInt64("max"), max)
	assert.Panic(t, func() ***REMOVED*** p.MustGetInt("min") ***REMOVED***, ".* out of range")
	assert.Panic(t, func() ***REMOVED*** p.MustGetInt("max") ***REMOVED***, ".* out of range")
***REMOVED***

func TestPanicOn32BitUintOverflow(t *testing.T) ***REMOVED***
	is32Bit = true
	var max uint64 = math.MaxUint32 + 1
	input := fmt.Sprintf("max=%d", max)
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetUint64("max"), max)
	assert.Panic(t, func() ***REMOVED*** p.MustGetUint("max") ***REMOVED***, ".* out of range")
***REMOVED***

func TestDeleteKey(t *testing.T) ***REMOVED***
	input := "#comments should also be gone\nkey=to-be-deleted\nsecond=key"
	p := mustParse(t, input)
	assert.Equal(t, len(p.m), 2)
	assert.Equal(t, len(p.c), 1)
	assert.Equal(t, len(p.k), 2)
	p.Delete("key")
	assert.Equal(t, len(p.m), 1)
	assert.Equal(t, len(p.c), 0)
	assert.Equal(t, len(p.k), 1)
	assert.Equal(t, p.k[0], "second")
	assert.Equal(t, p.m["second"], "key")
***REMOVED***

func TestDeleteUnknownKey(t *testing.T) ***REMOVED***
	input := "#comments should also be gone\nkey=to-be-deleted"
	p := mustParse(t, input)
	assert.Equal(t, len(p.m), 1)
	assert.Equal(t, len(p.c), 1)
	assert.Equal(t, len(p.k), 1)
	p.Delete("wrong-key")
	assert.Equal(t, len(p.m), 1)
	assert.Equal(t, len(p.c), 1)
	assert.Equal(t, len(p.k), 1)
***REMOVED***

func TestMerge(t *testing.T) ***REMOVED***
	input1 := "#comment\nkey=value\nkey2=value2"
	input2 := "#another comment\nkey=another value\nkey3=value3"
	p1 := mustParse(t, input1)
	p2 := mustParse(t, input2)
	p1.Merge(p2)
	assert.Equal(t, len(p1.m), 3)
	assert.Equal(t, len(p1.c), 1)
	assert.Equal(t, len(p1.k), 3)
	assert.Equal(t, p1.MustGet("key"), "another value")
	assert.Equal(t, p1.GetComment("key"), "another comment")
***REMOVED***

func TestMap(t *testing.T) ***REMOVED***
	input := "key=value\nabc=def"
	p := mustParse(t, input)
	m := map[string]string***REMOVED***"key": "value", "abc": "def"***REMOVED***
	assert.Equal(t, p.Map(), m)
***REMOVED***

func TestFilterFunc(t *testing.T) ***REMOVED***
	input := "key=value\nabc=def"
	p := mustParse(t, input)
	pp := p.FilterFunc(func(k, v string) bool ***REMOVED***
		return k != "abc"
	***REMOVED***)
	m := map[string]string***REMOVED***"key": "value"***REMOVED***
	assert.Equal(t, pp.Map(), m)
***REMOVED***

// ----------------------------------------------------------------------------

// tests all combinations of delimiters, leading and/or trailing whitespace and newlines.
func testWhitespaceAndDelimiterCombinations(t *testing.T, key, value string) ***REMOVED***
	whitespace := []string***REMOVED***"", " ", "\f", "\t"***REMOVED***
	delimiters := []string***REMOVED***"", " ", "=", ":"***REMOVED***
	newlines := []string***REMOVED***"", "\r", "\n", "\r\n"***REMOVED***
	for _, dl := range delimiters ***REMOVED***
		for _, ws1 := range whitespace ***REMOVED***
			for _, ws2 := range whitespace ***REMOVED***
				for _, nl := range newlines ***REMOVED***
					// skip the one case where there is nothing between a key and a value
					if ws1 == "" && dl == "" && ws2 == "" && value != "" ***REMOVED***
						continue
					***REMOVED***

					input := fmt.Sprintf("%s%s%s%s%s%s", key, ws1, dl, ws2, value, nl)
					testKeyValue(t, input, key, value)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// tests whether key/value pairs exist for a given input.
// keyvalues is expected to be an even number of strings of "key", "value", ...
func testKeyValue(t *testing.T, input string, keyvalues ...string) ***REMOVED***
	testKeyValuePrePostfix(t, "$***REMOVED***", "***REMOVED***", input, keyvalues...)
***REMOVED***

// tests whether key/value pairs exist for a given input.
// keyvalues is expected to be an even number of strings of "key", "value", ...
func testKeyValuePrePostfix(t *testing.T, prefix, postfix, input string, keyvalues ...string) ***REMOVED***
	p, err := Load([]byte(input), ISO_8859_1)
	assert.Equal(t, err, nil)
	p.Prefix = prefix
	p.Postfix = postfix
	assertKeyValues(t, input, p, keyvalues...)
***REMOVED***

// tests whether key/value pairs exist for a given input.
// keyvalues is expected to be an even number of strings of "key", "value", ...
func assertKeyValues(t *testing.T, input string, p *Properties, keyvalues ...string) ***REMOVED***
	assert.Equal(t, p != nil, true, "want properties")
	assert.Equal(t, 2*p.Len(), len(keyvalues), "Odd number of key/value pairs.")

	for i := 0; i < len(keyvalues); i += 2 ***REMOVED***
		key, value := keyvalues[i], keyvalues[i+1]
		v, ok := p.Get(key)
		if !ok ***REMOVED***
			t.Errorf("No key %q found (input=%q)", key, input)
		***REMOVED***
		if got, want := v, value; !reflect.DeepEqual(got, want) ***REMOVED***
			t.Errorf("Value %q does not match %q (input=%q)", v, value, input)
		***REMOVED***
	***REMOVED***
***REMOVED***

func mustParse(t *testing.T, s string) *Properties ***REMOVED***
	p, err := parse(s)
	if err != nil ***REMOVED***
		t.Fatalf("parse failed with %s", err)
	***REMOVED***
	return p
***REMOVED***

// prints to stderr if the -verbose flag was given.
func printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if *verbose ***REMOVED***
		fmt.Fprintf(os.Stderr, format, args...)
	***REMOVED***
***REMOVED***
