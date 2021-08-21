package object

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/kipukun/game/engine"
)

// Object represents some in-game object, with some
// bounding box and an ability to update its position.
type Object interface {
	Update()
	Size() (width, height int)
	Pos() (x, y float64)

	SetPosition(x, y float64)
	SetVelocity(dx, dy float64)
	SetAcceleration(ddx, ddy float64)

	Image() (*ebiten.Image, *ebiten.DrawImageOptions)
}

func FromAsset(e *engine.Engine, p string) (Object, error) {
	i := new(imgobj)
	f, err := e.Asset(p)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	i.img, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	return i, nil

}

type imgobj struct {
	img      *ebiten.Image // underlying image
	x, y     float64
	dx, dy   float64
	ddx, ddy float64
}

func (i *imgobj) Update() {
	i.dx += i.ddx
	i.dy += i.ddy
	i.x += i.dx
	i.y += i.dy
}

func (i *imgobj) Size() (int, int) {
	return i.img.Bounds().Max.X, i.img.Bounds().Max.Y
}

func (i *imgobj) Pos() (float64, float64) {
	return i.x, i.y
}

func (i *imgobj) SetPosition(x, y float64) {
	i.x = x
	i.y = y
}

func (i *imgobj) SetVelocity(dx, dy float64) {
	i.dx = dx
	i.dy = dy
}

func (i *imgobj) SetAcceleration(ddx, ddy float64) {
	i.ddx = ddx
	i.ddy = ddy
}

func (i *imgobj) Image() (*ebiten.Image, *ebiten.DrawImageOptions) {
	o := new(ebiten.DrawImageOptions)
	o.GeoM.Translate(i.x, i.y)
	return i.img, o
}
