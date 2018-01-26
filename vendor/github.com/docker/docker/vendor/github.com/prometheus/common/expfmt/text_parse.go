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

package expfmt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	dto "github.com/prometheus/client_model/go"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/common/model"
)

// A stateFn is a function that represents a state in a state machine. By
// executing it, the state is progressed to the next state. The stateFn returns
// another stateFn, which represents the new state. The end state is represented
// by nil.
type stateFn func() stateFn

// ParseError signals errors while parsing the simple and flat text-based
// exchange format.
type ParseError struct ***REMOVED***
	Line int
	Msg  string
***REMOVED***

// Error implements the error interface.
func (e ParseError) Error() string ***REMOVED***
	return fmt.Sprintf("text format parsing error in line %d: %s", e.Line, e.Msg)
***REMOVED***

// TextParser is used to parse the simple and flat text-based exchange format. Its
// nil value is ready to use.
type TextParser struct ***REMOVED***
	metricFamiliesByName map[string]*dto.MetricFamily
	buf                  *bufio.Reader // Where the parsed input is read through.
	err                  error         // Most recent error.
	lineCount            int           // Tracks the line count for error messages.
	currentByte          byte          // The most recent byte read.
	currentToken         bytes.Buffer  // Re-used each time a token has to be gathered from multiple bytes.
	currentMF            *dto.MetricFamily
	currentMetric        *dto.Metric
	currentLabelPair     *dto.LabelPair

	// The remaining member variables are only used for summaries/histograms.
	currentLabels map[string]string // All labels including '__name__' but excluding 'quantile'/'le'
	// Summary specific.
	summaries       map[uint64]*dto.Metric // Key is created with LabelsToSignature.
	currentQuantile float64
	// Histogram specific.
	histograms    map[uint64]*dto.Metric // Key is created with LabelsToSignature.
	currentBucket float64
	// These tell us if the currently processed line ends on '_count' or
	// '_sum' respectively and belong to a summary/histogram, representing the sample
	// count and sum of that summary/histogram.
	currentIsSummaryCount, currentIsSummarySum     bool
	currentIsHistogramCount, currentIsHistogramSum bool
***REMOVED***

// TextToMetricFamilies reads 'in' as the simple and flat text-based exchange
// format and creates MetricFamily proto messages. It returns the MetricFamily
// proto messages in a map where the metric names are the keys, along with any
// error encountered.
//
// If the input contains duplicate metrics (i.e. lines with the same metric name
// and exactly the same label set), the resulting MetricFamily will contain
// duplicate Metric proto messages. Similar is true for duplicate label
// names. Checks for duplicates have to be performed separately, if required.
// Also note that neither the metrics within each MetricFamily are sorted nor
// the label pairs within each Metric. Sorting is not required for the most
// frequent use of this method, which is sample ingestion in the Prometheus
// server. However, for presentation purposes, you might want to sort the
// metrics, and in some cases, you must sort the labels, e.g. for consumption by
// the metric family injection hook of the Prometheus registry.
//
// Summaries and histograms are rather special beasts. You would probably not
// use them in the simple text format anyway. This method can deal with
// summaries and histograms if they are presented in exactly the way the
// text.Create function creates them.
//
// This method must not be called concurrently. If you want to parse different
// input concurrently, instantiate a separate Parser for each goroutine.
func (p *TextParser) TextToMetricFamilies(in io.Reader) (map[string]*dto.MetricFamily, error) ***REMOVED***
	p.reset(in)
	for nextState := p.startOfLine; nextState != nil; nextState = nextState() ***REMOVED***
		// Magic happens here...
	***REMOVED***
	// Get rid of empty metric families.
	for k, mf := range p.metricFamiliesByName ***REMOVED***
		if len(mf.GetMetric()) == 0 ***REMOVED***
			delete(p.metricFamiliesByName, k)
		***REMOVED***
	***REMOVED***
	// If p.err is io.EOF now, we have run into a premature end of the input
	// stream. Turn this error into something nicer and more
	// meaningful. (io.EOF is often used as a signal for the legitimate end
	// of an input stream.)
	if p.err == io.EOF ***REMOVED***
		p.parseError("unexpected end of input stream")
	***REMOVED***
	return p.metricFamiliesByName, p.err
***REMOVED***

func (p *TextParser) reset(in io.Reader) ***REMOVED***
	p.metricFamiliesByName = map[string]*dto.MetricFamily***REMOVED******REMOVED***
	if p.buf == nil ***REMOVED***
		p.buf = bufio.NewReader(in)
	***REMOVED*** else ***REMOVED***
		p.buf.Reset(in)
	***REMOVED***
	p.err = nil
	p.lineCount = 0
	if p.summaries == nil || len(p.summaries) > 0 ***REMOVED***
		p.summaries = map[uint64]*dto.Metric***REMOVED******REMOVED***
	***REMOVED***
	if p.histograms == nil || len(p.histograms) > 0 ***REMOVED***
		p.histograms = map[uint64]*dto.Metric***REMOVED******REMOVED***
	***REMOVED***
	p.currentQuantile = math.NaN()
	p.currentBucket = math.NaN()
***REMOVED***

// startOfLine represents the state where the next byte read from p.buf is the
// start of a line (or whitespace leading up to it).
func (p *TextParser) startOfLine() stateFn ***REMOVED***
	p.lineCount++
	if p.skipBlankTab(); p.err != nil ***REMOVED***
		// End of input reached. This is the only case where
		// that is not an error but a signal that we are done.
		p.err = nil
		return nil
	***REMOVED***
	switch p.currentByte ***REMOVED***
	case '#':
		return p.startComment
	case '\n':
		return p.startOfLine // Empty line, start the next one.
	***REMOVED***
	return p.readingMetricName
***REMOVED***

// startComment represents the state where the next byte read from p.buf is the
// start of a comment (or whitespace leading up to it).
func (p *TextParser) startComment() stateFn ***REMOVED***
	if p.skipBlankTab(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	if p.currentByte == '\n' ***REMOVED***
		return p.startOfLine
	***REMOVED***
	if p.readTokenUntilWhitespace(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	// If we have hit the end of line already, there is nothing left
	// to do. This is not considered a syntax error.
	if p.currentByte == '\n' ***REMOVED***
		return p.startOfLine
	***REMOVED***
	keyword := p.currentToken.String()
	if keyword != "HELP" && keyword != "TYPE" ***REMOVED***
		// Generic comment, ignore by fast forwarding to end of line.
		for p.currentByte != '\n' ***REMOVED***
			if p.currentByte, p.err = p.buf.ReadByte(); p.err != nil ***REMOVED***
				return nil // Unexpected end of input.
			***REMOVED***
		***REMOVED***
		return p.startOfLine
	***REMOVED***
	// There is something. Next has to be a metric name.
	if p.skipBlankTab(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	if p.readTokenAsMetricName(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	if p.currentByte == '\n' ***REMOVED***
		// At the end of the line already.
		// Again, this is not considered a syntax error.
		return p.startOfLine
	***REMOVED***
	if !isBlankOrTab(p.currentByte) ***REMOVED***
		p.parseError("invalid metric name in comment")
		return nil
	***REMOVED***
	p.setOrCreateCurrentMF()
	if p.skipBlankTab(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	if p.currentByte == '\n' ***REMOVED***
		// At the end of the line already.
		// Again, this is not considered a syntax error.
		return p.startOfLine
	***REMOVED***
	switch keyword ***REMOVED***
	case "HELP":
		return p.readingHelp
	case "TYPE":
		return p.readingType
	***REMOVED***
	panic(fmt.Sprintf("code error: unexpected keyword %q", keyword))
***REMOVED***

// readingMetricName represents the state where the last byte read (now in
// p.currentByte) is the first byte of a metric name.
func (p *TextParser) readingMetricName() stateFn ***REMOVED***
	if p.readTokenAsMetricName(); p.err != nil ***REMOVED***
		return nil
	***REMOVED***
	if p.currentToken.Len() == 0 ***REMOVED***
		p.parseError("invalid metric name")
		return nil
	***REMOVED***
	p.setOrCreateCurrentMF()
	// Now is the time to fix the type if it hasn't happened yet.
	if p.currentMF.Type == nil ***REMOVED***
		p.currentMF.Type = dto.MetricType_UNTYPED.Enum()
	***REMOVED***
	p.currentMetric = &dto.Metric***REMOVED******REMOVED***
	// Do not append the newly created currentMetric to
	// currentMF.Metric right now. First wait if this is a summary,
	// and the metric exists already, which we can only know after
	// having read all the labels.
	if p.skipBlankTabIfCurrentBlankTab(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	return p.readingLabels
***REMOVED***

// readingLabels represents the state where the last byte read (now in
// p.currentByte) is either the first byte of the label set (i.e. a '***REMOVED***'), or the
// first byte of the value (otherwise).
func (p *TextParser) readingLabels() stateFn ***REMOVED***
	// Summaries/histograms are special. We have to reset the
	// currentLabels map, currentQuantile and currentBucket before starting to
	// read labels.
	if p.currentMF.GetType() == dto.MetricType_SUMMARY || p.currentMF.GetType() == dto.MetricType_HISTOGRAM ***REMOVED***
		p.currentLabels = map[string]string***REMOVED******REMOVED***
		p.currentLabels[string(model.MetricNameLabel)] = p.currentMF.GetName()
		p.currentQuantile = math.NaN()
		p.currentBucket = math.NaN()
	***REMOVED***
	if p.currentByte != '***REMOVED***' ***REMOVED***
		return p.readingValue
	***REMOVED***
	return p.startLabelName
***REMOVED***

// startLabelName represents the state where the next byte read from p.buf is
// the start of a label name (or whitespace leading up to it).
func (p *TextParser) startLabelName() stateFn ***REMOVED***
	if p.skipBlankTab(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	if p.currentByte == '***REMOVED***' ***REMOVED***
		if p.skipBlankTab(); p.err != nil ***REMOVED***
			return nil // Unexpected end of input.
		***REMOVED***
		return p.readingValue
	***REMOVED***
	if p.readTokenAsLabelName(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	if p.currentToken.Len() == 0 ***REMOVED***
		p.parseError(fmt.Sprintf("invalid label name for metric %q", p.currentMF.GetName()))
		return nil
	***REMOVED***
	p.currentLabelPair = &dto.LabelPair***REMOVED***Name: proto.String(p.currentToken.String())***REMOVED***
	if p.currentLabelPair.GetName() == string(model.MetricNameLabel) ***REMOVED***
		p.parseError(fmt.Sprintf("label name %q is reserved", model.MetricNameLabel))
		return nil
	***REMOVED***
	// Special summary/histogram treatment. Don't add 'quantile' and 'le'
	// labels to 'real' labels.
	if !(p.currentMF.GetType() == dto.MetricType_SUMMARY && p.currentLabelPair.GetName() == model.QuantileLabel) &&
		!(p.currentMF.GetType() == dto.MetricType_HISTOGRAM && p.currentLabelPair.GetName() == model.BucketLabel) ***REMOVED***
		p.currentMetric.Label = append(p.currentMetric.Label, p.currentLabelPair)
	***REMOVED***
	if p.skipBlankTabIfCurrentBlankTab(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	if p.currentByte != '=' ***REMOVED***
		p.parseError(fmt.Sprintf("expected '=' after label name, found %q", p.currentByte))
		return nil
	***REMOVED***
	return p.startLabelValue
***REMOVED***

// startLabelValue represents the state where the next byte read from p.buf is
// the start of a (quoted) label value (or whitespace leading up to it).
func (p *TextParser) startLabelValue() stateFn ***REMOVED***
	if p.skipBlankTab(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	if p.currentByte != '"' ***REMOVED***
		p.parseError(fmt.Sprintf("expected '\"' at start of label value, found %q", p.currentByte))
		return nil
	***REMOVED***
	if p.readTokenAsLabelValue(); p.err != nil ***REMOVED***
		return nil
	***REMOVED***
	p.currentLabelPair.Value = proto.String(p.currentToken.String())
	// Special treatment of summaries:
	// - Quantile labels are special, will result in dto.Quantile later.
	// - Other labels have to be added to currentLabels for signature calculation.
	if p.currentMF.GetType() == dto.MetricType_SUMMARY ***REMOVED***
		if p.currentLabelPair.GetName() == model.QuantileLabel ***REMOVED***
			if p.currentQuantile, p.err = strconv.ParseFloat(p.currentLabelPair.GetValue(), 64); p.err != nil ***REMOVED***
				// Create a more helpful error message.
				p.parseError(fmt.Sprintf("expected float as value for 'quantile' label, got %q", p.currentLabelPair.GetValue()))
				return nil
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			p.currentLabels[p.currentLabelPair.GetName()] = p.currentLabelPair.GetValue()
		***REMOVED***
	***REMOVED***
	// Similar special treatment of histograms.
	if p.currentMF.GetType() == dto.MetricType_HISTOGRAM ***REMOVED***
		if p.currentLabelPair.GetName() == model.BucketLabel ***REMOVED***
			if p.currentBucket, p.err = strconv.ParseFloat(p.currentLabelPair.GetValue(), 64); p.err != nil ***REMOVED***
				// Create a more helpful error message.
				p.parseError(fmt.Sprintf("expected float as value for 'le' label, got %q", p.currentLabelPair.GetValue()))
				return nil
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			p.currentLabels[p.currentLabelPair.GetName()] = p.currentLabelPair.GetValue()
		***REMOVED***
	***REMOVED***
	if p.skipBlankTab(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	switch p.currentByte ***REMOVED***
	case ',':
		return p.startLabelName

	case '***REMOVED***':
		if p.skipBlankTab(); p.err != nil ***REMOVED***
			return nil // Unexpected end of input.
		***REMOVED***
		return p.readingValue
	default:
		p.parseError(fmt.Sprintf("unexpected end of label value %q", p.currentLabelPair.Value))
		return nil
	***REMOVED***
***REMOVED***

// readingValue represents the state where the last byte read (now in
// p.currentByte) is the first byte of the sample value (i.e. a float).
func (p *TextParser) readingValue() stateFn ***REMOVED***
	// When we are here, we have read all the labels, so for the
	// special case of a summary/histogram, we can finally find out
	// if the metric already exists.
	if p.currentMF.GetType() == dto.MetricType_SUMMARY ***REMOVED***
		signature := model.LabelsToSignature(p.currentLabels)
		if summary := p.summaries[signature]; summary != nil ***REMOVED***
			p.currentMetric = summary
		***REMOVED*** else ***REMOVED***
			p.summaries[signature] = p.currentMetric
			p.currentMF.Metric = append(p.currentMF.Metric, p.currentMetric)
		***REMOVED***
	***REMOVED*** else if p.currentMF.GetType() == dto.MetricType_HISTOGRAM ***REMOVED***
		signature := model.LabelsToSignature(p.currentLabels)
		if histogram := p.histograms[signature]; histogram != nil ***REMOVED***
			p.currentMetric = histogram
		***REMOVED*** else ***REMOVED***
			p.histograms[signature] = p.currentMetric
			p.currentMF.Metric = append(p.currentMF.Metric, p.currentMetric)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		p.currentMF.Metric = append(p.currentMF.Metric, p.currentMetric)
	***REMOVED***
	if p.readTokenUntilWhitespace(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	value, err := strconv.ParseFloat(p.currentToken.String(), 64)
	if err != nil ***REMOVED***
		// Create a more helpful error message.
		p.parseError(fmt.Sprintf("expected float as value, got %q", p.currentToken.String()))
		return nil
	***REMOVED***
	switch p.currentMF.GetType() ***REMOVED***
	case dto.MetricType_COUNTER:
		p.currentMetric.Counter = &dto.Counter***REMOVED***Value: proto.Float64(value)***REMOVED***
	case dto.MetricType_GAUGE:
		p.currentMetric.Gauge = &dto.Gauge***REMOVED***Value: proto.Float64(value)***REMOVED***
	case dto.MetricType_UNTYPED:
		p.currentMetric.Untyped = &dto.Untyped***REMOVED***Value: proto.Float64(value)***REMOVED***
	case dto.MetricType_SUMMARY:
		// *sigh*
		if p.currentMetric.Summary == nil ***REMOVED***
			p.currentMetric.Summary = &dto.Summary***REMOVED******REMOVED***
		***REMOVED***
		switch ***REMOVED***
		case p.currentIsSummaryCount:
			p.currentMetric.Summary.SampleCount = proto.Uint64(uint64(value))
		case p.currentIsSummarySum:
			p.currentMetric.Summary.SampleSum = proto.Float64(value)
		case !math.IsNaN(p.currentQuantile):
			p.currentMetric.Summary.Quantile = append(
				p.currentMetric.Summary.Quantile,
				&dto.Quantile***REMOVED***
					Quantile: proto.Float64(p.currentQuantile),
					Value:    proto.Float64(value),
				***REMOVED***,
			)
		***REMOVED***
	case dto.MetricType_HISTOGRAM:
		// *sigh*
		if p.currentMetric.Histogram == nil ***REMOVED***
			p.currentMetric.Histogram = &dto.Histogram***REMOVED******REMOVED***
		***REMOVED***
		switch ***REMOVED***
		case p.currentIsHistogramCount:
			p.currentMetric.Histogram.SampleCount = proto.Uint64(uint64(value))
		case p.currentIsHistogramSum:
			p.currentMetric.Histogram.SampleSum = proto.Float64(value)
		case !math.IsNaN(p.currentBucket):
			p.currentMetric.Histogram.Bucket = append(
				p.currentMetric.Histogram.Bucket,
				&dto.Bucket***REMOVED***
					UpperBound:      proto.Float64(p.currentBucket),
					CumulativeCount: proto.Uint64(uint64(value)),
				***REMOVED***,
			)
		***REMOVED***
	default:
		p.err = fmt.Errorf("unexpected type for metric name %q", p.currentMF.GetName())
	***REMOVED***
	if p.currentByte == '\n' ***REMOVED***
		return p.startOfLine
	***REMOVED***
	return p.startTimestamp
***REMOVED***

// startTimestamp represents the state where the next byte read from p.buf is
// the start of the timestamp (or whitespace leading up to it).
func (p *TextParser) startTimestamp() stateFn ***REMOVED***
	if p.skipBlankTab(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	if p.readTokenUntilWhitespace(); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	timestamp, err := strconv.ParseInt(p.currentToken.String(), 10, 64)
	if err != nil ***REMOVED***
		// Create a more helpful error message.
		p.parseError(fmt.Sprintf("expected integer as timestamp, got %q", p.currentToken.String()))
		return nil
	***REMOVED***
	p.currentMetric.TimestampMs = proto.Int64(timestamp)
	if p.readTokenUntilNewline(false); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	if p.currentToken.Len() > 0 ***REMOVED***
		p.parseError(fmt.Sprintf("spurious string after timestamp: %q", p.currentToken.String()))
		return nil
	***REMOVED***
	return p.startOfLine
***REMOVED***

// readingHelp represents the state where the last byte read (now in
// p.currentByte) is the first byte of the docstring after 'HELP'.
func (p *TextParser) readingHelp() stateFn ***REMOVED***
	if p.currentMF.Help != nil ***REMOVED***
		p.parseError(fmt.Sprintf("second HELP line for metric name %q", p.currentMF.GetName()))
		return nil
	***REMOVED***
	// Rest of line is the docstring.
	if p.readTokenUntilNewline(true); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	p.currentMF.Help = proto.String(p.currentToken.String())
	return p.startOfLine
***REMOVED***

// readingType represents the state where the last byte read (now in
// p.currentByte) is the first byte of the type hint after 'HELP'.
func (p *TextParser) readingType() stateFn ***REMOVED***
	if p.currentMF.Type != nil ***REMOVED***
		p.parseError(fmt.Sprintf("second TYPE line for metric name %q, or TYPE reported after samples", p.currentMF.GetName()))
		return nil
	***REMOVED***
	// Rest of line is the type.
	if p.readTokenUntilNewline(false); p.err != nil ***REMOVED***
		return nil // Unexpected end of input.
	***REMOVED***
	metricType, ok := dto.MetricType_value[strings.ToUpper(p.currentToken.String())]
	if !ok ***REMOVED***
		p.parseError(fmt.Sprintf("unknown metric type %q", p.currentToken.String()))
		return nil
	***REMOVED***
	p.currentMF.Type = dto.MetricType(metricType).Enum()
	return p.startOfLine
***REMOVED***

// parseError sets p.err to a ParseError at the current line with the given
// message.
func (p *TextParser) parseError(msg string) ***REMOVED***
	p.err = ParseError***REMOVED***
		Line: p.lineCount,
		Msg:  msg,
	***REMOVED***
***REMOVED***

// skipBlankTab reads (and discards) bytes from p.buf until it encounters a byte
// that is neither ' ' nor '\t'. That byte is left in p.currentByte.
func (p *TextParser) skipBlankTab() ***REMOVED***
	for ***REMOVED***
		if p.currentByte, p.err = p.buf.ReadByte(); p.err != nil || !isBlankOrTab(p.currentByte) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// skipBlankTabIfCurrentBlankTab works exactly as skipBlankTab but doesn't do
// anything if p.currentByte is neither ' ' nor '\t'.
func (p *TextParser) skipBlankTabIfCurrentBlankTab() ***REMOVED***
	if isBlankOrTab(p.currentByte) ***REMOVED***
		p.skipBlankTab()
	***REMOVED***
***REMOVED***

// readTokenUntilWhitespace copies bytes from p.buf into p.currentToken.  The
// first byte considered is the byte already read (now in p.currentByte).  The
// first whitespace byte encountered is still copied into p.currentByte, but not
// into p.currentToken.
func (p *TextParser) readTokenUntilWhitespace() ***REMOVED***
	p.currentToken.Reset()
	for p.err == nil && !isBlankOrTab(p.currentByte) && p.currentByte != '\n' ***REMOVED***
		p.currentToken.WriteByte(p.currentByte)
		p.currentByte, p.err = p.buf.ReadByte()
	***REMOVED***
***REMOVED***

// readTokenUntilNewline copies bytes from p.buf into p.currentToken.  The first
// byte considered is the byte already read (now in p.currentByte).  The first
// newline byte encountered is still copied into p.currentByte, but not into
// p.currentToken. If recognizeEscapeSequence is true, two escape sequences are
// recognized: '\\' tranlates into '\', and '\n' into a line-feed character. All
// other escape sequences are invalid and cause an error.
func (p *TextParser) readTokenUntilNewline(recognizeEscapeSequence bool) ***REMOVED***
	p.currentToken.Reset()
	escaped := false
	for p.err == nil ***REMOVED***
		if recognizeEscapeSequence && escaped ***REMOVED***
			switch p.currentByte ***REMOVED***
			case '\\':
				p.currentToken.WriteByte(p.currentByte)
			case 'n':
				p.currentToken.WriteByte('\n')
			default:
				p.parseError(fmt.Sprintf("invalid escape sequence '\\%c'", p.currentByte))
				return
			***REMOVED***
			escaped = false
		***REMOVED*** else ***REMOVED***
			switch p.currentByte ***REMOVED***
			case '\n':
				return
			case '\\':
				escaped = true
			default:
				p.currentToken.WriteByte(p.currentByte)
			***REMOVED***
		***REMOVED***
		p.currentByte, p.err = p.buf.ReadByte()
	***REMOVED***
***REMOVED***

// readTokenAsMetricName copies a metric name from p.buf into p.currentToken.
// The first byte considered is the byte already read (now in p.currentByte).
// The first byte not part of a metric name is still copied into p.currentByte,
// but not into p.currentToken.
func (p *TextParser) readTokenAsMetricName() ***REMOVED***
	p.currentToken.Reset()
	if !isValidMetricNameStart(p.currentByte) ***REMOVED***
		return
	***REMOVED***
	for ***REMOVED***
		p.currentToken.WriteByte(p.currentByte)
		p.currentByte, p.err = p.buf.ReadByte()
		if p.err != nil || !isValidMetricNameContinuation(p.currentByte) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// readTokenAsLabelName copies a label name from p.buf into p.currentToken.
// The first byte considered is the byte already read (now in p.currentByte).
// The first byte not part of a label name is still copied into p.currentByte,
// but not into p.currentToken.
func (p *TextParser) readTokenAsLabelName() ***REMOVED***
	p.currentToken.Reset()
	if !isValidLabelNameStart(p.currentByte) ***REMOVED***
		return
	***REMOVED***
	for ***REMOVED***
		p.currentToken.WriteByte(p.currentByte)
		p.currentByte, p.err = p.buf.ReadByte()
		if p.err != nil || !isValidLabelNameContinuation(p.currentByte) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// readTokenAsLabelValue copies a label value from p.buf into p.currentToken.
// In contrast to the other 'readTokenAs...' functions, which start with the
// last read byte in p.currentByte, this method ignores p.currentByte and starts
// with reading a new byte from p.buf. The first byte not part of a label value
// is still copied into p.currentByte, but not into p.currentToken.
func (p *TextParser) readTokenAsLabelValue() ***REMOVED***
	p.currentToken.Reset()
	escaped := false
	for ***REMOVED***
		if p.currentByte, p.err = p.buf.ReadByte(); p.err != nil ***REMOVED***
			return
		***REMOVED***
		if escaped ***REMOVED***
			switch p.currentByte ***REMOVED***
			case '"', '\\':
				p.currentToken.WriteByte(p.currentByte)
			case 'n':
				p.currentToken.WriteByte('\n')
			default:
				p.parseError(fmt.Sprintf("invalid escape sequence '\\%c'", p.currentByte))
				return
			***REMOVED***
			escaped = false
			continue
		***REMOVED***
		switch p.currentByte ***REMOVED***
		case '"':
			return
		case '\n':
			p.parseError(fmt.Sprintf("label value %q contains unescaped new-line", p.currentToken.String()))
			return
		case '\\':
			escaped = true
		default:
			p.currentToken.WriteByte(p.currentByte)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *TextParser) setOrCreateCurrentMF() ***REMOVED***
	p.currentIsSummaryCount = false
	p.currentIsSummarySum = false
	p.currentIsHistogramCount = false
	p.currentIsHistogramSum = false
	name := p.currentToken.String()
	if p.currentMF = p.metricFamiliesByName[name]; p.currentMF != nil ***REMOVED***
		return
	***REMOVED***
	// Try out if this is a _sum or _count for a summary/histogram.
	summaryName := summaryMetricName(name)
	if p.currentMF = p.metricFamiliesByName[summaryName]; p.currentMF != nil ***REMOVED***
		if p.currentMF.GetType() == dto.MetricType_SUMMARY ***REMOVED***
			if isCount(name) ***REMOVED***
				p.currentIsSummaryCount = true
			***REMOVED***
			if isSum(name) ***REMOVED***
				p.currentIsSummarySum = true
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	histogramName := histogramMetricName(name)
	if p.currentMF = p.metricFamiliesByName[histogramName]; p.currentMF != nil ***REMOVED***
		if p.currentMF.GetType() == dto.MetricType_HISTOGRAM ***REMOVED***
			if isCount(name) ***REMOVED***
				p.currentIsHistogramCount = true
			***REMOVED***
			if isSum(name) ***REMOVED***
				p.currentIsHistogramSum = true
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	p.currentMF = &dto.MetricFamily***REMOVED***Name: proto.String(name)***REMOVED***
	p.metricFamiliesByName[name] = p.currentMF
***REMOVED***

func isValidLabelNameStart(b byte) bool ***REMOVED***
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
***REMOVED***

func isValidLabelNameContinuation(b byte) bool ***REMOVED***
	return isValidLabelNameStart(b) || (b >= '0' && b <= '9')
***REMOVED***

func isValidMetricNameStart(b byte) bool ***REMOVED***
	return isValidLabelNameStart(b) || b == ':'
***REMOVED***

func isValidMetricNameContinuation(b byte) bool ***REMOVED***
	return isValidLabelNameContinuation(b) || b == ':'
***REMOVED***

func isBlankOrTab(b byte) bool ***REMOVED***
	return b == ' ' || b == '\t'
***REMOVED***

func isCount(name string) bool ***REMOVED***
	return len(name) > 6 && name[len(name)-6:] == "_count"
***REMOVED***

func isSum(name string) bool ***REMOVED***
	return len(name) > 4 && name[len(name)-4:] == "_sum"
***REMOVED***

func isBucket(name string) bool ***REMOVED***
	return len(name) > 7 && name[len(name)-7:] == "_bucket"
***REMOVED***

func summaryMetricName(name string) string ***REMOVED***
	switch ***REMOVED***
	case isCount(name):
		return name[:len(name)-6]
	case isSum(name):
		return name[:len(name)-4]
	default:
		return name
	***REMOVED***
***REMOVED***

func histogramMetricName(name string) string ***REMOVED***
	switch ***REMOVED***
	case isCount(name):
		return name[:len(name)-6]
	case isSum(name):
		return name[:len(name)-4]
	case isBucket(name):
		return name[:len(name)-7]
	default:
		return name
	***REMOVED***
***REMOVED***
