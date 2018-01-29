package logrus

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestRegister(t *testing.T) ***REMOVED***
	current := len(handlers)
	RegisterExitHandler(func() ***REMOVED******REMOVED***)
	if len(handlers) != current+1 ***REMOVED***
		t.Fatalf("expected %d handlers, got %d", current+1, len(handlers))
	***REMOVED***
***REMOVED***

func TestHandler(t *testing.T) ***REMOVED***
	tempDir, err := ioutil.TempDir("", "test_handler")
	if err != nil ***REMOVED***
		log.Fatalf("can't create temp dir. %q", err)
	***REMOVED***
	defer os.RemoveAll(tempDir)

	gofile := filepath.Join(tempDir, "gofile.go")
	if err := ioutil.WriteFile(gofile, testprog, 0666); err != nil ***REMOVED***
		t.Fatalf("can't create go file. %q", err)
	***REMOVED***

	outfile := filepath.Join(tempDir, "outfile.out")
	arg := time.Now().UTC().String()
	err = exec.Command("go", "run", gofile, outfile, arg).Run()
	if err == nil ***REMOVED***
		t.Fatalf("completed normally, should have failed")
	***REMOVED***

	data, err := ioutil.ReadFile(outfile)
	if err != nil ***REMOVED***
		t.Fatalf("can't read output file %s. %q", outfile, err)
	***REMOVED***

	if string(data) != arg ***REMOVED***
		t.Fatalf("bad data. Expected %q, got %q", data, arg)
	***REMOVED***
***REMOVED***

var testprog = []byte(`
// Test program for atexit, gets output file and data as arguments and writes
// data to output file in atexit handler.
package main

import (
	"github.com/sirupsen/logrus"
	"flag"
	"fmt"
	"io/ioutil"
)

var outfile = ""
var data = ""

func handler() ***REMOVED***
	ioutil.WriteFile(outfile, []byte(data), 0666)
***REMOVED***

func badHandler() ***REMOVED***
	n := 0
	fmt.Println(1/n)
***REMOVED***

func main() ***REMOVED***
	flag.Parse()
	outfile = flag.Arg(0)
	data = flag.Arg(1)

	logrus.RegisterExitHandler(handler)
	logrus.RegisterExitHandler(badHandler)
	logrus.Fatal("Bye bye")
***REMOVED***
`)
