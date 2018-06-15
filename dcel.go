package voronoi

import (
	"fmt"
	"math"
	"sort"

	"github.com/quasoft/dcel"
)

// CloseTwins adds a vertex to the specified edges.
func (v *Voronoi) CloseTwins(list []*dcel.HalfEdge, vertex *dcel.Vertex) {
	for i := 0; i < len(list); i++ {
		he := list[i]
		if he.Twin != nil && he.Twin.Target == nil {
			he.Twin.Target = vertex
		} else if he.Target == nil {
			he.Target = vertex
		}
	}
}

// halfEdgesByCCW implements a slice of half-edges that sort in counter-clockwise order.
type halfEdgesByCCW []*dcel.HalfEdge

func (s halfEdgesByCCW) Len() int {
	return len(s)
}
func (s halfEdgesByCCW) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s halfEdgesByCCW) Less(i, j int) bool {
	if s[i].Target == nil || s[j].Target == nil {
		return false
	}

	// Find center of polygon
	var sumX int64
	var sumY int64
	var cnt int
	for _, v := range s {
		if v.Target != nil {
			sumX += int64(v.Target.X)
			sumY += int64(v.Target.Y)
			cnt++
		}
	}
	centerX := float64(sumX) / float64(cnt)
	centerY := float64(sumY) / float64(cnt)

	// Sort counter-clockwise
	a1 := math.Atan2(float64(s[i].Target.Y)-centerY, float64(s[i].Target.X)-centerX)
	a2 := math.Atan2(float64(s[j].Target.Y)-centerY, float64(s[j].Target.X)-centerX)
	return a1 >= a2
}

// GetFaceHalfEdges returns the half-edges that form the boundary of a face (cell).
func (v *Voronoi) GetFaceHalfEdges(face *dcel.Face) []*dcel.HalfEdge {
	var edges []*dcel.HalfEdge
	exists := make(map[string]bool)
	for _, edge := range v.DCEL.HalfEdges {
		if edge.Face.ID == face.ID {
			id := fmt.Sprintf("%v", edge.Target)
			if !exists[id] {
				exists[id] = true
				edges = append(edges, edge)
			}
		}
	}

	sort.Sort(halfEdgesByCCW(edges))
	return edges
}

// verticesByCCW implements a slice of vertices that sort in counter-clockwise order.
type verticesByCCW []*dcel.Vertex

func (s verticesByCCW) Len() int {
	return len(s)
}
func (s verticesByCCW) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s verticesByCCW) Less(i, j int) bool {
	// Find center of polygon
	var sumX float64
	var sumY float64
	for _, v := range s {
		sumX += float64(v.X)
		sumY += float64(v.Y)
	}
	centerX := sumX / float64(len(s))
	centerY := sumY / float64(len(s))

	// Sort counter-clockwise
	a1 := math.Atan2(float64(s[i].Y)-centerY, float64(s[i].X)-centerX)
	a2 := math.Atan2(float64(s[j].Y)-centerY, float64(s[j].X)-centerX)
	return a1 >= a2
}

// GetFaceVertices returns the vertices that form the boundary of a face (cell),
// sorted in counter-clockwise order.
func (v *Voronoi) GetFaceVertices(face *dcel.Face) []*dcel.Vertex {
	var vertices []*dcel.Vertex
	exists := make(map[string]bool)
	for _, edge := range v.DCEL.HalfEdges {
		if edge.Face.ID == face.ID {
			if edge.Target != nil {
				id := fmt.Sprintf("%v", edge.Target)
				if !exists[id] {
					exists[id] = true
					vertices = append(vertices, edge.Target)
				}
			}

			if edge.Twin != nil && edge.Twin.Target != nil {
				id := fmt.Sprintf("%v", edge.Twin.Target)
				if !exists[id] {
					exists[id] = true
					vertices = append(vertices, edge.Twin.Target)
				}
			}
		}
	}

	sort.Sort(verticesByCCW(vertices))
	return vertices
}
