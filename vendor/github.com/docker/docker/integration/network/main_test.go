package network

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/internal/test/environment"
)

var testEnv *environment.Execution

func TestMain(m *testing.M) ***REMOVED***
	var err error
	testEnv, err = environment.New()
	if err != nil ***REMOVED***
		fmt.Println(err)
		os.Exit(1)
	***REMOVED***
	err = environment.EnsureFrozenImagesLinux(testEnv)
	if err != nil ***REMOVED***
		fmt.Println(err)
		os.Exit(1)
	***REMOVED***

	testEnv.Print()
	os.Exit(m.Run())
***REMOVED***

func setupTest(t *testing.T) func() ***REMOVED***
	environment.ProtectAll(t, testEnv)
	return func() ***REMOVED*** testEnv.Clean(t) ***REMOVED***
***REMOVED***
