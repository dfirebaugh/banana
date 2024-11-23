package main

import (
	"bytes"
	"image"
	"image/color"
	"time"

	_ "image/png"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/assets"
	"github.com/dfirebaugh/banana/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 512
	screenHeight = 512
)

func main() {
	banana.SetWindowSize(screenWidth, screenHeight)
	banana.SetTitle("banana.sprite example")

	reader := bytes.NewReader(assets.BuddyDanceSpriteSheet)
	img, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	textureID := banana.UploadTexture(img)

	frameSize := image.Point{X: 32, Y: 32}
	sheetSize := image.Point{X: 4, Y: 1}

	frameIndex := 0
	lastFrameTime := time.Now()
	frameDuration := time.Millisecond * 200
	isBorderless := true
	isTextVisible := true
	banana.SetBorderlessWindowed(isBorderless)

	banana.Run(func() {
		if time.Since(lastFrameTime) >= frameDuration {
			lastFrameTime = time.Now()
			frameIndex = (frameIndex + 1) % (sheetSize.X * sheetSize.Y)
		}

		if banana.IsKeyJustPressed(input.KeyEscape) {
			banana.Close()
		}
		if banana.IsKeyJustPressed(input.KeyA) {
			isBorderless = !isBorderless
			banana.SetBorderlessWindowed(isBorderless)
		}
		if banana.IsKeyJustPressed(input.KeyD) {
			isTextVisible = !isTextVisible
		}
	}, func() {
		banana.Clear(color.RGBA{0, 0, 0, 0})

		frameX := (frameIndex % sheetSize.X) * frameSize.X
		frameY := (frameIndex / sheetSize.X) * frameSize.Y
		windowWidth, windowHeight := banana.GetWindowSize()

		options := &banana.TextureRenderOptions{
			X:             float32(screenWidth/2 - 256),
			Y:             float32(screenHeight/2 - 256),
			RectWidth:     float32(frameSize.X),
			RectHeight:    float32(frameSize.Y),
			Scale:         16, // 512 / 32 = 16
			RectX:         float32(frameX),
			RectY:         float32(frameY),
			Width:         float32(frameSize.X),
			Height:        float32(frameSize.Y),
			DesiredWidth:  float32(windowWidth),
			DesiredHeight: float32(windowHeight),
		}

		banana.RenderTexture(textureID, options)

		if isTextVisible {

			banana.RenderText("press escape to close", &banana.TextRenderOptions{
				X:     20,
				Y:     20,
				Size:  12,
				Color: colornames.Tomato,
			})
			banana.RenderText("press a for border", &banana.TextRenderOptions{
				X:     20,
				Y:     40,
				Size:  12,
				Color: colornames.Tomato,
			})
			banana.RenderText("press d to toggle text", &banana.TextRenderOptions{
				X:     20,
				Y:     60,
				Size:  12,
				Color: colornames.Tomato,
			})
		}
	})
}
