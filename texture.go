package banana

import (
	"image"

	"github.com/dfirebaugh/banana/graphics"
)

func UploadTexture(img image.Image) uint32 {
	return banana.graphicsBackend.UploadTexture(img)
}

func UpdateTexture(textureID uint32, img image.Image, xOffset, yOffset float32) {
	banana.graphicsBackend.UpdateTexture(textureID, img, int(xOffset), int(yOffset))
}

type TextureRenderOptions struct {
	TextureIndex                int
	X, Y                        float32
	RectWidth, RectHeight       float32
	DesiredWidth, DesiredHeight float32
	Scale                       float32
	RectX, RectY                float32
	Width, Height               float32
	FlipX, FlipY                bool
	Rotation                    float32
}

func RenderTexture(textureHandle uint32, options *TextureRenderOptions) {
	ensureSetupCompletion()
	banana.graphicsBackend.RenderTexture(
		textureHandle,
		&graphics.TextureRenderOptions{
			TextureIndex:  float32(options.TextureIndex),
			X:             options.X,
			Y:             options.Y,
			RectWidth:     options.RectWidth,
			RectHeight:    options.RectHeight,
			DesiredWidth:  options.DesiredWidth,
			DesiredHeight: options.DesiredHeight,
			Scale:         options.Scale,
			RectX:         options.RectX,
			RectY:         options.RectY,
			Width:         options.Width,
			Height:        options.Height,
			FlipX:         options.FlipX,
			FlipY:         options.FlipY,
			Rotation:      options.Rotation,
		})
}

type Framebuffer graphics.Framebuffer

func RenderFramebuffer(fb Framebuffer, options *TextureRenderOptions) {
	ensureSetupCompletion()
	graphicsOptions := &graphics.TextureRenderOptions{
		TextureIndex:  float32(options.TextureIndex),
		X:             options.X,
		Y:             options.Y,
		RectWidth:     options.RectWidth,
		RectHeight:    options.RectHeight,
		DesiredWidth:  options.DesiredWidth,
		DesiredHeight: options.DesiredHeight,
		Scale:         options.Scale,
		RectX:         options.RectX,
		RectY:         options.RectY,
		Width:         options.Width,
		Height:        options.Height,
		FlipX:         options.FlipX,
		FlipY:         options.FlipY,
		Rotation:      options.Rotation,
	}
	banana.graphicsBackend.RenderFramebuffer(fb, graphicsOptions)
}

func AddFramebuffer(width, height int) (Framebuffer, error) {
	ensureSetupCompletion()
	return banana.graphicsBackend.AddFramebuffer(width, height)
}

func ResizeFramebuffer(fb Framebuffer, width, height int) {
	fb.Resize(width, height)
}
