package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func do_client() int ***REMOVED***
	addr := *g_addr
	if *g_sock == "unix" ***REMOVED***
		addr = get_socket_filename()
	***REMOVED***

	// client
	client, err := rpc.Dial(*g_sock, addr)
	if err != nil ***REMOVED***
		if *g_sock == "unix" && file_exists(addr) ***REMOVED***
			os.Remove(addr)
		***REMOVED***

		err = try_run_server()
		if err != nil ***REMOVED***
			fmt.Printf("%s\n", err.Error())
			return 1
		***REMOVED***
		client, err = try_to_connect(*g_sock, addr)
		if err != nil ***REMOVED***
			fmt.Printf("%s\n", err.Error())
			return 1
		***REMOVED***
	***REMOVED***
	defer client.Close()

	if flag.NArg() > 0 ***REMOVED***
		switch flag.Arg(0) ***REMOVED***
		case "autocomplete":
			cmd_auto_complete(client)
		case "close":
			cmd_close(client)
		case "status":
			cmd_status(client)
		case "drop-cache":
			cmd_drop_cache(client)
		case "set":
			cmd_set(client)
		case "options":
			cmd_options(client)
		default:
			fmt.Printf("unknown argument: %q, try running \"gocode -h\"\n", flag.Arg(0))
			return 1
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func try_run_server() error ***REMOVED***
	path := get_executable_filename()
	args := []string***REMOVED***os.Args[0], "-s", "-sock", *g_sock, "-addr", *g_addr***REMOVED***
	cwd, _ := os.Getwd()

	var err error
	stdin, err := os.Open(os.DevNull)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	stdout, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	stderr, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	procattr := os.ProcAttr***REMOVED***Dir: cwd, Env: os.Environ(), Files: []*os.File***REMOVED***stdin, stdout, stderr***REMOVED******REMOVED***
	p, err := os.StartProcess(path, args, &procattr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return p.Release()
***REMOVED***

func try_to_connect(network, address string) (client *rpc.Client, err error) ***REMOVED***
	t := 0
	for ***REMOVED***
		client, err = rpc.Dial(network, address)
		if err != nil && t < 1000 ***REMOVED***
			time.Sleep(10 * time.Millisecond)
			t += 10
			continue
		***REMOVED***
		break
	***REMOVED***

	return
***REMOVED***

func prepare_file_filename_cursor() ([]byte, string, int) ***REMOVED***
	var file []byte
	var err error

	if *g_input != "" ***REMOVED***
		file, err = ioutil.ReadFile(*g_input)
	***REMOVED*** else ***REMOVED***
		file, err = ioutil.ReadAll(os.Stdin)
	***REMOVED***

	if err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***

	var skipped int
	file, skipped = filter_out_shebang(file)

	filename := *g_input
	cursor := -1

	offset := ""
	switch flag.NArg() ***REMOVED***
	case 2:
		offset = flag.Arg(1)
	case 3:
		filename = flag.Arg(1) // Override default filename
		offset = flag.Arg(2)
	***REMOVED***

	if offset != "" ***REMOVED***
		if offset[0] == 'c' || offset[0] == 'C' ***REMOVED***
			cursor, _ = strconv.Atoi(offset[1:])
			cursor = char_to_byte_offset(file, cursor)
		***REMOVED*** else ***REMOVED***
			cursor, _ = strconv.Atoi(offset)
		***REMOVED***
	***REMOVED***

	cursor -= skipped
	if filename != "" && !filepath.IsAbs(filename) ***REMOVED***
		cwd, _ := os.Getwd()
		filename = filepath.Join(cwd, filename)
	***REMOVED***
	return file, filename, cursor
***REMOVED***

//-------------------------------------------------------------------------
// commands
//-------------------------------------------------------------------------

func cmd_status(c *rpc.Client) ***REMOVED***
	fmt.Printf("%s\n", client_status(c, 0))
***REMOVED***

func cmd_auto_complete(c *rpc.Client) ***REMOVED***
	context := pack_build_context(&build.Default)
	file, filename, cursor := prepare_file_filename_cursor()
	f := get_formatter(*g_format)
	f.write_candidates(client_auto_complete(c, file, filename, cursor, context))
***REMOVED***

func cmd_close(c *rpc.Client) ***REMOVED***
	client_close(c, 0)
***REMOVED***

func cmd_drop_cache(c *rpc.Client) ***REMOVED***
	client_drop_cache(c, 0)
***REMOVED***

func cmd_set(c *rpc.Client) ***REMOVED***
	switch flag.NArg() ***REMOVED***
	case 1:
		fmt.Print(client_set(c, "\x00", "\x00"))
	case 2:
		fmt.Print(client_set(c, flag.Arg(1), "\x00"))
	case 3:
		fmt.Print(client_set(c, flag.Arg(1), flag.Arg(2)))
	***REMOVED***
***REMOVED***

func cmd_options(c *rpc.Client) ***REMOVED***
	fmt.Print(client_options(c, 0))
***REMOVED***
