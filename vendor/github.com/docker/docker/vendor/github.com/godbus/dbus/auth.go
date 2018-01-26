package dbus

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
)

// AuthStatus represents the Status of an authentication mechanism.
type AuthStatus byte

const (
	// AuthOk signals that authentication is finished; the next command
	// from the server should be an OK.
	AuthOk AuthStatus = iota

	// AuthContinue signals that additional data is needed; the next command
	// from the server should be a DATA.
	AuthContinue

	// AuthError signals an error; the server sent invalid data or some
	// other unexpected thing happened and the current authentication
	// process should be aborted.
	AuthError
)

type authState byte

const (
	waitingForData authState = iota
	waitingForOk
	waitingForReject
)

// Auth defines the behaviour of an authentication mechanism.
type Auth interface ***REMOVED***
	// Return the name of the mechnism, the argument to the first AUTH command
	// and the next status.
	FirstData() (name, resp []byte, status AuthStatus)

	// Process the given DATA command, and return the argument to the DATA
	// command and the next status. If len(resp) == 0, no DATA command is sent.
	HandleData(data []byte) (resp []byte, status AuthStatus)
***REMOVED***

// Auth authenticates the connection, trying the given list of authentication
// mechanisms (in that order). If nil is passed, the EXTERNAL and
// DBUS_COOKIE_SHA1 mechanisms are tried for the current user. For private
// connections, this method must be called before sending any messages to the
// bus. Auth must not be called on shared connections.
func (conn *Conn) Auth(methods []Auth) error ***REMOVED***
	if methods == nil ***REMOVED***
		uid := strconv.Itoa(os.Getuid())
		methods = []Auth***REMOVED***AuthExternal(uid), AuthCookieSha1(uid, getHomeDir())***REMOVED***
	***REMOVED***
	in := bufio.NewReader(conn.transport)
	err := conn.transport.SendNullByte()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = authWriteLine(conn.transport, []byte("AUTH"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	s, err := authReadLine(in)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(s) < 2 || !bytes.Equal(s[0], []byte("REJECTED")) ***REMOVED***
		return errors.New("dbus: authentication protocol error")
	***REMOVED***
	s = s[1:]
	for _, v := range s ***REMOVED***
		for _, m := range methods ***REMOVED***
			if name, data, status := m.FirstData(); bytes.Equal(v, name) ***REMOVED***
				var ok bool
				err = authWriteLine(conn.transport, []byte("AUTH"), []byte(v), data)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				switch status ***REMOVED***
				case AuthOk:
					err, ok = conn.tryAuth(m, waitingForOk, in)
				case AuthContinue:
					err, ok = conn.tryAuth(m, waitingForData, in)
				default:
					panic("dbus: invalid authentication status")
				***REMOVED***
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				if ok ***REMOVED***
					if conn.transport.SupportsUnixFDs() ***REMOVED***
						err = authWriteLine(conn, []byte("NEGOTIATE_UNIX_FD"))
						if err != nil ***REMOVED***
							return err
						***REMOVED***
						line, err := authReadLine(in)
						if err != nil ***REMOVED***
							return err
						***REMOVED***
						switch ***REMOVED***
						case bytes.Equal(line[0], []byte("AGREE_UNIX_FD")):
							conn.EnableUnixFDs()
							conn.unixFD = true
						case bytes.Equal(line[0], []byte("ERROR")):
						default:
							return errors.New("dbus: authentication protocol error")
						***REMOVED***
					***REMOVED***
					err = authWriteLine(conn.transport, []byte("BEGIN"))
					if err != nil ***REMOVED***
						return err
					***REMOVED***
					go conn.inWorker()
					go conn.outWorker()
					return nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return errors.New("dbus: authentication failed")
***REMOVED***

// tryAuth tries to authenticate with m as the mechanism, using state as the
// initial authState and in for reading input. It returns (nil, true) on
// success, (nil, false) on a REJECTED and (someErr, false) if some other
// error occured.
func (conn *Conn) tryAuth(m Auth, state authState, in *bufio.Reader) (error, bool) ***REMOVED***
	for ***REMOVED***
		s, err := authReadLine(in)
		if err != nil ***REMOVED***
			return err, false
		***REMOVED***
		switch ***REMOVED***
		case state == waitingForData && string(s[0]) == "DATA":
			if len(s) != 2 ***REMOVED***
				err = authWriteLine(conn.transport, []byte("ERROR"))
				if err != nil ***REMOVED***
					return err, false
				***REMOVED***
				continue
			***REMOVED***
			data, status := m.HandleData(s[1])
			switch status ***REMOVED***
			case AuthOk, AuthContinue:
				if len(data) != 0 ***REMOVED***
					err = authWriteLine(conn.transport, []byte("DATA"), data)
					if err != nil ***REMOVED***
						return err, false
					***REMOVED***
				***REMOVED***
				if status == AuthOk ***REMOVED***
					state = waitingForOk
				***REMOVED***
			case AuthError:
				err = authWriteLine(conn.transport, []byte("ERROR"))
				if err != nil ***REMOVED***
					return err, false
				***REMOVED***
			***REMOVED***
		case state == waitingForData && string(s[0]) == "REJECTED":
			return nil, false
		case state == waitingForData && string(s[0]) == "ERROR":
			err = authWriteLine(conn.transport, []byte("CANCEL"))
			if err != nil ***REMOVED***
				return err, false
			***REMOVED***
			state = waitingForReject
		case state == waitingForData && string(s[0]) == "OK":
			if len(s) != 2 ***REMOVED***
				err = authWriteLine(conn.transport, []byte("CANCEL"))
				if err != nil ***REMOVED***
					return err, false
				***REMOVED***
				state = waitingForReject
			***REMOVED***
			conn.uuid = string(s[1])
			return nil, true
		case state == waitingForData:
			err = authWriteLine(conn.transport, []byte("ERROR"))
			if err != nil ***REMOVED***
				return err, false
			***REMOVED***
		case state == waitingForOk && string(s[0]) == "OK":
			if len(s) != 2 ***REMOVED***
				err = authWriteLine(conn.transport, []byte("CANCEL"))
				if err != nil ***REMOVED***
					return err, false
				***REMOVED***
				state = waitingForReject
			***REMOVED***
			conn.uuid = string(s[1])
			return nil, true
		case state == waitingForOk && string(s[0]) == "REJECTED":
			return nil, false
		case state == waitingForOk && (string(s[0]) == "DATA" ||
			string(s[0]) == "ERROR"):

			err = authWriteLine(conn.transport, []byte("CANCEL"))
			if err != nil ***REMOVED***
				return err, false
			***REMOVED***
			state = waitingForReject
		case state == waitingForOk:
			err = authWriteLine(conn.transport, []byte("ERROR"))
			if err != nil ***REMOVED***
				return err, false
			***REMOVED***
		case state == waitingForReject && string(s[0]) == "REJECTED":
			return nil, false
		case state == waitingForReject:
			return errors.New("dbus: authentication protocol error"), false
		default:
			panic("dbus: invalid auth state")
		***REMOVED***
	***REMOVED***
***REMOVED***

// authReadLine reads a line and separates it into its fields.
func authReadLine(in *bufio.Reader) ([][]byte, error) ***REMOVED***
	data, err := in.ReadBytes('\n')
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	data = bytes.TrimSuffix(data, []byte("\r\n"))
	return bytes.Split(data, []byte***REMOVED***' '***REMOVED***), nil
***REMOVED***

// authWriteLine writes the given line in the authentication protocol format
// (elements of data separated by a " " and terminated by "\r\n").
func authWriteLine(out io.Writer, data ...[]byte) error ***REMOVED***
	buf := make([]byte, 0)
	for i, v := range data ***REMOVED***
		buf = append(buf, v...)
		if i != len(data)-1 ***REMOVED***
			buf = append(buf, ' ')
		***REMOVED***
	***REMOVED***
	buf = append(buf, '\r')
	buf = append(buf, '\n')
	n, err := out.Write(buf)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n != len(buf) ***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED***
	return nil
***REMOVED***
