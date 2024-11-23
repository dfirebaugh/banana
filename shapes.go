package banana

import (
	"image/color"

	"github.com/dfirebaugh/banana/graphics"
)

type TextRenderOptions struct {
	X, Y, Size float32
	Color      color.Color
}

func RenderText(text string, options *TextRenderOptions) {
	banana.graphicsBackend.RenderText(text, &graphics.TextRenderOptions{
		X:     options.X,
		Y:     options.Y,
		Size:  options.Size,
		Color: options.Color,
	})
}

type Renderable interface {
	GetVertices(screenWidth, screenHeight int) []graphics.Vertex
}

func RenderShape(shape Renderable) {
	banana.graphicsBackend.Render(graphics.Renderable(shape))
}

func normalizeCoordinates(x, y float32, screenWidth, screenHeight int) (float32, float32) {
	normX := (x/float32(screenWidth))*2.0 - 1.0
	normY := 1.0 - (y/float32(screenHeight))*2.0
	return normX, normY
}
