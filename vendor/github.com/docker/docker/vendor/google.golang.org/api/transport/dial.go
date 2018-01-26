// Copyright 2015 Google Inc. All Rights Reserved.
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

// Package transport supports network connections to HTTP and GRPC servers.
// This package is not intended for use by end developers. Use the
// google.golang.org/api/option package to configure API clients.
package transport

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	gtransport "google.golang.org/api/googleapi/transport"
	"google.golang.org/api/internal"
	"google.golang.org/api/option"
)

// NewHTTPClient returns an HTTP client for use communicating with a Google cloud
// service, configured with the given ClientOptions. It also returns the endpoint
// for the service as specified in the options.
func NewHTTPClient(ctx context.Context, opts ...option.ClientOption) (*http.Client, string, error) ***REMOVED***
	var o internal.DialSettings
	for _, opt := range opts ***REMOVED***
		opt.Apply(&o)
	***REMOVED***
	if o.GRPCConn != nil ***REMOVED***
		return nil, "", errors.New("unsupported gRPC connection specified")
	***REMOVED***
	// TODO(djd): Set UserAgent on all outgoing requests.
	if o.HTTPClient != nil ***REMOVED***
		return o.HTTPClient, o.Endpoint, nil
	***REMOVED***
	if o.APIKey != "" ***REMOVED***
		hc := &http.Client***REMOVED***
			Transport: &gtransport.APIKey***REMOVED***
				Key:       o.APIKey,
				Transport: http.DefaultTransport,
			***REMOVED***,
		***REMOVED***
		return hc, o.Endpoint, nil
	***REMOVED***
	if o.ServiceAccountJSONFilename != "" ***REMOVED***
		ts, err := serviceAcctTokenSource(ctx, o.ServiceAccountJSONFilename, o.Scopes...)
		if err != nil ***REMOVED***
			return nil, "", err
		***REMOVED***
		o.TokenSource = ts
	***REMOVED***
	if o.TokenSource == nil ***REMOVED***
		var err error
		o.TokenSource, err = google.DefaultTokenSource(ctx, o.Scopes...)
		if err != nil ***REMOVED***
			return nil, "", fmt.Errorf("google.DefaultTokenSource: %v", err)
		***REMOVED***
	***REMOVED***
	return oauth2.NewClient(ctx, o.TokenSource), o.Endpoint, nil
***REMOVED***

// Set at init time by dial_appengine.go. If nil, we're not on App Engine.
var appengineDialerHook func(context.Context) grpc.DialOption

// DialGRPC returns a GRPC connection for use communicating with a Google cloud
// service, configured with the given ClientOptions.
func DialGRPC(ctx context.Context, opts ...option.ClientOption) (*grpc.ClientConn, error) ***REMOVED***
	var o internal.DialSettings
	for _, opt := range opts ***REMOVED***
		opt.Apply(&o)
	***REMOVED***
	if o.HTTPClient != nil ***REMOVED***
		return nil, errors.New("unsupported HTTP client specified")
	***REMOVED***
	if o.GRPCConn != nil ***REMOVED***
		return o.GRPCConn, nil
	***REMOVED***
	if o.ServiceAccountJSONFilename != "" ***REMOVED***
		ts, err := serviceAcctTokenSource(ctx, o.ServiceAccountJSONFilename, o.Scopes...)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		o.TokenSource = ts
	***REMOVED***
	if o.TokenSource == nil ***REMOVED***
		var err error
		o.TokenSource, err = google.DefaultTokenSource(ctx, o.Scopes...)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("google.DefaultTokenSource: %v", err)
		***REMOVED***
	***REMOVED***
	grpcOpts := []grpc.DialOption***REMOVED***
		grpc.WithPerRPCCredentials(oauth.TokenSource***REMOVED***o.TokenSource***REMOVED***),
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	***REMOVED***
	if appengineDialerHook != nil ***REMOVED***
		// Use the Socket API on App Engine.
		grpcOpts = append(grpcOpts, appengineDialerHook(ctx))
	***REMOVED***
	grpcOpts = append(grpcOpts, o.GRPCDialOpts...)
	if o.UserAgent != "" ***REMOVED***
		grpcOpts = append(grpcOpts, grpc.WithUserAgent(o.UserAgent))
	***REMOVED***
	return grpc.DialContext(ctx, o.Endpoint, grpcOpts...)
***REMOVED***

func serviceAcctTokenSource(ctx context.Context, filename string, scope ...string) (oauth2.TokenSource, error) ***REMOVED***
	data, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("cannot read service account file: %v", err)
	***REMOVED***
	cfg, err := google.JWTConfigFromJSON(data, scope...)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	***REMOVED***
	return cfg.TokenSource(ctx), nil
***REMOVED***
