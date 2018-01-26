// +build linux freebsd darwin

package volume

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (p *linuxParser) HasResource(m *MountPoint, absolutePath string) bool ***REMOVED***
	relPath, err := filepath.Rel(m.Destination, absolutePath)
	return err == nil && relPath != ".." && !strings.HasPrefix(relPath, fmt.Sprintf("..%c", filepath.Separator))
***REMOVED***

func (p *windowsParser) HasResource(m *MountPoint, absolutePath string) bool ***REMOVED***
	return false
***REMOVED***
