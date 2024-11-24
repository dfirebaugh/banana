package gui

import (
	"github.com/dfirebaugh/banana"
)

const (
	drawOpCircle = iota
	drawOpRoundedRect
	drawOpTriangle
)

func (d *Draw) drawTriangle(x1, y1, x2, y2, x3, y3 int, op *DrawOptions) {
	banana.RenderShape(&banana.Polygon{
		Vertices: []banana.Vertex{
			{
				X:     float32(x1),
				Y:     float32(y1),
				Color: op.FillColor,
			},
			{
				X:     float32(x2),
				Y:     float32(y2),
				Color: op.FillColor,
			},
			{
				X:     float32(x3),
				Y:     float32(y3),
				Color: op.FillColor,
			},
		},
	})
}

func (d *Draw) drawRectangle(x, y, width, height int, op *DrawOptions) {
	banana.RenderShape(&banana.Rect{
		X:      float32(x),
		Y:      float32(y),
		Width:  float32(width),
		Height: float32(height),
		Radius: 1,
		Color:  op.FillColor,
	})
}

func (d *Draw) drawRoundedRectangle(x, y, width, height, radius int, op *DrawOptions) {
	banana.RenderShape(&banana.Rect{
		X:      float32(x),
		Y:      float32(y),
		Width:  float32(width),
		Height: float32(height),
		Radius: float32(radius),
		Color:  op.FillColor,
	})
}

func (d *Draw) drawCircle(x, y, radius int, op *DrawOptions) {
	banana.RenderShape(&banana.Circle{
		X:      float32(x),
		Y:      float32(y),
		Radius: float32(radius),
		Color:  op.FillColor,
	})
}

func (d *Draw) drawRoundedRectangleWithOutline(x, y, width, height, radius, outlineWidth int, op *DrawOptions) {
	outerX := x - outlineWidth
	outerY := y - outlineWidth
	outerWidth := width + 2*outlineWidth
	outerHeight := height + 2*outlineWidth

	banana.RenderShape(&banana.Rect{
		X:      float32(outerX),
		Y:      float32(outerY),
		Width:  float32(outerWidth),
		Height: float32(outerHeight),
		Radius: float32(radius + outlineWidth),
		Color:  op.OutlineColor,
	})

	banana.RenderShape(&banana.Rect{
		X:      float32(x),
		Y:      float32(y),
		Width:  float32(width),
		Height: float32(height),
		Radius: float32(radius),
		Color:  op.FillColor,
	})
}

func (d *Draw) drawCircleWithOutline(x, y, radius, outlineWidth int, op *DrawOptions) {
	banana.RenderShape(&banana.Circle{
		X:      float32(x),
		Y:      float32(y),
		Radius: float32(radius + outlineWidth),
		Color:  op.OutlineColor,
	})

	banana.RenderShape(&banana.Circle{
		X:      float32(x),
		Y:      float32(y),
		Radius: float32(radius),
		Color:  op.FillColor,
	})
}

func generateQuadraticBezierVertices(start, control, end Position) []Position {
	const steps = 100
	vertices := make([]Position, steps+1)

	for i := 0; i <= steps; i++ {
		t := float32(i) / float32(steps)
		x := (1-t)*(1-t)*float32(start.X) + 2*(1-t)*t*float32(control.X) + t*t*float32(end.X)
		y := (1-t)*(1-t)*float32(start.Y) + 2*(1-t)*t*float32(control.Y) + t*t*float32(end.Y)
		vertices[i] = Position{X: int(x), Y: int(y)}
	}

	return vertices
}

func generateWaveVertices(waveFunc func(x float64) float64, amplitude, frequency, phase float64, startX, endX, y int) []Position {
	const steps = 1000
	vertices := make([]Position, steps+1)
	stepSize := float64(endX-startX) / float64(steps)

	for i := 0; i <= steps; i++ {
		x := float64(startX) + stepSize*float64(i)
		yOffset := amplitude * waveFunc(frequency*x+phase)
		vertices[i] = Position{X: int(x), Y: y + int(yOffset)}
	}

	return vertices
}

func drawWave(vertices []Position, d *Draw, op *DrawOptions) {
	for i := 0; i < len(vertices)-1; i++ {
		d.DrawSegment(vertices[i].X, vertices[i].Y, vertices[i+1].X, vertices[i+1].Y, op)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d *Draw) drawText(text string, options *banana.TextRenderOptions) {
	// can do some theme stuff here
	banana.RenderText(text, options)
}
