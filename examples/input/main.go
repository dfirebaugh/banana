package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/pkg/input"
	"golang.org/x/image/colornames"
)

var colors = []color.Color{
	colornames.Red,
	colornames.Cadetblue,
	colornames.Black,
	colornames.Violet,
	colornames.Blue,
	colornames.Blueviolet,
	colornames.Orange,
}

func main() {
	banana.SetWindowSize(200, 200)
	var c color.Color = colornames.Red

	banana.Run(func() {
	}, func() {
		banana.Clear(colornames.Grey)
		if banana.IsKeyPressed(input.KeySpace) {
			c = colors[rand.Intn(len(colors))]
		}

		x, y := banana.GetCursorPosition()
		banana.RenderText(fmt.Sprintf("%d:%d", x, y), &banana.TextRenderOptions{
			X:     10,
			Y:     10,
			Size:  12,
			Color: colornames.Red,
		})

		banana.RenderShape(&banana.Polygon{
			Vertices: []banana.Vertex{
				{
					X:     0,
					Y:     200,
					Color: c,
				},
				{
					X:     100,
					Y:     0,
					Color: c,
				},
				{
					X:     200,
					Y:     200,
					Color: c,
				},
			},
		})
		banana.RenderShape(&banana.Circle{X: float32(x), Y: float32(y), Radius: 10, Color: colornames.Red})
		if banana.IsButtonPressed(input.MouseButtonLeft) {
			banana.RenderText("left mouse button pressed", &banana.TextRenderOptions{
				X:     10,
				Y:     25,
				Size:  12,
				Color: colornames.Red,
			})
		}
		if banana.IsButtonPressed(input.MouseButtonRight) {
			banana.RenderText("right mouse button pressed", &banana.TextRenderOptions{X: 10, Y: 35, Size: 12, Color: colornames.Red})
		}
		if banana.IsButtonPressed(input.MouseButtonMiddle) {
			banana.RenderText("middle mouse button pressed", &banana.TextRenderOptions{X: 10, Y: 45, Size: 12, Color: colornames.Red})
		}
		banana.RenderText("press space", &banana.TextRenderOptions{X: 45, Y: 180, Size: 12, Color: colornames.White})
	})
}
