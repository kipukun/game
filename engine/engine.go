package engine

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"log"
	"sync"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	fs                        fs.FS
	states                    []State
	audioCtx                  *audio.Context
	tf                        font.Face
	width, height, samplerate int
	ikh                       *inputHandler[ebiten.Key]
	igph                      *inputHandler[ebiten.GamepadButton]
	keepFlag                  bool
}

type input interface {
	ebiten.Key | ebiten.GamepadButton
}

type inputHandler[T input] struct {
	keys             map[T][]func(e *Engine)
	currentStateKeys []T
}

func (ih *inputHandler[T]) handle(k T, f func(e *Engine)) {
	if ih.keys[k] == nil {
		ih.keys[k] = make([]func(e *Engine), 0)
	}
	if ih.currentStateKeys == nil {
		ih.currentStateKeys = make([]T, 0)
	}
	ih.currentStateKeys = append(ih.currentStateKeys, k)
	ih.keys[k] = append(ih.keys[k], f)
}

func (ih *inputHandler[T]) deregister() {
	if len(ih.currentStateKeys) < 1 {
		return
	}
	for _, k := range ih.currentStateKeys {
		if len(ih.keys[k]) < 1 {
			return
		}
		ih.keys[k] = ih.keys[k][:len(ih.keys[k])-1]
	}
	ih.currentStateKeys = nil
}

func (ih *inputHandler[T]) run(e *Engine) {
	for k, v := range ih.keys {
		// we know that T can only be a type in input
		switch any(k).(type) {
		case ebiten.Key:
			if !inpututil.IsKeyJustPressed(ebiten.Key(k)) {
				continue
			}
		case ebiten.GamepadButton:
			if !inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton(k)) {
				continue
			}
		// if it's not, idk how it got past the compiler.
		// hopefully we can get compile-time guarantees on a type switch like this.
		default:
			panic("not a type in input constraint")
		}
		for _, f := range v {
			f(e)
		}
	}
}

func newInputHandler[T input]() *inputHandler[T] {
	ih := new(inputHandler[T])
	ih.keys = make(map[T][]func(e *Engine))
	return ih
}

// AudioPlayer is a concurrent-safe wrapper around an audio.Player.
type AudioPlayer struct {
	mu sync.Mutex
	p  *audio.Player
}

func (ap *AudioPlayer) Play() {
	ap.mu.Lock()
	ap.p.Pause()
	ap.p.Rewind()
	ap.p.Play()
	ap.mu.Unlock()
}

func (e *Engine) RegisterKey(k ebiten.Key, f func(e *Engine)) {
	e.ikh.handle(k, f)
}

func (e *Engine) RegisterButton(b ebiten.GamepadButton, f func(e *Engine)) {
	e.igph.handle(b, f)
}

// Deregister is called by a State when the engine should keep its Handlers even on change.
func (e *Engine) KeepHandlers() {
	e.keepFlag = true
}

func (e *Engine) Assets(fs fs.FS) {
	e.fs = fs
}

func (e *Engine) Asset(path string) (io.ReadSeekCloser, error) {
	bs, err := fs.ReadFile(e.fs, path)
	if err != nil {
		return nil, err
	}
	return &bscloser{bytes.NewReader(bs)}, nil
}

// Player is a helper method to give an audio.Player from a wav asset.
func (e *Engine) Player(path string) (*AudioPlayer, error) {
	fd, err := e.Asset(path)
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	mp3, err := mp3.DecodeWithSampleRate(e.samplerate, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	p, err := e.AudioCtx().NewPlayer(mp3)
	if err != nil {
		return nil, err
	}
	return &AudioPlayer{p: p}, nil
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
func (e *Engine) Init(ctx context.Context, name string, w, h, sr int) error {
	var err error
	e.audioCtx = audio.NewContext(sr)
	ebiten.SetWindowTitle(name)
	ebiten.SetWindowSize(w*2, h*2)
	e.width, e.height, e.samplerate = w, h, sr

	e.ikh = newInputHandler[ebiten.Key]()
	e.igph = newInputHandler[ebiten.GamepadButton]()

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
	if e.keepFlag {
		return
	}
	e.ikh.deregister()
	e.igph.deregister()
	e.keepFlag = false
}

// ChangeState sets the currently running state to s.
func (e *Engine) ChangeState(s State) {
	e.changed()
	s.Init(e)
	s.Register(e)
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
	idx := len(e.states) - 1
	e.states = e.states[:idx]
	head(e.states).Register(e)
}

// Update implements ebiten.Game
func (e *Engine) Update() error {
	// run key handlers
	e.ikh.run(e)
	e.igph.run(e)
	// let the current state update the engine.
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
