package main

import (
	"fmt"
	"image/color"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/exp/gui"
	"github.com/dfirebaugh/banana/exp/gui/components"
	"github.com/dfirebaugh/banana/pkg/input"
	"golang.org/x/image/colornames"
)

var (
	windowWidth  = 840
	windowHeight = 700
)

func main() {
	banana.SetWindowSize(windowWidth, windowHeight)
	banana.SetTitle("UI Components Example")
	banana.EnableFPS()

	banana.SetResizeCallback(func(physicalWidth, physicalHeight uint32) {
		windowWidth = int(physicalWidth)
		windowHeight = int(physicalHeight)
	})

	buttons := []*components.Button{}
	for i := 0; i < 16; i++ {
		x := 32 + (i%4)*200
		y := 60 + (i/4)*60
		index := i
		button := components.NewButton(x, y, 120, 40, 5, fmt.Sprintf("Btn %d", i+1), func() {
			fmt.Printf("Button %d clicked\n", index+1)
		})
		buttons = append(buttons, button)
	}

	sliders := []*components.Slider{}
	for i := 0; i < 10; i++ {
		x := 32 + (i%5)*160
		y := 310 + (i/5)*46
		slider := components.NewSlider(x, y, 100, 15, 0, 100, func(value float32) {
			fmt.Printf("Slider %d value: %f\n", i+1, value)
		})
		sliders = append(sliders, slider)
	}

	toggles := []*components.Toggle{}
	for i := 0; i < 10; i++ {
		x := 32 + (i%5)*160
		y := 390 + (i/5)*40
		toggle := components.NewToggle(x, y, 40, 20, false, func(isOn bool) {
			fmt.Printf("Toggle %d state: %t\n", i+1, isOn)
		})
		toggles = append(toggles, toggle)
	}

	rectColors := []color.Color{
		colornames.Mediumvioletred,
		colornames.Darkorange,
		colornames.Gold,
		colornames.Mediumseagreen,
	}

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
	banana.Run(func() {
		exampleControls()

		for _, button := range buttons {
			button.Update()
		}
		for _, slider := range sliders {
			slider.Update()
		}
		for _, toggle := range toggles {
			toggle.Update()
		}
	}, func() {
		banana.Clear(color.RGBA{0, 0, 0, 0})
		banana.RenderShape(&banana.Rect{
			X:      0,
			Y:      0,
			Width:  float32(windowWidth),
			Height: 10,
			Radius: 1,
			Color:  colornames.Gainsboro,
		})
		banana.RenderShape(&banana.Rect{
			X:      0,
			Y:      0,
			Width:  float32(windowWidth),
			Height: float32(windowHeight),
			Radius: 15,
			Color:  colornames.Gainsboro,
		})
		w, h := banana.GetWindowSize()
		drawContext := gui.NewDrawContext(w, h)

		for _, button := range buttons {
			button.Render(drawContext)
		}
		for _, slider := range sliders {
			slider.Render(drawContext)
		}
		for _, toggle := range toggles {
			toggle.Render(drawContext)
		}

		for i := 0; i < 8; i++ {
			x := 32 + (i%4)*200
			y := 490 + (i/4)*60
			color := rectColors[i%4]
			drawContext.DrawRectangle(x, y, 120, 60, &gui.DrawOptions{
				Style: gui.Style{
					FillColor:    color,
					OutlineColor: colornames.White,
					OutlineSize:  2,
					CornerRadius: 6,
				},
			})
		}
	})
}
