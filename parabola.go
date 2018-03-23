package goalgorithms

import (
	"fmt"
	"log"
	"math"

	"github.com/quasoft/btree"
)

// Arc represents an arc of a parabola (part of parabola), that lies on the beach line.
// Stores the site that created the arc, pointers to the circle events associated with it.
type Arc struct {
	// site is the focus of the parabola arc (the site which created this parabola)
	Site Site
	// events holds pointers to all circle events, in which this arc participates
	Events []*Event
}

func (a Arc) String() string {
	return fmt.Sprintf("%d:%d", a.Site.X, a.Site.Y)
}

// NewArc creates a new parabola arc for the given site
func NewArc(site Site) *Arc {
	return &Arc{Site: site}
}

// Less compares two nodes (a breakpoint or an arc) by the X value of the associated site.
// Not used in voronoi generator. Implemented just to fullfill ValueInterface.
func (a Arc) Less(value interface{}) bool {
	return a.Site.X < value.(Arc).Site.X
}

// NewArcNode creates a new tree node for the given site event.
func NewArcNode(event *Event) *btree.Node {
	arc := NewArc(event.site)
	return &btree.Node{Value: arc}
}

// GetParabolaABC returns the a, b and c coefficients of the standard form of
// a parabola equation, given only x and y of the focus and y of the directrix.
// Math behind this is explained at https://math.stackexchange.com/q/2700061/543428.
func GetParabolaABC(focus Site, yOfDirectrix int) (float64, float64, float64) {
	// a = 1 / 2(y_{f} - y_{d})
	// The formula for calculation of a coefficient is derived from the fact that
	// the distance (d) from the vertex of the parabola to its focus is 1/(4*a).
	// And the distance between the focus and the directrx is two times this distance,
	// so (y_{f} - y_{d}) = 1/(2a), which simplifies to the formula above.
	a := 1.0 / (2.0 * float64(focus.Y-yOfDirectrix))

	// b = -2ax_{f}
	// Calculation of b is based on the vertex form of the parabola equation: x_{0} = -b/(2a)
	b := -2.0 * a * float64(focus.X)

	// c = ax^2 + y_{f} - 1/(4a)
	// Formula for c is again derived from vertex form: c = ah^2 + k.
	// k is replace with (y_{f} - 1/(4а)), which is the distance between
	// y of vertex and y of focus.
	c := a*math.Pow(float64(focus.X), 2) + float64(focus.Y) - 1/(4*a)

	return a, b, c
}

// GetXOfIntersection returns the x of the intersection of two parabola arcs.
func GetXOfIntersection(node *btree.Node, sweepLine int) int {
	left := node.PrevLeaf()
	right := node.NextLeaf()

	leftFocus := left.Value.(*Arc).Site
	rightFocus := right.Value.(*Arc).Site

	// Determine the a, b and c coefficients for the two parabolas
	a1, b1, c1 := GetParabolaABC(leftFocus, sweepLine)
	a2, b2, c2 := GetParabolaABC(rightFocus, sweepLine)

	// Calculate the roots of the coefficients difference.
	a := a1 - a2
	b := b1 - b2
	c := c1 - c2

	discriminant := math.Pow(b, 2) - 4*a*c
	root1 := (-b + math.Sqrt(discriminant)) / (2 * a)
	root2 := (-b - math.Sqrt(discriminant)) / (2 * a)

	// X of the intersection is one of those roots.
	var x float64
	if leftFocus.Y < rightFocus.Y {
		x = math.Min(root1, root2)
	} else {
		x = math.Max(root1, root2)
	}

	log.Printf("X of S(%d,%d) and S(%d,%d) = %v\r\n", leftFocus.X, leftFocus.Y, rightFocus.X, rightFocus.Y, x)

	return int(x)
}

func GetXByY(focus Site, y int, sweepLine int) int {
	yf := float64(y)
	a, b, c := GetParabolaABC(focus, sweepLine)
	c -= yf

	discriminant := math.Pow(b, 2) - 4*a*c
	root1 := (-b + math.Sqrt(discriminant)) / (2 * a)
	root2 := (-b - math.Sqrt(discriminant)) / (2 * a)

	// X of the intersection is one of those roots.
	var x float64
	if focus.Y < focus.Y {
		x = math.Min(root1, root2)
	} else {
		x = math.Max(root1, root2)
	}

	if math.IsNaN(x) {
		x = float64(focus.X)
	}

	log.Printf("focus: %v:%v, y=%v, sweep line=%v \r\n", focus.X, focus.Y, yf, sweepLine)
	log.Printf("a=%v, b=%v, c=%v, x=%v \r\n", a, b, c, x)
	return int(x)
}

func GetYByX(focus Site, x int, sweepLine int) int {
	xf := float64(x)
	a, b, c := GetParabolaABC(focus, sweepLine)
	y := a*math.Pow(xf, 2) + b*xf + c

	if math.IsNaN(y) {
		y = 0
	}

	log.Printf("focus: %v:%v, x=%v, sweep line=%v \r\n", focus.X, focus.Y, xf, sweepLine)
	log.Printf("a=%v, b=%v, c=%v, y=%v \r\n", a, b, c, y)
	return int(y)
}
