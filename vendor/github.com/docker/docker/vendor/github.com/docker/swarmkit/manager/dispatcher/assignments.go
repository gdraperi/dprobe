package dispatcher

import (
	"fmt"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/equality"
	"github.com/docker/swarmkit/api/validation"
	"github.com/docker/swarmkit/manager/drivers"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/sirupsen/logrus"
)

type typeAndID struct ***REMOVED***
	id      string
	objType api.ResourceType
***REMOVED***

type assignmentSet struct ***REMOVED***
	dp                   *drivers.DriverProvider
	tasksMap             map[string]*api.Task
	tasksUsingDependency map[typeAndID]map[string]struct***REMOVED******REMOVED***
	changes              map[typeAndID]*api.AssignmentChange
	log                  *logrus.Entry
***REMOVED***

func newAssignmentSet(log *logrus.Entry, dp *drivers.DriverProvider) *assignmentSet ***REMOVED***
	return &assignmentSet***REMOVED***
		dp:                   dp,
		changes:              make(map[typeAndID]*api.AssignmentChange),
		tasksMap:             make(map[string]*api.Task),
		tasksUsingDependency: make(map[typeAndID]map[string]struct***REMOVED******REMOVED***),
		log:                  log,
	***REMOVED***
***REMOVED***

func assignSecret(a *assignmentSet, readTx store.ReadTx, mapKey typeAndID, t *api.Task) ***REMOVED***
	a.tasksUsingDependency[mapKey] = make(map[string]struct***REMOVED******REMOVED***)
	secret, err := a.secret(readTx, t, mapKey.id)
	if err != nil ***REMOVED***
		a.log.WithFields(logrus.Fields***REMOVED***
			"resource.type": "secret",
			"secret.id":     mapKey.id,
			"error":         err,
		***REMOVED***).Debug("failed to fetch secret")
		return
	***REMOVED***
	a.changes[mapKey] = &api.AssignmentChange***REMOVED***
		Assignment: &api.Assignment***REMOVED***
			Item: &api.Assignment_Secret***REMOVED***
				Secret: secret,
			***REMOVED***,
		***REMOVED***,
		Action: api.AssignmentChange_AssignmentActionUpdate,
	***REMOVED***
***REMOVED***

func assignConfig(a *assignmentSet, readTx store.ReadTx, mapKey typeAndID) ***REMOVED***
	a.tasksUsingDependency[mapKey] = make(map[string]struct***REMOVED******REMOVED***)
	config := store.GetConfig(readTx, mapKey.id)
	if config == nil ***REMOVED***
		a.log.WithFields(logrus.Fields***REMOVED***
			"resource.type": "config",
			"config.id":     mapKey.id,
		***REMOVED***).Debug("config not found")
		return
	***REMOVED***
	a.changes[mapKey] = &api.AssignmentChange***REMOVED***
		Assignment: &api.Assignment***REMOVED***
			Item: &api.Assignment_Config***REMOVED***
				Config: config,
			***REMOVED***,
		***REMOVED***,
		Action: api.AssignmentChange_AssignmentActionUpdate,
	***REMOVED***
***REMOVED***

func (a *assignmentSet) addTaskDependencies(readTx store.ReadTx, t *api.Task) ***REMOVED***
	for _, resourceRef := range t.Spec.ResourceReferences ***REMOVED***
		mapKey := typeAndID***REMOVED***objType: resourceRef.ResourceType, id: resourceRef.ResourceID***REMOVED***
		if len(a.tasksUsingDependency[mapKey]) == 0 ***REMOVED***
			switch resourceRef.ResourceType ***REMOVED***
			case api.ResourceType_SECRET:
				assignSecret(a, readTx, mapKey, t)
			case api.ResourceType_CONFIG:
				assignConfig(a, readTx, mapKey)
			default:
				a.log.WithField(
					"resource.type", resourceRef.ResourceType,
				).Debug("invalid resource type for a task dependency, skipping")
				continue
			***REMOVED***
		***REMOVED***
		a.tasksUsingDependency[mapKey][t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	var secrets []*api.SecretReference
	container := t.Spec.GetContainer()
	if container != nil ***REMOVED***
		secrets = container.Secrets
	***REMOVED***

	for _, secretRef := range secrets ***REMOVED***
		secretID := secretRef.SecretID
		mapKey := typeAndID***REMOVED***objType: api.ResourceType_SECRET, id: secretID***REMOVED***

		if len(a.tasksUsingDependency[mapKey]) == 0 ***REMOVED***
			assignSecret(a, readTx, mapKey, t)
		***REMOVED***
		a.tasksUsingDependency[mapKey][t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	var configs []*api.ConfigReference
	if container != nil ***REMOVED***
		configs = container.Configs
	***REMOVED***
	for _, configRef := range configs ***REMOVED***
		configID := configRef.ConfigID
		mapKey := typeAndID***REMOVED***objType: api.ResourceType_CONFIG, id: configID***REMOVED***

		if len(a.tasksUsingDependency[mapKey]) == 0 ***REMOVED***
			assignConfig(a, readTx, mapKey)
		***REMOVED***
		a.tasksUsingDependency[mapKey][t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

func (a *assignmentSet) releaseDependency(mapKey typeAndID, assignment *api.Assignment, taskID string) bool ***REMOVED***
	delete(a.tasksUsingDependency[mapKey], taskID)
	if len(a.tasksUsingDependency[mapKey]) != 0 ***REMOVED***
		return false
	***REMOVED***
	// No tasks are using the dependency anymore
	delete(a.tasksUsingDependency, mapKey)
	a.changes[mapKey] = &api.AssignmentChange***REMOVED***
		Assignment: assignment,
		Action:     api.AssignmentChange_AssignmentActionRemove,
	***REMOVED***
	return true
***REMOVED***

func (a *assignmentSet) releaseTaskDependencies(t *api.Task) bool ***REMOVED***
	var modified bool

	for _, resourceRef := range t.Spec.ResourceReferences ***REMOVED***
		var assignment *api.Assignment
		switch resourceRef.ResourceType ***REMOVED***
		case api.ResourceType_SECRET:
			assignment = &api.Assignment***REMOVED***
				Item: &api.Assignment_Secret***REMOVED***
					Secret: &api.Secret***REMOVED***ID: resourceRef.ResourceID***REMOVED***,
				***REMOVED***,
			***REMOVED***
		case api.ResourceType_CONFIG:
			assignment = &api.Assignment***REMOVED***
				Item: &api.Assignment_Config***REMOVED***
					Config: &api.Config***REMOVED***ID: resourceRef.ResourceID***REMOVED***,
				***REMOVED***,
			***REMOVED***
		default:
			a.log.WithField(
				"resource.type", resourceRef.ResourceType,
			).Debug("invalid resource type for a task dependency, skipping")
			continue
		***REMOVED***

		mapKey := typeAndID***REMOVED***objType: resourceRef.ResourceType, id: resourceRef.ResourceID***REMOVED***
		if a.releaseDependency(mapKey, assignment, t.ID) ***REMOVED***
			modified = true
		***REMOVED***
	***REMOVED***

	container := t.Spec.GetContainer()

	var secrets []*api.SecretReference
	if container != nil ***REMOVED***
		secrets = container.Secrets
	***REMOVED***

	for _, secretRef := range secrets ***REMOVED***
		secretID := secretRef.SecretID
		mapKey := typeAndID***REMOVED***objType: api.ResourceType_SECRET, id: secretID***REMOVED***
		assignment := &api.Assignment***REMOVED***
			Item: &api.Assignment_Secret***REMOVED***
				Secret: &api.Secret***REMOVED***ID: secretID***REMOVED***,
			***REMOVED***,
		***REMOVED***
		if a.releaseDependency(mapKey, assignment, t.ID) ***REMOVED***
			modified = true
		***REMOVED***
	***REMOVED***

	var configs []*api.ConfigReference
	if container != nil ***REMOVED***
		configs = container.Configs
	***REMOVED***

	for _, configRef := range configs ***REMOVED***
		configID := configRef.ConfigID
		mapKey := typeAndID***REMOVED***objType: api.ResourceType_CONFIG, id: configID***REMOVED***
		assignment := &api.Assignment***REMOVED***
			Item: &api.Assignment_Config***REMOVED***
				Config: &api.Config***REMOVED***ID: configID***REMOVED***,
			***REMOVED***,
		***REMOVED***
		if a.releaseDependency(mapKey, assignment, t.ID) ***REMOVED***
			modified = true
		***REMOVED***
	***REMOVED***

	return modified
***REMOVED***

func (a *assignmentSet) addOrUpdateTask(readTx store.ReadTx, t *api.Task) bool ***REMOVED***
	// We only care about tasks that are ASSIGNED or higher.
	if t.Status.State < api.TaskStateAssigned ***REMOVED***
		return false
	***REMOVED***

	if oldTask, exists := a.tasksMap[t.ID]; exists ***REMOVED***
		// States ASSIGNED and below are set by the orchestrator/scheduler,
		// not the agent, so tasks in these states need to be sent to the
		// agent even if nothing else has changed.
		if equality.TasksEqualStable(oldTask, t) && t.Status.State > api.TaskStateAssigned ***REMOVED***
			// this update should not trigger a task change for the agent
			a.tasksMap[t.ID] = t
			// If this task got updated to a final state, let's release
			// the dependencies that are being used by the task
			if t.Status.State > api.TaskStateRunning ***REMOVED***
				// If releasing the dependencies caused us to
				// remove something from the assignment set,
				// mark one modification.
				return a.releaseTaskDependencies(t)
			***REMOVED***
			return false
		***REMOVED***
	***REMOVED*** else if t.Status.State <= api.TaskStateRunning ***REMOVED***
		// If this task wasn't part of the assignment set before, and it's <= RUNNING
		// add the dependencies it references to the assignment.
		// Task states > RUNNING are worker reported only, are never created in
		// a > RUNNING state.
		a.addTaskDependencies(readTx, t)
	***REMOVED***
	a.tasksMap[t.ID] = t
	a.changes[typeAndID***REMOVED***objType: api.ResourceType_TASK, id: t.ID***REMOVED***] = &api.AssignmentChange***REMOVED***
		Assignment: &api.Assignment***REMOVED***
			Item: &api.Assignment_Task***REMOVED***
				Task: t,
			***REMOVED***,
		***REMOVED***,
		Action: api.AssignmentChange_AssignmentActionUpdate,
	***REMOVED***
	return true
***REMOVED***

func (a *assignmentSet) removeTask(t *api.Task) bool ***REMOVED***
	if _, exists := a.tasksMap[t.ID]; !exists ***REMOVED***
		return false
	***REMOVED***

	a.changes[typeAndID***REMOVED***objType: api.ResourceType_TASK, id: t.ID***REMOVED***] = &api.AssignmentChange***REMOVED***
		Assignment: &api.Assignment***REMOVED***
			Item: &api.Assignment_Task***REMOVED***
				Task: &api.Task***REMOVED***ID: t.ID***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Action: api.AssignmentChange_AssignmentActionRemove,
	***REMOVED***

	delete(a.tasksMap, t.ID)

	// Release the dependencies being used by this task.
	// Ignoring the return here. We will always mark this as a
	// modification, since a task is being removed.
	a.releaseTaskDependencies(t)
	return true
***REMOVED***

func (a *assignmentSet) message() api.AssignmentsMessage ***REMOVED***
	var message api.AssignmentsMessage
	for _, change := range a.changes ***REMOVED***
		message.Changes = append(message.Changes, change)
	***REMOVED***

	// The the set of changes is reinitialized to prepare for formation
	// of the next message.
	a.changes = make(map[typeAndID]*api.AssignmentChange)

	return message
***REMOVED***

// secret populates the secret value from raft store. For external secrets, the value is populated
// from the secret driver.
func (a *assignmentSet) secret(readTx store.ReadTx, task *api.Task, secretID string) (*api.Secret, error) ***REMOVED***
	secret := store.GetSecret(readTx, secretID)
	if secret == nil ***REMOVED***
		return nil, fmt.Errorf("secret not found")
	***REMOVED***
	if secret.Spec.Driver == nil ***REMOVED***
		return secret, nil
	***REMOVED***
	d, err := a.dp.NewSecretDriver(secret.Spec.Driver)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	value, err := d.Get(&secret.Spec, task)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := validation.ValidateSecretPayload(value); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Assign the secret
	secret.Spec.Data = value
	return secret, nil
***REMOVED***
