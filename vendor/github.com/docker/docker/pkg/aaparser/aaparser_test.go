package aaparser

import (
	"testing"
)

type versionExpected struct ***REMOVED***
	output  string
	version int
***REMOVED***

func TestParseVersion(t *testing.T) ***REMOVED***
	versions := []versionExpected***REMOVED***
		***REMOVED***
			output: `AppArmor parser version 2.10
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 210000,
		***REMOVED***,
		***REMOVED***
			output: `AppArmor parser version 2.8
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 208000,
		***REMOVED***,
		***REMOVED***
			output: `AppArmor parser version 2.20
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 220000,
		***REMOVED***,
		***REMOVED***
			output: `AppArmor parser version 2.05
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 205000,
		***REMOVED***,
		***REMOVED***
			output: `AppArmor parser version 2.9.95
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 209095,
		***REMOVED***,
		***REMOVED***
			output: `AppArmor parser version 3.14.159
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 314159,
		***REMOVED***,
	***REMOVED***

	for _, v := range versions ***REMOVED***
		version, err := parseVersion(v.output)
		if err != nil ***REMOVED***
			t.Fatalf("expected error to be nil for %#v, got: %v", v, err)
		***REMOVED***
		if version != v.version ***REMOVED***
			t.Fatalf("expected version to be %d, was %d, for: %#v\n", v.version, version, v)
		***REMOVED***
	***REMOVED***
***REMOVED***
