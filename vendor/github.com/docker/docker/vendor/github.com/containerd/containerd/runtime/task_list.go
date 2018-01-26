package runtime

import (
	"context"
	"sync"

	"github.com/containerd/containerd/namespaces"
	"github.com/pkg/errors"
)

var (
	// ErrTaskNotExists is returned when a task does not exist
	ErrTaskNotExists = errors.New("task does not exist")
	// ErrTaskAlreadyExists is returned when a task already exists
	ErrTaskAlreadyExists = errors.New("task already exists")
)

// NewTaskList returns a new TaskList
func NewTaskList() *TaskList ***REMOVED***
	return &TaskList***REMOVED***
		tasks: make(map[string]map[string]Task),
	***REMOVED***
***REMOVED***

// TaskList holds and provides locking around tasks
type TaskList struct ***REMOVED***
	mu    sync.Mutex
	tasks map[string]map[string]Task
***REMOVED***

// Get a task
func (l *TaskList) Get(ctx context.Context, id string) (Task, error) ***REMOVED***
	l.mu.Lock()
	defer l.mu.Unlock()
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	tasks, ok := l.tasks[namespace]
	if !ok ***REMOVED***
		return nil, ErrTaskNotExists
	***REMOVED***
	t, ok := tasks[id]
	if !ok ***REMOVED***
		return nil, ErrTaskNotExists
	***REMOVED***
	return t, nil
***REMOVED***

// GetAll tasks under a namespace
func (l *TaskList) GetAll(ctx context.Context) ([]Task, error) ***REMOVED***
	l.mu.Lock()
	defer l.mu.Unlock()
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var o []Task
	tasks, ok := l.tasks[namespace]
	if !ok ***REMOVED***
		return o, nil
	***REMOVED***
	for _, t := range tasks ***REMOVED***
		o = append(o, t)
	***REMOVED***
	return o, nil
***REMOVED***

// Add a task
func (l *TaskList) Add(ctx context.Context, t Task) error ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return l.AddWithNamespace(namespace, t)
***REMOVED***

// AddWithNamespace adds a task with the provided namespace
func (l *TaskList) AddWithNamespace(namespace string, t Task) error ***REMOVED***
	l.mu.Lock()
	defer l.mu.Unlock()

	id := t.ID()
	if _, ok := l.tasks[namespace]; !ok ***REMOVED***
		l.tasks[namespace] = make(map[string]Task)
	***REMOVED***
	if _, ok := l.tasks[namespace][id]; ok ***REMOVED***
		return errors.Wrap(ErrTaskAlreadyExists, id)
	***REMOVED***
	l.tasks[namespace][id] = t
	return nil
***REMOVED***

// Delete a task
func (l *TaskList) Delete(ctx context.Context, t Task) ***REMOVED***
	l.mu.Lock()
	defer l.mu.Unlock()
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	tasks, ok := l.tasks[namespace]
	if ok ***REMOVED***
		delete(tasks, t.ID())
	***REMOVED***
***REMOVED***
