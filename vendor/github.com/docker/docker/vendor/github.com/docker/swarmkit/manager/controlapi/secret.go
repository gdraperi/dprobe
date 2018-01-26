package controlapi

import (
	"crypto/subtle"
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/validation"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// assumes spec is not nil
func secretFromSecretSpec(spec *api.SecretSpec) *api.Secret ***REMOVED***
	return &api.Secret***REMOVED***
		ID:   identity.NewID(),
		Spec: *spec,
	***REMOVED***
***REMOVED***

// GetSecret returns a `GetSecretResponse` with a `Secret` with the same
// id as `GetSecretRequest.SecretID`
// - Returns `NotFound` if the Secret with the given id is not found.
// - Returns `InvalidArgument` if the `GetSecretRequest.SecretID` is empty.
// - Returns an error if getting fails.
func (s *Server) GetSecret(ctx context.Context, request *api.GetSecretRequest) (*api.GetSecretResponse, error) ***REMOVED***
	if request.SecretID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "secret ID must be provided")
	***REMOVED***

	var secret *api.Secret
	s.store.View(func(tx store.ReadTx) ***REMOVED***
		secret = store.GetSecret(tx, request.SecretID)
	***REMOVED***)

	if secret == nil ***REMOVED***
		return nil, status.Errorf(codes.NotFound, "secret %s not found", request.SecretID)
	***REMOVED***

	secret.Spec.Data = nil // clean the actual secret data so it's never returned
	return &api.GetSecretResponse***REMOVED***Secret: secret***REMOVED***, nil
***REMOVED***

// UpdateSecret updates a Secret referenced by SecretID with the given SecretSpec.
// - Returns `NotFound` if the Secret is not found.
// - Returns `InvalidArgument` if the SecretSpec is malformed or anything other than Labels is changed
// - Returns an error if the update fails.
func (s *Server) UpdateSecret(ctx context.Context, request *api.UpdateSecretRequest) (*api.UpdateSecretResponse, error) ***REMOVED***
	if request.SecretID == "" || request.SecretVersion == nil ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***
	var secret *api.Secret
	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		secret = store.GetSecret(tx, request.SecretID)
		if secret == nil ***REMOVED***
			return status.Errorf(codes.NotFound, "secret %s not found", request.SecretID)
		***REMOVED***

		// Check if the Name is different than the current name, or the secret is non-nil and different
		// than the current secret
		if secret.Spec.Annotations.Name != request.Spec.Annotations.Name ||
			(request.Spec.Data != nil && subtle.ConstantTimeCompare(request.Spec.Data, secret.Spec.Data) == 0) ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "only updates to Labels are allowed")
		***REMOVED***

		// We only allow updating Labels
		secret.Meta.Version = *request.SecretVersion
		secret.Spec.Annotations.Labels = request.Spec.Annotations.Labels

		return store.UpdateSecret(tx, secret)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	log.G(ctx).WithFields(logrus.Fields***REMOVED***
		"secret.ID":   request.SecretID,
		"secret.Name": request.Spec.Annotations.Name,
		"method":      "UpdateSecret",
	***REMOVED***).Debugf("secret updated")

	// WARN: we should never return the actual secret data here. We need to redact the private fields first.
	secret.Spec.Data = nil
	return &api.UpdateSecretResponse***REMOVED***
		Secret: secret,
	***REMOVED***, nil
***REMOVED***

// ListSecrets returns a `ListSecretResponse` with a list all non-internal `Secret`s being
// managed, or all secrets matching any name in `ListSecretsRequest.Names`, any
// name prefix in `ListSecretsRequest.NamePrefixes`, any id in
// `ListSecretsRequest.SecretIDs`, or any id prefix in `ListSecretsRequest.IDPrefixes`.
// - Returns an error if listing fails.
func (s *Server) ListSecrets(ctx context.Context, request *api.ListSecretsRequest) (*api.ListSecretsResponse, error) ***REMOVED***
	var (
		secrets     []*api.Secret
		respSecrets []*api.Secret
		err         error
		byFilters   []store.By
		by          store.By
		labels      map[string]string
	)

	// return all secrets that match either any of the names or any of the name prefixes (why would you give both?)
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
		secrets, err = store.FindSecrets(tx, by)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// strip secret data from the secret, filter by label, and filter out all internal secrets
	for _, secret := range secrets ***REMOVED***
		if secret.Internal || !filterMatchLabels(secret.Spec.Annotations.Labels, labels) ***REMOVED***
			continue
		***REMOVED***
		secret.Spec.Data = nil // clean the actual secret data so it's never returned
		respSecrets = append(respSecrets, secret)
	***REMOVED***

	return &api.ListSecretsResponse***REMOVED***Secrets: respSecrets***REMOVED***, nil
***REMOVED***

// CreateSecret creates and returns a `CreateSecretResponse` with a `Secret` based
// on the provided `CreateSecretRequest.SecretSpec`.
// - Returns `InvalidArgument` if the `CreateSecretRequest.SecretSpec` is malformed,
//   or if the secret data is too long or contains invalid characters.
// - Returns an error if the creation fails.
func (s *Server) CreateSecret(ctx context.Context, request *api.CreateSecretRequest) (*api.CreateSecretResponse, error) ***REMOVED***
	if err := validateSecretSpec(request.Spec); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if request.Spec.Driver != nil ***REMOVED*** // Check that the requested driver is valid
		if _, err := s.dr.NewSecretDriver(request.Spec.Driver); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	secret := secretFromSecretSpec(request.Spec) // the store will handle name conflicts
	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		return store.CreateSecret(tx, secret)
	***REMOVED***)

	switch err ***REMOVED***
	case store.ErrNameConflict:
		return nil, status.Errorf(codes.AlreadyExists, "secret %s already exists", request.Spec.Annotations.Name)
	case nil:
		secret.Spec.Data = nil // clean the actual secret data so it's never returned
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"secret.Name": request.Spec.Annotations.Name,
			"method":      "CreateSecret",
		***REMOVED***).Debugf("secret created")

		return &api.CreateSecretResponse***REMOVED***Secret: secret***REMOVED***, nil
	default:
		return nil, err
	***REMOVED***
***REMOVED***

// RemoveSecret removes the secret referenced by `RemoveSecretRequest.ID`.
// - Returns `InvalidArgument` if `RemoveSecretRequest.ID` is empty.
// - Returns `NotFound` if the a secret named `RemoveSecretRequest.ID` is not found.
// - Returns `SecretInUse` if the secret is currently in use
// - Returns an error if the deletion fails.
func (s *Server) RemoveSecret(ctx context.Context, request *api.RemoveSecretRequest) (*api.RemoveSecretResponse, error) ***REMOVED***
	if request.SecretID == "" ***REMOVED***
		return nil, status.Errorf(codes.InvalidArgument, "secret ID must be provided")
	***REMOVED***

	err := s.store.Update(func(tx store.Tx) error ***REMOVED***
		// Check if the secret exists
		secret := store.GetSecret(tx, request.SecretID)
		if secret == nil ***REMOVED***
			return status.Errorf(codes.NotFound, "could not find secret %s", request.SecretID)
		***REMOVED***

		// Check if any services currently reference this secret, return error if so
		services, err := store.FindServices(tx, store.ByReferencedSecretID(request.SecretID))
		if err != nil ***REMOVED***
			return status.Errorf(codes.Internal, "could not find services using secret %s: %v", request.SecretID, err)
		***REMOVED***

		if len(services) != 0 ***REMOVED***
			serviceNames := make([]string, 0, len(services))
			for _, service := range services ***REMOVED***
				serviceNames = append(serviceNames, service.Spec.Annotations.Name)
			***REMOVED***

			secretName := secret.Spec.Annotations.Name
			serviceNameStr := strings.Join(serviceNames, ", ")
			serviceStr := "services"
			if len(serviceNames) == 1 ***REMOVED***
				serviceStr = "service"
			***REMOVED***

			return status.Errorf(codes.InvalidArgument, "secret '%s' is in use by the following %s: %v", secretName, serviceStr, serviceNameStr)
		***REMOVED***

		return store.DeleteSecret(tx, request.SecretID)
	***REMOVED***)
	switch err ***REMOVED***
	case store.ErrNotExist:
		return nil, status.Errorf(codes.NotFound, "secret %s not found", request.SecretID)
	case nil:
		log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"secret.ID": request.SecretID,
			"method":    "RemoveSecret",
		***REMOVED***).Debugf("secret removed")

		return &api.RemoveSecretResponse***REMOVED******REMOVED***, nil
	default:
		return nil, err
	***REMOVED***
***REMOVED***

func validateSecretSpec(spec *api.SecretSpec) error ***REMOVED***
	if spec == nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	***REMOVED***
	if err := validateConfigOrSecretAnnotations(spec.Annotations); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Check if secret driver is defined
	if spec.Driver != nil ***REMOVED***
		// Ensure secret driver has a name
		if spec.Driver.Name == "" ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "secret driver must have a name")
		***REMOVED***
		return nil
	***REMOVED***
	if err := validation.ValidateSecretPayload(spec.Data); err != nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	***REMOVED***
	return nil
***REMOVED***
