package main

import (
	"image/color"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/gui"
	"golang.org/x/image/colornames"
)

func main() {
	width, height := 800, 600
	banana.SetWindowSize(width, height)
	ctx := gui.NewDrawContext(width, height)

	banana.Run(nil, func() {
		banana.Clear(colornames.Black)
		if ctx.Button(gui.ButtonOptions{
			ID:     "btn",
			X:      50,
			Y:      50,
			Width:  100,
			Height: 50,
			Label:  "Click Me",
		}) {
			println("button pressed")
		}

		ctx.Label(gui.LabelOptions{
			X:      50,
			Y:      120,
			Text:   "Hello, World!",
			Color:  color.White,
			PtSize: 14,
		})

		ctx.Slider(gui.SliderOptions{
			ID:       "slider",
			X:        50,
			Y:        150,
			Width:    200,
			Height:   20,
			Value:    50,
			MinValue: 0,
			MaxValue: 100,
			OnChange: func(value float32) {
				println("Slider value:", value)
			},
		})

		ctx.Toggle(gui.ToggleOptions{
			ID:     "toggle",
			X:      50,
			Y:      200,
			Width:  100,
			Height: 50,
			IsOn:   true,
			OnChange: func(isOn bool) {
				println("Toggle state:", isOn)
			},
		})

		ctx.TextBox(gui.TextBoxOptions{
			X:         50,
			Y:         270,
			Width:     200,
			Text:      "This is a\nmultiline text box.",
			TextColor: color.White,
			BgColor:   color.RGBA{30, 30, 30, 255},
			Padding:   10,
			PtSize:    14,
		})
	})
}
