package renderer

import (
	"github.com/DonutLaser/git-client/font"
	"github.com/DonutLaser/git-client/image"
	"github.com/veandco/go-sdl2/sdl"
)

func DrawRect(renderer *sdl.Renderer, rect *sdl.Rect, color sdl.Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.FillRect(rect)
}

func DrawRectOutline(renderer *sdl.Renderer, rect *sdl.Rect, color sdl.Color, outlineWidth int32) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)

	top := sdl.Rect{X: rect.X, Y: rect.Y, W: rect.W, H: outlineWidth}
	right := sdl.Rect{X: rect.X + rect.W - outlineWidth, Y: rect.Y, W: outlineWidth, H: rect.H}
	bottom := sdl.Rect{X: rect.X, Y: rect.Y + rect.H - outlineWidth, W: rect.W, H: outlineWidth}
	left := sdl.Rect{X: rect.X, Y: rect.Y, W: outlineWidth, H: rect.H}

	renderer.FillRect(&top)
	renderer.FillRect(&right)
	renderer.FillRect(&bottom)
	renderer.FillRect(&left)
}

func DrawRectTransparent(renderer *sdl.Renderer, rect *sdl.Rect, color sdl.Color) {
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	DrawRect(renderer, rect, color)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_NONE)
}

func DrawText(renderer *sdl.Renderer, ffont *font.Font, text string, rect *sdl.Rect, color sdl.Color) {
	surface, _ := ffont.Data.RenderUTF8Blended(text, color)
	defer surface.Free()

	texture, _ := renderer.CreateTextureFromSurface(surface)
	defer texture.Destroy()

	renderer.Copy(texture, nil, rect)
}

func DrawImage(renderer *sdl.Renderer, img *image.Image, position *sdl.Point, color sdl.Color) {
	rect := sdl.Rect{
		X: position.X,
		Y: position.Y,
		W: img.Width,
		H: img.Height,
	}

	img.Data.SetColorMod(color.R, color.G, color.B)
	renderer.Copy(img.Data, nil, &rect)
}
