package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/exp/gui"
	"github.com/dfirebaugh/banana/exp/gui/components"
	"github.com/dfirebaugh/banana/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	maxBalls = 275
)

var (
	balls          []Ball
	numBalls       int
	damping        = float32(0.9) // Damping factor to simulate energy loss during collisions
	enableFade     = true
	enableCollide  = true
	slider         *components.Slider
	fadeToggle     *components.Toggle
	lockToggle     *components.Toggle
	collideToggle  *components.Toggle
	resetButton    *components.Button
	ballCountLabel *components.Label
	drawContext    gui.DrawContext
)

func main() {
	banana.SetTitle("Bouncing Balls")
	banana.SetWindowSize(960, 640)
	banana.EnableFPS()

	numBalls = 3
	balls = make([]Ball, numBalls)
	for i := 0; i < numBalls; i++ {
		balls[i] = NewBall()
	}

	surface := setupUI()
	// Initialize the draw context
	drawContext = gui.NewDrawContext(960, 640)

	banana.Run(func() {
		handleInput()
		for i := range balls {
			balls[i].Update()
		}
	}, func() {
		surface.PrepareRender(drawContext)
		banana.Clear(colornames.Black)
		for i := range balls {
			balls[i].Render(drawContext)
		}
		surface.Render(drawContext)
	})
}

type Ball struct {
	X, Y                 float32
	VelocityX, VelocityY float32
	Radius               float32
	Color                color.RGBA
	TargetColor          color.RGBA
	ColorChangeSpeed     float32
}

func (b *Ball) Update() {
	b.X += b.VelocityX
	b.Y += b.VelocityY

	b.handleBorderCollision()

	if enableCollide {
		for i := range balls {
			if &balls[i] != b {
				b.checkCollision(&balls[i])
			}
		}
	}

	if enableFade {
		b.updateColor()
	}
}

func (b *Ball) updateColor() {
	b.Color.R = uint8(float32(b.Color.R) + (float32(b.TargetColor.R)-float32(b.Color.R))*b.ColorChangeSpeed)
	b.Color.G = uint8(float32(b.Color.G) + (float32(b.TargetColor.G)-float32(b.Color.G))*b.ColorChangeSpeed)
	b.Color.B = uint8(float32(b.Color.B) + (float32(b.TargetColor.B)-float32(b.Color.B))*b.ColorChangeSpeed)
}

func (b *Ball) Render(drawContext gui.DrawContext) {
	drawContext.DrawCircle(int(b.X), int(b.Y), int(b.Radius), &gui.DrawOptions{
		Style: gui.Style{
			FillColor: b.Color,
		},
	})
}

func NewBall() Ball {
	sw, sh := banana.GetWindowSize()
	radius := float32(rand.Intn(35) + 10)
	x := radius + float32(rand.Float64()*(float64(sw)-2*float64(radius)))
	y := radius + float32(rand.Float64()*(float64(sh)-2*float64(radius)))
	velocityX := float32(rand.Float64()*4 - 2)
	velocityY := float32(rand.Float64()*4 - 2)

	randomColor := color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}

	return Ball{
		X:                x,
		Y:                y,
		Radius:           radius,
		VelocityX:        velocityX,
		VelocityY:        velocityY,
		Color:            randomColor,
		TargetColor:      randomColor,
		ColorChangeSpeed: 0.1,
	}
}

func handleInput() {
	if fadeToggle != nil {
		enableFade = fadeToggle.IsOn()
	}

	if collideToggle != nil {
		enableCollide = collideToggle.IsOn()
	}

	if banana.IsButtonPressed(input.MouseButtonLeft) {
		x, y := banana.GetCursorPosition()
		for i := range balls {
			if balls[i].containsPoint(float32(x), float32(y)) {
				balls[i].VelocityX += float32(rand.Float64()*2-1) * 5
				balls[i].VelocityY += float32(rand.Float64()*2-1) * 5
			}
		}
	}
}

func (b *Ball) containsPoint(px, py float32) bool {
	dx := b.X - px
	dy := b.Y - py
	distance := math.Sqrt(float64(dx*dx + dy*dy))
	return distance < float64(b.Radius)
}

func (b *Ball) handleBorderCollision() {
	sw, sh := banana.GetWindowSize()
	if b.X-b.Radius < 0 {
		b.X = b.Radius
		b.VelocityX = -b.VelocityX * damping
	}
	if b.X+b.Radius > float32(sw) {
		b.X = float32(sw) - b.Radius
		b.VelocityX = -b.VelocityX * damping
	}
	if b.Y-b.Radius < 0 {
		b.Y = b.Radius
		b.VelocityY = -b.VelocityY * damping
	}
	if b.Y+b.Radius > float32(sh) {
		b.Y = float32(sh) - b.Radius
		b.VelocityY = -b.VelocityY * damping
	}
}

func (b *Ball) checkCollision(other *Ball) {
	dx := b.X - other.X
	dy := b.Y - other.Y
	distance := math.Sqrt(float64(dx*dx + dy*dy))

	if distance < float64(b.Radius+other.Radius) {
		b.resolveCollision(other)
	}
}

func (b *Ball) resolveCollision(other *Ball) {
	dx := b.X - other.X
	dy := b.Y - other.Y
	distance := math.Sqrt(float64(dx*dx + dy*dy))

	nx := dx / float32(distance)
	ny := dy / float32(distance)

	tx := -ny
	ty := nx

	dpTan1 := b.VelocityX*tx + b.VelocityY*ty
	dpTan2 := other.VelocityX*tx + other.VelocityY*ty

	dpNorm1 := b.VelocityX*nx + b.VelocityY*ny
	dpNorm2 := other.VelocityX*nx + other.VelocityY*ny

	m1 := (dpNorm1*(b.Radius-other.Radius) + 2.0*other.Radius*dpNorm2) / (b.Radius + other.Radius)
	m2 := (dpNorm2*(other.Radius-b.Radius) + 2.0*b.Radius*dpNorm1) / (b.Radius + other.Radius)

	b.VelocityX = (tx*dpTan1 + nx*m1) * damping
	b.VelocityY = (ty*dpTan1 + ny*m1) * damping
	other.VelocityX = (tx*dpTan2 + nx*m2) * damping
	other.VelocityY = (ty*dpTan2 + ny*m2) * damping

	overlap := 0.5 * (float32(distance) - b.Radius - other.Radius)
	b.X -= overlap * nx
	b.Y -= overlap * ny
	other.X += overlap * nx
	other.Y += overlap * ny

	if enableFade {
		b.TargetColor = randomColor()
		other.TargetColor = randomColor()
	}
}

func randomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

func adjustBallCount(targetCount int) {
	currentCount := len(balls)

	if currentCount < targetCount {
		for i := 0; i < targetCount-currentCount; i++ {
			balls = append(balls, NewBall())
		}
	} else if currentCount > targetCount {
		balls = balls[:targetCount]
	}
}

func setupUI() *components.Surface {
	surface := components.NewSurface(10, 10, 260, 220)

	surface.AppendChild(components.NewLabel(10, 30, "Ball Control"))

	slider = components.NewSlider(10, 60, 240, 20, 0, maxBalls, func(value float32) {
		adjustBallCount(int(value))
		if ballCountLabel != nil {
			ballCountLabel.SetText(fmt.Sprintf("Number of Balls: %d", int(value)))
		}
	})
	surface.AppendChild(slider)

	fadeToggle = components.NewToggle(200, 90, 30, 20, true, func(isOn bool) {
		enableFade = isOn
	})
	surface.AppendChild(fadeToggle)

	surface.AppendChild(components.NewLabel(10, 90, "Enable Color Fade"))

	collideToggle = components.NewToggle(200, 120, 30, 20, true, func(isOn bool) {
		enableCollide = isOn
	})
	surface.AppendChild(collideToggle)

	surface.AppendChild(components.NewLabel(10, 120, "Enable Collision"))

	resetButton = components.NewButton(10, 150, 240, 20, 5, "Reset Balls", func() {
		for i := range balls {
			balls[i] = NewBall()
		}
	})
	surface.AppendChild(resetButton)

	lockToggle = components.NewToggle(200, 180, 30, 20, false, func(isOn bool) {
		surface.IsLocked = isOn
	})
	surface.AppendChild(components.NewLabel(10, 180, "Lock Surface"))
	surface.AppendChild(lockToggle)

	return surface
}
