// Package graphics provides functionality for 2D graphics rendering,
// including textures, sprites, text, and shapes.
package graphics

import (
	"image"
	"image/color"

	"github.com/dfirebaugh/banana/pkg/input"
)

type Renderable interface {
	GetVertices(screenWidth, screenHeight int) []Vertex
}

type Font interface{}

type TextureRenderOptions struct {
	TextureIndex                float32
	X, Y                        float32
	RectWidth, RectHeight       float32
	DesiredWidth, DesiredHeight float32
	Scale                       float32
	RectX, RectY                float32
	Width, Height               float32
	FlipX, FlipY                bool
	Rotation                    float32
}

type TextRenderOptions struct {
	X, Y, Size float32
	Color      color.Color
}

type OpCode float32

const (
	OP_CODE_VERTEX  = 1.0
	OP_CODE_CIRCLE  = 2.0
	OP_CODE_RECT    = 3.0
	OP_CODE_TEXT    = 4.0
	OP_CODE_TEXTURE = 5.0
)

type Vertex struct {
	FsQuadPos    [2]float32
	ShapePos     [2]float32
	LocalPos     [2]float32
	OpCode       OpCode
	Radius       float32
	Width        float32
	Height       float32
	Color        [4]float32
	Resolution   [2]float32
	TexCoord     [2]float32
	TextureIndex float32
	FontIndex    float32
}

type Framebuffer interface {
	GetID() uint32
	GetTextureID() uint32
	Bind()
	Clear(c color.Color)
	Destroy()
	Unbind()
	GetWidth() int
	GetHeight() int
	Resize(width, height int)
	Draw(x, y, width, height int)
}

type GraphicsBackend interface {
	WindowManager
	InputManager
	EventManager
	TextureManager
	Clear(c color.Color)
	Close()
	Init()
	Draw()
	Render(shape Renderable)
	RenderText(text string, options *TextRenderOptions)
	LoadFont(fontPath []byte) (Font, error)
	SwapBuffers()
	GetViewportSize() (int, int)
	AddFramebuffer(width, height int) (Framebuffer, error)
	RenderFramebuffer(fb Framebuffer, options *TextureRenderOptions)
	BindFramebuffer(fb Framebuffer)
	UnbindFramebuffer()
	Begin()
	End()
}

type TextureManager interface {
	RenderTexture(textureHandle uint32, options *TextureRenderOptions)
	UploadTexture(image.Image) uint32
	UpdateTexture(textureID uint32, img image.Image, xOffset, yOffset int)
}

type WindowManager interface {
	DisableWindowResize()
	SetFullScreenBorderless(v bool)
	SetBorderlessWindowed(v bool)
	SetWindowTitle(title string)
	DestroyWindow()
	SetWindowSize(width int, height int)
	GetWindowSize() (int, int)
	GetWindowPosition() (x int, y int)
	SetWindowPosition(x, y int)
	SetResizedCallback(fn func(physicalWidth, physicalHeight uint32))
	ShouldClose() bool
	IsDisposed() bool
}

type EventManager interface {
	PollEvents() bool
}

type InputManager interface {
	SetInputCallback(fn func(eventChan chan input.Event))
}
