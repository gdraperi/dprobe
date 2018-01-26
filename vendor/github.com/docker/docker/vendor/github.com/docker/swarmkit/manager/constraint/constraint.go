package constraint

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/docker/swarmkit/api"
)

const (
	eq = iota
	noteq

	// NodeLabelPrefix is the constraint key prefix for node labels.
	NodeLabelPrefix = "node.labels."
	// EngineLabelPrefix is the constraint key prefix for engine labels.
	EngineLabelPrefix = "engine.labels."
)

var (
	alphaNumeric = regexp.MustCompile(`^(?i)[a-z_][a-z0-9\-_.]+$`)
	// value can be alphanumeric and some special characters. it shouldn't container
	// current or future operators like '>, <, ~', etc.
	valuePattern = regexp.MustCompile(`^(?i)[a-z0-9:\-_\s\.\*\(\)\?\+\[\]\\\^\$\|\/]+$`)

	// operators defines list of accepted operators
	operators = []string***REMOVED***"==", "!="***REMOVED***
)

// Constraint defines a constraint.
type Constraint struct ***REMOVED***
	key      string
	operator int
	exp      string
***REMOVED***

// Parse parses list of constraints.
func Parse(env []string) ([]Constraint, error) ***REMOVED***
	exprs := []Constraint***REMOVED******REMOVED***
	for _, e := range env ***REMOVED***
		found := false
		// each expr is in the form of "key op value"
		for i, op := range operators ***REMOVED***
			if !strings.Contains(e, op) ***REMOVED***
				continue
			***REMOVED***
			// split with the op
			parts := strings.SplitN(e, op, 2)

			if len(parts) < 2 ***REMOVED***
				return nil, fmt.Errorf("invalid expr: %s", e)
			***REMOVED***

			part0 := strings.TrimSpace(parts[0])
			// validate key
			matched := alphaNumeric.MatchString(part0)
			if matched == false ***REMOVED***
				return nil, fmt.Errorf("key '%s' is invalid", part0)
			***REMOVED***

			part1 := strings.TrimSpace(parts[1])

			// validate Value
			matched = valuePattern.MatchString(part1)
			if matched == false ***REMOVED***
				return nil, fmt.Errorf("value '%s' is invalid", part1)
			***REMOVED***
			// TODO(dongluochen): revisit requirements to see if globing or regex are useful
			exprs = append(exprs, Constraint***REMOVED***key: part0, operator: i, exp: part1***REMOVED***)

			found = true
			break // found an op, move to next entry
		***REMOVED***
		if !found ***REMOVED***
			return nil, fmt.Errorf("constraint expected one operator from %s", strings.Join(operators, ", "))
		***REMOVED***
	***REMOVED***
	return exprs, nil
***REMOVED***

// Match checks if the Constraint matches the target strings.
func (c *Constraint) Match(whats ...string) bool ***REMOVED***
	var match bool

	// full string match
	for _, what := range whats ***REMOVED***
		// case insensitive compare
		if strings.EqualFold(c.exp, what) ***REMOVED***
			match = true
			break
		***REMOVED***
	***REMOVED***

	switch c.operator ***REMOVED***
	case eq:
		return match
	case noteq:
		return !match
	***REMOVED***

	return false
***REMOVED***

// NodeMatches returns true if the node satisfies the given constraints.
func NodeMatches(constraints []Constraint, n *api.Node) bool ***REMOVED***
	for _, constraint := range constraints ***REMOVED***
		switch ***REMOVED***
		case strings.EqualFold(constraint.key, "node.id"):
			if !constraint.Match(n.ID) ***REMOVED***
				return false
			***REMOVED***
		case strings.EqualFold(constraint.key, "node.hostname"):
			// if this node doesn't have hostname
			// it's equivalent to match an empty hostname
			// where '==' would fail, '!=' matches
			if n.Description == nil ***REMOVED***
				if !constraint.Match("") ***REMOVED***
					return false
				***REMOVED***
				continue
			***REMOVED***
			if !constraint.Match(n.Description.Hostname) ***REMOVED***
				return false
			***REMOVED***
		case strings.EqualFold(constraint.key, "node.ip"):
			nodeIP := net.ParseIP(n.Status.Addr)
			// single IP address, node.ip == 2001:db8::2
			if ip := net.ParseIP(constraint.exp); ip != nil ***REMOVED***
				ipEq := ip.Equal(nodeIP)
				if (ipEq && constraint.operator != eq) || (!ipEq && constraint.operator == eq) ***REMOVED***
					return false
				***REMOVED***
				continue
			***REMOVED***
			// CIDR subnet, node.ip != 210.8.4.0/24
			if _, subnet, err := net.ParseCIDR(constraint.exp); err == nil ***REMOVED***
				within := subnet.Contains(nodeIP)
				if (within && constraint.operator != eq) || (!within && constraint.operator == eq) ***REMOVED***
					return false
				***REMOVED***
				continue
			***REMOVED***
			// reject constraint with malformed address/network
			return false
		case strings.EqualFold(constraint.key, "node.role"):
			if !constraint.Match(n.Role.String()) ***REMOVED***
				return false
			***REMOVED***
		case strings.EqualFold(constraint.key, "node.platform.os"):
			if n.Description == nil || n.Description.Platform == nil ***REMOVED***
				if !constraint.Match("") ***REMOVED***
					return false
				***REMOVED***
				continue
			***REMOVED***
			if !constraint.Match(n.Description.Platform.OS) ***REMOVED***
				return false
			***REMOVED***
		case strings.EqualFold(constraint.key, "node.platform.arch"):
			if n.Description == nil || n.Description.Platform == nil ***REMOVED***
				if !constraint.Match("") ***REMOVED***
					return false
				***REMOVED***
				continue
			***REMOVED***
			if !constraint.Match(n.Description.Platform.Architecture) ***REMOVED***
				return false
			***REMOVED***

		// node labels constraint in form like 'node.labels.key==value'
		case len(constraint.key) > len(NodeLabelPrefix) && strings.EqualFold(constraint.key[:len(NodeLabelPrefix)], NodeLabelPrefix):
			if n.Spec.Annotations.Labels == nil ***REMOVED***
				if !constraint.Match("") ***REMOVED***
					return false
				***REMOVED***
				continue
			***REMOVED***
			label := constraint.key[len(NodeLabelPrefix):]
			// label itself is case sensitive
			val := n.Spec.Annotations.Labels[label]
			if !constraint.Match(val) ***REMOVED***
				return false
			***REMOVED***

		// engine labels constraint in form like 'engine.labels.key!=value'
		case len(constraint.key) > len(EngineLabelPrefix) && strings.EqualFold(constraint.key[:len(EngineLabelPrefix)], EngineLabelPrefix):
			if n.Description == nil || n.Description.Engine == nil || n.Description.Engine.Labels == nil ***REMOVED***
				if !constraint.Match("") ***REMOVED***
					return false
				***REMOVED***
				continue
			***REMOVED***
			label := constraint.key[len(EngineLabelPrefix):]
			val := n.Description.Engine.Labels[label]
			if !constraint.Match(val) ***REMOVED***
				return false
			***REMOVED***
		default:
			// key doesn't match predefined syntax
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***
