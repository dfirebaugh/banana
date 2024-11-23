package shaders

import (
	_ "embed"
)

//go:embed primitive.vert
var VertexShaderSource string

//go:embed primitive.frag
var FragmentShaderSource string
