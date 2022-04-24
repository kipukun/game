package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine"
)

type BattleState struct{}

func (bs *BattleState) Update(e *engine.Engine) error {
	return nil
}

func (bs *BattleState) Draw(e *engine.Engine, s *ebiten.Image) {
}

// Init is called when the State is first pushed onto the engine stack.
func (bs *BattleState) Init(e *engine.Engine) {
}

// Register is called everytime this state becomes the active state.
func (bs *BattleState) Register(e *engine.Engine) {
}
