package voronoi

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

// Plotter draws the result of the voronoi diagram generator into an image.
type Plotter struct {
	voronoi         *Voronoi
	dst             *image.RGBA
	ctx             *draw.Context
	BackgroundColor color.RGBA
	VertexColor     color.RGBA
}

// NewPlotter creates a new voronoi diagram drawer.
func NewPlotter(voronoi *Voronoi, dst *image.RGBA) *Plotter {
	return &Plotter{
		voronoi,
		dst,
		draw.NewContext(dst),
		color.RGBA{255, 255, 255, 255}, // White
		color.RGBA{0, 0, 0, 255},       // Blue
	}
}

// Min returns the minimum point on the diagram.
func (p *Plotter) Min() image.Point {
	return p.dst.Bounds().Min
}

// Max returns the maximum point on the diagram.
func (p *Plotter) Max() image.Point {
	return p.dst.Bounds().Max
}

// SweepLine draws a sweep line with the given Y.
func (p *Plotter) SweepLine(y int) {
	p.ctx.SetPen(color.Black)
	p.ctx.Line(0, y, p.Max().X-1, y)
}

// Site draws the specified site with the given color.
func (p *Plotter) Site(site Site, clr color.Color) {
	p.ctx.SetPen(clr)

	p.ctx.Cross(site.X, site.Y, 2)
}

// Vertex draws the specified vertex.
func (p *Plotter) Vertex(vertex RVertex) {
	p.ctx.SetPen(p.VertexColor)
	p.ctx.Cross(vertex.X, vertex.Y, 2)
}

func (p *Plotter) colorOfSite(site Site) color.Color {
	siteIdx := 0
	for i, s := range p.voronoi.Sites {
		if site == s {
			siteIdx = i
			break
		}
	}
	return p.colorOfSiteIdx(siteIdx)
}

func (p *Plotter) colorOfSiteIdx(index int) color.Color {
	return colors[index%len(colors)]
}

// BeachLine draws the sequence of parabola arcs.
func (p *Plotter) BeachLine(tree *Node) {
	// Draw full parabolas with semi-transparent color
	first := tree.FirstArc()
	lastX := 0
	for first != nil {
		// Get parabola coefficients
		a, b, c := GetParabolaABC(first.Site, p.voronoi.SweepLine)

		cr, cg, cb, _ := p.colorOfSite(first.Site).RGBA()
		stclr := color.RGBA{uint8(cr), uint8(cg), uint8(cb), 75}
		p.ctx.SetPen(stclr)
		p.ctx.Parabola(a, b, c)

		first = first.NextArc()
	}

	// Draw parabola arcs with solid color
	first = tree.FirstArc()
	if first != nil {
		fmt.Printf("First: %v\r\n", first.Site)
	} else {
		fmt.Print("First: nil\r\n")
	}
	lastX = 0
	for first != nil {
		// Get parabola coefficients
		a, b, c := GetParabolaABC(first.Site, p.voronoi.SweepLine)

		clr := p.colorOfSite(first.Site)
		p.ctx.SetPen(clr)

		x := p.Max().X
		next := first.NextArc()
		if next != nil {
			intX, err := GetXOfIntersection(first, next, p.voronoi.SweepLine)
			if err == nil {
				x = intX
			}
		}

		if first.Site.Y == p.voronoi.SweepLine {
			p.ctx.Line(first.Site.X, 0, first.Site.X, first.Site.Y)
		} else {
			p.ctx.ParabolaArc(a, b, c, lastX, x)
		}
		lastX = x

		if next != nil {
			fmt.Printf("Next: %v\r\n", next.Site)
		} else {
			fmt.Print("Next: nil\r\n")
		}

		first = next
	}
	fmt.Print("Over\r\n")
}

// Plot paints the voronoi diagram over the given image.
func (p *Plotter) Plot() {
	// Draw border and fill with background color
	p.ctx.SetPen(color.Black)
	p.ctx.SetFill(color.White)
	p.ctx.Rect(0, 0, p.Max().X-1, p.Max().Y-1)

	// Draw beach line
	p.BeachLine(p.voronoi.ParabolaTree)

	// Draw sites and their labels
	for i, site := range p.voronoi.Sites {
		clr := p.colorOfSiteIdx(i)

		p.Site(site, clr)
		label := fmt.Sprintf("Site %d/%d", site.X, site.Y)
		p.ctx.SetTextColor(clr)
		p.ctx.Text(site.X-40, site.Y+15, label)
	}

	// Draw verteces
	for _, vertex := range p.voronoi.Result {
		p.Vertex(vertex)
	}

	// Draw sweep line with label
	p.SweepLine(p.voronoi.SweepLine)
	label := fmt.Sprintf("Sweep line = %d", p.voronoi.SweepLine)
	p.ctx.SetTextColor(color.Black)
	p.ctx.Text(p.Max().X-150, p.voronoi.SweepLine+15, label)
}

// Plot creates an image and paints a voronoi diagram over it.
func Plot(voronoi *Voronoi) *image.RGBA {
	img := image.NewRGBA(voronoi.Bounds)
	drawer := NewPlotter(voronoi, img)
	drawer.Plot()
	return img
}
