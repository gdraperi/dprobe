package ast

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/hcl/token"
)

func TestObjectListFilter(t *testing.T) ***REMOVED***
	var cases = []struct ***REMOVED***
		Filter []string
		Input  []*ObjectItem
		Output []*ObjectItem
	***REMOVED******REMOVED***
		***REMOVED***
			[]string***REMOVED***"foo"***REMOVED***,
			[]*ObjectItem***REMOVED***
				&ObjectItem***REMOVED***
					Keys: []*ObjectKey***REMOVED***
						&ObjectKey***REMOVED***
							Token: token.Token***REMOVED***Type: token.STRING, Text: `"foo"`***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			[]*ObjectItem***REMOVED***
				&ObjectItem***REMOVED***
					Keys: []*ObjectKey***REMOVED******REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		***REMOVED***
			[]string***REMOVED***"foo"***REMOVED***,
			[]*ObjectItem***REMOVED***
				&ObjectItem***REMOVED***
					Keys: []*ObjectKey***REMOVED***
						&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"foo"`***REMOVED******REMOVED***,
						&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"bar"`***REMOVED******REMOVED***,
					***REMOVED***,
				***REMOVED***,
				&ObjectItem***REMOVED***
					Keys: []*ObjectKey***REMOVED***
						&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"baz"`***REMOVED******REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			[]*ObjectItem***REMOVED***
				&ObjectItem***REMOVED***
					Keys: []*ObjectKey***REMOVED***
						&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"bar"`***REMOVED******REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, tc := range cases ***REMOVED***
		input := &ObjectList***REMOVED***Items: tc.Input***REMOVED***
		expected := &ObjectList***REMOVED***Items: tc.Output***REMOVED***
		if actual := input.Filter(tc.Filter...); !reflect.DeepEqual(actual, expected) ***REMOVED***
			t.Fatalf("in order: input, expected, actual\n\n%#v\n\n%#v\n\n%#v", input, expected, actual)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWalk(t *testing.T) ***REMOVED***
	items := []*ObjectItem***REMOVED***
		&ObjectItem***REMOVED***
			Keys: []*ObjectKey***REMOVED***
				&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"foo"`***REMOVED******REMOVED***,
				&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"bar"`***REMOVED******REMOVED***,
			***REMOVED***,
			Val: &LiteralType***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"example"`***REMOVED******REMOVED***,
		***REMOVED***,
		&ObjectItem***REMOVED***
			Keys: []*ObjectKey***REMOVED***
				&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"baz"`***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	node := &ObjectList***REMOVED***Items: items***REMOVED***

	order := []string***REMOVED***
		"*ast.ObjectList",
		"*ast.ObjectItem",
		"*ast.ObjectKey",
		"*ast.ObjectKey",
		"*ast.LiteralType",
		"*ast.ObjectItem",
		"*ast.ObjectKey",
	***REMOVED***
	count := 0

	Walk(node, func(n Node) (Node, bool) ***REMOVED***
		if n == nil ***REMOVED***
			return n, false
		***REMOVED***

		typeName := reflect.TypeOf(n).String()
		if order[count] != typeName ***REMOVED***
			t.Errorf("expected '%s' got: '%s'", order[count], typeName)
		***REMOVED***
		count++
		return n, true
	***REMOVED***)
***REMOVED***

func TestWalkEquality(t *testing.T) ***REMOVED***
	items := []*ObjectItem***REMOVED***
		&ObjectItem***REMOVED***
			Keys: []*ObjectKey***REMOVED***
				&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"foo"`***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		&ObjectItem***REMOVED***
			Keys: []*ObjectKey***REMOVED***
				&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"bar"`***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	node := &ObjectList***REMOVED***Items: items***REMOVED***

	rewritten := Walk(node, func(n Node) (Node, bool) ***REMOVED*** return n, true ***REMOVED***)

	newNode, ok := rewritten.(*ObjectList)
	if !ok ***REMOVED***
		t.Fatalf("expected Objectlist, got %T", rewritten)
	***REMOVED***

	if !reflect.DeepEqual(node, newNode) ***REMOVED***
		t.Fatal("rewritten node is not equal to the given node")
	***REMOVED***

	if len(newNode.Items) != 2 ***REMOVED***
		t.Error("expected newNode length 2, got: %d", len(newNode.Items))
	***REMOVED***

	expected := []string***REMOVED***
		`"foo"`,
		`"bar"`,
	***REMOVED***

	for i, item := range newNode.Items ***REMOVED***
		if len(item.Keys) != 1 ***REMOVED***
			t.Error("expected keys newNode length 1, got: %d", len(item.Keys))
		***REMOVED***

		if item.Keys[0].Token.Text != expected[i] ***REMOVED***
			t.Errorf("expected key %s, got %s", expected[i], item.Keys[0].Token.Text)
		***REMOVED***

		if item.Val != nil ***REMOVED***
			t.Errorf("expected item value should be nil")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWalkRewrite(t *testing.T) ***REMOVED***
	items := []*ObjectItem***REMOVED***
		&ObjectItem***REMOVED***
			Keys: []*ObjectKey***REMOVED***
				&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"foo"`***REMOVED******REMOVED***,
				&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"bar"`***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		&ObjectItem***REMOVED***
			Keys: []*ObjectKey***REMOVED***
				&ObjectKey***REMOVED***Token: token.Token***REMOVED***Type: token.STRING, Text: `"baz"`***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	node := &ObjectList***REMOVED***Items: items***REMOVED***

	suffix := "_example"
	node = Walk(node, func(n Node) (Node, bool) ***REMOVED***
		switch i := n.(type) ***REMOVED***
		case *ObjectKey:
			i.Token.Text = i.Token.Text + suffix
			n = i
		***REMOVED***
		return n, true
	***REMOVED***).(*ObjectList)

	Walk(node, func(n Node) (Node, bool) ***REMOVED***
		switch i := n.(type) ***REMOVED***
		case *ObjectKey:
			if !strings.HasSuffix(i.Token.Text, suffix) ***REMOVED***
				t.Errorf("Token '%s' should have suffix: %s", i.Token.Text, suffix)
			***REMOVED***
		***REMOVED***
		return n, true
	***REMOVED***)

***REMOVED***
