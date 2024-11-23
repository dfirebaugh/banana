package main

import (
	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/pkg/input"
	"golang.org/x/image/colornames"
)

var isFullScreen = false

func exampleControls() {
	if banana.IsKeyJustPressed(input.KeyF11) {
		isFullScreen = !isFullScreen
		banana.SetBorderlessWindowed(isFullScreen)
	}
	if banana.IsKeyJustPressed(input.KeyEscape) {
		banana.Close()
	}
}

func update() {
	exampleControls()
}

func render() {
	banana.Clear(colornames.White)
	rect := &banana.Rect{
		X:      0,
		Y:      0,
		Width:  50,
		Height: 50,
		Radius: 1,
		Color:  colornames.Yellow,
	}
	banana.RenderShape(rect)
	banana.RenderText("hello, world",
		&banana.TextRenderOptions{
			X:     10,
			Y:     30,
			Size:  12,
			Color: colornames.Red,
		},
	)
}

func main() {
	banana.SetWindowSize(200, 200)
	banana.Run(update, render)
}
