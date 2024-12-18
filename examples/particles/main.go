package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/exp/gui"
	"github.com/dfirebaugh/banana/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 800
	screenHeight = 600
	gravity      = 0.1
)

type Triangle struct {
	PositionX, PositionY int
	VelocityX, VelocityY float32
	Size                 int
	Age                  float64
	Lifetime             float64
	Options              gui.DrawOptions
}

var (
	triangles []Triangle
	frames    int
)

func newTriangleAtPosition(x, y int) Triangle {
	angle := rand.Float64() * 2 * math.Pi
	speed := rand.Float64()*2 + 2

	velocityX := float32(math.Cos(angle) * speed)
	velocityY := float32(math.Sin(angle) * speed)

	return Triangle{
		PositionX: x,
		PositionY: y,
		VelocityX: velocityX,
		VelocityY: velocityY,
		Size:      2,
		Lifetime:  2,
		Options: gui.DrawOptions{
			Style: gui.Style{
				FillColor: randomColor(),
			},
		},
	}
}

func randomColor() color.Color {
	return color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		255,
	}
}

var isFullScreen bool = false

func exampleControls() {
	if banana.IsKeyJustPressed(input.KeyF11) {
		isFullScreen = !isFullScreen
		banana.SetFullScreenBorderless(isFullScreen)
	}
	if banana.IsKeyJustPressed(input.KeyEscape) {
		banana.Close()
	}
}

func update() {
	exampleControls()
	for i := len(triangles) - 1; i >= 0; i-- {
		triangle := &triangles[i]
		triangle.VelocityY += gravity
		triangle.PositionX += int(triangle.VelocityX)
		triangle.PositionY += int(triangle.VelocityY)
		triangle.Age += 1.0 / 60.0

		if triangle.Age > triangle.Lifetime {
			triangles = append(triangles[:i], triangles[i+1:]...)
		}
	}

	if banana.IsButtonPressed(input.MouseButtonLeft) {
		x, y := banana.GetCursorPosition()
		for i := 0; i < 50; i++ {
			triangles = append(triangles, newTriangleAtPosition(x, y))
		}
	}

	frames++
	if frames%60 == 0 {
		fmt.Printf("Rendering %d triangles\n", len(triangles))
	}
}

func render() {
	banana.Clear(colornames.Black)
	w, h := banana.GetWindowSize()
	d := gui.Draw{
		ScreenWidth:  w,
		ScreenHeight: h,
	}

	for _, triangle := range triangles {
		d.DrawTriangle(
			triangle.PositionX,
			triangle.PositionY,
			triangle.PositionX+triangle.Size/2,
			triangle.PositionY-triangle.Size,
			triangle.PositionX+triangle.Size,
			triangle.PositionY,
			&triangle.Options,
		)
	}
	banana.RenderText(fmt.Sprintf("%d", len(triangles)),
		&banana.TextRenderOptions{
			X:     20,
			Y:     20,
			Size:  12,
			Color: colornames.Tomato,
		},
	)
}

func main() {
	banana.SetWindowSize(screenWidth, screenHeight)
	banana.EnableFPS()
	banana.Run(update, render)
}
