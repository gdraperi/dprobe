package types

import (
	"encoding/json"
	"fmt"
	"sort"
)

// PluginsListResponse contains the response for the Engine API
type PluginsListResponse []*Plugin

// UnmarshalJSON implements json.Unmarshaler for PluginInterfaceType
func (t *PluginInterfaceType) UnmarshalJSON(p []byte) error ***REMOVED***
	versionIndex := len(p)
	prefixIndex := 0
	if len(p) < 2 || p[0] != '"' || p[len(p)-1] != '"' ***REMOVED***
		return fmt.Errorf("%q is not a plugin interface type", p)
	***REMOVED***
	p = p[1 : len(p)-1]
loop:
	for i, b := range p ***REMOVED***
		switch b ***REMOVED***
		case '.':
			prefixIndex = i
		case '/':
			versionIndex = i
			break loop
		***REMOVED***
	***REMOVED***
	t.Prefix = string(p[:prefixIndex])
	t.Capability = string(p[prefixIndex+1 : versionIndex])
	if versionIndex < len(p) ***REMOVED***
		t.Version = string(p[versionIndex+1:])
	***REMOVED***
	return nil
***REMOVED***

// MarshalJSON implements json.Marshaler for PluginInterfaceType
func (t *PluginInterfaceType) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(t.String())
***REMOVED***

// String implements fmt.Stringer for PluginInterfaceType
func (t PluginInterfaceType) String() string ***REMOVED***
	return fmt.Sprintf("%s.%s/%s", t.Prefix, t.Capability, t.Version)
***REMOVED***

// PluginPrivilege describes a permission the user has to accept
// upon installing a plugin.
type PluginPrivilege struct ***REMOVED***
	Name        string
	Description string
	Value       []string
***REMOVED***

// PluginPrivileges is a list of PluginPrivilege
type PluginPrivileges []PluginPrivilege

func (s PluginPrivileges) Len() int ***REMOVED***
	return len(s)
***REMOVED***

func (s PluginPrivileges) Less(i, j int) bool ***REMOVED***
	return s[i].Name < s[j].Name
***REMOVED***

func (s PluginPrivileges) Swap(i, j int) ***REMOVED***
	sort.Strings(s[i].Value)
	sort.Strings(s[j].Value)
	s[i], s[j] = s[j], s[i]
***REMOVED***
