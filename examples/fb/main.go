package main

import (
	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	width  = 800
	height = 600
)

func main() {
	banana.SetWindowSize(width, height)
	isFullScreen := false
	exampleControls := func() {
		if banana.IsKeyJustPressed(input.KeyF11) {
			isFullScreen = !isFullScreen
			banana.SetBorderlessWindowed(isFullScreen)
		}
		if banana.IsKeyJustPressed(input.KeyEscape) {
			banana.Close()
		}
	}

	fb1, _ := banana.AddFramebuffer(width/2, height/2)
	banana.BindFramebuffer(fb1)
	banana.Clear(colornames.Darkblue)
	banana.RenderShape(&banana.Rect{
		X:      0,
		Y:      0,
		Width:  40,
		Height: 40,
		Color:  colornames.Seashell,
	})
	banana.RenderShape(&banana.Rect{
		X:      width/2 - 20,
		Y:      50,
		Width:  40,
		Height: 40,
		Radius: 10,
		Color:  colornames.Tomato,
	})
	banana.RenderText("hello, world", &banana.TextRenderOptions{
		X:     50,
		Y:     50,
		Size:  16,
		Color: colornames.Tomato,
	})
	fb1.Draw(100, 0, width/2, height/2)
	banana.UnbindFramebuffer()

	banana.Run(nil, func() {
		exampleControls()
		banana.Clear(colornames.Black)
		banana.RenderFramebuffer(fb1, &banana.TextureRenderOptions{
			X:             100,
			Y:             0,
			Width:         width / 2,
			Height:        height / 2,
			RectWidth:     width / 2,
			RectHeight:    height / 2,
			DesiredWidth:  width / 2,
			DesiredHeight: height / 2,
		})
	})
}
