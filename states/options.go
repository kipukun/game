package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kipukun/game/engine"
)

type OptionsState struct {
}

func (ops *OptionsState) Update(e *engine.Engine) error {
	return nil
}

func (ops *OptionsState) Draw(e *engine.Engine, s *ebiten.Image) {
	ebitenutil.DebugPrint(s, "options")
}

func (ops *OptionsState) Init(e *engine.Engine) {
	e.RegisterKey(ebiten.KeyBackspace, func(e *engine.Engine) {
		e.PopState()
	})
}
