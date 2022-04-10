package states

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/kipukun/game/engine"
)

type IntroState struct {
	music *audio.Player
}

func (i *IntroState) Init(e *engine.Engine) {
	var err error
	i.music, err = e.Player("assets/audio/lebron_james.wav")
	if err != nil {
		log.Fatal(err)
	}
	i.music.Play()
}

func (i *IntroState) Update(e *engine.Engine) error {
	return nil
}

func (i *IntroState) Draw(e *engine.Engine, s *ebiten.Image) {
	ebitenutil.DebugPrint(s, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}
