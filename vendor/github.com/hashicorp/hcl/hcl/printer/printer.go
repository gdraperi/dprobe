// Package printer implements printing of AST nodes to HCL format.
package printer

import (
	"bytes"
	"io"
	"text/tabwriter"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
)

var DefaultConfig = Config***REMOVED***
	SpacesWidth: 2,
***REMOVED***

// A Config node controls the output of Fprint.
type Config struct ***REMOVED***
	SpacesWidth int // if set, it will use spaces instead of tabs for alignment
***REMOVED***

func (c *Config) Fprint(output io.Writer, node ast.Node) error ***REMOVED***
	p := &printer***REMOVED***
		cfg:                *c,
		comments:           make([]*ast.CommentGroup, 0),
		standaloneComments: make([]*ast.CommentGroup, 0),
		// enableTrace:        true,
	***REMOVED***

	p.collectComments(node)

	if _, err := output.Write(p.unindent(p.output(node))); err != nil ***REMOVED***
		return err
	***REMOVED***

	// flush tabwriter, if any
	var err error
	if tw, _ := output.(*tabwriter.Writer); tw != nil ***REMOVED***
		err = tw.Flush()
	***REMOVED***

	return err
***REMOVED***

// Fprint "pretty-prints" an HCL node to output
// It calls Config.Fprint with default settings.
func Fprint(output io.Writer, node ast.Node) error ***REMOVED***
	return DefaultConfig.Fprint(output, node)
***REMOVED***

// Format formats src HCL and returns the result.
func Format(src []byte) ([]byte, error) ***REMOVED***
	node, err := parser.Parse(src)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var buf bytes.Buffer
	if err := DefaultConfig.Fprint(&buf, node); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Add trailing newline to result
	buf.WriteString("\n")
	return buf.Bytes(), nil
***REMOVED***
