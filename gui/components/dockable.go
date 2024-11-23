package components

import (
	"github.com/dfirebaugh/banana"
)

type Dockable struct {
	Node
	isDocked      bool
	dockX, dockY  int
	width, height int
	onDock        func(x, y int)
	screenWidth   int
	screenHeight  int
}

func NewDockable(width, height int, onDock func(x, y int), screenWidth, screenHeight int) *Dockable {
	return &Dockable{
		width:        width,
		height:       height,
		onDock:       onDock,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

func (d *Dockable) Update(globalX, globalY int) {
	if d.isDocked {
		if d.onDock != nil {
			d.onDock(d.dockX, d.dockY)
		}
		return
	}

	mouseX, mouseY := banana.GetCursorPosition()

	parent := d.GetParent()
	if parent != nil {
		parentX, parentY := parent.GetGlobalOffset()
		parentWidth, parentHeight := parent.GetDimensions()

		if mouseX < parentX || mouseX > parentX+parentWidth || mouseY < parentY || mouseY > parentY+parentHeight {
			return
		}
	}

	threshold := 20

	if mouseX < threshold {
		d.Dock(0, globalY)
		d.width = d.screenWidth / 3
		d.height = d.screenHeight
	} else if mouseX > d.screenWidth-threshold {
		d.Dock(d.screenWidth-d.width, globalY)
		d.width = d.screenWidth / 3
		d.height = d.screenHeight
	} else if mouseY < threshold {
		d.Dock(globalX, 0)
		d.width = d.screenWidth
		d.height = d.screenHeight / 3
	} else if mouseY > d.screenHeight-threshold {
		d.Dock(globalX, d.screenHeight-d.height)
		d.width = d.screenWidth
		d.height = d.screenHeight / 3
	}
}

func (d *Dockable) Dock(x, y int) {
	d.isDocked = true
	d.dockX = x
	d.dockY = y
}

func (d *Dockable) Undock() {
	d.isDocked = false
}
