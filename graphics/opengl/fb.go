package opengl

import (
	"fmt"
	"image/color"

	"github.com/dfirebaugh/banana/graphics"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Framebuffer struct {
	ID         uint32
	TextureID  uint32
	RenderID   uint32
	Width      int
	Height     int
	ClearColor [4]float32
	renderer   *FramebufferRenderer
}

func NewFramebuffer(width, height int, textureManager *TextureManager, mainRenderer *Renderer) (*Framebuffer, error) {
	var fb Framebuffer
	fb.Width = width
	fb.Height = height

	gl.GenFramebuffers(1, &fb.ID)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.ID)

	gl.GenTextures(1, &fb.TextureID)
	println("textureID: ", fb.TextureID)
	gl.BindTexture(gl.TEXTURE_2D, fb.TextureID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, fb.TextureID, 0)

	gl.GenRenderbuffers(1, &fb.RenderID)
	gl.BindRenderbuffer(gl.RENDERBUFFER, fb.RenderID)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, fb.RenderID)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		return nil, fmt.Errorf("framebuffer is not complete")
	}

	textureManager.RegisterFramebufferTexture(fb.TextureID, width, height)

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	fbRenderer, err := NewFramebufferRenderer(mainRenderer, &fb)
	if err != nil {
		return nil, err
	}
	fb.renderer = fbRenderer
	return &fb, nil
}

func (fb *Framebuffer) GetTextureID() uint32 {
	return fb.TextureID
}

func (fb *Framebuffer) GetID() uint32 {
	return fb.ID
}

func (fb *Framebuffer) GetWidth() int {
	return fb.Width
}

func (fb *Framebuffer) GetHeight() int {
	return fb.Height
}

func (fb *Framebuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.ID)
	gl.Viewport(0, 0, int32(fb.Width), int32(fb.Height))
}

func (fb *Framebuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (fb *Framebuffer) Clear(c color.Color) {
	rgba := toRGBA(c)
	fb.ClearColor = rgba
	gl.ClearColor(rgba[0], rgba[1], rgba[2], rgba[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
}

func (fb *Framebuffer) Resize(width, height int) {
	fb.Width = width
	fb.Height = height

	gl.BindTexture(gl.TEXTURE_2D, fb.TextureID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	gl.BindRenderbuffer(gl.RENDERBUFFER, fb.RenderID)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, int32(width), int32(height))

	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.ID)
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		fmt.Println("framebuffer is not complete after resize")
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (fb *Framebuffer) Destroy() {
	gl.DeleteFramebuffers(1, &fb.ID)
	gl.DeleteTextures(1, &fb.TextureID)
	gl.DeleteRenderbuffers(1, &fb.RenderID)
}

func (fb *Framebuffer) Draw(x, y, width, height int) {
	fb.renderer.Draw(x, y, width, height)
}

type FramebufferRenderer struct {
	*Renderer
	framebuffer *Framebuffer
}

func NewFramebufferRenderer(renderer *Renderer, framebuffer *Framebuffer) (*FramebufferRenderer, error) {
	return &FramebufferRenderer{
		Renderer:    renderer,
		framebuffer: framebuffer,
	}, nil
}

func (fr *FramebufferRenderer) Bind() {
	fr.framebuffer.Bind()
}

func (fr *FramebufferRenderer) Unbind() {
	fr.framebuffer.Unbind()
}

func (fr *FramebufferRenderer) Clear(c color.Color) {
	fr.framebuffer.Clear(c)
}

func (fr *FramebufferRenderer) Draw(x, y, width, height int) {
	fr.Bind()
	fr.Renderer.Draw()
	fr.Unbind()

	gl.Viewport(int32(x), int32(y), int32(width), int32(height))

	fr.Renderer.RenderFramebuffer(fr.framebuffer, &graphics.TextureRenderOptions{
		X:             float32(x),
		Y:             float32(y),
		Width:         float32(width),
		Height:        float32(height),
		DesiredWidth:  float32(width),
		DesiredHeight: float32(height),
	})
}
