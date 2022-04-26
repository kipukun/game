package states

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/engine/object"
)

type OptionsState struct {
	volume               float64
	vlabel, instructions object.ImageObject
	vbar                 object.Animation
}

func progressBar(e *engine.Engine, len int) object.Animation {
	fs := make([]object.ImageObject, 0)
	for i := 0; i <= len; i++ {
		f := object.FromText(e.Font(), "["+strings.Repeat("|", i), color.White)
		fs = append(fs, f)
	}
	return object.NewAnimationFromImages(fs)
}

func (ops *OptionsState) Update(e *engine.Engine, dt float64) error {
	return nil
}

func (ops *OptionsState) Draw(e *engine.Engine, s *ebiten.Image) {
	s.DrawImage(ops.instructions.Image())
	s.DrawImage(ops.vlabel.Image())
	s.DrawImage(ops.vbar.Image())
}

func (ops *OptionsState) Init(e *engine.Engine) error {
	w, h := e.Size()
	v := e.Registry.Get("volume", 10.0)
	ops.volume = v.(float64)
	ops.vbar = progressBar(e, 10)
	object.Middle(w, h, ops.vbar)
	for i := 0; i < int(ops.volume); i++ {
		ops.vbar.Progress()
	}
	ops.instructions = object.FromText(e.Font(), "ENTER TO CONFIRM, BACKSPACE TO EXIT", color.White)
	_, sy := ops.instructions.Size()
	ops.instructions.SetPosition(0, h-float64(sy)-10)
	object.CenterH(w, ops.instructions)
	ops.vlabel = object.FromText(e.Font(), "VOLUME", color.White)
	object.Middle(w, h, ops.vlabel)
	object.Offset(ops.vlabel, ops.vbar, 50)

	return nil
}

func (ops *OptionsState) Register(e *engine.Engine) error {
	e.RegisterKey(ebiten.KeyBackspace, func(e *engine.Engine) {
		e.PopState()
	})
	e.RegisterKey(ebiten.KeyLeft, ops.volumedown)
	e.RegisterKey(ebiten.KeyRight, ops.volumeup)
	e.RegisterKey(ebiten.KeyEnter, ops.save)

	return nil
}

func (ops *OptionsState) save(e *engine.Engine) {
	e.Registry.Save("volume", ops.volume)
}

func (ops *OptionsState) volumedown(e *engine.Engine) {
	if ops.volume == 0 {
		return
	}
	ops.volume -= 1
	e.Audio.SetVolume(ops.volume * 0.1)
	ops.vbar.Rewind()
}

func (ops *OptionsState) volumeup(e *engine.Engine) {
	if ops.volume == 10 {
		return
	}
	ops.volume += 1
	e.Audio.SetVolume(ops.volume * 0.1)
	ops.vbar.Progress()
}
