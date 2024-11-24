package gui

import (
	"image/color"

	"github.com/dfirebaugh/banana"
)

type (
	DrawContext interface {
		State
		Shapes
		Components
		DrawText(text string, options *banana.TextRenderOptions)
		GetTheme() *Theme
		SetTheme(t *Theme)
	}
	State interface {
		SetHot(id string)
		SetActive(id string)
		SetOwner(id string)
		IsHot(id string) bool
		IsActive(id string) bool
		IsOwner(id string) bool
	}
	Components interface {
		Label(options LabelOptions)
		Button(options ButtonOptions) bool
		Slider(options SliderOptions)
		TextBox(options TextBoxOptions)
		Toggle(options ToggleOptions)
	}
	Shapes interface {
		DrawCircle(x int, y int, radius int, op *DrawOptions)
		DrawLine(points []Position, op *DrawOptions)
		DrawCurve(start, control, end Position, op *DrawOptions)
		DrawWave(waveFunc func(x float64) float64, amplitude, frequency, phase float64, startX, endX, y int, op *DrawOptions)
		DrawRectangle(x int, y int, width int, height int, op *DrawOptions)
		DrawSegment(x1 int, y1 int, x2 int, y2 int, op *DrawOptions)
		DrawTriangle(x1, y1, x2, y2, x3, y3 int, op *DrawOptions)
	}
	Draw struct {
		ScreenWidth, ScreenHeight int
		*Theme
		hotID    string
		activeID string
		ownerID  string
	}
	DrawOptions struct {
		Style
	}
	Position struct {
		X, Y, Z int
	}
	Style struct {
		FillColor    color.Color
		OutlineColor color.Color
		OutlineSize  int
		CornerRadius int
	}
	ButtonOptions struct {
		ID                  string
		X, Y, Width, Height int
		Label               string
	}
	LabelOptions struct {
		ID     string
		X, Y   int
		Text   string
		Color  color.Color
		PtSize int
	}
	SliderOptions struct {
		ID                  string
		X, Y, Width, Height int
		Value               float32
		MinValue            float32
		MaxValue            float32
		OnChange            func(float32)
	}
	ToggleOptions struct {
		ID                  string
		X, Y, Width, Height int
		IsOn                bool
		OnChange            func(bool)
	}
	TextBoxOptions struct {
		ID          string
		X, Y, Width int
		Text        string
		TextColor   color.Color
		BgColor     color.Color
		Padding     int
		PtSize      int
	}
)

func NewDrawContext(width, height int) DrawContext {
	return &Draw{
		ScreenWidth:  width,
		ScreenHeight: height,
		Theme:        DefaultTheme(),
	}
}

func (d *Draw) SetHot(id string) {
	d.hotID = id
}

func (d *Draw) SetActive(id string) {
	d.activeID = id
}

func (d *Draw) SetOwner(id string) {
	d.ownerID = id
}

func (d *Draw) IsHot(id string) bool {
	return d.hotID == id
}

func (d *Draw) IsActive(id string) bool {
	return d.activeID == id
}

func (d *Draw) IsOwner(id string) bool {
	return d.ownerID == id
}

func (d *Draw) DrawTriangle(x1, y1, x2, y2, x3, y3 int, op *DrawOptions) {
	if op.OutlineSize > 0 {
		outlineOp := *op
		outlineOp.FillColor = op.OutlineColor
	}
	d.drawTriangle(x1, y1, x2, y2, x3, y3, op)
}

func (d *Draw) DrawRectangle(x, y, width, height int, op *DrawOptions) {
	if op.CornerRadius > 0 {
		if op.OutlineSize > 0 && op.CornerRadius > 0 {
			d.drawRoundedRectangleWithOutline(x, y, width, height, op.CornerRadius, op.OutlineSize, op)
		} else {
			d.drawRoundedRectangle(x, y, width, height, op.CornerRadius, op)
		}
	} else {
		if op.OutlineSize > 0 {
			outlineOp := *op
			outlineOp.FillColor = op.OutlineColor
		}
		d.drawRectangle(x, y, width, height, op)
	}
}

func (d *Draw) DrawCircle(x, y, radius int, op *DrawOptions) {
	if op.OutlineSize > 0 {
		d.drawCircleWithOutline(x, y, radius, op.OutlineSize, op)
	} else {
		d.drawCircle(x, y, radius, op)
	}
}

func (d *Draw) DrawSegment(x1, y1, x2, y2 int, op *DrawOptions) {
	banana.RenderShape(&banana.Segment{
		X1:    float32(x1),
		Y1:    float32(y1),
		X2:    float32(x2),
		Y2:    float32(y2),
		Width: float32(op.OutlineSize),
		Color: op.FillColor,
	})
}

func (d *Draw) DrawLine(points []Position, op *DrawOptions) {
	for i := 0; i < len(points)-1; i++ {
		d.DrawSegment(points[i].X, points[i].Y, points[i+1].X, points[i+1].Y, op)
	}
}

func (d *Draw) DrawCurve(start, control, end Position, op *DrawOptions) {
	vertices := generateQuadraticBezierVertices(start, control, end)
	for i := 0; i < len(vertices)-1; i++ {
		d.DrawSegment(int(vertices[i].X), int(vertices[i].Y), int(vertices[i+1].X), int(vertices[i+1].Y), op)
	}
}

func (d *Draw) DrawWave(waveFunc func(x float64) float64, amplitude, frequency, phase float64, startX, endX, y int, op *DrawOptions) {
	vertices := generateWaveVertices(waveFunc, amplitude, frequency, phase, startX, endX, y)
	drawWave(vertices, d, op)
}

func (d *Draw) DrawText(text string, options *banana.TextRenderOptions) {
	d.drawText(text, options)
}
