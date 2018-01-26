package runtime

// TaskMonitor provides an interface for monitoring of containers within containerd
type TaskMonitor interface ***REMOVED***
	// Monitor adds the provided container to the monitor
	Monitor(Task) error
	// Stop stops and removes the provided container from the monitor
	Stop(Task) error
***REMOVED***

// NewMultiTaskMonitor returns a new TaskMonitor broadcasting to the provided monitors
func NewMultiTaskMonitor(monitors ...TaskMonitor) TaskMonitor ***REMOVED***
	return &multiTaskMonitor***REMOVED***
		monitors: monitors,
	***REMOVED***
***REMOVED***

// NewNoopMonitor is a task monitor that does nothing
func NewNoopMonitor() TaskMonitor ***REMOVED***
	return &noopTaskMonitor***REMOVED******REMOVED***
***REMOVED***

type noopTaskMonitor struct ***REMOVED***
***REMOVED***

func (mm *noopTaskMonitor) Monitor(c Task) error ***REMOVED***
	return nil
***REMOVED***

func (mm *noopTaskMonitor) Stop(c Task) error ***REMOVED***
	return nil
***REMOVED***

type multiTaskMonitor struct ***REMOVED***
	monitors []TaskMonitor
***REMOVED***

func (mm *multiTaskMonitor) Monitor(c Task) error ***REMOVED***
	for _, m := range mm.monitors ***REMOVED***
		if err := m.Monitor(c); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (mm *multiTaskMonitor) Stop(c Task) error ***REMOVED***
	for _, m := range mm.monitors ***REMOVED***
		if err := m.Stop(c); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
