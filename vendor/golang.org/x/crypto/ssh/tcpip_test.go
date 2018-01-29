// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"testing"
)

func TestAutoPortListenBroken(t *testing.T) ***REMOVED***
	broken := "SSH-2.0-OpenSSH_5.9hh11"
	works := "SSH-2.0-OpenSSH_6.1"
	if !isBrokenOpenSSHVersion(broken) ***REMOVED***
		t.Errorf("version %q not marked as broken", broken)
	***REMOVED***
	if isBrokenOpenSSHVersion(works) ***REMOVED***
		t.Errorf("version %q marked as broken", works)
	***REMOVED***
***REMOVED***
