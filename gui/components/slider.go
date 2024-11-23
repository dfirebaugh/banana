package components

import (
	"fmt"
	"image/color"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/gui"
	"github.com/dfirebaugh/banana/pkg/input"
)

type Slider struct {
	Node
	x, y, width, height    int
	handleRadius           int
	value                  float32
	minValue               float32
	maxValue               float32
	isDragging             bool
	isHovered              bool
	prevMouseButtonPressed bool
	normalOutline          color.Color
	hoverOutline           color.Color
	onChange               func(float32)
}

func NewSlider(x, y, width, height int, minValue, maxValue float32, onChange func(float32)) *Slider {
	handleRadius := height / 2

	s := &Slider{
		x:            x,
		y:            y,
		width:        width,
		height:       height,
		handleRadius: handleRadius,
		value:        minValue,
		minValue:     minValue,
		maxValue:     maxValue,
		onChange:     onChange,
	}

	return s
}

func (s *Slider) Render(ctx gui.DrawContext) {
	s.update()

	ctx.DrawRectangle(s.x, s.y, s.width, s.height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().SecondaryColor,
			OutlineColor: s.getTrackOutlineColor(ctx),
			OutlineSize:  2,
			CornerRadius: 5,
		},
	})

	percentage := (s.value - s.minValue) / (s.maxValue - s.minValue)
	filledWidth := int(float32(s.width) * percentage)
	ctx.DrawRectangle(s.x, s.y, filledWidth, s.height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().PrimaryColor,
			CornerRadius: 5,
		},
	})

	handleX := s.x + filledWidth
	handleY := s.y + s.height/2
	ctx.DrawCircle(handleX, handleY, s.handleRadius, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().TextColor,
			OutlineColor: ctx.GetTheme().HandleColor,
			OutlineSize:  1,
		},
	})

	label := fmt.Sprintf("%d", int(s.value))
	labelX := float32(s.x + (s.width / 2) - len(label)*6)
	labelY := float32(s.y + s.height + 10)

	op := &banana.TextRenderOptions{}
	op.X = labelX
	op.Y = labelY - 12
	op.Size = 12
	op.Color = color.White
	ctx.DrawText(label, op)
}

func (s *Slider) update() {
	mouseX, mouseY := banana.GetCursorPosition()
	offsetX, offsetY := s.GetGlobalOffset()

	if !s.IsWithinBounds(mouseX, mouseY) {
		s.isHovered = false
		s.isDragging = false
		return
	}

	if s.isPointWithin(mouseX-offsetX, mouseY-offsetY) {
		s.isHovered = true
	} else {
		s.isHovered = false
	}

	currentMouseButtonPressed := banana.IsButtonPressed(input.MouseButtonLeft)

	if currentMouseButtonPressed && !s.prevMouseButtonPressed {
		if s.isHandlePointWithin(mouseX-offsetX, mouseY-offsetY) {
			s.isDragging = true
		} else {
			s.isDragging = false
		}
	} else if currentMouseButtonPressed && s.prevMouseButtonPressed {
		if s.isDragging {
			newX := mouseX - offsetX
			if newX < s.x {
				newX = s.x
			} else if newX > s.x+s.width {
				newX = s.x + s.width
			}

			percentage := float32(newX-s.x) / float32(s.width)
			actualValue := s.minValue + percentage*(s.maxValue-s.minValue)

			s.value = actualValue

			if s.onChange != nil {
				s.onChange(s.value)
			}
		}
	} else {
		s.isDragging = false
	}

	s.prevMouseButtonPressed = currentMouseButtonPressed
}

func (s *Slider) GetValue() float32 {
	return s.value
}

func (s *Slider) isPointWithin(x, y int) bool {
	return x >= s.x && x <= s.x+s.width && y >= s.y && y <= s.y+s.height
}

func (s *Slider) isHandlePointWithin(x, y int) bool {
	handleX := s.x + int(float32(s.width)*(s.value-s.minValue)/(s.maxValue-s.minValue)) - s.handleRadius
	handleY := s.y + s.height/2 - s.handleRadius
	return x >= handleX && x <= handleX+s.handleRadius*2 && y >= handleY && y <= handleY+s.handleRadius*2
}

func (s *Slider) getTrackOutlineColor(ctx gui.DrawContext) color.Color {
	if s.isHovered {
		return ctx.GetTheme().PrimaryColor
	}
	return ctx.GetTheme().HandleColor
}

func (s *Slider) SetValue(actualValue float32) {
	if actualValue < s.minValue {
		actualValue = s.minValue
	} else if actualValue > s.maxValue {
		actualValue = s.maxValue
	}

	s.value = actualValue

	if s.onChange != nil {
		s.onChange(s.value)
	}
}
