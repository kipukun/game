package transform

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine/object"
)

type TimeFunc func(t float64) float64

type ChangeFunc func(dt float64) float64

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

func Fader(io object.ImageObject, tf TimeFunc, dur time.Duration) ChangeFunc {
	elapsed := 0.0
	return func(dt float64) float64 {
		no := &ebiten.DrawImageOptions{}
		t := clamp(elapsed/dur.Seconds(), 0, 1)
		frac := tf(t)
		elapsed += dt
		no.ColorM.Translate(0, 0, 0, lerp(-1, 0, frac))
		io.SetOptions(no)
		return t
	}

}

// Ease takes in an Object and returns a function that when called, eases o by (tx, ty).
func Easer[O object.Object](o O, tx, ty float64, tf TimeFunc, dur time.Duration) ChangeFunc {
	elapsed := 0.0
	startx, starty := o.GetPosition()
	endx, endy := startx+tx, starty+ty
	return func(dt float64) float64 {
		t := clamp(elapsed/dur.Seconds(), 0, 1)
		frac := tf(t)
		elapsed += dt
		o.SetPosition(lerp(startx, endx, frac), lerp(starty, endy, frac))
		return t
	}
}
