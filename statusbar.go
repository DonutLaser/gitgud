package main

import (
	"github.com/DonutLaser/git-client/renderer"
	"github.com/veandco/go-sdl2/sdl"
)

type Statusbar struct {
	Rect *sdl.Rect

	RepoName    string
	BranchName  string
	StashExists bool

	StashExistsText string
}

func NewStatusbar(windowWidth int32, windowHeight int32) (result Statusbar) {
	result.Rect = &sdl.Rect{X: 0, Y: 0, W: windowWidth, H: 24}

	result.StashExists = false

	result.StashExistsText = "Stashed changes exist"

	return
}

func (statusbar *Statusbar) Resize(windowWidth int32, windowHeight int32) {
	statusbar.Rect.W = windowWidth
}

func (statusbar *Statusbar) ShowRepoName(name string) {
	statusbar.RepoName = name
}

func (statusbar *Statusbar) ShowBranchName(name string) {
	statusbar.BranchName = name
}

func (statusbar *Statusbar) ShowStashExists(exists bool) {
	statusbar.StashExists = exists
}

func (statusbar *Statusbar) Render(rend *sdl.Renderer, app *App) {
	renderer.DrawRect(rend, statusbar.Rect, sdl.Color{R: 47, G: 46, B: 47, A: 255})

	mainFont := app.Fonts["12"]
	repoIcon := app.Icons["repo"]
	branchIcon := app.Icons["branch"]

	repoNameWidth := mainFont.GetStringWidth(statusbar.RepoName)
	branchNameWidth := mainFont.GetStringWidth(statusbar.BranchName)

	if statusbar.StashExists {
		stashIcon := app.Icons["stash"]

		renderer.DrawImage(rend, &stashIcon, &sdl.Point{X: 5, Y: statusbar.Rect.Y + (statusbar.Rect.H-stashIcon.Height)/2}, sdl.Color{R: 171, G: 171, B: 171, A: 255})

		stashTextWidth := mainFont.GetStringWidth(statusbar.StashExistsText)

		stashRect := sdl.Rect{
			X: stashIcon.Width + 10,
			Y: statusbar.Rect.Y + (statusbar.Rect.H-mainFont.Size)/2 + 1,
			W: stashTextWidth,
			H: mainFont.Size,
		}
		renderer.DrawText(rend, &mainFont, statusbar.StashExistsText, &stashRect, sdl.Color{R: 171, G: 171, B: 171, A: 255})
	}

	totalWidth := (repoIcon.Width + 5 + repoNameWidth) + 20 + (branchIcon.Width + 5 + branchNameWidth)
	left := statusbar.Rect.X + (statusbar.Rect.W-totalWidth)/2

	{
		renderer.DrawImage(rend, &repoIcon, &sdl.Point{X: left, Y: statusbar.Rect.Y + (statusbar.Rect.H-repoIcon.Height)/2 + 2}, sdl.Color{R: 171, G: 171, B: 171, A: 255})
		left += repoIcon.Width + 5
	}

	{
		repoNameRect := sdl.Rect{
			X: left,
			Y: statusbar.Rect.Y + (statusbar.Rect.H-mainFont.Size)/2 + 1,
			W: repoNameWidth,
			H: mainFont.Size,
		}
		renderer.DrawText(rend, &mainFont, statusbar.RepoName, &repoNameRect, sdl.Color{R: 171, G: 171, B: 171, A: 255})

		left += repoNameRect.W + 20
	}

	{
		renderer.DrawImage(rend, &branchIcon, &sdl.Point{X: left, Y: statusbar.Rect.Y + (statusbar.Rect.H-branchIcon.Height)/2 + 2}, sdl.Color{R: 171, G: 171, B: 171, A: 255})
		left += branchIcon.Width + 5
	}

	{
		branchNameRect := sdl.Rect{
			X: left,
			Y: statusbar.Rect.Y + (statusbar.Rect.H-mainFont.Size)/2 + 1,
			W: branchNameWidth,
			H: mainFont.Size,
		}
		renderer.DrawText(rend, &mainFont, statusbar.BranchName, &branchNameRect, sdl.Color{R: 171, G: 171, B: 171, A: 255})

		left += branchNameRect.W + 20
	}

}
