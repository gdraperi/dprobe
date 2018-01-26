package runtime

import (
	"strconv"

	"github.com/containerd/typeurl"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func init() ***REMOVED***
	const prefix = "types.containerd.io"
	// register TypeUrls for commonly marshaled external types
	major := strconv.Itoa(specs.VersionMajor)
	typeurl.Register(&specs.Spec***REMOVED******REMOVED***, prefix, "opencontainers/runtime-spec", major, "Spec")
	typeurl.Register(&specs.Process***REMOVED******REMOVED***, prefix, "opencontainers/runtime-spec", major, "Process")
	typeurl.Register(&specs.LinuxResources***REMOVED******REMOVED***, prefix, "opencontainers/runtime-spec", major, "LinuxResources")
	typeurl.Register(&specs.WindowsResources***REMOVED******REMOVED***, prefix, "opencontainers/runtime-spec", major, "WindowsResources")
***REMOVED***
