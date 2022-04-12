package states

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/engine/object"
)

type IntroTitleState struct {
	music, menu_move *engine.AudioPlayer
	title            *object.Easer[object.ImageObject]
	menu             []*object.Fader
	pointer          *object.Fader
	px, py           float64
	index            int
}

func (its *IntroTitleState) Init(e *engine.Engine) {
	var err error
	its.music, err = e.Player("assets/audio/lebron_james.mp3")
	if err != nil {
		log.Fatal(err)
	}
	its.menu_move, err = e.Player("assets/audio/menu_move.mp3")
	if err != nil {
		log.Fatal(err)
	}
	menu := []string{"PLAY", "OPTIONS", "QUIT"}
	its.menu = make([]*object.Fader, 0)

	w, h := e.Size()

	for i, item := range menu {
		o := object.FromText(e, item, color.White)
		nx, ny := object.Middle(w, h, o)
		o.SetPosition(nx, ny+30*float64(i))
		its.menu = append(its.menu, object.NewFader(o))
	}

	pointer := object.FromText(e, ">", color.White)
	fromx, fromy := its.menu[0].GetPosition()
	pointer.SetPosition(fromx-30, fromy)

	its.px, its.py = pointer.GetPosition()
	its.pointer = object.NewFader(pointer)

	its.title = object.NewEaser(object.FromText(e, "JRPG", color.White), -h+40)
	nx, _ := object.CenterH(w, its.title.O)
	its.title.O.SetPosition(nx, h)

	e.RegisterKey(ebiten.KeyArrowDown, its.menudown)
	e.RegisterKey(ebiten.KeyArrowUp, its.menuup)

	its.music.Play()
}

func (its *IntroTitleState) Update(e *engine.Engine) error {
	its.title.Calculate(func() {
		for _, o := range its.menu {
			o.Calculate(nil)
		}
		its.pointer.Calculate(nil)
	})
	return nil
}

func (its *IntroTitleState) Draw(e *engine.Engine, s *ebiten.Image) {
	s.DrawImage(its.title.O.Image())
	for _, o := range its.menu {
		s.DrawImage(o.Image())
	}
	s.DrawImage(its.pointer.Image())
	ebitenutil.DebugPrint(s, fmt.Sprintf("%d", its.index))
}

func (its *IntroTitleState) menudown() {
	its.menu_move.Play()
	if its.index+1 > len(its.menu)-1 {
		its.index = 0
	} else {
		its.index += 1
	}
	its.pointer.SetPosition(its.px, its.py+float64(its.index)*30)
}

func (its *IntroTitleState) menuup() {
	its.menu_move.Play()
	if its.index-1 < 0 {
		its.index = len(its.menu) - 1
	} else {
		its.index -= 1
	}
	its.pointer.SetPosition(its.px, its.py+float64(its.index)*30)
}
