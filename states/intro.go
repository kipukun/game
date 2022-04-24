package states

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/engine/object"
	"github.com/kipukun/game/engine/translation"
)

type titleEntity struct {
	io    object.ImageObject
	easer translation.EaserFunc
}

type menuEntity struct {
	opt     *ebiten.DrawImageOptions
	options []object.ImageObject
}

type pointerEntity struct {
	io object.ImageObject
}

type IntroTitleState struct {
	music, menu_move, sel *engine.AudioPlayer
	title                 *titleEntity
	menu                  *menuEntity
	pointer               *pointerEntity
	fader                 translation.FaderFunc
	px, py                float64
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

	its.fader = translation.Fader(translation.Linear, 3*time.Second)

	w, h := e.Size()

	for i, item := range menu {
		o := object.FromText(e.Font(), item, color.White)
		nx, ny := object.Middle(w, h, o)
		o.SetPosition(nx, ny+30*float64(i))
		menuEntity.options = append(menuEntity.options, o)
	}
	its.menu = menuEntity

	its.pointer = new(pointerEntity)
	its.pointer.io = object.FromText(e.Font(), ">", color.White)
	fromx, fromy := its.menu.options[0].GetPosition()
	its.pointer.io.SetPosition(fromx-30, fromy)
	its.px, its.py = its.pointer.io.GetPosition()

	its.title = new(titleEntity)
	its.title.io = object.FromText(e.Font(), "JRPG", color.White)
	nx, _ := object.CenterH(w, its.title.io)
	its.title.io.SetPosition(nx, h)

	its.title.easer = translation.Easer(its.title.io, translation.EaseInOutCubic, 3*time.Second)

	its.music.Play()
}

func (its *IntroTitleState) Update(e *engine.Engine, dt float64) error {
	_, h := e.Size()
	its.title.easer(dt, 0, -h+40)
	its.fader(dt)
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
