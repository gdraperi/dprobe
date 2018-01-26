package environment

import (
	"os"

	"os/exec"

	"github.com/docker/docker/internal/test/environment"
)

var (
	// DefaultClientBinary is the name of the docker binary
	DefaultClientBinary = os.Getenv("TEST_CLIENT_BINARY")
)

func init() ***REMOVED***
	if DefaultClientBinary == "" ***REMOVED***
		DefaultClientBinary = "docker"
	***REMOVED***
***REMOVED***

// Execution contains information about the current test execution and daemon
// under test
type Execution struct ***REMOVED***
	environment.Execution
	dockerBinary string
***REMOVED***

// DockerBinary returns the docker binary for this testing environment
func (e *Execution) DockerBinary() string ***REMOVED***
	return e.dockerBinary
***REMOVED***

// New returns details about the testing environment
func New() (*Execution, error) ***REMOVED***
	env, err := environment.New()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	dockerBinary, err := exec.LookPath(DefaultClientBinary)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Execution***REMOVED***
		Execution:    *env,
		dockerBinary: dockerBinary,
	***REMOVED***, nil
***REMOVED***
