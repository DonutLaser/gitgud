package main

import (
	"strings"

	"github.com/DonutLaser/git-client/renderer"
	"github.com/veandco/go-sdl2/sdl"
)

type InputField struct {
	Rect *sdl.Rect

	Placeholder  string
	Value        strings.Builder
	ValueChanged bool
}

func NewInputField(rect *sdl.Rect) (result InputField) {
	result.Rect = rect

	return
}

func (field *InputField) Resize(rect *sdl.Rect) {
	field.Rect = rect
}

func (field *InputField) Tick(input *Input) {
	field.ValueChanged = false

	if input.TypedCharacter == 0 {
		if input.Backspace && field.Value.Len() > 0 {
			if input.Ctrl {
				field.Value.Reset()
			} else {
				currentValue := field.Value.String()
				field.Value.Reset()
				field.Value.WriteString(currentValue[:len(currentValue)-1])
			}

			field.ValueChanged = true
		}
	} else {
		if input.TypedCharacter != '\t' {
			field.Value.WriteByte(input.TypedCharacter)
			field.ValueChanged = true
		}
	}
}

func (field *InputField) Clear() {
	field.Value.Reset()
}

func (field *InputField) Render(rend *sdl.Renderer, app *App) {
	renderer.DrawRect(rend, field.Rect, sdl.Color{R: 32, G: 33, B: 35, A: 255})

	mainFont := app.Fonts["14"]
	value := field.Value.String()
	color := sdl.Color{R: 221, G: 221, B: 221, A: 255}
	if value == "" && field.Placeholder != "" {
		value = field.Placeholder
		color = sdl.Color{R: 127, G: 127, B: 127, A: 255}
	}

	valueWidth := mainFont.GetStringWidth(value)
	valueRect := sdl.Rect{
		X: field.Rect.X + 5,
		Y: field.Rect.Y + (field.Rect.H-mainFont.Size)/2,
		W: valueWidth,
		H: mainFont.Size,
	}
	renderer.DrawText(rend, &mainFont, value, &valueRect, color)

	var cursorLeft int32 = 0
	if field.Value.String() != "" {
		cursorLeft = valueWidth
	}

	cursorRect := sdl.Rect{
		X: field.Rect.X + 5 + cursorLeft - 1,
		Y: field.Rect.Y + 5,
		W: 1,
		H: field.Rect.H - 10,
	}
	renderer.DrawRect(rend, &cursorRect, sdl.Color{R: 221, G: 221, B: 221, A: 255})
}
