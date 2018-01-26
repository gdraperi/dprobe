package daemon

import (
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/libcontainerd"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func toContainerdResources(resources container.Resources) *libcontainerd.Resources ***REMOVED***
	var r libcontainerd.Resources

	r.BlockIO = &specs.LinuxBlockIO***REMOVED***
		Weight: &resources.BlkioWeight,
	***REMOVED***

	shares := uint64(resources.CPUShares)
	r.CPU = &specs.LinuxCPU***REMOVED***
		Shares: &shares,
		Cpus:   resources.CpusetCpus,
		Mems:   resources.CpusetMems,
	***REMOVED***

	var (
		period uint64
		quota  int64
	)
	if resources.NanoCPUs != 0 ***REMOVED***
		period = uint64(100 * time.Millisecond / time.Microsecond)
		quota = resources.NanoCPUs * int64(period) / 1e9
	***REMOVED***
	if quota == 0 && resources.CPUQuota != 0 ***REMOVED***
		quota = resources.CPUQuota
	***REMOVED***
	if period == 0 && resources.CPUPeriod != 0 ***REMOVED***
		period = uint64(resources.CPUPeriod)
	***REMOVED***

	r.CPU.Period = &period
	r.CPU.Quota = &quota

	r.Memory = &specs.LinuxMemory***REMOVED***
		Limit:       &resources.Memory,
		Reservation: &resources.MemoryReservation,
		Kernel:      &resources.KernelMemory,
	***REMOVED***

	if resources.MemorySwap > 0 ***REMOVED***
		r.Memory.Swap = &resources.MemorySwap
	***REMOVED***

	return &r
***REMOVED***
