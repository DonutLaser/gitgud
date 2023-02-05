package main

import (
	"github.com/DonutLaser/git-client/renderer"
	"github.com/veandco/go-sdl2/sdl"
)

type NoChanges struct {
	Rect *sdl.Rect
}

func NewNoChanges(windowWidth int32, windowHeight int32) (result NoChanges) {
	result.Rect = &sdl.Rect{X: 0, Y: 0, W: windowWidth, H: windowHeight - 24 - 1}

	return
}

func (nochanges *NoChanges) Resize(windowWidth int32, windowHeight int32) {
	nochanges.Rect.W = windowWidth
	nochanges.Rect.H = windowHeight - 24 - 1
}

func (nochanges *NoChanges) Render(rend *sdl.Renderer, app *App) {
	mainFont := app.Fonts["16"]

	text := "No changes to show"
	textWidth := mainFont.GetStringWidth(text)

	textRect := sdl.Rect{
		X: nochanges.Rect.X + (nochanges.Rect.W-textWidth)/2,
		Y: nochanges.Rect.Y + (nochanges.Rect.H-mainFont.Size)/2,
		W: textWidth,
		H: mainFont.Size,
	}
	renderer.DrawText(rend, &mainFont, text, &textRect, sdl.Color{R: 221, G: 221, B: 221, A: 255})
}
