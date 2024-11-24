package components

import (
	"github.com/dfirebaugh/banana/exp/gui"
)

type Component interface {
	Update()
	Render(ctx gui.DrawContext)
	SetParent(parent *Node)
	GetParent() *Node
	GetGlobalOffset() (int, int)
	GetDimensions() (int, int)
}

type Node struct {
	offsetX, offsetY int
	width, height    int
	children         []Component
	parent           *Node
}

func (n *Node) GetOffset() (int, int) {
	return n.offsetX, n.offsetY
}

func (n *Node) SetOffset(x int, y int) {
	n.offsetX = x
	n.offsetY = y
}

func (n *Node) AppendChild(child Component) {
	child.SetParent(n)
	n.children = append(n.children, child)
}

func (n *Node) Update() {
	for _, c := range n.children {
		c.Update()
	}
}

func (n *Node) Render(ctx gui.DrawContext) {
	for _, c := range n.children {
		c.Render(ctx)
	}
}

func (n *Node) SetParent(parent *Node) {
	n.parent = parent
}

func (n *Node) GetParent() *Node {
	return n.parent
}

func (n *Node) GetGlobalOffset() (int, int) {
	if n.parent == nil {
		return n.offsetX, n.offsetY
	}
	parentX, parentY := n.parent.GetGlobalOffset()
	return n.offsetX + parentX, n.offsetY + parentY
}

func (n *Node) GetGlobalOffsetX() int {
	globalX, _ := n.GetGlobalOffset()
	return globalX
}

func (n *Node) GetGlobalOffsetY() int {
	_, globalY := n.GetGlobalOffset()
	return globalY
}

func (n *Node) GetDimensions() (int, int) {
	return n.width, n.height
}

func (n *Node) IsWithinBounds(mouseX, mouseY int) bool {
	if n.parent == nil {
		return true
	}
	parentX, parentY := n.parent.GetGlobalOffset()
	parentWidth, parentHeight := n.parent.GetDimensions()

	return mouseX >= parentX && mouseX <= parentX+parentWidth && mouseY >= parentY && mouseY <= parentY+parentHeight
}
