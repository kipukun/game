package main

import (
	"embed"
	"log"

	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/states"
)

//go:embed assets/*
var assets embed.FS

const (
	WIDTH      = 320
	HEIGHT     = 240
	SAMPLERATE = 44100
)

func main() {
	e := new(engine.Engine)
	e.Assets(assets)
	err := e.Init("jayarrpeegee", WIDTH, HEIGHT, SAMPLERATE)
	if err != nil {
		log.Fatal(err)
	}

	i := new(states.IntroState)
	e.PushState(i)

	if err = e.Run(); err != nil {
		log.Fatal(err)
	}
}
