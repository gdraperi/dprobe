// +build !windows

package chrootarchive

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/pkg/reexec"
)

func init() ***REMOVED***
	reexec.Register("docker-applyLayer", applyLayer)
	reexec.Register("docker-untar", untar)
***REMOVED***

func fatal(err error) ***REMOVED***
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
***REMOVED***

// flush consumes all the bytes from the reader discarding
// any errors
func flush(r io.Reader) (bytes int64, err error) ***REMOVED***
	return io.Copy(ioutil.Discard, r)
***REMOVED***
