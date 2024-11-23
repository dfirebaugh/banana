package font

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"
)

const MaxLoadedGlyphs = 96

type Glyph struct {
	Rune         rune
	TexCoords    [4]float32
	AdvanceWidth float32
	BearingX     float32
	BearingY     float32
	SizeWidth    float32
	SizeHeight   float32
}

type Font struct {
	sfntFont    *sfnt.Font
	Atlas       uint32
	Glyphs      map[rune]*Glyph
	AtlasWidth  int
	AtlasHeight int
	atlasImage  *image.RGBA
}

func (font *Font) Destroy() {
	gl.DeleteTextures(1, &font.Atlas)
}

func (f *Font) Image() image.Image {
	return f.atlasImage
}

func (f *Font) GetKerning(r1, r2 rune, ppem fixed.Int26_6) (fixed.Int26_6, error) {
	var buf sfnt.Buffer
	g1, err := f.sfntFont.GlyphIndex(&buf, r1)
	if err != nil {
		return 0, fmt.Errorf("failed to get glyph index for rune '%c': %v", r1, err)
	}
	g2, err := f.sfntFont.GlyphIndex(&buf, r2)
	if err != nil {
		return 0, fmt.Errorf("failed to get glyph index for rune '%c': %v", r2, err)
	}
	kern, err := f.sfntFont.Kern(&buf, g1, g2, ppem, font.HintingNone)
	if err != nil {
		return 0, fmt.Errorf("failed to get kerning between '%c' and '%c': %v", r1, r2, err)
	}
	return kern, nil
}

func LoadFont(fontData []byte) (*Font, error) {
	sfntFont, err := sfnt.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %v", err)
	}

	f := &Font{
		sfntFont:    sfntFont,
		Glyphs:      make(map[rune]*Glyph),
		AtlasWidth:  512,
		AtlasHeight: 512,
	}

	atlasImg := image.NewRGBA(image.Rect(0, 0, f.AtlasWidth, f.AtlasHeight))
	draw.Draw(atlasImg, atlasImg.Bounds(), &image.Uniform{C: color.RGBA{0, 0, 0, 0}}, image.Point{}, draw.Src)

	var x, y, rowHeight int
	const padding = 1

	for r := rune(32); r < 32+MaxLoadedGlyphs; r++ {
		if r > 126 {
			break
		}

		var buf sfnt.Buffer

		glyphIndex, err := f.sfntFont.GlyphIndex(&buf, r)
		if err != nil || glyphIndex == 0 {
			log.Printf("Glyph not found for rune '%c'", r)
			continue
		}

		ppem := fixed.Int26_6(32 << 6)

		advance, err := f.sfntFont.GlyphAdvance(&buf, glyphIndex, ppem, font.HintingNone)
		if err != nil {
			log.Printf("Failed to get advance for rune '%c': %v", r, err)
			continue
		}

		bounds, _, err := f.sfntFont.GlyphBounds(&buf, glyphIndex, ppem, font.HintingNone)
		if err != nil {
			log.Printf("Failed to get bounds for rune '%c': %v", r, err)
			continue
		}

		segments, err := f.sfntFont.LoadGlyph(&buf, glyphIndex, ppem, nil)
		if err != nil {
			log.Printf("Failed to load glyph '%c': %v", r, err)
			continue
		}

		rWidth, rHeight, img, err := rasterizeGlyph(segments)
		if err != nil {
			log.Printf("Failed to rasterize glyph '%c': %v", r, err)
			continue
		}

		if x+rWidth+padding > f.AtlasWidth {
			x = 0
			y += rowHeight + padding
			rowHeight = 0
		}

		if y+rHeight+padding > f.AtlasHeight {
			return nil, fmt.Errorf("atlas size too small to fit all glyphs")
		}

		glyphRect := image.Rect(x, y, x+rWidth, y+rHeight)
		draw.Draw(atlasImg, glyphRect, img, image.Point{0, 0}, draw.Over)

		advanceWidth := float32(advance.Round())
		if r == ' ' {
			// how wide are spaces?
			advanceWidth = 5.0
		}

		glyph := &Glyph{
			Rune: r,
			TexCoords: [4]float32{
				float32(x) / float32(f.AtlasWidth), float32(y) / float32(f.AtlasHeight),
				float32(x+rWidth) / float32(f.AtlasWidth), float32(y+rHeight) / float32(f.AtlasHeight),
			},
			AdvanceWidth: advanceWidth,
			BearingX:     float32(bounds.Min.X.Round()),
			BearingY:     float32(bounds.Max.Y.Round()),
			SizeWidth:    float32(rWidth),
			SizeHeight:   float32(rHeight),
		}

		f.Glyphs[r] = glyph

		x += rWidth + padding
		if rHeight > rowHeight {
			rowHeight = rHeight
		}
	}

	f.atlasImage = atlasImg
	return f, nil
}

func rasterizeGlyph(segments []sfnt.Segment) (width, height int, img *image.RGBA, err error) {
	var minX, minY, maxX, maxY fixed.Int26_6
	for _, seg := range segments {
		for _, arg := range seg.Args {
			if arg.X < minX {
				minX = arg.X
			}
			if arg.Y < minY {
				minY = arg.Y
			}
			if arg.X > maxX {
				maxX = arg.X
			}
			if arg.Y > maxY {
				maxY = arg.Y
			}
		}
	}

	minXf := float32(minX) / 64.0
	minYf := float32(minY) / 64.0
	maxXf := float32(maxX) / 64.0
	maxYf := float32(maxY) / 64.0

	width = int(math.Ceil(float64(maxXf - minXf)))
	height = int(math.Ceil(float64(maxYf - minYf)))

	if width == 0 || height == 0 {
		width, height = 1, 1
	}

	raster := vector.NewRasterizer(width, height)
	raster.DrawOp = draw.Src

	translateX := -minXf
	translateY := -minYf

	for _, seg := range segments {
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			raster.MoveTo(float32(seg.Args[0].X)/64.0+translateX, float32(seg.Args[0].Y)/64.0+translateY)
		case sfnt.SegmentOpLineTo:
			raster.LineTo(float32(seg.Args[0].X)/64.0+translateX, float32(seg.Args[0].Y)/64.0+translateY)
		case sfnt.SegmentOpQuadTo:
			raster.QuadTo(
				float32(seg.Args[0].X)/64.0+translateX, float32(seg.Args[0].Y)/64.0+translateY,
				float32(seg.Args[1].X)/64.0+translateX, float32(seg.Args[1].Y)/64.0+translateY,
			)
		case sfnt.SegmentOpCubeTo:
			raster.CubeTo(
				float32(seg.Args[0].X)/64.0+translateX, float32(seg.Args[0].Y)/64.0+translateY,
				float32(seg.Args[1].X)/64.0+translateX, float32(seg.Args[1].Y)/64.0+translateY,
				float32(seg.Args[2].X)/64.0+translateX, float32(seg.Args[2].Y)/64.0+translateY,
			)
		}
	}

	raster.ClosePath()

	alphaImg := image.NewAlpha(image.Rect(0, 0, width, height))
	raster.Draw(alphaImg, alphaImg.Bounds(), image.White, image.Point{})

	alphaImg = applyGaussianBlur(alphaImg)

	rgbaImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			a := alphaImg.AlphaAt(x, y).A
			rgbaImg.SetRGBA(x, y, color.RGBA{255, 255, 255, a})
		}
	}

	return width, height, rgbaImg, nil
}

func applyGaussianBlur(img *image.Alpha) *image.Alpha {
	kernel := []float32{
		1 / 16.0, 2 / 16.0, 1 / 16.0,
		2 / 16.0, 4 / 16.0, 2 / 16.0,
		1 / 16.0, 2 / 16.0, 1 / 16.0,
	}
	kernelSize := 3

	blurred := image.NewAlpha(img.Bounds())

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			var sum float32
			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					nx := x + kx - 1
					ny := y + ky - 1
					if nx >= 0 && nx < img.Bounds().Dx() && ny >= 0 && ny < img.Bounds().Dy() {
						sum += float32(img.AlphaAt(nx, ny).A) * kernel[ky*kernelSize+kx]
					}
				}
			}
			blurred.SetAlpha(x, y, color.Alpha{A: uint8(math.Round(float64(sum)))})
		}
	}

	return blurred
}

func saveAtlasToFile(img *image.RGBA, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
