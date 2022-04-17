package engine

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

func (e *Engine) NewEbitenFromAsset(p string) *ebiten.Image {
	r, err := e.Asset(p)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(r)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}
