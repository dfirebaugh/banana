package main

import (
	"bytes"
	"fmt"
	"image"
	"time"

	_ "image/png"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/assets"
	"github.com/dfirebaugh/banana/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	windowWidth        = 800
	windowHeight       = 600
	playerSpeed        = 800.0  // pixels per second
	gravity            = 3200.0 // pixels per second^2
	jumpSpeed          = 1300.0 // pixels per second
	coyoteTimeDuration = 0.2    // seconds
	debug              = false
)

type Player struct {
	X              float32
	Y              float32
	W              float32
	H              float32
	VelY           float32
	Ground         bool
	TextureID      uint32
	FrameSize      image.Point
	SheetSize      image.Point
	FrameIndex     int
	LastFrame      time.Time
	LastUpdateTime time.Time
	platforms      []*Platform

	CoyoteTimeLeft float32
	Rect           *banana.Rect
}

func (p *Player) handleMovement(deltaTime float32) {
	if banana.IsKeyPressed(input.KeyA) || banana.IsKeyPressed(input.KeyLeft) {
		p.X -= playerSpeed * deltaTime
	}
	if banana.IsKeyPressed(input.KeyD) || banana.IsKeyPressed(input.KeyRight) {
		p.X += playerSpeed * deltaTime
	}
}

func (p *Player) handleCoyoteTime(deltaTime float32) {
	if p.Ground {
		p.CoyoteTimeLeft = coyoteTimeDuration
	} else {
		p.CoyoteTimeLeft -= deltaTime
		if p.CoyoteTimeLeft < 0 {
			p.CoyoteTimeLeft = 0
		}
	}
}

func (p *Player) handlePlatformCollision() {
	p.Ground = false

	playerLeft := p.X
	playerRight := p.X + p.W
	playerBottom := p.Y + p.H

	for _, pl := range p.platforms {
		platformLeft := pl.X
		platformRight := pl.X + pl.W
		platformTop := pl.Y
		platformBottom := pl.Y + pl.H

		if playerRight > platformLeft && playerLeft < platformRight {
			if playerBottom >= platformTop && playerBottom <= platformBottom && p.VelY >= 0 {
				p.Y = platformTop - p.H
				p.VelY = 0
				p.Ground = true
				break
			}
		}
	}
}

func (p *Player) handleGroundCollision() {
	if p.Y > float32(windowHeight)-p.H {
		p.Y = float32(windowHeight) - p.H
		p.VelY = 0
		p.Ground = true
	}
}

func (p *Player) updateVelocity(deltaTime float32) {
	if p.CoyoteTimeLeft <= 0 {
		p.VelY += gravity * deltaTime
	}
	p.Y += p.VelY * deltaTime
}

func (p *Player) updateSpriteFrame() {
	if time.Since(p.LastFrame) >= time.Millisecond*200 {
		p.LastFrame = time.Now()
		p.FrameIndex = (p.FrameIndex + 1) % (p.SheetSize.X * p.SheetSize.Y)
	}
}

func (p *Player) handleJump() {
	if banana.IsKeyPressed(input.KeySpace) && (p.Ground || p.CoyoteTimeLeft > 0) {
		p.VelY = -jumpSpeed
		p.Ground = false
		p.CoyoteTimeLeft = 0
	}
}

func (p *Player) Update(deltaTime float32) {
	p.updateVelocity(deltaTime)
	p.handleGroundCollision()
	p.handleMovement(deltaTime)
	p.handleCoyoteTime(deltaTime)

	p.handleJump()
	p.updateSpriteFrame()

	p.handlePlatformCollision()
	p.Rect.X = p.X
	p.Rect.Y = p.Y
}

func (p *Player) Render() {
	frameX := (p.FrameIndex % p.SheetSize.X) * p.FrameSize.X
	frameY := (p.FrameIndex / p.SheetSize.X) * p.FrameSize.Y

	options := &banana.TextureRenderOptions{
		X:          p.X,
		Y:          p.Y,
		RectWidth:  float32(p.FrameSize.X),
		RectHeight: float32(p.FrameSize.Y),
		Scale:      float32(p.W) / float32(p.FrameSize.X),
		RectX:      float32(frameX),
		RectY:      float32(frameY),
		Width:      float32(p.FrameSize.X),
		Height:     float32(p.FrameSize.Y),
	}

	banana.RenderTexture(p.TextureID, options)

	if debug {
		banana.RenderShape(p.Rect)
	}
}

func main() {
	banana.SetWindowSize(windowWidth, windowHeight)
	if debug {
		banana.EnableFPS()
	}

	reader := bytes.NewReader(assets.BuddyDanceSpriteSheet)
	img, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}
	textureID := banana.UploadTexture(img)
	frameSize := image.Point{X: 32, Y: 32}
	sheetSize := image.Point{X: 4, Y: 1}

	platforms := []*Platform{
		NewPlatform(200, 400, 100, 20),
		NewPlatform(400, 300, 150, 20),
	}

	player := &Player{
		X:              100,
		Y:              float32(windowHeight) - 100,
		W:              64,
		H:              64,
		TextureID:      textureID,
		FrameSize:      frameSize,
		SheetSize:      sheetSize,
		LastFrame:      time.Now(),
		LastUpdateTime: time.Now(),
		platforms:      platforms,
		Rect: &banana.Rect{
			X:      100,
			Y:      float32(windowHeight) - 100,
			Width:  64,
			Height: 64,
			Color:  colornames.Mediumpurple,
		},
	}

	banana.Run(func() {
		now := time.Now()
		deltaTime := now.Sub(player.LastUpdateTime).Seconds()
		player.LastUpdateTime = now

		player.Update(float32(deltaTime))
	}, func() {
		banana.Clear(colornames.White)
		player.Render()

		for _, pl := range platforms {
			pl.Render()
		}

		banana.RenderText(fmt.Sprintf("Player X: %d Y: %d", int(player.X), int(player.Y)),
			&banana.TextRenderOptions{
				X:     10,
				Y:     windowHeight - 20,
				Size:  16,
				Color: colornames.Black,
			})
		banana.RenderText(fmt.Sprintf("VelY: %.2f, Ground: %t, CoyoteTimeLeft: %.2f", player.VelY, player.Ground, player.CoyoteTimeLeft),
			&banana.TextRenderOptions{
				X:     10,
				Y:     windowHeight - 40,
				Size:  16,
				Color: colornames.Black,
			})
	})
}

type Platform struct {
	X    float32
	Y    float32
	W    float32
	H    float32
	Rect *banana.Rect
}

func NewPlatform(x, y, w, h float32) *Platform {
	rect := &banana.Rect{
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
		Color:  colornames.Royalblue,
	}
	return &Platform{X: x, Y: y, W: w, H: h, Rect: rect}
}

func (pl *Platform) Render() {
	banana.RenderShape(pl.Rect)
}
