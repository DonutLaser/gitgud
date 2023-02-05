package font

import (
	"github.com/veandco/go-sdl2/ttf"
)

type Font struct {
	Data           *ttf.Font
	Size           int32
	CharacterWidth int32
}

func LoadFont(path string, size int32) (result Font) {
	font, err := ttf.OpenFont(path, int(size))
	if err != nil {
		panic(err.Error())
	}

	// We assume that the font is going to always be monospaced
	metrics, err := font.GlyphMetrics('m')
	if err != nil {
		panic(err.Error())
	}

	result.Data = font
	result.Size = size
	result.CharacterWidth = int32(metrics.Advance)

	return
}

func (font *Font) GetStringWidth(text string) int32 {
	return int32(int32(len(text)) * font.CharacterWidth)
}

func (font *Font) Unload() {
	font.Data.Close()
}
