package dbus

func (t *unixTransport) SendNullByte() error ***REMOVED***
	_, err := t.Write([]byte***REMOVED***0***REMOVED***)
	return err
***REMOVED***
