// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package display

import (
	"fmt"
	"testing"

	"golang.org/x/text/internal/testtext"
)

func TestLinking(t *testing.T) ***REMOVED***
	base := getSize(t, `display.Tags(language.English).Name(language.English)`)
	compact := getSize(t, `display.English.Languages().Name(language.English)`)

	if d := base - compact; d < 1.5*1024*1024 ***REMOVED***
		t.Errorf("size(base) - size(compact) = %d - %d = was %d; want > 1.5MB", base, compact, d)
	***REMOVED***
***REMOVED***

func getSize(t *testing.T, main string) int ***REMOVED***
	size, err := testtext.CodeSize(fmt.Sprintf(body, main))
	if err != nil ***REMOVED***
		t.Skipf("skipping link size test; binary size could not be determined: %v", err)
	***REMOVED***
	return size
***REMOVED***

const body = `package main
import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)
func main() ***REMOVED***
	%s
***REMOVED***
`
