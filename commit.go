package main

import (
	"github.com/DonutLaser/git-client/renderer"
	"github.com/veandco/go-sdl2/sdl"
)

type Commit struct {
	BGRect *sdl.Rect
	Rect   *sdl.Rect
	Input  InputField

	Active         bool
	SubmitCallback func(string)
	Result         string
}

func NewCommit(windowWidth int32, windowHeight int32) (result Commit) {
	result.BGRect = &sdl.Rect{X: 0, Y: 0, W: windowWidth, H: windowHeight}
	result.Rect = &sdl.Rect{X: windowWidth/2 - 181, Y: 200, W: 543, H: 32}

	result.Input = NewInputField(&sdl.Rect{X: result.Rect.X + 2, Y: result.Rect.Y + 2, W: result.Rect.W - 4, H: 28})

	result.Active = false

	return
}

func (commit *Commit) Resize(windowWidth int32, windowHeight int32) {
	commit.BGRect.W = windowWidth
	commit.BGRect.H = windowHeight
	commit.Rect.X = windowWidth/2 - 181

	commit.Input.Resize(&sdl.Rect{X: commit.Rect.X + 2, Y: commit.Rect.Y + 2, W: commit.Rect.W - 4, H: 28})
}

func (commit *Commit) Tick(input *Input) {
	if input.Escape {
		commit.Active = false
		commit.Input.Clear()

		return
	}

	if input.TypedCharacter == '\n' {
		commit.Active = false
		commit.Input.Clear()

		if commit.Result != "" {
			commit.SubmitCallback(commit.Result)
		}

		return
	}

	commit.Input.Tick(input)

	if commit.Input.ValueChanged {
		commit.Result = commit.Input.Value.String()
	}
}

func (commit *Commit) Open(callback func(string)) {
	commit.Input.Placeholder = "Commit message"
	commit.SubmitCallback = callback
	commit.Result = ""

	commit.Active = true
}

func (commit *Commit) Render(rend *sdl.Renderer, app *App) {
	if !commit.Active {
		return
	}

	renderer.DrawRectTransparent(rend, commit.BGRect, sdl.Color{R: 0, G: 0, B: 0, A: 102})

	renderer.DrawRect(rend, commit.Rect, sdl.Color{R: 18, G: 17, B: 20, A: 255})

	commit.Input.Render(rend, app)
}
