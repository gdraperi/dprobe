//+build !windows

package chrootarchive

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/docker/pkg/system"
	rsystem "github.com/opencontainers/runc/libcontainer/system"
)

type applyLayerResponse struct ***REMOVED***
	LayerSize int64 `json:"layerSize"`
***REMOVED***

// applyLayer is the entry-point for docker-applylayer on re-exec. This is not
// used on Windows as it does not support chroot, hence no point sandboxing
// through chroot and rexec.
func applyLayer() ***REMOVED***

	var (
		tmpDir  string
		err     error
		options *archive.TarOptions
	)
	runtime.LockOSThread()
	flag.Parse()

	inUserns := rsystem.RunningInUserNS()
	if err := chroot(flag.Arg(0)); err != nil ***REMOVED***
		fatal(err)
	***REMOVED***

	// We need to be able to set any perms
	oldmask, err := system.Umask(0)
	defer system.Umask(oldmask)
	if err != nil ***REMOVED***
		fatal(err)
	***REMOVED***

	if err := json.Unmarshal([]byte(os.Getenv("OPT")), &options); err != nil ***REMOVED***
		fatal(err)
	***REMOVED***

	if inUserns ***REMOVED***
		options.InUserNS = true
	***REMOVED***

	if tmpDir, err = ioutil.TempDir("/", "temp-docker-extract"); err != nil ***REMOVED***
		fatal(err)
	***REMOVED***

	os.Setenv("TMPDIR", tmpDir)
	size, err := archive.UnpackLayer("/", os.Stdin, options)
	os.RemoveAll(tmpDir)
	if err != nil ***REMOVED***
		fatal(err)
	***REMOVED***

	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(applyLayerResponse***REMOVED***size***REMOVED***); err != nil ***REMOVED***
		fatal(fmt.Errorf("unable to encode layerSize JSON: %s", err))
	***REMOVED***

	if _, err := flush(os.Stdin); err != nil ***REMOVED***
		fatal(err)
	***REMOVED***

	os.Exit(0)
***REMOVED***

// applyLayerHandler parses a diff in the standard layer format from `layer`, and
// applies it to the directory `dest`. Returns the size in bytes of the
// contents of the layer.
func applyLayerHandler(dest string, layer io.Reader, options *archive.TarOptions, decompress bool) (size int64, err error) ***REMOVED***
	dest = filepath.Clean(dest)
	if decompress ***REMOVED***
		decompressed, err := archive.DecompressStream(layer)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		defer decompressed.Close()

		layer = decompressed
	***REMOVED***
	if options == nil ***REMOVED***
		options = &archive.TarOptions***REMOVED******REMOVED***
		if rsystem.RunningInUserNS() ***REMOVED***
			options.InUserNS = true
		***REMOVED***
	***REMOVED***
	if options.ExcludePatterns == nil ***REMOVED***
		options.ExcludePatterns = []string***REMOVED******REMOVED***
	***REMOVED***

	data, err := json.Marshal(options)
	if err != nil ***REMOVED***
		return 0, fmt.Errorf("ApplyLayer json encode: %v", err)
	***REMOVED***

	cmd := reexec.Command("docker-applyLayer", dest)
	cmd.Stdin = layer
	cmd.Env = append(cmd.Env, fmt.Sprintf("OPT=%s", data))

	outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout, cmd.Stderr = outBuf, errBuf

	if err = cmd.Run(); err != nil ***REMOVED***
		return 0, fmt.Errorf("ApplyLayer %s stdout: %s stderr: %s", err, outBuf, errBuf)
	***REMOVED***

	// Stdout should be a valid JSON struct representing an applyLayerResponse.
	response := applyLayerResponse***REMOVED******REMOVED***
	decoder := json.NewDecoder(outBuf)
	if err = decoder.Decode(&response); err != nil ***REMOVED***
		return 0, fmt.Errorf("unable to decode ApplyLayer JSON response: %s", err)
	***REMOVED***

	return response.LayerSize, nil
***REMOVED***
