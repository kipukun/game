package object

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

// Fader is able to fade an ImageObject from transparent to opaque.
type Fader struct {
	ImageObject
	t float64
}

func NewFader(img ImageObject) *Fader {
	f := new(Fader)
	f.ImageObject = img
	return f
}

// Calculate progresses the animation, and calls callback when the animation is finished.
func (f *Fader) Calculate(callback func()) {
	if f.t >= 1.0 {
		if callback != nil {
			callback()
		}
		return
	}
	f.t += 0.008
}

func (f *Fader) Image() (*ebiten.Image, *ebiten.DrawImageOptions) {
	img, o := f.ImageObject.Image()
	o.ColorM.Translate(0, 0, 0, -1+f.t)
	return img, o
}

// Ease takes in an Object and returns a function that when called, eases o by (tx, ty).
func Easer[O Object](o O) func(x, y float64) {
	t := 0.0
	startx, starty := o.GetPosition()
	return func(tx, ty float64) {
		if t >= 1 {
			return
		}
		t += 0.01
		frac := t * t * (3.0 - 2.0*t)
		o.SetPosition(startx+tx*frac, starty+ty*frac)
	}
}

func EaserTo(o Object) func(t Object) func() {
	e := Easer(o)
	ox, oy := o.GetPosition()
	return func(t Object) func() {
		return func() {
			tx, ty := t.GetPosition()
			e(tx-ox, ty-oy)
		}
	}
}
