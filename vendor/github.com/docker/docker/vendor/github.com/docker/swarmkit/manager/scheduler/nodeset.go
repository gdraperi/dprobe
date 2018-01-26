package scheduler

import (
	"container/heap"
	"errors"
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/constraint"
)

var errNodeNotFound = errors.New("node not found in scheduler dataset")

type nodeSet struct ***REMOVED***
	nodes map[string]NodeInfo // map from node id to node info
***REMOVED***

func (ns *nodeSet) alloc(n int) ***REMOVED***
	ns.nodes = make(map[string]NodeInfo, n)
***REMOVED***

// nodeInfo returns the NodeInfo struct for a given node identified by its ID.
func (ns *nodeSet) nodeInfo(nodeID string) (NodeInfo, error) ***REMOVED***
	node, ok := ns.nodes[nodeID]
	if ok ***REMOVED***
		return node, nil
	***REMOVED***
	return NodeInfo***REMOVED******REMOVED***, errNodeNotFound
***REMOVED***

// addOrUpdateNode sets the number of tasks for a given node. It adds the node
// to the set if it wasn't already tracked.
func (ns *nodeSet) addOrUpdateNode(n NodeInfo) ***REMOVED***
	ns.nodes[n.ID] = n
***REMOVED***

// updateNode sets the number of tasks for a given node. It ignores the update
// if the node isn't already tracked in the set.
func (ns *nodeSet) updateNode(n NodeInfo) ***REMOVED***
	_, ok := ns.nodes[n.ID]
	if ok ***REMOVED***
		ns.nodes[n.ID] = n
	***REMOVED***
***REMOVED***

func (ns *nodeSet) remove(nodeID string) ***REMOVED***
	delete(ns.nodes, nodeID)
***REMOVED***

func (ns *nodeSet) tree(serviceID string, preferences []*api.PlacementPreference, maxAssignments int, meetsConstraints func(*NodeInfo) bool, nodeLess func(*NodeInfo, *NodeInfo) bool) decisionTree ***REMOVED***
	var root decisionTree

	if maxAssignments == 0 ***REMOVED***
		return root
	***REMOVED***

	for _, node := range ns.nodes ***REMOVED***
		tree := &root
		for _, pref := range preferences ***REMOVED***
			// Only spread is supported so far
			spread := pref.GetSpread()
			if spread == nil ***REMOVED***
				continue
			***REMOVED***

			descriptor := spread.SpreadDescriptor
			var value string
			switch ***REMOVED***
			case len(descriptor) > len(constraint.NodeLabelPrefix) && strings.EqualFold(descriptor[:len(constraint.NodeLabelPrefix)], constraint.NodeLabelPrefix):
				if node.Spec.Annotations.Labels != nil ***REMOVED***
					value = node.Spec.Annotations.Labels[descriptor[len(constraint.NodeLabelPrefix):]]
				***REMOVED***
			case len(descriptor) > len(constraint.EngineLabelPrefix) && strings.EqualFold(descriptor[:len(constraint.EngineLabelPrefix)], constraint.EngineLabelPrefix):
				if node.Description != nil && node.Description.Engine != nil && node.Description.Engine.Labels != nil ***REMOVED***
					value = node.Description.Engine.Labels[descriptor[len(constraint.EngineLabelPrefix):]]
				***REMOVED***
			// TODO(aaronl): Support other items from constraint
			// syntax like node ID, hostname, os/arch, etc?
			default:
				continue
			***REMOVED***

			// If value is still uninitialized, the value used for
			// the node at this level of the tree is "". This makes
			// sure that the tree structure is not affected by
			// which properties nodes have and don't have.

			if node.ActiveTasksCountByService != nil ***REMOVED***
				tree.tasks += node.ActiveTasksCountByService[serviceID]
			***REMOVED***

			if tree.next == nil ***REMOVED***
				tree.next = make(map[string]*decisionTree)
			***REMOVED***
			next := tree.next[value]
			if next == nil ***REMOVED***
				next = &decisionTree***REMOVED******REMOVED***
				tree.next[value] = next
			***REMOVED***
			tree = next
		***REMOVED***

		if node.ActiveTasksCountByService != nil ***REMOVED***
			tree.tasks += node.ActiveTasksCountByService[serviceID]
		***REMOVED***

		if tree.nodeHeap.lessFunc == nil ***REMOVED***
			tree.nodeHeap.lessFunc = nodeLess
		***REMOVED***

		if tree.nodeHeap.Len() < maxAssignments ***REMOVED***
			if meetsConstraints(&node) ***REMOVED***
				heap.Push(&tree.nodeHeap, node)
			***REMOVED***
		***REMOVED*** else if nodeLess(&node, &tree.nodeHeap.nodes[0]) ***REMOVED***
			if meetsConstraints(&node) ***REMOVED***
				tree.nodeHeap.nodes[0] = node
				heap.Fix(&tree.nodeHeap, 0)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return root
***REMOVED***
