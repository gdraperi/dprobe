package dbus

import (
	"encoding/binary"
	"errors"
	"io"
)

type genericTransport struct ***REMOVED***
	io.ReadWriteCloser
***REMOVED***

func (t genericTransport) SendNullByte() error ***REMOVED***
	_, err := t.Write([]byte***REMOVED***0***REMOVED***)
	return err
***REMOVED***

func (t genericTransport) SupportsUnixFDs() bool ***REMOVED***
	return false
***REMOVED***

func (t genericTransport) EnableUnixFDs() ***REMOVED******REMOVED***

func (t genericTransport) ReadMessage() (*Message, error) ***REMOVED***
	return DecodeMessage(t)
***REMOVED***

func (t genericTransport) SendMessage(msg *Message) error ***REMOVED***
	for _, v := range msg.Body ***REMOVED***
		if _, ok := v.(UnixFD); ok ***REMOVED***
			return errors.New("dbus: unix fd passing not enabled")
		***REMOVED***
	***REMOVED***
	return msg.EncodeTo(t, binary.LittleEndian)
***REMOVED***
