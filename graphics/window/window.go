package window

import (
	"github.com/dfirebaugh/banana/pkg/input"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	*glfw.Window
	eventChan  chan input.Event
	isDisposed bool
}

func NewWindow(width, height int) (*Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	// glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	// glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)

	w := &Window{
		eventChan:  make(chan input.Event, 100),
		isDisposed: false,
	}

	win, err := glfw.CreateWindow(640, 480, "sample title", nil, nil)
	if err != nil {
		glfw.Terminate()
		w.DestroyWindow()
		return nil, err
	}

	win.MakeContextCurrent()
	w.Window = win
	win.SetSize(width, height)
	return w, nil
}

func (w *Window) SetBorderlessWindowed(v bool) {
	if v {
		width, height := w.GetSize()
		x, y := w.GetWindowPosition()

		w.SetAttrib(glfw.Decorated, glfw.False)
		w.SetSize(width, height)
		w.SetPos(x, y)
	} else {
		width, height := w.GetSize()
		x, y := w.GetWindowPosition()

		w.SetAttrib(glfw.Decorated, glfw.True)
		w.SetSize(width, height)
		w.SetPos(x, y)
	}
}

func (w *Window) SetFullScreenBorderless(v bool) {
	if v {
		monitor := glfw.GetPrimaryMonitor()
		mode := monitor.GetVideoMode()
		w.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
	} else {
		w.SetMonitor(nil, 100, 100, 800, 600, 0)
	}
}

func (w *Window) GetWindowPosition() (x int, y int) {
	return w.GetPos()
}

func (w *Window) SetAspectRatio(numerator, denominator int) {
	w.Window.SetAspectRatio(numerator, denominator)
}

func (w *Window) DisableWindowResize() {
	if w.Window != nil {
		w.SetAttrib(glfw.Resizable, glfw.False)
	}
}

func (w *Window) SetCloseCallback(fn func()) {
	w.Window.SetCloseCallback(func(w *glfw.Window) {
		defer w.Destroy()
		fn()
	})
}

func (w *Window) SetInputCallback(fn func(eventChan chan input.Event)) {
	w.SetMouseButtonCallback(func(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		var eventType input.EventType
		switch action {
		case glfw.Press:
			eventType = input.MousePress
		case glfw.Release:
			eventType = input.MouseRelease
		}
		w.eventChan <- input.Event{Type: eventType, MouseButton: input.MouseButton(button)}
		fn(w.eventChan)
	})

	w.SetCursorPosCallback(func(window *glfw.Window, xpos, ypos float64) {
		w.eventChan <- input.Event{Type: input.MouseMove, X: int(xpos), Y: int(ypos)}
		fn(w.eventChan)
	})

	w.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		var eventType input.EventType
		switch action {
		case glfw.Press:
			eventType = input.KeyPress
		case glfw.Release:
			eventType = input.KeyRelease
		}
		w.eventChan <- input.Event{Type: eventType, Key: input.Key(key)}
		fn(w.eventChan)
	})
}

func (w *Window) SetResizedCallback(fn func(physicalWidth, physicalHeight uint32)) {
	w.SetSizeCallback(func(window *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		fn(uint32(width), uint32(height))
	})
}

func (w *Window) SetWindowPosition(x, y int) {
	if w.Window != nil {
		w.SetPos(x, y)
	}
}

func (w *Window) SetCloseRequestedCallback(fn func()) {
	w.Window.SetCloseCallback(func(w *glfw.Window) {
		fn()
	})
}

func (w *Window) SetWindowTitle(title string) {
	w.SetTitle(title)
}

func (w *Window) GetSize() (int, int) {
	return w.GetWindowSize()
}

func (w *Window) GetWindowSize() (int, int) {
	if w == nil || w.isDisposed {
		return 0, 0
	}
	return w.Window.GetSize()
}

func (w *Window) SetWindowSize(width int, height int) {
	w.SetSize(width, height)
}

func (w *Window) DestroyWindow() {
	w.Destroy()
}

func (w *Window) Poll() bool {
	if w.isDisposed {
		return false
	}

	glfw.PollEvents()
	return true
}

func (w *Window) IsDisposed() bool {
	return w.isDisposed
}

func (w *Window) Destroy() {
	w.isDisposed = true
	if w.isDisposed {
		return
	}
	if w.Window != nil {
		w.Window.Destroy()
		w.Window = nil
	}
	glfw.Terminate()
}

func (w *Window) ShouldClose() bool {
	return w.Window.ShouldClose()
}
