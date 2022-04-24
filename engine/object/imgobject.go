package object

import (
	"image"
	"image/color"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// ImageObject is an Object that also has an associated ebiten.Image (i.e., can be drawn to the screen).
type ImageObject interface {
	Object

	Image() (*ebiten.Image, *ebiten.DrawImageOptions)

	SetOptions(opt *ebiten.DrawImageOptions)
}

type imgobj struct {
	opt *ebiten.DrawImageOptions
	img *ebiten.Image // underlying image
	Object
}

func (io *imgobj) Size() (int, int) {
	return io.img.Bounds().Max.X, io.img.Bounds().Max.Y
}

func (io *imgobj) Image() (*ebiten.Image, *ebiten.DrawImageOptions) {
	x, y := io.GetPosition()
	og := &ebiten.DrawImageOptions{}
	og.GeoM.Translate(x, y)
	io.opt.GeoM = og.GeoM
	return io.img, io.opt
}

func (io *imgobj) SetOptions(opt *ebiten.DrawImageOptions) {
	io.opt = opt
}

// FromText returns an ImageObject that is rendered using the engine's font and the supplied text.
func FromText(ft font.Face, t string, c color.Color) ImageObject {
	io := new(imgobj)
	io.opt = &ebiten.DrawImageOptions{}
	r := text.BoundString(ft, t)
	io.img = ebiten.NewImage(r.Dx(), r.Dy())
	io.Object = NewEmpty(r.Dx(), r.Dy())
	text.Draw(io.img, t, ft, 0, r.Dy(), c)
	return io
}

// FromAsset returns an object.ImageObject with the supplied asset, or nil and an error.
func FromReader(r io.Reader, p string) (ImageObject, error) {
	io := new(imgobj)
	io.opt = &ebiten.DrawImageOptions{}
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	io.img = ebiten.NewImageFromImage(img)
	io.Object = NewEmpty(img.Bounds().Dx(), img.Bounds().Dy())
	return io, nil

}

// FromImage returns an object.ImageObject from an image.Image.
func FromImage(img image.Image) ImageObject {
	io := new(imgobj)
	io.opt = &ebiten.DrawImageOptions{}
	io.Object = NewEmpty(img.Bounds().Dx(), img.Bounds().Dy())
	io.img = ebiten.NewImageFromImage(img)
	return io
}

func FromEbitenImage(eimg *ebiten.Image) ImageObject {
	io := new(imgobj)
	io.img = eimg
	io.Object = NewEmpty(eimg.Bounds().Dx(), eimg.Bounds().Dy())
	io.opt = &ebiten.DrawImageOptions{}
	return io
}

// Pinner pins an image to the screen.
type Pinner struct {
	io ImageObject
}

func (p *Pinner) Image(x, y float64) (*ebiten.Image, *ebiten.DrawImageOptions) {
	img, o := p.io.Image()
	o.GeoM.Translate(x, y)
	return img, o
}

// NewPinner returns a Pinner which pins io to the perspective of c.
func NewPinner(io ImageObject) *Pinner {
	p := new(Pinner)
	p.io = io
	return p
}
