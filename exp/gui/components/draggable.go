package components

import (
	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/pkg/input"
)

type Draggable struct {
	Node
	isDragging               bool
	dragOffsetX, dragOffsetY int
	onDrag                   func(x, y int)
	onDragStart              func(mouseX, mouseY int)
	onDragEnd                func()
}

func NewDraggable(width, height int, onDrag func(x, y int), onDragStart func(mouseX, mouseY int), onDragEnd func()) *Draggable {
	d := &Draggable{
		onDrag:      onDrag,
		onDragStart: onDragStart,
		onDragEnd:   onDragEnd,
	}
	d.width = width
	d.height = height
	return d
}

func (d *Draggable) Update(globalX, globalY int) {
	if banana.IsButtonPressed(input.MouseButtonLeft) {
		mouseX, mouseY := banana.GetCursorPosition()

		parent := d.GetParent()
		if parent != nil {
			parentX, parentY := parent.GetGlobalOffset()
			parentWidth, parentHeight := parent.GetDimensions()

			if mouseX < parentX || mouseX > parentX+parentWidth || mouseY < parentY || mouseY > parentY+parentHeight {
				d.isDragging = false
				return
			}
		}

		if !d.isDragging {
			if mouseX >= globalX && mouseX <= globalX+d.width &&
				mouseY >= globalY && mouseY <= globalY+d.height {
				d.isDragging = true
				d.dragOffsetX = mouseX - globalX
				d.dragOffsetY = mouseY - globalY
				if d.onDragStart != nil {
					d.onDragStart(mouseX, mouseY)
				}
			}
		} else {
			if d.onDrag != nil {
				d.onDrag(mouseX-d.dragOffsetX, mouseY-d.dragOffsetY)
			}
		}
	} else {
		if d.isDragging {
			d.isDragging = false
			if d.onDragEnd != nil {
				d.onDragEnd()
			}
		}
	}
}

func (d *Draggable) SetDimensions(width, height int) {
	d.width = width
	d.height = height
}
