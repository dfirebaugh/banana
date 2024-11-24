package main

import (
	"math"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/exp/gui"
	"github.com/dfirebaugh/banana/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 240
	screenHeight = 160
)

func main() {
	banana.SetWindowSize(screenWidth, screenHeight)
	banana.SetTitle("Shapes Example")
	isFullScreen := false
	setFullScreen := func() {
		if banana.IsKeyJustPressed(input.KeyA) {
			isFullScreen = !isFullScreen
			banana.SetBorderlessWindowed(isFullScreen)
		}
		if banana.IsKeyJustPressed(input.KeyEscape) {
			banana.Close()
		}
	}

	banana.Run(func() {
		setFullScreen()
	}, func() {
		banana.Clear(colornames.Skyblue)
		d := gui.Draw{
			ScreenWidth:  screenWidth,
			ScreenHeight: screenHeight,
		}
		d.DrawRectangle(20, 20, 20, 20, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Green,
				OutlineColor: colornames.Purple,
			},
		})
		d.DrawRectangle(50, 50, 70, 70, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Green,
				OutlineColor: colornames.Purple,
				CornerRadius: 2,
			},
		})
		d.DrawRectangle(140, 50, 40, 40, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Green,
				OutlineColor: colornames.Purple,
				CornerRadius: 2,
				OutlineSize:  4,
			},
		})

		d.DrawCircle(60, 100, 15, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Red,
				OutlineColor: colornames.Black,
				OutlineSize:  2,
			},
		})
		d.DrawCircle(100, 100, 20, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Blue,
				OutlineColor: colornames.White,
				OutlineSize:  4,
			},
		})

		d.DrawTriangle(160, 120, 180, 140, 140, 140, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Orange,
				OutlineColor: colornames.Black,
				OutlineSize:  1,
			},
		})
		d.DrawTriangle(180, 60, 200, 80, 160, 80, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Yellow,
				OutlineColor: colornames.Black,
				OutlineSize:  1,
			},
		})

		points := []gui.Position{
			{X: 20, Y: 150},
			{X: 60, Y: 130},
			{X: 100, Y: 150},
			{X: 140, Y: 130},
			{X: 180, Y: 150},
		}

		for i := 0; i < len(points)-1; i++ {
			d.DrawLine([]gui.Position{
				points[i],
				points[i+1],
			}, &gui.DrawOptions{
				Style: gui.Style{
					FillColor:    colornames.Red,
					OutlineColor: colornames.Black,
					OutlineSize:  2,
				},
			})
		}
		d.DrawCurve(gui.Position{X: 20, Y: 50}, gui.Position{X: 120, Y: 10}, gui.Position{X: 220, Y: 50}, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Blue,
				OutlineColor: colornames.Black,
				OutlineSize:  2,
			},
		})

		d.DrawWave(sineWave, 20, 0.1, 0, 0, screenWidth, 20, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Blue,
				OutlineColor: colornames.Black,
				OutlineSize:  2,
			},
		})

		d.DrawWave(sawtoothWave, 20, 0.1, 0, 0, screenWidth, 40, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Green,
				OutlineColor: colornames.Black,
				OutlineSize:  2,
			},
		})

		d.DrawWave(triangleWave, 20, 0.1, 0, 0, screenWidth, 60, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Red,
				OutlineColor: colornames.Black,
				OutlineSize:  2,
			},
		})

		d.DrawWave(squareWave, 20, 0.1, 0, 0, screenWidth, 80, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Yellow,
				OutlineColor: colornames.Black,
				OutlineSize:  2,
			},
		})

		d.DrawWave(func(x float64) float64 {
			return math.Sin(x) * math.Cos(0.05*x)
		}, 20, 0.1, 0, 0, screenWidth, 100, &gui.DrawOptions{
			Style: gui.Style{
				FillColor:    colornames.Purple,
				OutlineColor: colornames.Black,
				OutlineSize:  2,
			},
		})
	})
}

func sineWave(x float64) float64 {
	return math.Sin(x)
}

func sawtoothWave(x float64) float64 {
	return 2 * (x - math.Floor(x+0.5))
}

func triangleWave(x float64) float64 {
	return 2*math.Abs(2*(x-math.Floor(x+0.5))) - 1
}

func squareWave(x float64) float64 {
	return math.Copysign(1, math.Sin(x))
}
