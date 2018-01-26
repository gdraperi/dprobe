// Package le_go provides a Golang client library for logging to
// logentries.com over a TCP connection.
//
// it uses an access token for sending log events.
package le_go

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Logger represents a Logentries logger,
// it holds the open TCP connection, access token, prefix and flags.
//
// all Logger operations are thread safe and blocking,
// log operations can be invoked in a non-blocking way by calling them from
// a goroutine.
type Logger struct ***REMOVED***
	conn   net.Conn
	flag   int
	mu     sync.Mutex
	prefix string
	token  string
	buf    []byte
***REMOVED***

const lineSep = "\n"

// Connect creates a new Logger instance and opens a TCP connection to
// logentries.com,
// The token can be generated at logentries.com by adding a new log,
// choosing manual configuration and token based TCP connection.
func Connect(token string) (*Logger, error) ***REMOVED***
	logger := Logger***REMOVED***
		token: token,
	***REMOVED***

	if err := logger.openConnection(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &logger, nil
***REMOVED***

// Close closes the TCP connection to logentries.com
func (logger *Logger) Close() error ***REMOVED***
	if logger.conn != nil ***REMOVED***
		return logger.conn.Close()
	***REMOVED***

	return nil
***REMOVED***

// Opens a TCP connection to logentries.com
func (logger *Logger) openConnection() error ***REMOVED***
	conn, err := tls.Dial("tcp", "data.logentries.com:443", &tls.Config***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	logger.conn = conn
	return nil
***REMOVED***

// It returns if the TCP connection to logentries.com is open
func (logger *Logger) isOpenConnection() bool ***REMOVED***
	if logger.conn == nil ***REMOVED***
		return false
	***REMOVED***

	buf := make([]byte, 1)

	logger.conn.SetReadDeadline(time.Now())

	_, err := logger.conn.Read(buf)

	switch err.(type) ***REMOVED***
	case net.Error:
		if err.(net.Error).Timeout() == true ***REMOVED***
			logger.conn.SetReadDeadline(time.Time***REMOVED******REMOVED***)

			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// It ensures that the TCP connection to logentries.com is open.
// If the connection is closed, a new one is opened.
func (logger *Logger) ensureOpenConnection() error ***REMOVED***
	if !logger.isOpenConnection() ***REMOVED***
		if err := logger.openConnection(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Fatal is same as Print() but calls to os.Exit(1)
func (logger *Logger) Fatal(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
***REMOVED***

// Fatalf is same as Printf() but calls to os.Exit(1)
func (logger *Logger) Fatalf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
***REMOVED***

// Fatalln is same as Println() but calls to os.Exit(1)
func (logger *Logger) Fatalln(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
***REMOVED***

// Flags returns the logger flags
func (logger *Logger) Flags() int ***REMOVED***
	return logger.flag
***REMOVED***

// Output does the actual writing to the TCP connection
func (logger *Logger) Output(calldepth int, s string) error ***REMOVED***
	_, err := logger.Write([]byte(s))

	return err
***REMOVED***

// Panic is same as Print() but calls to panic
func (logger *Logger) Panic(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	s := fmt.Sprint(v...)
	logger.Output(2, s)
	panic(s)
***REMOVED***

// Panicf is same as Printf() but calls to panic
func (logger *Logger) Panicf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	s := fmt.Sprintf(format, v...)
	logger.Output(2, s)
	panic(s)
***REMOVED***

// Panicln is same as Println() but calls to panic
func (logger *Logger) Panicln(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	s := fmt.Sprintln(v...)
	logger.Output(2, s)
	panic(s)
***REMOVED***

// Prefix returns the logger prefix
func (logger *Logger) Prefix() string ***REMOVED***
	return logger.prefix
***REMOVED***

// Print logs a message
func (logger *Logger) Print(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Output(2, fmt.Sprint(v...))
***REMOVED***

// Printf logs a formatted message
func (logger *Logger) Printf(format string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Output(2, fmt.Sprintf(format, v...))
***REMOVED***

// Println logs a message with a linebreak
func (logger *Logger) Println(v ...interface***REMOVED******REMOVED***) ***REMOVED***
	logger.Output(2, fmt.Sprintln(v...))
***REMOVED***

// SetFlags sets the logger flags
func (logger *Logger) SetFlags(flag int) ***REMOVED***
	logger.flag = flag
***REMOVED***

// SetPrefix sets the logger prefix
func (logger *Logger) SetPrefix(prefix string) ***REMOVED***
	logger.prefix = prefix
***REMOVED***

// Write writes a bytes array to the Logentries TCP connection,
// it adds the access token and prefix and also replaces
// line breaks with the unicode \u2028 character
func (logger *Logger) Write(p []byte) (n int, err error) ***REMOVED***
	if err := logger.ensureOpenConnection(); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	logger.mu.Lock()
	defer logger.mu.Unlock()

	logger.makeBuf(p)

	return logger.conn.Write(logger.buf)
***REMOVED***

// makeBuf constructs the logger buffer
// it is not safe to be used from within multiple concurrent goroutines
func (logger *Logger) makeBuf(p []byte) ***REMOVED***
	count := strings.Count(string(p), lineSep)
	p = []byte(strings.Replace(string(p), lineSep, "\u2028", count-1))

	logger.buf = logger.buf[:0]
	logger.buf = append(logger.buf, (logger.token + " ")...)
	logger.buf = append(logger.buf, (logger.prefix + " ")...)
	logger.buf = append(logger.buf, p...)

	if !strings.HasSuffix(string(logger.buf), lineSep) ***REMOVED***
		logger.buf = append(logger.buf, (lineSep)...)
	***REMOVED***
***REMOVED***
