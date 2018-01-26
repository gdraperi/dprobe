package client

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// ImageBuild sends request to the daemon to build images.
// The Body in the response implement an io.ReadCloser and it's up to the caller to
// close it.
func (cli *Client) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) ***REMOVED***
	query, err := cli.imageBuildOptionsToQuery(options)
	if err != nil ***REMOVED***
		return types.ImageBuildResponse***REMOVED******REMOVED***, err
	***REMOVED***

	headers := http.Header(make(map[string][]string))
	buf, err := json.Marshal(options.AuthConfigs)
	if err != nil ***REMOVED***
		return types.ImageBuildResponse***REMOVED******REMOVED***, err
	***REMOVED***
	headers.Add("X-Registry-Config", base64.URLEncoding.EncodeToString(buf))

	if options.Platform != "" ***REMOVED***
		if err := cli.NewVersionError("1.32", "platform"); err != nil ***REMOVED***
			return types.ImageBuildResponse***REMOVED******REMOVED***, err
		***REMOVED***
		query.Set("platform", options.Platform)
	***REMOVED***
	headers.Set("Content-Type", "application/x-tar")

	serverResp, err := cli.postRaw(ctx, "/build", query, buildContext, headers)
	if err != nil ***REMOVED***
		return types.ImageBuildResponse***REMOVED******REMOVED***, err
	***REMOVED***

	osType := getDockerOS(serverResp.header.Get("Server"))

	return types.ImageBuildResponse***REMOVED***
		Body:   serverResp.body,
		OSType: osType,
	***REMOVED***, nil
***REMOVED***

func (cli *Client) imageBuildOptionsToQuery(options types.ImageBuildOptions) (url.Values, error) ***REMOVED***
	query := url.Values***REMOVED***
		"t":           options.Tags,
		"securityopt": options.SecurityOpt,
		"extrahosts":  options.ExtraHosts,
	***REMOVED***
	if options.SuppressOutput ***REMOVED***
		query.Set("q", "1")
	***REMOVED***
	if options.RemoteContext != "" ***REMOVED***
		query.Set("remote", options.RemoteContext)
	***REMOVED***
	if options.NoCache ***REMOVED***
		query.Set("nocache", "1")
	***REMOVED***
	if options.Remove ***REMOVED***
		query.Set("rm", "1")
	***REMOVED*** else ***REMOVED***
		query.Set("rm", "0")
	***REMOVED***

	if options.ForceRemove ***REMOVED***
		query.Set("forcerm", "1")
	***REMOVED***

	if options.PullParent ***REMOVED***
		query.Set("pull", "1")
	***REMOVED***

	if options.Squash ***REMOVED***
		if err := cli.NewVersionError("1.25", "squash"); err != nil ***REMOVED***
			return query, err
		***REMOVED***
		query.Set("squash", "1")
	***REMOVED***

	if !container.Isolation.IsDefault(options.Isolation) ***REMOVED***
		query.Set("isolation", string(options.Isolation))
	***REMOVED***

	query.Set("cpusetcpus", options.CPUSetCPUs)
	query.Set("networkmode", options.NetworkMode)
	query.Set("cpusetmems", options.CPUSetMems)
	query.Set("cpushares", strconv.FormatInt(options.CPUShares, 10))
	query.Set("cpuquota", strconv.FormatInt(options.CPUQuota, 10))
	query.Set("cpuperiod", strconv.FormatInt(options.CPUPeriod, 10))
	query.Set("memory", strconv.FormatInt(options.Memory, 10))
	query.Set("memswap", strconv.FormatInt(options.MemorySwap, 10))
	query.Set("cgroupparent", options.CgroupParent)
	query.Set("shmsize", strconv.FormatInt(options.ShmSize, 10))
	query.Set("dockerfile", options.Dockerfile)
	query.Set("target", options.Target)

	ulimitsJSON, err := json.Marshal(options.Ulimits)
	if err != nil ***REMOVED***
		return query, err
	***REMOVED***
	query.Set("ulimits", string(ulimitsJSON))

	buildArgsJSON, err := json.Marshal(options.BuildArgs)
	if err != nil ***REMOVED***
		return query, err
	***REMOVED***
	query.Set("buildargs", string(buildArgsJSON))

	labelsJSON, err := json.Marshal(options.Labels)
	if err != nil ***REMOVED***
		return query, err
	***REMOVED***
	query.Set("labels", string(labelsJSON))

	cacheFromJSON, err := json.Marshal(options.CacheFrom)
	if err != nil ***REMOVED***
		return query, err
	***REMOVED***
	query.Set("cachefrom", string(cacheFromJSON))
	if options.SessionID != "" ***REMOVED***
		query.Set("session", options.SessionID)
	***REMOVED***
	if options.Platform != "" ***REMOVED***
		query.Set("platform", strings.ToLower(options.Platform))
	***REMOVED***
	return query, nil
***REMOVED***
