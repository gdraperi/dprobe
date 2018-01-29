// WARNING! Autogenerated by goremote, don't touch.

package main

import (
	"net/rpc"
)

type RPC struct ***REMOVED***
***REMOVED***

// wrapper for: server_auto_complete

type Args_auto_complete struct ***REMOVED***
	Arg0 []byte
	Arg1 string
	Arg2 int
	Arg3 go_build_context
***REMOVED***
type Reply_auto_complete struct ***REMOVED***
	Arg0 []candidate
	Arg1 int
***REMOVED***

func (r *RPC) RPC_auto_complete(args *Args_auto_complete, reply *Reply_auto_complete) error ***REMOVED***
	reply.Arg0, reply.Arg1 = server_auto_complete(args.Arg0, args.Arg1, args.Arg2, args.Arg3)
	return nil
***REMOVED***
func client_auto_complete(cli *rpc.Client, Arg0 []byte, Arg1 string, Arg2 int, Arg3 go_build_context) (c []candidate, d int) ***REMOVED***
	var args Args_auto_complete
	var reply Reply_auto_complete
	args.Arg0 = Arg0
	args.Arg1 = Arg1
	args.Arg2 = Arg2
	args.Arg3 = Arg3
	err := cli.Call("RPC.RPC_auto_complete", &args, &reply)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return reply.Arg0, reply.Arg1
***REMOVED***

// wrapper for: server_close

type Args_close struct ***REMOVED***
	Arg0 int
***REMOVED***
type Reply_close struct ***REMOVED***
	Arg0 int
***REMOVED***

func (r *RPC) RPC_close(args *Args_close, reply *Reply_close) error ***REMOVED***
	reply.Arg0 = server_close(args.Arg0)
	return nil
***REMOVED***
func client_close(cli *rpc.Client, Arg0 int) int ***REMOVED***
	var args Args_close
	var reply Reply_close
	args.Arg0 = Arg0
	err := cli.Call("RPC.RPC_close", &args, &reply)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return reply.Arg0
***REMOVED***

// wrapper for: server_status

type Args_status struct ***REMOVED***
	Arg0 int
***REMOVED***
type Reply_status struct ***REMOVED***
	Arg0 string
***REMOVED***

func (r *RPC) RPC_status(args *Args_status, reply *Reply_status) error ***REMOVED***
	reply.Arg0 = server_status(args.Arg0)
	return nil
***REMOVED***
func client_status(cli *rpc.Client, Arg0 int) string ***REMOVED***
	var args Args_status
	var reply Reply_status
	args.Arg0 = Arg0
	err := cli.Call("RPC.RPC_status", &args, &reply)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return reply.Arg0
***REMOVED***

// wrapper for: server_drop_cache

type Args_drop_cache struct ***REMOVED***
	Arg0 int
***REMOVED***
type Reply_drop_cache struct ***REMOVED***
	Arg0 int
***REMOVED***

func (r *RPC) RPC_drop_cache(args *Args_drop_cache, reply *Reply_drop_cache) error ***REMOVED***
	reply.Arg0 = server_drop_cache(args.Arg0)
	return nil
***REMOVED***
func client_drop_cache(cli *rpc.Client, Arg0 int) int ***REMOVED***
	var args Args_drop_cache
	var reply Reply_drop_cache
	args.Arg0 = Arg0
	err := cli.Call("RPC.RPC_drop_cache", &args, &reply)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return reply.Arg0
***REMOVED***

// wrapper for: server_set

type Args_set struct ***REMOVED***
	Arg0, Arg1 string
***REMOVED***
type Reply_set struct ***REMOVED***
	Arg0 string
***REMOVED***

func (r *RPC) RPC_set(args *Args_set, reply *Reply_set) error ***REMOVED***
	reply.Arg0 = server_set(args.Arg0, args.Arg1)
	return nil
***REMOVED***
func client_set(cli *rpc.Client, Arg0, Arg1 string) string ***REMOVED***
	var args Args_set
	var reply Reply_set
	args.Arg0 = Arg0
	args.Arg1 = Arg1
	err := cli.Call("RPC.RPC_set", &args, &reply)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return reply.Arg0
***REMOVED***

// wrapper for: server_options

type Args_options struct ***REMOVED***
	Arg0 int
***REMOVED***
type Reply_options struct ***REMOVED***
	Arg0 string
***REMOVED***

func (r *RPC) RPC_options(args *Args_options, reply *Reply_options) error ***REMOVED***
	reply.Arg0 = server_options(args.Arg0)
	return nil
***REMOVED***
func client_options(cli *rpc.Client, Arg0 int) string ***REMOVED***
	var args Args_options
	var reply Reply_options
	args.Arg0 = Arg0
	err := cli.Call("RPC.RPC_options", &args, &reply)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return reply.Arg0
***REMOVED***
