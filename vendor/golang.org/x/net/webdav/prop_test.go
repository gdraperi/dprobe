// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webdav

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sort"
	"testing"

	"golang.org/x/net/context"
)

func TestMemPS(t *testing.T) ***REMOVED***
	ctx := context.Background()
	// calcProps calculates the getlastmodified and getetag DAV: property
	// values in pstats for resource name in file-system fs.
	calcProps := func(name string, fs FileSystem, ls LockSystem, pstats []Propstat) error ***REMOVED***
		fi, err := fs.Stat(ctx, name)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, pst := range pstats ***REMOVED***
			for i, p := range pst.Props ***REMOVED***
				switch p.XMLName ***REMOVED***
				case xml.Name***REMOVED***Space: "DAV:", Local: "getlastmodified"***REMOVED***:
					p.InnerXML = []byte(fi.ModTime().Format(http.TimeFormat))
					pst.Props[i] = p
				case xml.Name***REMOVED***Space: "DAV:", Local: "getetag"***REMOVED***:
					if fi.IsDir() ***REMOVED***
						continue
					***REMOVED***
					etag, err := findETag(ctx, fs, ls, name, fi)
					if err != nil ***REMOVED***
						return err
					***REMOVED***
					p.InnerXML = []byte(etag)
					pst.Props[i] = p
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	const (
		lockEntry = `` +
			`<D:lockentry xmlns:D="DAV:">` +
			`<D:lockscope><D:exclusive/></D:lockscope>` +
			`<D:locktype><D:write/></D:locktype>` +
			`</D:lockentry>`
		statForbiddenError = `<D:cannot-modify-protected-property xmlns:D="DAV:"/>`
	)

	type propOp struct ***REMOVED***
		op            string
		name          string
		pnames        []xml.Name
		patches       []Proppatch
		wantPnames    []xml.Name
		wantPropstats []Propstat
	***REMOVED***

	testCases := []struct ***REMOVED***
		desc        string
		noDeadProps bool
		buildfs     []string
		propOp      []propOp
	***REMOVED******REMOVED******REMOVED***
		desc:    "propname",
		buildfs: []string***REMOVED***"mkdir /dir", "touch /file"***REMOVED***,
		propOp: []propOp***REMOVED******REMOVED***
			op:   "propname",
			name: "/dir",
			wantPnames: []xml.Name***REMOVED***
				***REMOVED***Space: "DAV:", Local: "resourcetype"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "displayname"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "supportedlock"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "getlastmodified"***REMOVED***,
			***REMOVED***,
		***REMOVED***, ***REMOVED***
			op:   "propname",
			name: "/file",
			wantPnames: []xml.Name***REMOVED***
				***REMOVED***Space: "DAV:", Local: "resourcetype"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "displayname"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "getcontentlength"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "getlastmodified"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "getcontenttype"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "getetag"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "supportedlock"***REMOVED***,
			***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:    "allprop dir and file",
		buildfs: []string***REMOVED***"mkdir /dir", "write /file foobarbaz"***REMOVED***,
		propOp: []propOp***REMOVED******REMOVED***
			op:   "allprop",
			name: "/dir",
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "resourcetype"***REMOVED***,
					InnerXML: []byte(`<D:collection xmlns:D="DAV:"/>`),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "displayname"***REMOVED***,
					InnerXML: []byte("dir"),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "getlastmodified"***REMOVED***,
					InnerXML: nil, // Calculated during test.
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "supportedlock"***REMOVED***,
					InnerXML: []byte(lockEntry),
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:   "allprop",
			name: "/file",
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "resourcetype"***REMOVED***,
					InnerXML: []byte(""),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "displayname"***REMOVED***,
					InnerXML: []byte("file"),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "getcontentlength"***REMOVED***,
					InnerXML: []byte("9"),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "getlastmodified"***REMOVED***,
					InnerXML: nil, // Calculated during test.
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "getcontenttype"***REMOVED***,
					InnerXML: []byte("text/plain; charset=utf-8"),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "getetag"***REMOVED***,
					InnerXML: nil, // Calculated during test.
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "supportedlock"***REMOVED***,
					InnerXML: []byte(lockEntry),
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:   "allprop",
			name: "/file",
			pnames: []xml.Name***REMOVED***
				***REMOVED***"DAV:", "resourcetype"***REMOVED***,
				***REMOVED***"foo", "bar"***REMOVED***,
			***REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "resourcetype"***REMOVED***,
					InnerXML: []byte(""),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "displayname"***REMOVED***,
					InnerXML: []byte("file"),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "getcontentlength"***REMOVED***,
					InnerXML: []byte("9"),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "getlastmodified"***REMOVED***,
					InnerXML: nil, // Calculated during test.
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "getcontenttype"***REMOVED***,
					InnerXML: []byte("text/plain; charset=utf-8"),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "getetag"***REMOVED***,
					InnerXML: nil, // Calculated during test.
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "supportedlock"***REMOVED***,
					InnerXML: []byte(lockEntry),
				***REMOVED******REMOVED******REMOVED***, ***REMOVED***
				Status: http.StatusNotFound,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:    "propfind DAV:resourcetype",
		buildfs: []string***REMOVED***"mkdir /dir", "touch /file"***REMOVED***,
		propOp: []propOp***REMOVED******REMOVED***
			op:     "propfind",
			name:   "/dir",
			pnames: []xml.Name***REMOVED******REMOVED***"DAV:", "resourcetype"***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "resourcetype"***REMOVED***,
					InnerXML: []byte(`<D:collection xmlns:D="DAV:"/>`),
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:     "propfind",
			name:   "/file",
			pnames: []xml.Name***REMOVED******REMOVED***"DAV:", "resourcetype"***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "resourcetype"***REMOVED***,
					InnerXML: []byte(""),
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:    "propfind unsupported DAV properties",
		buildfs: []string***REMOVED***"mkdir /dir"***REMOVED***,
		propOp: []propOp***REMOVED******REMOVED***
			op:     "propfind",
			name:   "/dir",
			pnames: []xml.Name***REMOVED******REMOVED***"DAV:", "getcontentlanguage"***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusNotFound,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "DAV:", Local: "getcontentlanguage"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:     "propfind",
			name:   "/dir",
			pnames: []xml.Name***REMOVED******REMOVED***"DAV:", "creationdate"***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusNotFound,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "DAV:", Local: "creationdate"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:    "propfind getetag for files but not for directories",
		buildfs: []string***REMOVED***"mkdir /dir", "touch /file"***REMOVED***,
		propOp: []propOp***REMOVED******REMOVED***
			op:     "propfind",
			name:   "/dir",
			pnames: []xml.Name***REMOVED******REMOVED***"DAV:", "getetag"***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusNotFound,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "DAV:", Local: "getetag"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:     "propfind",
			name:   "/file",
			pnames: []xml.Name***REMOVED******REMOVED***"DAV:", "getetag"***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "getetag"***REMOVED***,
					InnerXML: nil, // Calculated during test.
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:        "proppatch property on no-dead-properties file system",
		buildfs:     []string***REMOVED***"mkdir /dir"***REMOVED***,
		noDeadProps: true,
		propOp: []propOp***REMOVED******REMOVED***
			op:   "proppatch",
			name: "/dir",
			patches: []Proppatch***REMOVED******REMOVED***
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusForbidden,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:   "proppatch",
			name: "/dir",
			patches: []Proppatch***REMOVED******REMOVED***
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "DAV:", Local: "getetag"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status:   http.StatusForbidden,
				XMLError: statForbiddenError,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "DAV:", Local: "getetag"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:    "proppatch dead property",
		buildfs: []string***REMOVED***"mkdir /dir"***REMOVED***,
		propOp: []propOp***REMOVED******REMOVED***
			op:   "proppatch",
			name: "/dir",
			patches: []Proppatch***REMOVED******REMOVED***
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
					InnerXML: []byte("baz"),
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:     "propfind",
			name:   "/dir",
			pnames: []xml.Name***REMOVED******REMOVED***Space: "foo", Local: "bar"***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
					InnerXML: []byte("baz"),
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:    "proppatch dead property with failed dependency",
		buildfs: []string***REMOVED***"mkdir /dir"***REMOVED***,
		propOp: []propOp***REMOVED******REMOVED***
			op:   "proppatch",
			name: "/dir",
			patches: []Proppatch***REMOVED******REMOVED***
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
					InnerXML: []byte("baz"),
				***REMOVED******REMOVED***,
			***REMOVED***, ***REMOVED***
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "DAV:", Local: "displayname"***REMOVED***,
					InnerXML: []byte("xxx"),
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status:   http.StatusForbidden,
				XMLError: statForbiddenError,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "DAV:", Local: "displayname"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED***, ***REMOVED***
				Status: StatusFailedDependency,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:     "propfind",
			name:   "/dir",
			pnames: []xml.Name***REMOVED******REMOVED***Space: "foo", Local: "bar"***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusNotFound,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:    "proppatch remove dead property",
		buildfs: []string***REMOVED***"mkdir /dir"***REMOVED***,
		propOp: []propOp***REMOVED******REMOVED***
			op:   "proppatch",
			name: "/dir",
			patches: []Proppatch***REMOVED******REMOVED***
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
					InnerXML: []byte("baz"),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "spam", Local: "ham"***REMOVED***,
					InnerXML: []byte("eggs"),
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED***, ***REMOVED***
					XMLName: xml.Name***REMOVED***Space: "spam", Local: "ham"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:   "propfind",
			name: "/dir",
			pnames: []xml.Name***REMOVED***
				***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED***Space: "spam", Local: "ham"***REMOVED***,
			***REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
					InnerXML: []byte("baz"),
				***REMOVED***, ***REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "spam", Local: "ham"***REMOVED***,
					InnerXML: []byte("eggs"),
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:   "proppatch",
			name: "/dir",
			patches: []Proppatch***REMOVED******REMOVED***
				Remove: true,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:   "propfind",
			name: "/dir",
			pnames: []xml.Name***REMOVED***
				***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED***Space: "spam", Local: "ham"***REMOVED***,
			***REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusNotFound,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED***, ***REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "spam", Local: "ham"***REMOVED***,
					InnerXML: []byte("eggs"),
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:    "propname with dead property",
		buildfs: []string***REMOVED***"touch /file"***REMOVED***,
		propOp: []propOp***REMOVED******REMOVED***
			op:   "proppatch",
			name: "/file",
			patches: []Proppatch***REMOVED******REMOVED***
				Props: []Property***REMOVED******REMOVED***
					XMLName:  xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
					InnerXML: []byte("baz"),
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED***, ***REMOVED***
			op:   "propname",
			name: "/file",
			wantPnames: []xml.Name***REMOVED***
				***REMOVED***Space: "DAV:", Local: "resourcetype"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "displayname"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "getcontentlength"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "getlastmodified"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "getcontenttype"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "getetag"***REMOVED***,
				***REMOVED***Space: "DAV:", Local: "supportedlock"***REMOVED***,
				***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
			***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:    "proppatch remove unknown dead property",
		buildfs: []string***REMOVED***"mkdir /dir"***REMOVED***,
		propOp: []propOp***REMOVED******REMOVED***
			op:   "proppatch",
			name: "/dir",
			patches: []Proppatch***REMOVED******REMOVED***
				Remove: true,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusOK,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:    "bad: propfind unknown property",
		buildfs: []string***REMOVED***"mkdir /dir"***REMOVED***,
		propOp: []propOp***REMOVED******REMOVED***
			op:     "propfind",
			name:   "/dir",
			pnames: []xml.Name***REMOVED******REMOVED***"foo:", "bar"***REMOVED******REMOVED***,
			wantPropstats: []Propstat***REMOVED******REMOVED***
				Status: http.StatusNotFound,
				Props: []Property***REMOVED******REMOVED***
					XMLName: xml.Name***REMOVED***Space: "foo:", Local: "bar"***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED******REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED******REMOVED***

	for _, tc := range testCases ***REMOVED***
		fs, err := buildTestFS(tc.buildfs)
		if err != nil ***REMOVED***
			t.Fatalf("%s: cannot create test filesystem: %v", tc.desc, err)
		***REMOVED***
		if tc.noDeadProps ***REMOVED***
			fs = noDeadPropsFS***REMOVED***fs***REMOVED***
		***REMOVED***
		ls := NewMemLS()
		for _, op := range tc.propOp ***REMOVED***
			desc := fmt.Sprintf("%s: %s %s", tc.desc, op.op, op.name)
			if err = calcProps(op.name, fs, ls, op.wantPropstats); err != nil ***REMOVED***
				t.Fatalf("%s: calcProps: %v", desc, err)
			***REMOVED***

			// Call property system.
			var propstats []Propstat
			switch op.op ***REMOVED***
			case "propname":
				pnames, err := propnames(ctx, fs, ls, op.name)
				if err != nil ***REMOVED***
					t.Errorf("%s: got error %v, want nil", desc, err)
					continue
				***REMOVED***
				sort.Sort(byXMLName(pnames))
				sort.Sort(byXMLName(op.wantPnames))
				if !reflect.DeepEqual(pnames, op.wantPnames) ***REMOVED***
					t.Errorf("%s: pnames\ngot  %q\nwant %q", desc, pnames, op.wantPnames)
				***REMOVED***
				continue
			case "allprop":
				propstats, err = allprop(ctx, fs, ls, op.name, op.pnames)
			case "propfind":
				propstats, err = props(ctx, fs, ls, op.name, op.pnames)
			case "proppatch":
				propstats, err = patch(ctx, fs, ls, op.name, op.patches)
			default:
				t.Fatalf("%s: %s not implemented", desc, op.op)
			***REMOVED***
			if err != nil ***REMOVED***
				t.Errorf("%s: got error %v, want nil", desc, err)
				continue
			***REMOVED***
			// Compare return values from allprop, propfind or proppatch.
			for _, pst := range propstats ***REMOVED***
				sort.Sort(byPropname(pst.Props))
			***REMOVED***
			for _, pst := range op.wantPropstats ***REMOVED***
				sort.Sort(byPropname(pst.Props))
			***REMOVED***
			sort.Sort(byStatus(propstats))
			sort.Sort(byStatus(op.wantPropstats))
			if !reflect.DeepEqual(propstats, op.wantPropstats) ***REMOVED***
				t.Errorf("%s: propstat\ngot  %q\nwant %q", desc, propstats, op.wantPropstats)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func cmpXMLName(a, b xml.Name) bool ***REMOVED***
	if a.Space != b.Space ***REMOVED***
		return a.Space < b.Space
	***REMOVED***
	return a.Local < b.Local
***REMOVED***

type byXMLName []xml.Name

func (b byXMLName) Len() int           ***REMOVED*** return len(b) ***REMOVED***
func (b byXMLName) Swap(i, j int)      ***REMOVED*** b[i], b[j] = b[j], b[i] ***REMOVED***
func (b byXMLName) Less(i, j int) bool ***REMOVED*** return cmpXMLName(b[i], b[j]) ***REMOVED***

type byPropname []Property

func (b byPropname) Len() int           ***REMOVED*** return len(b) ***REMOVED***
func (b byPropname) Swap(i, j int)      ***REMOVED*** b[i], b[j] = b[j], b[i] ***REMOVED***
func (b byPropname) Less(i, j int) bool ***REMOVED*** return cmpXMLName(b[i].XMLName, b[j].XMLName) ***REMOVED***

type byStatus []Propstat

func (b byStatus) Len() int           ***REMOVED*** return len(b) ***REMOVED***
func (b byStatus) Swap(i, j int)      ***REMOVED*** b[i], b[j] = b[j], b[i] ***REMOVED***
func (b byStatus) Less(i, j int) bool ***REMOVED*** return b[i].Status < b[j].Status ***REMOVED***

type noDeadPropsFS struct ***REMOVED***
	FileSystem
***REMOVED***

func (fs noDeadPropsFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (File, error) ***REMOVED***
	f, err := fs.FileSystem.OpenFile(ctx, name, flag, perm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return noDeadPropsFile***REMOVED***f***REMOVED***, nil
***REMOVED***

// noDeadPropsFile wraps a File but strips any optional DeadPropsHolder methods
// provided by the underlying File implementation.
type noDeadPropsFile struct ***REMOVED***
	f File
***REMOVED***

func (f noDeadPropsFile) Close() error                              ***REMOVED*** return f.f.Close() ***REMOVED***
func (f noDeadPropsFile) Read(p []byte) (int, error)                ***REMOVED*** return f.f.Read(p) ***REMOVED***
func (f noDeadPropsFile) Readdir(count int) ([]os.FileInfo, error)  ***REMOVED*** return f.f.Readdir(count) ***REMOVED***
func (f noDeadPropsFile) Seek(off int64, whence int) (int64, error) ***REMOVED*** return f.f.Seek(off, whence) ***REMOVED***
func (f noDeadPropsFile) Stat() (os.FileInfo, error)                ***REMOVED*** return f.f.Stat() ***REMOVED***
func (f noDeadPropsFile) Write(p []byte) (int, error)               ***REMOVED*** return f.f.Write(p) ***REMOVED***
