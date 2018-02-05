// Copyright 2017 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import "flag"

// MustFlag sets flags that are skipped by dst.Parse when p contains
// the respective key for flag.Flag.Name.
//
// It's use is recommended with command line arguments as in:
// 	flag.Parse()
// 	p.MustFlag(flag.CommandLine)
func (p *Properties) MustFlag(dst *flag.FlagSet) ***REMOVED***
	m := make(map[string]*flag.Flag)
	dst.VisitAll(func(f *flag.Flag) ***REMOVED***
		m[f.Name] = f
	***REMOVED***)
	dst.Visit(func(f *flag.Flag) ***REMOVED***
		delete(m, f.Name) // overridden
	***REMOVED***)

	for name, f := range m ***REMOVED***
		v, ok := p.Get(name)
		if !ok ***REMOVED***
			continue
		***REMOVED***

		if err := f.Value.Set(v); err != nil ***REMOVED***
			ErrorHandler(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
