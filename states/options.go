package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kipukun/game/engine"
)

type OptionsState struct {
	volume float64
}

func (ops *OptionsState) Update(e *engine.Engine) error {
	return nil
}

func (ops *OptionsState) Draw(e *engine.Engine, s *ebiten.Image) {
	ebitenutil.DebugPrint(s, fmt.Sprintf("options, v: %f", ops.volume))
}

func (ops *OptionsState) Init(e *engine.Engine) {
	v := e.Registry.Get("volume", 10.0)
	ops.volume = v.(float64)
}

func (ops *OptionsState) Register(e *engine.Engine) {
	e.RegisterKey(ebiten.KeyBackspace, func(e *engine.Engine) {
		e.PopState()
	})
	e.RegisterKey(ebiten.KeyLeft, ops.volumedown)
	e.RegisterKey(ebiten.KeyRight, ops.volumeup)
	e.RegisterKey(ebiten.KeyEnter, ops.save)
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
}

func (ops *OptionsState) volumeup(e *engine.Engine) {
	if ops.volume == 10 {
		return
	}
	ops.volume += 1
	e.Audio.SetVolume(ops.volume * 0.1)
}
