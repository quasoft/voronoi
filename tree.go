package voronoi

import "fmt"

// VNode represent an element in a binary tree.
// Each Leaf in the tree represents an arc of a parabola (part of parabola),
// that lies on the beach line. Leaf nodes store the site that created the arc
// and pointers to the circle events associated with it.
// Internal nodes represent intersections (breakpoints) between the arcs and
// store no values.
type VNode struct {
	// Site is the focus of the parabola arc (the site which created the parabola).
	// Not used for internal nodes.
	Site Site
	// Events hold pointers to all circle events, in which this arc participates.
	// Not used for internal nodes.
	Events []*Event
	// Pointer to the parent node.
	Parent *VNode
	// Left stores a subtree of arcs with smaller X values.
	Left *VNode
	// Right stores a subtree of arcs with larger X values.
	Right *VNode
}

// String method from https://github.com/golang/tour/blob/master/tree/tree.go
func (n *VNode) String() string {
	if n == nil {
		return "()"
	}
	s := ""
	if n.Left != nil {
		s += n.Left.String() + " "
	}
	s += fmt.Sprint(n.Site)
	if n.Right != nil {
		s += " " + n.Right.String()
	}
	return "(" + s + ")"
}

// IsLeaf returns true if the TreeNode has no left or right children.
// A single root node is also considered a leaf.
func (n *VNode) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

// PrevChildArc returns the node for the previous arc.
func (n *VNode) PrevChildArc() *VNode {
	left := n.Left
	for !left.IsLeaf() {
		left = left.Right
	}
	return left
}

// NextChildArc returns the node for the next arc.
func (n *VNode) NextChildArc() *VNode {
	right := n.Right
	for !right.IsLeaf() {
		right = right.Left
	}
	return right
}

// PrevArc returns the node for the previous arc.
func (n *VNode) PrevArc() *VNode {
	// If an internal node, traverse down
	if !n.IsLeaf() {
		return n.PrevChildArc()
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

	return parent.Left.NextChildArc()
}

// NextArc returns the node for the next arc.
func (n *VNode) NextArc() *VNode {
	// If an internal node, traverse down
	if !n.IsLeaf() {
		return n.NextChildArc()
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

	return parent.Right.NextChildArc()
}

// AddEvent pushes a pointer to an event in the Events list of the node.
func (n *VNode) AddEvent(event *Event) {
	n.Events = append(n.Events, event)
}
