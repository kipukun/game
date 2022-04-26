package engine

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine/errors"
)

func (e *Engine) NewEbitenFromAsset(p string) (*ebiten.Image, error) {
	var op errors.Op = "NewEbitenFromAsset"

	r, err := e.Asset(p)
	if err != nil {
		return nil, errors.Error(op, "error retrieving asset", err)
	}
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, errors.Error(op, "error decoding image", err)
	}
	return ebiten.NewImageFromImage(img), nil
}
