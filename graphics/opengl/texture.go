package opengl

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/sirupsen/logrus"
)

type TextureAtlas struct {
	ID     uint32
	Width  int
	Height int
	Image  *image.RGBA
}

func (atlas *TextureAtlas) canFit(width, height int) bool {
	return width <= atlas.Width && height <= atlas.Height
}

func (atlas *TextureAtlas) grow(newImgWidth, newImgHeight int) {
	requiredWidth := atlas.Width + newImgWidth
	requiredHeight := atlas.Height + newImgHeight

	newWidth := nextPowerOfTwo(requiredWidth)
	newHeight := nextPowerOfTwo(requiredHeight)

	if newWidth <= atlas.Width {
		newWidth = atlas.Width * 2
	}
	if newHeight <= atlas.Height {
		newHeight = atlas.Height * 2
	}

	newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.Draw(newImage, atlas.Image.Bounds(), atlas.Image, image.Point{}, draw.Src)

	atlas.Width = newWidth
	atlas.Height = newHeight
	atlas.Image = newImage
}

func nextPowerOfTwo(n int) int {
	if n <= 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

func (atlas *TextureAtlas) findPlace(width, height int) (int, int) {
	for y := 0; y <= atlas.Height-height; y++ {
		for x := 0; x <= atlas.Width-width; x++ {
			if atlas.isAreaFree(x, y, width, height) {
				return x, y
			}
		}
	}
	return -1, -1
}

func (atlas *TextureAtlas) isAreaFree(x, y, width, height int) bool {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if atlas.Image.At(x+j, y+i) != (color.RGBA{0, 0, 0, 0}) {
				return false
			}
		}
	}
	return true
}

func (atlas *TextureAtlas) updateTexture() {
	if atlas.ID == 0 {
		gl.GenTextures(1, &atlas.ID)
	}
	gl.BindTexture(gl.TEXTURE_2D, atlas.ID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(atlas.Width), int32(atlas.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(atlas.Image.Pix))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

type TextureManager struct {
	atlas         *TextureAtlas
	renderer      *Renderer
	textureBounds map[uint32]image.Rectangle
	textureIDs    map[uint32]uint32
	glTextureIDs  map[uint32]image.Rectangle
}

func NewTextureManager(renderer *Renderer) *TextureManager {
	initialSize := 512
	atlas := &TextureAtlas{
		Width:  initialSize,
		Height: initialSize,
		Image:  image.NewRGBA(image.Rect(0, 0, initialSize, initialSize)),
	}
	return &TextureManager{
		atlas:         atlas,
		renderer:      renderer,
		textureBounds: make(map[uint32]image.Rectangle),
		textureIDs:    make(map[uint32]uint32),
		glTextureIDs:  make(map[uint32]image.Rectangle),
	}
}

func (tm *TextureManager) RegisterTexture(textureID uint32) uint32 {
	managerID := uint32(len(tm.textureBounds) + 1)
	tm.textureIDs[textureID] = managerID
	return managerID
}

func (tm *TextureManager) GetManagerID(textureID uint32) (uint32, bool) {
	managerID, exists := tm.textureIDs[textureID]
	return managerID, exists
}

func (tm *TextureManager) RegisterFramebufferTexture(textureID uint32, width, height int) {
	textureID += 2
	tm.textureBounds[textureID] = image.Rect(0, 0, width, height)
	tm.textureIDs[textureID] = textureID                         // Direct mapping for framebuffer textures
	tm.glTextureIDs[textureID] = image.Rect(0, 0, width, height) // Store the location in the atlas
	logrus.Infof("Registered framebuffer texture with ID %d, width %d, height %d", textureID, width, height)
}

func (tm *TextureManager) UploadTexture(img image.Image) uint32 {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	img = flipImageVertically(img)

	x, y := tm.atlas.findPlace(width, height)

	if x == -1 || y == -1 {
		tm.atlas.grow(width, height)
		x, y = tm.atlas.findPlace(width, height)
		if x == -1 || y == -1 {
			logrus.Error("Failed to find place for new texture after growing the atlas")
			return 0
		}
	}

	draw.Draw(tm.atlas.Image, image.Rect(x, y, x+width, y+height), img, bounds.Min, draw.Src)

	tm.atlas.updateTexture()

	textureID := uint32(len(tm.textureBounds) + 1)
	tm.textureBounds[textureID] = image.Rect(x, y, x+width, y+height)
	tm.glTextureIDs[textureID] = image.Rect(x, y, x+width, y+height)

	logrus.Infof("Uploaded texture with ID %d at position (%d, %d)", textureID, x, y)

	return textureID
}

func (tm *TextureManager) UpdateTexture(textureID uint32, img image.Image, xOffset, yOffset int) {
	bounds := img.Bounds()
	existingBounds, exists := tm.textureBounds[textureID]
	if !exists {
		logrus.Error("Texture ID not found")
		return
	}

	draw.Draw(tm.atlas.Image, existingBounds.Add(image.Pt(xOffset, yOffset)), img, bounds.Min, draw.Src)

	tm.atlas.updateTexture()
}

func flipImageVertically(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		rgbaImg = image.NewRGBA(bounds)
		draw.Draw(rgbaImg, bounds, img, image.Point{}, draw.Src)
	}

	height := bounds.Dy()
	stride := rgbaImg.Stride
	temp := make([]byte, stride)

	for y := 0; y < height/2; y++ {
		top := rgbaImg.Pix[y*stride : (y+1)*stride]
		bottom := rgbaImg.Pix[(height-y-1)*stride : (height-y)*stride]
		copy(temp, top)
		copy(top, bottom)
		copy(bottom, temp)
	}

	return rgbaImg
}
