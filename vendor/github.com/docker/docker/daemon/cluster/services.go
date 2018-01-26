package cluster

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	apitypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	types "github.com/docker/docker/api/types/swarm"
	timetypes "github.com/docker/docker/api/types/time"
	"github.com/docker/docker/daemon/cluster/convert"
	"github.com/docker/docker/errdefs"
	runconfigopts "github.com/docker/docker/runconfig/opts"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// GetServices returns all services of a managed swarm cluster.
func (c *Cluster) GetServices(options apitypes.ServiceListOptions) ([]types.Service, error) ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() ***REMOVED***
		return nil, c.errNoManager(state)
	***REMOVED***

	// We move the accepted filter check here as "mode" filter
	// is processed in the daemon, not in SwarmKit. So it might
	// be good to have accepted file check in the same file as
	// the filter processing (in the for loop below).
	accepted := map[string]bool***REMOVED***
		"name":    true,
		"id":      true,
		"label":   true,
		"mode":    true,
		"runtime": true,
	***REMOVED***
	if err := options.Filters.Validate(accepted); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(options.Filters.Get("runtime")) == 0 ***REMOVED***
		// Default to using the container runtime filter
		options.Filters.Add("runtime", string(types.RuntimeContainer))
	***REMOVED***

	filters := &swarmapi.ListServicesRequest_Filters***REMOVED***
		NamePrefixes: options.Filters.Get("name"),
		IDPrefixes:   options.Filters.Get("id"),
		Labels:       runconfigopts.ConvertKVStringsToMap(options.Filters.Get("label")),
		Runtimes:     options.Filters.Get("runtime"),
	***REMOVED***

	ctx, cancel := c.getRequestContext()
	defer cancel()

	r, err := state.controlClient.ListServices(
		ctx,
		&swarmapi.ListServicesRequest***REMOVED***Filters: filters***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	services := make([]types.Service, 0, len(r.Services))

	for _, service := range r.Services ***REMOVED***
		if options.Filters.Contains("mode") ***REMOVED***
			var mode string
			switch service.Spec.GetMode().(type) ***REMOVED***
			case *swarmapi.ServiceSpec_Global:
				mode = "global"
			case *swarmapi.ServiceSpec_Replicated:
				mode = "replicated"
			***REMOVED***

			if !options.Filters.ExactMatch("mode", mode) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		svcs, err := convert.ServiceFromGRPC(*service)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		services = append(services, svcs)
	***REMOVED***

	return services, nil
***REMOVED***

// GetService returns a service based on an ID or name.
func (c *Cluster) GetService(input string, insertDefaults bool) (types.Service, error) ***REMOVED***
	var service *swarmapi.Service
	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		s, err := getService(ctx, state.controlClient, input, insertDefaults)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		service = s
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return types.Service***REMOVED******REMOVED***, err
	***REMOVED***
	svc, err := convert.ServiceFromGRPC(*service)
	if err != nil ***REMOVED***
		return types.Service***REMOVED******REMOVED***, err
	***REMOVED***
	return svc, nil
***REMOVED***

// CreateService creates a new service in a managed swarm cluster.
func (c *Cluster) CreateService(s types.ServiceSpec, encodedAuth string, queryRegistry bool) (*apitypes.ServiceCreateResponse, error) ***REMOVED***
	var resp *apitypes.ServiceCreateResponse
	err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		err := c.populateNetworkID(ctx, state.controlClient, &s)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		serviceSpec, err := convert.ServiceSpecToGRPC(s)
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***

		resp = &apitypes.ServiceCreateResponse***REMOVED******REMOVED***

		switch serviceSpec.Task.Runtime.(type) ***REMOVED***
		// handle other runtimes here
		case *swarmapi.TaskSpec_Generic:
			switch serviceSpec.Task.GetGeneric().Kind ***REMOVED***
			case string(types.RuntimePlugin):
				info, _ := c.config.Backend.SystemInfo()
				if !info.ExperimentalBuild ***REMOVED***
					return fmt.Errorf("runtime type %q only supported in experimental", types.RuntimePlugin)
				***REMOVED***
				if s.TaskTemplate.PluginSpec == nil ***REMOVED***
					return errors.New("plugin spec must be set")
				***REMOVED***

			default:
				return fmt.Errorf("unsupported runtime type: %q", serviceSpec.Task.GetGeneric().Kind)
			***REMOVED***

			r, err := state.controlClient.CreateService(ctx, &swarmapi.CreateServiceRequest***REMOVED***Spec: &serviceSpec***REMOVED***)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			resp.ID = r.Service.ID
		case *swarmapi.TaskSpec_Container:
			ctnr := serviceSpec.Task.GetContainer()
			if ctnr == nil ***REMOVED***
				return errors.New("service does not use container tasks")
			***REMOVED***
			if encodedAuth != "" ***REMOVED***
				ctnr.PullOptions = &swarmapi.ContainerSpec_PullOptions***REMOVED***RegistryAuth: encodedAuth***REMOVED***
			***REMOVED***

			// retrieve auth config from encoded auth
			authConfig := &apitypes.AuthConfig***REMOVED******REMOVED***
			if encodedAuth != "" ***REMOVED***
				authReader := strings.NewReader(encodedAuth)
				dec := json.NewDecoder(base64.NewDecoder(base64.URLEncoding, authReader))
				if err := dec.Decode(authConfig); err != nil ***REMOVED***
					logrus.Warnf("invalid authconfig: %v", err)
				***REMOVED***
			***REMOVED***

			// pin image by digest for API versions < 1.30
			// TODO(nishanttotla): The check on "DOCKER_SERVICE_PREFER_OFFLINE_IMAGE"
			// should be removed in the future. Since integration tests only use the
			// latest API version, so this is no longer required.
			if os.Getenv("DOCKER_SERVICE_PREFER_OFFLINE_IMAGE") != "1" && queryRegistry ***REMOVED***
				digestImage, err := c.imageWithDigestString(ctx, ctnr.Image, authConfig)
				if err != nil ***REMOVED***
					logrus.Warnf("unable to pin image %s to digest: %s", ctnr.Image, err.Error())
					// warning in the client response should be concise
					resp.Warnings = append(resp.Warnings, digestWarning(ctnr.Image))

				***REMOVED*** else if ctnr.Image != digestImage ***REMOVED***
					logrus.Debugf("pinning image %s by digest: %s", ctnr.Image, digestImage)
					ctnr.Image = digestImage

				***REMOVED*** else ***REMOVED***
					logrus.Debugf("creating service using supplied digest reference %s", ctnr.Image)

				***REMOVED***

				// Replace the context with a fresh one.
				// If we timed out while communicating with the
				// registry, then "ctx" will already be expired, which
				// would cause UpdateService below to fail. Reusing
				// "ctx" could make it impossible to create a service
				// if the registry is slow or unresponsive.
				var cancel func()
				ctx, cancel = c.getRequestContext()
				defer cancel()
			***REMOVED***

			r, err := state.controlClient.CreateService(ctx, &swarmapi.CreateServiceRequest***REMOVED***Spec: &serviceSpec***REMOVED***)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			resp.ID = r.Service.ID
		***REMOVED***
		return nil
	***REMOVED***)

	return resp, err
***REMOVED***

// UpdateService updates existing service to match new properties.
func (c *Cluster) UpdateService(serviceIDOrName string, version uint64, spec types.ServiceSpec, flags apitypes.ServiceUpdateOptions, queryRegistry bool) (*apitypes.ServiceUpdateResponse, error) ***REMOVED***
	var resp *apitypes.ServiceUpdateResponse

	err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***

		err := c.populateNetworkID(ctx, state.controlClient, &spec)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		serviceSpec, err := convert.ServiceSpecToGRPC(spec)
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***

		currentService, err := getService(ctx, state.controlClient, serviceIDOrName, false)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		resp = &apitypes.ServiceUpdateResponse***REMOVED******REMOVED***

		switch serviceSpec.Task.Runtime.(type) ***REMOVED***
		case *swarmapi.TaskSpec_Generic:
			switch serviceSpec.Task.GetGeneric().Kind ***REMOVED***
			case string(types.RuntimePlugin):
				if spec.TaskTemplate.PluginSpec == nil ***REMOVED***
					return errors.New("plugin spec must be set")
				***REMOVED***
			***REMOVED***
		case *swarmapi.TaskSpec_Container:
			newCtnr := serviceSpec.Task.GetContainer()
			if newCtnr == nil ***REMOVED***
				return errors.New("service does not use container tasks")
			***REMOVED***

			encodedAuth := flags.EncodedRegistryAuth
			if encodedAuth != "" ***REMOVED***
				newCtnr.PullOptions = &swarmapi.ContainerSpec_PullOptions***REMOVED***RegistryAuth: encodedAuth***REMOVED***
			***REMOVED*** else ***REMOVED***
				// this is needed because if the encodedAuth isn't being updated then we
				// shouldn't lose it, and continue to use the one that was already present
				var ctnr *swarmapi.ContainerSpec
				switch flags.RegistryAuthFrom ***REMOVED***
				case apitypes.RegistryAuthFromSpec, "":
					ctnr = currentService.Spec.Task.GetContainer()
				case apitypes.RegistryAuthFromPreviousSpec:
					if currentService.PreviousSpec == nil ***REMOVED***
						return errors.New("service does not have a previous spec")
					***REMOVED***
					ctnr = currentService.PreviousSpec.Task.GetContainer()
				default:
					return errors.New("unsupported registryAuthFrom value")
				***REMOVED***
				if ctnr == nil ***REMOVED***
					return errors.New("service does not use container tasks")
				***REMOVED***
				newCtnr.PullOptions = ctnr.PullOptions
				// update encodedAuth so it can be used to pin image by digest
				if ctnr.PullOptions != nil ***REMOVED***
					encodedAuth = ctnr.PullOptions.RegistryAuth
				***REMOVED***
			***REMOVED***

			// retrieve auth config from encoded auth
			authConfig := &apitypes.AuthConfig***REMOVED******REMOVED***
			if encodedAuth != "" ***REMOVED***
				if err := json.NewDecoder(base64.NewDecoder(base64.URLEncoding, strings.NewReader(encodedAuth))).Decode(authConfig); err != nil ***REMOVED***
					logrus.Warnf("invalid authconfig: %v", err)
				***REMOVED***
			***REMOVED***

			// pin image by digest for API versions < 1.30
			// TODO(nishanttotla): The check on "DOCKER_SERVICE_PREFER_OFFLINE_IMAGE"
			// should be removed in the future. Since integration tests only use the
			// latest API version, so this is no longer required.
			if os.Getenv("DOCKER_SERVICE_PREFER_OFFLINE_IMAGE") != "1" && queryRegistry ***REMOVED***
				digestImage, err := c.imageWithDigestString(ctx, newCtnr.Image, authConfig)
				if err != nil ***REMOVED***
					logrus.Warnf("unable to pin image %s to digest: %s", newCtnr.Image, err.Error())
					// warning in the client response should be concise
					resp.Warnings = append(resp.Warnings, digestWarning(newCtnr.Image))
				***REMOVED*** else if newCtnr.Image != digestImage ***REMOVED***
					logrus.Debugf("pinning image %s by digest: %s", newCtnr.Image, digestImage)
					newCtnr.Image = digestImage
				***REMOVED*** else ***REMOVED***
					logrus.Debugf("updating service using supplied digest reference %s", newCtnr.Image)
				***REMOVED***

				// Replace the context with a fresh one.
				// If we timed out while communicating with the
				// registry, then "ctx" will already be expired, which
				// would cause UpdateService below to fail. Reusing
				// "ctx" could make it impossible to update a service
				// if the registry is slow or unresponsive.
				var cancel func()
				ctx, cancel = c.getRequestContext()
				defer cancel()
			***REMOVED***
		***REMOVED***

		var rollback swarmapi.UpdateServiceRequest_Rollback
		switch flags.Rollback ***REMOVED***
		case "", "none":
			rollback = swarmapi.UpdateServiceRequest_NONE
		case "previous":
			rollback = swarmapi.UpdateServiceRequest_PREVIOUS
		default:
			return fmt.Errorf("unrecognized rollback option %s", flags.Rollback)
		***REMOVED***

		_, err = state.controlClient.UpdateService(
			ctx,
			&swarmapi.UpdateServiceRequest***REMOVED***
				ServiceID: currentService.ID,
				Spec:      &serviceSpec,
				ServiceVersion: &swarmapi.Version***REMOVED***
					Index: version,
				***REMOVED***,
				Rollback: rollback,
			***REMOVED***,
		)
		return err
	***REMOVED***)
	return resp, err
***REMOVED***

// RemoveService removes a service from a managed swarm cluster.
func (c *Cluster) RemoveService(input string) error ***REMOVED***
	return c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		service, err := getService(ctx, state.controlClient, input, false)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		_, err = state.controlClient.RemoveService(ctx, &swarmapi.RemoveServiceRequest***REMOVED***ServiceID: service.ID***REMOVED***)
		return err
	***REMOVED***)
***REMOVED***

// ServiceLogs collects service logs and writes them back to `config.OutStream`
func (c *Cluster) ServiceLogs(ctx context.Context, selector *backend.LogSelector, config *apitypes.ContainerLogsOptions) (<-chan *backend.LogMessage, error) ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() ***REMOVED***
		return nil, c.errNoManager(state)
	***REMOVED***

	swarmSelector, err := convertSelector(ctx, state.controlClient, selector)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "error making log selector")
	***REMOVED***

	// set the streams we'll use
	stdStreams := []swarmapi.LogStream***REMOVED******REMOVED***
	if config.ShowStdout ***REMOVED***
		stdStreams = append(stdStreams, swarmapi.LogStreamStdout)
	***REMOVED***
	if config.ShowStderr ***REMOVED***
		stdStreams = append(stdStreams, swarmapi.LogStreamStderr)
	***REMOVED***

	// Get tail value squared away - the number of previous log lines we look at
	var tail int64
	// in ContainerLogs, if the tail value is ANYTHING non-integer, we just set
	// it to -1 (all). i don't agree with that, but i also think no tail value
	// should be legitimate. if you don't pass tail, we assume you want "all"
	if config.Tail == "all" || config.Tail == "" ***REMOVED***
		// tail of 0 means send all logs on the swarmkit side
		tail = 0
	***REMOVED*** else ***REMOVED***
		t, err := strconv.Atoi(config.Tail)
		if err != nil ***REMOVED***
			return nil, errors.New("tail value must be a positive integer or \"all\"")
		***REMOVED***
		if t < 0 ***REMOVED***
			return nil, errors.New("negative tail values not supported")
		***REMOVED***
		// we actually use negative tail in swarmkit to represent messages
		// backwards starting from the beginning. also, -1 means no logs. so,
		// basically, for api compat with docker container logs, add one and
		// flip the sign. we error above if you try to negative tail, which
		// isn't supported by docker (and would error deeper in the stack
		// anyway)
		//
		// See the logs protobuf for more information
		tail = int64(-(t + 1))
	***REMOVED***

	// get the since value - the time in the past we're looking at logs starting from
	var sinceProto *gogotypes.Timestamp
	if config.Since != "" ***REMOVED***
		s, n, err := timetypes.ParseTimestamps(config.Since, 0)
		if err != nil ***REMOVED***
			return nil, errors.Wrap(err, "could not parse since timestamp")
		***REMOVED***
		since := time.Unix(s, n)
		sinceProto, err = gogotypes.TimestampProto(since)
		if err != nil ***REMOVED***
			return nil, errors.Wrap(err, "could not parse timestamp to proto")
		***REMOVED***
	***REMOVED***

	stream, err := state.logsClient.SubscribeLogs(ctx, &swarmapi.SubscribeLogsRequest***REMOVED***
		Selector: swarmSelector,
		Options: &swarmapi.LogSubscriptionOptions***REMOVED***
			Follow:  config.Follow,
			Streams: stdStreams,
			Tail:    tail,
			Since:   sinceProto,
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	messageChan := make(chan *backend.LogMessage, 1)
	go func() ***REMOVED***
		defer close(messageChan)
		for ***REMOVED***
			// Check the context before doing anything.
			select ***REMOVED***
			case <-ctx.Done():
				return
			default:
			***REMOVED***
			subscribeMsg, err := stream.Recv()
			if err == io.EOF ***REMOVED***
				return
			***REMOVED***
			// if we're not io.EOF, push the message in and return
			if err != nil ***REMOVED***
				select ***REMOVED***
				case <-ctx.Done():
				case messageChan <- &backend.LogMessage***REMOVED***Err: err***REMOVED***:
				***REMOVED***
				return
			***REMOVED***

			for _, msg := range subscribeMsg.Messages ***REMOVED***
				// make a new message
				m := new(backend.LogMessage)
				m.Attrs = make([]backend.LogAttr, 0, len(msg.Attrs)+3)
				// add the timestamp, adding the error if it fails
				m.Timestamp, err = gogotypes.TimestampFromProto(msg.Timestamp)
				if err != nil ***REMOVED***
					m.Err = err
				***REMOVED***

				nodeKey := contextPrefix + ".node.id"
				serviceKey := contextPrefix + ".service.id"
				taskKey := contextPrefix + ".task.id"

				// copy over all of the details
				for _, d := range msg.Attrs ***REMOVED***
					switch d.Key ***REMOVED***
					case nodeKey, serviceKey, taskKey:
						// we have the final say over context details (in case there
						// is a conflict (if the user added a detail with a context's
						// key for some reason))
					default:
						m.Attrs = append(m.Attrs, backend.LogAttr***REMOVED***Key: d.Key, Value: d.Value***REMOVED***)
					***REMOVED***
				***REMOVED***
				m.Attrs = append(m.Attrs,
					backend.LogAttr***REMOVED***Key: nodeKey, Value: msg.Context.NodeID***REMOVED***,
					backend.LogAttr***REMOVED***Key: serviceKey, Value: msg.Context.ServiceID***REMOVED***,
					backend.LogAttr***REMOVED***Key: taskKey, Value: msg.Context.TaskID***REMOVED***,
				)

				switch msg.Stream ***REMOVED***
				case swarmapi.LogStreamStdout:
					m.Source = "stdout"
				case swarmapi.LogStreamStderr:
					m.Source = "stderr"
				***REMOVED***
				m.Line = msg.Data

				// there could be a case where the reader stops accepting
				// messages and the context is canceled. we need to check that
				// here, or otherwise we risk blocking forever on the message
				// send.
				select ***REMOVED***
				case <-ctx.Done():
					return
				case messageChan <- m:
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	return messageChan, nil
***REMOVED***

// convertSelector takes a backend.LogSelector, which contains raw names that
// may or may not be valid, and converts them to an api.LogSelector proto. It
// returns an error if something fails
func convertSelector(ctx context.Context, cc swarmapi.ControlClient, selector *backend.LogSelector) (*swarmapi.LogSelector, error) ***REMOVED***
	// don't rely on swarmkit to resolve IDs, do it ourselves
	swarmSelector := &swarmapi.LogSelector***REMOVED******REMOVED***
	for _, s := range selector.Services ***REMOVED***
		service, err := getService(ctx, cc, s, false)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		c := service.Spec.Task.GetContainer()
		if c == nil ***REMOVED***
			return nil, errors.New("logs only supported on container tasks")
		***REMOVED***
		swarmSelector.ServiceIDs = append(swarmSelector.ServiceIDs, service.ID)
	***REMOVED***
	for _, t := range selector.Tasks ***REMOVED***
		task, err := getTask(ctx, cc, t)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		c := task.Spec.GetContainer()
		if c == nil ***REMOVED***
			return nil, errors.New("logs only supported on container tasks")
		***REMOVED***
		swarmSelector.TaskIDs = append(swarmSelector.TaskIDs, task.ID)
	***REMOVED***
	return swarmSelector, nil
***REMOVED***

// imageWithDigestString takes an image such as name or name:tag
// and returns the image pinned to a digest, such as name@sha256:34234
func (c *Cluster) imageWithDigestString(ctx context.Context, image string, authConfig *apitypes.AuthConfig) (string, error) ***REMOVED***
	ref, err := reference.ParseAnyReference(image)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	namedRef, ok := ref.(reference.Named)
	if !ok ***REMOVED***
		if _, ok := ref.(reference.Digested); ok ***REMOVED***
			return image, nil
		***REMOVED***
		return "", errors.Errorf("unknown image reference format: %s", image)
	***REMOVED***
	// only query registry if not a canonical reference (i.e. with digest)
	if _, ok := namedRef.(reference.Canonical); !ok ***REMOVED***
		namedRef = reference.TagNameOnly(namedRef)

		taggedRef, ok := namedRef.(reference.NamedTagged)
		if !ok ***REMOVED***
			return "", errors.Errorf("image reference not tagged: %s", image)
		***REMOVED***

		repo, _, err := c.config.Backend.GetRepository(ctx, taggedRef, authConfig)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		dscrptr, err := repo.Tags(ctx).Get(ctx, taggedRef.Tag())
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		namedDigestedRef, err := reference.WithDigest(taggedRef, dscrptr.Digest)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		// return familiar form until interface updated to return type
		return reference.FamiliarString(namedDigestedRef), nil
	***REMOVED***
	// reference already contains a digest, so just return it
	return reference.FamiliarString(ref), nil
***REMOVED***

// digestWarning constructs a formatted warning string
// using the image name that could not be pinned by digest. The
// formatting is hardcoded, but could me made smarter in the future
func digestWarning(image string) string ***REMOVED***
	return fmt.Sprintf("image %s could not be accessed on a registry to record\nits digest. Each node will access %s independently,\npossibly leading to different nodes running different\nversions of the image.\n", image, image)
***REMOVED***
