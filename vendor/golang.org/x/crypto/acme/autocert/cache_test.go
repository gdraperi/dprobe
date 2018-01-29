// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autocert

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// make sure DirCache satisfies Cache interface
var _ Cache = DirCache("/")

func TestDirCache(t *testing.T) ***REMOVED***
	dir, err := ioutil.TempDir("", "autocert")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir)
	dir = filepath.Join(dir, "certs") // a nonexistent dir
	cache := DirCache(dir)
	ctx := context.Background()

	// test cache miss
	if _, err := cache.Get(ctx, "nonexistent"); err != ErrCacheMiss ***REMOVED***
		t.Errorf("get: %v; want ErrCacheMiss", err)
	***REMOVED***

	// test put/get
	b1 := []byte***REMOVED***1***REMOVED***
	if err := cache.Put(ctx, "dummy", b1); err != nil ***REMOVED***
		t.Fatalf("put: %v", err)
	***REMOVED***
	b2, err := cache.Get(ctx, "dummy")
	if err != nil ***REMOVED***
		t.Fatalf("get: %v", err)
	***REMOVED***
	if !reflect.DeepEqual(b1, b2) ***REMOVED***
		t.Errorf("b1 = %v; want %v", b1, b2)
	***REMOVED***
	name := filepath.Join(dir, "dummy")
	if _, err := os.Stat(name); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	// test delete
	if err := cache.Delete(ctx, "dummy"); err != nil ***REMOVED***
		t.Fatalf("delete: %v", err)
	***REMOVED***
	if _, err := cache.Get(ctx, "dummy"); err != ErrCacheMiss ***REMOVED***
		t.Errorf("get: %v; want ErrCacheMiss", err)
	***REMOVED***
***REMOVED***
