package voronoi

import (
	"container/heap"
	"fmt"
	"image"
	"log"
	"math"

	"github.com/quasoft/dcel"
)

// Voronoi implements Fortune's algorithm for voronoi diagram generation.
type Voronoi struct {
	Bounds       image.Rectangle
	Sites        SiteSlice
	EventQueue   EventQueue
	ParabolaTree *Node
	SweepLine    int // tracks the current position of the sweep line; updated when a new site is added.
	DCEL         *dcel.DCEL
}

// New creates a voronoi diagram generator for a list of sites and within the specified bounds.
func New(sites SiteSlice, bounds image.Rectangle) *Voronoi {
	voronoi := &Voronoi{Bounds: bounds}
	voronoi.Sites = make(SiteSlice, len(sites), len(sites))
	copy(voronoi.Sites, sites)
	voronoi.init()
	return voronoi
}

// NewFromPoints creates a voronoi diagram generator for a list of points within the specified bounds.
func NewFromPoints(points []image.Point, bounds image.Rectangle) *Voronoi {
	var sites SiteSlice
	var id int64
	for _, point := range points {
		sites = append(sites, Site{
			X:  point.X,
			Y:  point.Y,
			ID: id,
		})
		id++
	}
	return New(sites, bounds)
}

func (v *Voronoi) init() {
	// 1. Push sites to a priority queue, sorted by by Y
	// 2. Create empty binary tree for parabola arcs
	// 3. Create empty doubly-connected edge list (DCEL) for the voronoi diagram

	// 1. Push sites to a priority queue, sorted by by Y
	v.EventQueue = NewEventQueue(v.Sites)

	// 2. Create empty binary tree for parabola arcs
	v.ParabolaTree = nil

	// 3. Create empty doubly-connected edge list (DCEL) for the voronoi diagram
	v.DCEL = dcel.NewDCEL()
}

// Reset clears the state of the voronoi generator.
func (v *Voronoi) Reset() {
	v.EventQueue = NewEventQueue(v.Sites)
	v.ParabolaTree = nil
	v.SweepLine = 0
	v.DCEL = dcel.NewDCEL()
}

// HandleNextEvent processes the next event from the internal event queue.
// Used from the player application while developing the algorithm.
func (v *Voronoi) HandleNextEvent() {
	if v.EventQueue.Len() <= 0 {
		return
	}

	// Process events by Y (priority)
	event := heap.Pop(&v.EventQueue).(*Event)

	// Event with Y above the sweep line should be ignored.
	if event.Y < v.SweepLine {
		log.Printf("Ignoring event with Y %d as it's above the sweep line (%d)\r\n", event.Y, v.SweepLine)
		return
	}

	v.SweepLine = event.Y
	if event.EventType == EventSite {
		v.handleSiteEvent(event)
	} else {
		v.handleCircleEvent(event)
	}
}

// Generate runs the algorithm for the given sites and bounds, creating a voronoi diagram.
func (v *Voronoi) Generate() {
	v.Reset()

	// While queue is not empty
	for v.EventQueue.Len() > 0 {
		v.HandleNextEvent()
	}
}

// findNodeAbove finds the node for the parabola that is vertically above the specified site.
func (v *Voronoi) findNodeAbove(site *Site) *Node {
	node := v.ParabolaTree

	for !node.IsLeaf() {
		if node.IsLeaf() {
			log.Printf("At leaf %v\r\n", node)
		} else {
			log.Printf("At internal node %v <-> %v\r\n", node.PrevChildArc(), node.NextChildArc())
		}

		x, err := GetXOfInternalNode(node, v.SweepLine)
		if err != nil {
			panic(fmt.Errorf("could not find arc above %v - this should never happen", node))
		}
		if site.X < x {
			log.Printf("site.X (%d) < x (%d), going left\r\n", site.X, x)
			node = node.Left
		} else {
			log.Printf("site.X (%d) >= x (%d), going right\r\n", site.X, x)
			node = node.Right
		}
		if node.IsLeaf() {
			log.Printf("X of intersection: %d\r\n", x)
		}
	}

	return node
}

func (v *Voronoi) handleSiteEvent(event *Event) {
	log.Println()
	log.Printf("Handling site event %d,%d\r\n", event.X, event.Y)
	log.Printf("Sweep line: %d", v.SweepLine)
	log.Printf("Tree: %v", v.ParabolaTree)

	// Create a face for this site and link it to it
	face := v.DCEL.NewFace()
	face.ID = event.Site.ID
	face.Data = event.Site
	event.Site.Face = face

	// If the binary tree is empty, just add an arc for this site as the only leaf in the tree
	if v.ParabolaTree == nil {
		log.Print("Adding event as root\r\n")
		v.ParabolaTree = &Node{Site: event.Site}
		return
	}

	// If the tree is not empty, find the arc vertically above the new site
	arcAbove := v.findNodeAbove(event.Site)
	if arcAbove == nil {
		log.Print("Could not find arc above event site!\r\n")
		// Do something
		return
	}
	log.Printf("Arc above: %v\r\n", arcAbove)

	v.removeCircleEvent(arcAbove)

	y := GetYByX(arcAbove.Site, event.Site.X, v.SweepLine)
	vertex := v.DCEL.NewVertex(event.Site.X, y)
	log.Printf("Y of intersection = %d,%d\r\n", vertex.X, vertex.Y)

	// The node above (NA) is replaced wit ha branch with one internal node and three leafs.
	// The middle leaf stores the new parabola and the other two store the one being split.
	//    (NA)
	//   /   \
	//  (  )  [old]
	// /    \
	//[old]  [new]

	// Copy of the old arc
	arcAbove.Right = &Node{
		Site:   arcAbove.Site,
		Events: arcAbove.Events,
		Parent: arcAbove,
	}
	oldArcRight := arcAbove.Right
	oldArcRight.RightEdges = make([]*dcel.HalfEdge, len(arcAbove.RightEdges))
	copy(oldArcRight.RightEdges, arcAbove.RightEdges)

	// Internal node
	arcAbove.Left = &Node{Parent: arcAbove}

	// The new arc
	arcAbove.Left.Right = &Node{
		Site:   event.Site,
		Parent: arcAbove.Left,
	}
	newArc := arcAbove.Left.Right

	// Copy of the old arc
	arcAbove.Left.Left = &Node{
		Site:   arcAbove.Site,
		Events: arcAbove.Events,
		Parent: arcAbove.Left,
	}
	oldArcLeft := arcAbove.Left.Left
	oldArcLeft.LeftEdges = make([]*dcel.HalfEdge, len(arcAbove.LeftEdges))
	copy(oldArcLeft.LeftEdges, arcAbove.LeftEdges)

	// Internal nodes have no site
	arcAbove.Site = nil
	arcAbove.Events = nil

	// Add four new half-edges in DCEL and add a pointer to those
	// half-edges from the arcs which are tracing them.
	edge1, edge2 := v.DCEL.NewEdge(oldArcLeft.Site.Face, newArc.Site.Face, vertex)
	oldArcLeft.RightEdges = append(oldArcLeft.RightEdges, edge1)
	newArc.LeftEdges = append(newArc.LeftEdges, edge2)

	edge3, edge4 := v.DCEL.NewEdge(newArc.Site.Face, oldArcRight.Site.Face, vertex)
	newArc.RightEdges = append(newArc.RightEdges, edge3)
	oldArcRight.LeftEdges = append(oldArcRight.LeftEdges, edge4)

	// Check for circle events where the new arc is the right most arc
	prevArc := newArc.PrevArc()
	log.Printf("Prev arc for %v is %v\r\n", newArc, prevArc)
	prevPrevArc := prevArc.PrevArc()
	log.Printf("Prev->prev arc for %v is %v\r\n", newArc, prevPrevArc)
	v.addCircleEvent(prevPrevArc, prevArc, newArc)

	// Check for circle events where the new arc is the left most arc
	nextArc := newArc.NextArc()
	nextNextArc := nextArc.NextArc()
	v.addCircleEvent(newArc, nextArc, nextNextArc)
}

// calcCircle checks if the circle passing through three sites is counter-clockwise,
// and retunrs the center of the circle and it's radius if it is.
func (v *Voronoi) calcCircle(site1, site2, site3 *Site) (x int, y int, r int, err error) {
	// Solution by https://math.stackexchange.com/a/1268279/543428
	// Explanation at http://mathforum.org/library/drmath/view/55002.html
	x = 0
	y = 0
	r = 0
	err = nil

	x1 := float64(site1.X)
	y1 := float64(site1.Y)

	x2 := float64(site2.X)
	y2 := float64(site2.Y)

	x3 := float64(site3.X)
	y3 := float64(site3.Y)

	// If circle is oriented clockwise (there is a circle, but the sites are in reverse order),
	// then ignore this circle.
	// Code for that part adapted from: https://github.com/gorhill/Javascript-Voronoi/blob/master/rhill-voronoi-core.js
	// Explanation at https://en.wikipedia.org/wiki/Curve_orientation#Orientation_of_a_simple_polygon
	determinant := (x2*y3 + x1*y2 + y1*x3) - (y1*x2 + y2*x3 + x1*y3)
	if determinant < 0 {
		log.Printf("Sites are in reversed order, so circle would be clockwise")
		err = fmt.Errorf("circle is clockwise - sites %f,%f %f,%f %f,%f are in reversed order", x1, y1, x2, y2, x3, y3)
		return
	}

	if x2-x1 == 0 || x3-x2 == 0 {
		log.Printf("Ignoring circle, division by zero")
		err = fmt.Errorf("no circle found connecting points %f,%f %f,%f and %f,%f", x1, y1, x2, y2, x3, y3)
		return
	}

	mr := (y2 - y1) / (x2 - x1)
	mt := (y3 - y2) / (x3 - x2)

	if mr == mt || mr-mt == 0 || mr == 0 {
		log.Printf("Ignoring circle, division by zero")
		err = fmt.Errorf("no circle found connecting points %f,%f %f,%f and %f,%f", x1, y1, x2, y2, x3, y3)
		return
	}

	cx := (mr*mt*(y3-y1) + mr*(x2+x3) - mt*(x1+x2)) / (2 * (mr - mt))
	cy := (y1+y2)/2 - (cx-(x1+x2)/2)/mr
	cr := math.Pow((math.Pow((x2-cx), 2) + math.Pow((y2-cy), 2)), 0.5)

	x = int(cx + 0.5)
	y = int(cy + 0.5)
	r = int(cr + 0.5)

	return
}

// addCircleEvent adds a circle event for arc2 to the queue if the bottom point of
// the circle is below the sweep line.
func (v *Voronoi) addCircleEvent(arc1, arc2, arc3 *Node) {
	if arc1 == nil || arc2 == nil || arc3 == nil {
		return
	}

	log.Printf("Checking for circle at %v %v %v\r\n", arc1, arc2, arc3)
	x, y, r, err := v.calcCircle(arc1.Site, arc2.Site, arc3.Site)
	if err != nil {
		return
	}

	// Only add events with bottom point below the sweep line
	bottomY := y + r
	if bottomY <= v.SweepLine {
		return
	}

	event := &Event{
		EventType: EventCircle,
		X:         x,
		Y:         bottomY,
		Radius:    r,
	}
	v.EventQueue.Push(event)

	arc1.AddEvent(event)
	arc2.AddEvent(event)
	arc3.AddEvent(event)
	event.Node = arc2

	log.Printf("Added circle with center %d,%d, r=%d and bottom Y=%d\r\n", x, y, r, bottomY)
}

func (v *Voronoi) handleCircleEvent(event *Event) {
	log.Println()
	log.Printf("Handling circle event %d,%d with radius %d\r\n", event.X, event.Y, event.Radius)
	log.Printf("Sweep line: %d", v.SweepLine)
	log.Printf("Tree: %v", v.ParabolaTree)

	log.Printf("Node to be removed: %v", event.Node)

	// Add center of circle as vertex
	vertex := v.DCEL.NewVertex(event.X, event.Y-event.Radius)
	log.Printf("Vertex at %d,%d (center of circle)\r\n", vertex.X, vertex.Y)

	// Finish edges for the node that is about to be removed
	v.CloseTwins(event.Node.LeftEdges, vertex)
	v.CloseTwins(event.Node.RightEdges, vertex)

	// Delete the arc for event.Node from the tree
	prevArc := event.Node.PrevArc()
	nextArc := event.Node.NextArc()
	log.Printf("Removing arc %v between %v and %v", event.Node, prevArc, nextArc)
	log.Printf("Previous arc: %v", prevArc)
	log.Printf("Next arc: %v", nextArc)
	// TODO: Remove any common half-edges between event.Node and prevArc or nextArc
	v.removeArc(event.Node)

	// Remove circle events
	v.removeCircleEvent(prevArc)
	v.removeCircleEvent(event.Node)
	v.removeCircleEvent(nextArc)

	// Check for new circle events where the former left arc is the middle
	prevPrevArc := prevArc.PrevArc()
	v.addCircleEvent(prevPrevArc, prevArc, nextArc)

	// Check for new circle events where the former right arc is the middle
	nextNextArc := nextArc.NextArc()
	v.addCircleEvent(prevArc, nextArc, nextNextArc)

	// Finish edges for the neighbouring edges.
	v.CloseTwins(prevArc.RightEdges, vertex)
	v.CloseTwins(nextArc.LeftEdges, vertex)

	// Create a new edge in DCEL with this vertex as a target.
	// Attach the half edges to the corresponding arc.
	edge1, edge2 := v.DCEL.NewEdge(prevArc.Site.Face, nextArc.Site.Face, vertex)
	prevArc.RightEdges = append(prevArc.RightEdges, edge1)
	nextArc.LeftEdges = append(nextArc.LeftEdges, edge2)

	// TODO: Check for other circles? Rebuild circle events with prevArc and nextArc.

	return
}

// removeArc removes the given arc leaf from the binary tree.
func (v *Voronoi) removeArc(node *Node) {
	// TODO: Consider the case of removing several arcs at once
	parent := node.Parent
	other := (*Node)(nil)
	if parent.Left == node {
		other = parent.Right
	} else {
		other = parent.Left
	}
	grandParent := parent.Parent
	if grandParent == nil {
		v.ParabolaTree = other
		v.ParabolaTree.Parent = nil
		return
	}

	if grandParent.Left == parent {
		grandParent.Left = other
		grandParent.Left.Parent = grandParent
	} else if grandParent.Right == parent {
		grandParent.Right = other
		grandParent.Right.Parent = grandParent
	}
}

// removeCircleEvent remove the circle event where the specified node represents the middle arc.
func (v *Voronoi) removeCircleEvent(node *Node) {
	if node == nil {
		return
	}

	if len(node.Events) > 0 {
		log.Printf("Removing %d events from queue for arc %v.\r\n", len(node.Events), node.Site)

		for _, e := range node.Events {
			if e.index > -1 {
				v.EventQueue.Remove(e)
			}
		}
		node.Events = nil
	}
}
