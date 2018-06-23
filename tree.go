package voronoi

import (
	"fmt"

	"github.com/quasoft/dcel"
)

// CircleEvents represents a list of pointers to circle events in which the node participates.
type CircleEvents []*Event

// RemoveEvent removes the given event from the list.
func (ce *CircleEvents) RemoveEvent(event *Event) {
	for i := len(*ce) - 1; i >= 0; i-- {
		if (*ce)[i] == event {
			(*ce)[i] = (*ce)[len(*ce)-1]
			*ce = (*ce)[:len(*ce)-1]
		}
	}
}

// HasEvent tests if the node has a pointer to the given event.
func (ce *CircleEvents) HasEvent(event *Event) bool {
	for i := 0; i < len(*ce); i++ {
		if (*ce)[i] == event {
			return true
		}
	}
	return false
}

// Node represent an element in a binary tree.
// Each Leaf in the tree represents an arc of a parabola (part of parabola),
// that lies on the beach line. Leaf nodes store the site that created the arc
// and pointers to the circle events associated with it.
// Internal nodes represent intersections (breakpoints) between the arcs and
// store no values.
type Node struct {
	// Site is the focus of the parabola arc (the site which created the parabola).
	// Not used for internal nodes.
	Site *Site
	// Events hold pointers to circle events, in which this arc is the left most, middle or right-most arc.
	// Not used for internal nodes.
	LeftEvents, MiddleEvents, RightEvents CircleEvents
	// Pointer to the parent node.
	Parent *Node
	// Left stores a subtree of arcs with smaller X values.
	Left *Node
	// Right stores a subtree of arcs with larger X values.
	Right *Node

	LeftEdges  []*dcel.HalfEdge
	RightEdges []*dcel.HalfEdge
}

// String method from https://github.com/golang/tour/blob/master/tree/tree.go
func (n *Node) String() string {
	if n == nil {
		return "()"
	}
	s := ""
	if n.Left != nil {
		s += n.Left.String() + " "
	}

	if n.IsLeaf() {
		s += "[" + fmt.Sprint(n.Site) + "]"
	} else if n.Parent == nil {
		s += "<root>"
	} else {
		s += "<int>" // internal
	}

	if n.Right != nil {
		s += " " + n.Right.String()
	}

	return "(" + s + ")"
}

// IsLeaf returns true if the TreeNode has no left or right children.
// A single root node is also considered a leaf.
func (n *Node) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

// PrevChildArc returns the node for the previous arc.
func (n *Node) PrevChildArc() *Node {
	left := n.Left
	for !left.IsLeaf() {
		left = left.Right
	}
	return left
}

// NextChildArc returns the node for the next arc.
func (n *Node) NextChildArc() *Node {
	right := n.Right
	for !right.IsLeaf() {
		right = right.Left
	}
	return right
}

// PrevArc returns the node for the previous arc.
func (n *Node) PrevArc() *Node {
	if n == nil {
		return nil
	}

	// If an internal node, traverse down
	if !n.IsLeaf() {
		return n.LastArc()
	}

	// If a leaf, traverse up
	if n.Parent == nil {
		return nil
	}

	parent := n.Parent
	node := n
	for parent.Left == node {
		if parent.Parent == nil {
			return nil
		}
		node = parent
		parent = parent.Parent
	}

	if parent.Left.IsLeaf() {
		return parent.Left
	}

	return parent.Left.LastArc()
}

// NextArc returns the node for the next arc.
func (n *Node) NextArc() *Node {
	if n == nil {
		return nil
	}

	// If an internal node, traverse down
	if !n.IsLeaf() {
		return n.FirstArc()
	}

	// If a leaf, traverse up
	if n.Parent == nil {
		return nil
	}

	parent := n.Parent
	node := n
	for parent.Right == node {
		if parent.Parent == nil {
			return nil
		}
		node = parent
		parent = parent.Parent
	}

	if parent.Right.IsLeaf() {
		return parent.Right
	}

	return parent.Right.FirstArc()
}

// FirstArc returns the  left-most arc (leaf) in the tree.
func (n *Node) FirstArc() *Node {
	first := n
	for first != nil && !first.IsLeaf() {
		if first.Left != nil {
			first = first.Left
		} else {
			first = first.Right
		}
	}
	return first
}

// LastArc returns the right-most arc (leaf) in the tree.
func (n *Node) LastArc() *Node {
	last := n
	for last != nil && !last.IsLeaf() {
		if last.Right != nil {
			last = last.Right
		} else {
			last = last.Left
		}
	}
	return last
}

// AddLeftEvent pushes a pointer to an event for which this is the left-most node.
func (n *Node) AddLeftEvent(event *Event) {
	n.LeftEvents = append(n.LeftEvents, event)
}

// AddMiddleEvent pushes a pointer to an event for which this is the left-most node.
func (n *Node) AddMiddleEvent(event *Event) {
	n.MiddleEvents = append(n.MiddleEvents, event)
}

// AddRightEvent pushes a pointer to an event for which this is the left-most node.
func (n *Node) AddRightEvent(event *Event) {
	n.RightEvents = append(n.RightEvents, event)
}

// RemoveEvent removes the given event from the lists of the node.
func (n *Node) RemoveEvent(event *Event) {
	n.LeftEvents.RemoveEvent(event)
	n.MiddleEvents.RemoveEvent(event)
	n.RightEvents.RemoveEvent(event)
}

// HasEvent tests if the node has a pointer to the given event.
func (n *Node) HasEvent(event *Event) bool {
	return n.LeftEvents.HasEvent(event) || n.MiddleEvents.HasEvent(event) ||
		n.RightEvents.HasEvent(event)
}
