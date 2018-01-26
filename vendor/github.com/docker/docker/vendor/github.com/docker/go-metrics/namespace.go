package metrics

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Labels map[string]string

// NewNamespace returns a namespaces that is responsible for managing a collection of
// metrics for a particual namespace and subsystem
//
// labels allows const labels to be added to all metrics created in this namespace
// and are commonly used for data like application version and git commit
func NewNamespace(name, subsystem string, labels Labels) *Namespace ***REMOVED***
	if labels == nil ***REMOVED***
		labels = make(map[string]string)
	***REMOVED***
	return &Namespace***REMOVED***
		name:      name,
		subsystem: subsystem,
		labels:    labels,
	***REMOVED***
***REMOVED***

// Namespace describes a set of metrics that share a namespace and subsystem.
type Namespace struct ***REMOVED***
	name      string
	subsystem string
	labels    Labels
	mu        sync.Mutex
	metrics   []prometheus.Collector
***REMOVED***

// WithConstLabels returns a namespace with the provided set of labels merged
// with the existing constant labels on the namespace.
//
//  Only metrics created with the returned namespace will get the new constant
//  labels.  The returned namespace must be registered separately.
func (n *Namespace) WithConstLabels(labels Labels) *Namespace ***REMOVED***
	n.mu.Lock()
	ns := &Namespace***REMOVED***
		name:      n.name,
		subsystem: n.subsystem,
		labels:    mergeLabels(n.labels, labels),
	***REMOVED***
	n.mu.Unlock()
	return ns
***REMOVED***

func (n *Namespace) NewCounter(name, help string) Counter ***REMOVED***
	c := &counter***REMOVED***pc: prometheus.NewCounter(n.newCounterOpts(name, help))***REMOVED***
	n.Add(c)
	return c
***REMOVED***

func (n *Namespace) NewLabeledCounter(name, help string, labels ...string) LabeledCounter ***REMOVED***
	c := &labeledCounter***REMOVED***pc: prometheus.NewCounterVec(n.newCounterOpts(name, help), labels)***REMOVED***
	n.Add(c)
	return c
***REMOVED***

func (n *Namespace) newCounterOpts(name, help string) prometheus.CounterOpts ***REMOVED***
	return prometheus.CounterOpts***REMOVED***
		Namespace:   n.name,
		Subsystem:   n.subsystem,
		Name:        makeName(name, Total),
		Help:        help,
		ConstLabels: prometheus.Labels(n.labels),
	***REMOVED***
***REMOVED***

func (n *Namespace) NewTimer(name, help string) Timer ***REMOVED***
	t := &timer***REMOVED***
		m: prometheus.NewHistogram(n.newTimerOpts(name, help)),
	***REMOVED***
	n.Add(t)
	return t
***REMOVED***

func (n *Namespace) NewLabeledTimer(name, help string, labels ...string) LabeledTimer ***REMOVED***
	t := &labeledTimer***REMOVED***
		m: prometheus.NewHistogramVec(n.newTimerOpts(name, help), labels),
	***REMOVED***
	n.Add(t)
	return t
***REMOVED***

func (n *Namespace) newTimerOpts(name, help string) prometheus.HistogramOpts ***REMOVED***
	return prometheus.HistogramOpts***REMOVED***
		Namespace:   n.name,
		Subsystem:   n.subsystem,
		Name:        makeName(name, Seconds),
		Help:        help,
		ConstLabels: prometheus.Labels(n.labels),
	***REMOVED***
***REMOVED***

func (n *Namespace) NewGauge(name, help string, unit Unit) Gauge ***REMOVED***
	g := &gauge***REMOVED***
		pg: prometheus.NewGauge(n.newGaugeOpts(name, help, unit)),
	***REMOVED***
	n.Add(g)
	return g
***REMOVED***

func (n *Namespace) NewLabeledGauge(name, help string, unit Unit, labels ...string) LabeledGauge ***REMOVED***
	g := &labeledGauge***REMOVED***
		pg: prometheus.NewGaugeVec(n.newGaugeOpts(name, help, unit), labels),
	***REMOVED***
	n.Add(g)
	return g
***REMOVED***

func (n *Namespace) newGaugeOpts(name, help string, unit Unit) prometheus.GaugeOpts ***REMOVED***
	return prometheus.GaugeOpts***REMOVED***
		Namespace:   n.name,
		Subsystem:   n.subsystem,
		Name:        makeName(name, unit),
		Help:        help,
		ConstLabels: prometheus.Labels(n.labels),
	***REMOVED***
***REMOVED***

func (n *Namespace) Describe(ch chan<- *prometheus.Desc) ***REMOVED***
	n.mu.Lock()
	defer n.mu.Unlock()

	for _, metric := range n.metrics ***REMOVED***
		metric.Describe(ch)
	***REMOVED***
***REMOVED***

func (n *Namespace) Collect(ch chan<- prometheus.Metric) ***REMOVED***
	n.mu.Lock()
	defer n.mu.Unlock()

	for _, metric := range n.metrics ***REMOVED***
		metric.Collect(ch)
	***REMOVED***
***REMOVED***

func (n *Namespace) Add(collector prometheus.Collector) ***REMOVED***
	n.mu.Lock()
	n.metrics = append(n.metrics, collector)
	n.mu.Unlock()
***REMOVED***

func (n *Namespace) NewDesc(name, help string, unit Unit, labels ...string) *prometheus.Desc ***REMOVED***
	name = makeName(name, unit)
	namespace := n.name
	if n.subsystem != "" ***REMOVED***
		namespace = fmt.Sprintf("%s_%s", namespace, n.subsystem)
	***REMOVED***
	name = fmt.Sprintf("%s_%s", namespace, name)
	return prometheus.NewDesc(name, help, labels, prometheus.Labels(n.labels))
***REMOVED***

// mergeLabels merges two or more labels objects into a single map, favoring
// the later labels.
func mergeLabels(lbs ...Labels) Labels ***REMOVED***
	merged := make(Labels)

	for _, target := range lbs ***REMOVED***
		for k, v := range target ***REMOVED***
			merged[k] = v
		***REMOVED***
	***REMOVED***

	return merged
***REMOVED***

func makeName(name string, unit Unit) string ***REMOVED***
	if unit == "" ***REMOVED***
		return name
	***REMOVED***

	return fmt.Sprintf("%s_%s", name, unit)
***REMOVED***
