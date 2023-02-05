package image

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Image struct {
	Data   *sdl.Texture
	Width  int32
	Height int32
}

func LoadImage(path string, renderer *sdl.Renderer) (result Image) {
	image, err := img.Load(path)
	if err != nil {
		panic(err.Error())
	}

	texture, err := renderer.CreateTextureFromSurface(image)
	if err != nil {
		panic(err.Error())
	}

	result = Image{
		Data:   texture,
		Width:  image.W,
		Height: image.H,
	}

	image.Free()

	return
}

func (image *Image) Unload() {
	image.Data.Destroy()
}
