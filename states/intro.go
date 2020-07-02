package states

import (
	"log"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/kipukun/game/engine"
	"github.com/markbates/pkger"
)

type IntroState struct {
	music *audio.Player
}

func (i *IntroState) Init(e *engine.Engine) {
	fd, err := pkger.Open("/assets/audio/lebron_james.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	d, err := wav.Decode(e.AudioCtx(), fd)
	if err != nil {
		log.Fatal(err)
	}
	i.music, err = audio.NewPlayer(e.AudioCtx(), d)
	if err != nil {
		log.Fatal(err)
	}
}
