package engine

import (
	"io/ioutil"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/markbates/pkger"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
)

// TileSize is the size in pixels of each side of a tile.
// A tile is not necessary, but helps keep graphics consistent in size.
const TileSize = 16

// State defines the state of the game engine at some time.
// A State is able to initialize itself, and change the state of the engine.
type State interface {
	Update(e *Engine) error
	Draw(e *Engine, s *ebiten.Image)
	Init(e *Engine)
}

// Engine is the main game engine, which implements
// the ebiten.Game interface and maintains a stack of states.
type Engine struct {
	states                              []State
	audioCtx                            *audio.Context
	tf                                  font.Face
	width, height, tilesize, samplerate int
}

// AudioCtx returns the engine's audio context.
func (e *Engine) AudioCtx() *audio.Context {
	return e.audioCtx
}

// Font returns the configured font-face for the engine.
func (e *Engine) Font() font.Face {
	return e.tf
}

func (e *Engine) Size() (float64, float64) {
	return float64(e.width), float64(e.height)
}

// Init initializes the game window.
func (e *Engine) Init(name string, w, h, sr int) error {
	e.audioCtx, _ = audio.NewContext(sr)
	ebiten.SetWindowTitle(name)
	ebiten.SetWindowSize(w*2, h*2)
	e.width, e.height, e.samplerate = w, h, sr
	e.tilesize = 16

	f, err := pkger.Open("/assets/fonts/font.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	ttf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	tt, err := truetype.Parse(ttf)
	if err != nil {
		return err
	}
	e.tf = truetype.NewFace(tt, &truetype.Options{
		Size: 9,
		DPI:  72,
	})
	return nil
}

// ChangeState sets the currently running state to s.
func (e *Engine) ChangeState(s State) {
	s.Init(e)
	e.states[len(e.states)-1] = s
}

// PushState appends s to the top of the stack.
func (e *Engine) PushState(s State) {
	s.Init(e)
	e.states = append(e.states, s)
}

// PopState removes the state at the top of the stack.
func (e *Engine) PopState() {
	idx := len(e.states) - 1
	e.states = e.states[:idx]
}

// Update implements ebiten.Game
func (e *Engine) Update(screen *ebiten.Image) error {
	// let the current state draw to the screen.
	return e.states[len(e.states)-1].Update(e)
}

// Draw implements ebiten.Game (kind of)
func (e *Engine) Draw(screen *ebiten.Image) {
	// let the current state draw to the screen.
	e.states[len(e.states)-1].Draw(e, screen)
}

// Layout implements ebiten.Game
func (e *Engine) Layout(ow, oh int) (int, int) {
	return e.width, e.height
}

// Run runs the Engine.
func (e *Engine) Run() error {
	return ebiten.RunGame(e)
}
