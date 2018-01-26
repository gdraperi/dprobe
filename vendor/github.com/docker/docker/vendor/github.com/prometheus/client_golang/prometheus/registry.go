// Copyright 2014 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Copyright (c) 2013, The Prometheus Authors
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package prometheus

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/common/expfmt"

	dto "github.com/prometheus/client_model/go"
)

var (
	defRegistry   = newDefaultRegistry()
	errAlreadyReg = errors.New("duplicate metrics collector registration attempted")
)

// Constants relevant to the HTTP interface.
const (
	// APIVersion is the version of the format of the exported data.  This
	// will match this library's version, which subscribes to the Semantic
	// Versioning scheme.
	APIVersion = "0.0.4"

	// DelimitedTelemetryContentType is the content type set on telemetry
	// data responses in delimited protobuf format.
	DelimitedTelemetryContentType = `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited`
	// TextTelemetryContentType is the content type set on telemetry data
	// responses in text format.
	TextTelemetryContentType = `text/plain; version=` + APIVersion
	// ProtoTextTelemetryContentType is the content type set on telemetry
	// data responses in protobuf text format.  (Only used for debugging.)
	ProtoTextTelemetryContentType = `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=text`
	// ProtoCompactTextTelemetryContentType is the content type set on
	// telemetry data responses in protobuf compact text format.  (Only used
	// for debugging.)
	ProtoCompactTextTelemetryContentType = `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=compact-text`

	// Constants for object pools.
	numBufs           = 4
	numMetricFamilies = 1000
	numMetrics        = 10000

	// Capacity for the channel to collect metrics and descriptors.
	capMetricChan = 1000
	capDescChan   = 10

	contentTypeHeader     = "Content-Type"
	contentLengthHeader   = "Content-Length"
	contentEncodingHeader = "Content-Encoding"

	acceptEncodingHeader = "Accept-Encoding"
	acceptHeader         = "Accept"
)

// Handler returns the HTTP handler for the global Prometheus registry. It is
// already instrumented with InstrumentHandler (using "prometheus" as handler
// name). Usually the handler is used to handle the "/metrics" endpoint.
//
// Please note the issues described in the doc comment of InstrumentHandler. You
// might want to consider using UninstrumentedHandler instead.
func Handler() http.Handler ***REMOVED***
	return InstrumentHandler("prometheus", defRegistry)
***REMOVED***

// UninstrumentedHandler works in the same way as Handler, but the returned HTTP
// handler is not instrumented. This is useful if no instrumentation is desired
// (for whatever reason) or if the instrumentation has to happen with a
// different handler name (or with a different instrumentation approach
// altogether). See the InstrumentHandler example.
func UninstrumentedHandler() http.Handler ***REMOVED***
	return defRegistry
***REMOVED***

// Register registers a new Collector to be included in metrics collection. It
// returns an error if the descriptors provided by the Collector are invalid or
// if they - in combination with descriptors of already registered Collectors -
// do not fulfill the consistency and uniqueness criteria described in the Desc
// documentation.
//
// Do not register the same Collector multiple times concurrently. (Registering
// the same Collector twice would result in an error anyway, but on top of that,
// it is not safe to do so concurrently.)
func Register(m Collector) error ***REMOVED***
	_, err := defRegistry.Register(m)
	return err
***REMOVED***

// MustRegister works like Register but panics where Register would have
// returned an error. MustRegister is also Variadic, where Register only
// accepts a single Collector to register.
func MustRegister(m ...Collector) ***REMOVED***
	for i := range m ***REMOVED***
		if err := Register(m[i]); err != nil ***REMOVED***
			panic(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// RegisterOrGet works like Register but does not return an error if a Collector
// is registered that equals a previously registered Collector. (Two Collectors
// are considered equal if their Describe method yields the same set of
// descriptors.) Instead, the previously registered Collector is returned (which
// is helpful if the new and previously registered Collectors are equal but not
// identical, i.e. not pointers to the same object).
//
// As for Register, it is still not safe to call RegisterOrGet with the same
// Collector multiple times concurrently.
func RegisterOrGet(m Collector) (Collector, error) ***REMOVED***
	return defRegistry.RegisterOrGet(m)
***REMOVED***

// MustRegisterOrGet works like RegisterOrGet but panics where RegisterOrGet
// would have returned an error.
func MustRegisterOrGet(m Collector) Collector ***REMOVED***
	existing, err := RegisterOrGet(m)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return existing
***REMOVED***

// Unregister unregisters the Collector that equals the Collector passed in as
// an argument. (Two Collectors are considered equal if their Describe method
// yields the same set of descriptors.) The function returns whether a Collector
// was unregistered.
func Unregister(c Collector) bool ***REMOVED***
	return defRegistry.Unregister(c)
***REMOVED***

// SetMetricFamilyInjectionHook sets a function that is called whenever metrics
// are collected. The hook function must be set before metrics collection begins
// (i.e. call SetMetricFamilyInjectionHook before setting the HTTP handler.) The
// MetricFamily protobufs returned by the hook function are merged with the
// metrics collected in the usual way.
//
// This is a way to directly inject MetricFamily protobufs managed and owned by
// the caller. The caller has full responsibility. As no registration of the
// injected metrics has happened, there is no descriptor to check against, and
// there are no registration-time checks. If collect-time checks are disabled
// (see function EnableCollectChecks), no sanity checks are performed on the
// returned protobufs at all. If collect-checks are enabled, type and uniqueness
// checks are performed, but no further consistency checks (which would require
// knowledge of a metric descriptor).
//
// Sorting concerns: The caller is responsible for sorting the label pairs in
// each metric. However, the order of metrics will be sorted by the registry as
// it is required anyway after merging with the metric families collected
// conventionally.
//
// The function must be callable at any time and concurrently.
func SetMetricFamilyInjectionHook(hook func() []*dto.MetricFamily) ***REMOVED***
	defRegistry.metricFamilyInjectionHook = hook
***REMOVED***

// PanicOnCollectError sets the behavior whether a panic is caused upon an error
// while metrics are collected and served to the HTTP endpoint. By default, an
// internal server error (status code 500) is served with an error message.
func PanicOnCollectError(b bool) ***REMOVED***
	defRegistry.panicOnCollectError = b
***REMOVED***

// EnableCollectChecks enables (or disables) additional consistency checks
// during metrics collection. These additional checks are not enabled by default
// because they inflict a performance penalty and the errors they check for can
// only happen if the used Metric and Collector types have internal programming
// errors. It can be helpful to enable these checks while working with custom
// Collectors or Metrics whose correctness is not well established yet.
func EnableCollectChecks(b bool) ***REMOVED***
	defRegistry.collectChecksEnabled = b
***REMOVED***

// encoder is a function that writes a dto.MetricFamily to an io.Writer in a
// certain encoding. It returns the number of bytes written and any error
// encountered.  Note that pbutil.WriteDelimited and pbutil.MetricFamilyToText
// are encoders.
type encoder func(io.Writer, *dto.MetricFamily) (int, error)

type registry struct ***REMOVED***
	mtx                       sync.RWMutex
	collectorsByID            map[uint64]Collector // ID is a hash of the descIDs.
	descIDs                   map[uint64]struct***REMOVED******REMOVED***
	dimHashesByName           map[string]uint64
	bufPool                   chan *bytes.Buffer
	metricFamilyPool          chan *dto.MetricFamily
	metricPool                chan *dto.Metric
	metricFamilyInjectionHook func() []*dto.MetricFamily

	panicOnCollectError, collectChecksEnabled bool
***REMOVED***

func (r *registry) Register(c Collector) (Collector, error) ***REMOVED***
	descChan := make(chan *Desc, capDescChan)
	go func() ***REMOVED***
		c.Describe(descChan)
		close(descChan)
	***REMOVED***()

	newDescIDs := map[uint64]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	newDimHashesByName := map[string]uint64***REMOVED******REMOVED***
	var collectorID uint64 // Just a sum of all desc IDs.
	var duplicateDescErr error

	r.mtx.Lock()
	defer r.mtx.Unlock()
	// Coduct various tests...
	for desc := range descChan ***REMOVED***

		// Is the descriptor valid at all?
		if desc.err != nil ***REMOVED***
			return c, fmt.Errorf("descriptor %s is invalid: %s", desc, desc.err)
		***REMOVED***

		// Is the descID unique?
		// (In other words: Is the fqName + constLabel combination unique?)
		if _, exists := r.descIDs[desc.id]; exists ***REMOVED***
			duplicateDescErr = fmt.Errorf("descriptor %s already exists with the same fully-qualified name and const label values", desc)
		***REMOVED***
		// If it is not a duplicate desc in this collector, add it to
		// the collectorID.  (We allow duplicate descs within the same
		// collector, but their existence must be a no-op.)
		if _, exists := newDescIDs[desc.id]; !exists ***REMOVED***
			newDescIDs[desc.id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			collectorID += desc.id
		***REMOVED***

		// Are all the label names and the help string consistent with
		// previous descriptors of the same name?
		// First check existing descriptors...
		if dimHash, exists := r.dimHashesByName[desc.fqName]; exists ***REMOVED***
			if dimHash != desc.dimHash ***REMOVED***
				return nil, fmt.Errorf("a previously registered descriptor with the same fully-qualified name as %s has different label names or a different help string", desc)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// ...then check the new descriptors already seen.
			if dimHash, exists := newDimHashesByName[desc.fqName]; exists ***REMOVED***
				if dimHash != desc.dimHash ***REMOVED***
					return nil, fmt.Errorf("descriptors reported by collector have inconsistent label names or help strings for the same fully-qualified name, offender is %s", desc)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				newDimHashesByName[desc.fqName] = desc.dimHash
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Did anything happen at all?
	if len(newDescIDs) == 0 ***REMOVED***
		return nil, errors.New("collector has no descriptors")
	***REMOVED***
	if existing, exists := r.collectorsByID[collectorID]; exists ***REMOVED***
		return existing, errAlreadyReg
	***REMOVED***
	// If the collectorID is new, but at least one of the descs existed
	// before, we are in trouble.
	if duplicateDescErr != nil ***REMOVED***
		return nil, duplicateDescErr
	***REMOVED***

	// Only after all tests have passed, actually register.
	r.collectorsByID[collectorID] = c
	for hash := range newDescIDs ***REMOVED***
		r.descIDs[hash] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	for name, dimHash := range newDimHashesByName ***REMOVED***
		r.dimHashesByName[name] = dimHash
	***REMOVED***
	return c, nil
***REMOVED***

func (r *registry) RegisterOrGet(m Collector) (Collector, error) ***REMOVED***
	existing, err := r.Register(m)
	if err != nil && err != errAlreadyReg ***REMOVED***
		return nil, err
	***REMOVED***
	return existing, nil
***REMOVED***

func (r *registry) Unregister(c Collector) bool ***REMOVED***
	descChan := make(chan *Desc, capDescChan)
	go func() ***REMOVED***
		c.Describe(descChan)
		close(descChan)
	***REMOVED***()

	descIDs := map[uint64]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	var collectorID uint64 // Just a sum of the desc IDs.
	for desc := range descChan ***REMOVED***
		if _, exists := descIDs[desc.id]; !exists ***REMOVED***
			collectorID += desc.id
			descIDs[desc.id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	r.mtx.RLock()
	if _, exists := r.collectorsByID[collectorID]; !exists ***REMOVED***
		r.mtx.RUnlock()
		return false
	***REMOVED***
	r.mtx.RUnlock()

	r.mtx.Lock()
	defer r.mtx.Unlock()

	delete(r.collectorsByID, collectorID)
	for id := range descIDs ***REMOVED***
		delete(r.descIDs, id)
	***REMOVED***
	// dimHashesByName is left untouched as those must be consistent
	// throughout the lifetime of a program.
	return true
***REMOVED***

func (r *registry) Push(job, instance, pushURL, method string) error ***REMOVED***
	if !strings.Contains(pushURL, "://") ***REMOVED***
		pushURL = "http://" + pushURL
	***REMOVED***
	if strings.HasSuffix(pushURL, "/") ***REMOVED***
		pushURL = pushURL[:len(pushURL)-1]
	***REMOVED***
	pushURL = fmt.Sprintf("%s/metrics/jobs/%s", pushURL, url.QueryEscape(job))
	if instance != "" ***REMOVED***
		pushURL += "/instances/" + url.QueryEscape(instance)
	***REMOVED***
	buf := r.getBuf()
	defer r.giveBuf(buf)
	if err := r.writePB(expfmt.NewEncoder(buf, expfmt.FmtProtoDelim)); err != nil ***REMOVED***
		if r.panicOnCollectError ***REMOVED***
			panic(err)
		***REMOVED***
		return err
	***REMOVED***
	req, err := http.NewRequest(method, pushURL, buf)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	req.Header.Set(contentTypeHeader, DelimitedTelemetryContentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer resp.Body.Close()
	if resp.StatusCode != 202 ***REMOVED***
		return fmt.Errorf("unexpected status code %d while pushing to %s", resp.StatusCode, pushURL)
	***REMOVED***
	return nil
***REMOVED***

func (r *registry) ServeHTTP(w http.ResponseWriter, req *http.Request) ***REMOVED***
	contentType := expfmt.Negotiate(req.Header)
	buf := r.getBuf()
	defer r.giveBuf(buf)
	writer, encoding := decorateWriter(req, buf)
	if err := r.writePB(expfmt.NewEncoder(writer, contentType)); err != nil ***REMOVED***
		if r.panicOnCollectError ***REMOVED***
			panic(err)
		***REMOVED***
		http.Error(w, "An error has occurred:\n\n"+err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	if closer, ok := writer.(io.Closer); ok ***REMOVED***
		closer.Close()
	***REMOVED***
	header := w.Header()
	header.Set(contentTypeHeader, string(contentType))
	header.Set(contentLengthHeader, fmt.Sprint(buf.Len()))
	if encoding != "" ***REMOVED***
		header.Set(contentEncodingHeader, encoding)
	***REMOVED***
	w.Write(buf.Bytes())
***REMOVED***

func (r *registry) writePB(encoder expfmt.Encoder) error ***REMOVED***
	var metricHashes map[uint64]struct***REMOVED******REMOVED***
	if r.collectChecksEnabled ***REMOVED***
		metricHashes = make(map[uint64]struct***REMOVED******REMOVED***)
	***REMOVED***
	metricChan := make(chan Metric, capMetricChan)
	wg := sync.WaitGroup***REMOVED******REMOVED***

	r.mtx.RLock()
	metricFamiliesByName := make(map[string]*dto.MetricFamily, len(r.dimHashesByName))

	// Scatter.
	// (Collectors could be complex and slow, so we call them all at once.)
	wg.Add(len(r.collectorsByID))
	go func() ***REMOVED***
		wg.Wait()
		close(metricChan)
	***REMOVED***()
	for _, collector := range r.collectorsByID ***REMOVED***
		go func(collector Collector) ***REMOVED***
			defer wg.Done()
			collector.Collect(metricChan)
		***REMOVED***(collector)
	***REMOVED***
	r.mtx.RUnlock()

	// Drain metricChan in case of premature return.
	defer func() ***REMOVED***
		for _ = range metricChan ***REMOVED***
		***REMOVED***
	***REMOVED***()

	// Gather.
	for metric := range metricChan ***REMOVED***
		// This could be done concurrently, too, but it required locking
		// of metricFamiliesByName (and of metricHashes if checks are
		// enabled). Most likely not worth it.
		desc := metric.Desc()
		metricFamily, ok := metricFamiliesByName[desc.fqName]
		if !ok ***REMOVED***
			metricFamily = r.getMetricFamily()
			defer r.giveMetricFamily(metricFamily)
			metricFamily.Name = proto.String(desc.fqName)
			metricFamily.Help = proto.String(desc.help)
			metricFamiliesByName[desc.fqName] = metricFamily
		***REMOVED***
		dtoMetric := r.getMetric()
		defer r.giveMetric(dtoMetric)
		if err := metric.Write(dtoMetric); err != nil ***REMOVED***
			// TODO: Consider different means of error reporting so
			// that a single erroneous metric could be skipped
			// instead of blowing up the whole collection.
			return fmt.Errorf("error collecting metric %v: %s", desc, err)
		***REMOVED***
		switch ***REMOVED***
		case metricFamily.Type != nil:
			// Type already set. We are good.
		case dtoMetric.Gauge != nil:
			metricFamily.Type = dto.MetricType_GAUGE.Enum()
		case dtoMetric.Counter != nil:
			metricFamily.Type = dto.MetricType_COUNTER.Enum()
		case dtoMetric.Summary != nil:
			metricFamily.Type = dto.MetricType_SUMMARY.Enum()
		case dtoMetric.Untyped != nil:
			metricFamily.Type = dto.MetricType_UNTYPED.Enum()
		case dtoMetric.Histogram != nil:
			metricFamily.Type = dto.MetricType_HISTOGRAM.Enum()
		default:
			return fmt.Errorf("empty metric collected: %s", dtoMetric)
		***REMOVED***
		if r.collectChecksEnabled ***REMOVED***
			if err := r.checkConsistency(metricFamily, dtoMetric, desc, metricHashes); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		metricFamily.Metric = append(metricFamily.Metric, dtoMetric)
	***REMOVED***

	if r.metricFamilyInjectionHook != nil ***REMOVED***
		for _, mf := range r.metricFamilyInjectionHook() ***REMOVED***
			existingMF, exists := metricFamiliesByName[mf.GetName()]
			if !exists ***REMOVED***
				metricFamiliesByName[mf.GetName()] = mf
				if r.collectChecksEnabled ***REMOVED***
					for _, m := range mf.Metric ***REMOVED***
						if err := r.checkConsistency(mf, m, nil, metricHashes); err != nil ***REMOVED***
							return err
						***REMOVED***
					***REMOVED***
				***REMOVED***
				continue
			***REMOVED***
			for _, m := range mf.Metric ***REMOVED***
				if r.collectChecksEnabled ***REMOVED***
					if err := r.checkConsistency(existingMF, m, nil, metricHashes); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				existingMF.Metric = append(existingMF.Metric, m)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Now that MetricFamilies are all set, sort their Metrics
	// lexicographically by their label values.
	for _, mf := range metricFamiliesByName ***REMOVED***
		sort.Sort(metricSorter(mf.Metric))
	***REMOVED***

	// Write out MetricFamilies sorted by their name.
	names := make([]string, 0, len(metricFamiliesByName))
	for name := range metricFamiliesByName ***REMOVED***
		names = append(names, name)
	***REMOVED***
	sort.Strings(names)

	for _, name := range names ***REMOVED***
		if err := encoder.Encode(metricFamiliesByName[name]); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *registry) checkConsistency(metricFamily *dto.MetricFamily, dtoMetric *dto.Metric, desc *Desc, metricHashes map[uint64]struct***REMOVED******REMOVED***) error ***REMOVED***

	// Type consistency with metric family.
	if metricFamily.GetType() == dto.MetricType_GAUGE && dtoMetric.Gauge == nil ||
		metricFamily.GetType() == dto.MetricType_COUNTER && dtoMetric.Counter == nil ||
		metricFamily.GetType() == dto.MetricType_SUMMARY && dtoMetric.Summary == nil ||
		metricFamily.GetType() == dto.MetricType_HISTOGRAM && dtoMetric.Histogram == nil ||
		metricFamily.GetType() == dto.MetricType_UNTYPED && dtoMetric.Untyped == nil ***REMOVED***
		return fmt.Errorf(
			"collected metric %s %s is not a %s",
			metricFamily.GetName(), dtoMetric, metricFamily.GetType(),
		)
	***REMOVED***

	// Is the metric unique (i.e. no other metric with the same name and the same label values)?
	h := hashNew()
	h = hashAdd(h, metricFamily.GetName())
	h = hashAddByte(h, separatorByte)
	// Make sure label pairs are sorted. We depend on it for the consistency
	// check. Label pairs must be sorted by contract. But the point of this
	// method is to check for contract violations. So we better do the sort
	// now.
	sort.Sort(LabelPairSorter(dtoMetric.Label))
	for _, lp := range dtoMetric.Label ***REMOVED***
		h = hashAdd(h, lp.GetValue())
		h = hashAddByte(h, separatorByte)
	***REMOVED***
	if _, exists := metricHashes[h]; exists ***REMOVED***
		return fmt.Errorf(
			"collected metric %s %s was collected before with the same name and label values",
			metricFamily.GetName(), dtoMetric,
		)
	***REMOVED***
	metricHashes[h] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	if desc == nil ***REMOVED***
		return nil // Nothing left to check if we have no desc.
	***REMOVED***

	// Desc consistency with metric family.
	if metricFamily.GetName() != desc.fqName ***REMOVED***
		return fmt.Errorf(
			"collected metric %s %s has name %q but should have %q",
			metricFamily.GetName(), dtoMetric, metricFamily.GetName(), desc.fqName,
		)
	***REMOVED***
	if metricFamily.GetHelp() != desc.help ***REMOVED***
		return fmt.Errorf(
			"collected metric %s %s has help %q but should have %q",
			metricFamily.GetName(), dtoMetric, metricFamily.GetHelp(), desc.help,
		)
	***REMOVED***

	// Is the desc consistent with the content of the metric?
	lpsFromDesc := make([]*dto.LabelPair, 0, len(dtoMetric.Label))
	lpsFromDesc = append(lpsFromDesc, desc.constLabelPairs...)
	for _, l := range desc.variableLabels ***REMOVED***
		lpsFromDesc = append(lpsFromDesc, &dto.LabelPair***REMOVED***
			Name: proto.String(l),
		***REMOVED***)
	***REMOVED***
	if len(lpsFromDesc) != len(dtoMetric.Label) ***REMOVED***
		return fmt.Errorf(
			"labels in collected metric %s %s are inconsistent with descriptor %s",
			metricFamily.GetName(), dtoMetric, desc,
		)
	***REMOVED***
	sort.Sort(LabelPairSorter(lpsFromDesc))
	for i, lpFromDesc := range lpsFromDesc ***REMOVED***
		lpFromMetric := dtoMetric.Label[i]
		if lpFromDesc.GetName() != lpFromMetric.GetName() ||
			lpFromDesc.Value != nil && lpFromDesc.GetValue() != lpFromMetric.GetValue() ***REMOVED***
			return fmt.Errorf(
				"labels in collected metric %s %s are inconsistent with descriptor %s",
				metricFamily.GetName(), dtoMetric, desc,
			)
		***REMOVED***
	***REMOVED***

	r.mtx.RLock() // Remaining checks need the read lock.
	defer r.mtx.RUnlock()

	// Is the desc registered?
	if _, exist := r.descIDs[desc.id]; !exist ***REMOVED***
		return fmt.Errorf(
			"collected metric %s %s with unregistered descriptor %s",
			metricFamily.GetName(), dtoMetric, desc,
		)
	***REMOVED***

	return nil
***REMOVED***

func (r *registry) getBuf() *bytes.Buffer ***REMOVED***
	select ***REMOVED***
	case buf := <-r.bufPool:
		return buf
	default:
		return &bytes.Buffer***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

func (r *registry) giveBuf(buf *bytes.Buffer) ***REMOVED***
	buf.Reset()
	select ***REMOVED***
	case r.bufPool <- buf:
	default:
	***REMOVED***
***REMOVED***

func (r *registry) getMetricFamily() *dto.MetricFamily ***REMOVED***
	select ***REMOVED***
	case mf := <-r.metricFamilyPool:
		return mf
	default:
		return &dto.MetricFamily***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

func (r *registry) giveMetricFamily(mf *dto.MetricFamily) ***REMOVED***
	mf.Reset()
	select ***REMOVED***
	case r.metricFamilyPool <- mf:
	default:
	***REMOVED***
***REMOVED***

func (r *registry) getMetric() *dto.Metric ***REMOVED***
	select ***REMOVED***
	case m := <-r.metricPool:
		return m
	default:
		return &dto.Metric***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

func (r *registry) giveMetric(m *dto.Metric) ***REMOVED***
	m.Reset()
	select ***REMOVED***
	case r.metricPool <- m:
	default:
	***REMOVED***
***REMOVED***

func newRegistry() *registry ***REMOVED***
	return &registry***REMOVED***
		collectorsByID:   map[uint64]Collector***REMOVED******REMOVED***,
		descIDs:          map[uint64]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		dimHashesByName:  map[string]uint64***REMOVED******REMOVED***,
		bufPool:          make(chan *bytes.Buffer, numBufs),
		metricFamilyPool: make(chan *dto.MetricFamily, numMetricFamilies),
		metricPool:       make(chan *dto.Metric, numMetrics),
	***REMOVED***
***REMOVED***

func newDefaultRegistry() *registry ***REMOVED***
	r := newRegistry()
	r.Register(NewProcessCollector(os.Getpid(), ""))
	r.Register(NewGoCollector())
	return r
***REMOVED***

// decorateWriter wraps a writer to handle gzip compression if requested.  It
// returns the decorated writer and the appropriate "Content-Encoding" header
// (which is empty if no compression is enabled).
func decorateWriter(request *http.Request, writer io.Writer) (io.Writer, string) ***REMOVED***
	header := request.Header.Get(acceptEncodingHeader)
	parts := strings.Split(header, ",")
	for _, part := range parts ***REMOVED***
		part := strings.TrimSpace(part)
		if part == "gzip" || strings.HasPrefix(part, "gzip;") ***REMOVED***
			return gzip.NewWriter(writer), "gzip"
		***REMOVED***
	***REMOVED***
	return writer, ""
***REMOVED***

type metricSorter []*dto.Metric

func (s metricSorter) Len() int ***REMOVED***
	return len(s)
***REMOVED***

func (s metricSorter) Swap(i, j int) ***REMOVED***
	s[i], s[j] = s[j], s[i]
***REMOVED***

func (s metricSorter) Less(i, j int) bool ***REMOVED***
	if len(s[i].Label) != len(s[j].Label) ***REMOVED***
		// This should not happen. The metrics are
		// inconsistent. However, we have to deal with the fact, as
		// people might use custom collectors or metric family injection
		// to create inconsistent metrics. So let's simply compare the
		// number of labels in this case. That will still yield
		// reproducible sorting.
		return len(s[i].Label) < len(s[j].Label)
	***REMOVED***
	for n, lp := range s[i].Label ***REMOVED***
		vi := lp.GetValue()
		vj := s[j].Label[n].GetValue()
		if vi != vj ***REMOVED***
			return vi < vj
		***REMOVED***
	***REMOVED***

	// We should never arrive here. Multiple metrics with the same
	// label set in the same scrape will lead to undefined ingestion
	// behavior. However, as above, we have to provide stable sorting
	// here, even for inconsistent metrics. So sort equal metrics
	// by their timestamp, with missing timestamps (implying "now")
	// coming last.
	if s[i].TimestampMs == nil ***REMOVED***
		return false
	***REMOVED***
	if s[j].TimestampMs == nil ***REMOVED***
		return true
	***REMOVED***
	return s[i].GetTimestampMs() < s[j].GetTimestampMs()
***REMOVED***
