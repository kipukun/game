package engine

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kipukun/game/engine/object"
)

type imgobj struct {
	img *ebiten.Image // underlying image
	object.Object
}

func (io *imgobj) Size() (int, int) {
	return io.img.Bounds().Max.X, io.img.Bounds().Max.Y
}

func (io *imgobj) Image() (*ebiten.Image, *ebiten.DrawImageOptions) {
	x, y := io.GetPosition()
	o := new(ebiten.DrawImageOptions)
	o.GeoM.Translate(x, y)
	return io.img, o
}

// FromText returns an ImageObject that is rendered using the engine's font and the supplied text.
func FromText(e *Engine, t string, c color.Color) object.ImageObject {
	io := new(imgobj)
	r := text.BoundString(e.Font(), t)
	io.img = ebiten.NewImage(r.Dx(), r.Dy())
	io.Object = object.NewEmpty(r.Dx(), r.Dy())
	text.Draw(io.img, t, e.Font(), 0, r.Dy(), c)
	return io
}

// FromAsset returns an object.ImageObject with the supplied asset, or nil and an error.
func FromAsset(e *Engine, p string) (object.ImageObject, error) {
	io := new(imgobj)
	f, err := e.Asset(p)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	io.img = ebiten.NewImageFromImage(img)
	io.Object = object.NewEmpty(img.Bounds().Dx(), img.Bounds().Dy())
	return io, nil

}

// FromImage returns an object.ImageObject from an image.Image.
func FromImage(img image.Image) object.ImageObject {
	io := new(imgobj)
	io.Object = object.NewEmpty(img.Bounds().Dx(), img.Bounds().Dy())
	io.img = ebiten.NewImageFromImage(img)
	return io
}

func FromEbitenImage(eimg *ebiten.Image) object.ImageObject {
	io := new(imgobj)
	io.img = eimg
	io.Object = object.NewEmpty(eimg.Bounds().Dx(), eimg.Bounds().Dy())
	return io
}

// Pinner pins an image to the screen.
type Pinner struct {
	io object.ImageObject
	c  *Camera
}

func (p *Pinner) Image() (*ebiten.Image, *ebiten.DrawImageOptions) {
	cx, cy := p.c.GetPosition()
	img, o := p.io.Image()
	o.GeoM.Translate(cx, cy)
	return img, o
}

// NewPinner returns a Pinner which pins io to the perspective of c.
func NewPinner(c *Camera, io object.ImageObject) *Pinner {
	p := new(Pinner)
	p.io = io
	p.c = c
	return p
}
