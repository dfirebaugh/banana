package assets

import _ "embed"

//go:embed models/the-utah-teapot/source/teapot.obj
var TeaPot string

//go:embed models/the-utah-teapot/source/default.png
var DefaultTextureImage []byte

//go:embed  images/buddy_dance.png
var BuddyDanceSpriteSheet []byte

//go:embed fonts/ComicShannsMono/ComicShannsMonoNerdFont-Regular.otf
var ComicShannsMonoNerdFontRegular []byte

//go:embed fonts/Noto/NotoSansMNerdFontMono-Regular.ttf
var NotoSansNerdFontMonoRegular []byte

//go:embed fonts/Lato/Lato-Regular.ttf
var LatoRegular []byte
