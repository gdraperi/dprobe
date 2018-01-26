// Copyright 2015 The Prometheus Authors
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

package expfmt

import (
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"

	dto "github.com/prometheus/client_model/go"

	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	"github.com/prometheus/common/model"
)

// Decoder types decode an input stream into metric families.
type Decoder interface ***REMOVED***
	Decode(*dto.MetricFamily) error
***REMOVED***

type DecodeOptions struct ***REMOVED***
	// Timestamp is added to each value from the stream that has no explicit timestamp set.
	Timestamp model.Time
***REMOVED***

// ResponseFormat extracts the correct format from a HTTP response header.
// If no matching format can be found FormatUnknown is returned.
func ResponseFormat(h http.Header) Format ***REMOVED***
	ct := h.Get(hdrContentType)

	mediatype, params, err := mime.ParseMediaType(ct)
	if err != nil ***REMOVED***
		return FmtUnknown
	***REMOVED***

	const textType = "text/plain"

	switch mediatype ***REMOVED***
	case ProtoType:
		if p, ok := params["proto"]; ok && p != ProtoProtocol ***REMOVED***
			return FmtUnknown
		***REMOVED***
		if e, ok := params["encoding"]; ok && e != "delimited" ***REMOVED***
			return FmtUnknown
		***REMOVED***
		return FmtProtoDelim

	case textType:
		if v, ok := params["version"]; ok && v != TextVersion ***REMOVED***
			return FmtUnknown
		***REMOVED***
		return FmtText
	***REMOVED***

	return FmtUnknown
***REMOVED***

// NewDecoder returns a new decoder based on the given input format.
// If the input format does not imply otherwise, a text format decoder is returned.
func NewDecoder(r io.Reader, format Format) Decoder ***REMOVED***
	switch format ***REMOVED***
	case FmtProtoDelim:
		return &protoDecoder***REMOVED***r: r***REMOVED***
	***REMOVED***
	return &textDecoder***REMOVED***r: r***REMOVED***
***REMOVED***

// protoDecoder implements the Decoder interface for protocol buffers.
type protoDecoder struct ***REMOVED***
	r io.Reader
***REMOVED***

// Decode implements the Decoder interface.
func (d *protoDecoder) Decode(v *dto.MetricFamily) error ***REMOVED***
	_, err := pbutil.ReadDelimited(d.r, v)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !model.IsValidMetricName(model.LabelValue(v.GetName())) ***REMOVED***
		return fmt.Errorf("invalid metric name %q", v.GetName())
	***REMOVED***
	for _, m := range v.GetMetric() ***REMOVED***
		if m == nil ***REMOVED***
			continue
		***REMOVED***
		for _, l := range m.GetLabel() ***REMOVED***
			if l == nil ***REMOVED***
				continue
			***REMOVED***
			if !model.LabelValue(l.GetValue()).IsValid() ***REMOVED***
				return fmt.Errorf("invalid label value %q", l.GetValue())
			***REMOVED***
			if !model.LabelName(l.GetName()).IsValid() ***REMOVED***
				return fmt.Errorf("invalid label name %q", l.GetName())
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// textDecoder implements the Decoder interface for the text protocol.
type textDecoder struct ***REMOVED***
	r    io.Reader
	p    TextParser
	fams []*dto.MetricFamily
***REMOVED***

// Decode implements the Decoder interface.
func (d *textDecoder) Decode(v *dto.MetricFamily) error ***REMOVED***
	// TODO(fabxc): Wrap this as a line reader to make streaming safer.
	if len(d.fams) == 0 ***REMOVED***
		// No cached metric families, read everything and parse metrics.
		fams, err := d.p.TextToMetricFamilies(d.r)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if len(fams) == 0 ***REMOVED***
			return io.EOF
		***REMOVED***
		d.fams = make([]*dto.MetricFamily, 0, len(fams))
		for _, f := range fams ***REMOVED***
			d.fams = append(d.fams, f)
		***REMOVED***
	***REMOVED***

	*v = *d.fams[0]
	d.fams = d.fams[1:]

	return nil
***REMOVED***

type SampleDecoder struct ***REMOVED***
	Dec  Decoder
	Opts *DecodeOptions

	f dto.MetricFamily
***REMOVED***

func (sd *SampleDecoder) Decode(s *model.Vector) error ***REMOVED***
	if err := sd.Dec.Decode(&sd.f); err != nil ***REMOVED***
		return err
	***REMOVED***
	*s = extractSamples(&sd.f, sd.Opts)
	return nil
***REMOVED***

// Extract samples builds a slice of samples from the provided metric families.
func ExtractSamples(o *DecodeOptions, fams ...*dto.MetricFamily) model.Vector ***REMOVED***
	var all model.Vector
	for _, f := range fams ***REMOVED***
		all = append(all, extractSamples(f, o)...)
	***REMOVED***
	return all
***REMOVED***

func extractSamples(f *dto.MetricFamily, o *DecodeOptions) model.Vector ***REMOVED***
	switch f.GetType() ***REMOVED***
	case dto.MetricType_COUNTER:
		return extractCounter(o, f)
	case dto.MetricType_GAUGE:
		return extractGauge(o, f)
	case dto.MetricType_SUMMARY:
		return extractSummary(o, f)
	case dto.MetricType_UNTYPED:
		return extractUntyped(o, f)
	case dto.MetricType_HISTOGRAM:
		return extractHistogram(o, f)
	***REMOVED***
	panic("expfmt.extractSamples: unknown metric family type")
***REMOVED***

func extractCounter(o *DecodeOptions, f *dto.MetricFamily) model.Vector ***REMOVED***
	samples := make(model.Vector, 0, len(f.Metric))

	for _, m := range f.Metric ***REMOVED***
		if m.Counter == nil ***REMOVED***
			continue
		***REMOVED***

		lset := make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label ***REMOVED***
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		***REMOVED***
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName())

		smpl := &model.Sample***REMOVED***
			Metric: model.Metric(lset),
			Value:  model.SampleValue(m.Counter.GetValue()),
		***REMOVED***

		if m.TimestampMs != nil ***REMOVED***
			smpl.Timestamp = model.TimeFromUnixNano(*m.TimestampMs * 1000000)
		***REMOVED*** else ***REMOVED***
			smpl.Timestamp = o.Timestamp
		***REMOVED***

		samples = append(samples, smpl)
	***REMOVED***

	return samples
***REMOVED***

func extractGauge(o *DecodeOptions, f *dto.MetricFamily) model.Vector ***REMOVED***
	samples := make(model.Vector, 0, len(f.Metric))

	for _, m := range f.Metric ***REMOVED***
		if m.Gauge == nil ***REMOVED***
			continue
		***REMOVED***

		lset := make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label ***REMOVED***
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		***REMOVED***
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName())

		smpl := &model.Sample***REMOVED***
			Metric: model.Metric(lset),
			Value:  model.SampleValue(m.Gauge.GetValue()),
		***REMOVED***

		if m.TimestampMs != nil ***REMOVED***
			smpl.Timestamp = model.TimeFromUnixNano(*m.TimestampMs * 1000000)
		***REMOVED*** else ***REMOVED***
			smpl.Timestamp = o.Timestamp
		***REMOVED***

		samples = append(samples, smpl)
	***REMOVED***

	return samples
***REMOVED***

func extractUntyped(o *DecodeOptions, f *dto.MetricFamily) model.Vector ***REMOVED***
	samples := make(model.Vector, 0, len(f.Metric))

	for _, m := range f.Metric ***REMOVED***
		if m.Untyped == nil ***REMOVED***
			continue
		***REMOVED***

		lset := make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label ***REMOVED***
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		***REMOVED***
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName())

		smpl := &model.Sample***REMOVED***
			Metric: model.Metric(lset),
			Value:  model.SampleValue(m.Untyped.GetValue()),
		***REMOVED***

		if m.TimestampMs != nil ***REMOVED***
			smpl.Timestamp = model.TimeFromUnixNano(*m.TimestampMs * 1000000)
		***REMOVED*** else ***REMOVED***
			smpl.Timestamp = o.Timestamp
		***REMOVED***

		samples = append(samples, smpl)
	***REMOVED***

	return samples
***REMOVED***

func extractSummary(o *DecodeOptions, f *dto.MetricFamily) model.Vector ***REMOVED***
	samples := make(model.Vector, 0, len(f.Metric))

	for _, m := range f.Metric ***REMOVED***
		if m.Summary == nil ***REMOVED***
			continue
		***REMOVED***

		timestamp := o.Timestamp
		if m.TimestampMs != nil ***REMOVED***
			timestamp = model.TimeFromUnixNano(*m.TimestampMs * 1000000)
		***REMOVED***

		for _, q := range m.Summary.Quantile ***REMOVED***
			lset := make(model.LabelSet, len(m.Label)+2)
			for _, p := range m.Label ***REMOVED***
				lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
			***REMOVED***
			// BUG(matt): Update other names to "quantile".
			lset[model.LabelName(model.QuantileLabel)] = model.LabelValue(fmt.Sprint(q.GetQuantile()))
			lset[model.MetricNameLabel] = model.LabelValue(f.GetName())

			samples = append(samples, &model.Sample***REMOVED***
				Metric:    model.Metric(lset),
				Value:     model.SampleValue(q.GetValue()),
				Timestamp: timestamp,
			***REMOVED***)
		***REMOVED***

		lset := make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label ***REMOVED***
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		***REMOVED***
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_sum")

		samples = append(samples, &model.Sample***REMOVED***
			Metric:    model.Metric(lset),
			Value:     model.SampleValue(m.Summary.GetSampleSum()),
			Timestamp: timestamp,
		***REMOVED***)

		lset = make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label ***REMOVED***
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		***REMOVED***
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_count")

		samples = append(samples, &model.Sample***REMOVED***
			Metric:    model.Metric(lset),
			Value:     model.SampleValue(m.Summary.GetSampleCount()),
			Timestamp: timestamp,
		***REMOVED***)
	***REMOVED***

	return samples
***REMOVED***

func extractHistogram(o *DecodeOptions, f *dto.MetricFamily) model.Vector ***REMOVED***
	samples := make(model.Vector, 0, len(f.Metric))

	for _, m := range f.Metric ***REMOVED***
		if m.Histogram == nil ***REMOVED***
			continue
		***REMOVED***

		timestamp := o.Timestamp
		if m.TimestampMs != nil ***REMOVED***
			timestamp = model.TimeFromUnixNano(*m.TimestampMs * 1000000)
		***REMOVED***

		infSeen := false

		for _, q := range m.Histogram.Bucket ***REMOVED***
			lset := make(model.LabelSet, len(m.Label)+2)
			for _, p := range m.Label ***REMOVED***
				lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
			***REMOVED***
			lset[model.LabelName(model.BucketLabel)] = model.LabelValue(fmt.Sprint(q.GetUpperBound()))
			lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_bucket")

			if math.IsInf(q.GetUpperBound(), +1) ***REMOVED***
				infSeen = true
			***REMOVED***

			samples = append(samples, &model.Sample***REMOVED***
				Metric:    model.Metric(lset),
				Value:     model.SampleValue(q.GetCumulativeCount()),
				Timestamp: timestamp,
			***REMOVED***)
		***REMOVED***

		lset := make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label ***REMOVED***
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		***REMOVED***
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_sum")

		samples = append(samples, &model.Sample***REMOVED***
			Metric:    model.Metric(lset),
			Value:     model.SampleValue(m.Histogram.GetSampleSum()),
			Timestamp: timestamp,
		***REMOVED***)

		lset = make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label ***REMOVED***
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		***REMOVED***
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_count")

		count := &model.Sample***REMOVED***
			Metric:    model.Metric(lset),
			Value:     model.SampleValue(m.Histogram.GetSampleCount()),
			Timestamp: timestamp,
		***REMOVED***
		samples = append(samples, count)

		if !infSeen ***REMOVED***
			// Append an infinity bucket sample.
			lset := make(model.LabelSet, len(m.Label)+2)
			for _, p := range m.Label ***REMOVED***
				lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
			***REMOVED***
			lset[model.LabelName(model.BucketLabel)] = model.LabelValue("+Inf")
			lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_bucket")

			samples = append(samples, &model.Sample***REMOVED***
				Metric:    model.Metric(lset),
				Value:     count.Value,
				Timestamp: timestamp,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	return samples
***REMOVED***
