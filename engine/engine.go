package engine

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type bscloser struct {
	*bytes.Reader
}

func (b *bscloser) Close() error {
	return nil
}

// State defines the state of the game engine at some time.
// A State is able to initialize itself, and change the state of the engine.
type State interface {
	Update(e *Engine) error
	Draw(e *Engine, s *ebiten.Image)
	// Init is called when the State is first pushed onto the Engine stack.
	Init(e *Engine)

	// Register is called everytime State becomes the active state.
	Register(e *Engine)
}

// Engine is the main game engine, which implements
// the ebiten.Game interface and maintains a stack of states.
type Engine struct {
	conf   *Config
	fs     fs.FS
	states []State
	Audio  *audioHandler
	tf     font.Face
	ikh    *inputHandler[ebiten.Key]
	igph   *inputHandler[ebiten.GamepadButton]
	*Registry
}

func (e *Engine) RegisterKey(k ebiten.Key, f func(e *Engine)) {
	e.ikh.handle(k, f)
}

func (e *Engine) RegisterButton(b ebiten.GamepadButton, f func(e *Engine)) {
	e.igph.handle(b, f)
}

// Deregister is called by a State when the engine should keep its handlers on state change.
func (e *Engine) KeepHandlers() {
	e.ikh.keepFlag = true
	e.igph.keepFlag = true
}

func (e *Engine) Asset(path string) (io.ReadSeekCloser, error) {
	return asset(e.fs, path)
}

// Font returns the configured font-face for the engine.
func (e *Engine) Font() font.Face {
	return e.tf
}

func (e *Engine) Size() (float64, float64) {
	return float64(e.conf.Width), float64(e.conf.Height)
}

// Init initializes the game window.
func (e *Engine) Init(ctx context.Context, c *Config, fsys fs.FS) error {
	var err error
	e.fs = fsys
	e.conf = c
	ah := &audioHandler{
		actx:  audio.NewContext(c.Samplerate),
		sr:    c.Samplerate,
		files: e.fs,
	}
	e.Audio = ah
	ebiten.SetWindowTitle(c.Name)
	ebiten.SetWindowSize(c.Width*2, c.Height*2)

	e.ikh = newInputHandler[ebiten.Key]()
	e.igph = newInputHandler[ebiten.GamepadButton]()

	e.Registry = newRegistry(e.conf.SaveFile)

	f, err := e.fs.Open("assets/fonts/font.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	ttf, err := io.ReadAll(f)
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

func (e *Engine) changed() {
	e.ikh.deregister()
	e.igph.deregister()
}

// ChangeState sets the currently running state to s.
func (e *Engine) ChangeState(s State) {
	e.changed()
	s.Init(e)
	s.Register(e)
	if len(e.states) == 0 {
		e.states = make([]State, 1)
	}
	e.states[len(e.states)-1] = s
}

// PushState appends s to the top of the stack.
func (e *Engine) PushState(s State) {
	e.changed()
	s.Init(e)
	s.Register(e)
	e.states = append(e.states, s)
}

// PopState removes the state at the top of the stack.
func (e *Engine) PopState() {
	e.changed()
	e.states = e.states[:len(e.states)-1]
	head(e.states).Register(e)
}

// Update implements ebiten.Game
func (e *Engine) Update() error {
	// run key handlers
	e.ikh.run(e)
	e.igph.run(e)
	// let the current state update the engine.
	return head(e.states).Update(e)
}

// Draw implements ebiten.Game (kind of)
func (e *Engine) Draw(screen *ebiten.Image) {
	// let the current state draw to the screen.
	head(e.states).Draw(e, screen)
}

// Layout implements ebiten.Game
func (e *Engine) Layout(ow, oh int) (int, int) {
	return e.conf.Width, e.conf.Height
}

// Run runs the Engine.
func (e *Engine) Run() error {
	return ebiten.RunGame(e)
}
