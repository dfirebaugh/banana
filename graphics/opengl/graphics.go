package opengl

import (
	"log"

	"github.com/dfirebaugh/banana/graphics/window"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type GraphicsBackend struct {
	*window.Window
	*Renderer
}

func NewGraphicsBackend(width, height int) (*GraphicsBackend, error) {
	w, err := window.NewWindow(width, height)
	if err != nil {
		panic(err)
	}

	gb := &GraphicsBackend{
		Window: w,
	}
	if err := gl.Init(); err != nil {
		log.Fatalf("Failed to initialize OpenGL bindings: %s\n", err)
	}
	renderer := NewRenderer()

	renderer.Init()
	gb.Renderer = renderer

	w.SetResizedCallback(func(physicalWidth, physicalHeight uint32) {
		gl.Viewport(0, 0, int32(physicalWidth), int32(physicalHeight))
	})

	w.SetCloseRequestedCallback(func() {
		w.Destroy()
		gb.Close()
	})

	return gb, nil
}

func (backend *GraphicsBackend) Close() {
	backend.Renderer.Destroy()
	backend.Window.Destroy()
}

func (backend *GraphicsBackend) PollEvents() bool {
	return backend.Poll()
}
