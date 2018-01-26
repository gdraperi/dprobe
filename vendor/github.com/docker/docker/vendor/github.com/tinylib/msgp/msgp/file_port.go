// +build windows appengine

package msgp

import (
	"io/ioutil"
	"os"
)

// MarshalSizer is the combination
// of the Marshaler and Sizer
// interfaces.
type MarshalSizer interface ***REMOVED***
	Marshaler
	Sizer
***REMOVED***

func ReadFile(dst Unmarshaler, file *os.File) error ***REMOVED***
	if u, ok := dst.(Decodable); ok ***REMOVED***
		return u.DecodeMsg(NewReader(file))
	***REMOVED***

	data, err := ioutil.ReadAll(file)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = dst.UnmarshalMsg(data)
	return err
***REMOVED***

func WriteFile(src MarshalSizer, file *os.File) error ***REMOVED***
	if e, ok := src.(Encodable); ok ***REMOVED***
		w := NewWriter(file)
		err := e.EncodeMsg(w)
		if err == nil ***REMOVED***
			err = w.Flush()
		***REMOVED***
		return err
	***REMOVED***

	raw, err := src.MarshalMsg(nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = file.Write(raw)
	return err
***REMOVED***
