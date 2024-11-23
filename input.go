package banana

import "github.com/dfirebaugh/banana/pkg/input"

// IsKeyPressed checks if a key is currently pressed
func IsKeyPressed(keyCode input.Key) bool {
	return banana.inputState.IsKeyPressed(keyCode)
}

func IsKeyJustPressed(keyCode input.Key) bool {
	return banana.inputState.IsKeyJustPressed(keyCode)
}

func PressKey(keyCode input.Key) {
	banana.inputState.PressKey(keyCode)
}

func ReleaseKey(keyCode input.Key) {
	banana.inputState.ReleaseKey(keyCode)
}

func IsButtonPressed(buttonCode input.MouseButton) bool {
	return banana.inputState.IsButtonPressed(buttonCode)
}

// IsButtonJustPressed checks if a mouse button was just pressed
func IsButtonJustPressed(buttonCode input.MouseButton) bool {
	return banana.inputState.IsButtonJustPressed(buttonCode)
}

// PressButton simulates a mouse button press
func PressButton(buttonCode input.MouseButton) {
	banana.inputState.PressButton(buttonCode)
}

func ReleaseButton(buttonCode input.MouseButton) {
	banana.inputState.ReleaseButton(buttonCode)
}

func GetCursorPosition() (int, int) {
	return banana.inputState.GetCursorPosition()
}

func SetScrollCallback(cb func(x float64, y float64)) {
	banana.inputState.SetScrollCallback(cb)
}
