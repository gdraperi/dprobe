package controlapi

import (
	"bytes"
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MaxConfigSize is the maximum byte length of the `Config.Spec.Data` field.
const MaxConfigSize = 500 * 1024 // 500KB

// assumes spec is not nil
func configFromConfigSpec(spec *api.ConfigSpec) *api.Config ***REMOVED***
	return &api.Config***REMOVED***
		ID:   identity.NewID(),
		Spec: *spec,
	***REMOVED***
***REMOVED***

// GetConfig returns a `GetConfigResponse` with a `Config` with the same
// id as `GetConfigRequest.ConfigID`
// - Returns `NotFound` if the Config with the given id is not found.
// - Returns `InvalidArgument` if the `GetConfigRequest.ConfigID` is empty.
// - Returns an error if getting fails.
func (s *Server) GetConfig(ctx context.Context, request *api.GetConfigRequest) (*api.GetConfigResponse, error) ***REMOVED***
	if request.ConfigID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "config ID must be provided")
	***REMOVED***

	var config *api.Config
	s.store.View(func(tx store.ReadTx) ***REMOVED***
		config = store.GetConfig(tx, request.ConfigID)
	***REMOVED***)

	if config == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "config %s not found", request.ConfigID)
	***REMOVED***

	return &api.GetConfigResponse***REMOVED***Config: config***REMOVED***, nil
***REMOVED***

// UpdateConfig updates a Config referenced by ConfigID with the given ConfigSpec.
// - Returns `NotFound` if the Config is not found.
// - Returns `InvalidArgument` if the ConfigSpec is malformed or anything other than Labels is changed
// - Returns an error if the update fails.
func (s *Server) UpdateConfig(ctx context.Context, request *api.UpdateConfigRequest) (*api.UpdateConfigResponse, error) ***REMOVED***
	if request.ConfigID == "" || request.ConfigVersion == nil ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***

	var config *api.Config
	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		config = store.GetConfig(tx, request.ConfigID)
		if config == nil ***REMOVED***
			return status.Errorf(codes.NotFound, "config %s not found", request.ConfigID)
		***REMOVED***

		// Check if the Name is different than the current name, or the config is non-nil and different
		// than the current config
		if config.Spec.Annotations.Name != request.Spec.Annotations.Name ||
			(request.Spec.Data != nil && !bytes.Equal(request.Spec.Data, config.Spec.Data)) ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "only updates to Labels are allowed")
		***REMOVED***

		// We only allow updating Labels
		config.Meta.Version = *request.ConfigVersion
		config.Spec.Annotations.Labels = request.Spec.Annotations.Labels

		return store.UpdateConfig(tx, config)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"config.ID":   request.ConfigID,
		"config.Name": request.Spec.Annotations.Name,
		"method":      "UpdateConfig",
	***REMOVED***).Debugf("config updated")

	return &api.UpdateConfigResponse***REMOVED***
		Config: config,
	***REMOVED***, nil
***REMOVED***

// ListConfigs returns a `ListConfigResponse` with a list all non-internal `Config`s being
// managed, or all configs matching any name in `ListConfigsRequest.Names`, any
// name prefix in `ListConfigsRequest.NamePrefixes`, any id in
// `ListConfigsRequest.ConfigIDs`, or any id prefix in `ListConfigsRequest.IDPrefixes`.
// - Returns an error if listing fails.
func (s *Server) ListConfigs(ctx context.Context, request *api.ListConfigsRequest) (*api.ListConfigsResponse, error) ***REMOVED***
	var (
		configs     []*api.Config
		respConfigs []*api.Config
		err         error
		byFilters   []store.By
		by          store.By
		labels      map[string]string
	)

	// return all configs that match either any of the names or any of the name prefixes (why would you give both?)
	if request.Filters != nil ***REMOVED***
		for _, name := range request.Filters.Names ***REMOVED***
			byFilters = append(byFilters, store.ByName(name))
		***REMOVED***
		for _, prefix := range request.Filters.NamePrefixes ***REMOVED***
			byFilters = append(byFilters, store.ByNamePrefix(prefix))
		***REMOVED***
		for _, prefix := range request.Filters.IDPrefixes ***REMOVED***
			byFilters = append(byFilters, store.ByIDPrefix(prefix))
		***REMOVED***
		labels = request.Filters.Labels
	***REMOVED***

	switch len(byFilters) ***REMOVED***
	case 0:
		by = store.All
	case 1:
		by = byFilters[0]
	default:
		by = store.Or(byFilters...)
	***REMOVED***

	s.store.View(func(tx store.ReadTx) ***REMOVED***
		configs, err = store.FindConfigs(tx, by)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// filter by label
	for _, config := range configs ***REMOVED***
		if !filterMatchLabels(config.Spec.Annotations.Labels, labels) ***REMOVED***
			continue
		***REMOVED***
		respConfigs = append(respConfigs, config)
	***REMOVED***

	return &api.ListConfigsResponse***REMOVED***Configs: respConfigs***REMOVED***, nil
***REMOVED***

// CreateConfig creates and returns a `CreateConfigResponse` with a `Config` based
// on the provided `CreateConfigRequest.ConfigSpec`.
// - Returns `InvalidArgument` if the `CreateConfigRequest.ConfigSpec` is malformed,
//   or if the config data is too long or contains invalid characters.
// - Returns an error if the creation fails.
func (s *Server) CreateConfig(ctx context.Context, request *api.CreateConfigRequest) (*api.CreateConfigResponse, error) ***REMOVED***
	if err := validateConfigSpec(request.Spec); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	config := configFromConfigSpec(request.Spec) // the store will handle name conflicts
	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		return store.CreateConfig(tx, config)
	***REMOVED***)

	switch err ***REMOVED***
	case store.ErrNameConflict:
		return nil, status.Errorf(codes.AlreadyExists, "config %s already exists", request.Spec.Annotations.Name)
	case nil:
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"config.Name": request.Spec.Annotations.Name,
			"method":      "CreateConfig",
		***REMOVED***).Debugf("config created")

		return &api.CreateConfigResponse***REMOVED***Config: config***REMOVED***, nil
	default:
		return nil, err
	***REMOVED***
***REMOVED***

// RemoveConfig removes the config referenced by `RemoveConfigRequest.ID`.
// - Returns `InvalidArgument` if `RemoveConfigRequest.ID` is empty.
// - Returns `NotFound` if the a config named `RemoveConfigRequest.ID` is not found.
// - Returns `ConfigInUse` if the config is currently in use
// - Returns an error if the deletion fails.
func (s *Server) RemoveConfig(ctx context.Context, request *api.RemoveConfigRequest) (*api.RemoveConfigResponse, error) ***REMOVED***
	if request.ConfigID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "config ID must be provided")
	***REMOVED***

	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		// Check if the config exists
		config := store.GetConfig(tx, request.ConfigID)
		if config == nil ***REMOVED***
			return status.Errorf(codes.NotFound, "could not find config %s", request.ConfigID)
		***REMOVED***

		// Check if any services currently reference this config, return error if so
		services, err := store.FindServices(tx, store.ByReferencedConfigID(request.ConfigID))
		if err != nil ***REMOVED***
			return status.Errorf(codes.Internal, "could not find services using config %s: %v", request.ConfigID, err)
		***REMOVED***

		if len(services) != 0 ***REMOVED***
			serviceNames := make([]string, 0, len(services))
			for _, service := range services ***REMOVED***
				serviceNames = append(serviceNames, service.Spec.Annotations.Name)
			***REMOVED***

			configName := config.Spec.Annotations.Name
			serviceNameStr := strings.Join(serviceNames, ", ")
			serviceStr := "services"
			if len(serviceNames) == 1 ***REMOVED***
				serviceStr = "service"
			***REMOVED***

			return status.Errorf(codes.InvalidArgument, "config '%s' is in use by the following %s: %v", configName, serviceStr, serviceNameStr)
		***REMOVED***

		return store.DeleteConfig(tx, request.ConfigID)
	***REMOVED***)
	switch err ***REMOVED***
	case store.ErrNotExist:
		return nil, status.Errorf(codes.NotFound, "config %s not found", request.ConfigID)
	case nil:
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"config.ID": request.ConfigID,
			"method":    "RemoveConfig",
		***REMOVED***).Debugf("config removed")

		return &api.RemoveConfigResponse***REMOVED******REMOVED***, nil
	default:
		return nil, err
	***REMOVED***
***REMOVED***

func validateConfigSpec(spec *api.ConfigSpec) error ***REMOVED***
	if spec == nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***
	if err := validateConfigOrSecretAnnotations(spec.Annotations); err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(spec.Data) >= MaxConfigSize || len(spec.Data) < 1 ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "config data must be larger than 0 and less than %d bytes", MaxConfigSize)
	***REMOVED***
	return nil
***REMOVED***
