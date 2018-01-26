package main

import (
	"os"
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

// export an image and try to import it into a new one
func (s *DockerSuite) TestExportContainerAndImportImage(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	containerID := "testexportcontainerandimportimage"

	dockerCmd(c, "run", "--name", containerID, "busybox", "true")

	out, _ := dockerCmd(c, "export", containerID)

	result := icmd.RunCmd(icmd.Cmd***REMOVED***
		Command: []string***REMOVED***dockerBinary, "import", "-", "repo/testexp:v1"***REMOVED***,
		Stdin:   strings.NewReader(out),
	***REMOVED***)
	result.Assert(c, icmd.Success)

	cleanedImageID := strings.TrimSpace(result.Combined())
	c.Assert(cleanedImageID, checker.Not(checker.Equals), "", check.Commentf("output should have been an image id"))
***REMOVED***

// Used to test output flag in the export command
func (s *DockerSuite) TestExportContainerWithOutputAndImportImage(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	containerID := "testexportcontainerwithoutputandimportimage"

	dockerCmd(c, "run", "--name", containerID, "busybox", "true")
	dockerCmd(c, "export", "--output=testexp.tar", containerID)
	defer os.Remove("testexp.tar")

	resultCat := icmd.RunCommand("cat", "testexp.tar")
	resultCat.Assert(c, icmd.Success)

	result := icmd.RunCmd(icmd.Cmd***REMOVED***
		Command: []string***REMOVED***dockerBinary, "import", "-", "repo/testexp:v1"***REMOVED***,
		Stdin:   strings.NewReader(resultCat.Combined()),
	***REMOVED***)
	result.Assert(c, icmd.Success)

	cleanedImageID := strings.TrimSpace(result.Combined())
	c.Assert(cleanedImageID, checker.Not(checker.Equals), "", check.Commentf("output should have been an image id"))
***REMOVED***
