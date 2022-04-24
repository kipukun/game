package states

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/engine/object"
	"github.com/kipukun/game/engine/translation"
)

type IntroTitleState struct {
	music, menu_move, sel *engine.AudioPlayer
	title                 object.ImageObject
	titleEaser            func(dt, x, y float64)
	menu                  []*translation.Fader
	pointer               *translation.Fader
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
	its.menu = make([]*translation.Fader, 0)

	w, h := e.Size()

	for i, item := range menu {
		o := object.FromText(e.Font(), item, color.White)
		nx, ny := object.Middle(w, h, o)
		o.SetPosition(nx, ny+30*float64(i))
		its.menu = append(its.menu, translation.NewFader(o))
	}

	pointer := object.FromText(e.Font(), ">", color.White)
	fromx, fromy := its.menu[0].GetPosition()
	pointer.SetPosition(fromx-30, fromy)

	its.px, its.py = pointer.GetPosition()
	its.pointer = translation.NewFader(pointer)

	its.title = object.FromText(e.Font(), "JRPG", color.White)
	nx, _ := object.CenterH(w, its.title)
	its.title.SetPosition(nx, h)
	its.titleEaser = translation.Easer(its.title, translation.EaseInOutCubic, 3*time.Second)

	its.music.Play()
}

func (its *IntroTitleState) Update(e *engine.Engine, dt float64) error {
	_, h := e.Size()

	its.titleEaser(dt, 0, -h+40)

	for _, o := range its.menu {
		o.Calculate(nil)
	}
	its.pointer.Calculate(nil)
	return nil
}

func (its *IntroTitleState) Draw(e *engine.Engine, s *ebiten.Image) {
	s.DrawImage(its.title.Image())
	for _, o := range its.menu {
		s.DrawImage(o.Image())
	}
	s.DrawImage(its.pointer.Image())
	ebitenutil.DebugPrint(s, fmt.Sprintf("%s", ebiten.GamepadName(0)))
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
	if its.index+1 > len(its.menu)-1 {
		its.index = 0
	} else {
		its.index += 1
	}
	its.pointer.SetPosition(its.px, its.py+float64(its.index)*30)
}

func (its *IntroTitleState) menuup(e *engine.Engine) {
	its.menu_move.Play()
	if its.index-1 < 0 {
		its.index = len(its.menu) - 1
	} else {
		its.index -= 1
	}
	its.pointer.SetPosition(its.px, its.py+float64(its.index)*30)
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
