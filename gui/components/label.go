package components

import (
	"image/color"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/gui"
)

type Label struct {
	Node
	x, y   int
	text   string
	color  color.Color
	ptSize int
}

func NewLabel(x, y int, text string) *Label {
	return &Label{
		x:      x,
		y:      y,
		text:   text,
		color:  color.White,
		ptSize: 12,
	}
}

func (l *Label) Render(ctx gui.DrawContext) {
	// x, y := l.GetOffset()
	op := &banana.TextRenderOptions{}
	op.X = float32(l.x)
	op.Y = float32(l.y + l.ptSize)
	op.Size = float32(l.ptSize)
	op.Color = l.color
	op.Color = ctx.GetTheme().TextColor
	ctx.DrawText(l.text, op)
}

func (l *Label) SetText(text string) {
	l.text = text
}

func (l *Label) SetColor(color color.Color) {
	l.color = color
}

func (l *Label) SetPtSize(ptSize int) {
	l.ptSize = ptSize
}
