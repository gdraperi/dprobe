package namesgenerator

import (
	"strings"
	"testing"
)

func TestNameFormat(t *testing.T) ***REMOVED***
	name := GetRandomName(0)
	if !strings.Contains(name, "_") ***REMOVED***
		t.Fatalf("Generated name does not contain an underscore")
	***REMOVED***
	if strings.ContainsAny(name, "0123456789") ***REMOVED***
		t.Fatalf("Generated name contains numbers!")
	***REMOVED***
***REMOVED***

func TestNameRetries(t *testing.T) ***REMOVED***
	name := GetRandomName(1)
	if !strings.Contains(name, "_") ***REMOVED***
		t.Fatalf("Generated name does not contain an underscore")
	***REMOVED***
	if !strings.ContainsAny(name, "0123456789") ***REMOVED***
		t.Fatalf("Generated name doesn't contain a number")
	***REMOVED***

***REMOVED***
