package longpath

import (
	"strings"
	"testing"
)

func TestStandardLongPath(t *testing.T) ***REMOVED***
	c := `C:\simple\path`
	longC := AddPrefix(c)
	if !strings.EqualFold(longC, `\\?\C:\simple\path`) ***REMOVED***
		t.Errorf("Wrong long path returned. Original = %s ; Long = %s", c, longC)
	***REMOVED***
***REMOVED***

func TestUNCLongPath(t *testing.T) ***REMOVED***
	c := `\\server\share\path`
	longC := AddPrefix(c)
	if !strings.EqualFold(longC, `\\?\UNC\server\share\path`) ***REMOVED***
		t.Errorf("Wrong UNC long path returned. Original = %s ; Long = %s", c, longC)
	***REMOVED***
***REMOVED***
