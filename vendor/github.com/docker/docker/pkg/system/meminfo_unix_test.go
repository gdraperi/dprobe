// +build linux freebsd

package system

import (
	"strings"
	"testing"

	"github.com/docker/go-units"
)

// TestMemInfo tests parseMemInfo with a static meminfo string
func TestMemInfo(t *testing.T) ***REMOVED***
	const input = `
	MemTotal:      1 kB
	MemFree:       2 kB
	SwapTotal:     3 kB
	SwapFree:      4 kB
	Malformed1:
	Malformed2:    1
	Malformed3:    2 MB
	Malformed4:    X kB
	`
	meminfo, err := parseMemInfo(strings.NewReader(input))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if meminfo.MemTotal != 1*units.KiB ***REMOVED***
		t.Fatalf("Unexpected MemTotal: %d", meminfo.MemTotal)
	***REMOVED***
	if meminfo.MemFree != 2*units.KiB ***REMOVED***
		t.Fatalf("Unexpected MemFree: %d", meminfo.MemFree)
	***REMOVED***
	if meminfo.SwapTotal != 3*units.KiB ***REMOVED***
		t.Fatalf("Unexpected SwapTotal: %d", meminfo.SwapTotal)
	***REMOVED***
	if meminfo.SwapFree != 4*units.KiB ***REMOVED***
		t.Fatalf("Unexpected SwapFree: %d", meminfo.SwapFree)
	***REMOVED***
***REMOVED***
