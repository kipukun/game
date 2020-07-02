package main

import (
	"log"

	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/states"
)

const (
	WIDTH      = 320
	HEIGHT     = 240
	SAMPLERATE = 44100
)

func main() {
	e := new(engine.Engine)
	err := e.Init("jayarrpeegee", WIDTH, HEIGHT, SAMPLERATE)
	if err != nil {
		log.Fatal(err)
	}

	i := new(states.PlayState)
	e.PushState(i)

	if err = e.Run(); err != nil {
		log.Fatal(err)
	}
}
