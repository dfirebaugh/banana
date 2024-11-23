package opengl

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"unsafe"

	"github.com/dfirebaugh/banana/assets"
	"github.com/dfirebaugh/banana/graphics"
	"github.com/dfirebaugh/banana/graphics/font"
	"github.com/dfirebaugh/banana/graphics/opengl/shaders"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/math/fixed"
)

type AttribLocation uint32

const (
	ATTRIB_POS_LOCATION           AttribLocation = 0
	ATTRIB_SHAPE_POS_LOCATION     AttribLocation = 1
	ATTRIB_LOCAL_POS_LOCATION     AttribLocation = 2
	ATTRIB_OPCODE_LOCATION        AttribLocation = 3
	ATTRIB_RADIUS_LOCATION        AttribLocation = 4
	ATTRIB_COLOR_LOCATION         AttribLocation = 5
	ATTRIB_WIDTH_LOCATION         AttribLocation = 6
	ATTRIB_HEIGHT_LOCATION        AttribLocation = 7
	ATTRIB_TEX_COORD_LOCATION     AttribLocation = 8
	ATTRIB_RESOLUTION_LOCATION    AttribLocation = 9
	ATTRIB_TEXTURE_INDEX_LOCATION AttribLocation = 10
	ATTRIB_FONT_INDEX_LOCATION    AttribLocation = 11
)

const (
	AtlasWidth  = 512
	AtlasHeight = 512
	MaxTextures = 100
	MaxVertices = 600000
)

type Renderer struct {
	Vertices       []graphics.Vertex
	Framebuffers   []*Framebuffer
	VertexCount    int
	Textures       []TextureAtlas
	TextureCount   int
	VAO            uint32
	VBO            uint32
	BufferCapacity int
	ShaderProgram  uint32
	Font           *font.Font
	FontTextureID  uint32
	*TextureManager
}

func NewRenderer() *Renderer {
	initialCapacity := 1024
	renderer := &Renderer{
		Vertices:       make([]graphics.Vertex, initialCapacity),
		Framebuffers:   make([]*Framebuffer, 0),
		VertexCount:    0,
		Textures:       make([]TextureAtlas, MaxTextures),
		Font:           &font.Font{},
		BufferCapacity: initialCapacity,
	}
	renderer.TextureManager = NewTextureManager(renderer)
	return renderer
}

func (renderer *Renderer) Init() {
	var err error
	renderer.ShaderProgram, err = newShaderProgram(shaders.VertexShaderSource, shaders.FragmentShaderSource)
	if err != nil {
		fmt.Printf("Shader compilation or linking error: %s\n", err)
		return
	}

	renderer.Font, err = font.LoadFont(assets.LatoRegular)
	if err != nil {
		fmt.Printf("Failed to load font: %s\n", err)
		return
	}

	fontImg := renderer.Font.Image()
	fontImg = flipImageVertically(fontImg)
	renderer.FontTextureID = renderer.TextureManager.UploadTexture(fontImg)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.GenVertexArrays(1, &renderer.VAO)
	gl.GenBuffers(1, &renderer.VBO)

	gl.BindVertexArray(renderer.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, renderer.BufferCapacity*int(unsafe.Sizeof(graphics.Vertex{})), nil, gl.DYNAMIC_DRAW)

	stride := int32(unsafe.Sizeof(graphics.Vertex{}))

	gl.EnableVertexAttribArray(uint32(ATTRIB_POS_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_POS_LOCATION), 2, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.FsQuadPos))

	gl.EnableVertexAttribArray(uint32(ATTRIB_SHAPE_POS_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_SHAPE_POS_LOCATION), 2, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.ShapePos))
	gl.EnableVertexAttribArray(uint32(ATTRIB_LOCAL_POS_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_LOCAL_POS_LOCATION), 2, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.LocalPos))

	gl.EnableVertexAttribArray(uint32(ATTRIB_OPCODE_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_OPCODE_LOCATION), 1, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.OpCode))

	gl.EnableVertexAttribArray(uint32(ATTRIB_RADIUS_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_RADIUS_LOCATION), 1, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.Radius))

	gl.EnableVertexAttribArray(uint32(ATTRIB_COLOR_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_COLOR_LOCATION), 4, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.Color))

	gl.EnableVertexAttribArray(uint32(ATTRIB_WIDTH_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_WIDTH_LOCATION), 1, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.Width))

	gl.EnableVertexAttribArray(uint32(ATTRIB_HEIGHT_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_HEIGHT_LOCATION), 1, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.Height))

	gl.EnableVertexAttribArray(uint32(ATTRIB_TEX_COORD_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_TEX_COORD_LOCATION), 2, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.TexCoord))

	gl.EnableVertexAttribArray(uint32(ATTRIB_RESOLUTION_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_RESOLUTION_LOCATION), 2, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.Resolution))

	gl.EnableVertexAttribArray(uint32(ATTRIB_TEXTURE_INDEX_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_TEXTURE_INDEX_LOCATION), 1, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.TextureIndex))

	gl.EnableVertexAttribArray(uint32(ATTRIB_FONT_INDEX_LOCATION))
	gl.VertexAttribPointerWithOffset(uint32(ATTRIB_FONT_INDEX_LOCATION), 1, gl.FLOAT, false, stride, unsafe.Offsetof(graphics.Vertex{}.FontIndex))

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (renderer *Renderer) ensureCapacityForVertices(additionalVertices int) error {
	requiredVertices := renderer.VertexCount + additionalVertices
	if requiredVertices <= len(renderer.Vertices) {
		return nil
	}

	newCapacity := len(renderer.Vertices)
	if newCapacity == 0 {
		newCapacity = 1024
	}
	for newCapacity < requiredVertices {
		newCapacity *= 2
	}

	newVertices := make([]graphics.Vertex, newCapacity)
	copy(newVertices, renderer.Vertices[:renderer.VertexCount])
	renderer.Vertices = newVertices

	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, newCapacity*int(unsafe.Sizeof(graphics.Vertex{})), nil, gl.DYNAMIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, renderer.VertexCount*int(unsafe.Sizeof(graphics.Vertex{})), unsafe.Pointer(&renderer.Vertices[0]))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	renderer.BufferCapacity = newCapacity

	return nil
}

func (renderer *Renderer) AddFramebuffer(width, height int) (graphics.Framebuffer, error) {
	fb, err := NewFramebuffer(width, height, renderer.TextureManager, renderer)
	if err != nil {
		return nil, err
	}
	renderer.Framebuffers = append(renderer.Framebuffers, fb)
	return fb, nil
}

func (renderer *Renderer) Destroy() {
	gl.DeleteVertexArrays(1, &renderer.VAO)
	gl.DeleteBuffers(1, &renderer.VBO)
	gl.DeleteProgram(renderer.ShaderProgram)

	for i := 0; i < renderer.TextureCount; i++ {
		gl.DeleteTextures(1, &renderer.Textures[i].ID)
	}

	renderer.Font.Destroy()
}

func (renderer *Renderer) Clear(c color.Color) {
	rgba := toRGBA(c)
	gl.ClearColor(rgba[0], rgba[1], rgba[2], rgba[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	renderer.VertexCount = 0
	renderer.Vertices = renderer.Vertices[:0]
}

func (renderer *Renderer) Begin() {
	renderer.VertexCount = 0
	renderer.Vertices = renderer.Vertices[:0]
}

func (renderer *Renderer) End() {
	// cleanup
}

func (renderer *Renderer) Draw() {
	gl.UseProgram(renderer.ShaderProgram)
	width, height := renderer.GetViewportSize()
	gl.Uniform2f(gl.GetUniformLocation(renderer.ShaderProgram, gl.Str("u_resolution\x00")), float32(width), float32(height))

	gl.BindVertexArray(renderer.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.VBO)

	size := renderer.VertexCount * int(unsafe.Sizeof(graphics.Vertex{}))
	if size > 0 {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, size, unsafe.Pointer(&renderer.Vertices[0]))
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, renderer.FontTextureID)
	gl.Uniform1i(gl.GetUniformLocation(renderer.ShaderProgram, gl.Str("samplers[0]\x00")), 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, renderer.TextureManager.atlas.ID)
	gl.Uniform1i(gl.GetUniformLocation(renderer.ShaderProgram, gl.Str("samplers[1]\x00")), 1)

	for _, fb := range renderer.Framebuffers {
		textureUnit := fb.GetTextureID()
		gl.ActiveTexture(gl.TEXTURE0 + textureUnit)
		gl.BindTexture(gl.TEXTURE_2D, fb.TextureID)
		samplerName := fmt.Sprintf("samplers[%d]\x00", textureUnit)
		gl.Uniform1i(gl.GetUniformLocation(renderer.ShaderProgram, gl.Str(samplerName)), int32(textureUnit))
	}

	if renderer.VertexCount > 0 {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(renderer.VertexCount))
	}

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (renderer *Renderer) Render(shape graphics.Renderable) {
	vertices := shape.GetVertices(renderer.GetViewportSize())
	if len(vertices) == 0 {
		return
	}

	additionalVertices := len(vertices)

	if err := renderer.ensureCapacityForVertices(additionalVertices); err != nil {
		logrus.Errorf("Failed to ensure capacity: %v", err)
		return
	}

	copy(renderer.Vertices[renderer.VertexCount:], vertices)
	renderer.VertexCount += additionalVertices
}

func bruteForceFixFloaters(r rune, ypos float32, ptSize float32) float32 {
	if r == '^' || r == '\'' {
		return ypos + ptSize
	}
	if r == '\'' || r == '"' || r == '`' {
		return ypos + ptSize/2
	}
	return ypos
}

func (renderer *Renderer) RenderText(text string, options *graphics.TextRenderOptions) {
	r, g, b, a := options.Color.RGBA()
	colorVec := [4]float32{
		float32(r) / 65535.0,
		float32(g) / 65535.0,
		float32(b) / 65535.0,
		float32(a) / 65535.0,
	}

	const dpi = 96.0
	scale := options.Size / 72.0 * dpi / 32.0

	cursorX := options.X
	cursorY := options.Y
	ppem := fixed.Int26_6(32 << 6)

	var prevRune rune

	width, height := renderer.GetViewportSize()
	for _, r := range text {
		if r == '\n' {
			continue
		}

		glyph, exists := renderer.Font.Glyphs[r]
		if !exists {
			log.Printf("Glyph for rune '%c' not found", r)
			continue
		}

		if prevRune != 0 {
			kern, err := renderer.Font.GetKerning(prevRune, r, ppem)
			if err == nil {
				kerning := float32(kern) / 64.0 * scale
				cursorX += kerning
			}
		}

		xpos := cursorX + glyph.BearingX*scale
		ypos := cursorY - (glyph.SizeHeight-glyph.BearingY)*scale
		ypos = bruteForceFixFloaters(glyph.Rune, ypos, options.Size)

		w := glyph.SizeWidth * scale
		h := glyph.SizeHeight * scale

		u0, v0, u1, v1 := glyph.TexCoords[0], glyph.TexCoords[1], glyph.TexCoords[2], glyph.TexCoords[3]
		v0, v1 = v1, v0

		normX0 := (xpos/float32(width))*2.0 - 1.0
		normY0 := 1.0 - (ypos/float32(height))*2.0
		normX1 := ((xpos+w)/float32(width))*2.0 - 1.0
		normY1 := 1.0 - ((ypos+h)/float32(height))*2.0

		vertices := []graphics.Vertex{
			// Triangle 1
			{
				FsQuadPos:  [2]float32{normX0, normY0},
				LocalPos:   [2]float32{0, 0},
				OpCode:     graphics.OP_CODE_TEXT,
				Color:      colorVec,
				Resolution: [2]float32{float32(width), float32(height)},
				TexCoord:   [2]float32{u0, v1},
			},
			{
				FsQuadPos:  [2]float32{normX0, normY1},
				LocalPos:   [2]float32{0, h},
				OpCode:     graphics.OP_CODE_TEXT,
				Color:      colorVec,
				Resolution: [2]float32{float32(width), float32(height)},
				TexCoord:   [2]float32{u0, v0},
			},
			{
				FsQuadPos:  [2]float32{normX1, normY1},
				LocalPos:   [2]float32{w, h},
				OpCode:     graphics.OP_CODE_TEXT,
				Color:      colorVec,
				Resolution: [2]float32{float32(width), float32(height)},
				TexCoord:   [2]float32{u1, v0},
			},
			// Triangle 2
			{
				FsQuadPos:  [2]float32{normX0, normY0},
				LocalPos:   [2]float32{0, 0},
				OpCode:     graphics.OP_CODE_TEXT,
				Color:      colorVec,
				Resolution: [2]float32{float32(width), float32(height)},
				TexCoord:   [2]float32{u0, v1},
			},
			{
				FsQuadPos:  [2]float32{normX1, normY1},
				LocalPos:   [2]float32{w, h},
				OpCode:     graphics.OP_CODE_TEXT,
				Color:      colorVec,
				Resolution: [2]float32{float32(width), float32(height)},
				TexCoord:   [2]float32{u1, v0},
			},
			{
				FsQuadPos:  [2]float32{normX1, normY0},
				LocalPos:   [2]float32{w, 0},
				OpCode:     graphics.OP_CODE_TEXT,
				Color:      colorVec,
				Resolution: [2]float32{float32(width), float32(height)},
				TexCoord:   [2]float32{u1, v1},
			},
		}

		if err := renderer.ensureCapacityForVertices(len(vertices)); err != nil {
			logrus.Errorf("Failed to ensure capacity: %v", err)
			break
		}

		baseIndex := renderer.VertexCount
		if baseIndex+len(vertices) > len(renderer.Vertices) {
			log.Println("Not enough space in Vertices slice to add text")
			return
		}

		copy(renderer.Vertices[baseIndex:], vertices)
		renderer.VertexCount += len(vertices)

		cursorX += glyph.AdvanceWidth * scale

		prevRune = r
	}
}

func (renderer *Renderer) RenderFramebuffer(fb graphics.Framebuffer, options *graphics.TextureRenderOptions) {
	options.TextureIndex = float32(fb.GetTextureID())
	renderer.renderTexture(options)
}

func (renderer *Renderer) RenderTexture(textureID uint32, options *graphics.TextureRenderOptions) {
	tm := renderer.TextureManager
	bounds, exists := tm.textureBounds[textureID]
	if !exists {
		logrus.Error("Texture handle not found")
		return
	}
	options.RectX = float32(bounds.Min.X) + options.RectX
	options.RectY = float32(bounds.Min.Y) + options.RectY
	options.Width = float32(tm.atlas.Width)
	options.Height = float32(tm.atlas.Height)
	renderer.renderTexture(options)
}

func (renderer *Renderer) renderTexture(options *graphics.TextureRenderOptions) {
	const verticesPerTexture = 6
	if renderer.VertexCount+verticesPerTexture > len(renderer.Vertices) {
		if err := renderer.ensureCapacityForVertices(verticesPerTexture); err != nil {
			logrus.Errorf("Failed to ensure capacity: %v", err)
			return
		}
	}

	screenWidth, screenHeight := renderer.GetViewportSize()
	normX, normY := normalizeCoordinates(options.X, options.Y, screenWidth, screenHeight)

	width := options.DesiredWidth
	if width == 0 {
		width = options.RectWidth * options.Scale
	}

	height := options.DesiredHeight
	if height == 0 {
		height = options.RectHeight * options.Scale
	}

	u0 := options.RectX / options.Width
	v0 := options.RectY / options.Height
	u1 := (options.RectX + options.RectWidth) / options.Width
	v1 := (options.RectY + options.RectHeight) / options.Height

	if options.FlipX {
		u0, u1 = u1, u0
	}

	if options.FlipY {
		v0, v1 = v1, v0
	}

	vertices := []graphics.Vertex{
		{
			FsQuadPos:  [2]float32{normX, normY},
			LocalPos:   [2]float32{0, 0},
			TexCoord:   [2]float32{u0, v1},
			OpCode:     graphics.OP_CODE_TEXTURE,
			Color:      [4]float32{1.0, 1.0, 1.0, 1.0},
			Resolution: [2]float32{float32(screenWidth), float32(screenHeight)},
		},
		{
			FsQuadPos:  [2]float32{normX, normY - height},
			LocalPos:   [2]float32{0, -height},
			TexCoord:   [2]float32{u0, v0},
			OpCode:     graphics.OP_CODE_TEXTURE,
			Color:      [4]float32{1.0, 1.0, 1.0, 1.0},
			Resolution: [2]float32{float32(screenWidth), float32(screenHeight)},
		},
		{
			FsQuadPos:  [2]float32{normX + width, normY - height},
			LocalPos:   [2]float32{width, -height},
			TexCoord:   [2]float32{u1, v0},
			OpCode:     graphics.OP_CODE_TEXTURE,
			Color:      [4]float32{1.0, 1.0, 1.0, 1.0},
			Resolution: [2]float32{float32(screenWidth), float32(screenHeight)},
		},
		{
			FsQuadPos:  [2]float32{normX, normY},
			LocalPos:   [2]float32{0, 0},
			TexCoord:   [2]float32{u0, v1},
			OpCode:     graphics.OP_CODE_TEXTURE,
			Color:      [4]float32{1.0, 1.0, 1.0, 1.0},
			Resolution: [2]float32{float32(screenWidth), float32(screenHeight)},
		},
		{
			FsQuadPos:  [2]float32{normX + width, normY - height},
			LocalPos:   [2]float32{width, -height},
			TexCoord:   [2]float32{u1, v0},
			OpCode:     graphics.OP_CODE_TEXTURE,
			Color:      [4]float32{1.0, 1.0, 1.0, 1.0},
			Resolution: [2]float32{float32(screenWidth), float32(screenHeight)},
		},
		{
			FsQuadPos:  [2]float32{normX + width, normY},
			LocalPos:   [2]float32{width, 0},
			TexCoord:   [2]float32{u1, v1},
			OpCode:     graphics.OP_CODE_TEXTURE,
			Color:      [4]float32{1.0, 1.0, 1.0, 1.0},
			Resolution: [2]float32{float32(screenWidth), float32(screenHeight)},
		},
	}

	cosTheta := float32(math.Cos(float64(options.Rotation)))
	sinTheta := float32(math.Sin(float64(options.Rotation)))

	for i := 0; i < verticesPerTexture; i++ {
		v := &vertices[i]
		v.TextureIndex = options.TextureIndex

		localX := v.LocalPos[0]
		localY := v.LocalPos[1]

		rotatedX := localX*cosTheta - localY*sinTheta
		rotatedY := localX*sinTheta + localY*cosTheta

		v.FsQuadPos[0] = normX + rotatedX/float32(screenWidth)*2.0
		v.FsQuadPos[1] = normY + rotatedY/float32(screenHeight)*2.0

		v.LocalPos[0] = rotatedX
		v.LocalPos[1] = rotatedY
	}

	copy(renderer.Vertices[renderer.VertexCount:], vertices)
	renderer.VertexCount += verticesPerTexture
}

func normalizeCoordinates(x, y float32, screenWidth, screenHeight int) (float32, float32) {
	normX := (x/float32(screenWidth))*2.0 - 1.0
	normY := 1.0 - (y/float32(screenHeight))*2.0
	return normX, normY
}

func (renderer *Renderer) GetViewportSize() (int, int) {
	var viewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &viewport[0])
	return int(viewport[2]), int(viewport[3])
}

func (renderer *Renderer) LoadFont(fontData []byte) (graphics.Font, error) {
	font, err := font.LoadFont(fontData)
	if err != nil {
		logrus.Error(err)
	}
	renderer.Font = font
	return font, err
}

func (renderer *Renderer) BindFramebuffer(fb graphics.Framebuffer) {
	if fb != nil {
		fb.Bind()
		gl.Viewport(0, 0, int32(fb.GetWidth()), int32(fb.GetHeight()))
	} else {
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		width, height := renderer.GetViewportSize()
		gl.Viewport(0, 0, int32(width), int32(height))
	}
}

func (renderer *Renderer) UnbindFramebuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	width, height := renderer.GetViewportSize()
	gl.Viewport(0, 0, int32(width), int32(height))
}

func toRGBA(c color.Color) [4]float32 {
	r, g, b, a := c.RGBA()
	return [4]float32{
		float32(r) / 0xFFFF,
		float32(g) / 0xFFFF,
		float32(b) / 0xFFFF,
		float32(a) / 0xFFFF,
	}
}
