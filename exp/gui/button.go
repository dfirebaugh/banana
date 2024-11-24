package gui

import (
	"image/color"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/pkg/input"
)

func (d *Draw) buttonUpdate(options ButtonOptions) {
	mouseX, mouseY := banana.GetCursorPosition()

	if mouseX >= options.X && mouseX <= options.X+options.Width &&
		mouseY >= options.Y && mouseY <= options.Y+options.Height {
		d.SetHot(options.ID)
	} else {
		d.SetHot("")
	}

	if banana.IsButtonJustPressed(input.MouseButtonLeft) && d.IsHot(options.ID) {
		d.SetActive(options.ID)
		d.SetOwner(options.ID)
	} else if d.IsOwner(options.ID) {
		d.SetActive("")
		d.SetOwner("")
	}
}

func (d *Draw) buttonRender(options ButtonOptions) {
	fillColor := d.PrimaryColor
	if d.IsActive(options.ID) {
		fillColor = d.SecondaryColor
	} else if d.IsHot(options.ID) {
		fillColor = d.HandleColor
	}
	if banana.IsButtonJustPressed(input.MouseButtonLeft) && d.IsHot(options.ID) {
		fillColor = d.SecondaryColor
	}

	d.DrawRectangle(options.X, options.Y, options.Width, options.Height, &DrawOptions{
		Style: Style{
			FillColor:    fillColor,
			OutlineColor: d.SecondaryColor,
			OutlineSize:  2,
			CornerRadius: 5,
		},
	})

	op := &banana.TextRenderOptions{
		X:     float32(options.X + (options.Width / 2) - (len(options.Label) * 6) + 14),
		Y:     float32(options.Y + (options.Height / 2) + 6),
		Size:  12,
		Color: color.White,
	}
	d.DrawText(options.Label, op)
}

func (d *Draw) Button(options ButtonOptions) bool {
	d.buttonUpdate(options)
	d.buttonRender(options)

	return d.IsHot(options.ID) && d.IsActive(options.ID) && d.IsOwner(options.ID)
}
