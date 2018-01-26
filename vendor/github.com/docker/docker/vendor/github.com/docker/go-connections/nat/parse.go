package nat

import (
	"fmt"
	"strconv"
	"strings"
)

// PartParser parses and validates the specified string (data) using the specified template
// e.g. ip:public:private -> 192.168.0.1:80:8000
// DEPRECATED: do not use, this function may be removed in a future version
func PartParser(template, data string) (map[string]string, error) ***REMOVED***
	// ip:public:private
	var (
		templateParts = strings.Split(template, ":")
		parts         = strings.Split(data, ":")
		out           = make(map[string]string, len(templateParts))
	)
	if len(parts) != len(templateParts) ***REMOVED***
		return nil, fmt.Errorf("Invalid format to parse. %s should match template %s", data, template)
	***REMOVED***

	for i, t := range templateParts ***REMOVED***
		value := ""
		if len(parts) > i ***REMOVED***
			value = parts[i]
		***REMOVED***
		out[t] = value
	***REMOVED***
	return out, nil
***REMOVED***

// ParsePortRange parses and validates the specified string as a port-range (8000-9000)
func ParsePortRange(ports string) (uint64, uint64, error) ***REMOVED***
	if ports == "" ***REMOVED***
		return 0, 0, fmt.Errorf("Empty string specified for ports.")
	***REMOVED***
	if !strings.Contains(ports, "-") ***REMOVED***
		start, err := strconv.ParseUint(ports, 10, 16)
		end := start
		return start, end, err
	***REMOVED***

	parts := strings.Split(ports, "-")
	start, err := strconv.ParseUint(parts[0], 10, 16)
	if err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***
	end, err := strconv.ParseUint(parts[1], 10, 16)
	if err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***
	if end < start ***REMOVED***
		return 0, 0, fmt.Errorf("Invalid range specified for the Port: %s", ports)
	***REMOVED***
	return start, end, nil
***REMOVED***
