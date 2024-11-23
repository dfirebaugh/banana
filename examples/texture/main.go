package main

import (
	"image"
	"net/http"

	"golang.org/x/image/colornames"
	_ "golang.org/x/image/webp"

	"github.com/dfirebaugh/banana"
	"github.com/dfirebaugh/banana/pkg/input"
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

func main() {
	banana.SetWindowSize(600, 412)
	banana.SetTitle("banana.texture example")

	tImage := downloadImage(`https://parade.com/.image/c_limit%2Ccs_srgb%2Cq_auto:good%2Cw_620/MTkwNTgxNDg5MjU4ODY1Nzg5/nick-offerman-donkey-thoughts.webp`)
	t := banana.UploadTexture(
		tImage,
	)
	mountainImage := downloadImage(`https://www.gstatic.com/webp/gallery/1.webp`)
	mountain := banana.UploadTexture(
		mountainImage,
	)

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
	banana.Run(func() {
		exampleControls()
	}, func() {
		banana.Clear(colornames.Black)
		banana.RenderTexture(mountain, &banana.TextureRenderOptions{
			X:          200,
			Y:          200,
			RectWidth:  float32(mountainImage.Bounds().Dx()),
			RectHeight: float32(mountainImage.Bounds().Dy()),
			Width:      float32(mountainImage.Bounds().Dx()),
			Height:     float32(mountainImage.Bounds().Dy()),
			Rotation:   .9,
			Scale:      1,
		})
		banana.RenderTexture(t, &banana.TextureRenderOptions{
			X:          0,
			Y:          0,
			RectWidth:  float32(tImage.Bounds().Dx()),
			RectHeight: float32(tImage.Bounds().Dy()),
			Width:      float32(tImage.Bounds().Dx()),
			Height:     float32(tImage.Bounds().Dy()),
			Scale:      .4,
			// Rotation:   .9,
		})
	})
}
