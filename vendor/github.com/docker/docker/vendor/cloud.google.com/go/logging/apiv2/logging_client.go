// Copyright 2016, Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// AUTO-GENERATED CODE. DO NOT EDIT.

package logging

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"time"

	gax "github.com/googleapis/gax-go"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var (
	loggingParentPathTemplate = gax.MustCompilePathTemplate("projects/***REMOVED***project***REMOVED***")
	loggingLogPathTemplate    = gax.MustCompilePathTemplate("projects/***REMOVED***project***REMOVED***/logs/***REMOVED***log***REMOVED***")
)

// CallOptions contains the retry settings for each method of Client.
type CallOptions struct ***REMOVED***
	DeleteLog                        []gax.CallOption
	WriteLogEntries                  []gax.CallOption
	ListLogEntries                   []gax.CallOption
	ListMonitoredResourceDescriptors []gax.CallOption
***REMOVED***

func defaultClientOptions() []option.ClientOption ***REMOVED***
	return []option.ClientOption***REMOVED***
		option.WithEndpoint("logging.googleapis.com:443"),
		option.WithScopes(
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/cloud-platform.read-only",
			"https://www.googleapis.com/auth/logging.admin",
			"https://www.googleapis.com/auth/logging.read",
			"https://www.googleapis.com/auth/logging.write",
		),
	***REMOVED***
***REMOVED***

func defaultCallOptions() *CallOptions ***REMOVED***
	retry := map[[2]string][]gax.CallOption***REMOVED***
		***REMOVED***"default", "idempotent"***REMOVED***: ***REMOVED***
			gax.WithRetry(func() gax.Retryer ***REMOVED***
				return gax.OnCodes([]codes.Code***REMOVED***
					codes.DeadlineExceeded,
					codes.Unavailable,
				***REMOVED***, gax.Backoff***REMOVED***
					Initial:    100 * time.Millisecond,
					Max:        1000 * time.Millisecond,
					Multiplier: 1.2,
				***REMOVED***)
			***REMOVED***),
		***REMOVED***,
		***REMOVED***"list", "idempotent"***REMOVED***: ***REMOVED***
			gax.WithRetry(func() gax.Retryer ***REMOVED***
				return gax.OnCodes([]codes.Code***REMOVED***
					codes.DeadlineExceeded,
					codes.Unavailable,
				***REMOVED***, gax.Backoff***REMOVED***
					Initial:    100 * time.Millisecond,
					Max:        1000 * time.Millisecond,
					Multiplier: 1.2,
				***REMOVED***)
			***REMOVED***),
		***REMOVED***,
	***REMOVED***
	return &CallOptions***REMOVED***
		DeleteLog:                        retry[[2]string***REMOVED***"default", "idempotent"***REMOVED***],
		WriteLogEntries:                  retry[[2]string***REMOVED***"default", "non_idempotent"***REMOVED***],
		ListLogEntries:                   retry[[2]string***REMOVED***"list", "idempotent"***REMOVED***],
		ListMonitoredResourceDescriptors: retry[[2]string***REMOVED***"default", "idempotent"***REMOVED***],
	***REMOVED***
***REMOVED***

// Client is a client for interacting with Stackdriver Logging API.
type Client struct ***REMOVED***
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	client loggingpb.LoggingServiceV2Client

	// The call options for this service.
	CallOptions *CallOptions

	// The metadata to be sent with each request.
	metadata metadata.MD
***REMOVED***

// NewClient creates a new logging service v2 client.
//
// Service for ingesting and querying logs.
func NewClient(ctx context.Context, opts ...option.ClientOption) (*Client, error) ***REMOVED***
	conn, err := transport.DialGRPC(ctx, append(defaultClientOptions(), opts...)...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c := &Client***REMOVED***
		conn:        conn,
		CallOptions: defaultCallOptions(),

		client: loggingpb.NewLoggingServiceV2Client(conn),
	***REMOVED***
	c.SetGoogleClientInfo("gax", gax.Version)
	return c, nil
***REMOVED***

// Connection returns the client's connection to the API service.
func (c *Client) Connection() *grpc.ClientConn ***REMOVED***
	return c.conn
***REMOVED***

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *Client) Close() error ***REMOVED***
	return c.conn.Close()
***REMOVED***

// SetGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *Client) SetGoogleClientInfo(name, version string) ***REMOVED***
	goVersion := strings.Replace(runtime.Version(), " ", "_", -1)
	v := fmt.Sprintf("%s/%s %s gax/%s go/%s", name, version, gapicNameVersion, gax.Version, goVersion)
	c.metadata = metadata.Pairs("x-goog-api-client", v)
***REMOVED***

// LoggingParentPath returns the path for the parent resource.
func LoggingParentPath(project string) string ***REMOVED***
	path, err := loggingParentPathTemplate.Render(map[string]string***REMOVED***
		"project": project,
	***REMOVED***)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return path
***REMOVED***

// LoggingLogPath returns the path for the log resource.
func LoggingLogPath(project, log string) string ***REMOVED***
	path, err := loggingLogPathTemplate.Render(map[string]string***REMOVED***
		"project": project,
		"log":     log,
	***REMOVED***)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return path
***REMOVED***

// DeleteLog deletes all the log entries in a log.
// The log reappears if it receives new entries.
func (c *Client) DeleteLog(ctx context.Context, req *loggingpb.DeleteLogRequest) error ***REMOVED***
	md, _ := metadata.FromContext(ctx)
	ctx = metadata.NewContext(ctx, metadata.Join(md, c.metadata))
	err := gax.Invoke(ctx, func(ctx context.Context) error ***REMOVED***
		var err error
		_, err = c.client.DeleteLog(ctx, req)
		return err
	***REMOVED***, c.CallOptions.DeleteLog...)
	return err
***REMOVED***

// WriteLogEntries writes log entries to Stackdriver Logging.  All log entries are
// written by this method.
func (c *Client) WriteLogEntries(ctx context.Context, req *loggingpb.WriteLogEntriesRequest) (*loggingpb.WriteLogEntriesResponse, error) ***REMOVED***
	md, _ := metadata.FromContext(ctx)
	ctx = metadata.NewContext(ctx, metadata.Join(md, c.metadata))
	var resp *loggingpb.WriteLogEntriesResponse
	err := gax.Invoke(ctx, func(ctx context.Context) error ***REMOVED***
		var err error
		resp, err = c.client.WriteLogEntries(ctx, req)
		return err
	***REMOVED***, c.CallOptions.WriteLogEntries...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp, nil
***REMOVED***

// ListLogEntries lists log entries.  Use this method to retrieve log entries from Cloud
// Logging.  For ways to export log entries, see
// [Exporting Logs](/logging/docs/export).
func (c *Client) ListLogEntries(ctx context.Context, req *loggingpb.ListLogEntriesRequest) *LogEntryIterator ***REMOVED***
	md, _ := metadata.FromContext(ctx)
	ctx = metadata.NewContext(ctx, metadata.Join(md, c.metadata))
	it := &LogEntryIterator***REMOVED******REMOVED***
	it.InternalFetch = func(pageSize int, pageToken string) ([]*loggingpb.LogEntry, string, error) ***REMOVED***
		var resp *loggingpb.ListLogEntriesResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 ***REMOVED***
			req.PageSize = math.MaxInt32
		***REMOVED*** else ***REMOVED***
			req.PageSize = int32(pageSize)
		***REMOVED***
		err := gax.Invoke(ctx, func(ctx context.Context) error ***REMOVED***
			var err error
			resp, err = c.client.ListLogEntries(ctx, req)
			return err
		***REMOVED***, c.CallOptions.ListLogEntries...)
		if err != nil ***REMOVED***
			return nil, "", err
		***REMOVED***
		return resp.Entries, resp.NextPageToken, nil
	***REMOVED***
	fetch := func(pageSize int, pageToken string) (string, error) ***REMOVED***
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		it.items = append(it.items, items...)
		return nextPageToken, nil
	***REMOVED***
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	return it
***REMOVED***

// ListMonitoredResourceDescriptors lists the monitored resource descriptors used by Stackdriver Logging.
func (c *Client) ListMonitoredResourceDescriptors(ctx context.Context, req *loggingpb.ListMonitoredResourceDescriptorsRequest) *MonitoredResourceDescriptorIterator ***REMOVED***
	md, _ := metadata.FromContext(ctx)
	ctx = metadata.NewContext(ctx, metadata.Join(md, c.metadata))
	it := &MonitoredResourceDescriptorIterator***REMOVED******REMOVED***
	it.InternalFetch = func(pageSize int, pageToken string) ([]*monitoredrespb.MonitoredResourceDescriptor, string, error) ***REMOVED***
		var resp *loggingpb.ListMonitoredResourceDescriptorsResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 ***REMOVED***
			req.PageSize = math.MaxInt32
		***REMOVED*** else ***REMOVED***
			req.PageSize = int32(pageSize)
		***REMOVED***
		err := gax.Invoke(ctx, func(ctx context.Context) error ***REMOVED***
			var err error
			resp, err = c.client.ListMonitoredResourceDescriptors(ctx, req)
			return err
		***REMOVED***, c.CallOptions.ListMonitoredResourceDescriptors...)
		if err != nil ***REMOVED***
			return nil, "", err
		***REMOVED***
		return resp.ResourceDescriptors, resp.NextPageToken, nil
	***REMOVED***
	fetch := func(pageSize int, pageToken string) (string, error) ***REMOVED***
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		it.items = append(it.items, items...)
		return nextPageToken, nil
	***REMOVED***
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	return it
***REMOVED***

// LogEntryIterator manages a stream of *loggingpb.LogEntry.
type LogEntryIterator struct ***REMOVED***
	items    []*loggingpb.LogEntry
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*loggingpb.LogEntry, nextPageToken string, err error)
***REMOVED***

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *LogEntryIterator) PageInfo() *iterator.PageInfo ***REMOVED***
	return it.pageInfo
***REMOVED***

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *LogEntryIterator) Next() (*loggingpb.LogEntry, error) ***REMOVED***
	var item *loggingpb.LogEntry
	if err := it.nextFunc(); err != nil ***REMOVED***
		return item, err
	***REMOVED***
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
***REMOVED***

func (it *LogEntryIterator) bufLen() int ***REMOVED***
	return len(it.items)
***REMOVED***

func (it *LogEntryIterator) takeBuf() interface***REMOVED******REMOVED*** ***REMOVED***
	b := it.items
	it.items = nil
	return b
***REMOVED***

// MonitoredResourceDescriptorIterator manages a stream of *monitoredrespb.MonitoredResourceDescriptor.
type MonitoredResourceDescriptorIterator struct ***REMOVED***
	items    []*monitoredrespb.MonitoredResourceDescriptor
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*monitoredrespb.MonitoredResourceDescriptor, nextPageToken string, err error)
***REMOVED***

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *MonitoredResourceDescriptorIterator) PageInfo() *iterator.PageInfo ***REMOVED***
	return it.pageInfo
***REMOVED***

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *MonitoredResourceDescriptorIterator) Next() (*monitoredrespb.MonitoredResourceDescriptor, error) ***REMOVED***
	var item *monitoredrespb.MonitoredResourceDescriptor
	if err := it.nextFunc(); err != nil ***REMOVED***
		return item, err
	***REMOVED***
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
***REMOVED***

func (it *MonitoredResourceDescriptorIterator) bufLen() int ***REMOVED***
	return len(it.items)
***REMOVED***

func (it *MonitoredResourceDescriptorIterator) takeBuf() interface***REMOVED******REMOVED*** ***REMOVED***
	b := it.items
	it.items = nil
	return b
***REMOVED***
