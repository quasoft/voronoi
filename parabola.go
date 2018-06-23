package voronoi

import (
	"fmt"
	"log"
	"math"
)

// GetParabolaABC returns the a, b and c coefficients of the standard form of
// a parabola equation, given only x and y of the focus and y of the directrix.
// Math behind this is explained at https://math.stackexchange.com/q/2700061/543428.
func GetParabolaABC(focus *Site, yOfDirectrix int) (float64, float64, float64) {
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

// GetXOfInternalNode returns the x of the intersection of the two parabola arcs below an internal node.
func GetXOfInternalNode(node *Node, directrix int) (int, error) {
	left := node.PrevChildArc()
	right := node.NextChildArc()

	return GetXOfIntersection(left, right, directrix)
}

// GetXOfIntersection returns the x of the intersection of two parabola arcs.
func GetXOfIntersection(left *Node, right *Node, directrix int) (int, error) {
	leftFocus := left.Site
	rightFocus := right.Site

	// If two parabolas have the same Y, then the intersection lies exactly at the
	// middle between them.
	if leftFocus.Y == rightFocus.Y {
		return (leftFocus.X + rightFocus.X) / 2, nil
	}

	// Handle the degenerate case where one or both of the sites have the same Y value.
	// In this case the focus of one or both sites and the directrix would be equal.
	if leftFocus.Y == directrix {
		return leftFocus.X, nil
	} else if rightFocus.Y == directrix {
		return rightFocus.X, nil
	}

	// Determine the a, b and c coefficients for the two parabolas
	a1, b1, c1 := GetParabolaABC(leftFocus, directrix)
	a2, b2, c2 := GetParabolaABC(rightFocus, directrix)

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

	log.Printf("X of %v and %v = %v\r\n", leftFocus, rightFocus, x)
	if math.IsNaN(x) {
		return 0, fmt.Errorf("there is no intersection between S(%v) and S(%v)", leftFocus, rightFocus)
	}

	return int(x), nil
}

// GetYByX calculates the Y value for the parabola with the given focus and directrix (the sweep line)
func GetYByX(focus *Site, x int, directrix int) int {
	xf := float64(x)
	a, b, c := GetParabolaABC(focus, directrix)
	y := a*math.Pow(xf, 2) + b*xf + c

	if math.IsNaN(y) {
		y = 0
	}

	log.Printf("focus: %v:%v, x=%v, directrix=%v \r\n", focus.X, focus.Y, xf, directrix)
	log.Printf("a=%v, b=%v, c=%v, y=%v \r\n", a, b, c, y)
	return int(y)
}
