package voronoi

import (
	"fmt"

	"github.com/quasoft/dcel"
)

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
	// Events hold pointers to all circle events, in which this arc participates.
	// Not used for internal nodes.
	Events []*Event
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

// AddEvent pushes a pointer to an event in the Events list of the node.
func (n *Node) AddEvent(event *Event) {
	n.Events = append(n.Events, event)
}
