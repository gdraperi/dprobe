package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/hcl/token"
)

// PosError is a parse error that contains a position.
type PosError struct ***REMOVED***
	Pos token.Pos
	Err error
***REMOVED***

func (e *PosError) Error() string ***REMOVED***
	return fmt.Sprintf("At %s: %s", e.Pos, e.Err)
***REMOVED***
