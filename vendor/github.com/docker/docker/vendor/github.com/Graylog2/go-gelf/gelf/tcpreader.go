package gelf

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type TCPReader struct ***REMOVED***
	listener *net.TCPListener
	conn     net.Conn
	messages chan []byte
***REMOVED***

type connChannels struct ***REMOVED***
	drop    chan string
	confirm chan string
***REMOVED***

func newTCPReader(addr string) (*TCPReader, chan string, chan string, error) ***REMOVED***
	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil ***REMOVED***
		return nil, nil, nil, fmt.Errorf("ResolveTCPAddr('%s'): %s", addr, err)
	***REMOVED***

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil ***REMOVED***
		return nil, nil, nil, fmt.Errorf("ListenTCP: %s", err)
	***REMOVED***

	r := &TCPReader***REMOVED***
		listener: listener,
		messages: make(chan []byte, 100), // Make a buffered channel with at most 100 messages
	***REMOVED***

	closeSignal := make(chan string, 1)
	doneSignal := make(chan string, 1)

	go r.listenUntilCloseSignal(closeSignal, doneSignal)

	return r, closeSignal, doneSignal, nil
***REMOVED***

func (r *TCPReader) accepter(connections chan net.Conn) ***REMOVED***
	for ***REMOVED***
		conn, err := r.listener.Accept()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		connections <- conn
	***REMOVED***
***REMOVED***

func (r *TCPReader) listenUntilCloseSignal(closeSignal chan string, doneSignal chan string) ***REMOVED***
	defer func() ***REMOVED*** doneSignal <- "done" ***REMOVED***()
	defer r.listener.Close()
	var conns []connChannels
	connectionsChannel := make(chan net.Conn, 1)
	go r.accepter(connectionsChannel)
	for ***REMOVED***
		select ***REMOVED***
		case conn := <-connectionsChannel:
			dropSignal := make(chan string, 1)
			dropConfirm := make(chan string, 1)
			channels := connChannels***REMOVED***drop: dropSignal, confirm: dropConfirm***REMOVED***
			go handleConnection(conn, r.messages, dropSignal, dropConfirm)
			conns = append(conns, channels)
		default:
		***REMOVED***

		select ***REMOVED***
		case sig := <-closeSignal:
			if sig == "stop" || sig == "drop" ***REMOVED***
				if len(conns) >= 1 ***REMOVED***
					for _, s := range conns ***REMOVED***
						if s.drop != nil ***REMOVED***
							s.drop <- "drop"
							<-s.confirm
							conns = append(conns[:0], conns[1:]...)
						***REMOVED***
					***REMOVED***
					if sig == "stop" ***REMOVED***
						return
					***REMOVED***
				***REMOVED*** else if sig == "stop" ***REMOVED***
					closeSignal <- "stop"
				***REMOVED***
				if sig == "drop" ***REMOVED***
					doneSignal <- "done"
				***REMOVED***
			***REMOVED***
		default:
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *TCPReader) addr() string ***REMOVED***
	return r.listener.Addr().String()
***REMOVED***

func handleConnection(conn net.Conn, messages chan<- []byte, dropSignal chan string, dropConfirm chan string) ***REMOVED***
	defer func() ***REMOVED*** dropConfirm <- "done" ***REMOVED***()
	defer conn.Close()
	reader := bufio.NewReader(conn)

	var b []byte
	var err error
	drop := false
	canDrop := false

	for ***REMOVED***
		conn.SetDeadline(time.Now().Add(2 * time.Second))
		if b, err = reader.ReadBytes(0); err != nil ***REMOVED***
			if drop ***REMOVED***
				return
			***REMOVED***
		***REMOVED*** else if len(b) > 0 ***REMOVED***
			messages <- b
			canDrop = true
			if drop ***REMOVED***
				return
			***REMOVED***
		***REMOVED*** else if drop ***REMOVED***
			return
		***REMOVED***
		select ***REMOVED***
		case sig := <-dropSignal:
			if sig == "drop" ***REMOVED***
				drop = true
				time.Sleep(1 * time.Second)
				if canDrop ***REMOVED***
					return
				***REMOVED***
			***REMOVED***
		default:
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *TCPReader) readMessage() (*Message, error) ***REMOVED***
	b := <-r.messages

	var msg Message
	if err := json.Unmarshal(b[:len(b)-1], &msg); err != nil ***REMOVED***
		return nil, fmt.Errorf("json.Unmarshal: %s", err)
	***REMOVED***

	return &msg, nil
***REMOVED***

func (r *TCPReader) Close() ***REMOVED***
	r.listener.Close()
***REMOVED***
