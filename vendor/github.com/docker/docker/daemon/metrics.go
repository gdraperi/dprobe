package daemon

import (
	"path/filepath"
	"sync"

	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/plugingetter"
	metrics "github.com/docker/go-metrics"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const metricsPluginType = "MetricsCollector"

var (
	containerActions          metrics.LabeledTimer
	imageActions              metrics.LabeledTimer
	networkActions            metrics.LabeledTimer
	engineInfo                metrics.LabeledGauge
	engineCpus                metrics.Gauge
	engineMemory              metrics.Gauge
	healthChecksCounter       metrics.Counter
	healthChecksFailedCounter metrics.Counter

	stateCtr *stateCounter
)

func init() ***REMOVED***
	ns := metrics.NewNamespace("engine", "daemon", nil)
	containerActions = ns.NewLabeledTimer("container_actions", "The number of seconds it takes to process each container action", "action")
	for _, a := range []string***REMOVED***
		"start",
		"changes",
		"commit",
		"create",
		"delete",
	***REMOVED*** ***REMOVED***
		containerActions.WithValues(a).Update(0)
	***REMOVED***

	networkActions = ns.NewLabeledTimer("network_actions", "The number of seconds it takes to process each network action", "action")
	engineInfo = ns.NewLabeledGauge("engine", "The information related to the engine and the OS it is running on", metrics.Unit("info"),
		"version",
		"commit",
		"architecture",
		"graphdriver",
		"kernel", "os",
		"os_type",
		"daemon_id", // ID is a randomly generated unique identifier (e.g. UUID4)
	)
	engineCpus = ns.NewGauge("engine_cpus", "The number of cpus that the host system of the engine has", metrics.Unit("cpus"))
	engineMemory = ns.NewGauge("engine_memory", "The number of bytes of memory that the host system of the engine has", metrics.Bytes)
	healthChecksCounter = ns.NewCounter("health_checks", "The total number of health checks")
	healthChecksFailedCounter = ns.NewCounter("health_checks_failed", "The total number of failed health checks")
	imageActions = ns.NewLabeledTimer("image_actions", "The number of seconds it takes to process each image action", "action")

	stateCtr = newStateCounter(ns.NewDesc("container_states", "The count of containers in various states", metrics.Unit("containers"), "state"))
	ns.Add(stateCtr)

	metrics.Register(ns)
***REMOVED***

type stateCounter struct ***REMOVED***
	mu     sync.Mutex
	states map[string]string
	desc   *prometheus.Desc
***REMOVED***

func newStateCounter(desc *prometheus.Desc) *stateCounter ***REMOVED***
	return &stateCounter***REMOVED***
		states: make(map[string]string),
		desc:   desc,
	***REMOVED***
***REMOVED***

func (ctr *stateCounter) get() (running int, paused int, stopped int) ***REMOVED***
	ctr.mu.Lock()
	defer ctr.mu.Unlock()

	states := map[string]int***REMOVED***
		"running": 0,
		"paused":  0,
		"stopped": 0,
	***REMOVED***
	for _, state := range ctr.states ***REMOVED***
		states[state]++
	***REMOVED***
	return states["running"], states["paused"], states["stopped"]
***REMOVED***

func (ctr *stateCounter) set(id, label string) ***REMOVED***
	ctr.mu.Lock()
	ctr.states[id] = label
	ctr.mu.Unlock()
***REMOVED***

func (ctr *stateCounter) del(id string) ***REMOVED***
	ctr.mu.Lock()
	delete(ctr.states, id)
	ctr.mu.Unlock()
***REMOVED***

func (ctr *stateCounter) Describe(ch chan<- *prometheus.Desc) ***REMOVED***
	ch <- ctr.desc
***REMOVED***

func (ctr *stateCounter) Collect(ch chan<- prometheus.Metric) ***REMOVED***
	running, paused, stopped := ctr.get()
	ch <- prometheus.MustNewConstMetric(ctr.desc, prometheus.GaugeValue, float64(running), "running")
	ch <- prometheus.MustNewConstMetric(ctr.desc, prometheus.GaugeValue, float64(paused), "paused")
	ch <- prometheus.MustNewConstMetric(ctr.desc, prometheus.GaugeValue, float64(stopped), "stopped")
***REMOVED***

func (d *Daemon) cleanupMetricsPlugins() ***REMOVED***
	ls := d.PluginStore.GetAllManagedPluginsByCap(metricsPluginType)
	var wg sync.WaitGroup
	wg.Add(len(ls))

	for _, plugin := range ls ***REMOVED***
		p := plugin
		go func() ***REMOVED***
			defer wg.Done()
			pluginStopMetricsCollection(p)
		***REMOVED***()
	***REMOVED***
	wg.Wait()

	if d.metricsPluginListener != nil ***REMOVED***
		d.metricsPluginListener.Close()
	***REMOVED***
***REMOVED***

type metricsPlugin struct ***REMOVED***
	plugingetter.CompatPlugin
***REMOVED***

func (p metricsPlugin) sock() string ***REMOVED***
	return "metrics.sock"
***REMOVED***

func (p metricsPlugin) sockBase() string ***REMOVED***
	return filepath.Join(p.BasePath(), "run", "docker")
***REMOVED***

func pluginStartMetricsCollection(p plugingetter.CompatPlugin) error ***REMOVED***
	type metricsPluginResponse struct ***REMOVED***
		Err string
	***REMOVED***
	var res metricsPluginResponse
	if err := p.Client().Call(metricsPluginType+".StartMetrics", nil, &res); err != nil ***REMOVED***
		return errors.Wrap(err, "could not start metrics plugin")
	***REMOVED***
	if res.Err != "" ***REMOVED***
		return errors.New(res.Err)
	***REMOVED***
	return nil
***REMOVED***

func pluginStopMetricsCollection(p plugingetter.CompatPlugin) ***REMOVED***
	if err := p.Client().Call(metricsPluginType+".StopMetrics", nil, nil); err != nil ***REMOVED***
		logrus.WithError(err).WithField("name", p.Name()).Error("error stopping metrics collector")
	***REMOVED***

	mp := metricsPlugin***REMOVED***p***REMOVED***
	sockPath := filepath.Join(mp.sockBase(), mp.sock())
	if err := mount.Unmount(sockPath); err != nil ***REMOVED***
		if mounted, _ := mount.Mounted(sockPath); mounted ***REMOVED***
			logrus.WithError(err).WithField("name", p.Name()).WithField("socket", sockPath).Error("error unmounting metrics socket for plugin")
		***REMOVED***
	***REMOVED***
***REMOVED***
