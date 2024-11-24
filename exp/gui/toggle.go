package gui

import (
	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/pkg/input"
)

func (d *Draw) toggleUpdate(options ToggleOptions) {
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
		options.IsOn = !options.IsOn
	} else if d.IsOwner(options.ID) {
		d.SetActive("")
		d.SetOwner("")
	}
}

func (d *Draw) toggleRender(options ToggleOptions) {
	fillColor := d.SecondaryColor
	if options.IsOn {
		fillColor = d.PrimaryColor
	}

	d.DrawRectangle(options.X, options.Y, options.Width, options.Height, &DrawOptions{
		Style: Style{
			FillColor:    fillColor,
			OutlineColor: d.HandleColor,
			CornerRadius: 5,
		},
	})

	handleX := options.X
	if options.IsOn {
		handleX = options.X + options.Width - options.Height
	}
	handleY := options.Y
	d.DrawRectangle(handleX, handleY, options.Height, options.Height, &DrawOptions{
		Style: Style{
			FillColor:    d.TextColor,
			OutlineColor: d.HandleColor,
			CornerRadius: 5,
			OutlineSize:  1,
		},
	})
}

func (d *Draw) Toggle(options ToggleOptions) {
	d.toggleUpdate(options)
	d.toggleRender(options)
}
