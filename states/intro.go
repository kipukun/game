package states

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/engine/object"
	"github.com/kipukun/game/engine/transform"
)

type titleEntity struct {
	io    object.ImageObject
	easer transform.ChangeFunc
}

type menuEntity struct {
	options []object.ImageObject
	faders  []transform.ChangeFunc
}

type pointerEntity struct {
	io    object.ImageObject
	fader transform.ChangeFunc
}

type IntroTitleState struct {
	music, menu_move, sel *engine.AudioPlayer
	title                 *titleEntity
	menu                  *menuEntity
	pointer               *pointerEntity
	px, py                float64
	ch                    transform.ChangeFunc
	index                 int
}

func (its *IntroTitleState) Init(e *engine.Engine) {
	var err error
	its.music, err = e.Audio.Player("assets/audio/lebron_james.mp3")
	if err != nil {
		log.Fatal(err)
	}
	its.menu_move, err = e.Audio.Player("assets/audio/menu_move.mp3")
	if err != nil {
		log.Fatal(err)
	}
	its.sel, err = e.Audio.Player("assets/audio/select.mp3")
	if err != nil {
		log.Fatal(err)
	}
	menu := []string{"PLAY", "OPTIONS", "QUIT"}
	menuEntity := new(menuEntity)
	menuEntity.options = make([]object.ImageObject, 0)
	menuEntity.faders = make([]transform.ChangeFunc, 0)

	w, h := e.Size()

	for i, item := range menu {
		o := object.FromText(e.Font(), item, color.White)
		nx, ny := object.Middle(w, h, o)
		o.SetPosition(nx, ny+30*float64(i))
		menuEntity.options = append(menuEntity.options, o)
		menuEntity.faders = append(menuEntity.faders, transform.Fader(o, transform.Linear, time.Second))
		object.Transparent(o)
	}
	its.menu = menuEntity

	its.pointer = new(pointerEntity)
	its.pointer.io = object.FromText(e.Font(), ">", color.White)
	object.Transparent(its.pointer.io)
	its.pointer.fader = transform.Fader(its.pointer.io, transform.Linear, time.Second)
	fromx, fromy := its.menu.options[0].GetPosition()
	its.pointer.io.SetPosition(fromx-30, fromy)
	its.px, its.py = its.pointer.io.GetPosition()

	its.title = new(titleEntity)
	its.title.io = object.FromText(e.Font(), "JRPG", color.White)
	nx, _ := object.CenterH(w, its.title.io)
	its.title.io.SetPosition(nx, h)

	its.title.easer = transform.Easer(its.title.io, 0, -h+40, transform.EaseInOutCubic, 5*time.Second)

	its.ch = transform.Chain(
		its.title.easer,
		transform.Combine(
			its.pointer.fader,
			transform.Combine(its.menu.faders...)),
	)

	its.music.Play()
}

func (its *IntroTitleState) Update(e *engine.Engine, dt float64) error {
	its.ch(dt)
	return nil
}

func (its *IntroTitleState) Draw(e *engine.Engine, s *ebiten.Image) {
	s.DrawImage(its.title.io.Image())
	for _, o := range its.menu.options {
		s.DrawImage(o.Image())
	}
	s.DrawImage(its.pointer.io.Image())
	ebitenutil.DebugPrint(s, ebiten.GamepadName(0))
}

func (its *IntroTitleState) Register(e *engine.Engine) {
	e.Camera.Reset()

	e.RegisterKey(ebiten.KeyArrowDown, its.menudown)
	e.RegisterKey(ebiten.KeyArrowUp, its.menuup)
	e.RegisterKey(ebiten.KeyEnter, its.o)
	e.RegisterButton(ebiten.GamepadButton0, its.menudown)
	e.RegisterButton(ebiten.GamepadButton1, its.menuup)
	e.RegisterButton(ebiten.GamepadButton2, its.o)
}

func (its *IntroTitleState) menudown(e *engine.Engine) {
	its.menu_move.Play()
	if its.index+1 > len(its.menu.options)-1 {
		its.index = 0
	} else {
		its.index += 1
	}
	its.pointer.io.SetPosition(its.px, its.py+float64(its.index)*30)
}

func (its *IntroTitleState) menuup(e *engine.Engine) {
	its.menu_move.Play()
	if its.index-1 < 0 {
		its.index = len(its.menu.options) - 1
	} else {
		its.index -= 1
	}
	its.pointer.io.SetPosition(its.px, its.py+float64(its.index)*30)
}

func (its *IntroTitleState) o(e *engine.Engine) {
	its.sel.Play()
	var s engine.State
	switch its.index {
	case 0:
		s = new(PlayState)
	case 1:
		s = new(OptionsState)
	default:
		return
	}
	e.PushState(s)
}
