package dockerfile

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func strPtr(source string) *string ***REMOVED***
	return &source
***REMOVED***

func TestGetAllAllowed(t *testing.T) ***REMOVED***
	buildArgs := newBuildArgs(map[string]*string***REMOVED***
		"ArgNotUsedInDockerfile":              strPtr("fromopt1"),
		"ArgOverriddenByOptions":              strPtr("fromopt2"),
		"ArgNoDefaultInDockerfileFromOptions": strPtr("fromopt3"),
		"HTTP_PROXY":                          strPtr("theproxy"),
	***REMOVED***)

	buildArgs.AddMetaArg("ArgFromMeta", strPtr("frommeta1"))
	buildArgs.AddMetaArg("ArgFromMetaOverridden", strPtr("frommeta2"))
	buildArgs.AddMetaArg("ArgFromMetaNotUsed", strPtr("frommeta3"))

	buildArgs.AddArg("ArgOverriddenByOptions", strPtr("fromdockerfile2"))
	buildArgs.AddArg("ArgWithDefaultInDockerfile", strPtr("fromdockerfile1"))
	buildArgs.AddArg("ArgNoDefaultInDockerfile", nil)
	buildArgs.AddArg("ArgNoDefaultInDockerfileFromOptions", nil)
	buildArgs.AddArg("ArgFromMeta", nil)
	buildArgs.AddArg("ArgFromMetaOverridden", strPtr("fromdockerfile3"))

	all := buildArgs.GetAllAllowed()
	expected := map[string]string***REMOVED***
		"HTTP_PROXY":                          "theproxy",
		"ArgOverriddenByOptions":              "fromopt2",
		"ArgWithDefaultInDockerfile":          "fromdockerfile1",
		"ArgNoDefaultInDockerfileFromOptions": "fromopt3",
		"ArgFromMeta":                         "frommeta1",
		"ArgFromMetaOverridden":               "fromdockerfile3",
	***REMOVED***
	assert.Equal(t, expected, all)
***REMOVED***

func TestGetAllMeta(t *testing.T) ***REMOVED***
	buildArgs := newBuildArgs(map[string]*string***REMOVED***
		"ArgNotUsedInDockerfile":        strPtr("fromopt1"),
		"ArgOverriddenByOptions":        strPtr("fromopt2"),
		"ArgNoDefaultInMetaFromOptions": strPtr("fromopt3"),
		"HTTP_PROXY":                    strPtr("theproxy"),
	***REMOVED***)

	buildArgs.AddMetaArg("ArgFromMeta", strPtr("frommeta1"))
	buildArgs.AddMetaArg("ArgOverriddenByOptions", strPtr("frommeta2"))
	buildArgs.AddMetaArg("ArgNoDefaultInMetaFromOptions", nil)

	all := buildArgs.GetAllMeta()
	expected := map[string]string***REMOVED***
		"HTTP_PROXY":                    "theproxy",
		"ArgFromMeta":                   "frommeta1",
		"ArgOverriddenByOptions":        "fromopt2",
		"ArgNoDefaultInMetaFromOptions": "fromopt3",
	***REMOVED***
	assert.Equal(t, expected, all)
***REMOVED***

func TestWarnOnUnusedBuildArgs(t *testing.T) ***REMOVED***
	buildArgs := newBuildArgs(map[string]*string***REMOVED***
		"ThisArgIsUsed":    strPtr("fromopt1"),
		"ThisArgIsNotUsed": strPtr("fromopt2"),
		"HTTPS_PROXY":      strPtr("referenced builtin"),
		"HTTP_PROXY":       strPtr("unreferenced builtin"),
	***REMOVED***)
	buildArgs.AddArg("ThisArgIsUsed", nil)
	buildArgs.AddArg("HTTPS_PROXY", nil)

	buffer := new(bytes.Buffer)
	buildArgs.WarnOnUnusedBuildArgs(buffer)
	out := buffer.String()
	assert.NotContains(t, out, "ThisArgIsUsed")
	assert.NotContains(t, out, "HTTPS_PROXY")
	assert.NotContains(t, out, "HTTP_PROXY")
	assert.Contains(t, out, "ThisArgIsNotUsed")
***REMOVED***

func TestIsUnreferencedBuiltin(t *testing.T) ***REMOVED***
	buildArgs := newBuildArgs(map[string]*string***REMOVED***
		"ThisArgIsUsed":    strPtr("fromopt1"),
		"ThisArgIsNotUsed": strPtr("fromopt2"),
		"HTTPS_PROXY":      strPtr("referenced builtin"),
		"HTTP_PROXY":       strPtr("unreferenced builtin"),
	***REMOVED***)
	buildArgs.AddArg("ThisArgIsUsed", nil)
	buildArgs.AddArg("HTTPS_PROXY", nil)

	assert.True(t, buildArgs.IsReferencedOrNotBuiltin("ThisArgIsUsed"))
	assert.True(t, buildArgs.IsReferencedOrNotBuiltin("ThisArgIsNotUsed"))
	assert.True(t, buildArgs.IsReferencedOrNotBuiltin("HTTPS_PROXY"))
	assert.False(t, buildArgs.IsReferencedOrNotBuiltin("HTTP_PROXY"))
***REMOVED***
