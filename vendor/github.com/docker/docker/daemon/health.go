package daemon

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/exec"
	"github.com/sirupsen/logrus"
)

const (
	// Longest healthcheck probe output message to store. Longer messages will be truncated.
	maxOutputLen = 4096

	// Default interval between probe runs (from the end of the first to the start of the second).
	// Also the time before the first probe.
	defaultProbeInterval = 30 * time.Second

	// The maximum length of time a single probe run should take. If the probe takes longer
	// than this, the check is considered to have failed.
	defaultProbeTimeout = 30 * time.Second

	// The time given for the container to start before the health check starts considering
	// the container unstable. Defaults to none.
	defaultStartPeriod = 0 * time.Second

	// Default number of consecutive failures of the health check
	// for the container to be considered unhealthy.
	defaultProbeRetries = 3

	// Maximum number of entries to record
	maxLogEntries = 5
)

const (
	// Exit status codes that can be returned by the probe command.

	exitStatusHealthy = 0 // Container is healthy
)

// probe implementations know how to run a particular type of probe.
type probe interface ***REMOVED***
	// Perform one run of the check. Returns the exit code and an optional
	// short diagnostic string.
	run(context.Context, *Daemon, *container.Container) (*types.HealthcheckResult, error)
***REMOVED***

// cmdProbe implements the "CMD" probe type.
type cmdProbe struct ***REMOVED***
	// Run the command with the system's default shell instead of execing it directly.
	shell bool
***REMOVED***

// exec the healthcheck command in the container.
// Returns the exit code and probe output (if any)
func (p *cmdProbe) run(ctx context.Context, d *Daemon, cntr *container.Container) (*types.HealthcheckResult, error) ***REMOVED***
	cmdSlice := strslice.StrSlice(cntr.Config.Healthcheck.Test)[1:]
	if p.shell ***REMOVED***
		cmdSlice = append(getShell(cntr.Config), cmdSlice...)
	***REMOVED***
	entrypoint, args := d.getEntrypointAndArgs(strslice.StrSlice***REMOVED******REMOVED***, cmdSlice)
	execConfig := exec.NewConfig()
	execConfig.OpenStdin = false
	execConfig.OpenStdout = true
	execConfig.OpenStderr = true
	execConfig.ContainerID = cntr.ID
	execConfig.DetachKeys = []byte***REMOVED******REMOVED***
	execConfig.Entrypoint = entrypoint
	execConfig.Args = args
	execConfig.Tty = false
	execConfig.Privileged = false
	execConfig.User = cntr.Config.User
	execConfig.WorkingDir = cntr.Config.WorkingDir

	linkedEnv, err := d.setupLinkedContainers(cntr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	execConfig.Env = container.ReplaceOrAppendEnvValues(cntr.CreateDaemonEnvironment(execConfig.Tty, linkedEnv), execConfig.Env)

	d.registerExecCommand(cntr, execConfig)
	attributes := map[string]string***REMOVED***
		"execID": execConfig.ID,
	***REMOVED***
	d.LogContainerEventWithAttributes(cntr, "exec_create: "+execConfig.Entrypoint+" "+strings.Join(execConfig.Args, " "), attributes)

	output := &limitedBuffer***REMOVED******REMOVED***
	err = d.ContainerExecStart(ctx, execConfig.ID, nil, output, output)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	info, err := d.getExecConfig(execConfig.ID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if info.ExitCode == nil ***REMOVED***
		return nil, fmt.Errorf("healthcheck for container %s has no exit code", cntr.ID)
	***REMOVED***
	// Note: Go's json package will handle invalid UTF-8 for us
	out := output.String()
	return &types.HealthcheckResult***REMOVED***
		End:      time.Now(),
		ExitCode: *info.ExitCode,
		Output:   out,
	***REMOVED***, nil
***REMOVED***

// Update the container's Status.Health struct based on the latest probe's result.
func handleProbeResult(d *Daemon, c *container.Container, result *types.HealthcheckResult, done chan struct***REMOVED******REMOVED***) ***REMOVED***
	c.Lock()
	defer c.Unlock()

	// probe may have been cancelled while waiting on lock. Ignore result then
	select ***REMOVED***
	case <-done:
		return
	default:
	***REMOVED***

	retries := c.Config.Healthcheck.Retries
	if retries <= 0 ***REMOVED***
		retries = defaultProbeRetries
	***REMOVED***

	h := c.State.Health
	oldStatus := h.Status()

	if len(h.Log) >= maxLogEntries ***REMOVED***
		h.Log = append(h.Log[len(h.Log)+1-maxLogEntries:], result)
	***REMOVED*** else ***REMOVED***
		h.Log = append(h.Log, result)
	***REMOVED***

	if result.ExitCode == exitStatusHealthy ***REMOVED***
		h.FailingStreak = 0
		h.SetStatus(types.Healthy)
	***REMOVED*** else ***REMOVED*** // Failure (including invalid exit code)
		shouldIncrementStreak := true

		// If the container is starting (i.e. we never had a successful health check)
		// then we check if we are within the start period of the container in which
		// case we do not increment the failure streak.
		if h.Status() == types.Starting ***REMOVED***
			startPeriod := timeoutWithDefault(c.Config.Healthcheck.StartPeriod, defaultStartPeriod)
			timeSinceStart := result.Start.Sub(c.State.StartedAt)

			// If still within the start period, then don't increment failing streak.
			if timeSinceStart < startPeriod ***REMOVED***
				shouldIncrementStreak = false
			***REMOVED***
		***REMOVED***

		if shouldIncrementStreak ***REMOVED***
			h.FailingStreak++

			if h.FailingStreak >= retries ***REMOVED***
				h.SetStatus(types.Unhealthy)
			***REMOVED***
		***REMOVED***
		// Else we're starting or healthy. Stay in that state.
	***REMOVED***

	// replicate Health status changes
	if err := c.CheckpointTo(d.containersReplica); err != nil ***REMOVED***
		// queries will be inconsistent until the next probe runs or other state mutations
		// checkpoint the container
		logrus.Errorf("Error replicating health state for container %s: %v", c.ID, err)
	***REMOVED***

	current := h.Status()
	if oldStatus != current ***REMOVED***
		d.LogContainerEvent(c, "health_status: "+current)
	***REMOVED***
***REMOVED***

// Run the container's monitoring thread until notified via "stop".
// There is never more than one monitor thread running per container at a time.
func monitor(d *Daemon, c *container.Container, stop chan struct***REMOVED******REMOVED***, probe probe) ***REMOVED***
	probeTimeout := timeoutWithDefault(c.Config.Healthcheck.Timeout, defaultProbeTimeout)
	probeInterval := timeoutWithDefault(c.Config.Healthcheck.Interval, defaultProbeInterval)
	for ***REMOVED***
		select ***REMOVED***
		case <-stop:
			logrus.Debugf("Stop healthcheck monitoring for container %s (received while idle)", c.ID)
			return
		case <-time.After(probeInterval):
			logrus.Debugf("Running health check for container %s ...", c.ID)
			startTime := time.Now()
			ctx, cancelProbe := context.WithTimeout(context.Background(), probeTimeout)
			results := make(chan *types.HealthcheckResult, 1)
			go func() ***REMOVED***
				healthChecksCounter.Inc()
				result, err := probe.run(ctx, d, c)
				if err != nil ***REMOVED***
					healthChecksFailedCounter.Inc()
					logrus.Warnf("Health check for container %s error: %v", c.ID, err)
					results <- &types.HealthcheckResult***REMOVED***
						ExitCode: -1,
						Output:   err.Error(),
						Start:    startTime,
						End:      time.Now(),
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					result.Start = startTime
					logrus.Debugf("Health check for container %s done (exitCode=%d)", c.ID, result.ExitCode)
					results <- result
				***REMOVED***
				close(results)
			***REMOVED***()
			select ***REMOVED***
			case <-stop:
				logrus.Debugf("Stop healthcheck monitoring for container %s (received while probing)", c.ID)
				cancelProbe()
				// Wait for probe to exit (it might take a while to respond to the TERM
				// signal and we don't want dying probes to pile up).
				<-results
				return
			case result := <-results:
				handleProbeResult(d, c, result, stop)
				// Stop timeout
				cancelProbe()
			case <-ctx.Done():
				logrus.Debugf("Health check for container %s taking too long", c.ID)
				handleProbeResult(d, c, &types.HealthcheckResult***REMOVED***
					ExitCode: -1,
					Output:   fmt.Sprintf("Health check exceeded timeout (%v)", probeTimeout),
					Start:    startTime,
					End:      time.Now(),
				***REMOVED***, stop)
				cancelProbe()
				// Wait for probe to exit (it might take a while to respond to the TERM
				// signal and we don't want dying probes to pile up).
				<-results
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Get a suitable probe implementation for the container's healthcheck configuration.
// Nil will be returned if no healthcheck was configured or NONE was set.
func getProbe(c *container.Container) probe ***REMOVED***
	config := c.Config.Healthcheck
	if config == nil || len(config.Test) == 0 ***REMOVED***
		return nil
	***REMOVED***
	switch config.Test[0] ***REMOVED***
	case "CMD":
		return &cmdProbe***REMOVED***shell: false***REMOVED***
	case "CMD-SHELL":
		return &cmdProbe***REMOVED***shell: true***REMOVED***
	case "NONE":
		return nil
	default:
		logrus.Warnf("Unknown healthcheck type '%s' (expected 'CMD') in container %s", config.Test[0], c.ID)
		return nil
	***REMOVED***
***REMOVED***

// Ensure the health-check monitor is running or not, depending on the current
// state of the container.
// Called from monitor.go, with c locked.
func (d *Daemon) updateHealthMonitor(c *container.Container) ***REMOVED***
	h := c.State.Health
	if h == nil ***REMOVED***
		return // No healthcheck configured
	***REMOVED***

	probe := getProbe(c)
	wantRunning := c.Running && !c.Paused && probe != nil
	if wantRunning ***REMOVED***
		if stop := h.OpenMonitorChannel(); stop != nil ***REMOVED***
			go monitor(d, c, stop, probe)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		h.CloseMonitorChannel()
	***REMOVED***
***REMOVED***

// Reset the health state for a newly-started, restarted or restored container.
// initHealthMonitor is called from monitor.go and we should never be running
// two instances at once.
// Called with c locked.
func (d *Daemon) initHealthMonitor(c *container.Container) ***REMOVED***
	// If no healthcheck is setup then don't init the monitor
	if getProbe(c) == nil ***REMOVED***
		return
	***REMOVED***

	// This is needed in case we're auto-restarting
	d.stopHealthchecks(c)

	if h := c.State.Health; h != nil ***REMOVED***
		h.SetStatus(types.Starting)
		h.FailingStreak = 0
	***REMOVED*** else ***REMOVED***
		h := &container.Health***REMOVED******REMOVED***
		h.SetStatus(types.Starting)
		c.State.Health = h
	***REMOVED***

	d.updateHealthMonitor(c)
***REMOVED***

// Called when the container is being stopped (whether because the health check is
// failing or for any other reason).
func (d *Daemon) stopHealthchecks(c *container.Container) ***REMOVED***
	h := c.State.Health
	if h != nil ***REMOVED***
		h.CloseMonitorChannel()
	***REMOVED***
***REMOVED***

// Buffer up to maxOutputLen bytes. Further data is discarded.
type limitedBuffer struct ***REMOVED***
	buf       bytes.Buffer
	mu        sync.Mutex
	truncated bool // indicates that data has been lost
***REMOVED***

// Append to limitedBuffer while there is room.
func (b *limitedBuffer) Write(data []byte) (int, error) ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()

	bufLen := b.buf.Len()
	dataLen := len(data)
	keep := min(maxOutputLen-bufLen, dataLen)
	if keep > 0 ***REMOVED***
		b.buf.Write(data[:keep])
	***REMOVED***
	if keep < dataLen ***REMOVED***
		b.truncated = true
	***REMOVED***
	return dataLen, nil
***REMOVED***

// The contents of the buffer, with "..." appended if it overflowed.
func (b *limitedBuffer) String() string ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()

	out := b.buf.String()
	if b.truncated ***REMOVED***
		out = out + "..."
	***REMOVED***
	return out
***REMOVED***

// If configuredValue is zero, use defaultValue instead.
func timeoutWithDefault(configuredValue time.Duration, defaultValue time.Duration) time.Duration ***REMOVED***
	if configuredValue == 0 ***REMOVED***
		return defaultValue
	***REMOVED***
	return configuredValue
***REMOVED***

func min(x, y int) int ***REMOVED***
	if x < y ***REMOVED***
		return x
	***REMOVED***
	return y
***REMOVED***

func getShell(config *containertypes.Config) []string ***REMOVED***
	if len(config.Shell) != 0 ***REMOVED***
		return config.Shell
	***REMOVED***
	if runtime.GOOS != "windows" ***REMOVED***
		return []string***REMOVED***"/bin/sh", "-c"***REMOVED***
	***REMOVED***
	return []string***REMOVED***"cmd", "/S", "/C"***REMOVED***
***REMOVED***
