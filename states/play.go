package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kipukun/game/engine"
)

type PlayState struct {
}

func (ps *PlayState) Update(e *engine.Engine) error {
	return nil
}

func (ps *PlayState) Draw(e *engine.Engine, s *ebiten.Image) {
	ebitenutil.DebugPrint(s, "play")
}

func (ps *PlayState) Init(e *engine.Engine) {
	e.RegisterKey(ebiten.KeyBackspace, func(e *engine.Engine) {
		e.PopState()
	})
}
