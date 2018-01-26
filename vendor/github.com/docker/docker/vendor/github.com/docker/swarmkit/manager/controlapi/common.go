package controlapi

import (
	"regexp"
	"strings"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/allocator"
	"github.com/docker/swarmkit/manager/state/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var isValidDNSName = regexp.MustCompile(`^[a-zA-Z0-9](?:[-_]*[A-Za-z0-9]+)*$`)

// configs and secrets have different naming requirements from tasks and services
var isValidConfigOrSecretName = regexp.MustCompile(`^[a-zA-Z0-9]+(?:[a-zA-Z0-9-_.]*[a-zA-Z0-9])?$`)

func buildFilters(by func(string) store.By, values []string) store.By ***REMOVED***
	filters := make([]store.By, 0, len(values))
	for _, v := range values ***REMOVED***
		filters = append(filters, by(v))
	***REMOVED***
	return store.Or(filters...)
***REMOVED***

func filterContains(match string, candidates []string) bool ***REMOVED***
	if len(candidates) == 0 ***REMOVED***
		return true
	***REMOVED***
	for _, c := range candidates ***REMOVED***
		if c == match ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func filterContainsPrefix(match string, candidates []string) bool ***REMOVED***
	if len(candidates) == 0 ***REMOVED***
		return true
	***REMOVED***
	for _, c := range candidates ***REMOVED***
		if strings.HasPrefix(match, c) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func filterMatchLabels(match map[string]string, candidates map[string]string) bool ***REMOVED***
	if len(candidates) == 0 ***REMOVED***
		return true
	***REMOVED***

	for k, v := range candidates ***REMOVED***
		c, ok := match[k]
		if !ok ***REMOVED***
			return false
		***REMOVED***
		if v != "" && v != c ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func validateAnnotations(m api.Annotations) error ***REMOVED***
	if m.Name == "" ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "meta: name must be provided")
	***REMOVED***
	if !isValidDNSName.MatchString(m.Name) ***REMOVED***
		// if the name doesn't match the regex
		return status.Errorf(codes.InvalidArgument, "name must be valid as a DNS name component")
	***REMOVED***
	if len(m.Name) > 63 ***REMOVED***
		// DNS labels are limited to 63 characters
		return status.Errorf(codes.InvalidArgument, "name must be 63 characters or fewer")
	***REMOVED***
	return nil
***REMOVED***

func validateConfigOrSecretAnnotations(m api.Annotations) error ***REMOVED***
	if m.Name == "" ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "name must be provided")
	***REMOVED*** else if len(m.Name) > 64 || !isValidConfigOrSecretName.MatchString(m.Name) ***REMOVED***
		// if the name doesn't match the regex
		return status.Errorf(codes.InvalidArgument,
			"invalid name, only 64 [a-zA-Z0-9-_.] characters allowed, and the start and end character must be [a-zA-Z0-9]")
	***REMOVED***
	return nil
***REMOVED***

func validateDriver(driver *api.Driver, pg plugingetter.PluginGetter, pluginType string) error ***REMOVED***
	if driver == nil ***REMOVED***
		// It is ok to not specify the driver. We will choose
		// a default driver.
		return nil
	***REMOVED***

	if driver.Name == "" ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "driver name: if driver is specified name is required")
	***REMOVED***

	// First check against the known drivers
	switch pluginType ***REMOVED***
	case ipamapi.PluginEndpointType:
		if strings.ToLower(driver.Name) == ipamapi.DefaultIPAM ***REMOVED***
			return nil
		***REMOVED***
	case driverapi.NetworkPluginEndpointType:
		if allocator.IsBuiltInNetworkDriver(driver.Name) ***REMOVED***
			return nil
		***REMOVED***
	default:
	***REMOVED***

	if pg == nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "plugin %s not supported", driver.Name)
	***REMOVED***

	p, err := pg.Get(driver.Name, pluginType, plugingetter.Lookup)
	if err != nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "error during lookup of plugin %s", driver.Name)
	***REMOVED***

	if p.IsV1() ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "legacy plugin %s of type %s is not supported in swarm mode", driver.Name, pluginType)
	***REMOVED***

	return nil
***REMOVED***
