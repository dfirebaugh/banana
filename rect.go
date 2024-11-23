package banana

import (
	"image/color"

	"github.com/dfirebaugh/banana/graphics"
)

type Rect struct {
	X, Y, Width, Height, Radius float32
	Color                       color.Color
}

func (r *Rect) GetVertices(screenWidth, screenHeight int) []graphics.Vertex {
	red, g, b, a := r.Color.RGBA()
	color := [4]float32{
		float32(red) / 65535.0,
		float32(g) / 65535.0,
		float32(b) / 65535.0,
		float32(a) / 65535.0,
	}
	normX, normY := normalizeCoordinates(r.X, r.Y, screenWidth, screenHeight)

	halfWidth := r.Width * 0.5
	halfHeight := r.Height * 0.5

	vertices := []float32{
		-halfWidth, halfHeight,
		-halfWidth, -halfHeight,
		halfWidth, -halfHeight,
		-halfWidth, halfHeight,
		halfWidth, -halfHeight,
		halfWidth, halfHeight,
	}

	var result []graphics.Vertex
	for i := 0; i < 6; i++ {
		if r.Radius == 0 {
			r.Radius = 1
		}
		v := graphics.Vertex{
			FsQuadPos:  [2]float32{vertices[i*2], vertices[i*2+1]},
			ShapePos:   [2]float32{normX + halfWidth/float32(screenWidth)*2.0, normY - halfHeight/float32(screenHeight)*2.0},
			LocalPos:   [2]float32{vertices[i*2], vertices[i*2+1]},
			OpCode:     graphics.OP_CODE_RECT,
			Radius:     r.Radius,
			Width:      r.Width,
			Height:     r.Height,
			Color:      color,
			Resolution: [2]float32{float32(screenWidth), float32(screenHeight)},
		}
		result = append(result, v)
	}
	return result
}

func (r *Rect) IsWithinBounds(px, py float32) bool {
	left := r.X - r.Width/2
	right := r.X + r.Width/2
	top := r.Y - r.Height/2
	bottom := r.Y + r.Height/2

	return px >= left && px <= right && py >= top && py <= bottom
}
