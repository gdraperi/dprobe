package links

import (
	"fmt"
	"path"
	"strings"

	"github.com/docker/go-connections/nat"
)

// Link struct holds informations about parent/child linked container
type Link struct ***REMOVED***
	// Parent container IP address
	ParentIP string
	// Child container IP address
	ChildIP string
	// Link name
	Name string
	// Child environments variables
	ChildEnvironment []string
	// Child exposed ports
	Ports []nat.Port
***REMOVED***

// NewLink initializes a new Link struct with the provided options.
func NewLink(parentIP, childIP, name string, env []string, exposedPorts map[nat.Port]struct***REMOVED******REMOVED***) *Link ***REMOVED***
	var (
		i     int
		ports = make([]nat.Port, len(exposedPorts))
	)

	for p := range exposedPorts ***REMOVED***
		ports[i] = p
		i++
	***REMOVED***

	return &Link***REMOVED***
		Name:             name,
		ChildIP:          childIP,
		ParentIP:         parentIP,
		ChildEnvironment: env,
		Ports:            ports,
	***REMOVED***
***REMOVED***

// ToEnv creates a string's slice containing child container informations in
// the form of environment variables which will be later exported on container
// startup.
func (l *Link) ToEnv() []string ***REMOVED***
	env := []string***REMOVED******REMOVED***

	_, n := path.Split(l.Name)
	alias := strings.Replace(strings.ToUpper(n), "-", "_", -1)

	if p := l.getDefaultPort(); p != nil ***REMOVED***
		env = append(env, fmt.Sprintf("%s_PORT=%s://%s:%s", alias, p.Proto(), l.ChildIP, p.Port()))
	***REMOVED***

	//sort the ports so that we can bulk the continuous ports together
	nat.Sort(l.Ports, func(ip, jp nat.Port) bool ***REMOVED***
		// If the two ports have the same number, tcp takes priority
		// Sort in desc order
		return ip.Int() < jp.Int() || (ip.Int() == jp.Int() && strings.ToLower(ip.Proto()) == "tcp")
	***REMOVED***)

	for i := 0; i < len(l.Ports); ***REMOVED***
		p := l.Ports[i]
		j := nextContiguous(l.Ports, p.Int(), i)
		if j > i+1 ***REMOVED***
			env = append(env, fmt.Sprintf("%s_PORT_%s_%s_START=%s://%s:%s", alias, p.Port(), strings.ToUpper(p.Proto()), p.Proto(), l.ChildIP, p.Port()))
			env = append(env, fmt.Sprintf("%s_PORT_%s_%s_ADDR=%s", alias, p.Port(), strings.ToUpper(p.Proto()), l.ChildIP))
			env = append(env, fmt.Sprintf("%s_PORT_%s_%s_PROTO=%s", alias, p.Port(), strings.ToUpper(p.Proto()), p.Proto()))
			env = append(env, fmt.Sprintf("%s_PORT_%s_%s_PORT_START=%s", alias, p.Port(), strings.ToUpper(p.Proto()), p.Port()))

			q := l.Ports[j]
			env = append(env, fmt.Sprintf("%s_PORT_%s_%s_END=%s://%s:%s", alias, p.Port(), strings.ToUpper(q.Proto()), q.Proto(), l.ChildIP, q.Port()))
			env = append(env, fmt.Sprintf("%s_PORT_%s_%s_PORT_END=%s", alias, p.Port(), strings.ToUpper(q.Proto()), q.Port()))

			i = j + 1
			continue
		***REMOVED*** else ***REMOVED***
			i++
		***REMOVED***
	***REMOVED***
	for _, p := range l.Ports ***REMOVED***
		env = append(env, fmt.Sprintf("%s_PORT_%s_%s=%s://%s:%s", alias, p.Port(), strings.ToUpper(p.Proto()), p.Proto(), l.ChildIP, p.Port()))
		env = append(env, fmt.Sprintf("%s_PORT_%s_%s_ADDR=%s", alias, p.Port(), strings.ToUpper(p.Proto()), l.ChildIP))
		env = append(env, fmt.Sprintf("%s_PORT_%s_%s_PORT=%s", alias, p.Port(), strings.ToUpper(p.Proto()), p.Port()))
		env = append(env, fmt.Sprintf("%s_PORT_%s_%s_PROTO=%s", alias, p.Port(), strings.ToUpper(p.Proto()), p.Proto()))
	***REMOVED***

	// Load the linked container's name into the environment
	env = append(env, fmt.Sprintf("%s_NAME=%s", alias, l.Name))

	if l.ChildEnvironment != nil ***REMOVED***
		for _, v := range l.ChildEnvironment ***REMOVED***
			parts := strings.SplitN(v, "=", 2)
			if len(parts) < 2 ***REMOVED***
				continue
			***REMOVED***
			// Ignore a few variables that are added during docker build (and not really relevant to linked containers)
			if parts[0] == "HOME" || parts[0] == "PATH" ***REMOVED***
				continue
			***REMOVED***
			env = append(env, fmt.Sprintf("%s_ENV_%s=%s", alias, parts[0], parts[1]))
		***REMOVED***
	***REMOVED***
	return env
***REMOVED***

func nextContiguous(ports []nat.Port, value int, index int) int ***REMOVED***
	if index+1 == len(ports) ***REMOVED***
		return index
	***REMOVED***
	for i := index + 1; i < len(ports); i++ ***REMOVED***
		if ports[i].Int() > value+1 ***REMOVED***
			return i - 1
		***REMOVED***

		value++
	***REMOVED***
	return len(ports) - 1
***REMOVED***

// Default port rules
func (l *Link) getDefaultPort() *nat.Port ***REMOVED***
	var p nat.Port
	i := len(l.Ports)

	if i == 0 ***REMOVED***
		return nil
	***REMOVED*** else if i > 1 ***REMOVED***
		nat.Sort(l.Ports, func(ip, jp nat.Port) bool ***REMOVED***
			// If the two ports have the same number, tcp takes priority
			// Sort in desc order
			return ip.Int() < jp.Int() || (ip.Int() == jp.Int() && strings.ToLower(ip.Proto()) == "tcp")
		***REMOVED***)
	***REMOVED***
	p = l.Ports[0]
	return &p
***REMOVED***
