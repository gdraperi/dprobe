// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package webdav provides a WebDAV server implementation.
package webdav // import "golang.org/x/net/webdav"

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Handler struct ***REMOVED***
	// Prefix is the URL path prefix to strip from WebDAV resource paths.
	Prefix string
	// FileSystem is the virtual file system.
	FileSystem FileSystem
	// LockSystem is the lock management system.
	LockSystem LockSystem
	// Logger is an optional error logger. If non-nil, it will be called
	// for all HTTP requests.
	Logger func(*http.Request, error)
***REMOVED***

func (h *Handler) stripPrefix(p string) (string, int, error) ***REMOVED***
	if h.Prefix == "" ***REMOVED***
		return p, http.StatusOK, nil
	***REMOVED***
	if r := strings.TrimPrefix(p, h.Prefix); len(r) < len(p) ***REMOVED***
		return r, http.StatusOK, nil
	***REMOVED***
	return p, http.StatusNotFound, errPrefixMismatch
***REMOVED***

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	status, err := http.StatusBadRequest, errUnsupportedMethod
	if h.FileSystem == nil ***REMOVED***
		status, err = http.StatusInternalServerError, errNoFileSystem
	***REMOVED*** else if h.LockSystem == nil ***REMOVED***
		status, err = http.StatusInternalServerError, errNoLockSystem
	***REMOVED*** else ***REMOVED***
		switch r.Method ***REMOVED***
		case "OPTIONS":
			status, err = h.handleOptions(w, r)
		case "GET", "HEAD", "POST":
			status, err = h.handleGetHeadPost(w, r)
		case "DELETE":
			status, err = h.handleDelete(w, r)
		case "PUT":
			status, err = h.handlePut(w, r)
		case "MKCOL":
			status, err = h.handleMkcol(w, r)
		case "COPY", "MOVE":
			status, err = h.handleCopyMove(w, r)
		case "LOCK":
			status, err = h.handleLock(w, r)
		case "UNLOCK":
			status, err = h.handleUnlock(w, r)
		case "PROPFIND":
			status, err = h.handlePropfind(w, r)
		case "PROPPATCH":
			status, err = h.handleProppatch(w, r)
		***REMOVED***
	***REMOVED***

	if status != 0 ***REMOVED***
		w.WriteHeader(status)
		if status != http.StatusNoContent ***REMOVED***
			w.Write([]byte(StatusText(status)))
		***REMOVED***
	***REMOVED***
	if h.Logger != nil ***REMOVED***
		h.Logger(r, err)
	***REMOVED***
***REMOVED***

func (h *Handler) lock(now time.Time, root string) (token string, status int, err error) ***REMOVED***
	token, err = h.LockSystem.Create(now, LockDetails***REMOVED***
		Root:      root,
		Duration:  infiniteTimeout,
		ZeroDepth: true,
	***REMOVED***)
	if err != nil ***REMOVED***
		if err == ErrLocked ***REMOVED***
			return "", StatusLocked, err
		***REMOVED***
		return "", http.StatusInternalServerError, err
	***REMOVED***
	return token, 0, nil
***REMOVED***

func (h *Handler) confirmLocks(r *http.Request, src, dst string) (release func(), status int, err error) ***REMOVED***
	hdr := r.Header.Get("If")
	if hdr == "" ***REMOVED***
		// An empty If header means that the client hasn't previously created locks.
		// Even if this client doesn't care about locks, we still need to check that
		// the resources aren't locked by another client, so we create temporary
		// locks that would conflict with another client's locks. These temporary
		// locks are unlocked at the end of the HTTP request.
		now, srcToken, dstToken := time.Now(), "", ""
		if src != "" ***REMOVED***
			srcToken, status, err = h.lock(now, src)
			if err != nil ***REMOVED***
				return nil, status, err
			***REMOVED***
		***REMOVED***
		if dst != "" ***REMOVED***
			dstToken, status, err = h.lock(now, dst)
			if err != nil ***REMOVED***
				if srcToken != "" ***REMOVED***
					h.LockSystem.Unlock(now, srcToken)
				***REMOVED***
				return nil, status, err
			***REMOVED***
		***REMOVED***

		return func() ***REMOVED***
			if dstToken != "" ***REMOVED***
				h.LockSystem.Unlock(now, dstToken)
			***REMOVED***
			if srcToken != "" ***REMOVED***
				h.LockSystem.Unlock(now, srcToken)
			***REMOVED***
		***REMOVED***, 0, nil
	***REMOVED***

	ih, ok := parseIfHeader(hdr)
	if !ok ***REMOVED***
		return nil, http.StatusBadRequest, errInvalidIfHeader
	***REMOVED***
	// ih is a disjunction (OR) of ifLists, so any ifList will do.
	for _, l := range ih.lists ***REMOVED***
		lsrc := l.resourceTag
		if lsrc == "" ***REMOVED***
			lsrc = src
		***REMOVED*** else ***REMOVED***
			u, err := url.Parse(lsrc)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			if u.Host != r.Host ***REMOVED***
				continue
			***REMOVED***
			lsrc, status, err = h.stripPrefix(u.Path)
			if err != nil ***REMOVED***
				return nil, status, err
			***REMOVED***
		***REMOVED***
		release, err = h.LockSystem.Confirm(time.Now(), lsrc, dst, l.conditions...)
		if err == ErrConfirmationFailed ***REMOVED***
			continue
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, http.StatusInternalServerError, err
		***REMOVED***
		return release, 0, nil
	***REMOVED***
	// Section 10.4.1 says that "If this header is evaluated and all state lists
	// fail, then the request must fail with a 412 (Precondition Failed) status."
	// We follow the spec even though the cond_put_corrupt_token test case from
	// the litmus test warns on seeing a 412 instead of a 423 (Locked).
	return nil, http.StatusPreconditionFailed, ErrLocked
***REMOVED***

func (h *Handler) handleOptions(w http.ResponseWriter, r *http.Request) (status int, err error) ***REMOVED***
	reqPath, status, err := h.stripPrefix(r.URL.Path)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	ctx := getContext(r)
	allow := "OPTIONS, LOCK, PUT, MKCOL"
	if fi, err := h.FileSystem.Stat(ctx, reqPath); err == nil ***REMOVED***
		if fi.IsDir() ***REMOVED***
			allow = "OPTIONS, LOCK, DELETE, PROPPATCH, COPY, MOVE, UNLOCK, PROPFIND"
		***REMOVED*** else ***REMOVED***
			allow = "OPTIONS, LOCK, GET, HEAD, POST, DELETE, PROPPATCH, COPY, MOVE, UNLOCK, PROPFIND, PUT"
		***REMOVED***
	***REMOVED***
	w.Header().Set("Allow", allow)
	// http://www.webdav.org/specs/rfc4918.html#dav.compliance.classes
	w.Header().Set("DAV", "1, 2")
	// http://msdn.microsoft.com/en-au/library/cc250217.aspx
	w.Header().Set("MS-Author-Via", "DAV")
	return 0, nil
***REMOVED***

func (h *Handler) handleGetHeadPost(w http.ResponseWriter, r *http.Request) (status int, err error) ***REMOVED***
	reqPath, status, err := h.stripPrefix(r.URL.Path)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	// TODO: check locks for read-only access??
	ctx := getContext(r)
	f, err := h.FileSystem.OpenFile(ctx, reqPath, os.O_RDONLY, 0)
	if err != nil ***REMOVED***
		return http.StatusNotFound, err
	***REMOVED***
	defer f.Close()
	fi, err := f.Stat()
	if err != nil ***REMOVED***
		return http.StatusNotFound, err
	***REMOVED***
	if fi.IsDir() ***REMOVED***
		return http.StatusMethodNotAllowed, nil
	***REMOVED***
	etag, err := findETag(ctx, h.FileSystem, h.LockSystem, reqPath, fi)
	if err != nil ***REMOVED***
		return http.StatusInternalServerError, err
	***REMOVED***
	w.Header().Set("ETag", etag)
	// Let ServeContent determine the Content-Type header.
	http.ServeContent(w, r, reqPath, fi.ModTime(), f)
	return 0, nil
***REMOVED***

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) (status int, err error) ***REMOVED***
	reqPath, status, err := h.stripPrefix(r.URL.Path)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	release, status, err := h.confirmLocks(r, reqPath, "")
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	defer release()

	ctx := getContext(r)

	// TODO: return MultiStatus where appropriate.

	// "godoc os RemoveAll" says that "If the path does not exist, RemoveAll
	// returns nil (no error)." WebDAV semantics are that it should return a
	// "404 Not Found". We therefore have to Stat before we RemoveAll.
	if _, err := h.FileSystem.Stat(ctx, reqPath); err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return http.StatusNotFound, err
		***REMOVED***
		return http.StatusMethodNotAllowed, err
	***REMOVED***
	if err := h.FileSystem.RemoveAll(ctx, reqPath); err != nil ***REMOVED***
		return http.StatusMethodNotAllowed, err
	***REMOVED***
	return http.StatusNoContent, nil
***REMOVED***

func (h *Handler) handlePut(w http.ResponseWriter, r *http.Request) (status int, err error) ***REMOVED***
	reqPath, status, err := h.stripPrefix(r.URL.Path)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	release, status, err := h.confirmLocks(r, reqPath, "")
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	defer release()
	// TODO(rost): Support the If-Match, If-None-Match headers? See bradfitz'
	// comments in http.checkEtag.
	ctx := getContext(r)

	f, err := h.FileSystem.OpenFile(ctx, reqPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil ***REMOVED***
		return http.StatusNotFound, err
	***REMOVED***
	_, copyErr := io.Copy(f, r.Body)
	fi, statErr := f.Stat()
	closeErr := f.Close()
	// TODO(rost): Returning 405 Method Not Allowed might not be appropriate.
	if copyErr != nil ***REMOVED***
		return http.StatusMethodNotAllowed, copyErr
	***REMOVED***
	if statErr != nil ***REMOVED***
		return http.StatusMethodNotAllowed, statErr
	***REMOVED***
	if closeErr != nil ***REMOVED***
		return http.StatusMethodNotAllowed, closeErr
	***REMOVED***
	etag, err := findETag(ctx, h.FileSystem, h.LockSystem, reqPath, fi)
	if err != nil ***REMOVED***
		return http.StatusInternalServerError, err
	***REMOVED***
	w.Header().Set("ETag", etag)
	return http.StatusCreated, nil
***REMOVED***

func (h *Handler) handleMkcol(w http.ResponseWriter, r *http.Request) (status int, err error) ***REMOVED***
	reqPath, status, err := h.stripPrefix(r.URL.Path)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	release, status, err := h.confirmLocks(r, reqPath, "")
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	defer release()

	ctx := getContext(r)

	if r.ContentLength > 0 ***REMOVED***
		return http.StatusUnsupportedMediaType, nil
	***REMOVED***
	if err := h.FileSystem.Mkdir(ctx, reqPath, 0777); err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return http.StatusConflict, err
		***REMOVED***
		return http.StatusMethodNotAllowed, err
	***REMOVED***
	return http.StatusCreated, nil
***REMOVED***

func (h *Handler) handleCopyMove(w http.ResponseWriter, r *http.Request) (status int, err error) ***REMOVED***
	hdr := r.Header.Get("Destination")
	if hdr == "" ***REMOVED***
		return http.StatusBadRequest, errInvalidDestination
	***REMOVED***
	u, err := url.Parse(hdr)
	if err != nil ***REMOVED***
		return http.StatusBadRequest, errInvalidDestination
	***REMOVED***
	if u.Host != r.Host ***REMOVED***
		return http.StatusBadGateway, errInvalidDestination
	***REMOVED***

	src, status, err := h.stripPrefix(r.URL.Path)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***

	dst, status, err := h.stripPrefix(u.Path)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***

	if dst == "" ***REMOVED***
		return http.StatusBadGateway, errInvalidDestination
	***REMOVED***
	if dst == src ***REMOVED***
		return http.StatusForbidden, errDestinationEqualsSource
	***REMOVED***

	ctx := getContext(r)

	if r.Method == "COPY" ***REMOVED***
		// Section 7.5.1 says that a COPY only needs to lock the destination,
		// not both destination and source. Strictly speaking, this is racy,
		// even though a COPY doesn't modify the source, if a concurrent
		// operation modifies the source. However, the litmus test explicitly
		// checks that COPYing a locked-by-another source is OK.
		release, status, err := h.confirmLocks(r, "", dst)
		if err != nil ***REMOVED***
			return status, err
		***REMOVED***
		defer release()

		// Section 9.8.3 says that "The COPY method on a collection without a Depth
		// header must act as if a Depth header with value "infinity" was included".
		depth := infiniteDepth
		if hdr := r.Header.Get("Depth"); hdr != "" ***REMOVED***
			depth = parseDepth(hdr)
			if depth != 0 && depth != infiniteDepth ***REMOVED***
				// Section 9.8.3 says that "A client may submit a Depth header on a
				// COPY on a collection with a value of "0" or "infinity"."
				return http.StatusBadRequest, errInvalidDepth
			***REMOVED***
		***REMOVED***
		return copyFiles(ctx, h.FileSystem, src, dst, r.Header.Get("Overwrite") != "F", depth, 0)
	***REMOVED***

	release, status, err := h.confirmLocks(r, src, dst)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	defer release()

	// Section 9.9.2 says that "The MOVE method on a collection must act as if
	// a "Depth: infinity" header was used on it. A client must not submit a
	// Depth header on a MOVE on a collection with any value but "infinity"."
	if hdr := r.Header.Get("Depth"); hdr != "" ***REMOVED***
		if parseDepth(hdr) != infiniteDepth ***REMOVED***
			return http.StatusBadRequest, errInvalidDepth
		***REMOVED***
	***REMOVED***
	return moveFiles(ctx, h.FileSystem, src, dst, r.Header.Get("Overwrite") == "T")
***REMOVED***

func (h *Handler) handleLock(w http.ResponseWriter, r *http.Request) (retStatus int, retErr error) ***REMOVED***
	duration, err := parseTimeout(r.Header.Get("Timeout"))
	if err != nil ***REMOVED***
		return http.StatusBadRequest, err
	***REMOVED***
	li, status, err := readLockInfo(r.Body)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***

	ctx := getContext(r)
	token, ld, now, created := "", LockDetails***REMOVED******REMOVED***, time.Now(), false
	if li == (lockInfo***REMOVED******REMOVED***) ***REMOVED***
		// An empty lockInfo means to refresh the lock.
		ih, ok := parseIfHeader(r.Header.Get("If"))
		if !ok ***REMOVED***
			return http.StatusBadRequest, errInvalidIfHeader
		***REMOVED***
		if len(ih.lists) == 1 && len(ih.lists[0].conditions) == 1 ***REMOVED***
			token = ih.lists[0].conditions[0].Token
		***REMOVED***
		if token == "" ***REMOVED***
			return http.StatusBadRequest, errInvalidLockToken
		***REMOVED***
		ld, err = h.LockSystem.Refresh(now, token, duration)
		if err != nil ***REMOVED***
			if err == ErrNoSuchLock ***REMOVED***
				return http.StatusPreconditionFailed, err
			***REMOVED***
			return http.StatusInternalServerError, err
		***REMOVED***

	***REMOVED*** else ***REMOVED***
		// Section 9.10.3 says that "If no Depth header is submitted on a LOCK request,
		// then the request MUST act as if a "Depth:infinity" had been submitted."
		depth := infiniteDepth
		if hdr := r.Header.Get("Depth"); hdr != "" ***REMOVED***
			depth = parseDepth(hdr)
			if depth != 0 && depth != infiniteDepth ***REMOVED***
				// Section 9.10.3 says that "Values other than 0 or infinity must not be
				// used with the Depth header on a LOCK method".
				return http.StatusBadRequest, errInvalidDepth
			***REMOVED***
		***REMOVED***
		reqPath, status, err := h.stripPrefix(r.URL.Path)
		if err != nil ***REMOVED***
			return status, err
		***REMOVED***
		ld = LockDetails***REMOVED***
			Root:      reqPath,
			Duration:  duration,
			OwnerXML:  li.Owner.InnerXML,
			ZeroDepth: depth == 0,
		***REMOVED***
		token, err = h.LockSystem.Create(now, ld)
		if err != nil ***REMOVED***
			if err == ErrLocked ***REMOVED***
				return StatusLocked, err
			***REMOVED***
			return http.StatusInternalServerError, err
		***REMOVED***
		defer func() ***REMOVED***
			if retErr != nil ***REMOVED***
				h.LockSystem.Unlock(now, token)
			***REMOVED***
		***REMOVED***()

		// Create the resource if it didn't previously exist.
		if _, err := h.FileSystem.Stat(ctx, reqPath); err != nil ***REMOVED***
			f, err := h.FileSystem.OpenFile(ctx, reqPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil ***REMOVED***
				// TODO: detect missing intermediate dirs and return http.StatusConflict?
				return http.StatusInternalServerError, err
			***REMOVED***
			f.Close()
			created = true
		***REMOVED***

		// http://www.webdav.org/specs/rfc4918.html#HEADER_Lock-Token says that the
		// Lock-Token value is a Coded-URL. We add angle brackets.
		w.Header().Set("Lock-Token", "<"+token+">")
	***REMOVED***

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	if created ***REMOVED***
		// This is "w.WriteHeader(http.StatusCreated)" and not "return
		// http.StatusCreated, nil" because we write our own (XML) response to w
		// and Handler.ServeHTTP would otherwise write "Created".
		w.WriteHeader(http.StatusCreated)
	***REMOVED***
	writeLockInfo(w, token, ld)
	return 0, nil
***REMOVED***

func (h *Handler) handleUnlock(w http.ResponseWriter, r *http.Request) (status int, err error) ***REMOVED***
	// http://www.webdav.org/specs/rfc4918.html#HEADER_Lock-Token says that the
	// Lock-Token value is a Coded-URL. We strip its angle brackets.
	t := r.Header.Get("Lock-Token")
	if len(t) < 2 || t[0] != '<' || t[len(t)-1] != '>' ***REMOVED***
		return http.StatusBadRequest, errInvalidLockToken
	***REMOVED***
	t = t[1 : len(t)-1]

	switch err = h.LockSystem.Unlock(time.Now(), t); err ***REMOVED***
	case nil:
		return http.StatusNoContent, err
	case ErrForbidden:
		return http.StatusForbidden, err
	case ErrLocked:
		return StatusLocked, err
	case ErrNoSuchLock:
		return http.StatusConflict, err
	default:
		return http.StatusInternalServerError, err
	***REMOVED***
***REMOVED***

func (h *Handler) handlePropfind(w http.ResponseWriter, r *http.Request) (status int, err error) ***REMOVED***
	reqPath, status, err := h.stripPrefix(r.URL.Path)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	ctx := getContext(r)
	fi, err := h.FileSystem.Stat(ctx, reqPath)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return http.StatusNotFound, err
		***REMOVED***
		return http.StatusMethodNotAllowed, err
	***REMOVED***
	depth := infiniteDepth
	if hdr := r.Header.Get("Depth"); hdr != "" ***REMOVED***
		depth = parseDepth(hdr)
		if depth == invalidDepth ***REMOVED***
			return http.StatusBadRequest, errInvalidDepth
		***REMOVED***
	***REMOVED***
	pf, status, err := readPropfind(r.Body)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***

	mw := multistatusWriter***REMOVED***w: w***REMOVED***

	walkFn := func(reqPath string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		var pstats []Propstat
		if pf.Propname != nil ***REMOVED***
			pnames, err := propnames(ctx, h.FileSystem, h.LockSystem, reqPath)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			pstat := Propstat***REMOVED***Status: http.StatusOK***REMOVED***
			for _, xmlname := range pnames ***REMOVED***
				pstat.Props = append(pstat.Props, Property***REMOVED***XMLName: xmlname***REMOVED***)
			***REMOVED***
			pstats = append(pstats, pstat)
		***REMOVED*** else if pf.Allprop != nil ***REMOVED***
			pstats, err = allprop(ctx, h.FileSystem, h.LockSystem, reqPath, pf.Prop)
		***REMOVED*** else ***REMOVED***
			pstats, err = props(ctx, h.FileSystem, h.LockSystem, reqPath, pf.Prop)
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return mw.write(makePropstatResponse(path.Join(h.Prefix, reqPath), pstats))
	***REMOVED***

	walkErr := walkFS(ctx, h.FileSystem, depth, reqPath, fi, walkFn)
	closeErr := mw.close()
	if walkErr != nil ***REMOVED***
		return http.StatusInternalServerError, walkErr
	***REMOVED***
	if closeErr != nil ***REMOVED***
		return http.StatusInternalServerError, closeErr
	***REMOVED***
	return 0, nil
***REMOVED***

func (h *Handler) handleProppatch(w http.ResponseWriter, r *http.Request) (status int, err error) ***REMOVED***
	reqPath, status, err := h.stripPrefix(r.URL.Path)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	release, status, err := h.confirmLocks(r, reqPath, "")
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	defer release()

	ctx := getContext(r)

	if _, err := h.FileSystem.Stat(ctx, reqPath); err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return http.StatusNotFound, err
		***REMOVED***
		return http.StatusMethodNotAllowed, err
	***REMOVED***
	patches, status, err := readProppatch(r.Body)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	pstats, err := patch(ctx, h.FileSystem, h.LockSystem, reqPath, patches)
	if err != nil ***REMOVED***
		return http.StatusInternalServerError, err
	***REMOVED***
	mw := multistatusWriter***REMOVED***w: w***REMOVED***
	writeErr := mw.write(makePropstatResponse(r.URL.Path, pstats))
	closeErr := mw.close()
	if writeErr != nil ***REMOVED***
		return http.StatusInternalServerError, writeErr
	***REMOVED***
	if closeErr != nil ***REMOVED***
		return http.StatusInternalServerError, closeErr
	***REMOVED***
	return 0, nil
***REMOVED***

func makePropstatResponse(href string, pstats []Propstat) *response ***REMOVED***
	resp := response***REMOVED***
		Href:     []string***REMOVED***(&url.URL***REMOVED***Path: href***REMOVED***).EscapedPath()***REMOVED***,
		Propstat: make([]propstat, 0, len(pstats)),
	***REMOVED***
	for _, p := range pstats ***REMOVED***
		var xmlErr *xmlError
		if p.XMLError != "" ***REMOVED***
			xmlErr = &xmlError***REMOVED***InnerXML: []byte(p.XMLError)***REMOVED***
		***REMOVED***
		resp.Propstat = append(resp.Propstat, propstat***REMOVED***
			Status:              fmt.Sprintf("HTTP/1.1 %d %s", p.Status, StatusText(p.Status)),
			Prop:                p.Props,
			ResponseDescription: p.ResponseDescription,
			Error:               xmlErr,
		***REMOVED***)
	***REMOVED***
	return &resp
***REMOVED***

const (
	infiniteDepth = -1
	invalidDepth  = -2
)

// parseDepth maps the strings "0", "1" and "infinity" to 0, 1 and
// infiniteDepth. Parsing any other string returns invalidDepth.
//
// Different WebDAV methods have further constraints on valid depths:
//	- PROPFIND has no further restrictions, as per section 9.1.
//	- COPY accepts only "0" or "infinity", as per section 9.8.3.
//	- MOVE accepts only "infinity", as per section 9.9.2.
//	- LOCK accepts only "0" or "infinity", as per section 9.10.3.
// These constraints are enforced by the handleXxx methods.
func parseDepth(s string) int ***REMOVED***
	switch s ***REMOVED***
	case "0":
		return 0
	case "1":
		return 1
	case "infinity":
		return infiniteDepth
	***REMOVED***
	return invalidDepth
***REMOVED***

// http://www.webdav.org/specs/rfc4918.html#status.code.extensions.to.http11
const (
	StatusMulti               = 207
	StatusUnprocessableEntity = 422
	StatusLocked              = 423
	StatusFailedDependency    = 424
	StatusInsufficientStorage = 507
)

func StatusText(code int) string ***REMOVED***
	switch code ***REMOVED***
	case StatusMulti:
		return "Multi-Status"
	case StatusUnprocessableEntity:
		return "Unprocessable Entity"
	case StatusLocked:
		return "Locked"
	case StatusFailedDependency:
		return "Failed Dependency"
	case StatusInsufficientStorage:
		return "Insufficient Storage"
	***REMOVED***
	return http.StatusText(code)
***REMOVED***

var (
	errDestinationEqualsSource = errors.New("webdav: destination equals source")
	errDirectoryNotEmpty       = errors.New("webdav: directory not empty")
	errInvalidDepth            = errors.New("webdav: invalid depth")
	errInvalidDestination      = errors.New("webdav: invalid destination")
	errInvalidIfHeader         = errors.New("webdav: invalid If header")
	errInvalidLockInfo         = errors.New("webdav: invalid lock info")
	errInvalidLockToken        = errors.New("webdav: invalid lock token")
	errInvalidPropfind         = errors.New("webdav: invalid propfind")
	errInvalidProppatch        = errors.New("webdav: invalid proppatch")
	errInvalidResponse         = errors.New("webdav: invalid response")
	errInvalidTimeout          = errors.New("webdav: invalid timeout")
	errNoFileSystem            = errors.New("webdav: no file system")
	errNoLockSystem            = errors.New("webdav: no lock system")
	errNotADirectory           = errors.New("webdav: not a directory")
	errPrefixMismatch          = errors.New("webdav: prefix mismatch")
	errRecursionTooDeep        = errors.New("webdav: recursion too deep")
	errUnsupportedLockInfo     = errors.New("webdav: unsupported lock info")
	errUnsupportedMethod       = errors.New("webdav: unsupported method")
)
