package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type errorState struct {
	err error
}

func (es *errorState) Update(e *Engine, dt float64) error {
	return es.err
}

func (es *errorState) Draw(e *Engine, s *ebiten.Image) {
}

// Init is called when the State is first pushed onto the engine stack.
func (es *errorState) Init(e *Engine) error {
	return nil
}

// Register is called everytime this state becomes the active state.
func (es *errorState) Register(e *Engine) error {
	return nil
}
