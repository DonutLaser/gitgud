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
	result.ModalRect = &sdl.Rect{X: windowWidth/2 - 200, Y: windowHeight/2 - 150, W: 400, H: 0}

	result.Input = NewInputField(&sdl.Rect{X: result.ModalRect.X + 5, Y: result.ModalRect.Y + 5, W: result.ModalRect.W - 10, H: 28})

	result.Active = false

	return
}

func (search *QuickSearch) Resize(windowWidth int32, windowHeight int32) {
	search.BGRect.W = windowWidth
	search.BGRect.H = windowHeight
	search.ModalRect = &sdl.Rect{X: windowWidth/2 - 200, Y: windowHeight/2 - 150, W: 400, H: 0}

	search.Input.Resize(&sdl.Rect{X: search.ModalRect.X + 5, Y: search.ModalRect.Y + 5, W: search.ModalRect.W - 10, H: 28})
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
	search.SearchResult = make([]string, 0)
	search.ActiveResult = -1
	search.Method = searchMethod
	search.SubmitCallback = callback
}

func (search *QuickSearch) Render(rend *sdl.Renderer, app *App) {
	if !search.Active {
		return
	}

	mainFont := app.Fonts["16"]
	itemHeight := mainFont.Size + 10

	renderer.DrawRectTransparent(rend, search.BGRect, sdl.Color{R: 0, G: 0, B: 0, A: 122})
	if len(search.SearchResult) == 0 {
		search.ModalRect.H = 28 + 10
	} else {
		search.ModalRect.H = 28 + 10 + int32(len(search.SearchResult))*itemHeight + int32(len(search.SearchResult)) + 10
	}
	renderer.DrawRect(rend, search.ModalRect, sdl.Color{R: 48, G: 51, B: 59, A: 255})
	search.Input.Render(rend, app)

	if len(search.SearchResult) > 0 {
		resultsRect := sdl.Rect{
			X: search.ModalRect.X + 5,
			Y: search.ModalRect.Y + 28 + 10,
			W: search.ModalRect.W - 10,
			H: itemHeight*int32(len(search.SearchResult)) + int32(len(search.SearchResult)-1) + 2,
		}
		renderer.DrawRect(rend, &resultsRect, sdl.Color{R: 32, G: 33, B: 35, A: 255})

		for index, item := range search.SearchResult {
			itemBGRect := sdl.Rect{
				X: resultsRect.X + 1,
				Y: resultsRect.Y + (mainFont.Size+10)*int32(index) + int32(index) + 1,
				W: resultsRect.W - 2,
				H: mainFont.Size + 10,
			}

			itemWidth := mainFont.GetStringWidth(item)
			itemRect := sdl.Rect{
				X: itemBGRect.X + 5,
				Y: itemBGRect.Y + (itemBGRect.H-mainFont.Size)/2,
				W: itemWidth,
				H: mainFont.Size,
			}

			renderer.DrawRect(rend, &itemBGRect, sdl.Color{R: 61, G: 64, B: 71, A: 255})
			renderer.DrawText(rend, &mainFont, item, &itemRect, sdl.Color{R: 221, G: 221, B: 221, A: 255})

			if index == search.ActiveResult {
				renderer.DrawRectOutline(rend, &itemBGRect, sdl.Color{R: 221, G: 221, B: 221, A: 255}, 1)
			}
		}
	}
}
