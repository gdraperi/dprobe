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

package prometheus

import (
	"encoding/json"
	"expvar"
)

// ExpvarCollector collects metrics from the expvar interface. It provides a
// quick way to expose numeric values that are already exported via expvar as
// Prometheus metrics. Note that the data models of expvar and Prometheus are
// fundamentally different, and that the ExpvarCollector is inherently
// slow. Thus, the ExpvarCollector is probably great for experiments and
// prototying, but you should seriously consider a more direct implementation of
// Prometheus metrics for monitoring production systems.
//
// Use NewExpvarCollector to create new instances.
type ExpvarCollector struct ***REMOVED***
	exports map[string]*Desc
***REMOVED***

// NewExpvarCollector returns a newly allocated ExpvarCollector that still has
// to be registered with the Prometheus registry.
//
// The exports map has the following meaning:
//
// The keys in the map correspond to expvar keys, i.e. for every expvar key you
// want to export as Prometheus metric, you need an entry in the exports
// map. The descriptor mapped to each key describes how to export the expvar
// value. It defines the name and the help string of the Prometheus metric
// proxying the expvar value. The type will always be Untyped.
//
// For descriptors without variable labels, the expvar value must be a number or
// a bool. The number is then directly exported as the Prometheus sample
// value. (For a bool, 'false' translates to 0 and 'true' to 1). Expvar values
// that are not numbers or bools are silently ignored.
//
// If the descriptor has one variable label, the expvar value must be an expvar
// map. The keys in the expvar map become the various values of the one
// Prometheus label. The values in the expvar map must be numbers or bools again
// as above.
//
// For descriptors with more than one variable label, the expvar must be a
// nested expvar map, i.e. where the values of the topmost map are maps again
// etc. until a depth is reached that corresponds to the number of labels. The
// leaves of that structure must be numbers or bools as above to serve as the
// sample values.
//
// Anything that does not fit into the scheme above is silently ignored.
func NewExpvarCollector(exports map[string]*Desc) *ExpvarCollector ***REMOVED***
	return &ExpvarCollector***REMOVED***
		exports: exports,
	***REMOVED***
***REMOVED***

// Describe implements Collector.
func (e *ExpvarCollector) Describe(ch chan<- *Desc) ***REMOVED***
	for _, desc := range e.exports ***REMOVED***
		ch <- desc
	***REMOVED***
***REMOVED***

// Collect implements Collector.
func (e *ExpvarCollector) Collect(ch chan<- Metric) ***REMOVED***
	for name, desc := range e.exports ***REMOVED***
		var m Metric
		expVar := expvar.Get(name)
		if expVar == nil ***REMOVED***
			continue
		***REMOVED***
		var v interface***REMOVED******REMOVED***
		labels := make([]string, len(desc.variableLabels))
		if err := json.Unmarshal([]byte(expVar.String()), &v); err != nil ***REMOVED***
			ch <- NewInvalidMetric(desc, err)
			continue
		***REMOVED***
		var processValue func(v interface***REMOVED******REMOVED***, i int)
		processValue = func(v interface***REMOVED******REMOVED***, i int) ***REMOVED***
			if i >= len(labels) ***REMOVED***
				copiedLabels := append(make([]string, 0, len(labels)), labels...)
				switch v := v.(type) ***REMOVED***
				case float64:
					m = MustNewConstMetric(desc, UntypedValue, v, copiedLabels...)
				case bool:
					if v ***REMOVED***
						m = MustNewConstMetric(desc, UntypedValue, 1, copiedLabels...)
					***REMOVED*** else ***REMOVED***
						m = MustNewConstMetric(desc, UntypedValue, 0, copiedLabels...)
					***REMOVED***
				default:
					return
				***REMOVED***
				ch <- m
				return
			***REMOVED***
			vm, ok := v.(map[string]interface***REMOVED******REMOVED***)
			if !ok ***REMOVED***
				return
			***REMOVED***
			for lv, val := range vm ***REMOVED***
				labels[i] = lv
				processValue(val, i+1)
			***REMOVED***
		***REMOVED***
		processValue(v, 0)
	***REMOVED***
***REMOVED***
