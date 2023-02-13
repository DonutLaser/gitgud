package main

import (
	"github.com/DonutLaser/git-client/git"
	"github.com/DonutLaser/git-client/renderer"
	"github.com/veandco/go-sdl2/sdl"
)

type Staging struct {
	Rect *sdl.Rect

	Entries     []git.GitStatusEntry
	ActiveEntry int
	AllSelected bool
}

func NewStaging(windowHeight int32) (result Staging) {
	result.Rect = &sdl.Rect{X: 0, Y: 24 + 2, W: 280, H: windowHeight - 24 - 2}

	result.ActiveEntry = 0

	return
}

func (staging *Staging) Resize(windowHeight int32) {
	staging.Rect.H = windowHeight - 24 - 2
}

func (staging *Staging) ShowEntries(entries []git.GitStatusEntry) {
	staging.Entries = entries
	staging.ActiveEntry = 0
	staging.AllSelected = true
}

func (staging *Staging) GoToNextEntry() {
	staging.ActiveEntry += 1
	if int(staging.ActiveEntry) == len(staging.Entries) {
		staging.ActiveEntry = len(staging.Entries) - 1
	}
}

func (staging *Staging) GoToPrevEntry() {
	staging.ActiveEntry -= 1
	if staging.ActiveEntry < 0 {
		staging.ActiveEntry = 0
	}
}

func (staging *Staging) ToggleEntrySelected() {
	staging.Entries[staging.ActiveEntry].Selected = !staging.Entries[staging.ActiveEntry].Selected
}

func (staging *Staging) ToggleAllEntriesSelected() {
	staging.AllSelected = !staging.AllSelected

	for index := range staging.Entries {
		staging.Entries[index].Selected = staging.AllSelected
	}
}

func (staging *Staging) DiscardActiveEntry() {
	staging.Entries = append(staging.Entries[0:staging.ActiveEntry], staging.Entries[staging.ActiveEntry+1:]...)

	if staging.ActiveEntry > 0 && staging.ActiveEntry >= len(staging.Entries) {
		staging.ActiveEntry = len(staging.Entries) - 1
	}
}

func (staging *Staging) GetActiveEntry() git.GitStatusEntry {
	return staging.Entries[staging.ActiveEntry]
}

func (staging *Staging) GetActiveEntryFileName() string {
	return staging.Entries[staging.ActiveEntry].Filename
}

func (staging *Staging) Render(rend *sdl.Renderer, app *App) {
	renderer.DrawRect(rend, staging.Rect, sdl.Color{R: 47, G: 46, B: 47, A: 255})

	mainFont := app.Fonts["14"]
	onIcon := app.Icons["entry_on"]
	offIcon := app.Icons["entry_off"]

	top := staging.Rect.Y

	var entryHeight int32 = 28

	for index, entry := range staging.Entries {
		bgRect := sdl.Rect{
			X: staging.Rect.X,
			Y: top,
			W: staging.Rect.W,
			H: entryHeight,
		}

		bgColor := sdl.Color{R: 63, G: 63, B: 63, A: 255}
		if index == staging.ActiveEntry {
			bgColor = sdl.Color{R: 77, G: 77, B: 77, A: 255}
		}

		renderer.DrawRect(rend, &bgRect, bgColor)
		if index == staging.ActiveEntry {
			renderer.DrawRectOutline(rend, &bgRect, sdl.Color{R: 92, G: 91, B: 92, A: 255}, 1)
		}

		nameWidth := mainFont.GetStringWidth(entry.Filename)
		nameRect := sdl.Rect{
			X: bgRect.X + 10,
			Y: bgRect.Y + (bgRect.H-mainFont.Size)/2,
			W: nameWidth,
			H: mainFont.Size,
		}

		nameColor := sdl.Color{R: 171, G: 171, B: 171, A: 255}
		if !entry.Selected {
			nameColor = sdl.Color{R: 93, G: 93, B: 93, A: 255}
		}

		renderer.DrawText(rend, &mainFont, entry.Filename, &nameRect, nameColor)

		icon := onIcon
		if !entry.Selected {
			icon = offIcon
		}

		iconColor := staging.changeTypeToColor(entry.Type)
		renderer.DrawImage(rend, &icon, &sdl.Point{X: bgRect.X + bgRect.W - 10 - icon.Width, Y: bgRect.Y + (bgRect.H-icon.Height)/2}, iconColor)

		top += entryHeight + 2
	}
}

func (staging *Staging) changeTypeToColor(t git.GitStatusEntryType) sdl.Color {
	switch t {
	case git.GIT_ENTRY_MODIFIED:
		return sdl.Color{R: 207, G: 173, B: 16, A: 255}
	case git.GIT_ENTRY_NEW_UNSTAGED:
		fallthrough
	case git.GIT_ENTRY_NEW:
		return sdl.Color{R: 82, G: 153, B: 19, A: 255}
	case git.GIT_ENTRY_DELETED:
		return sdl.Color{R: 169, G: 26, B: 23, A: 255}
	default:
		panic("Unreachable")
	}
}
