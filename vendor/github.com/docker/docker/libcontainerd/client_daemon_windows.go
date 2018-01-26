package libcontainerd

import (
	"fmt"
	"path/filepath"

	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/windows/hcsshimtypes"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

func summaryFromInterface(i interface***REMOVED******REMOVED***) (*Summary, error) ***REMOVED***
	switch pd := i.(type) ***REMOVED***
	case *hcsshimtypes.ProcessDetails:
		return &Summary***REMOVED***
			CreateTimestamp:              pd.CreatedAt,
			ImageName:                    pd.ImageName,
			KernelTime100ns:              pd.KernelTime_100Ns,
			MemoryCommitBytes:            pd.MemoryCommitBytes,
			MemoryWorkingSetPrivateBytes: pd.MemoryWorkingSetPrivateBytes,
			MemoryWorkingSetSharedBytes:  pd.MemoryWorkingSetSharedBytes,
			ProcessId:                    pd.ProcessID,
			UserTime100ns:                pd.UserTime_100Ns,
		***REMOVED***, nil
	default:
		return nil, errors.Errorf("Unknown process details type %T", pd)
	***REMOVED***
***REMOVED***

func prepareBundleDir(bundleDir string, ociSpec *specs.Spec) (string, error) ***REMOVED***
	return bundleDir, nil
***REMOVED***

func pipeName(containerID, processID, name string) string ***REMOVED***
	return fmt.Sprintf(`\\.\pipe\containerd-%s-%s-%s`, containerID, processID, name)
***REMOVED***

func newFIFOSet(bundleDir, processID string, withStdin, withTerminal bool) *cio.FIFOSet ***REMOVED***
	containerID := filepath.Base(bundleDir)
	config := cio.Config***REMOVED***
		Terminal: withTerminal,
		Stdout:   pipeName(containerID, processID, "stdout"),
	***REMOVED***

	if withStdin ***REMOVED***
		config.Stdin = pipeName(containerID, processID, "stdin")
	***REMOVED***

	if !config.Terminal ***REMOVED***
		config.Stderr = pipeName(containerID, processID, "stderr")
	***REMOVED***

	return cio.NewFIFOSet(config, nil)
***REMOVED***
