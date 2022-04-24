package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine"
)

type InitState struct {
}

func (is *InitState) Update(e *engine.Engine, dt float64) error {
	its := new(IntroTitleState)
	e.ChangeState(its)
	return nil
}

func (is *InitState) Draw(e *engine.Engine, s *ebiten.Image) {
}

// Init is called when the State is first pushed onto the Engine stack.
func (is *InitState) Init(e *engine.Engine) {
	e.Registry.Load()
	v := e.Registry.Get("volume", 10.0)
	e.Audio.SetVolume(v.(float64) * 0.1)
}

// Register is called everytime State becomes the active state.
func (is *InitState) Register(e *engine.Engine) {
}
