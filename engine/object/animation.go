package object

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Animation describes objects that can animate themselves.
type Animation interface {
	Progress()
	Rewind()
	ImageObject
}

type animatedObject struct {
	i      int
	frames []ImageObject
}

// NewAnimationFromImages returns an Animation that can animate the frames in frames.
func NewAnimationFromImages(frames []ImageObject) Animation {
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
