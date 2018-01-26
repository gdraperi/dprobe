package zk

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

var (
	requests     = make(map[int32]int32) // Map of Xid -> Opcode
	requestsLock = &sync.Mutex***REMOVED******REMOVED***
)

func trace(conn1, conn2 net.Conn, client bool) ***REMOVED***
	defer conn1.Close()
	defer conn2.Close()
	buf := make([]byte, 10*1024)
	init := true
	for ***REMOVED***
		_, err := io.ReadFull(conn1, buf[:4])
		if err != nil ***REMOVED***
			fmt.Println("1>", client, err)
			return
		***REMOVED***

		blen := int(binary.BigEndian.Uint32(buf[:4]))

		_, err = io.ReadFull(conn1, buf[4:4+blen])
		if err != nil ***REMOVED***
			fmt.Println("2>", client, err)
			return
		***REMOVED***

		var cr interface***REMOVED******REMOVED***
		opcode := int32(-1)
		readHeader := true
		if client ***REMOVED***
			if init ***REMOVED***
				cr = &connectRequest***REMOVED******REMOVED***
				readHeader = false
			***REMOVED*** else ***REMOVED***
				xid := int32(binary.BigEndian.Uint32(buf[4:8]))
				opcode = int32(binary.BigEndian.Uint32(buf[8:12]))
				requestsLock.Lock()
				requests[xid] = opcode
				requestsLock.Unlock()
				cr = requestStructForOp(opcode)
				if cr == nil ***REMOVED***
					fmt.Printf("Unknown opcode %d\n", opcode)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if init ***REMOVED***
				cr = &connectResponse***REMOVED******REMOVED***
				readHeader = false
			***REMOVED*** else ***REMOVED***
				xid := int32(binary.BigEndian.Uint32(buf[4:8]))
				zxid := int64(binary.BigEndian.Uint64(buf[8:16]))
				errnum := int32(binary.BigEndian.Uint32(buf[16:20]))
				if xid != -1 || zxid != -1 ***REMOVED***
					requestsLock.Lock()
					found := false
					opcode, found = requests[xid]
					if !found ***REMOVED***
						opcode = 0
					***REMOVED***
					delete(requests, xid)
					requestsLock.Unlock()
				***REMOVED*** else ***REMOVED***
					opcode = opWatcherEvent
				***REMOVED***
				cr = responseStructForOp(opcode)
				if cr == nil ***REMOVED***
					fmt.Printf("Unknown opcode %d\n", opcode)
				***REMOVED***
				if errnum != 0 ***REMOVED***
					cr = &struct***REMOVED******REMOVED******REMOVED******REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		opname := "."
		if opcode != -1 ***REMOVED***
			opname = opNames[opcode]
		***REMOVED***
		if cr == nil ***REMOVED***
			fmt.Printf("%+v %s %+v\n", client, opname, buf[4:4+blen])
		***REMOVED*** else ***REMOVED***
			n := 4
			hdrStr := ""
			if readHeader ***REMOVED***
				var hdr interface***REMOVED******REMOVED***
				if client ***REMOVED***
					hdr = &requestHeader***REMOVED******REMOVED***
				***REMOVED*** else ***REMOVED***
					hdr = &responseHeader***REMOVED******REMOVED***
				***REMOVED***
				if n2, err := decodePacket(buf[n:n+blen], hdr); err != nil ***REMOVED***
					fmt.Println(err)
				***REMOVED*** else ***REMOVED***
					n += n2
				***REMOVED***
				hdrStr = fmt.Sprintf(" %+v", hdr)
			***REMOVED***
			if _, err := decodePacket(buf[n:n+blen], cr); err != nil ***REMOVED***
				fmt.Println(err)
			***REMOVED***
			fmt.Printf("%+v %s%s %+v\n", client, opname, hdrStr, cr)
		***REMOVED***

		init = false

		written, err := conn2.Write(buf[:4+blen])
		if err != nil ***REMOVED***
			fmt.Println("3>", client, err)
			return
		***REMOVED*** else if written != 4+blen ***REMOVED***
			fmt.Printf("Written != read: %d != %d\n", written, blen)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func handleConnection(addr string, conn net.Conn) ***REMOVED***
	zkConn, err := net.Dial("tcp", addr)
	if err != nil ***REMOVED***
		fmt.Println(err)
		return
	***REMOVED***
	go trace(conn, zkConn, true)
	trace(zkConn, conn, false)
***REMOVED***

func StartTracer(listenAddr, serverAddr string) ***REMOVED***
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	for ***REMOVED***
		conn, err := ln.Accept()
		if err != nil ***REMOVED***
			fmt.Println(err)
			continue
		***REMOVED***
		go handleConnection(serverAddr, conn)
	***REMOVED***
***REMOVED***
