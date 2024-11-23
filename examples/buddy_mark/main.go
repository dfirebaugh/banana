package main

import (
	"bytes"
	"fmt"
	"image"
	"math/rand"
	"net/http"
	"time"

	_ "image/png"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/assets"
	"github.com/dfirebaugh/banana/pkg/input"
	"golang.org/x/image/colornames"
)

const (
	gravity      = 0.1
	damping      = 0.9
	buddyWidth   = 32
	buddyHeight  = 32
	screenWidth  = 800
	screenHeight = 600
)

// downloadImage fetches the image from the given URL and returns it as an image.Image
func downloadImage(url string) image.Image {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		panic(err)
	}

	return img
}

type Buddy struct {
	X, Y                 float32
	VelocityX, VelocityY float32
	TextureID            uint32
	Frame                int
}

var (
	buddies   []*Buddy
	img       image.Image
	textureID uint32
)

func main() {
	var err error
	banana.SetWindowSize(screenWidth, screenHeight)
	banana.SetTitle("BuddyMark Stress Test")
	banana.EnableFPS()

	isFullScreen := false
	exampleControls := func() {
		if banana.IsKeyJustPressed(input.KeyA) {
			isFullScreen = !isFullScreen
			banana.SetBorderlessWindowed(isFullScreen)
		}
		if banana.IsKeyJustPressed(input.KeyEscape) {
			banana.Close()
		}
	}
	reader := bytes.NewReader(assets.BuddyDanceSpriteSheet)
	img, _, err = image.Decode(reader)
	if err != nil {
		panic(err)
	}

	textureID = banana.UploadTexture(img)

	lastFrameTime := time.Now()
	frameDuration := time.Millisecond * 100
	// backgroundImg := downloadImage(`https://www.gstatic.com/webp/gallery/1.webp`)

	// background, _ := banana.CreateTextureFromImage(
	// 	backgroundImg,
	// )
	// rq := banana.CreateRenderQueue()
	// rq.SetPriority(50)
	// background.RenderToQueue(rq)
	//
	// background.Resize(screenWidth, screenHeight)
	// background.Move(0, 0)

	banana.Run(func() {
		exampleControls()
		if time.Since(lastFrameTime) >= frameDuration {
			lastFrameTime = time.Now()
			for _, buddy := range buddies {
				buddy.Frame = (buddy.Frame + 1) % 4
			}
		}

		handleInput()

		for _, buddy := range buddies {
			buddy.Update()
		}
	}, func() {
		banana.Clear(colornames.Skyblue)
		for _, buddy := range buddies {
			buddy.Render()
		}
		banana.RenderText(fmt.Sprintf("buddies: %d", len(buddies)), &banana.TextRenderOptions{
			X:     20,
			Y:     20,
			Size:  12,
			Color: colornames.Red,
		},
		)
	})
}

func handleInput() {
	if banana.IsButtonPressed(input.MouseButtonLeft) {
		if banana.GetFPS() < 40 {
			return
		}
		x, y := banana.GetCursorPosition()
		for i := 0; i < 5; i++ {
			buddy := NewBuddy(float32(x), float32(y))
			buddies = append(buddies, buddy)
		}
	}

	if banana.IsButtonPressed(input.MouseButtonRight) {
		buddies = []*Buddy{}
	}
}

func NewBuddy(x, y float32) *Buddy {
	return &Buddy{
		X:         x,
		Y:         y,
		VelocityX: rand.Float32()*10 - 5,
		VelocityY: rand.Float32()*10 - 5,
		TextureID: textureID,
		Frame:     0,
	}
}

func (b *Buddy) Update() {
	b.VelocityY += gravity

	b.X += b.VelocityX
	b.Y += b.VelocityY

	if b.X < 0 {
		b.X = 0
		b.VelocityX *= -damping
	} else if b.X > screenWidth-buddyWidth {
		b.X = screenWidth - buddyWidth
		b.VelocityX *= -damping
	}

	if b.Y > screenHeight-buddyHeight {
		b.Y = screenHeight - buddyHeight
		b.VelocityY *= -damping
	}
}

func (b *Buddy) Render() {
	frameX := float32(b.Frame * buddyWidth)
	banana.RenderTexture(b.TextureID, &banana.TextureRenderOptions{
		X:          b.X,
		Y:          b.Y,
		RectX:      frameX,
		RectY:      0,
		RectWidth:  buddyWidth,
		RectHeight: buddyHeight,
		Scale:      2.0,
		Width:      float32(img.Bounds().Dx()),
		Height:     float32(img.Bounds().Dy()),
		FlipX:      false,
		FlipY:      false,
		Rotation:   0,
	})
}
