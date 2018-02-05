package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
)

func TestType(t *testing.T) ***REMOVED***
	var literals = []struct ***REMOVED***
		typ token.Type
		src string
	***REMOVED******REMOVED***
		***REMOVED***token.STRING, `foo = "foo"`***REMOVED***,
		***REMOVED***token.NUMBER, `foo = 123`***REMOVED***,
		***REMOVED***token.NUMBER, `foo = -29`***REMOVED***,
		***REMOVED***token.FLOAT, `foo = 123.12`***REMOVED***,
		***REMOVED***token.FLOAT, `foo = -123.12`***REMOVED***,
		***REMOVED***token.BOOL, `foo = true`***REMOVED***,
		***REMOVED***token.HEREDOC, "foo = <<EOF\nHello\nWorld\nEOF"***REMOVED***,
	***REMOVED***

	for _, l := range literals ***REMOVED***
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
			`foo = ["123", 123]`,
			[]token.Type***REMOVED***token.STRING, token.NUMBER***REMOVED***,
		***REMOVED***,
		***REMOVED***
			`foo = [123, "123",]`,
			[]token.Type***REMOVED***token.NUMBER, token.STRING***REMOVED***,
		***REMOVED***,
		***REMOVED***
			`foo = [false]`,
			[]token.Type***REMOVED***token.BOOL***REMOVED***,
		***REMOVED***,
		***REMOVED***
			`foo = []`,
			[]token.Type***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			`foo = [1,
"string",
<<EOF
heredoc contents
EOF
]`,
			[]token.Type***REMOVED***token.NUMBER, token.STRING, token.HEREDOC***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, l := range literals ***REMOVED***
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
			if tp, ok := li.(*ast.LiteralType); ok ***REMOVED***
				tokens = append(tokens, tp.Token.Type)
			***REMOVED***
		***REMOVED***

		equals(t, l.tokens, tokens)
	***REMOVED***
***REMOVED***

func TestListOfMaps(t *testing.T) ***REMOVED***
	src := `foo = [
    ***REMOVED***key = "bar"***REMOVED***,
    ***REMOVED***key = "baz", key2 = "qux"***REMOVED***,
  ]`
	p := newParser([]byte(src))

	file, err := p.Parse()
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	// Here we make all sorts of assumptions about the input structure w/ type
	// assertions. The intent is only for this to be a "smoke test" ensuring
	// parsing actually performed its duty - giving this test something a bit
	// more robust than _just_ "no error occurred".
	expected := []string***REMOVED***`"bar"`, `"baz"`, `"qux"`***REMOVED***
	actual := make([]string, 0, 3)
	ol := file.Node.(*ast.ObjectList)
	objItem := ol.Items[0]
	list := objItem.Val.(*ast.ListType)
	for _, node := range list.List ***REMOVED***
		obj := node.(*ast.ObjectType)
		for _, item := range obj.List.Items ***REMOVED***
			val := item.Val.(*ast.LiteralType)
			actual = append(actual, val.Token.Text)
		***REMOVED***

	***REMOVED***
	if !reflect.DeepEqual(expected, actual) ***REMOVED***
		t.Fatalf("Expected: %#v, got %#v", expected, actual)
	***REMOVED***
***REMOVED***

func TestListOfMaps_requiresComma(t *testing.T) ***REMOVED***
	src := `foo = [
    ***REMOVED***key = "bar"***REMOVED***
    ***REMOVED***key = "baz"***REMOVED***
  ]`
	p := newParser([]byte(src))

	_, err := p.Parse()
	if err == nil ***REMOVED***
		t.Fatalf("Expected error, got none!")
	***REMOVED***

	expected := "error parsing list, expected comma or list end"
	if !strings.Contains(err.Error(), expected) ***REMOVED***
		t.Fatalf("Expected err:\n  %s\nTo contain:\n  %s\n", err, expected)
	***REMOVED***
***REMOVED***

func TestListType_leadComment(t *testing.T) ***REMOVED***
	var literals = []struct ***REMOVED***
		src     string
		comment []string
	***REMOVED******REMOVED***
		***REMOVED***
			`foo = [
			1,
			# bar
			2,
			3,
			]`,
			[]string***REMOVED***"", "# bar", ""***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, l := range literals ***REMOVED***
		p := newParser([]byte(l.src))
		item, err := p.objectItem()
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		list, ok := item.Val.(*ast.ListType)
		if !ok ***REMOVED***
			t.Fatalf("node should be of type LiteralType, got: %T", item.Val)
		***REMOVED***

		if len(list.List) != len(l.comment) ***REMOVED***
			t.Fatalf("bad: %d", len(list.List))
		***REMOVED***

		for i, li := range list.List ***REMOVED***
			lt := li.(*ast.LiteralType)
			comment := l.comment[i]

			if (lt.LeadComment == nil) != (comment == "") ***REMOVED***
				t.Fatalf("bad: %#v", lt)
			***REMOVED***

			if comment == "" ***REMOVED***
				continue
			***REMOVED***

			actual := lt.LeadComment.List[0].Text
			if actual != comment ***REMOVED***
				t.Fatalf("bad: %q %q", actual, comment)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestListType_lineComment(t *testing.T) ***REMOVED***
	var literals = []struct ***REMOVED***
		src     string
		comment []string
	***REMOVED******REMOVED***
		***REMOVED***
			`foo = [
			1,
			2, # bar
			3,
			]`,
			[]string***REMOVED***"", "# bar", ""***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, l := range literals ***REMOVED***
		p := newParser([]byte(l.src))
		item, err := p.objectItem()
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		list, ok := item.Val.(*ast.ListType)
		if !ok ***REMOVED***
			t.Fatalf("node should be of type LiteralType, got: %T", item.Val)
		***REMOVED***

		if len(list.List) != len(l.comment) ***REMOVED***
			t.Fatalf("bad: %d", len(list.List))
		***REMOVED***

		for i, li := range list.List ***REMOVED***
			lt := li.(*ast.LiteralType)
			comment := l.comment[i]

			if (lt.LineComment == nil) != (comment == "") ***REMOVED***
				t.Fatalf("bad: %s", lt)
			***REMOVED***

			if comment == "" ***REMOVED***
				continue
			***REMOVED***

			actual := lt.LineComment.List[0].Text
			if actual != comment ***REMOVED***
				t.Fatalf("bad: %q %q", actual, comment)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestObjectType(t *testing.T) ***REMOVED***
	var literals = []struct ***REMOVED***
		src      string
		nodeType []ast.Node
		itemLen  int
	***REMOVED******REMOVED***
		***REMOVED***
			`foo = ***REMOVED******REMOVED***`,
			nil,
			0,
		***REMOVED***,
		***REMOVED***
			`foo = ***REMOVED***
				bar = "fatih"
			 ***REMOVED***`,
			[]ast.Node***REMOVED***&ast.LiteralType***REMOVED******REMOVED******REMOVED***,
			1,
		***REMOVED***,
		***REMOVED***
			`foo = ***REMOVED***
				bar = "fatih"
				baz = ["arslan"]
			 ***REMOVED***`,
			[]ast.Node***REMOVED***
				&ast.LiteralType***REMOVED******REMOVED***,
				&ast.ListType***REMOVED******REMOVED***,
			***REMOVED***,
			2,
		***REMOVED***,
		***REMOVED***
			`foo = ***REMOVED***
				bar ***REMOVED******REMOVED***
			 ***REMOVED***`,
			[]ast.Node***REMOVED***
				&ast.ObjectType***REMOVED******REMOVED***,
			***REMOVED***,
			1,
		***REMOVED***,
		***REMOVED***
			`foo ***REMOVED***
				bar ***REMOVED******REMOVED***
				foo = true
			 ***REMOVED***`,
			[]ast.Node***REMOVED***
				&ast.ObjectType***REMOVED******REMOVED***,
				&ast.LiteralType***REMOVED******REMOVED***,
			***REMOVED***,
			2,
		***REMOVED***,
	***REMOVED***

	for _, l := range literals ***REMOVED***
		t.Logf("Source: %s", l.src)

		p := newParser([]byte(l.src))
		// p.enableTrace = true
		item, err := p.objectItem()
		if err != nil ***REMOVED***
			t.Error(err)
			continue
		***REMOVED***

		// we know that the ObjectKey name is foo for all cases, what matters
		// is the object
		obj, ok := item.Val.(*ast.ObjectType)
		if !ok ***REMOVED***
			t.Errorf("node should be of type LiteralType, got: %T", item.Val)
			continue
		***REMOVED***

		// check if the total length of items are correct
		equals(t, l.itemLen, len(obj.List.Items))

		// check if the types are correct
		for i, item := range obj.List.Items ***REMOVED***
			equals(t, reflect.TypeOf(l.nodeType[i]), reflect.TypeOf(item.Val))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestObjectKey(t *testing.T) ***REMOVED***
	keys := []struct ***REMOVED***
		exp []token.Type
		src string
	***REMOVED******REMOVED***
		***REMOVED***[]token.Type***REMOVED***token.IDENT***REMOVED***, `foo ***REMOVED******REMOVED***`***REMOVED***,
		***REMOVED***[]token.Type***REMOVED***token.IDENT***REMOVED***, `foo = ***REMOVED******REMOVED***`***REMOVED***,
		***REMOVED***[]token.Type***REMOVED***token.IDENT***REMOVED***, `foo = bar`***REMOVED***,
		***REMOVED***[]token.Type***REMOVED***token.IDENT***REMOVED***, `foo = 123`***REMOVED***,
		***REMOVED***[]token.Type***REMOVED***token.IDENT***REMOVED***, `foo = "$***REMOVED***var.bar***REMOVED***`***REMOVED***,
		***REMOVED***[]token.Type***REMOVED***token.STRING***REMOVED***, `"foo" ***REMOVED******REMOVED***`***REMOVED***,
		***REMOVED***[]token.Type***REMOVED***token.STRING***REMOVED***, `"foo" = ***REMOVED******REMOVED***`***REMOVED***,
		***REMOVED***[]token.Type***REMOVED***token.STRING***REMOVED***, `"foo" = "$***REMOVED***var.bar***REMOVED***`***REMOVED***,
		***REMOVED***[]token.Type***REMOVED***token.IDENT, token.IDENT***REMOVED***, `foo bar ***REMOVED******REMOVED***`***REMOVED***,
		***REMOVED***[]token.Type***REMOVED***token.IDENT, token.STRING***REMOVED***, `foo "bar" ***REMOVED******REMOVED***`***REMOVED***,
		***REMOVED***[]token.Type***REMOVED***token.STRING, token.IDENT***REMOVED***, `"foo" bar ***REMOVED******REMOVED***`***REMOVED***,
		***REMOVED***[]token.Type***REMOVED***token.IDENT, token.IDENT, token.IDENT***REMOVED***, `foo bar baz ***REMOVED******REMOVED***`***REMOVED***,
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

func TestCommentGroup(t *testing.T) ***REMOVED***
	var cases = []struct ***REMOVED***
		src    string
		groups int
	***REMOVED******REMOVED***
		***REMOVED***"# Hello\n# World", 1***REMOVED***,
		***REMOVED***"# Hello\r\n# Windows", 1***REMOVED***,
	***REMOVED***

	for _, tc := range cases ***REMOVED***
		t.Run(tc.src, func(t *testing.T) ***REMOVED***
			p := newParser([]byte(tc.src))
			file, err := p.Parse()
			if err != nil ***REMOVED***
				t.Fatalf("parse error: %s", err)
			***REMOVED***

			if len(file.Comments) != tc.groups ***REMOVED***
				t.Fatalf("bad: %#v", file.Comments)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

// Official HCL tests
func TestParse(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		Name string
		Err  bool
	***REMOVED******REMOVED***
		***REMOVED***
			"assign_colon.hcl",
			true,
		***REMOVED***,
		***REMOVED***
			"comment.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"comment_crlf.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"comment_lastline.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"comment_single.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"empty.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"list_comma.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"multiple.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"object_list_comma.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"structure.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"structure_basic.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"structure_empty.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"complex.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"complex_crlf.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"types.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"array_comment.hcl",
			false,
		***REMOVED***,
		***REMOVED***
			"array_comment_2.hcl",
			true,
		***REMOVED***,
		***REMOVED***
			"missing_braces.hcl",
			true,
		***REMOVED***,
		***REMOVED***
			"unterminated_object.hcl",
			true,
		***REMOVED***,
		***REMOVED***
			"unterminated_object_2.hcl",
			true,
		***REMOVED***,
		***REMOVED***
			"key_without_value.hcl",
			true,
		***REMOVED***,
		***REMOVED***
			"object_key_without_value.hcl",
			true,
		***REMOVED***,
		***REMOVED***
			"object_key_assign_without_value.hcl",
			true,
		***REMOVED***,
		***REMOVED***
			"object_key_assign_without_value2.hcl",
			true,
		***REMOVED***,
		***REMOVED***
			"object_key_assign_without_value3.hcl",
			true,
		***REMOVED***,
		***REMOVED***
			"git_crypt.hcl",
			true,
		***REMOVED***,
	***REMOVED***

	const fixtureDir = "./test-fixtures"

	for _, tc := range cases ***REMOVED***
		t.Run(tc.Name, func(t *testing.T) ***REMOVED***
			d, err := ioutil.ReadFile(filepath.Join(fixtureDir, tc.Name))
			if err != nil ***REMOVED***
				t.Fatalf("err: %s", err)
			***REMOVED***

			v, err := Parse(d)
			if (err != nil) != tc.Err ***REMOVED***
				t.Fatalf("Input: %s\n\nError: %s\n\nAST: %#v", tc.Name, err, v)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestParse_inline(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		Value string
		Err   bool
	***REMOVED******REMOVED***
		***REMOVED***"t t e***REMOVED******REMOVED******REMOVED******REMOVED***", true***REMOVED***,
		***REMOVED***"o***REMOVED******REMOVED******REMOVED******REMOVED***", true***REMOVED***,
		***REMOVED***"t t e d N***REMOVED******REMOVED******REMOVED******REMOVED***", true***REMOVED***,
		***REMOVED***"t t e d***REMOVED******REMOVED******REMOVED******REMOVED***", true***REMOVED***,
		***REMOVED***"N***REMOVED******REMOVED***N***REMOVED******REMOVED******REMOVED******REMOVED***", true***REMOVED***,
		***REMOVED***"v\nN***REMOVED******REMOVED******REMOVED******REMOVED***", true***REMOVED***,
		***REMOVED***"v=/\n[,", true***REMOVED***,
		***REMOVED***"v=10kb", true***REMOVED***,
		***REMOVED***"v=/foo", true***REMOVED***,
	***REMOVED***

	for _, tc := range cases ***REMOVED***
		t.Logf("Testing: %q", tc.Value)
		ast, err := Parse([]byte(tc.Value))
		if (err != nil) != tc.Err ***REMOVED***
			t.Fatalf("Input: %q\n\nError: %s\n\nAST: %#v", tc.Value, err, ast)
		***REMOVED***
	***REMOVED***
***REMOVED***

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface***REMOVED******REMOVED***) ***REMOVED***
	if !reflect.DeepEqual(exp, act) ***REMOVED***
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	***REMOVED***
***REMOVED***
