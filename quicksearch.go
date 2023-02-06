package main

import (
	"strings"

	"github.com/DonutLaser/git-client/renderer"
	"github.com/veandco/go-sdl2/sdl"
)

type SearchMethod uint8

const (
	SEARCH_BEGINS_WITH SearchMethod = iota
	SEARCH_INCLUDES
)

type QuickSearch struct {
	BGRect    *sdl.Rect
	ModalRect *sdl.Rect

	Input InputField

	Active         bool
	ItemsToSearch  []string
	SearchResult   []string
	ActiveResult   int
	Method         SearchMethod
	SubmitCallback func(string)
}

func NewQuickSearch(windowWidth int32, windowHeight int32) (result QuickSearch) {
	result.BGRect = &sdl.Rect{X: 0, Y: 0, W: windowWidth, H: windowHeight}
	result.ModalRect = &sdl.Rect{X: windowWidth/2 - 181, Y: 200, W: 362, H: 323}

	result.Input = NewInputField(&sdl.Rect{X: result.ModalRect.X, Y: result.ModalRect.Y, W: result.ModalRect.W, H: 28})

	result.Active = false

	return
}

func (search *QuickSearch) Resize(windowWidth int32, windowHeight int32) {
	search.BGRect.W = windowWidth
	search.BGRect.H = windowHeight
	search.ModalRect = &sdl.Rect{X: windowWidth/2 - 181, Y: 200, W: 362, H: 323}

	search.Input.Resize(&sdl.Rect{X: search.ModalRect.X, Y: search.ModalRect.Y, W: search.ModalRect.W, H: 28})
}

func (search *QuickSearch) Tick(input *Input) {
	if input.Escape {
		search.Active = false
		search.Input.Clear()

		return
	}

	if input.TypedCharacter == '\n' {
		search.Active = false
		search.Input.Clear()

		if search.ActiveResult >= 0 {
			search.SubmitCallback(search.SearchResult[search.ActiveResult])
		}

		return
	}

	if input.Alt {
		if input.TypedCharacter == 'j' {
			search.ActiveResult += 1
			if search.ActiveResult == len(search.SearchResult) {
				search.ActiveResult = len(search.SearchResult) - 1
			}
		}

		if input.TypedCharacter == 'k' {
			search.ActiveResult -= 1
			if search.ActiveResult < 0 {
				search.ActiveResult = 0
			}
		}

		return
	} else {
		if search.ActiveResult > 0 {
			search.Active = false

			search.Input.Clear()
			search.SubmitCallback(search.SearchResult[search.ActiveResult])

			return
		}
	}

	search.Input.Tick(input)

	if search.Input.ValueChanged {
		search.SearchResult = make([]string, 0)
		query := search.Input.Value.String()

		if query != "" {
			if search.Method == SEARCH_BEGINS_WITH {
				for _, item := range search.ItemsToSearch {
					if strings.HasPrefix(strings.ToLower(item), strings.ToLower(query)) {
						search.SearchResult = append(search.SearchResult, item)
					}
				}
			} else if search.Method == SEARCH_INCLUDES {
				for _, item := range search.ItemsToSearch {
					if strings.Contains(strings.ToLower(item), strings.ToLower(query)) {
						search.SearchResult = append(search.SearchResult, item)
					}
				}
			} else {
				panic("Unreachable")
			}
		} else {
			search.SearchResult = search.ItemsToSearch
		}

		if len(search.SearchResult) > 0 {
			search.ActiveResult = 0
		} else {
			search.ActiveResult = -1
		}

		if len(search.SearchResult) > 5 {
			search.SearchResult = search.SearchResult[0:5]
		}
	}
}

func (search *QuickSearch) Open(inputPlaceholder string, itemsToSearch []string, searchMethod SearchMethod, callback func(string)) {
	search.Input.Placeholder = inputPlaceholder
	search.Active = true
	search.ItemsToSearch = itemsToSearch
	search.SearchResult = search.ItemsToSearch
	search.ActiveResult = 0
	search.Method = searchMethod
	search.SubmitCallback = callback
}

func (search *QuickSearch) Render(rend *sdl.Renderer, app *App) {
	if !search.Active {
		return
	}

	renderer.DrawRectTransparent(rend, search.BGRect, sdl.Color{R: 0, G: 0, B: 0, A: 102})

	borderRect := sdl.Rect{
		X: search.ModalRect.X - 2,
		Y: search.ModalRect.Y - 2,
		W: search.ModalRect.W + 4,
		H: search.ModalRect.H + 4,
	}
	renderer.DrawRect(rend, &borderRect, sdl.Color{R: 18, G: 17, B: 20, A: 255})

	search.Input.Render(rend, app)

	resultsRect := sdl.Rect{
		X: search.ModalRect.X,
		Y: search.ModalRect.Y + 28 + 2,
		W: search.ModalRect.W,
		H: search.ModalRect.H - 28 - 2,
	}
	renderer.DrawRect(rend, &resultsRect, sdl.Color{R: 47, G: 46, B: 47, A: 255})

	mainFont := app.Fonts["16"]

	itemTop := resultsRect.Y
	for index, item := range search.SearchResult {
		itemBGRect := sdl.Rect{
			X: resultsRect.X,
			Y: itemTop,
			W: resultsRect.W,
			H: 28 + 2,
		}

		itemWidth := mainFont.GetStringWidth(item)
		itemRect := sdl.Rect{
			X: itemBGRect.X + 10,
			Y: itemBGRect.Y + (itemBGRect.H-mainFont.Size)/2,
			W: itemWidth,
			H: mainFont.Size,
		}

		bgColor := sdl.Color{R: 63, G: 63, B: 63, A: 255}
		if index == search.ActiveResult {
			bgColor = sdl.Color{R: 77, G: 77, B: 77, A: 255}
		}

		renderer.DrawRect(rend, &itemBGRect, bgColor)
		renderer.DrawText(rend, &mainFont, item, &itemRect, sdl.Color{R: 171, G: 171, B: 171, A: 255})

		if index == search.ActiveResult {
			renderer.DrawRectOutline(rend, &itemBGRect, sdl.Color{R: 92, G: 91, B: 92, A: 255}, 1)
		}

		itemTop += 28 + 2
	}
}
