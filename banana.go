package banana

import (
	"fmt"
	"image/color"
	"runtime"
	"time"

	"github.com/dfirebaugh/banana/graphics"
	"github.com/dfirebaugh/banana/graphics/opengl"
	"github.com/dfirebaugh/banana/pkg/input"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Game interface {
	Update()
	Render()
}

var (
	windowWidth  = 240
	windowHeight = 160
)

type engine struct {
	graphicsBackend   graphics.GraphicsBackend
	inputState        *input.InputState
	windowTitle       string
	fpsCounter        *fpsCounter
	hasSetupCompleted bool
}

var banana = &engine{}

func setup() {
	runtime.LockOSThread()
	banana.inputState = input.NewInputState()
	var err error
	gb, err := opengl.NewGraphicsBackend(windowWidth, windowHeight)
	if err != nil {
		panic(err.Error())
	}
	banana.graphicsBackend = gb
	banana.hasSetupCompleted = true
}

func initWindow() {
	SetWindowSize(windowWidth, windowHeight)
	SetTitle("banana")
}

func ensureSetupCompletion() {
	if banana.hasSetupCompleted {
		return
	}
	setup()
	initWindow()
}

func run(updateFn func(), renderFn func()) {
	ensureSetupCompletion()
	defer close()
	banana.fpsCounter = newFPSCounter()

	// width, height := GetWindowSize()
	// gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	// gl.Viewport(0, 0, int32(width), int32(height))
	banana.graphicsBackend.SetInputCallback(func(eventChan chan input.Event) {
		evt := <-eventChan
		handleEvent(evt, banana.inputState)
	})
	targetFPS := 120.0
	targetFrameDuration := time.Second / time.Duration(targetFPS)

	var lastUpdateTime time.Time
	var accumulator time.Duration

	lastUpdateTime = time.Now()
	for banana.graphicsBackend.PollEvents() {
		currentTime := time.Now()
		deltaTime := currentTime.Sub(lastUpdateTime)
		lastUpdateTime = currentTime
		accumulator += deltaTime

		frameRendered := false
		for accumulator >= targetFrameDuration {
			if updateFn != nil {
				updateFn()
			}
			accumulator -= targetFrameDuration
			frameRendered = true

		}

		if frameRendered {
			if renderFn != nil {
				renderFn()
			}

			calculateFPS()
			banana.graphicsBackend.Draw()
			banana.graphicsBackend.SwapBuffers()
			banana.inputState.ResetJustPressed()
		}
	}
}

func RunGame(game Game) {
	run(game.Update, game.Render)
}

func RunApp(game Game) {
	run(game.Update, game.Render)
}

// Run is the main update function called to refresh the engine state.
func Run(updateFn func(), renderFn func()) {
	run(updateFn, renderFn)
}

func calculateFPS() {
	banana.fpsCounter.Frame()
	fps := banana.fpsCounter.GetFPS()
	title := banana.windowTitle
	if fps != 0 && fpsEnabled {
		title = fmt.Sprintf("%s -- %d\n", title, int(fps))
	}
	if banana.graphicsBackend.IsDisposed() {
		return
	}
	banana.graphicsBackend.SetWindowTitle(title)
}

func GetFPS() float64 {
	return banana.fpsCounter.GetFPS()
}

func Close() {
	close()
}

func close() {
	banana.graphicsBackend.Close()
}

func Draw() {
	ensureSetupCompletion()
	banana.graphicsBackend.Draw()
}

// Clear clears the screen with the specified color.
func Clear(c color.Color) {
	ensureSetupCompletion()
	banana.graphicsBackend.Clear(c)
}

// SetTitle sets the title of the window.
func SetTitle(title string) {
	ensureSetupCompletion()
	banana.windowTitle = title
}

// GetWindowSize retrieves the current window size.
func GetWindowSize() (int, int) {
	ensureSetupCompletion()
	return banana.graphicsBackend.GetWindowSize()
}

func GetWindowWidth() int {
	ensureSetupCompletion()
	w, _ := banana.graphicsBackend.GetWindowSize()
	return w
}

func GetWindowHeight() int {
	ensureSetupCompletion()
	_, h := banana.graphicsBackend.GetWindowSize()
	return h
}

func GetWindowPosition() (int, int) {
	ensureSetupCompletion()
	return banana.graphicsBackend.GetWindowPosition()
}

func GetViewportSize() (int, int) {
	ensureSetupCompletion()
	return banana.graphicsBackend.GetViewportSize()
}

func SetWindowPosition(x, y int) {
	ensureSetupCompletion()
	banana.graphicsBackend.SetWindowPosition(x, y)
}

func SetWindowSize(width, height int) {
	ensureSetupCompletion()

	banana.graphicsBackend.SetWindowSize(width, height)
	windowWidth = width
	windowHeight = height
}

func DisableWindowResize() {
	ensureSetupCompletion()

	banana.graphicsBackend.DisableWindowResize()

	SetWindowSize(windowWidth, windowHeight)
}

func SetBorderlessWindowed(v bool) {
	ensureSetupCompletion()
	banana.graphicsBackend.SetBorderlessWindowed(v)
}

func SetFullScreenBorderless(v bool) {
	ensureSetupCompletion()
	banana.graphicsBackend.SetFullScreenBorderless(v)
}

func SetResizeCallback(fn func(physicalWidth, physicalHeight uint32)) {
	ensureSetupCompletion()
	banana.graphicsBackend.SetResizedCallback(fn)
}

func BindFramebuffer(fb graphics.Framebuffer) {
	ensureSetupCompletion()
	banana.graphicsBackend.BindFramebuffer(fb)
}

func UnbindFramebuffer() {
	ensureSetupCompletion()
	banana.graphicsBackend.UnbindFramebuffer()
	windowWidth, windowHeight := GetWindowSize()
	Viewport(0, 0, int32(windowWidth), int32(windowHeight))
}

func Viewport(x int32, y int32, width int32, height int32) {
	ensureSetupCompletion()
	gl.Viewport(x, y, width, height)
}
