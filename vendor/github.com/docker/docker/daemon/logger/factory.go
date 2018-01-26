package logger

import (
	"fmt"
	"sort"
	"sync"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/plugingetter"
	units "github.com/docker/go-units"
	"github.com/pkg/errors"
)

// Creator builds a logging driver instance with given context.
type Creator func(Info) (Logger, error)

// LogOptValidator checks the options specific to the underlying
// logging implementation.
type LogOptValidator func(cfg map[string]string) error

type logdriverFactory struct ***REMOVED***
	registry     map[string]Creator
	optValidator map[string]LogOptValidator
	m            sync.Mutex
***REMOVED***

func (lf *logdriverFactory) list() []string ***REMOVED***
	ls := make([]string, 0, len(lf.registry))
	lf.m.Lock()
	for name := range lf.registry ***REMOVED***
		ls = append(ls, name)
	***REMOVED***
	lf.m.Unlock()
	sort.Strings(ls)
	return ls
***REMOVED***

// ListDrivers gets the list of registered log driver names
func ListDrivers() []string ***REMOVED***
	return factory.list()
***REMOVED***

func (lf *logdriverFactory) register(name string, c Creator) error ***REMOVED***
	if lf.driverRegistered(name) ***REMOVED***
		return fmt.Errorf("logger: log driver named '%s' is already registered", name)
	***REMOVED***

	lf.m.Lock()
	lf.registry[name] = c
	lf.m.Unlock()
	return nil
***REMOVED***

func (lf *logdriverFactory) driverRegistered(name string) bool ***REMOVED***
	lf.m.Lock()
	_, ok := lf.registry[name]
	lf.m.Unlock()
	if !ok ***REMOVED***
		if pluginGetter != nil ***REMOVED*** // this can be nil when the init functions are running
			if l, _ := getPlugin(name, plugingetter.Lookup); l != nil ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ok
***REMOVED***

func (lf *logdriverFactory) registerLogOptValidator(name string, l LogOptValidator) error ***REMOVED***
	lf.m.Lock()
	defer lf.m.Unlock()

	if _, ok := lf.optValidator[name]; ok ***REMOVED***
		return fmt.Errorf("logger: log validator named '%s' is already registered", name)
	***REMOVED***
	lf.optValidator[name] = l
	return nil
***REMOVED***

func (lf *logdriverFactory) get(name string) (Creator, error) ***REMOVED***
	lf.m.Lock()
	defer lf.m.Unlock()

	c, ok := lf.registry[name]
	if ok ***REMOVED***
		return c, nil
	***REMOVED***

	c, err := getPlugin(name, plugingetter.Acquire)
	return c, errors.Wrapf(err, "logger: no log driver named '%s' is registered", name)
***REMOVED***

func (lf *logdriverFactory) getLogOptValidator(name string) LogOptValidator ***REMOVED***
	lf.m.Lock()
	defer lf.m.Unlock()

	c := lf.optValidator[name]
	return c
***REMOVED***

var factory = &logdriverFactory***REMOVED***registry: make(map[string]Creator), optValidator: make(map[string]LogOptValidator)***REMOVED*** // global factory instance

// RegisterLogDriver registers the given logging driver builder with given logging
// driver name.
func RegisterLogDriver(name string, c Creator) error ***REMOVED***
	return factory.register(name, c)
***REMOVED***

// RegisterLogOptValidator registers the logging option validator with
// the given logging driver name.
func RegisterLogOptValidator(name string, l LogOptValidator) error ***REMOVED***
	return factory.registerLogOptValidator(name, l)
***REMOVED***

// GetLogDriver provides the logging driver builder for a logging driver name.
func GetLogDriver(name string) (Creator, error) ***REMOVED***
	return factory.get(name)
***REMOVED***

var builtInLogOpts = map[string]bool***REMOVED***
	"mode":            true,
	"max-buffer-size": true,
***REMOVED***

// ValidateLogOpts checks the options for the given log driver. The
// options supported are specific to the LogDriver implementation.
func ValidateLogOpts(name string, cfg map[string]string) error ***REMOVED***
	if name == "none" ***REMOVED***
		return nil
	***REMOVED***

	switch containertypes.LogMode(cfg["mode"]) ***REMOVED***
	case containertypes.LogModeBlocking, containertypes.LogModeNonBlock, containertypes.LogModeUnset:
	default:
		return fmt.Errorf("logger: logging mode not supported: %s", cfg["mode"])
	***REMOVED***

	if s, ok := cfg["max-buffer-size"]; ok ***REMOVED***
		if containertypes.LogMode(cfg["mode"]) != containertypes.LogModeNonBlock ***REMOVED***
			return fmt.Errorf("logger: max-buffer-size option is only supported with 'mode=%s'", containertypes.LogModeNonBlock)
		***REMOVED***
		if _, err := units.RAMInBytes(s); err != nil ***REMOVED***
			return errors.Wrap(err, "error parsing option max-buffer-size")
		***REMOVED***
	***REMOVED***

	if !factory.driverRegistered(name) ***REMOVED***
		return fmt.Errorf("logger: no log driver named '%s' is registered", name)
	***REMOVED***

	filteredOpts := make(map[string]string, len(builtInLogOpts))
	for k, v := range cfg ***REMOVED***
		if !builtInLogOpts[k] ***REMOVED***
			filteredOpts[k] = v
		***REMOVED***
	***REMOVED***

	validator := factory.getLogOptValidator(name)
	if validator != nil ***REMOVED***
		return validator(filteredOpts)
	***REMOVED***
	return nil
***REMOVED***
