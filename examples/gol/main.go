package main

import (
	"math/rand"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/pkg/fb"
	"github.com/dfirebaugh/banana/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	gridWidth  = 400
	gridHeight = 300
	cellSize   = 1
)

var (
	grid     [gridWidth][gridHeight]bool
	nextGrid [gridWidth][gridHeight]bool
)

func initGrid() {
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			grid[x][y] = rand.Float64() < 0.2
		}
	}
}

func countNeighbors(x, y int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			neighborX := (x + i + gridWidth) % gridWidth
			neighborY := (y + j + gridHeight) % gridHeight
			if grid[neighborX][neighborY] {
				count++
			}
		}
	}
	return count
}

func updateGrid() {
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			neighbors := countNeighbors(x, y)
			if grid[x][y] {
				if neighbors < 2 || neighbors > 3 {
					nextGrid[x][y] = false
				} else {
					nextGrid[x][y] = true
				}
			} else {
				if neighbors == 3 {
					nextGrid[x][y] = true
				} else {
					nextGrid[x][y] = false
				}
			}
		}
	}

	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			grid[x][y] = nextGrid[x][y]
		}
	}
}

func renderGrid(fb *fb.ImageFB) {
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			cellColor := colornames.Black
			if grid[x][y] {
				cellColor = colornames.White
			}
			for i := 0; i < cellSize; i++ {
				for j := 0; j < cellSize; j++ {
					fb.SetPixel(int16(x*cellSize+i), int16(y*cellSize+j), cellColor)
				}
			}
		}
	}
}

func main() {
	banana.SetWindowSize(gridWidth, gridHeight)
	banana.EnableFPS()
	initGrid()
	framebuffer := fb.New(gridWidth*cellSize, gridHeight*cellSize)

	isFullScreen := false
	exampleControls := func() {
		if banana.IsKeyJustPressed(input.KeyF11) {
			isFullScreen = !isFullScreen
			banana.SetBorderlessWindowed(isFullScreen)
		}
		if banana.IsKeyJustPressed(input.KeyEscape) {
			banana.Close()
		}
	}
	renderGrid(framebuffer)
	texture := banana.UploadTexture(framebuffer.ToImage())

	banana.Run(func() {
		exampleControls()
		if banana.IsKeyPressed(input.KeyR) {
			initGrid()
			renderGrid(framebuffer)
			banana.UpdateTexture(texture, framebuffer.ToImage(), 0, 0)
		}
		updateGrid()
		renderGrid(framebuffer)
		banana.UpdateTexture(texture, framebuffer.ToImage(), 0, 0)
	}, func() {
		banana.Clear(colornames.Black)
		banana.RenderTexture(texture, &banana.TextureRenderOptions{
			X:          0,
			Y:          0,
			RectWidth:  gridWidth,
			RectHeight: gridHeight,
			Width:      gridWidth,
			Height:     gridHeight,
			Scale:      1.0,
		})
	})
}
