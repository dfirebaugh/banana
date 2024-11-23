package banana

import (
	"image/color"

	"github.com/dfirebaugh/banana/graphics"
)

type Circle struct {
	X, Y, Radius float32
	Color        color.Color
}

func (c *Circle) GetVertices(screenWidth, screenHeight int) []graphics.Vertex {
	normX, normY := normalizeCoordinates(c.X, c.Y, screenWidth, screenHeight)

	r, g, b, a := c.Color.RGBA()
	color := [4]float32{
		float32(r) / 65535.0,
		float32(g) / 65535.0,
		float32(b) / 65535.0,
		float32(a) / 65535.0,
	}
	vertices := []float32{
		-c.Radius, c.Radius,
		-c.Radius, -c.Radius,
		c.Radius, -c.Radius,
		-c.Radius, c.Radius,
		c.Radius, -c.Radius,
		c.Radius, c.Radius,
	}

	var result []graphics.Vertex
	for i := 0; i < 6; i++ {
		v := graphics.Vertex{
			FsQuadPos:  [2]float32{vertices[i*2], vertices[i*2+1]},
			ShapePos:   [2]float32{normX, normY},
			LocalPos:   [2]float32{vertices[i*2], vertices[i*2+1]},
			OpCode:     graphics.OP_CODE_CIRCLE,
			Radius:     c.Radius,
			Width:      c.Radius * 2.0,
			Height:     c.Radius * 2.0,
			Color:      color,
			Resolution: [2]float32{float32(screenWidth), float32(screenHeight)},
		}
		result = append(result, v)
	}
	return result
}
