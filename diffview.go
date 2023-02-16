package main

import (
	"strconv"

	"github.com/DonutLaser/git-client/git"
	"github.com/DonutLaser/git-client/renderer"
	"github.com/veandco/go-sdl2/sdl"
)

type DiffView struct {
	OldRect *sdl.Rect
	NewRect *sdl.Rect

	Data  git.GitDiff
	Entry git.GitStatusEntry
}

func NewDiffView(windowWidth int32, windowHeight int32) (result DiffView) {
	width := (windowWidth - 280 - 2) / 2
	height := windowHeight - 24 - 2

	result.OldRect = &sdl.Rect{X: 280 + 2, Y: 24 + 2, W: width, H: height}
	result.NewRect = &sdl.Rect{X: result.OldRect.X + result.OldRect.W + 2, Y: 24 + 2, W: width, H: height}

	return
}

func (diff *DiffView) Resize(windowWidth int32, windowHeight int32) {
	width := (windowWidth - 280 - 2) / 2
	height := windowHeight - 24 - 2

	diff.OldRect.W = width
	diff.OldRect.H = height
	diff.NewRect.X = diff.OldRect.X + diff.OldRect.W + 2
	diff.NewRect.W = width
	diff.NewRect.H = height
}

func (diff *DiffView) ShowDiff(data git.GitDiff, entry git.GitStatusEntry) {
	diff.Data = data
	diff.Entry = entry
}

func (diff *DiffView) Render(rend *sdl.Renderer, app *App) {
	diff.renderOld(rend, app)
	diff.renderNew(rend, app)
}

func (diff *DiffView) renderOld(rend *sdl.Renderer, app *App) {
	renderer.DrawRect(rend, diff.OldRect, sdl.Color{R: 47, G: 46, B: 47, A: 255})

	if len(diff.Data.Chunks) == 1 && diff.Data.Chunks[0].Old.StartLine == 0 && diff.Data.Chunks[0].Old.EndLine == 0 {
		return
	}

	numbersRect := sdl.Rect{
		X: diff.OldRect.X,
		Y: diff.OldRect.Y,
		W: 40,
		H: diff.OldRect.H,
	}
	renderer.DrawRect(rend, &numbersRect, sdl.Color{R: 30, G: 30, B: 30, A: 255})

	mainFont := app.Fonts["12"]

	var lineHeight int32 = 23

	for _, chunk := range diff.Data.Chunks {
		lineTop := diff.OldRect.Y
		lineNumber := chunk.Old.StartLine

		for _, line := range chunk.Old.Lines {
			if line.Type != git.GIT_LINE_UNMODIFIED {
				bgRect := sdl.Rect{
					X: numbersRect.X + numbersRect.W,
					Y: lineTop,
					W: diff.OldRect.W - numbersRect.W,
					H: lineHeight,
				}

				bgColor := diff.diffLineTypeToColor(line.Type)

				renderer.DrawRectTransparent(rend, &bgRect, bgColor)

				lineNumberBgRect := sdl.Rect{
					X: numbersRect.X,
					Y: lineTop,
					W: numbersRect.W,
					H: lineHeight,
				}

				renderer.DrawRectTransparent(rend, &lineNumberBgRect, bgColor)
			}

			lineNumberStr := strconv.Itoa(int(lineNumber))

			lineNumberWidth := mainFont.GetStringWidth(lineNumberStr)
			lineNumberRect := sdl.Rect{
				X: numbersRect.X + numbersRect.W - lineNumberWidth - 10,
				Y: lineTop + (lineHeight-mainFont.Size)/2,
				W: lineNumberWidth,
				H: mainFont.Size,
			}
			renderer.DrawText(rend, &mainFont, lineNumberStr, &lineNumberRect, sdl.Color{R: 171, G: 171, B: 171, A: 255})

			textWidth := mainFont.GetStringWidth(line.Text)
			textRect := sdl.Rect{
				X: numbersRect.X + numbersRect.W + 10,
				Y: lineTop + (lineHeight-mainFont.Size)/2,
				W: textWidth,
				H: mainFont.Size,
			}
			renderer.DrawText(rend, &mainFont, line.Text, &textRect, sdl.Color{R: 171, G: 171, B: 171, A: 255})

			lineTop += lineHeight

			lineNumber += 1
		}
	}
}

func (diff *DiffView) renderNew(rend *sdl.Renderer, app *App) {
	renderer.DrawRect(rend, diff.NewRect, sdl.Color{R: 47, G: 46, B: 47, A: 255})

	if diff.Entry.Type == git.GIT_ENTRY_DELETED {
		renderer.DrawRectTransparent(rend, diff.NewRect, sdl.Color{R: 169, G: 26, B: 23, A: 49})

		text := "File was removed"
		font := app.Fonts["24"]

		textWidth := font.GetStringWidth(text)
		textRect := sdl.Rect{
			X: diff.NewRect.X + (diff.NewRect.W-textWidth)/2,
			Y: diff.NewRect.Y + (diff.NewRect.H-font.Size)/2,
			W: textWidth,
			H: font.Size,
		}

		renderer.DrawText(rend, &font, text, &textRect, sdl.Color{R: 171, G: 171, B: 171, A: 255})
		return
	}

	numbersRect := sdl.Rect{
		X: diff.NewRect.X,
		Y: diff.NewRect.Y,
		W: 40,
		H: diff.NewRect.H,
	}
	renderer.DrawRect(rend, &numbersRect, sdl.Color{R: 30, G: 30, B: 30, A: 255})

	mainFont := app.Fonts["12"]

	var lineHeight int32 = 23

	for _, chunk := range diff.Data.Chunks {
		lineTop := diff.NewRect.Y
		lineNumber := chunk.New.StartLine

		for _, line := range chunk.New.Lines {
			if line.Type != git.GIT_LINE_UNMODIFIED {
				bgRect := sdl.Rect{
					X: numbersRect.X + numbersRect.W,
					Y: lineTop,
					W: diff.OldRect.W - numbersRect.W,
					H: lineHeight,
				}

				bgColor := diff.diffLineTypeToColor(line.Type)

				renderer.DrawRectTransparent(rend, &bgRect, bgColor)

				lineNumberBgRect := sdl.Rect{
					X: numbersRect.X,
					Y: lineTop,
					W: numbersRect.W,
					H: lineHeight,
				}

				renderer.DrawRectTransparent(rend, &lineNumberBgRect, bgColor)
			}

			lineNumberStr := strconv.Itoa(int(lineNumber))

			lineNumberWidth := mainFont.GetStringWidth(lineNumberStr)
			lineNumberRect := sdl.Rect{
				X: numbersRect.X + numbersRect.W - lineNumberWidth - 10,
				Y: lineTop + (lineHeight-mainFont.Size)/2,
				W: lineNumberWidth,
				H: mainFont.Size,
			}
			renderer.DrawText(rend, &mainFont, lineNumberStr, &lineNumberRect, sdl.Color{R: 171, G: 171, B: 171, A: 255})

			textWidth := mainFont.GetStringWidth(line.Text)
			textRect := sdl.Rect{
				X: numbersRect.X + numbersRect.W + 10,
				Y: lineTop + (lineHeight-mainFont.Size)/2,
				W: textWidth,
				H: mainFont.Size,
			}
			renderer.DrawText(rend, &mainFont, line.Text, &textRect, sdl.Color{R: 171, G: 171, B: 171, A: 255})

			lineTop += lineHeight

			lineNumber += 1
		}
	}
}

func (diff *DiffView) diffLineTypeToColor(t git.GitDiffLineType) sdl.Color {
	switch t {
	case git.GIT_LINE_NEW:
		return sdl.Color{R: 82, G: 153, B: 19, A: 49}
	case git.GIT_LINE_REMOVED:
		return sdl.Color{R: 169, G: 26, B: 23, A: 49}
	case git.GIT_LINE_UNMODIFIED:
		panic("Unmodified line should not have any background")
	default:
		panic("Unreachable")
	}
}
