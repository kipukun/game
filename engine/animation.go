package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine/object"
)

// Animation describes objects that can animate themselves.
type Animation interface {
	Progress()
	Rewind()
	object.ImageObject
}

type animatedObject struct {
	i      int
	frames []object.ImageObject
}

// NewAnimationFromImage returns a new Animation that creates frames from sheet of size size.
func NewAnimationFromImage(sheet object.ImageObject, size image.Rectangle, nframes int) Animation {
	ao := new(animatedObject)
	img, _ := sheet.Image()
	frames := make([]object.ImageObject, 0)
	for i := 0; i < nframes; i++ {
		simg := img.SubImage(size.Add(image.Pt(size.Dx()*i, 0)))
		frames = append(frames, FromImage(simg))
	}
	ao.frames = frames
	return ao
}

// NewAnimationFromImages returns an Animation that can animate the frames in frames.
func NewAnimationFromImages(frames []object.ImageObject) Animation {
	ao := new(animatedObject)
	ao.frames = frames
	return ao
}

func (ao *animatedObject) Progress() {
	if ao.i+1 > len(ao.frames)-1 {
		ao.i = len(ao.frames) - 1
		return
	}
	ao.i += 1
}

func (ao *animatedObject) Rewind() {
	if ao.i-1 < 0 {
		ao.i = 0
		return
	}
	ao.i -= 1
}

func (ao *animatedObject) Image() (*ebiten.Image, *ebiten.DrawImageOptions) {
	return ao.frames[ao.i].Image()
}

func (ao *animatedObject) Update() {
	for _, o := range ao.frames {
		o.Update()
	}
}

func (ao *animatedObject) Size() (width int, height int) {
	return ao.frames[ao.i].Size()
}

func (ao *animatedObject) GetPosition() (x float64, y float64) {
	return ao.frames[ao.i].GetPosition()
}

func (ao *animatedObject) GetVelocity() (dx float64, dy float64) {
	return ao.frames[ao.i].GetVelocity()
}

func (ao *animatedObject) GetAcceleration() (ddx float64, ddy float64) {
	return ao.frames[ao.i].GetAcceleration()
}

func (ao *animatedObject) SetPosition(x float64, y float64) {
	for _, o := range ao.frames {
		o.SetPosition(x, y)
	}
}

func (ao *animatedObject) SetVelocity(dx float64, dy float64) {
	for _, o := range ao.frames {
		o.SetVelocity(dx, dy)
	}
}

func (ao *animatedObject) SetAcceleration(ddx float64, ddy float64) {
	for _, o := range ao.frames {
		o.SetAcceleration(ddx, ddy)
	}
}
