package components

import (
	"image/color"
	"strings"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/exp/gui"
)

type TextBox struct {
	Node
	x, y, width, height int
	text                string
	textColor           color.Color
	bgColor             color.Color
	padding             int
	ptSize              int
}

func NewTextBox(x, y, width int, textColor, bgColor color.Color) *TextBox {
	return &TextBox{
		x:         x,
		y:         y,
		width:     width,
		height:    0,
		textColor: textColor,
		bgColor:   bgColor,
		padding:   10,
		ptSize:    14,
	}
}

func (tb *TextBox) SetText(text string) {
	tb.text = text
	tb.updateHeight()
}

func (tb *TextBox) updateHeight() {
	lines := strings.Split(tb.text, "\n")
	tb.height = tb.padding*2 + len(lines)*tb.ptSize
}

func (tb *TextBox) Height() int {
	lines := strings.Split(tb.text, "\n")
	return tb.padding*2 + len(lines)*tb.ptSize
}

func (tb *TextBox) Render(ctx gui.DrawContext) {
	tb.update()
	globalX, globalY := tb.GetGlobalOffset()

	ctx.DrawRectangle(globalX+tb.x, globalY+tb.y, tb.width, tb.height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    tb.bgColor,
			OutlineColor: ctx.GetTheme().PrimaryColor,
			OutlineSize:  2,
			CornerRadius: 5,
		},
	})

	lines := strings.Split(tb.text, "\n")
	for i, line := range lines {
		op := &banana.TextRenderOptions{
			X:     float32(globalX + tb.x + tb.padding),
			Y:     float32(globalY + tb.y + tb.padding + i*tb.ptSize + tb.ptSize),
			Size:  float32(tb.ptSize),
			Color: tb.textColor,
		}
		ctx.DrawText(line, op)
	}
}

func (tb *TextBox) Update() {
}

func (tb *TextBox) update() {
	tb.Node.Update()
}
