package printer

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/hcl/parser"
)

var update = flag.Bool("update", false, "update golden files")

const (
	dataDir = "testdata"
)

type entry struct ***REMOVED***
	source, golden string
***REMOVED***

// Use go test -update to create/update the respective golden files.
var data = []entry***REMOVED***
	***REMOVED***"complexhcl.input", "complexhcl.golden"***REMOVED***,
	***REMOVED***"list.input", "list.golden"***REMOVED***,
	***REMOVED***"list_comment.input", "list_comment.golden"***REMOVED***,
	***REMOVED***"comment.input", "comment.golden"***REMOVED***,
	***REMOVED***"comment_crlf.input", "comment.golden"***REMOVED***,
	***REMOVED***"comment_aligned.input", "comment_aligned.golden"***REMOVED***,
	***REMOVED***"comment_array.input", "comment_array.golden"***REMOVED***,
	***REMOVED***"comment_end_file.input", "comment_end_file.golden"***REMOVED***,
	***REMOVED***"comment_multiline_indent.input", "comment_multiline_indent.golden"***REMOVED***,
	***REMOVED***"comment_multiline_no_stanza.input", "comment_multiline_no_stanza.golden"***REMOVED***,
	***REMOVED***"comment_multiline_stanza.input", "comment_multiline_stanza.golden"***REMOVED***,
	***REMOVED***"comment_newline.input", "comment_newline.golden"***REMOVED***,
	***REMOVED***"comment_object_multi.input", "comment_object_multi.golden"***REMOVED***,
	***REMOVED***"comment_standalone.input", "comment_standalone.golden"***REMOVED***,
	***REMOVED***"empty_block.input", "empty_block.golden"***REMOVED***,
	***REMOVED***"list_of_objects.input", "list_of_objects.golden"***REMOVED***,
	***REMOVED***"multiline_string.input", "multiline_string.golden"***REMOVED***,
	***REMOVED***"object_singleline.input", "object_singleline.golden"***REMOVED***,
	***REMOVED***"object_with_heredoc.input", "object_with_heredoc.golden"***REMOVED***,
***REMOVED***

func TestFiles(t *testing.T) ***REMOVED***
	for _, e := range data ***REMOVED***
		source := filepath.Join(dataDir, e.source)
		golden := filepath.Join(dataDir, e.golden)
		t.Run(e.source, func(t *testing.T) ***REMOVED***
			check(t, source, golden)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func check(t *testing.T, source, golden string) ***REMOVED***
	src, err := ioutil.ReadFile(source)
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***

	res, err := format(src)
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***

	// update golden files if necessary
	if *update ***REMOVED***
		if err := ioutil.WriteFile(golden, res, 0644); err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
		return
	***REMOVED***

	// get golden
	gld, err := ioutil.ReadFile(golden)
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***

	// formatted source and golden must be the same
	if err := diff(source, golden, res, gld); err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
***REMOVED***

// diff compares a and b.
func diff(aname, bname string, a, b []byte) error ***REMOVED***
	var buf bytes.Buffer // holding long error message

	// compare lengths
	if len(a) != len(b) ***REMOVED***
		fmt.Fprintf(&buf, "\nlength changed: len(%s) = %d, len(%s) = %d", aname, len(a), bname, len(b))
	***REMOVED***

	// compare contents
	line := 1
	offs := 1
	for i := 0; i < len(a) && i < len(b); i++ ***REMOVED***
		ch := a[i]
		if ch != b[i] ***REMOVED***
			fmt.Fprintf(&buf, "\n%s:%d:%d: %q", aname, line, i-offs+1, lineAt(a, offs))
			fmt.Fprintf(&buf, "\n%s:%d:%d: %q", bname, line, i-offs+1, lineAt(b, offs))
			fmt.Fprintf(&buf, "\n\n")
			break
		***REMOVED***
		if ch == '\n' ***REMOVED***
			line++
			offs = i + 1
		***REMOVED***
	***REMOVED***

	if buf.Len() > 0 ***REMOVED***
		return errors.New(buf.String())
	***REMOVED***
	return nil
***REMOVED***

// format parses src, prints the corresponding AST, verifies the resulting
// src is syntactically correct, and returns the resulting src or an error
// if any.
func format(src []byte) ([]byte, error) ***REMOVED***
	formatted, err := Format(src)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// make sure formatted output is syntactically correct
	if _, err := parser.Parse(formatted); err != nil ***REMOVED***
		return nil, fmt.Errorf("parse: %s\n%s", err, formatted)
	***REMOVED***

	return formatted, nil
***REMOVED***

// lineAt returns the line in text starting at offset offs.
func lineAt(text []byte, offs int) []byte ***REMOVED***
	i := offs
	for i < len(text) && text[i] != '\n' ***REMOVED***
		i++
	***REMOVED***
	return text[offs:i]
***REMOVED***
