# banana engine
banana engine is a 2d graphics engine.

currently supported
- basic shape rendering
- texture rendering
- rendering to a framebuffer
- input detection (keyboard and mouse)
- bitmap font

I'm currently researching strategies to implement a gui system.

Check out the `examples` dir...

### drawing a triangle

`go run ./examples/triangle`

```golang
import (
	"github.com/dfirebaugh/banana"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 240
	screenHeight = 160
)

func main() {
	banana.SetWindowSize(screenWidth, screenHeight)
	banana.Run(nil, func() {
		banana.Clear(colornames.Skyblue)
		banana.RenderShape(&banana.Polygon{
			Vertices: []banana.Vertex{
				{
					X:     0,
					Y:     float32(screenHeight),
					Color: colornames.Red,
				},
				{
					X:     float32(screenWidth / 2),
					Y:     0,
					Color: colornames.Green,
				},
				{
					X:     float32(screenWidth),
					Y:     float32(screenHeight),
					Color: colornames.Blue,
				},
			},
		})
	})
}
```

![color triangle](./assets/images/color_triangle_example00.png)

