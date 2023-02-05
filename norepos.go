package main

import (
	"github.com/DonutLaser/git-client/renderer"
	"github.com/veandco/go-sdl2/sdl"
)

type NoRepos struct {
	Rect *sdl.Rect
}

func NewNoRepos(windowWidth int32, windowHeight int32) (result NoRepos) {
	result.Rect = &sdl.Rect{X: 0, Y: 0, W: windowWidth, H: windowHeight}

	return
}

func (norepos *NoRepos) Resize(windowWidth int32, windowHeight int32) {
	norepos.Rect.W = windowWidth
	norepos.Rect.H = windowHeight
}

func (norepos *NoRepos) Render(rend *sdl.Renderer, app *App) {
	mainFont := app.Fonts["16"]

	text := "No repositories added to the client. Press `Ctrl + Shift + O` to add a repository"
	textWidth := mainFont.GetStringWidth(text)

	textRect := sdl.Rect{
		X: norepos.Rect.X + (norepos.Rect.W-textWidth)/2,
		Y: norepos.Rect.Y + (norepos.Rect.H-mainFont.Size)/2,
		W: textWidth,
		H: mainFont.Size,
	}
	renderer.DrawText(rend, &mainFont, text, &textRect, sdl.Color{R: 221, G: 221, B: 221, A: 255})
}
