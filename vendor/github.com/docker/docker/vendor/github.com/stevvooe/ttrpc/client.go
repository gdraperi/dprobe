package ttrpc

import (
	"context"
	"net"
	"sync"

	"github.com/containerd/containerd/log"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

type Client struct ***REMOVED***
	codec   codec
	conn    net.Conn
	channel *channel
	calls   chan *callRequest

	closed    chan struct***REMOVED******REMOVED***
	closeOnce sync.Once
	done      chan struct***REMOVED******REMOVED***
	err       error
***REMOVED***

func NewClient(conn net.Conn) *Client ***REMOVED***
	c := &Client***REMOVED***
		codec:   codec***REMOVED******REMOVED***,
		conn:    conn,
		channel: newChannel(conn, conn),
		calls:   make(chan *callRequest),
		closed:  make(chan struct***REMOVED******REMOVED***),
		done:    make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	go c.run()
	return c
***REMOVED***

type callRequest struct ***REMOVED***
	ctx  context.Context
	req  *Request
	resp *Response  // response will be written back here
	errs chan error // error written here on completion
***REMOVED***

func (c *Client) Call(ctx context.Context, service, method string, req, resp interface***REMOVED******REMOVED***) error ***REMOVED***
	payload, err := c.codec.Marshal(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var (
		creq = &Request***REMOVED***
			Service: service,
			Method:  method,
			Payload: payload,
		***REMOVED***

		cresp = &Response***REMOVED******REMOVED***
	)

	if err := c.dispatch(ctx, creq, cresp); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := c.codec.Unmarshal(cresp.Payload, resp); err != nil ***REMOVED***
		return err
	***REMOVED***

	if cresp.Status == nil ***REMOVED***
		return errors.New("no status provided on response")
	***REMOVED***

	return status.ErrorProto(cresp.Status)
***REMOVED***

func (c *Client) dispatch(ctx context.Context, req *Request, resp *Response) error ***REMOVED***
	errs := make(chan error, 1)
	call := &callRequest***REMOVED***
		req:  req,
		resp: resp,
		errs: errs,
	***REMOVED***

	select ***REMOVED***
	case c.calls <- call:
	case <-c.done:
		return c.err
	***REMOVED***

	select ***REMOVED***
	case err := <-errs:
		return err
	case <-c.done:
		return c.err
	***REMOVED***
***REMOVED***

func (c *Client) Close() error ***REMOVED***
	c.closeOnce.Do(func() ***REMOVED***
		close(c.closed)
	***REMOVED***)

	return nil
***REMOVED***

type message struct ***REMOVED***
	messageHeader
	p   []byte
	err error
***REMOVED***

func (c *Client) run() ***REMOVED***
	var (
		streamID    uint32 = 1
		waiters            = make(map[uint32]*callRequest)
		calls              = c.calls
		incoming           = make(chan *message)
		shutdown           = make(chan struct***REMOVED******REMOVED***)
		shutdownErr error
	)

	go func() ***REMOVED***
		defer close(shutdown)

		// start one more goroutine to recv messages without blocking.
		for ***REMOVED***
			mh, p, err := c.channel.recv(context.TODO())
			if err != nil ***REMOVED***
				_, ok := status.FromError(err)
				if !ok ***REMOVED***
					// treat all errors that are not an rpc status as terminal.
					// all others poison the connection.
					shutdownErr = err
					return
				***REMOVED***
			***REMOVED***
			select ***REMOVED***
			case incoming <- &message***REMOVED***
				messageHeader: mh,
				p:             p[:mh.Length],
				err:           err,
			***REMOVED***:
			case <-c.done:
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	defer c.conn.Close()
	defer close(c.done)

	for ***REMOVED***
		select ***REMOVED***
		case call := <-calls:
			if err := c.send(call.ctx, streamID, messageTypeRequest, call.req); err != nil ***REMOVED***
				call.errs <- err
				continue
			***REMOVED***

			waiters[streamID] = call
			streamID += 2 // enforce odd client initiated request ids
		case msg := <-incoming:
			call, ok := waiters[msg.StreamID]
			if !ok ***REMOVED***
				log.L.Errorf("ttrpc: received message for unknown channel %v", msg.StreamID)
				continue
			***REMOVED***

			call.errs <- c.recv(call.resp, msg)
			delete(waiters, msg.StreamID)
		case <-shutdown:
			shutdownErr = errors.Wrapf(shutdownErr, "ttrpc: client shutting down")
			c.err = shutdownErr
			for _, waiter := range waiters ***REMOVED***
				waiter.errs <- shutdownErr
			***REMOVED***
			c.Close()
			return
		case <-c.closed:
			// broadcast the shutdown error to the remaining waiters.
			for _, waiter := range waiters ***REMOVED***
				waiter.errs <- shutdownErr
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Client) send(ctx context.Context, streamID uint32, mtype messageType, msg interface***REMOVED******REMOVED***) error ***REMOVED***
	p, err := c.codec.Marshal(msg)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return c.channel.send(ctx, streamID, mtype, p)
***REMOVED***

func (c *Client) recv(resp *Response, msg *message) error ***REMOVED***
	if msg.err != nil ***REMOVED***
		return msg.err
	***REMOVED***

	if msg.Type != messageTypeResponse ***REMOVED***
		return errors.New("unkown message type received")
	***REMOVED***

	defer c.channel.putmbuf(msg.p)
	return proto.Unmarshal(msg.p, resp)
***REMOVED***
