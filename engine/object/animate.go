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

type Fader struct {
	img ImageObject
	t   float64
}

func NewFader(img ImageObject) *Fader {
	f := new(Fader)
	f.img = img
	return f
}

func (f *Fader) Calculate(callback func()) {
	f.t += 0.008
	if f.t >= 1.0 {
		return
	}
}

func (f *Fader) Image() (*ebiten.Image, *ebiten.DrawImageOptions) {
	img, o := f.img.Image()
	o.ColorM.Translate(0, 0, 0, -1+f.t)
	return img, o
}

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
		callback()
		return
	}
	x, _ := e.O.GetPosition()
	e.t += 0.01
	frac := e.t * e.t * (3.0 - 2.0*e.t)
	e.O.SetPosition(x, e.start+e.end*frac)
}
