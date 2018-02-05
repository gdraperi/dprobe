package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
)

func TestType(t *testing.T) ***REMOVED***
	var literals = []struct ***REMOVED***
		typ token.Type
		src string
	***REMOVED******REMOVED***
		***REMOVED***token.STRING, `"foo": "bar"`***REMOVED***,
		***REMOVED***token.NUMBER, `"foo": 123`***REMOVED***,
		***REMOVED***token.FLOAT, `"foo": 123.12`***REMOVED***,
		***REMOVED***token.FLOAT, `"foo": -123.12`***REMOVED***,
		***REMOVED***token.BOOL, `"foo": true`***REMOVED***,
		***REMOVED***token.STRING, `"foo": null`***REMOVED***,
	***REMOVED***

	for _, l := range literals ***REMOVED***
		t.Logf("Testing: %s", l.src)

		p := newParser([]byte(l.src))
		item, err := p.objectItem()
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***

		lit, ok := item.Val.(*ast.LiteralType)
		if !ok ***REMOVED***
			t.Errorf("node should be of type LiteralType, got: %T", item.Val)
		***REMOVED***

		if lit.Token.Type != l.typ ***REMOVED***
			t.Errorf("want: %s, got: %s", l.typ, lit.Token.Type)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestListType(t *testing.T) ***REMOVED***
	var literals = []struct ***REMOVED***
		src    string
		tokens []token.Type
	***REMOVED******REMOVED***
		***REMOVED***
			`"foo": ["123", 123]`,
			[]token.Type***REMOVED***token.STRING, token.NUMBER***REMOVED***,
		***REMOVED***,
		***REMOVED***
			`"foo": [123, "123",]`,
			[]token.Type***REMOVED***token.NUMBER, token.STRING***REMOVED***,
		***REMOVED***,
		***REMOVED***
			`"foo": []`,
			[]token.Type***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			`"foo": ["123", 123]`,
			[]token.Type***REMOVED***token.STRING, token.NUMBER***REMOVED***,
		***REMOVED***,
		***REMOVED***
			`"foo": ["123", ***REMOVED******REMOVED***]`,
			[]token.Type***REMOVED***token.STRING, token.LBRACE***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, l := range literals ***REMOVED***
		t.Logf("Testing: %s", l.src)

		p := newParser([]byte(l.src))
		item, err := p.objectItem()
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***

		list, ok := item.Val.(*ast.ListType)
		if !ok ***REMOVED***
			t.Errorf("node should be of type LiteralType, got: %T", item.Val)
		***REMOVED***

		tokens := []token.Type***REMOVED******REMOVED***
		for _, li := range list.List ***REMOVED***
			switch v := li.(type) ***REMOVED***
			case *ast.LiteralType:
				tokens = append(tokens, v.Token.Type)
			case *ast.ObjectType:
				tokens = append(tokens, token.LBRACE)
			***REMOVED***
		***REMOVED***

		equals(t, l.tokens, tokens)
	***REMOVED***
***REMOVED***

func TestObjectType(t *testing.T) ***REMOVED***
	var literals = []struct ***REMOVED***
		src      string
		nodeType []ast.Node
		itemLen  int
	***REMOVED******REMOVED***
		***REMOVED***
			`"foo": ***REMOVED******REMOVED***`,
			nil,
			0,
		***REMOVED***,
		***REMOVED***
			`"foo": ***REMOVED***
				"bar": "fatih"
			 ***REMOVED***`,
			[]ast.Node***REMOVED***&ast.LiteralType***REMOVED******REMOVED******REMOVED***,
			1,
		***REMOVED***,
		***REMOVED***
			`"foo": ***REMOVED***
				"bar": "fatih",
				"baz": ["arslan"]
			 ***REMOVED***`,
			[]ast.Node***REMOVED***
				&ast.LiteralType***REMOVED******REMOVED***,
				&ast.ListType***REMOVED******REMOVED***,
			***REMOVED***,
			2,
		***REMOVED***,
		***REMOVED***
			`"foo": ***REMOVED***
				"bar": ***REMOVED******REMOVED***
			 ***REMOVED***`,
			[]ast.Node***REMOVED***
				&ast.ObjectType***REMOVED******REMOVED***,
			***REMOVED***,
			1,
		***REMOVED***,
		***REMOVED***
			`"foo": ***REMOVED***
				"bar": ***REMOVED******REMOVED***,
				"foo": true
			 ***REMOVED***`,
			[]ast.Node***REMOVED***
				&ast.ObjectType***REMOVED******REMOVED***,
				&ast.LiteralType***REMOVED******REMOVED***,
			***REMOVED***,
			2,
		***REMOVED***,
	***REMOVED***

	for _, l := range literals ***REMOVED***
		t.Logf("Testing:\n%s\n", l.src)

		p := newParser([]byte(l.src))
		// p.enableTrace = true
		item, err := p.objectItem()
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***

		// we know that the ObjectKey name is foo for all cases, what matters
		// is the object
		obj, ok := item.Val.(*ast.ObjectType)
		if !ok ***REMOVED***
			t.Errorf("node should be of type LiteralType, got: %T", item.Val)
		***REMOVED***

		// check if the total length of items are correct
		equals(t, l.itemLen, len(obj.List.Items))

		// check if the types are correct
		for i, item := range obj.List.Items ***REMOVED***
			equals(t, reflect.TypeOf(l.nodeType[i]), reflect.TypeOf(item.Val))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFlattenObjects(t *testing.T) ***REMOVED***
	var literals = []struct ***REMOVED***
		src      string
		nodeType []ast.Node
		itemLen  int
	***REMOVED******REMOVED***
		***REMOVED***
			`***REMOVED***
					"foo": [
						***REMOVED***
							"foo": "svh",
							"bar": "fatih"
						***REMOVED***
					]
				***REMOVED***`,
			[]ast.Node***REMOVED***
				&ast.ObjectType***REMOVED******REMOVED***,
				&ast.LiteralType***REMOVED******REMOVED***,
				&ast.LiteralType***REMOVED******REMOVED***,
			***REMOVED***,
			3,
		***REMOVED***,
		***REMOVED***
			`***REMOVED***
					"variable": ***REMOVED***
						"foo": ***REMOVED******REMOVED***
					***REMOVED***
				***REMOVED***`,
			[]ast.Node***REMOVED***
				&ast.ObjectType***REMOVED******REMOVED***,
			***REMOVED***,
			1,
		***REMOVED***,
		***REMOVED***
			`***REMOVED***
				"empty": []
			***REMOVED***`,
			[]ast.Node***REMOVED***
				&ast.ListType***REMOVED******REMOVED***,
			***REMOVED***,
			1,
		***REMOVED***,
		***REMOVED***
			`***REMOVED***
				"basic": [1, 2, 3]
			***REMOVED***`,
			[]ast.Node***REMOVED***
				&ast.ListType***REMOVED******REMOVED***,
			***REMOVED***,
			1,
		***REMOVED***,
	***REMOVED***

	for _, l := range literals ***REMOVED***
		t.Logf("Testing:\n%s\n", l.src)

		f, err := Parse([]byte(l.src))
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***

		// the first object is always an ObjectList so just assert that one
		// so we can use it as such
		obj, ok := f.Node.(*ast.ObjectList)
		if !ok ***REMOVED***
			t.Errorf("node should be *ast.ObjectList, got: %T", f.Node)
		***REMOVED***

		// check if the types are correct
		var i int
		for _, item := range obj.Items ***REMOVED***
			equals(t, reflect.TypeOf(l.nodeType[i]), reflect.TypeOf(item.Val))
			i++

			if obj, ok := item.Val.(*ast.ObjectType); ok ***REMOVED***
				for _, item := range obj.List.Items ***REMOVED***
					equals(t, reflect.TypeOf(l.nodeType[i]), reflect.TypeOf(item.Val))
					i++
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// check if the number of items is correct
		equals(t, l.itemLen, i)

	***REMOVED***
***REMOVED***

func TestObjectKey(t *testing.T) ***REMOVED***
	keys := []struct ***REMOVED***
		exp []token.Type
		src string
	***REMOVED******REMOVED***
		***REMOVED***[]token.Type***REMOVED***token.STRING***REMOVED***, `"foo": ***REMOVED******REMOVED***`***REMOVED***,
	***REMOVED***

	for _, k := range keys ***REMOVED***
		p := newParser([]byte(k.src))
		keys, err := p.objectKey()
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		tokens := []token.Type***REMOVED******REMOVED***
		for _, o := range keys ***REMOVED***
			tokens = append(tokens, o.Token.Type)
		***REMOVED***

		equals(t, k.exp, tokens)
	***REMOVED***

	errKeys := []struct ***REMOVED***
		src string
	***REMOVED******REMOVED***
		***REMOVED***`foo 12 ***REMOVED******REMOVED***`***REMOVED***,
		***REMOVED***`foo bar = ***REMOVED******REMOVED***`***REMOVED***,
		***REMOVED***`foo []`***REMOVED***,
		***REMOVED***`12 ***REMOVED******REMOVED***`***REMOVED***,
	***REMOVED***

	for _, k := range errKeys ***REMOVED***
		p := newParser([]byte(k.src))
		_, err := p.objectKey()
		if err == nil ***REMOVED***
			t.Errorf("case '%s' should give an error", k.src)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Official HCL tests
func TestParse(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		Name string
		Err  bool
	***REMOVED******REMOVED***
		***REMOVED***
			"array.json",
			false,
		***REMOVED***,
		***REMOVED***
			"basic.json",
			false,
		***REMOVED***,
		***REMOVED***
			"object.json",
			false,
		***REMOVED***,
		***REMOVED***
			"types.json",
			false,
		***REMOVED***,
		***REMOVED***
			"bad_input_128.json",
			true,
		***REMOVED***,
		***REMOVED***
			"bad_input_tf_8110.json",
			true,
		***REMOVED***,
		***REMOVED***
			"good_input_tf_8110.json",
			false,
		***REMOVED***,
	***REMOVED***

	const fixtureDir = "./test-fixtures"

	for _, tc := range cases ***REMOVED***
		d, err := ioutil.ReadFile(filepath.Join(fixtureDir, tc.Name))
		if err != nil ***REMOVED***
			t.Fatalf("err: %s", err)
		***REMOVED***

		_, err = Parse(d)
		if (err != nil) != tc.Err ***REMOVED***
			t.Fatalf("Input: %s\n\nError: %s", tc.Name, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParse_inline(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		Value string
		Err   bool
	***REMOVED******REMOVED***
		***REMOVED***"***REMOVED***:***REMOVED***", true***REMOVED***,
	***REMOVED***

	for _, tc := range cases ***REMOVED***
		_, err := Parse([]byte(tc.Value))
		if (err != nil) != tc.Err ***REMOVED***
			t.Fatalf("Input: %q\n\nError: %s", tc.Value, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface***REMOVED******REMOVED***) ***REMOVED***
	if !reflect.DeepEqual(exp, act) ***REMOVED***
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %s\n\n\tgot: %s\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	***REMOVED***
***REMOVED***
