package object

import (
	"sync"

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

// Easer implements an ease transition, moving an Object to a position using a parametric Bezier curve.
type Easer[O Object] struct {
	O     O
	t     float64
	end   float64
	start float64
	once  sync.Once
}

func NewEaser[O Object](o O, end float64) *Easer[O] {
	e := new(Easer[O])
	e.O = o
	e.end = end
	return e
}

func (e *Easer[O]) Calculate(callback func()) {
	e.once.Do(func() {
		_, y := e.O.GetPosition()
		e.start = y
	})
	if e.t >= 1.0 {
		if callback != nil {
			callback()
		}
		return
	}
	x, _ := e.O.GetPosition()
	e.t += 0.01
	frac := e.t * e.t * (3.0 - 2.0*e.t)
	e.O.SetPosition(x, e.start+e.end*frac)
}
