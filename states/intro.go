package states

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/engine/object"
)

type IntroTitleState struct {
	music *audio.Player
	title *object.Easer[object.ImageObject]
	menu  []*object.Fader
	index int
}

func (its *IntroTitleState) Init(e *engine.Engine) {
	var err error
	its.music, err = e.Player("assets/audio/lebron_james.wav")
	if err != nil {
		log.Fatal(err)
	}
	menu := []string{"PLAY", "OPTIONS", "QUIT"}

	w, h := e.Size()

	for i, item := range menu {
		o := object.FromText(e, item, color.White)
		nx, ny := object.Middle(w, h, o)
		o.SetPosition(nx, ny+30*float64(i))
		its.menu = append(its.menu, object.NewFader(o))
	}

	its.title = object.NewEaser(object.FromText(e, "JRPG", color.White), -h+40)
	nx, _ := object.CenterH(w, its.title.O)
	its.title.O.SetPosition(nx, h)

	its.music.Play()
}

func (its *IntroTitleState) Update(e *engine.Engine) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if its.index+1 > len(its.menu)-1 {
			its.index = 0
		} else {
			its.index += 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if its.index-1 < 0 {
			its.index = len(its.menu) - 1
		} else {
			its.index -= 1
		}
	}
	its.title.Calculate(func() {
		for _, o := range its.menu {
			o.Calculate(nil)
		}
	})
	return nil
}

func (its *IntroTitleState) Draw(e *engine.Engine, s *ebiten.Image) {
	s.DrawImage(its.title.O.Image())
	for _, o := range its.menu {
		s.DrawImage(o.Image())
	}
	ebitenutil.DebugPrint(s, fmt.Sprintf("%d", its.index))
}
