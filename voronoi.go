package goalgorithms

import (
	"container/heap"
	"image"
	"log"

	"github.com/quasoft/btree"
)

// RVertex repsesent a vertex on the resulting voronoi diagram.
// Will be replaced with a double-connected edge list structure in a later
// version of the library.
type RVertex struct {
	X int
	Y int
}

type Voronoi struct {
	Bounds       image.Rectangle
	Sites        SiteSlice
	EventQueue   EventQueue
	ParabolaTree *btree.Node
	SweepLine    int // tracks the current position of the sweep line; updated when a new site is added.
	Result       []RVertex
}

func NewVoronoi(sites SiteSlice, bounds image.Rectangle) *Voronoi {
	voronoi := &Voronoi{Bounds: bounds}
	voronoi.Sites = make(SiteSlice, len(sites), len(sites))
	copy(voronoi.Sites, sites)
	voronoi.init()
	return voronoi
}

func NewFromPoints(points []image.Point, bounds image.Rectangle) *Voronoi {
	var sites SiteSlice
	for _, point := range points {
		sites = append(sites, Site{point.X, point.Y})
	}
	return NewVoronoi(sites, bounds)
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
	// TODO: Create DCEL list
}

func (v *Voronoi) Reset() {
	v.EventQueue = NewEventQueue(v.Sites)
	v.ParabolaTree = nil
	v.Result = make([]RVertex, 0)
	v.SweepLine = 0
}

func (v *Voronoi) HandleNextEvent() {
	if v.EventQueue.Len() > 0 {
		// Process events by Y (priority)
		event := heap.Pop(&v.EventQueue).(*Event)
		v.SweepLine = event.site.Y
		if event.EventType == EventSite {
			v.handleSiteEvent(event)
		} else {
			v.handleCircleEvent(event)
		}
	}
}

func (v *Voronoi) Generate() {
	v.Reset()

	// While queue is not empty
	for v.EventQueue.Len() > 0 {
		v.HandleNextEvent()
	}
}

// findNodeAbove finds the node for the parabola that is vertically above the specified site.
func (v *Voronoi) findNodeAbove(site Site) *btree.Node {
	node := v.ParabolaTree

	for !node.IsLeaf() {
		x := GetXOfIntersection(node, v.SweepLine)
		if site.X < x {
			node = node.Left
		} else {
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
	log.Printf("Handling event %d:%d of type %d\r\n", event.site.X, event.site.Y, event.EventType)
	log.Printf("Sweep line: %d", v.SweepLine)
	log.Printf("Tree: %v", v.ParabolaTree)

	// Event with Y above the sweep line should be ignored
	if event.site.Y < v.SweepLine {
		log.Printf("Ignoring event as it's above the sweep line (%d)\r\n", v.SweepLine)
		return
	}

	// If the binary tree is empty, just add an arc for this site as the only leaf in the tree
	if v.ParabolaTree == nil {
		log.Print("Adding event as root\r\n")
		v.ParabolaTree = NewArcNode(event)
		return
	}

	// If the tree is not empty, find the arc vertically above the new site
	arcNode := v.findNodeAbove(event.site)
	if arcNode == nil {
		log.Print("Could not find arc above event site!\r\n")
		// Do something
		return
	}

	arc := arcNode.Value.(*Arc)
	log.Printf("Arc above: %d:%d\r\n", arc.Site.X, arc.Site.Y)

	if len(arc.Events) > 0 {
		log.Printf("Removing %d events from queue.\r\n", len(arc.Events))

		// Remove false circle events from queue
		for _, e := range arc.Events {
			v.EventQueue.Remove(e)
		}
		arc.Events = nil
	}

	y := GetYByX(arc.Site, event.site.X, v.SweepLine)
	point := RVertex{event.site.X, y}
	log.Printf("Y of intersection = %d:%d\r\n", point.X, point.Y)
	v.Result = append(v.Result, point)

	// The node above (NA) is replaced wit ha branch with one internal node and three leafs.
	// The middle leaf stores the new parabola and the other two store the one being split.
	//    (NA)
	//   /   \
	//  (  )  [old]
	// /    \
	//[old]  [new]
	arcNode.Right = btree.New(arc)                         // Copy of the old arc
	arcNode.Left = btree.New(&Arc{})                       // Internal node
	arcNode.Left.Left = btree.New(arc)                     // Copy of the old arc
	arcNode.Left.Right = btree.New(&Arc{Site: event.site}) // The new arc
	arcNode.Value = &Arc{}                                 // Remove the value of the internal node
}

func (v *Voronoi) handleCircleEvent(event *Event) {
	return
}
