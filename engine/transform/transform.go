package transform

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine/object"
)

type AnyFunc interface {
	FaderFunc | EaserFunc
}

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

type FaderFunc = func(dt float64)

func Fader(io object.ImageObject, tf TimeFunc, dur time.Duration) FaderFunc {
	elapsed := 0.0
	return func(dt float64) {
		no := &ebiten.DrawImageOptions{}
		t := clamp(elapsed/dur.Seconds(), 0, 1)
		frac := tf(t)
		elapsed += dt
		no.ColorM.Translate(0, 0, 0, lerp(-1, 0, frac))
		io.SetOptions(no)
	}

}

type EaserFunc = func(dt, x, y float64)

// Ease takes in an Object and returns a function that when called, eases o by (tx, ty).
func Easer[O object.Object](o O, tf TimeFunc, dur time.Duration) func(dt, x, y float64) {
	elapsed := 0.0
	startx, starty := o.GetPosition()
	return func(dt, tx, ty float64) {
		t := clamp(elapsed/dur.Seconds(), 0, 1)
		frac := tf(t)
		elapsed += dt
		o.SetPosition(lerp(startx, startx+tx, frac), lerp(starty, starty+ty, frac))
	}
}

type EaserToFunc = func(t object.Object) func(dt float64)

func EaserTo(o object.Object, tf TimeFunc, dur time.Duration) EaserToFunc {
	e := Easer(o, tf, dur)
	ox, oy := o.GetPosition()
	return func(t object.Object) func(dt float64) {
		return func(dt float64) {
			tx, ty := t.GetPosition()
			e(dt, tx-ox, ty-oy)
		}
	}
}
