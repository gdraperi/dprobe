package sysx

import (
	"bytes"
	"fmt"
	"syscall"
)

const defaultXattrBufferSize = 5

var ErrNotSupported = fmt.Errorf("not supported")

type listxattrFunc func(path string, dest []byte) (int, error)

func listxattrAll(path string, listFunc listxattrFunc) ([]string, error) ***REMOVED***
	var p []byte // nil on first execution

	for ***REMOVED***
		n, err := listFunc(path, p) // first call gets buffer size.
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if n > len(p) ***REMOVED***
			p = make([]byte, n)
			continue
		***REMOVED***

		p = p[:n]

		ps := bytes.Split(bytes.TrimSuffix(p, []byte***REMOVED***0***REMOVED***), []byte***REMOVED***0***REMOVED***)
		var entries []string
		for _, p := range ps ***REMOVED***
			s := string(p)
			if s != "" ***REMOVED***
				entries = append(entries, s)
			***REMOVED***
		***REMOVED***

		return entries, nil
	***REMOVED***
***REMOVED***

type getxattrFunc func(string, string, []byte) (int, error)

func getxattrAll(path, attr string, getFunc getxattrFunc) ([]byte, error) ***REMOVED***
	p := make([]byte, defaultXattrBufferSize)
	for ***REMOVED***
		n, err := getFunc(path, attr, p)
		if err != nil ***REMOVED***
			if errno, ok := err.(syscall.Errno); ok && errno == syscall.ERANGE ***REMOVED***
				p = make([]byte, len(p)*2) // this can't be ideal.
				continue                   // try again!
			***REMOVED***

			return nil, err
		***REMOVED***

		// realloc to correct size and repeat
		if n > len(p) ***REMOVED***
			p = make([]byte, n)
			continue
		***REMOVED***

		return p[:n], nil
	***REMOVED***
***REMOVED***
