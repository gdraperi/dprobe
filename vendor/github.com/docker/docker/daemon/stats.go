package daemon

import (
	"encoding/json"
	"errors"
	"runtime"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/api/types/versions/v1p20"
	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/ioutils"
)

// ContainerStats writes information about the container to the stream
// given in the config object.
func (daemon *Daemon) ContainerStats(ctx context.Context, prefixOrName string, config *backend.ContainerStatsConfig) error ***REMOVED***
	// Engine API version (used for backwards compatibility)
	apiVersion := config.Version

	container, err := daemon.GetContainer(prefixOrName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If the container is either not running or restarting and requires no stream, return an empty stats.
	if (!container.IsRunning() || container.IsRestarting()) && !config.Stream ***REMOVED***
		return json.NewEncoder(config.OutStream).Encode(&types.StatsJSON***REMOVED***
			Name: container.Name,
			ID:   container.ID***REMOVED***)
	***REMOVED***

	outStream := config.OutStream
	if config.Stream ***REMOVED***
		wf := ioutils.NewWriteFlusher(outStream)
		defer wf.Close()
		wf.Flush()
		outStream = wf
	***REMOVED***

	var preCPUStats types.CPUStats
	var preRead time.Time
	getStatJSON := func(v interface***REMOVED******REMOVED***) *types.StatsJSON ***REMOVED***
		ss := v.(types.StatsJSON)
		ss.Name = container.Name
		ss.ID = container.ID
		ss.PreCPUStats = preCPUStats
		ss.PreRead = preRead
		preCPUStats = ss.CPUStats
		preRead = ss.Read
		return &ss
	***REMOVED***

	enc := json.NewEncoder(outStream)

	updates := daemon.subscribeToContainerStats(container)
	defer daemon.unsubscribeToContainerStats(container, updates)

	noStreamFirstFrame := true
	for ***REMOVED***
		select ***REMOVED***
		case v, ok := <-updates:
			if !ok ***REMOVED***
				return nil
			***REMOVED***

			var statsJSON interface***REMOVED******REMOVED***
			statsJSONPost120 := getStatJSON(v)
			if versions.LessThan(apiVersion, "1.21") ***REMOVED***
				if runtime.GOOS == "windows" ***REMOVED***
					return errors.New("API versions pre v1.21 do not support stats on Windows")
				***REMOVED***
				var (
					rxBytes   uint64
					rxPackets uint64
					rxErrors  uint64
					rxDropped uint64
					txBytes   uint64
					txPackets uint64
					txErrors  uint64
					txDropped uint64
				)
				for _, v := range statsJSONPost120.Networks ***REMOVED***
					rxBytes += v.RxBytes
					rxPackets += v.RxPackets
					rxErrors += v.RxErrors
					rxDropped += v.RxDropped
					txBytes += v.TxBytes
					txPackets += v.TxPackets
					txErrors += v.TxErrors
					txDropped += v.TxDropped
				***REMOVED***
				statsJSON = &v1p20.StatsJSON***REMOVED***
					Stats: statsJSONPost120.Stats,
					Network: types.NetworkStats***REMOVED***
						RxBytes:   rxBytes,
						RxPackets: rxPackets,
						RxErrors:  rxErrors,
						RxDropped: rxDropped,
						TxBytes:   txBytes,
						TxPackets: txPackets,
						TxErrors:  txErrors,
						TxDropped: txDropped,
					***REMOVED***,
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				statsJSON = statsJSONPost120
			***REMOVED***

			if !config.Stream && noStreamFirstFrame ***REMOVED***
				// prime the cpu stats so they aren't 0 in the final output
				noStreamFirstFrame = false
				continue
			***REMOVED***

			if err := enc.Encode(statsJSON); err != nil ***REMOVED***
				return err
			***REMOVED***

			if !config.Stream ***REMOVED***
				return nil
			***REMOVED***
		case <-ctx.Done():
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (daemon *Daemon) subscribeToContainerStats(c *container.Container) chan interface***REMOVED******REMOVED*** ***REMOVED***
	return daemon.statsCollector.Collect(c)
***REMOVED***

func (daemon *Daemon) unsubscribeToContainerStats(c *container.Container, ch chan interface***REMOVED******REMOVED***) ***REMOVED***
	daemon.statsCollector.Unsubscribe(c, ch)
***REMOVED***

// GetContainerStats collects all the stats published by a container
func (daemon *Daemon) GetContainerStats(container *container.Container) (*types.StatsJSON, error) ***REMOVED***
	stats, err := daemon.stats(container)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// We already have the network stats on Windows directly from HCS.
	if !container.Config.NetworkDisabled && runtime.GOOS != "windows" ***REMOVED***
		if stats.Networks, err = daemon.getNetworkStats(container); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return stats, nil
***REMOVED***
