package main

import (
	"fmt"
	"os"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func main() ***REMOVED***
	var f *os.File
	var err error

	if len(os.Args) < 2 ***REMOVED***
		fmt.Println("please supply filename(s)")
		os.Exit(1)
	***REMOVED***

	for _, fn := range os.Args[1:] ***REMOVED***
		f, err = os.Open(fn)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		defer f.Close()

		result, err := parser.Parse(f)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		fmt.Println(result.AST.Dump())
	***REMOVED***
***REMOVED***
