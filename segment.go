package banana

import (
	"image/color"
	"math"

	"github.com/dfirebaugh/banana/graphics"
)

type Segment struct {
	X1, Y1, X2, Y2, Width float32
	Color                 color.Color
	StencilWriteValue     uint8
	StencilTestValue      uint8
}

func (l *Segment) GetVertices(screenWidth, screenHeight int) []graphics.Vertex {
	color := colorToVec(l.Color)

	angle := float32(math.Atan2(float64(l.Y2-l.Y1), float64(l.X2-l.X1)))

	halfWidth := l.Width / 2.0

	cosAngle := float32(math.Cos(float64(angle)))
	sinAngle := float32(math.Sin(float64(angle)))

	dx := halfWidth * sinAngle
	dy := halfWidth * cosAngle

	x1Left := l.X1 - dx
	y1Left := l.Y1 + dy
	x1Right := l.X1 + dx
	y1Right := l.Y1 - dy

	x2Left := l.X2 - dx
	y2Left := l.Y2 + dy
	x2Right := l.X2 + dx
	y2Right := l.Y2 - dy

	normX1Left, normY1Left := normalizeCoordinates(x1Left, y1Left, screenWidth, screenHeight)
	normX1Right, normY1Right := normalizeCoordinates(x1Right, y1Right, screenWidth, screenHeight)
	normX2Left, normY2Left := normalizeCoordinates(x2Left, y2Left, screenWidth, screenHeight)
	normX2Right, normY2Right := normalizeCoordinates(x2Right, y2Right, screenWidth, screenHeight)

	vertices := []graphics.Vertex{
		{FsQuadPos: [2]float32{normX1Left, normY1Left}, ShapePos: [2]float32{normX1Left, normY1Left}, OpCode: graphics.OP_CODE_VERTEX, Color: color, Resolution: [2]float32{float32(screenWidth), float32(screenHeight)}},
		{FsQuadPos: [2]float32{normX2Left, normY2Left}, ShapePos: [2]float32{normX2Left, normY2Left}, OpCode: graphics.OP_CODE_VERTEX, Color: color, Resolution: [2]float32{float32(screenWidth), float32(screenHeight)}},
		{FsQuadPos: [2]float32{normX2Right, normY2Right}, ShapePos: [2]float32{normX2Right, normY2Right}, OpCode: graphics.OP_CODE_VERTEX, Color: color, Resolution: [2]float32{float32(screenWidth), float32(screenHeight)}},
		{FsQuadPos: [2]float32{normX1Left, normY1Left}, ShapePos: [2]float32{normX1Left, normY1Left}, OpCode: graphics.OP_CODE_VERTEX, Color: color, Resolution: [2]float32{float32(screenWidth), float32(screenHeight)}},
		{FsQuadPos: [2]float32{normX2Right, normY2Right}, ShapePos: [2]float32{normX2Right, normY2Right}, OpCode: graphics.OP_CODE_VERTEX, Color: color, Resolution: [2]float32{float32(screenWidth), float32(screenHeight)}},
		{FsQuadPos: [2]float32{normX1Right, normY1Right}, ShapePos: [2]float32{normX1Right, normY1Right}, OpCode: graphics.OP_CODE_VERTEX, Color: color, Resolution: [2]float32{float32(screenWidth), float32(screenHeight)}},
	}

	return vertices
}

func (l *Segment) GetStencilTestValue() uint8 {
	return l.StencilTestValue
}

func (l *Segment) GetStencilWriteValue() uint8 {
	return l.StencilWriteValue
}
