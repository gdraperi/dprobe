package homedir

import (
	"path/filepath"
	"testing"
)

func TestGet(t *testing.T) ***REMOVED***
	home := Get()
	if home == "" ***REMOVED***
		t.Fatal("returned home directory is empty")
	***REMOVED***

	if !filepath.IsAbs(home) ***REMOVED***
		t.Fatalf("returned path is not absolute: %s", home)
	***REMOVED***
***REMOVED***

func TestGetShortcutString(t *testing.T) ***REMOVED***
	shortcut := GetShortcutString()
	if shortcut == "" ***REMOVED***
		t.Fatal("returned shortcut string is empty")
	***REMOVED***
***REMOVED***
