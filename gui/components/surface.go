package components

import (
	"fmt"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/gui"
)

type surfaceHandle struct {
	Node
	*Draggable
	surface *Surface
}

func newSurfaceHandle(width, height int, onDrag func(x, y int), surface *Surface) *surfaceHandle {
	h := &surfaceHandle{
		surface: surface,
	}
	h.Node.width = width
	h.Node.height = height
	h.Draggable = NewDraggable(width, height, onDrag, h.OnDragStart, surface.OnDragEnd)
	return h
}

func (h *surfaceHandle) OnDragStart(mouseX, mouseY int) {
	h.surface.OnDragStart(mouseX, mouseY)
}

func (h *surfaceHandle) Render(ctx gui.DrawContext) {
	if h.surface.IsLocked {
		return
	}
	globalX, globalY := h.GetOffset()
	ctx.DrawRectangle(globalX, globalY, h.width, h.height-5, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().HandleColor,
			OutlineColor: ctx.GetTheme().PrimaryColor,
			OutlineSize:  0,
			CornerRadius: 5,
		},
	})
	ctx.DrawRectangle(globalX, globalY+10, h.width, h.height-10, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().HandleColor,
			OutlineColor: ctx.GetTheme().PrimaryColor,
			OutlineSize:  0,
			CornerRadius: 0,
		},
	})
}

func (h *surfaceHandle) Update() {
	if h.surface.IsLocked {
		return
	}
	globalX, globalY := h.GetGlobalOffset()
	h.Draggable.Update(globalX, globalY)
}

func (h *surfaceHandle) SetDimensions(width, height int) {
	h.width = width
	h.height = height
	h.Draggable.SetDimensions(width, height)
}

type surfaceBorder struct {
	Node
	*Draggable
	width, height int
	surface       *Surface
}

func newSurfaceBorder(width, height int, onDrag func(x, y int), surface *Surface) *surfaceBorder {
	return &surfaceBorder{
		width:     width,
		height:    height,
		Draggable: NewDraggable(width, height, onDrag, nil, nil),
		surface:   surface,
	}
}

func (b *surfaceBorder) Render(ctx gui.DrawContext) {
	if true {
		return
	}
	globalX, globalY := b.GetGlobalOffset()
	ctx.DrawRectangle(globalX, globalY, b.width, b.height, &gui.DrawOptions{
		Style: gui.Style{
			FillColor:    ctx.GetTheme().HandleColor,
			OutlineColor: ctx.GetTheme().PrimaryColor,
			OutlineSize:  0,
			CornerRadius: 0,
		},
	})
}

func (b *surfaceBorder) Update() {
	if b.surface.IsLocked {
		return
	}
	globalX, globalY := b.GetGlobalOffset()
	b.Draggable.Update(globalX, globalY)
}

func (b *surfaceBorder) SetDimensions(width, height int) {
	b.width = width
	b.height = height
	b.Draggable.SetDimensions(width, height)
}

type Surface struct {
	Node
	*surfaceHandle
	originalWidth    int
	originalHeight   int
	isDocked         bool
	dockedEdge       string
	prevWindowWidth  int
	prevWindowHeight int
	borders          []*surfaceBorder
	IsLocked         bool
	framebuffer      banana.Framebuffer
}

func NewSurface(x, y, width, height int) *Surface {
	s := &Surface{
		originalWidth:  width,
		originalHeight: height,
	}
	s.Node.width = width
	s.Node.height = height
	s.SetOffset(x, y)

	s.prevWindowWidth, s.prevWindowHeight = banana.GetWindowSize()

	s.surfaceHandle = newSurfaceHandle(width, 20, s.SetOffset, s)

	s.surfaceHandle.SetParent(&s.Node)
	s.surfaceHandle.SetOffset(0, 0)
	s.AppendChild(s.surfaceHandle)

	s.borders = []*surfaceBorder{
		newSurfaceBorder(width, 5, s.ResizeTop, s),
		newSurfaceBorder(width, 5, s.ResizeBottom, s),
		newSurfaceBorder(5, height, s.ResizeLeft, s),
		newSurfaceBorder(5, height, s.ResizeRight, s),
	}

	s.borders[0].SetParent(&s.Node)
	s.borders[0].SetOffset(0, -5)
	s.AppendChild(s.borders[0])

	s.borders[1].SetParent(&s.Node)
	s.borders[1].SetOffset(0, s.height)
	s.AppendChild(s.borders[1])

	s.borders[2].SetParent(&s.Node)
	s.borders[2].SetOffset(-5, 0)
	s.AppendChild(s.borders[2])

	s.borders[3].SetParent(&s.Node)
	s.borders[3].SetOffset(s.width, 0)
	s.AppendChild(s.borders[3])

	fb, err := banana.AddFramebuffer(width, height)
	if err != nil {
		fmt.Println("Error creating framebuffer:", err)
	}
	s.framebuffer = fb

	return s
}

func (s *Surface) PrepareRender(ctx gui.DrawContext) {
	banana.BindFramebuffer(s.framebuffer)
	banana.Clear(ctx.GetTheme().BackgroundColor)
	for _, child := range s.children {
		child.Render(ctx)
	}
	if !s.IsLocked {
		s.surfaceHandle.Render(ctx)
		for _, border := range s.borders {
			border.Render(ctx)
		}
	}
	s.framebuffer.Draw(s.offsetX, s.offsetY, s.width, s.height)
	banana.UnbindFramebuffer()
}

func (s *Surface) Render(ctx gui.DrawContext) {
	s.update()
	banana.RenderShape(&banana.Rect{
		X:      float32(s.offsetX),
		Y:      float32(s.offsetY),
		Width:  float32(s.width),
		Height: float32(s.height),
		Color:  ctx.GetTheme().BackgroundColor,
	})
	globalX, globalY := s.GetGlobalOffset()

	banana.RenderFramebuffer(s.framebuffer, &banana.TextureRenderOptions{
		X:             float32(globalX),
		Y:             float32(globalY),
		Width:         float32(s.width),
		Height:        float32(s.height),
		RectWidth:     float32(s.width),
		RectHeight:    float32(s.height),
		DesiredWidth:  float32(s.width),
		DesiredHeight: float32(s.height),
	})
}

func (s *Surface) update() {
	for _, child := range s.children {
		child.Update()
	}

	s.surfaceHandle.SetDimensions(s.width, 20)

	s.borders[0].SetDimensions(s.width, 5)
	s.borders[0].SetOffset(0, -5)
	s.borders[1].SetDimensions(s.width, 5)
	s.borders[1].SetOffset(0, s.height)
	s.borders[2].SetDimensions(5, s.height)
	s.borders[2].SetOffset(-5, 0)
	s.borders[3].SetDimensions(5, s.height)
	s.borders[3].SetOffset(s.width, 0)

	currentWindowWidth, currentWindowHeight := banana.GetWindowSize()
	if currentWindowWidth != s.prevWindowWidth || currentWindowHeight != s.prevWindowHeight {
		s.OnWindowResize(currentWindowWidth, currentWindowHeight)
		s.prevWindowWidth = currentWindowWidth
		s.prevWindowHeight = currentWindowHeight
	}
}

func (s *Surface) OnDragStart(mouseX, mouseY int) {
	if s.IsLocked {
		return
	}
	if s.isDocked {
		s.isDocked = false
		s.dockedEdge = ""
		s.Resize(s.originalWidth, s.originalHeight)
	}

	handleX, handleY := s.surfaceHandle.GetGlobalOffset()
	handleWidth := s.surfaceHandle.width
	handleHeight := s.surfaceHandle.height

	if mouseX >= handleX && mouseX <= handleX+handleWidth &&
		mouseY >= handleY && mouseY <= handleY+handleHeight {
		s.surfaceHandle.Draggable.isDragging = true
		s.surfaceHandle.Draggable.dragOffsetX = mouseX - handleX
		s.surfaceHandle.Draggable.dragOffsetY = mouseY - handleY
	} else {
		s.SetOffset(mouseX-s.width/2, mouseY-s.surfaceHandle.height/2)

		handleX, handleY = s.surfaceHandle.GetGlobalOffset()

		s.surfaceHandle.Draggable.isDragging = true
		s.surfaceHandle.Draggable.dragOffsetX = mouseX - handleX
		s.surfaceHandle.Draggable.dragOffsetY = mouseY - handleY
	}
}

func (s *Surface) OnDragEnd() {
	if s.IsLocked {
		return
	}
	screenWidth, screenHeight := banana.GetWindowSize()
	radius := 50

	mouseX, mouseY := banana.GetCursorPosition()
	globalX, globalY := s.GetGlobalOffset()

	if mouseX < radius {
		s.Dock(0, globalY)
		s.Resize(s.width, screenHeight)
		s.dockedEdge = "left"
	} else if mouseX > screenWidth-radius {
		s.Dock(screenWidth-s.width, globalY)
		s.Resize(s.width, screenHeight)
		s.dockedEdge = "right"
	} else if mouseY < radius {
		s.Dock(globalX, 0)
		s.Resize(screenWidth, s.height)
		s.dockedEdge = "top"
	} else if mouseY > screenHeight-radius {
		s.Dock(globalX, screenHeight-s.height)
		s.Resize(screenWidth, s.height)
		s.dockedEdge = "bottom"
	} else {
		s.isDocked = false
		s.dockedEdge = ""
		s.Resize(s.originalWidth, s.originalHeight)
	}

	s.Reposition()
}

func (s *Surface) Dock(x, y int) {
	s.isDocked = true
	s.SetOffset(x, y)
}

func (s *Surface) Undock() {
	s.isDocked = false
	s.Resize(s.originalWidth, s.originalHeight)
	s.Reposition()
}

func (s *Surface) Resize(width, height int) {
	s.width = width
	s.height = height

	s.surfaceHandle.SetDimensions(width, 20)

	s.borders[0].SetDimensions(width, 5)
	s.borders[1].SetDimensions(width, 5)
	s.borders[2].SetDimensions(5, height)
	s.borders[3].SetDimensions(5, height)

	// Resize framebuffer
	banana.ResizeFramebuffer(s.framebuffer, width, height)
}

func (s *Surface) Reposition() {
	screenWidth, screenHeight := banana.GetWindowSize()
	globalX, globalY := s.GetGlobalOffset()

	if globalX < 0 {
		s.SetOffset(0, globalY)
	} else if globalX+s.width > screenWidth {
		s.SetOffset(screenWidth-s.width, globalY)
	}

	if globalY < 0 {
		s.SetOffset(globalX, 0)
	} else if globalY+s.height > screenHeight {
		s.SetOffset(globalX, screenHeight-s.height)
	}
}

func (s *Surface) ResizeTop(x, y int) {
	if s.IsLocked {
		return
	}
	newHeight := s.height + (s.GetGlobalOffsetY() - y)
	if newHeight > 20 {
		s.SetOffset(s.GetGlobalOffsetX(), y)
		s.Resize(s.width, newHeight)
	}
}

func (s *Surface) ResizeBottom(x, y int) {
	if s.IsLocked {
		return
	}
	newHeight := y - s.GetGlobalOffsetY()
	if newHeight > 20 {
		s.Resize(s.width, newHeight)
	}
}

func (s *Surface) ResizeLeft(x, y int) {
	if s.IsLocked {
		return
	}
	newWidth := s.width + (s.GetGlobalOffsetX() - x)
	if newWidth > 20 {
		s.SetOffset(x, s.GetGlobalOffsetY())
		s.Resize(newWidth, s.height)
	}
}

func (s *Surface) ResizeRight(x, y int) {
	if s.IsLocked {
		return
	}
	newWidth := x - s.GetGlobalOffsetX()
	if newWidth > 20 {
		s.Resize(newWidth, s.height)
	}
}

func (s *Surface) OnWindowResize(newWindowWidth, newWindowHeight int) {
	if s.isDocked {
		switch s.dockedEdge {
		case "left":
			s.SetOffset(0, 0)
			s.Resize(s.width, newWindowHeight)
		case "right":
			s.SetOffset(newWindowWidth-s.width, 0)
			s.Resize(s.width, newWindowHeight)
		case "top":
			s.SetOffset(0, 0)
			s.Resize(newWindowWidth, s.height)
		case "bottom":
			s.SetOffset(0, newWindowHeight-s.height)
			s.Resize(newWindowWidth, s.height)
		}
	}
}

func (s *Surface) DockLeft() {
	_, screenHeight := banana.GetWindowSize()
	s.Dock(0, 0)
	s.Resize(s.width, screenHeight)
	s.dockedEdge = "left"
}

func (s *Surface) DockRight() {
	screenWidth, screenHeight := banana.GetWindowSize()
	s.Dock(screenWidth-s.width, 0)
	s.Resize(s.width, screenHeight)
	s.dockedEdge = "right"
}

func (s *Surface) DockTop() {
	screenWidth, _ := banana.GetWindowSize()
	s.Dock(0, 0)
	s.Resize(screenWidth, s.height)
	s.dockedEdge = "top"
}

func (s *Surface) DockBottom() {
	screenWidth, screenHeight := banana.GetWindowSize()
	s.Dock(0, screenHeight-s.height)
	s.Resize(screenWidth, s.height)
	s.dockedEdge = "bottom"
}
