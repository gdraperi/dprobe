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
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var (
	configParentPathTemplate = gax.MustCompilePathTemplate("projects/***REMOVED***project***REMOVED***")
	configSinkPathTemplate   = gax.MustCompilePathTemplate("projects/***REMOVED***project***REMOVED***/sinks/***REMOVED***sink***REMOVED***")
)

// ConfigCallOptions contains the retry settings for each method of ConfigClient.
type ConfigCallOptions struct ***REMOVED***
	ListSinks  []gax.CallOption
	GetSink    []gax.CallOption
	CreateSink []gax.CallOption
	UpdateSink []gax.CallOption
	DeleteSink []gax.CallOption
***REMOVED***

func defaultConfigClientOptions() []option.ClientOption ***REMOVED***
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

func defaultConfigCallOptions() *ConfigCallOptions ***REMOVED***
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
	***REMOVED***
	return &ConfigCallOptions***REMOVED***
		ListSinks:  retry[[2]string***REMOVED***"default", "idempotent"***REMOVED***],
		GetSink:    retry[[2]string***REMOVED***"default", "idempotent"***REMOVED***],
		CreateSink: retry[[2]string***REMOVED***"default", "non_idempotent"***REMOVED***],
		UpdateSink: retry[[2]string***REMOVED***"default", "non_idempotent"***REMOVED***],
		DeleteSink: retry[[2]string***REMOVED***"default", "idempotent"***REMOVED***],
	***REMOVED***
***REMOVED***

// ConfigClient is a client for interacting with Stackdriver Logging API.
type ConfigClient struct ***REMOVED***
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	configClient loggingpb.ConfigServiceV2Client

	// The call options for this service.
	CallOptions *ConfigCallOptions

	// The metadata to be sent with each request.
	metadata metadata.MD
***REMOVED***

// NewConfigClient creates a new config service v2 client.
//
// Service for configuring sinks used to export log entries outside of
// Stackdriver Logging.
func NewConfigClient(ctx context.Context, opts ...option.ClientOption) (*ConfigClient, error) ***REMOVED***
	conn, err := transport.DialGRPC(ctx, append(defaultConfigClientOptions(), opts...)...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c := &ConfigClient***REMOVED***
		conn:        conn,
		CallOptions: defaultConfigCallOptions(),

		configClient: loggingpb.NewConfigServiceV2Client(conn),
	***REMOVED***
	c.SetGoogleClientInfo("gax", gax.Version)
	return c, nil
***REMOVED***

// Connection returns the client's connection to the API service.
func (c *ConfigClient) Connection() *grpc.ClientConn ***REMOVED***
	return c.conn
***REMOVED***

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *ConfigClient) Close() error ***REMOVED***
	return c.conn.Close()
***REMOVED***

// SetGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *ConfigClient) SetGoogleClientInfo(name, version string) ***REMOVED***
	goVersion := strings.Replace(runtime.Version(), " ", "_", -1)
	v := fmt.Sprintf("%s/%s %s gax/%s go/%s", name, version, gapicNameVersion, gax.Version, goVersion)
	c.metadata = metadata.Pairs("x-goog-api-client", v)
***REMOVED***

// ConfigParentPath returns the path for the parent resource.
func ConfigParentPath(project string) string ***REMOVED***
	path, err := configParentPathTemplate.Render(map[string]string***REMOVED***
		"project": project,
	***REMOVED***)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return path
***REMOVED***

// ConfigSinkPath returns the path for the sink resource.
func ConfigSinkPath(project, sink string) string ***REMOVED***
	path, err := configSinkPathTemplate.Render(map[string]string***REMOVED***
		"project": project,
		"sink":    sink,
	***REMOVED***)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return path
***REMOVED***

// ListSinks lists sinks.
func (c *ConfigClient) ListSinks(ctx context.Context, req *loggingpb.ListSinksRequest) *LogSinkIterator ***REMOVED***
	md, _ := metadata.FromContext(ctx)
	ctx = metadata.NewContext(ctx, metadata.Join(md, c.metadata))
	it := &LogSinkIterator***REMOVED******REMOVED***
	it.InternalFetch = func(pageSize int, pageToken string) ([]*loggingpb.LogSink, string, error) ***REMOVED***
		var resp *loggingpb.ListSinksResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 ***REMOVED***
			req.PageSize = math.MaxInt32
		***REMOVED*** else ***REMOVED***
			req.PageSize = int32(pageSize)
		***REMOVED***
		err := gax.Invoke(ctx, func(ctx context.Context) error ***REMOVED***
			var err error
			resp, err = c.configClient.ListSinks(ctx, req)
			return err
		***REMOVED***, c.CallOptions.ListSinks...)
		if err != nil ***REMOVED***
			return nil, "", err
		***REMOVED***
		return resp.Sinks, resp.NextPageToken, nil
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

// GetSink gets a sink.
func (c *ConfigClient) GetSink(ctx context.Context, req *loggingpb.GetSinkRequest) (*loggingpb.LogSink, error) ***REMOVED***
	md, _ := metadata.FromContext(ctx)
	ctx = metadata.NewContext(ctx, metadata.Join(md, c.metadata))
	var resp *loggingpb.LogSink
	err := gax.Invoke(ctx, func(ctx context.Context) error ***REMOVED***
		var err error
		resp, err = c.configClient.GetSink(ctx, req)
		return err
	***REMOVED***, c.CallOptions.GetSink...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp, nil
***REMOVED***

// CreateSink creates a sink.
func (c *ConfigClient) CreateSink(ctx context.Context, req *loggingpb.CreateSinkRequest) (*loggingpb.LogSink, error) ***REMOVED***
	md, _ := metadata.FromContext(ctx)
	ctx = metadata.NewContext(ctx, metadata.Join(md, c.metadata))
	var resp *loggingpb.LogSink
	err := gax.Invoke(ctx, func(ctx context.Context) error ***REMOVED***
		var err error
		resp, err = c.configClient.CreateSink(ctx, req)
		return err
	***REMOVED***, c.CallOptions.CreateSink...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp, nil
***REMOVED***

// UpdateSink updates or creates a sink.
func (c *ConfigClient) UpdateSink(ctx context.Context, req *loggingpb.UpdateSinkRequest) (*loggingpb.LogSink, error) ***REMOVED***
	md, _ := metadata.FromContext(ctx)
	ctx = metadata.NewContext(ctx, metadata.Join(md, c.metadata))
	var resp *loggingpb.LogSink
	err := gax.Invoke(ctx, func(ctx context.Context) error ***REMOVED***
		var err error
		resp, err = c.configClient.UpdateSink(ctx, req)
		return err
	***REMOVED***, c.CallOptions.UpdateSink...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return resp, nil
***REMOVED***

// DeleteSink deletes a sink.
func (c *ConfigClient) DeleteSink(ctx context.Context, req *loggingpb.DeleteSinkRequest) error ***REMOVED***
	md, _ := metadata.FromContext(ctx)
	ctx = metadata.NewContext(ctx, metadata.Join(md, c.metadata))
	err := gax.Invoke(ctx, func(ctx context.Context) error ***REMOVED***
		var err error
		_, err = c.configClient.DeleteSink(ctx, req)
		return err
	***REMOVED***, c.CallOptions.DeleteSink...)
	return err
***REMOVED***

// LogSinkIterator manages a stream of *loggingpb.LogSink.
type LogSinkIterator struct ***REMOVED***
	items    []*loggingpb.LogSink
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*loggingpb.LogSink, nextPageToken string, err error)
***REMOVED***

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *LogSinkIterator) PageInfo() *iterator.PageInfo ***REMOVED***
	return it.pageInfo
***REMOVED***

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *LogSinkIterator) Next() (*loggingpb.LogSink, error) ***REMOVED***
	var item *loggingpb.LogSink
	if err := it.nextFunc(); err != nil ***REMOVED***
		return item, err
	***REMOVED***
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
***REMOVED***

func (it *LogSinkIterator) bufLen() int ***REMOVED***
	return len(it.items)
***REMOVED***

func (it *LogSinkIterator) takeBuf() interface***REMOVED******REMOVED*** ***REMOVED***
	b := it.items
	it.items = nil
	return b
***REMOVED***
