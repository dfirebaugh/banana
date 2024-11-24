package components

import (
	"image/color"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/exp/gui"
	"github.com/dfirebaugh/banana/pkg/input"
)

type Button struct {
	Node
	label               string
	labelPtSize         int
	x, y, width, height int
	cornerRadius        int
	onClick             func()
	isHovered           bool
	isPressed           bool
}

func NewButton(x, y, width, height int, cornerRadius int, label string, onClick func()) *Button {
	b := &Button{
		label:        label,
		labelPtSize:  12,
		x:            x,
		y:            y,
		width:        width,
		height:       height,
		cornerRadius: cornerRadius,
		onClick:      onClick,
	}

	return b
}

func (b *Button) Render(ctx gui.DrawContext) {
	b.update()
	offsetX, offsetY := b.GetOffset()
	ctx.DrawRectangle(offsetX+b.x, offsetY+b.y, b.width, b.height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    b.getFillColor(ctx),
			OutlineColor: ctx.GetTheme().SecondaryColor,
			OutlineSize:  2,
			CornerRadius: b.cornerRadius,
		},
	})

	op := &banana.TextRenderOptions{}
	op.X = float32(offsetX + b.x + (b.width / 2) - (len(b.label) * b.labelPtSize / 3))
	op.Y = float32(offsetY + b.y + (b.height / 2) + b.labelPtSize/2)
	op.Size = float32(b.labelPtSize)
	op.Color = color.White
	ctx.DrawText(b.label, op)

	b.Node.Render(ctx)
}

func (b *Button) update() {
	b.Node.Update()
	mouseX, mouseY := banana.GetCursorPosition()
	offsetX, offsetY := b.GetGlobalOffset()

	if !b.IsWithinBounds(mouseX, mouseY) {
		b.isHovered = false
		b.isPressed = false
		return
	}

	b.isHovered = b.isPointWithin(mouseX-offsetX, mouseY-offsetY)

	if banana.IsButtonJustPressed(input.MouseButtonLeft) && b.isHovered {
		b.onClick()
	} else if banana.IsButtonPressed(input.MouseButtonLeft) && b.isHovered {
		b.isPressed = true
	} else {
		b.isPressed = false
	}
}

func (b *Button) getFillColor(ctx gui.DrawContext) color.Color {
	if b.isPressed {
		return ctx.GetTheme().HandleColor
	}
	if b.isHovered {
		return ctx.GetTheme().BackgroundColor
	}
	return ctx.GetTheme().PrimaryColor
}

func (b *Button) isPointWithin(x, y int) bool {
	return x >= b.x && x <= b.x+b.width && y >= b.y && y <= b.y+b.height
}

func (b *Button) IsPressed() bool {
	return b.isPressed
}

func (b *Button) SetText(text string) {
	b.label = text
}
