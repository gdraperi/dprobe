// Copyright 2016, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package gax

import (
	"time"

	"golang.org/x/net/context"
)

// A user defined call stub.
type APICall func(context.Context) error

// Invoke calls the given APICall,
// performing retries as specified by opts, if any.
func Invoke(ctx context.Context, call APICall, opts ...CallOption) error ***REMOVED***
	var settings CallSettings
	for _, opt := range opts ***REMOVED***
		opt.Resolve(&settings)
	***REMOVED***
	return invoke(ctx, call, settings, Sleep)
***REMOVED***

// Sleep is similar to time.Sleep, but it can be interrupted by ctx.Done() closing.
// If interrupted, Sleep returns ctx.Err().
func Sleep(ctx context.Context, d time.Duration) error ***REMOVED***
	t := time.NewTimer(d)
	select ***REMOVED***
	case <-ctx.Done():
		t.Stop()
		return ctx.Err()
	case <-t.C:
		return nil
	***REMOVED***
***REMOVED***

type sleeper func(ctx context.Context, d time.Duration) error

// invoke implements Invoke, taking an additional sleeper argument for testing.
func invoke(ctx context.Context, call APICall, settings CallSettings, sp sleeper) error ***REMOVED***
	var retryer Retryer
	for ***REMOVED***
		err := call(ctx)
		if err == nil ***REMOVED***
			return nil
		***REMOVED***
		if settings.Retry == nil ***REMOVED***
			return err
		***REMOVED***
		if retryer == nil ***REMOVED***
			if r := settings.Retry(); r != nil ***REMOVED***
				retryer = r
			***REMOVED*** else ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if d, ok := retryer.Retry(err); !ok ***REMOVED***
			return err
		***REMOVED*** else if err = sp(ctx, d); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
***REMOVED***
