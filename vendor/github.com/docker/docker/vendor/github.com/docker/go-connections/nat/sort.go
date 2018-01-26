package nat

import (
	"sort"
	"strings"
)

type portSorter struct ***REMOVED***
	ports []Port
	by    func(i, j Port) bool
***REMOVED***

func (s *portSorter) Len() int ***REMOVED***
	return len(s.ports)
***REMOVED***

func (s *portSorter) Swap(i, j int) ***REMOVED***
	s.ports[i], s.ports[j] = s.ports[j], s.ports[i]
***REMOVED***

func (s *portSorter) Less(i, j int) bool ***REMOVED***
	ip := s.ports[i]
	jp := s.ports[j]

	return s.by(ip, jp)
***REMOVED***

// Sort sorts a list of ports using the provided predicate
// This function should compare `i` and `j`, returning true if `i` is
// considered to be less than `j`
func Sort(ports []Port, predicate func(i, j Port) bool) ***REMOVED***
	s := &portSorter***REMOVED***ports, predicate***REMOVED***
	sort.Sort(s)
***REMOVED***

type portMapEntry struct ***REMOVED***
	port    Port
	binding PortBinding
***REMOVED***

type portMapSorter []portMapEntry

func (s portMapSorter) Len() int      ***REMOVED*** return len(s) ***REMOVED***
func (s portMapSorter) Swap(i, j int) ***REMOVED*** s[i], s[j] = s[j], s[i] ***REMOVED***

// sort the port so that the order is:
// 1. port with larger specified bindings
// 2. larger port
// 3. port with tcp protocol
func (s portMapSorter) Less(i, j int) bool ***REMOVED***
	pi, pj := s[i].port, s[j].port
	hpi, hpj := toInt(s[i].binding.HostPort), toInt(s[j].binding.HostPort)
	return hpi > hpj || pi.Int() > pj.Int() || (pi.Int() == pj.Int() && strings.ToLower(pi.Proto()) == "tcp")
***REMOVED***

// SortPortMap sorts the list of ports and their respected mapping. The ports
// will explicit HostPort will be placed first.
func SortPortMap(ports []Port, bindings PortMap) ***REMOVED***
	s := portMapSorter***REMOVED******REMOVED***
	for _, p := range ports ***REMOVED***
		if binding, ok := bindings[p]; ok ***REMOVED***
			for _, b := range binding ***REMOVED***
				s = append(s, portMapEntry***REMOVED***port: p, binding: b***REMOVED***)
			***REMOVED***
			bindings[p] = []PortBinding***REMOVED******REMOVED***
		***REMOVED*** else ***REMOVED***
			s = append(s, portMapEntry***REMOVED***port: p***REMOVED***)
		***REMOVED***
	***REMOVED***

	sort.Sort(s)
	var (
		i  int
		pm = make(map[Port]struct***REMOVED******REMOVED***)
	)
	// reorder ports
	for _, entry := range s ***REMOVED***
		if _, ok := pm[entry.port]; !ok ***REMOVED***
			ports[i] = entry.port
			pm[entry.port] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			i++
		***REMOVED***
		// reorder bindings for this port
		if _, ok := bindings[entry.port]; ok ***REMOVED***
			bindings[entry.port] = append(bindings[entry.port], entry.binding)
		***REMOVED***
	***REMOVED***
***REMOVED***

func toInt(s string) uint64 ***REMOVED***
	i, _, err := ParsePortRange(s)
	if err != nil ***REMOVED***
		i = 0
	***REMOVED***
	return i
***REMOVED***
