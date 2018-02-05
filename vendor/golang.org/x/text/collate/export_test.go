// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collate

// Export for testing.
// TODO: no longer necessary. Remove at some point.

import (
	"fmt"

	"golang.org/x/text/internal/colltab"
)

const (
	defaultSecondary = 0x20
	defaultTertiary  = 0x2
)

type Weights struct ***REMOVED***
	Primary, Secondary, Tertiary, Quaternary int
***REMOVED***

func W(ce ...int) Weights ***REMOVED***
	w := Weights***REMOVED***ce[0], defaultSecondary, defaultTertiary, 0***REMOVED***
	if len(ce) > 1 ***REMOVED***
		w.Secondary = ce[1]
	***REMOVED***
	if len(ce) > 2 ***REMOVED***
		w.Tertiary = ce[2]
	***REMOVED***
	if len(ce) > 3 ***REMOVED***
		w.Quaternary = ce[3]
	***REMOVED***
	return w
***REMOVED***
func (w Weights) String() string ***REMOVED***
	return fmt.Sprintf("[%X.%X.%X.%X]", w.Primary, w.Secondary, w.Tertiary, w.Quaternary)
***REMOVED***

func convertFromWeights(ws []Weights) []colltab.Elem ***REMOVED***
	out := make([]colltab.Elem, len(ws))
	for i, w := range ws ***REMOVED***
		out[i], _ = colltab.MakeElem(w.Primary, w.Secondary, w.Tertiary, 0)
		if out[i] == colltab.Ignore && w.Quaternary > 0 ***REMOVED***
			out[i] = colltab.MakeQuaternary(w.Quaternary)
		***REMOVED***
	***REMOVED***
	return out
***REMOVED***
