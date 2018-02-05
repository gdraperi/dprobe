package cobra

import (
	"bytes"
	"strings"
	"testing"
)

func TestZshCompletion(t *testing.T) ***REMOVED***
	tcs := []struct ***REMOVED***
		name                string
		root                *Command
		expectedExpressions []string
	***REMOVED******REMOVED***
		***REMOVED***
			name:                "trivial",
			root:                &Command***REMOVED***Use: "trivialapp"***REMOVED***,
			expectedExpressions: []string***REMOVED***"#compdef trivial"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "linear",
			root: func() *Command ***REMOVED***
				r := &Command***REMOVED***Use: "linear"***REMOVED***

				sub1 := &Command***REMOVED***Use: "sub1"***REMOVED***
				r.AddCommand(sub1)

				sub2 := &Command***REMOVED***Use: "sub2"***REMOVED***
				sub1.AddCommand(sub2)

				sub3 := &Command***REMOVED***Use: "sub3"***REMOVED***
				sub2.AddCommand(sub3)
				return r
			***REMOVED***(),
			expectedExpressions: []string***REMOVED***"sub1", "sub2", "sub3"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "flat",
			root: func() *Command ***REMOVED***
				r := &Command***REMOVED***Use: "flat"***REMOVED***
				r.AddCommand(&Command***REMOVED***Use: "c1"***REMOVED***)
				r.AddCommand(&Command***REMOVED***Use: "c2"***REMOVED***)
				return r
			***REMOVED***(),
			expectedExpressions: []string***REMOVED***"(c1 c2)"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "tree",
			root: func() *Command ***REMOVED***
				r := &Command***REMOVED***Use: "tree"***REMOVED***

				sub1 := &Command***REMOVED***Use: "sub1"***REMOVED***
				r.AddCommand(sub1)

				sub11 := &Command***REMOVED***Use: "sub11"***REMOVED***
				sub12 := &Command***REMOVED***Use: "sub12"***REMOVED***

				sub1.AddCommand(sub11)
				sub1.AddCommand(sub12)

				sub2 := &Command***REMOVED***Use: "sub2"***REMOVED***
				r.AddCommand(sub2)

				sub21 := &Command***REMOVED***Use: "sub21"***REMOVED***
				sub22 := &Command***REMOVED***Use: "sub22"***REMOVED***

				sub2.AddCommand(sub21)
				sub2.AddCommand(sub22)

				return r
			***REMOVED***(),
			expectedExpressions: []string***REMOVED***"(sub11 sub12)", "(sub21 sub22)"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, tc := range tcs ***REMOVED***
		t.Run(tc.name, func(t *testing.T) ***REMOVED***
			buf := new(bytes.Buffer)
			tc.root.GenZshCompletion(buf)
			output := buf.String()

			for _, expectedExpression := range tc.expectedExpressions ***REMOVED***
				if !strings.Contains(output, expectedExpression) ***REMOVED***
					t.Errorf("Expected completion to contain %q somewhere; got %q", expectedExpression, output)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
