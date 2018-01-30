// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9,go1.7

package ctxhttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
)

func TestGo17Context(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		io.WriteString(w, "ok")
	***REMOVED***))
	defer ts.Close()
	ctx := context.Background()
	resp, err := Get(ctx, http.DefaultClient, ts.URL)
	if resp == nil || err != nil ***REMOVED***
		t.Fatalf("error received from client: %v %v", err, resp)
	***REMOVED***
	resp.Body.Close()
***REMOVED***
