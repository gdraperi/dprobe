// Position support for go-toml

package toml

import (
	"fmt"
)

// Position of a document element within a TOML document.
//
// Line and Col are both 1-indexed positions for the element's line number and
// column number, respectively.  Values of zero or less will cause Invalid(),
// to return true.
type Position struct ***REMOVED***
	Line int // line within the document
	Col  int // column within the line
***REMOVED***

// String representation of the position.
// Displays 1-indexed line and column numbers.
func (p Position) String() string ***REMOVED***
	return fmt.Sprintf("(%d, %d)", p.Line, p.Col)
***REMOVED***

// Invalid returns whether or not the position is valid (i.e. with negative or
// null values)
func (p Position) Invalid() bool ***REMOVED***
	return p.Line <= 0 || p.Col <= 0
***REMOVED***
