// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// API/gRPC features intentionally missing from this client:
// - You cannot have the server pick the time of the entry. This client
//   always sends a time.
// - There is no way to provide a protocol buffer payload.
// - No support for the "partial success" feature when writing log entries.

// TODO(jba): test whether forward-slash characters in the log ID must be URL-encoded.
// These features are missing now, but will likely be added:
// - There is no way to specify CallOptions.

package logging

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	vkit "cloud.google.com/go/logging/apiv2"
	"cloud.google.com/go/logging/internal"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/support/bundler"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
	logtypepb "google.golang.org/genproto/googleapis/logging/type"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
)

const (
	// Scope for reading from the logging service.
	ReadScope = "https://www.googleapis.com/auth/logging.read"

	// Scope for writing to the logging service.
	WriteScope = "https://www.googleapis.com/auth/logging.write"

	// Scope for administrative actions on the logging service.
	AdminScope = "https://www.googleapis.com/auth/logging.admin"
)

const (
	// defaultErrorCapacity is the capacity of the channel used to deliver
	// errors to the OnError function.
	defaultErrorCapacity = 10

	// DefaultDelayThreshold is the default value for the DelayThreshold LoggerOption.
	DefaultDelayThreshold = time.Second

	// DefaultEntryCountThreshold is the default value for the EntryCountThreshold LoggerOption.
	DefaultEntryCountThreshold = 1000

	// DefaultEntryByteThreshold is the default value for the EntryByteThreshold LoggerOption.
	DefaultEntryByteThreshold = 1 << 20 // 1MiB

	// DefaultBufferedByteLimit is the default value for the BufferedByteLimit LoggerOption.
	DefaultBufferedByteLimit = 1 << 30 // 1GiB
)

// For testing:
var now = time.Now

// ErrOverflow signals that the number of buffered entries for a Logger
// exceeds its BufferLimit.
var ErrOverflow = errors.New("logging: log entry overflowed buffer limits")

// Client is a Logging client. A Client is associated with a single Cloud project.
type Client struct ***REMOVED***
	client    *vkit.Client // client for the logging service
	projectID string
	errc      chan error     // should be buffered to minimize dropped errors
	donec     chan struct***REMOVED******REMOVED***  // closed on Client.Close to close Logger bundlers
	loggers   sync.WaitGroup // so we can wait for loggers to close
	closed    bool

	// OnError is called when an error occurs in a call to Log or Flush. The
	// error may be due to an invalid Entry, an overflow because BufferLimit
	// was reached (in which case the error will be ErrOverflow) or an error
	// communicating with the logging service. OnError is called with errors
	// from all Loggers. It is never called concurrently. OnError is expected
	// to return quickly; if errors occur while OnError is running, some may
	// not be reported. The default behavior is to call log.Printf.
	//
	// This field should be set only once, before any method of Client is called.
	OnError func(err error)
***REMOVED***

// NewClient returns a new logging client associated with the provided project ID.
//
// By default NewClient uses WriteScope. To use a different scope, call
// NewClient using a WithScopes option (see https://godoc.org/google.golang.org/api/option#WithScopes).
func NewClient(ctx context.Context, projectID string, opts ...option.ClientOption) (*Client, error) ***REMOVED***
	// Check for '/' in project ID to reserve the ability to support various owning resources,
	// in the form "***REMOVED***Collection***REMOVED***/***REMOVED***Name***REMOVED***", for instance "organizations/my-org".
	if strings.ContainsRune(projectID, '/') ***REMOVED***
		return nil, errors.New("logging: project ID contains '/'")
	***REMOVED***
	opts = append([]option.ClientOption***REMOVED***
		option.WithEndpoint(internal.ProdAddr),
		option.WithScopes(WriteScope),
	***REMOVED***, opts...)
	c, err := vkit.NewClient(ctx, opts...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c.SetGoogleClientInfo("logging", internal.Version)
	client := &Client***REMOVED***
		client:    c,
		projectID: projectID,
		errc:      make(chan error, defaultErrorCapacity), // create a small buffer for errors
		donec:     make(chan struct***REMOVED******REMOVED***),
		OnError:   func(e error) ***REMOVED*** log.Printf("logging client: %v", e) ***REMOVED***,
	***REMOVED***
	// Call the user's function synchronously, to make life easier for them.
	go func() ***REMOVED***
		for err := range client.errc ***REMOVED***
			// This reference to OnError is memory-safe if the user sets OnError before
			// calling any client methods. The reference happens before the first read from
			// client.errc, which happens before the first write to client.errc, which
			// happens before any call, which happens before the user sets OnError.
			if fn := client.OnError; fn != nil ***REMOVED***
				fn(err)
			***REMOVED*** else ***REMOVED***
				log.Printf("logging (project ID %q): %v", projectID, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	return client, nil
***REMOVED***

// parent returns the string used in many RPCs to denote the parent resource of the log.
func (c *Client) parent() string ***REMOVED***
	return "projects/" + c.projectID
***REMOVED***

var unixZeroTimestamp *tspb.Timestamp

func init() ***REMOVED***
	var err error
	unixZeroTimestamp, err = ptypes.TimestampProto(time.Unix(0, 0))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

// Ping reports whether the client's connection to the logging service and the
// authentication configuration are valid. To accomplish this, Ping writes a
// log entry "ping" to a log named "ping".
func (c *Client) Ping(ctx context.Context) error ***REMOVED***
	ent := &logpb.LogEntry***REMOVED***
		Payload:   &logpb.LogEntry_TextPayload***REMOVED***"ping"***REMOVED***,
		Timestamp: unixZeroTimestamp, // Identical timestamps and insert IDs are both
		InsertId:  "ping",            // necessary for the service to dedup these entries.
	***REMOVED***
	_, err := c.client.WriteLogEntries(ctx, &logpb.WriteLogEntriesRequest***REMOVED***
		LogName:  internal.LogPath(c.parent(), "ping"),
		Resource: &mrpb.MonitoredResource***REMOVED***Type: "global"***REMOVED***,
		Entries:  []*logpb.LogEntry***REMOVED***ent***REMOVED***,
	***REMOVED***)
	return err
***REMOVED***

// A Logger is used to write log messages to a single log. It can be configured
// with a log ID, common monitored resource, and a set of common labels.
type Logger struct ***REMOVED***
	client     *Client
	logName    string // "projects/***REMOVED***projectID***REMOVED***/logs/***REMOVED***logID***REMOVED***"
	stdLoggers map[Severity]*log.Logger
	bundler    *bundler.Bundler

	// Options
	commonResource *mrpb.MonitoredResource
	commonLabels   map[string]string
***REMOVED***

// A LoggerOption is a configuration option for a Logger.
type LoggerOption interface ***REMOVED***
	set(*Logger)
***REMOVED***

// CommonResource sets the monitored resource associated with all log entries
// written from a Logger. If not provided, a resource of type "global" is used.
// This value can be overridden by setting an Entry's Resource field.
func CommonResource(r *mrpb.MonitoredResource) LoggerOption ***REMOVED*** return commonResource***REMOVED***r***REMOVED*** ***REMOVED***

type commonResource struct***REMOVED*** *mrpb.MonitoredResource ***REMOVED***

func (r commonResource) set(l *Logger) ***REMOVED*** l.commonResource = r.MonitoredResource ***REMOVED***

// CommonLabels are labels that apply to all log entries written from a Logger,
// so that you don't have to repeat them in each log entry's Labels field. If
// any of the log entries contains a (key, value) with the same key that is in
// CommonLabels, then the entry's (key, value) overrides the one in
// CommonLabels.
func CommonLabels(m map[string]string) LoggerOption ***REMOVED*** return commonLabels(m) ***REMOVED***

type commonLabels map[string]string

func (c commonLabels) set(l *Logger) ***REMOVED*** l.commonLabels = c ***REMOVED***

// DelayThreshold is the maximum amount of time that an entry should remain
// buffered in memory before a call to the logging service is triggered. Larger
// values of DelayThreshold will generally result in fewer calls to the logging
// service, while increasing the risk that log entries will be lost if the
// process crashes.
// The default is DefaultDelayThreshold.
func DelayThreshold(d time.Duration) LoggerOption ***REMOVED*** return delayThreshold(d) ***REMOVED***

type delayThreshold time.Duration

func (d delayThreshold) set(l *Logger) ***REMOVED*** l.bundler.DelayThreshold = time.Duration(d) ***REMOVED***

// EntryCountThreshold is the maximum number of entries that will be buffered
// in memory before a call to the logging service is triggered. Larger values
// will generally result in fewer calls to the logging service, while
// increasing both memory consumption and the risk that log entries will be
// lost if the process crashes.
// The default is DefaultEntryCountThreshold.
func EntryCountThreshold(n int) LoggerOption ***REMOVED*** return entryCountThreshold(n) ***REMOVED***

type entryCountThreshold int

func (e entryCountThreshold) set(l *Logger) ***REMOVED*** l.bundler.BundleCountThreshold = int(e) ***REMOVED***

// EntryByteThreshold is the maximum number of bytes of entries that will be
// buffered in memory before a call to the logging service is triggered. See
// EntryCountThreshold for a discussion of the tradeoffs involved in setting
// this option.
// The default is DefaultEntryByteThreshold.
func EntryByteThreshold(n int) LoggerOption ***REMOVED*** return entryByteThreshold(n) ***REMOVED***

type entryByteThreshold int

func (e entryByteThreshold) set(l *Logger) ***REMOVED*** l.bundler.BundleByteThreshold = int(e) ***REMOVED***

// EntryByteLimit is the maximum number of bytes of entries that will be sent
// in a single call to the logging service. This option limits the size of a
// single RPC payload, to account for network or service issues with large
// RPCs. If EntryByteLimit is smaller than EntryByteThreshold, the latter has
// no effect.
// The default is zero, meaning there is no limit.
func EntryByteLimit(n int) LoggerOption ***REMOVED*** return entryByteLimit(n) ***REMOVED***

type entryByteLimit int

func (e entryByteLimit) set(l *Logger) ***REMOVED*** l.bundler.BundleByteLimit = int(e) ***REMOVED***

// BufferedByteLimit is the maximum number of bytes that the Logger will keep
// in memory before returning ErrOverflow. This option limits the total memory
// consumption of the Logger (but note that each Logger has its own, separate
// limit). It is possible to reach BufferedByteLimit even if it is larger than
// EntryByteThreshold or EntryByteLimit, because calls triggered by the latter
// two options may be enqueued (and hence occupying memory) while new log
// entries are being added.
// The default is DefaultBufferedByteLimit.
func BufferedByteLimit(n int) LoggerOption ***REMOVED*** return bufferedByteLimit(n) ***REMOVED***

type bufferedByteLimit int

func (b bufferedByteLimit) set(l *Logger) ***REMOVED*** l.bundler.BufferedByteLimit = int(b) ***REMOVED***

// Logger returns a Logger that will write entries with the given log ID, such as
// "syslog". A log ID must be less than 512 characters long and can only
// include the following characters: upper and lower case alphanumeric
// characters: [A-Za-z0-9]; and punctuation characters: forward-slash,
// underscore, hyphen, and period.
func (c *Client) Logger(logID string, opts ...LoggerOption) *Logger ***REMOVED***
	l := &Logger***REMOVED***
		client:         c,
		logName:        internal.LogPath(c.parent(), logID),
		commonResource: &mrpb.MonitoredResource***REMOVED***Type: "global"***REMOVED***,
	***REMOVED***
	// TODO(jba): determine the right context for the bundle handler.
	ctx := context.TODO()
	l.bundler = bundler.NewBundler(&logpb.LogEntry***REMOVED******REMOVED***, func(entries interface***REMOVED******REMOVED***) ***REMOVED***
		l.writeLogEntries(ctx, entries.([]*logpb.LogEntry))
	***REMOVED***)
	l.bundler.DelayThreshold = DefaultDelayThreshold
	l.bundler.BundleCountThreshold = DefaultEntryCountThreshold
	l.bundler.BundleByteThreshold = DefaultEntryByteThreshold
	l.bundler.BufferedByteLimit = DefaultBufferedByteLimit
	for _, opt := range opts ***REMOVED***
		opt.set(l)
	***REMOVED***

	l.stdLoggers = map[Severity]*log.Logger***REMOVED******REMOVED***
	for s := range severityName ***REMOVED***
		l.stdLoggers[s] = log.New(severityWriter***REMOVED***l, s***REMOVED***, "", 0)
	***REMOVED***
	c.loggers.Add(1)
	go func() ***REMOVED***
		defer c.loggers.Done()
		<-c.donec
		l.bundler.Close()
	***REMOVED***()
	return l
***REMOVED***

type severityWriter struct ***REMOVED***
	l *Logger
	s Severity
***REMOVED***

func (w severityWriter) Write(p []byte) (n int, err error) ***REMOVED***
	w.l.Log(Entry***REMOVED***
		Severity: w.s,
		Payload:  string(p),
	***REMOVED***)
	return len(p), nil
***REMOVED***

// Close closes the client.
func (c *Client) Close() error ***REMOVED***
	if c.closed ***REMOVED***
		return nil
	***REMOVED***
	close(c.donec)   // close Logger bundlers
	c.loggers.Wait() // wait for all bundlers to flush and close
	// Now there can be no more errors.
	close(c.errc) // terminate error goroutine
	// Return only the first error. Since all clients share an underlying connection,
	// Closes after the first always report a "connection is closing" error.
	err := c.client.Close()
	c.closed = true
	return err
***REMOVED***

// Severity is the severity of the event described in a log entry. These
// guideline severity levels are ordered, with numerically smaller levels
// treated as less severe than numerically larger levels.
type Severity int

const (
	// Default means the log entry has no assigned severity level.
	Default = Severity(logtypepb.LogSeverity_DEFAULT)
	// Debug means debug or trace information.
	Debug = Severity(logtypepb.LogSeverity_DEBUG)
	// Info means routine information, such as ongoing status or performance.
	Info = Severity(logtypepb.LogSeverity_INFO)
	// Notice means normal but significant events, such as start up, shut down, or configuration.
	Notice = Severity(logtypepb.LogSeverity_NOTICE)
	// Warning means events that might cause problems.
	Warning = Severity(logtypepb.LogSeverity_WARNING)
	// Error means events that are likely to cause problems.
	Error = Severity(logtypepb.LogSeverity_ERROR)
	// Critical means events that cause more severe problems or brief outages.
	Critical = Severity(logtypepb.LogSeverity_CRITICAL)
	// Alert means a person must take an action immediately.
	Alert = Severity(logtypepb.LogSeverity_ALERT)
	// Emergency means one or more systems are unusable.
	Emergency = Severity(logtypepb.LogSeverity_EMERGENCY)
)

var severityName = map[Severity]string***REMOVED***
	Default:   "Default",
	Debug:     "Debug",
	Info:      "Info",
	Notice:    "Notice",
	Warning:   "Warning",
	Error:     "Error",
	Critical:  "Critical",
	Alert:     "Alert",
	Emergency: "Emergency",
***REMOVED***

// String converts a severity level to a string.
func (v Severity) String() string ***REMOVED***
	// same as proto.EnumName
	s, ok := severityName[v]
	if ok ***REMOVED***
		return s
	***REMOVED***
	return strconv.Itoa(int(v))
***REMOVED***

// ParseSeverity returns the Severity whose name equals s, ignoring case. It
// returns Default if no Severity matches.
func ParseSeverity(s string) Severity ***REMOVED***
	sl := strings.ToLower(s)
	for sev, name := range severityName ***REMOVED***
		if strings.ToLower(name) == sl ***REMOVED***
			return sev
		***REMOVED***
	***REMOVED***
	return Default
***REMOVED***

// Entry is a log entry.
// See https://cloud.google.com/logging/docs/view/logs_index for more about entries.
type Entry struct ***REMOVED***
	// Timestamp is the time of the entry. If zero, the current time is used.
	Timestamp time.Time

	// Severity is the entry's severity level.
	// The zero value is Default.
	Severity Severity

	// Payload must be either a string or something that
	// marshals via the encoding/json package to a JSON object
	// (and not any other type of JSON value).
	Payload interface***REMOVED******REMOVED***

	// Labels optionally specifies key/value labels for the log entry.
	// The Logger.Log method takes ownership of this map. See Logger.CommonLabels
	// for more about labels.
	Labels map[string]string

	// InsertID is a unique ID for the log entry. If you provide this field,
	// the logging service considers other log entries in the same log with the
	// same ID as duplicates which can be removed. If omitted, the logging
	// service will generate a unique ID for this log entry. Note that because
	// this client retries RPCs automatically, it is possible (though unlikely)
	// that an Entry without an InsertID will be written more than once.
	InsertID string

	// HTTPRequest optionally specifies metadata about the HTTP request
	// associated with this log entry, if applicable. It is optional.
	HTTPRequest *HTTPRequest

	// Operation optionally provides information about an operation associated
	// with the log entry, if applicable.
	Operation *logpb.LogEntryOperation

	// LogName is the full log name, in the form
	// "projects/***REMOVED***ProjectID***REMOVED***/logs/***REMOVED***LogID***REMOVED***". It is set by the client when
	// reading entries. It is an error to set it when writing entries.
	LogName string

	// Resource is the monitored resource associated with the entry. It is set
	// by the client when reading entries. It is an error to set it when
	// writing entries.
	Resource *mrpb.MonitoredResource
***REMOVED***

// HTTPRequest contains an http.Request as well as additional
// information about the request and its response.
type HTTPRequest struct ***REMOVED***
	// Request is the http.Request passed to the handler.
	Request *http.Request

	// RequestSize is the size of the HTTP request message in bytes, including
	// the request headers and the request body.
	RequestSize int64

	// Status is the response code indicating the status of the response.
	// Examples: 200, 404.
	Status int

	// ResponseSize is the size of the HTTP response message sent back to the client, in bytes,
	// including the response headers and the response body.
	ResponseSize int64

	// Latency is the request processing latency on the server, from the time the request was
	// received until the response was sent.
	Latency time.Duration

	// RemoteIP is the IP address (IPv4 or IPv6) of the client that issued the
	// HTTP request. Examples: "192.168.1.1", "FE80::0202:B3FF:FE1E:8329".
	RemoteIP string

	// CacheHit reports whether an entity was served from cache (with or without
	// validation).
	CacheHit bool

	// CacheValidatedWithOriginServer reports whether the response was
	// validated with the origin server before being served from cache. This
	// field is only meaningful if CacheHit is true.
	CacheValidatedWithOriginServer bool
***REMOVED***

func fromHTTPRequest(r *HTTPRequest) *logtypepb.HttpRequest ***REMOVED***
	if r == nil ***REMOVED***
		return nil
	***REMOVED***
	if r.Request == nil ***REMOVED***
		panic("HTTPRequest must have a non-nil Request")
	***REMOVED***
	u := *r.Request.URL
	u.Fragment = ""
	return &logtypepb.HttpRequest***REMOVED***
		RequestMethod:                  r.Request.Method,
		RequestUrl:                     u.String(),
		RequestSize:                    r.RequestSize,
		Status:                         int32(r.Status),
		ResponseSize:                   r.ResponseSize,
		Latency:                        ptypes.DurationProto(r.Latency),
		UserAgent:                      r.Request.UserAgent(),
		RemoteIp:                       r.RemoteIP, // TODO(jba): attempt to parse http.Request.RemoteAddr?
		Referer:                        r.Request.Referer(),
		CacheHit:                       r.CacheHit,
		CacheValidatedWithOriginServer: r.CacheValidatedWithOriginServer,
	***REMOVED***
***REMOVED***

// toProtoStruct converts v, which must marshal into a JSON object,
// into a Google Struct proto.
func toProtoStruct(v interface***REMOVED******REMOVED***) (*structpb.Struct, error) ***REMOVED***
	// Fast path: if v is already a *structpb.Struct, nothing to do.
	if s, ok := v.(*structpb.Struct); ok ***REMOVED***
		return s, nil
	***REMOVED***
	// v is a Go struct that supports JSON marshalling. We want a Struct
	// protobuf. Some day we may have a more direct way to get there, but right
	// now the only way is to marshal the Go struct to JSON, unmarshal into a
	// map, and then build the Struct proto from the map.
	jb, err := json.Marshal(v)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("logging: json.Marshal: %v", err)
	***REMOVED***
	var m map[string]interface***REMOVED******REMOVED***
	err = json.Unmarshal(jb, &m)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("logging: json.Unmarshal: %v", err)
	***REMOVED***
	return jsonMapToProtoStruct(m), nil
***REMOVED***

func jsonMapToProtoStruct(m map[string]interface***REMOVED******REMOVED***) *structpb.Struct ***REMOVED***
	fields := map[string]*structpb.Value***REMOVED******REMOVED***
	for k, v := range m ***REMOVED***
		fields[k] = jsonValueToStructValue(v)
	***REMOVED***
	return &structpb.Struct***REMOVED***Fields: fields***REMOVED***
***REMOVED***

func jsonValueToStructValue(v interface***REMOVED******REMOVED***) *structpb.Value ***REMOVED***
	switch x := v.(type) ***REMOVED***
	case bool:
		return &structpb.Value***REMOVED***Kind: &structpb.Value_BoolValue***REMOVED***x***REMOVED******REMOVED***
	case float64:
		return &structpb.Value***REMOVED***Kind: &structpb.Value_NumberValue***REMOVED***x***REMOVED******REMOVED***
	case string:
		return &structpb.Value***REMOVED***Kind: &structpb.Value_StringValue***REMOVED***x***REMOVED******REMOVED***
	case nil:
		return &structpb.Value***REMOVED***Kind: &structpb.Value_NullValue***REMOVED******REMOVED******REMOVED***
	case map[string]interface***REMOVED******REMOVED***:
		return &structpb.Value***REMOVED***Kind: &structpb.Value_StructValue***REMOVED***jsonMapToProtoStruct(x)***REMOVED******REMOVED***
	case []interface***REMOVED******REMOVED***:
		var vals []*structpb.Value
		for _, e := range x ***REMOVED***
			vals = append(vals, jsonValueToStructValue(e))
		***REMOVED***
		return &structpb.Value***REMOVED***Kind: &structpb.Value_ListValue***REMOVED***&structpb.ListValue***REMOVED***vals***REMOVED******REMOVED******REMOVED***
	default:
		panic(fmt.Sprintf("bad type %T for JSON value", v))
	***REMOVED***
***REMOVED***

// LogSync logs the Entry synchronously without any buffering. Because LogSync is slow
// and will block, it is intended primarily for debugging or critical errors.
// Prefer Log for most uses.
// TODO(jba): come up with a better name (LogNow?) or eliminate.
func (l *Logger) LogSync(ctx context.Context, e Entry) error ***REMOVED***
	ent, err := toLogEntry(e)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = l.client.client.WriteLogEntries(ctx, &logpb.WriteLogEntriesRequest***REMOVED***
		LogName:  l.logName,
		Resource: l.commonResource,
		Labels:   l.commonLabels,
		Entries:  []*logpb.LogEntry***REMOVED***ent***REMOVED***,
	***REMOVED***)
	return err
***REMOVED***

// Log buffers the Entry for output to the logging service. It never blocks.
func (l *Logger) Log(e Entry) ***REMOVED***
	ent, err := toLogEntry(e)
	if err != nil ***REMOVED***
		l.error(err)
		return
	***REMOVED***
	if err := l.bundler.Add(ent, proto.Size(ent)); err != nil ***REMOVED***
		l.error(err)
	***REMOVED***
***REMOVED***

// Flush blocks until all currently buffered log entries are sent.
func (l *Logger) Flush() ***REMOVED***
	l.bundler.Flush()
***REMOVED***

func (l *Logger) writeLogEntries(ctx context.Context, entries []*logpb.LogEntry) ***REMOVED***
	req := &logpb.WriteLogEntriesRequest***REMOVED***
		LogName:  l.logName,
		Resource: l.commonResource,
		Labels:   l.commonLabels,
		Entries:  entries,
	***REMOVED***
	_, err := l.client.client.WriteLogEntries(ctx, req)
	if err != nil ***REMOVED***
		l.error(err)
	***REMOVED***
***REMOVED***

// error puts the error on the client's error channel
// without blocking.
func (l *Logger) error(err error) ***REMOVED***
	select ***REMOVED***
	case l.client.errc <- err:
	default:
	***REMOVED***
***REMOVED***

// StandardLogger returns a *log.Logger for the provided severity.
//
// This method is cheap. A single log.Logger is pre-allocated for each
// severity level in each Logger. Callers may mutate the returned log.Logger
// (for example by calling SetFlags or SetPrefix).
func (l *Logger) StandardLogger(s Severity) *log.Logger ***REMOVED*** return l.stdLoggers[s] ***REMOVED***

func trunc32(i int) int32 ***REMOVED***
	if i > math.MaxInt32 ***REMOVED***
		i = math.MaxInt32
	***REMOVED***
	return int32(i)
***REMOVED***

func toLogEntry(e Entry) (*logpb.LogEntry, error) ***REMOVED***
	if e.LogName != "" ***REMOVED***
		return nil, errors.New("logging: Entry.LogName should be not be set when writing")
	***REMOVED***
	t := e.Timestamp
	if t.IsZero() ***REMOVED***
		t = now()
	***REMOVED***
	ts, err := ptypes.TimestampProto(t)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ent := &logpb.LogEntry***REMOVED***
		Timestamp:   ts,
		Severity:    logtypepb.LogSeverity(e.Severity),
		InsertId:    e.InsertID,
		HttpRequest: fromHTTPRequest(e.HTTPRequest),
		Operation:   e.Operation,
		Labels:      e.Labels,
	***REMOVED***

	switch p := e.Payload.(type) ***REMOVED***
	case string:
		ent.Payload = &logpb.LogEntry_TextPayload***REMOVED***p***REMOVED***
	default:
		s, err := toProtoStruct(p)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ent.Payload = &logpb.LogEntry_JsonPayload***REMOVED***s***REMOVED***
	***REMOVED***
	return ent, nil
***REMOVED***
