package gui

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/pkg/input"
)

func (d *Draw) Label(options LabelOptions) {
	op := &banana.TextRenderOptions{
		X:     float32(options.X),
		Y:     float32(options.Y + options.PtSize),
		Size:  float32(options.PtSize),
		Color: options.Color,
	}
	d.DrawText(options.Text, op)
}

func (d *Draw) Slider(options SliderOptions) {
	mouseX, mouseY := banana.GetCursorPosition()

	if mouseX >= options.X && mouseX <= options.X+options.Width &&
		mouseY >= options.Y && mouseY <= options.Y+options.Height {
		d.SetHot(options.ID)
	} else if d.IsHot(options.ID) {
		d.SetHot("")
	}

	if banana.IsButtonJustPressed(input.MouseButtonLeft) && d.IsHot(options.ID) {
		d.SetActive(options.ID)
		d.SetOwner(options.ID)
	} else if d.IsOwner(options.ID) {
		d.SetActive("")
		d.SetOwner("")
	}

	d.DrawRectangle(options.X, options.Y, options.Width, options.Height, &DrawOptions{
		Style: Style{
			FillColor:    d.SecondaryColor,
			OutlineColor: d.HandleColor,
			OutlineSize:  2,
			CornerRadius: 5,
		},
	})

	percentage := (options.Value - options.MinValue) / (options.MaxValue - options.MinValue)
	filledWidth := int(float32(options.Width) * percentage)
	d.DrawRectangle(options.X, options.Y, filledWidth, options.Height, &DrawOptions{
		Style: Style{
			FillColor:    d.PrimaryColor,
			CornerRadius: 5,
		},
	})

	handleX := options.X + filledWidth
	handleY := options.Y + options.Height/2
	d.DrawCircle(handleX, handleY, options.Height/2, &DrawOptions{
		Style: Style{
			FillColor:    d.TextColor,
			OutlineColor: d.HandleColor,
			OutlineSize:  1,
		},
	})

	if d.IsActive(options.ID) {
		newX := mouseX - options.X
		if newX < 0 {
			newX = 0
		} else if newX > options.Width {
			newX = options.Width
		}

		percentage := float32(newX) / float32(options.Width)
		actualValue := options.MinValue + percentage*(options.MaxValue-options.MinValue)

		options.Value = actualValue

		if options.OnChange != nil {
			options.OnChange(options.Value)
		}
	}

	label := fmt.Sprintf("%d", int(options.Value))
	labelX := float32(options.X + (options.Width / 2) - len(label)*6)
	labelY := float32(options.Y + options.Height + 10)

	op := &banana.TextRenderOptions{}
	op.X = labelX
	op.Y = labelY - 12
	op.Size = 12
	op.Color = color.White
	d.DrawText(label, op)
}

func (d *Draw) TextBox(options TextBoxOptions) {
	lines := strings.Split(options.Text, "\n")
	height := options.Padding*2 + len(lines)*options.PtSize

	d.DrawRectangle(options.X, options.Y, options.Width, height, &DrawOptions{
		Style: Style{
			FillColor:    options.BgColor,
			OutlineColor: d.PrimaryColor,
			OutlineSize:  2,
			CornerRadius: 5,
		},
	})

	for i, line := range lines {
		op := &banana.TextRenderOptions{
			X:     float32(options.X + options.Padding),
			Y:     float32(options.Y + options.Padding + i*options.PtSize + options.PtSize),
			Size:  float32(options.PtSize),
			Color: options.TextColor,
		}
		d.DrawText(line, op)
	}
}
