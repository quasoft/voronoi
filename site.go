package voronoi

import "fmt"

// Site is a prerequisute for computing a voronoi diagram.
// Site is the point (also called seed or generator) in a voronoi diagram,
// around which a cell (subset of the plane) is formed, with such a property
// that every point in the cell is closer to this site than any other site.
type Site struct {
	X, Y int
}

func (s Site) String() string { return fmt.Sprintf("%d,%d", s.X, s.Y) }

// SiteSlice is a slice of Site values, sortable by Y
type SiteSlice []Site

func (s SiteSlice) Len() int           { return len(s) }
func (s SiteSlice) Less(i, j int) bool { return s[i].Y < s[j].Y }
func (s SiteSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
