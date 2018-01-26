package platforms

import (
	"runtime"
	"strings"
)

// These function are generated from from https://golang.org/src/go/build/syslist.go.
//
// We use switch statements because they are slightly faster than map lookups
// and use a little less memory.

// isKnownOS returns true if we know about the operating system.
//
// The OS value should be normalized before calling this function.
func isKnownOS(os string) bool ***REMOVED***
	switch os ***REMOVED***
	case "android", "darwin", "dragonfly", "freebsd", "linux", "nacl", "netbsd", "openbsd", "plan9", "solaris", "windows", "zos":
		return true
	***REMOVED***
	return false
***REMOVED***

// isKnownArch returns true if we know about the architecture.
//
// The arch value should be normalized before being passed to this function.
func isKnownArch(arch string) bool ***REMOVED***
	switch arch ***REMOVED***
	case "386", "amd64", "amd64p32", "arm", "armbe", "arm64", "arm64be", "ppc64", "ppc64le", "mips", "mipsle", "mips64", "mips64le", "mips64p32", "mips64p32le", "ppc", "s390", "s390x", "sparc", "sparc64":
		return true
	***REMOVED***
	return false
***REMOVED***

func normalizeOS(os string) string ***REMOVED***
	if os == "" ***REMOVED***
		return runtime.GOOS
	***REMOVED***
	os = strings.ToLower(os)

	switch os ***REMOVED***
	case "macos":
		os = "darwin"
	***REMOVED***
	return os
***REMOVED***

// normalizeArch normalizes the architecture.
func normalizeArch(arch, variant string) (string, string) ***REMOVED***
	arch, variant = strings.ToLower(arch), strings.ToLower(variant)
	switch arch ***REMOVED***
	case "i386":
		arch = "386"
		variant = ""
	case "x86_64", "x86-64":
		arch = "amd64"
		variant = ""
	case "aarch64":
		arch = "arm64"
		variant = "" // v8 is implied
	case "armhf":
		arch = "arm"
		variant = ""
	case "armel":
		arch = "arm"
		variant = "v6"
	case "arm":
		switch variant ***REMOVED***
		case "v7", "7":
			variant = "v7"
		case "5", "6", "8":
			variant = "v" + variant
		***REMOVED***
	***REMOVED***

	return arch, variant
***REMOVED***
