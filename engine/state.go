package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// State defines the state of the game engine at some time.
// A State is able to initialize itself, and change the state of the engine.
type State interface {
	Update(e *Engine, dt float64) error
	Draw(e *Engine, s *ebiten.Image)

	// Init is called when the State is first pushed onto the engine stack.
	Init(e *Engine) error

	// Register is called everytime this state becomes the active state.
	Register(e *Engine) error

	fmt.Stringer
}
