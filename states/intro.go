package states

import (
	"log"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/kipukun/game/engine"
)

type IntroState struct {
	music *audio.Player
}

func (i *IntroState) Init(e *engine.Engine) {
	fd, err := e.Asset("assets/audio/lebron_james.wav")
	if err != nil {
		log.Fatal(err)
	}
	d, err := wav.Decode(e.AudioCtx(), fd)
	if err != nil {
		log.Fatal(err)
	}
	i.music, err = audio.NewPlayer(e.AudioCtx(), d)
	if err != nil {
		log.Fatal(err)
	}
}
