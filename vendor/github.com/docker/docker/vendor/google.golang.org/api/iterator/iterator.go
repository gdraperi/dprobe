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

// Package iterator provides support for standard Google API iterators.
// See https://github.com/GoogleCloudPlatform/gcloud-golang/wiki/Iterator-Guidelines.
package iterator

import (
	"errors"
	"fmt"
	"reflect"
)

// Done is returned by an iterator's Next method when the iteration is
// complete; when there are no more items to return.
var Done = errors.New("no more items in iterator")

// We don't support mixed calls to Next and NextPage because they play
// with the paging state in incompatible ways.
var errMixed = errors.New("iterator: Next and NextPage called on same iterator")

// PageInfo contains information about an iterator's paging state.
type PageInfo struct ***REMOVED***
	// Token is the token used to retrieve the next page of items from the
	// API. You may set Token immediately after creating an iterator to
	// begin iteration at a particular point. If Token is the empty string,
	// the iterator will begin with the first eligible item.
	//
	// The result of setting Token after the first call to Next is undefined.
	//
	// After the underlying API method is called to retrieve a page of items,
	// Token is set to the next-page token in the response.
	Token string

	// MaxSize is the maximum number of items returned by a call to the API.
	// Set MaxSize as a hint to optimize the buffering behavior of the iterator.
	// If zero, the page size is determined by the underlying service.
	//
	// Use Pager to retrieve a page of a specific, exact size.
	MaxSize int

	// The error state of the iterator. Manipulated by PageInfo.next and Pager.
	// This is a latch: it starts as nil, and once set should never change.
	err error

	// If true, no more calls to fetch should be made. Set to true when fetch
	// returns an empty page token. The iterator is Done when this is true AND
	// the buffer is empty.
	atEnd bool

	// Function that fetches a page from the underlying service. It should pass
	// the pageSize and pageToken arguments to the service, fill the buffer
	// with the results from the call, and return the next-page token returned
	// by the service. The function must not remove any existing items from the
	// buffer. If the underlying RPC takes an int32 page size, pageSize should
	// be silently truncated.
	fetch func(pageSize int, pageToken string) (nextPageToken string, err error)

	// Function that clears the iterator's buffer, returning any currently buffered items.
	bufLen func() int

	// Function that returns the buffer, after setting the buffer variable to nil.
	takeBuf func() interface***REMOVED******REMOVED***

	// Set to true on first call to PageInfo.next or Pager.NextPage. Used to check
	// for calls to both Next and NextPage with the same iterator.
	nextCalled, nextPageCalled bool
***REMOVED***

// NewPageInfo exposes internals for iterator implementations.
// It is not a stable interface.
var NewPageInfo = newPageInfo

// If an iterator can support paging, its iterator-creating method should call
// this (via the NewPageInfo variable above).
//
// The fetch, bufLen and takeBuf arguments provide access to the
// iterator's internal slice of buffered items. They behave as described in
// PageInfo, above.
//
// The return value is the PageInfo.next method bound to the returned PageInfo value.
// (Returning it avoids exporting PageInfo.next.)
func newPageInfo(fetch func(int, string) (string, error), bufLen func() int, takeBuf func() interface***REMOVED******REMOVED***) (*PageInfo, func() error) ***REMOVED***
	pi := &PageInfo***REMOVED***
		fetch:   fetch,
		bufLen:  bufLen,
		takeBuf: takeBuf,
	***REMOVED***
	return pi, pi.next
***REMOVED***

// Remaining returns the number of items available before the iterator makes another API call.
func (pi *PageInfo) Remaining() int ***REMOVED*** return pi.bufLen() ***REMOVED***

// next provides support for an iterator's Next function. An iterator's Next
// should return the error returned by next if non-nil; else it can assume
// there is at least one item in its buffer, and it should return that item and
// remove it from the buffer.
func (pi *PageInfo) next() error ***REMOVED***
	pi.nextCalled = true
	if pi.err != nil ***REMOVED*** // Once we get an error, always return it.
		// TODO(jba): fix so users can retry on transient errors? Probably not worth it.
		return pi.err
	***REMOVED***
	if pi.nextPageCalled ***REMOVED***
		pi.err = errMixed
		return pi.err
	***REMOVED***
	// Loop until we get some items or reach the end.
	for pi.bufLen() == 0 && !pi.atEnd ***REMOVED***
		if err := pi.fill(pi.MaxSize); err != nil ***REMOVED***
			pi.err = err
			return pi.err
		***REMOVED***
		if pi.Token == "" ***REMOVED***
			pi.atEnd = true
		***REMOVED***
	***REMOVED***
	// Either the buffer is non-empty or pi.atEnd is true (or both).
	if pi.bufLen() == 0 ***REMOVED***
		// The buffer is empty and pi.atEnd is true, i.e. the service has no
		// more items.
		pi.err = Done
	***REMOVED***
	return pi.err
***REMOVED***

// Call the service to fill the buffer, using size and pi.Token. Set pi.Token to the
// next-page token returned by the call.
// If fill returns a non-nil error, the buffer will be empty.
func (pi *PageInfo) fill(size int) error ***REMOVED***
	tok, err := pi.fetch(size, pi.Token)
	if err != nil ***REMOVED***
		pi.takeBuf() // clear the buffer
		return err
	***REMOVED***
	pi.Token = tok
	return nil
***REMOVED***

// Pageable is implemented by iterators that support paging.
type Pageable interface ***REMOVED***
	// PageInfo returns paging information associated with the iterator.
	PageInfo() *PageInfo
***REMOVED***

// Pager supports retrieving iterator items a page at a time.
type Pager struct ***REMOVED***
	pageInfo *PageInfo
	pageSize int
***REMOVED***

// NewPager returns a pager that uses iter. Calls to its NextPage method will
// obtain exactly pageSize items, unless fewer remain. The pageToken argument
// indicates where to start the iteration. Pass the empty string to start at
// the beginning, or pass a token retrieved from a call to Pager.NextPage.
//
// If you use an iterator with a Pager, you must not call Next on the iterator.
func NewPager(iter Pageable, pageSize int, pageToken string) *Pager ***REMOVED***
	p := &Pager***REMOVED***
		pageInfo: iter.PageInfo(),
		pageSize: pageSize,
	***REMOVED***
	p.pageInfo.Token = pageToken
	if pageSize <= 0 ***REMOVED***
		p.pageInfo.err = errors.New("iterator: page size must be positive")
	***REMOVED***
	return p
***REMOVED***

// NextPage retrieves a sequence of items from the iterator and appends them
// to slicep, which must be a pointer to a slice of the iterator's item type.
// Exactly p.pageSize items will be appended, unless fewer remain.
//
// The first return value is the page token to use for the next page of items.
// If empty, there are no more pages. Aside from checking for the end of the
// iteration, the returned page token is only needed if the iteration is to be
// resumed a later time, in another context (possibly another process).
//
// The second return value is non-nil if an error occurred. It will never be
// the special iterator sentinel value Done. To recognize the end of the
// iteration, compare nextPageToken to the empty string.
//
// It is possible for NextPage to return a single zero-length page along with
// an empty page token when there are no more items in the iteration.
func (p *Pager) NextPage(slicep interface***REMOVED******REMOVED***) (nextPageToken string, err error) ***REMOVED***
	p.pageInfo.nextPageCalled = true
	if p.pageInfo.err != nil ***REMOVED***
		return "", p.pageInfo.err
	***REMOVED***
	if p.pageInfo.nextCalled ***REMOVED***
		p.pageInfo.err = errMixed
		return "", p.pageInfo.err
	***REMOVED***
	if p.pageInfo.bufLen() > 0 ***REMOVED***
		return "", errors.New("must call NextPage with an empty buffer")
	***REMOVED***
	// The buffer must be empty here, so takeBuf is a no-op. We call it just to get
	// the buffer's type.
	wantSliceType := reflect.PtrTo(reflect.ValueOf(p.pageInfo.takeBuf()).Type())
	if slicep == nil ***REMOVED***
		return "", errors.New("nil passed to Pager.NextPage")
	***REMOVED***
	vslicep := reflect.ValueOf(slicep)
	if vslicep.Type() != wantSliceType ***REMOVED***
		return "", fmt.Errorf("slicep should be of type %s, got %T", wantSliceType, slicep)
	***REMOVED***
	for p.pageInfo.bufLen() < p.pageSize ***REMOVED***
		if err := p.pageInfo.fill(p.pageSize - p.pageInfo.bufLen()); err != nil ***REMOVED***
			p.pageInfo.err = err
			return "", p.pageInfo.err
		***REMOVED***
		if p.pageInfo.Token == "" ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	e := vslicep.Elem()
	e.Set(reflect.AppendSlice(e, reflect.ValueOf(p.pageInfo.takeBuf())))
	return p.pageInfo.Token, nil
***REMOVED***
