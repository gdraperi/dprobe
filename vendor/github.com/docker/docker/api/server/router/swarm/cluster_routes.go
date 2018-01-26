package swarm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/docker/docker/api/server/httputils"
	basictypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/filters"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/errdefs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (sr *swarmRouter) initCluster(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var req types.InitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ***REMOVED***
		return err
	***REMOVED***
	nodeID, err := sr.backend.Init(req)
	if err != nil ***REMOVED***
		logrus.Errorf("Error initializing swarm: %v", err)
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, nodeID)
***REMOVED***

func (sr *swarmRouter) joinCluster(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var req types.JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ***REMOVED***
		return err
	***REMOVED***
	return sr.backend.Join(req)
***REMOVED***

func (sr *swarmRouter) leaveCluster(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	force := httputils.BoolValue(r, "force")
	return sr.backend.Leave(force)
***REMOVED***

func (sr *swarmRouter) inspectCluster(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	swarm, err := sr.backend.Inspect()
	if err != nil ***REMOVED***
		logrus.Errorf("Error getting swarm: %v", err)
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, swarm)
***REMOVED***

func (sr *swarmRouter) updateCluster(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var swarm types.Spec
	if err := json.NewDecoder(r.Body).Decode(&swarm); err != nil ***REMOVED***
		return err
	***REMOVED***

	rawVersion := r.URL.Query().Get("version")
	version, err := strconv.ParseUint(rawVersion, 10, 64)
	if err != nil ***REMOVED***
		err := fmt.Errorf("invalid swarm version '%s': %v", rawVersion, err)
		return errdefs.InvalidParameter(err)
	***REMOVED***

	var flags types.UpdateFlags

	if value := r.URL.Query().Get("rotateWorkerToken"); value != "" ***REMOVED***
		rot, err := strconv.ParseBool(value)
		if err != nil ***REMOVED***
			err := fmt.Errorf("invalid value for rotateWorkerToken: %s", value)
			return errdefs.InvalidParameter(err)
		***REMOVED***

		flags.RotateWorkerToken = rot
	***REMOVED***

	if value := r.URL.Query().Get("rotateManagerToken"); value != "" ***REMOVED***
		rot, err := strconv.ParseBool(value)
		if err != nil ***REMOVED***
			err := fmt.Errorf("invalid value for rotateManagerToken: %s", value)
			return errdefs.InvalidParameter(err)
		***REMOVED***

		flags.RotateManagerToken = rot
	***REMOVED***

	if value := r.URL.Query().Get("rotateManagerUnlockKey"); value != "" ***REMOVED***
		rot, err := strconv.ParseBool(value)
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(fmt.Errorf("invalid value for rotateManagerUnlockKey: %s", value))
		***REMOVED***

		flags.RotateManagerUnlockKey = rot
	***REMOVED***

	if err := sr.backend.Update(version, swarm, flags); err != nil ***REMOVED***
		logrus.Errorf("Error configuring swarm: %v", err)
		return err
	***REMOVED***
	return nil
***REMOVED***

func (sr *swarmRouter) unlockCluster(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var req types.UnlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := sr.backend.UnlockSwarm(req); err != nil ***REMOVED***
		logrus.Errorf("Error unlocking swarm: %v", err)
		return err
	***REMOVED***
	return nil
***REMOVED***

func (sr *swarmRouter) getUnlockKey(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	unlockKey, err := sr.backend.GetUnlockKey()
	if err != nil ***REMOVED***
		logrus.WithError(err).Errorf("Error retrieving swarm unlock key")
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, &basictypes.SwarmUnlockKeyResponse***REMOVED***
		UnlockKey: unlockKey,
	***REMOVED***)
***REMOVED***

func (sr *swarmRouter) getServices(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	filter, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***

	services, err := sr.backend.GetServices(basictypes.ServiceListOptions***REMOVED***Filters: filter***REMOVED***)
	if err != nil ***REMOVED***
		logrus.Errorf("Error getting services: %v", err)
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, services)
***REMOVED***

func (sr *swarmRouter) getService(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var insertDefaults bool
	if value := r.URL.Query().Get("insertDefaults"); value != "" ***REMOVED***
		var err error
		insertDefaults, err = strconv.ParseBool(value)
		if err != nil ***REMOVED***
			err := fmt.Errorf("invalid value for insertDefaults: %s", value)
			return errors.Wrapf(errdefs.InvalidParameter(err), "invalid value for insertDefaults: %s", value)
		***REMOVED***
	***REMOVED***

	service, err := sr.backend.GetService(vars["id"], insertDefaults)
	if err != nil ***REMOVED***
		logrus.Errorf("Error getting service %s: %v", vars["id"], err)
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, service)
***REMOVED***

func (sr *swarmRouter) createService(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var service types.ServiceSpec
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Get returns "" if the header does not exist
	encodedAuth := r.Header.Get("X-Registry-Auth")
	cliVersion := r.Header.Get("version")
	queryRegistry := false
	if cliVersion != "" && versions.LessThan(cliVersion, "1.30") ***REMOVED***
		queryRegistry = true
	***REMOVED***

	resp, err := sr.backend.CreateService(service, encodedAuth, queryRegistry)
	if err != nil ***REMOVED***
		logrus.Errorf("Error creating service %s: %v", service.Name, err)
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusCreated, resp)
***REMOVED***

func (sr *swarmRouter) updateService(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var service types.ServiceSpec
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil ***REMOVED***
		return err
	***REMOVED***

	rawVersion := r.URL.Query().Get("version")
	version, err := strconv.ParseUint(rawVersion, 10, 64)
	if err != nil ***REMOVED***
		err := fmt.Errorf("invalid service version '%s': %v", rawVersion, err)
		return errdefs.InvalidParameter(err)
	***REMOVED***

	var flags basictypes.ServiceUpdateOptions

	// Get returns "" if the header does not exist
	flags.EncodedRegistryAuth = r.Header.Get("X-Registry-Auth")
	flags.RegistryAuthFrom = r.URL.Query().Get("registryAuthFrom")
	flags.Rollback = r.URL.Query().Get("rollback")
	cliVersion := r.Header.Get("version")
	queryRegistry := false
	if cliVersion != "" && versions.LessThan(cliVersion, "1.30") ***REMOVED***
		queryRegistry = true
	***REMOVED***

	resp, err := sr.backend.UpdateService(vars["id"], version, service, flags, queryRegistry)
	if err != nil ***REMOVED***
		logrus.Errorf("Error updating service %s: %v", vars["id"], err)
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, resp)
***REMOVED***

func (sr *swarmRouter) removeService(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := sr.backend.RemoveService(vars["id"]); err != nil ***REMOVED***
		logrus.Errorf("Error removing service %s: %v", vars["id"], err)
		return err
	***REMOVED***
	return nil
***REMOVED***

func (sr *swarmRouter) getTaskLogs(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	// make a selector to pass to the helper function
	selector := &backend.LogSelector***REMOVED***
		Tasks: []string***REMOVED***vars["id"]***REMOVED***,
	***REMOVED***
	return sr.swarmLogs(ctx, w, r, selector)
***REMOVED***

func (sr *swarmRouter) getServiceLogs(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	// make a selector to pass to the helper function
	selector := &backend.LogSelector***REMOVED***
		Services: []string***REMOVED***vars["id"]***REMOVED***,
	***REMOVED***
	return sr.swarmLogs(ctx, w, r, selector)
***REMOVED***

func (sr *swarmRouter) getNodes(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	filter, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	nodes, err := sr.backend.GetNodes(basictypes.NodeListOptions***REMOVED***Filters: filter***REMOVED***)
	if err != nil ***REMOVED***
		logrus.Errorf("Error getting nodes: %v", err)
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, nodes)
***REMOVED***

func (sr *swarmRouter) getNode(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	node, err := sr.backend.GetNode(vars["id"])
	if err != nil ***REMOVED***
		logrus.Errorf("Error getting node %s: %v", vars["id"], err)
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, node)
***REMOVED***

func (sr *swarmRouter) updateNode(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var node types.NodeSpec
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil ***REMOVED***
		return err
	***REMOVED***

	rawVersion := r.URL.Query().Get("version")
	version, err := strconv.ParseUint(rawVersion, 10, 64)
	if err != nil ***REMOVED***
		err := fmt.Errorf("invalid node version '%s': %v", rawVersion, err)
		return errdefs.InvalidParameter(err)
	***REMOVED***

	if err := sr.backend.UpdateNode(vars["id"], version, node); err != nil ***REMOVED***
		logrus.Errorf("Error updating node %s: %v", vars["id"], err)
		return err
	***REMOVED***
	return nil
***REMOVED***

func (sr *swarmRouter) removeNode(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	force := httputils.BoolValue(r, "force")

	if err := sr.backend.RemoveNode(vars["id"], force); err != nil ***REMOVED***
		logrus.Errorf("Error removing node %s: %v", vars["id"], err)
		return err
	***REMOVED***
	return nil
***REMOVED***

func (sr *swarmRouter) getTasks(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	filter, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	tasks, err := sr.backend.GetTasks(basictypes.TaskListOptions***REMOVED***Filters: filter***REMOVED***)
	if err != nil ***REMOVED***
		logrus.Errorf("Error getting tasks: %v", err)
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, tasks)
***REMOVED***

func (sr *swarmRouter) getTask(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	task, err := sr.backend.GetTask(vars["id"])
	if err != nil ***REMOVED***
		logrus.Errorf("Error getting task %s: %v", vars["id"], err)
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, task)
***REMOVED***

func (sr *swarmRouter) getSecrets(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	filters, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	secrets, err := sr.backend.GetSecrets(basictypes.SecretListOptions***REMOVED***Filters: filters***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, secrets)
***REMOVED***

func (sr *swarmRouter) createSecret(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var secret types.SecretSpec
	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil ***REMOVED***
		return err
	***REMOVED***

	id, err := sr.backend.CreateSecret(secret)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusCreated, &basictypes.SecretCreateResponse***REMOVED***
		ID: id,
	***REMOVED***)
***REMOVED***

func (sr *swarmRouter) removeSecret(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := sr.backend.RemoveSecret(vars["id"]); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.WriteHeader(http.StatusNoContent)

	return nil
***REMOVED***

func (sr *swarmRouter) getSecret(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	secret, err := sr.backend.GetSecret(vars["id"])
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, secret)
***REMOVED***

func (sr *swarmRouter) updateSecret(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var secret types.SecretSpec
	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***

	rawVersion := r.URL.Query().Get("version")
	version, err := strconv.ParseUint(rawVersion, 10, 64)
	if err != nil ***REMOVED***
		return errdefs.InvalidParameter(fmt.Errorf("invalid secret version"))
	***REMOVED***

	id := vars["id"]
	return sr.backend.UpdateSecret(id, version, secret)
***REMOVED***

func (sr *swarmRouter) getConfigs(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	filters, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	configs, err := sr.backend.GetConfigs(basictypes.ConfigListOptions***REMOVED***Filters: filters***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, configs)
***REMOVED***

func (sr *swarmRouter) createConfig(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var config types.ConfigSpec
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil ***REMOVED***
		return err
	***REMOVED***

	id, err := sr.backend.CreateConfig(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusCreated, &basictypes.ConfigCreateResponse***REMOVED***
		ID: id,
	***REMOVED***)
***REMOVED***

func (sr *swarmRouter) removeConfig(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := sr.backend.RemoveConfig(vars["id"]); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.WriteHeader(http.StatusNoContent)

	return nil
***REMOVED***

func (sr *swarmRouter) getConfig(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	config, err := sr.backend.GetConfig(vars["id"])
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, config)
***REMOVED***

func (sr *swarmRouter) updateConfig(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var config types.ConfigSpec
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***

	rawVersion := r.URL.Query().Get("version")
	version, err := strconv.ParseUint(rawVersion, 10, 64)
	if err != nil ***REMOVED***
		return errdefs.InvalidParameter(fmt.Errorf("invalid config version"))
	***REMOVED***

	id := vars["id"]
	return sr.backend.UpdateConfig(id, version, config)
***REMOVED***
