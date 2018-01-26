package dockerversion

import (
	"fmt"
	"runtime"

	"github.com/docker/docker/pkg/parsers/kernel"
	"github.com/docker/docker/pkg/useragent"
	"golang.org/x/net/context"
)

// UAStringKey is used as key type for user-agent string in net/context struct
const UAStringKey = "upstream-user-agent"

// DockerUserAgent is the User-Agent the Docker client uses to identify itself.
// In accordance with RFC 7231 (5.5.3) is of the form:
//    [docker client's UA] UpstreamClient([upstream client's UA])
func DockerUserAgent(ctx context.Context) string ***REMOVED***
	httpVersion := make([]useragent.VersionInfo, 0, 6)
	httpVersion = append(httpVersion, useragent.VersionInfo***REMOVED***Name: "docker", Version: Version***REMOVED***)
	httpVersion = append(httpVersion, useragent.VersionInfo***REMOVED***Name: "go", Version: runtime.Version()***REMOVED***)
	httpVersion = append(httpVersion, useragent.VersionInfo***REMOVED***Name: "git-commit", Version: GitCommit***REMOVED***)
	if kernelVersion, err := kernel.GetKernelVersion(); err == nil ***REMOVED***
		httpVersion = append(httpVersion, useragent.VersionInfo***REMOVED***Name: "kernel", Version: kernelVersion.String()***REMOVED***)
	***REMOVED***
	httpVersion = append(httpVersion, useragent.VersionInfo***REMOVED***Name: "os", Version: runtime.GOOS***REMOVED***)
	httpVersion = append(httpVersion, useragent.VersionInfo***REMOVED***Name: "arch", Version: runtime.GOARCH***REMOVED***)

	dockerUA := useragent.AppendVersions("", httpVersion...)
	upstreamUA := getUserAgentFromContext(ctx)
	if len(upstreamUA) > 0 ***REMOVED***
		ret := insertUpstreamUserAgent(upstreamUA, dockerUA)
		return ret
	***REMOVED***
	return dockerUA
***REMOVED***

// getUserAgentFromContext returns the previously saved user-agent context stored in ctx, if one exists
func getUserAgentFromContext(ctx context.Context) string ***REMOVED***
	var upstreamUA string
	if ctx != nil ***REMOVED***
		var ki interface***REMOVED******REMOVED*** = ctx.Value(UAStringKey)
		if ki != nil ***REMOVED***
			upstreamUA = ctx.Value(UAStringKey).(string)
		***REMOVED***
	***REMOVED***
	return upstreamUA
***REMOVED***

// escapeStr returns s with every rune in charsToEscape escaped by a backslash
func escapeStr(s string, charsToEscape string) string ***REMOVED***
	var ret string
	for _, currRune := range s ***REMOVED***
		appended := false
		for _, escapableRune := range charsToEscape ***REMOVED***
			if currRune == escapableRune ***REMOVED***
				ret += `\` + string(currRune)
				appended = true
				break
			***REMOVED***
		***REMOVED***
		if !appended ***REMOVED***
			ret += string(currRune)
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

// insertUpstreamUserAgent adds the upstream client useragent to create a user-agent
// string of the form:
//   $dockerUA UpstreamClient($upstreamUA)
func insertUpstreamUserAgent(upstreamUA string, dockerUA string) string ***REMOVED***
	charsToEscape := `();\`
	upstreamUAEscaped := escapeStr(upstreamUA, charsToEscape)
	return fmt.Sprintf("%s UpstreamClient(%s)", dockerUA, upstreamUAEscaped)
***REMOVED***
