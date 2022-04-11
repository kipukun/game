package engine

import (
	"bytes"
	"io"
	"io/fs"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type Bytes struct {
	*bytes.Reader
}

func (b *Bytes) Close() error {
	return nil
}

// TileSize is the square size in pixels of a tile.
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
	fs                        fs.FS
	states                    []State
	audioCtx                  *audio.Context
	tf                        font.Face
	width, height, samplerate int
}

func (e *Engine) Assets(fs fs.FS) {
	e.fs = fs
}

func (e *Engine) Asset(path string) (io.ReadSeekCloser, error) {
	bs, err := fs.ReadFile(e.fs, path)
	if err != nil {
		return nil, err
	}
	return &Bytes{bytes.NewReader(bs)}, nil
}

// Player is a helper method to give an audio.Player from a wav asset.
func (e *Engine) Player(path string) (*audio.Player, error) {
	fd, err := e.Asset(path)
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	return e.AudioCtx().NewPlayerFromBytes(bs), nil
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
	var err error
	e.audioCtx = audio.NewContext(sr)
	ebiten.SetWindowTitle(name)
	ebiten.SetWindowSize(w*2, h*2)
	e.width, e.height, e.samplerate = w, h, sr

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
func (e *Engine) Update() error {
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
