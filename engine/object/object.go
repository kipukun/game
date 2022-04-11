package object

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kipukun/game/engine"
)

// An ImageObject is an Object that also has an associated ebiten.Image (i.e., can be drawn to the screen).
type ImageObject interface {
	Object

	Image() (*ebiten.Image, *ebiten.DrawImageOptions)
}

// Object represents some in-game object, with some
// bounding box and an ability to update its position.
type Object interface {
	Update()
	Size() (width, height int)
	Pos() (x, y float64)

	GetPosition() (x, y float64)
	GetVelocity() (dx, dy float64)
	GetAcceleration() (ddx, ddy float64)

	SetPosition(x, y float64)
	SetVelocity(dx, dy float64)
	SetAcceleration(ddx, ddy float64)
}

// FromText returns an ImageObject that is rendered using the engine's font and the supplied text.
func FromText(e *engine.Engine, t string, c color.Color) ImageObject {
	i := new(imgobj)
	r := text.BoundString(e.Font(), t)
	i.img = ebiten.NewImage(r.Dx(), r.Dy())
	text.Draw(i.img, t, e.Font(), 0, r.Dy(), c)
	return i
}

func FromAsset(e *engine.Engine, p string) (ImageObject, error) {
	i := new(imgobj)
	f, err := e.Asset(p)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	i.img = ebiten.NewImageFromImage(img)
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

func (i *imgobj) GetPosition() (x, y float64) {
	return i.x, i.y
}

func (i *imgobj) GetVelocity() (dx, dy float64) {
	return i.dx, i.dy
}
func (i *imgobj) GetAcceleration() (ddx, ddy float64) {
	return i.dx, i.dy
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
