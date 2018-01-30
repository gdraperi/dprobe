// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.7

package webdav

import (
	"net/http"

	"golang.org/x/net/context"
)

func getContext(r *http.Request) context.Context ***REMOVED***
	return context.Background()
***REMOVED***
