package translation

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine/object"
)

type TimeFunc func(t float64) float64

var (
	Linear TimeFunc = func(t float64) float64 {
		return t
	}
	EaseInSine TimeFunc = func(t float64) float64 {
		return 1 - math.Cos((t*math.Pi)/2)
	}
	EaseInOutCubic TimeFunc = func(t float64) float64 {
		if t < 0.5 {
			return 4 * t * t * t
		} else {
			return 1 - math.Pow(-2*t+2, 3)/2
		}
	}
)

// Fader is able to fade an ImageObject from transparent to opaque.
type Fader struct {
	object.ImageObject
	t float64
}

func NewFader(img object.ImageObject) *Fader {
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
func Easer[O object.Object](o O, tf TimeFunc, dur time.Duration) func(dt, x, y float64) {
	elapsed := 0.0
	startx, starty := o.GetPosition()
	return func(dt, tx, ty float64) {
		if elapsed >= dur.Seconds() {
			elapsed = dur.Seconds()
		}
		elapsed += dt
		t := tf(elapsed / dur.Seconds())
		o.SetPosition(lerp(startx, startx+tx, t), lerp(starty, starty+ty, t))
		fmt.Printf("t: %f\r", t)
	}
}

func EaserTo(o object.Object, tf TimeFunc, dur time.Duration) func(t object.Object) func(dt float64) {
	e := Easer(o, tf, dur)
	ox, oy := o.GetPosition()
	return func(t object.Object) func(dt float64) {
		return func(dt float64) {
			tx, ty := t.GetPosition()
			e(dt, tx-ox, ty-oy)
		}
	}
}
