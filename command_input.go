package main

import (
	"github.com/DonutLaser/git-client/renderer"
	"github.com/veandco/go-sdl2/sdl"
)

type CommandInput struct {
	BGRect *sdl.Rect
	Rect   *sdl.Rect
	Input  InputField

	Active         bool
	SubmitCallback func(string)
	Result         string
}

func NewCommandInput(windowWidth int32, windowHeight int32) (result CommandInput) {
	result.BGRect = &sdl.Rect{X: 0, Y: 0, W: windowWidth, H: windowHeight}
	result.Rect = &sdl.Rect{X: windowWidth/2 - 181, Y: 200, W: 543, H: 32}

	result.Input = NewInputField(&sdl.Rect{X: result.Rect.X + 2, Y: result.Rect.Y + 2, W: result.Rect.W - 4, H: 28})

	result.Active = false

	return
}

func (ci *CommandInput) Resize(windowWidth int32, windowHeight int32) {
	ci.BGRect.W = windowWidth
	ci.BGRect.H = windowHeight
	ci.Rect.X = windowWidth/2 - 181

	ci.Input.Resize(&sdl.Rect{X: ci.Rect.X + 2, Y: ci.Rect.Y + 2, W: ci.Rect.W - 4, H: 28})
}

func (ci *CommandInput) Tick(input *Input) {
	if input.Escape {
		ci.Active = false
		ci.Input.Clear()

		return
	}

	if input.TypedCharacter == '\n' {
		ci.Active = false
		ci.Input.Clear()

		if ci.Result != "" {
			ci.SubmitCallback(ci.Result)
		}

		return
	}

	ci.Input.Tick(input)

	if ci.Input.ValueChanged {
		ci.Result = ci.Input.Value.String()
	}
}

func (ci *CommandInput) Open(placeholder string, callback func(string)) {
	ci.Input.Placeholder = placeholder
	ci.SubmitCallback = callback
	ci.Result = ""

	ci.Active = true
}

func (ci *CommandInput) Render(rend *sdl.Renderer, app *App) {
	if !ci.Active {
		return
	}

	renderer.DrawRectTransparent(rend, ci.BGRect, sdl.Color{R: 0, G: 0, B: 0, A: 102})

	renderer.DrawRect(rend, ci.Rect, sdl.Color{R: 18, G: 17, B: 20, A: 255})

	ci.Input.Render(rend, app)
}
