package banana

import (
	"image/color"

	"github.com/dfirebaugh/banana/graphics"
)

type Vertex struct {
	X, Y  float32
	Color color.Color
}

type Polygon struct {
	Vertices []Vertex
}

func colorToVec(c color.Color) [4]float32 {
	if c == nil {
		return [4]float32{0, 0, 0, 0}
	}
	r, g, b, a := c.RGBA()
	return [4]float32{
		float32(r) / 65535.0,
		float32(g) / 65535.0,
		float32(b) / 65535.0,
		float32(a) / 65535.0,
	}
}

func (t *Polygon) GetVertices(screenWidth, screenHeight int) []graphics.Vertex {
	if len(t.Vertices) == 0 {
		return []graphics.Vertex{}
	}

	result := make([]graphics.Vertex, 3)
	for i, v := range t.Vertices {
		normX := (float32(v.X)/float32(screenWidth))*2.0 - 1.0
		normY := 1.0 - (float32(v.Y)/float32(screenHeight))*2.0
		result[i] = graphics.Vertex{
			FsQuadPos:  [2]float32{normX, normY},
			ShapePos:   [2]float32{normX, normY},
			OpCode:     graphics.OP_CODE_VERTEX,
			Color:      colorToVec(v.Color),
			Resolution: [2]float32{float32(screenWidth), float32(screenHeight)},
		}
	}

	return result
}
