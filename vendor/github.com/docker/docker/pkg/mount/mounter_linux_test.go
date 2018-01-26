// +build linux

package mount

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	selinux "github.com/opencontainers/selinux/go-selinux"
)

func TestMount(t *testing.T) ***REMOVED***
	if os.Getuid() != 0 ***REMOVED***
		t.Skip("not root tests would fail")
	***REMOVED***

	source, err := ioutil.TempDir("", "mount-test-source-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(source)

	// Ensure we have a known start point by mounting tmpfs with given options
	if err := Mount("tmpfs", source, "tmpfs", "private"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer ensureUnmount(t, source)
	validateMount(t, source, "", "", "")
	if t.Failed() ***REMOVED***
		t.FailNow()
	***REMOVED***

	target, err := ioutil.TempDir("", "mount-test-target-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(target)

	tests := []struct ***REMOVED***
		source           string
		ftype            string
		options          string
		expectedOpts     string
		expectedOptional string
		expectedVFS      string
	***REMOVED******REMOVED***
		// No options
		***REMOVED***"tmpfs", "tmpfs", "", "", "", ""***REMOVED***,
		// Default rw / ro test
		***REMOVED***source, "", "bind", "", "", ""***REMOVED***,
		***REMOVED***source, "", "bind,private", "", "", ""***REMOVED***,
		***REMOVED***source, "", "bind,shared", "", "shared", ""***REMOVED***,
		***REMOVED***source, "", "bind,slave", "", "master", ""***REMOVED***,
		***REMOVED***source, "", "bind,unbindable", "", "unbindable", ""***REMOVED***,
		// Read Write tests
		***REMOVED***source, "", "bind,rw", "rw", "", ""***REMOVED***,
		***REMOVED***source, "", "bind,rw,private", "rw", "", ""***REMOVED***,
		***REMOVED***source, "", "bind,rw,shared", "rw", "shared", ""***REMOVED***,
		***REMOVED***source, "", "bind,rw,slave", "rw", "master", ""***REMOVED***,
		***REMOVED***source, "", "bind,rw,unbindable", "rw", "unbindable", ""***REMOVED***,
		// Read Only tests
		***REMOVED***source, "", "bind,ro", "ro", "", ""***REMOVED***,
		***REMOVED***source, "", "bind,ro,private", "ro", "", ""***REMOVED***,
		***REMOVED***source, "", "bind,ro,shared", "ro", "shared", ""***REMOVED***,
		***REMOVED***source, "", "bind,ro,slave", "ro", "master", ""***REMOVED***,
		***REMOVED***source, "", "bind,ro,unbindable", "ro", "unbindable", ""***REMOVED***,
		// Remount tests to change per filesystem options
		***REMOVED***"", "", "remount,size=128k", "rw", "", "rw,size=128k"***REMOVED***,
		***REMOVED***"", "", "remount,ro,size=128k", "ro", "", "ro,size=128k"***REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		ftype, options := tc.ftype, tc.options
		if tc.ftype == "" ***REMOVED***
			ftype = "none"
		***REMOVED***
		if tc.options == "" ***REMOVED***
			options = "none"
		***REMOVED***

		t.Run(fmt.Sprintf("%v-%v", ftype, options), func(t *testing.T) ***REMOVED***
			if strings.Contains(tc.options, "slave") ***REMOVED***
				// Slave requires a shared source
				if err := MakeShared(source); err != nil ***REMOVED***
					t.Fatal(err)
				***REMOVED***
				defer func() ***REMOVED***
					if err := MakePrivate(source); err != nil ***REMOVED***
						t.Fatal(err)
					***REMOVED***
				***REMOVED***()
			***REMOVED***
			if strings.Contains(tc.options, "remount") ***REMOVED***
				// create a new mount to remount first
				if err := Mount("tmpfs", target, "tmpfs", ""); err != nil ***REMOVED***
					t.Fatal(err)
				***REMOVED***
			***REMOVED***
			if err := Mount(tc.source, target, tc.ftype, tc.options); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			defer ensureUnmount(t, target)
			expectedVFS := tc.expectedVFS
			if selinux.GetEnabled() && expectedVFS != "" ***REMOVED***
				expectedVFS = expectedVFS + ",seclabel"
			***REMOVED***
			validateMount(t, target, tc.expectedOpts, tc.expectedOptional, expectedVFS)
		***REMOVED***)
	***REMOVED***
***REMOVED***

// ensureUnmount umounts mnt checking for errors
func ensureUnmount(t *testing.T, mnt string) ***REMOVED***
	if err := Unmount(mnt); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

// validateMount checks that mnt has the given options
func validateMount(t *testing.T, mnt string, opts, optional, vfs string) ***REMOVED***
	info, err := GetMounts()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	wantedOpts := make(map[string]struct***REMOVED******REMOVED***)
	if opts != "" ***REMOVED***
		for _, opt := range strings.Split(opts, ",") ***REMOVED***
			wantedOpts[opt] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	wantedOptional := make(map[string]struct***REMOVED******REMOVED***)
	if optional != "" ***REMOVED***
		for _, opt := range strings.Split(optional, ",") ***REMOVED***
			wantedOptional[opt] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	wantedVFS := make(map[string]struct***REMOVED******REMOVED***)
	if vfs != "" ***REMOVED***
		for _, opt := range strings.Split(vfs, ",") ***REMOVED***
			wantedVFS[opt] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	mnts := make(map[int]*Info, len(info))
	for _, mi := range info ***REMOVED***
		mnts[mi.ID] = mi
	***REMOVED***

	for _, mi := range info ***REMOVED***
		if mi.Mountpoint != mnt ***REMOVED***
			continue
		***REMOVED***

		// Use parent info as the defaults
		p := mnts[mi.Parent]
		pOpts := make(map[string]struct***REMOVED******REMOVED***)
		if p.Opts != "" ***REMOVED***
			for _, opt := range strings.Split(p.Opts, ",") ***REMOVED***
				pOpts[clean(opt)] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
		pOptional := make(map[string]struct***REMOVED******REMOVED***)
		if p.Optional != "" ***REMOVED***
			for _, field := range strings.Split(p.Optional, ",") ***REMOVED***
				pOptional[clean(field)] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***

		// Validate Opts
		if mi.Opts != "" ***REMOVED***
			for _, opt := range strings.Split(mi.Opts, ",") ***REMOVED***
				opt = clean(opt)
				if !has(wantedOpts, opt) && !has(pOpts, opt) ***REMOVED***
					t.Errorf("unexpected mount option %q expected %q", opt, opts)
				***REMOVED***
				delete(wantedOpts, opt)
			***REMOVED***
		***REMOVED***
		for opt := range wantedOpts ***REMOVED***
			t.Errorf("missing mount option %q found %q", opt, mi.Opts)
		***REMOVED***

		// Validate Optional
		if mi.Optional != "" ***REMOVED***
			for _, field := range strings.Split(mi.Optional, ",") ***REMOVED***
				field = clean(field)
				if !has(wantedOptional, field) && !has(pOptional, field) ***REMOVED***
					t.Errorf("unexpected optional failed %q expected %q", field, optional)
				***REMOVED***
				delete(wantedOptional, field)
			***REMOVED***
		***REMOVED***
		for field := range wantedOptional ***REMOVED***
			t.Errorf("missing optional field %q found %q", field, mi.Optional)
		***REMOVED***

		// Validate VFS if set
		if vfs != "" ***REMOVED***
			if mi.VfsOpts != "" ***REMOVED***
				for _, opt := range strings.Split(mi.VfsOpts, ",") ***REMOVED***
					opt = clean(opt)
					if !has(wantedVFS, opt) ***REMOVED***
						t.Errorf("unexpected mount option %q expected %q", opt, vfs)
					***REMOVED***
					delete(wantedVFS, opt)
				***REMOVED***
			***REMOVED***
			for opt := range wantedVFS ***REMOVED***
				t.Errorf("missing mount option %q found %q", opt, mi.VfsOpts)
			***REMOVED***
		***REMOVED***

		return
	***REMOVED***

	t.Errorf("failed to find mount %q", mnt)
***REMOVED***

// clean strips off any value param after the colon
func clean(v string) string ***REMOVED***
	return strings.SplitN(v, ":", 2)[0]
***REMOVED***

// has returns true if key is a member of m
func has(m map[string]struct***REMOVED******REMOVED***, key string) bool ***REMOVED***
	_, ok := m[key]
	return ok
***REMOVED***
