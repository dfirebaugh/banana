package components

import (
	"image/color"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/gui"
	"github.com/dfirebaugh/banana/pkg/input"
)

type Toggle struct {
	Node
	x, y, width, height int
	isOn                bool
	isHovered           bool
	normalOutline       color.Color
	hoverOutline        color.Color
	onChange            func(bool)
}

func NewToggle(x, y, width, height int, initialState bool, onChange func(bool)) *Toggle {
	return &Toggle{
		x:        x,
		y:        y,
		width:    width,
		height:   height,
		isOn:     initialState,
		onChange: onChange,
	}
}

func (t *Toggle) updateHandlePosition() int {
	if t.isOn {
		return t.x + t.width - t.height
	}
	return t.x
}

func (t *Toggle) Render(ctx gui.DrawContext) {
	t.update()
	offsetX, offsetY := t.GetOffset()
	ctx.DrawRectangle(offsetX+t.x, offsetY+t.y, t.width, t.height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().SecondaryColor,
			OutlineColor: t.normalOutline,
			CornerRadius: 5,
		},
	})

	filledWidth := 0
	if t.isOn {
		filledWidth = t.width
	}
	ctx.DrawRectangle(offsetX+t.x, offsetY+t.y, filledWidth, t.height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().PrimaryColor,
			OutlineColor: color.Transparent,
			CornerRadius: 5,
		},
	})

	handleX := offsetX + t.updateHandlePosition()
	ctx.DrawRectangle(handleX, offsetY+t.y, t.height, t.height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().TextColor,
			OutlineColor: ctx.GetTheme().HandleColor,
			CornerRadius: 5,
			OutlineSize:  1,
		},
	})
}

func (t *Toggle) update() {
	mouseX, mouseY := banana.GetCursorPosition()
	offsetX, offsetY := t.GetGlobalOffset()

	if !t.IsWithinBounds(mouseX, mouseY) {
		t.isHovered = false
		// t.isPressed = false
		return
	}

	if mouseX >= offsetX+t.x && mouseX <= offsetX+t.x+t.width && mouseY >= offsetY+t.y && mouseY <= offsetY+t.y+t.height {
		if !t.isHovered {
			t.isHovered = true
		}
	} else {
		if t.isHovered {
			t.isHovered = false
		}
	}

	if banana.IsButtonJustPressed(input.MouseButtonLeft) {
		if mouseX >= offsetX+t.x && mouseX <= offsetX+t.x+t.width && mouseY >= offsetY+t.y && mouseY <= offsetY+t.y+t.height {
			t.isOn = !t.isOn

			if t.onChange != nil {
				t.onChange(t.isOn)
			}
		}
	}
}

func (t *Toggle) IsOn() bool {
	return t.isOn
}
