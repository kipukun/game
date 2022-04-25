package engine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	timeNow, timeLast time.Time
	delta             time.Duration
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
	Update(e *Engine, dt float64) error
	Draw(e *Engine, s *ebiten.Image)

	// Init is called when the State is first pushed onto the engine stack.
	Init(e *Engine)

	// Register is called everytime this state becomes the active state.
	Register(e *Engine)
}

// Engine is the main game engine, which implements
// the ebiten.Game interface and maintains a stack of states.
type Engine struct {
	conf        *Config
	debug       bool
	fs          fs.FS
	states      []State
	Audio       *audioHandler
	tf          font.Face
	keyh, hkeyh *inputHandler[ebiten.Key]
	gph, hgph   *inputHandler[ebiten.GamepadButton]
	Camera      *Camera
	vscreen     *ebiten.Image
	*Registry
}

// RegisterKey registers f to be called when k is pressed during the current frame.
func (e *Engine) RegisterKey(k ebiten.Key, f func(e *Engine)) {
	e.keyh.handle(k, f)
}

// RegisterKey registers f to be called when b is pressed during the current frame.
func (e *Engine) RegisterButton(b ebiten.GamepadButton, f func(e *Engine)) {
	e.gph.handle(b, f)
}

// RegisterHeldKey registers f to be called when k is held down.
func (e *Engine) RegisterHeldKey(k ebiten.Key, f func(e *Engine)) {
	e.hkeyh.handle(k, f)
}

// RegisterHeldKey registers f to be called when b is held down.
func (e *Engine) RegisterHeldButton(b ebiten.GamepadButton, f func(e *Engine)) {
	e.hgph.handle(b, f)
}

// Deregister is called by a State when the engine should keep its handlers on state change.
func (e *Engine) KeepHandlers() {
	e.keyh.keepFlag = true
	e.gph.keepFlag = true
	e.hkeyh.keepFlag = true
	e.hgph.keepFlag = true
}

func (e *Engine) Asset(path string) io.ReadSeekCloser {
	return asset(e.fs, path)
}

// Font returns the configured font-face for the engine.
func (e *Engine) Font() font.Face {
	return e.tf
}

func (e *Engine) Size() (float64, float64) {
	return float64(e.conf.Width), float64(e.conf.Height)
}

// Init initializes the entire Engine.
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

	e.keyh = newInputHandler[ebiten.Key](false)
	e.gph = newInputHandler[ebiten.GamepadButton](false)
	e.hkeyh = newInputHandler[ebiten.Key](true)
	e.hgph = newInputHandler[ebiten.GamepadButton](true)

	e.Camera = NewCamera()
	e.vscreen = ebiten.NewImage(c.Width*2, c.Height*2)

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

	e.KeepHandlers()
	e.RegisterKey(ebiten.KeyEscape, func(e *Engine) {
		e.ToggleDebugMode()
	})

	return nil
}

func (e *Engine) changed() {
	e.keyh.deregister()
	e.gph.deregister()
	e.hkeyh.deregister()
	e.hgph.deregister()
}

func (e *Engine) ToggleDebugMode() {
	e.debug = !e.debug
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
	e.keyh.run(e)
	e.gph.run(e)
	e.hkeyh.run(e)
	e.hgph.run(e)
	timeNow = time.Now()
	delta = timeNow.Sub(timeLast)
	// let the current state update the engine.
	err := head(e.states).Update(e, delta.Seconds())
	timeLast = timeNow
	return err
}

// Draw implements ebiten.Game (kind of)
func (e *Engine) Draw(screen *ebiten.Image) {
	// let the current state draw to the screen.
	e.vscreen.Clear()
	head(e.states).Draw(e, e.vscreen)
	screen.DrawImage(e.vscreen, e.Camera.View())

	if e.debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f, TPS: %f, DT: %f", ebiten.CurrentFPS(), ebiten.CurrentTPS(), delta.Seconds()))
	}
}

// Layout implements ebiten.Game
func (e *Engine) Layout(ow, oh int) (int, int) {
	return e.conf.Width, e.conf.Height
}

// Run runs the Engine.
func (e *Engine) Run() error {
	return ebiten.RunGame(e)
}
