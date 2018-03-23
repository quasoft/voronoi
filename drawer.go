package goalgorithms

import (
	"fmt"
	"image"
	"image/color"

	"github.com/quasoft/draw"
)

var colors = []color.Color{
	color.RGBA{0xff, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0xff, 0xff},
	color.RGBA{0x00, 0xff, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x8b, 0xff},
	color.RGBA{0x00, 0x8b, 0x8b, 0xff},
	color.RGBA{0xb8, 0x86, 0x0b, 0xff},
	color.RGBA{0x00, 0x64, 0x00, 0xff},
	color.RGBA{0xbd, 0xb7, 0x6b, 0xff},
	color.RGBA{0x8b, 0x00, 0x8b, 0xff},
	color.RGBA{0x55, 0x6b, 0x2f, 0xff},
	color.RGBA{0xff, 0x8c, 0x00, 0xff},
	color.RGBA{0x99, 0x32, 0xcc, 0xff},
	color.RGBA{0x8b, 0x00, 0x00, 0xff},
	color.RGBA{0xe9, 0x96, 0x7a, 0xff},
	color.RGBA{0x8f, 0xbc, 0x8f, 0xff},
	color.RGBA{0x48, 0x3d, 0x8b, 0xff},
	color.RGBA{0x2f, 0x4f, 0x4f, 0xff},
	color.RGBA{0x2f, 0x4f, 0x4f, 0xff},
	color.RGBA{0x00, 0xce, 0xd1, 0xff},
	color.RGBA{0x94, 0x00, 0xd3, 0xff},
	color.RGBA{0xff, 0x14, 0x93, 0xff},
	color.RGBA{0x00, 0xbf, 0xff, 0xff},
	color.RGBA{0x00, 0xff, 0x00, 0xff},
	color.RGBA{0xff, 0x00, 0xff, 0xff},
	color.RGBA{0x80, 0x00, 0x00, 0xff},
	color.RGBA{0x66, 0xcd, 0xaa, 0xff},
	color.RGBA{0x00, 0x00, 0xcd, 0xff},
	color.RGBA{0xba, 0x55, 0xd3, 0xff},
	color.RGBA{0x93, 0x70, 0xdb, 0xff},
	color.RGBA{0x3c, 0xb3, 0x71, 0xff},
	color.RGBA{0x7b, 0x68, 0xee, 0xff},
	color.RGBA{0x00, 0xfa, 0x9a, 0xff},
	color.RGBA{0x48, 0xd1, 0xcc, 0xff},
	color.RGBA{0xc7, 0x15, 0x85, 0xff},
	color.RGBA{0x00, 0x80, 0x80, 0xff},
	color.RGBA{0x40, 0xe0, 0xd0, 0xff},
	color.RGBA{0xee, 0x82, 0xee, 0xff},
}

// VoronoiDrawing draws the result of the voronoi diagram generator into an image.
type VoronoiDrawing struct {
	voronoi         *Voronoi
	dst             *image.RGBA
	ctx             *draw.Context
	BackgroundColor color.RGBA
	VertexColor     color.RGBA
}

// NewVoronoiDrawing creates a new voronoi diagram drawer.
func NewVoronoiDrawing(voronoi *Voronoi, dst *image.RGBA) *VoronoiDrawing {
	return &VoronoiDrawing{
		voronoi,
		dst,
		draw.NewContext(dst),
		color.RGBA{255, 255, 255, 255}, // White
		color.RGBA{0, 0, 0, 255},       // Blue
	}
}

// Min returns the minimum point on the diagram.
func (d *VoronoiDrawing) Min() image.Point {
	return d.dst.Bounds().Min
}

// Max returns the maximum point on the diagram.
func (d *VoronoiDrawing) Max() image.Point {
	return d.dst.Bounds().Max
}

// SweepLine draws a sweep line with the given Y.
func (d *VoronoiDrawing) SweepLine(y int) {
	d.ctx.SetPen(color.Black)
	d.ctx.Line(0, y, d.Max().X-1, y)
}

// Site draws the specified site with the given color.
func (d *VoronoiDrawing) Site(site Site, clr color.Color) {
	d.ctx.SetPen(clr)

	d.ctx.Cross(site.X, site.Y, 2)

	if site.Y == d.voronoi.SweepLine {
		d.ctx.Line(site.X, 0, site.X, site.Y)
	} else if site.Y < d.voronoi.SweepLine {
		a, b, c := GetParabolaABC(site, d.voronoi.SweepLine)
		d.ctx.Parabola(a, b, c)
	}
}

// Vertex draws the specified vertex.
func (d *VoronoiDrawing) Vertex(vertex RVertex) {
	d.ctx.SetPen(d.VertexColor)
	d.ctx.Cross(vertex.X, vertex.Y, 2)
}

// Plot paints the voronoi diagram over the given image.
func (d *VoronoiDrawing) Plot() {
	// Draw border and fill with background color
	d.ctx.SetPen(color.Black)
	d.ctx.SetFill(color.White)
	d.ctx.Rect(0, 0, d.Max().X-1, d.Max().Y-1)

	// Draw sites and their labels
	for i, site := range d.voronoi.Sites {
		clr := colors[i]

		d.Site(site, clr)
		label := fmt.Sprintf("Site %d/%d", site.X, site.Y)
		d.ctx.SetTextColor(clr)
		d.ctx.Text(site.X-40, site.Y+15, label)
	}

	// Draw verteces
	for _, vertex := range d.voronoi.Result {
		d.Vertex(vertex)
	}

	// Draw sweep line with label
	d.SweepLine(d.voronoi.SweepLine)
	label := fmt.Sprintf("Sweep line = %d", d.voronoi.SweepLine)
	d.ctx.SetTextColor(color.Black)
	d.ctx.Text(d.Max().X-150, d.voronoi.SweepLine+15, label)
}

// Plot creates an image and paints a voronoi diagram over it.
func Plot(voronoi *Voronoi) *image.RGBA {
	img := image.NewRGBA(voronoi.Bounds)
	drawer := NewVoronoiDrawing(voronoi, img)
	drawer.Plot()
	return img
}
