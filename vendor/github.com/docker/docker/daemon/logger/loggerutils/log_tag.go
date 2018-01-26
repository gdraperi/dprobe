package loggerutils

import (
	"bytes"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/templates"
)

// DefaultTemplate defines the defaults template logger should use.
const DefaultTemplate = "***REMOVED******REMOVED***.ID***REMOVED******REMOVED***"

// ParseLogTag generates a context aware tag for consistency across different
// log drivers based on the context of the running container.
func ParseLogTag(info logger.Info, defaultTemplate string) (string, error) ***REMOVED***
	tagTemplate := info.Config["tag"]
	if tagTemplate == "" ***REMOVED***
		tagTemplate = defaultTemplate
	***REMOVED***

	tmpl, err := templates.NewParse("log-tag", tagTemplate)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, &info); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return buf.String(), nil
***REMOVED***
