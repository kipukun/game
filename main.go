package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"

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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill)
	defer cancel()

	err := e.Init(ctx, "jayarrpeegee", WIDTH, HEIGHT, SAMPLERATE)
	if err != nil {
		log.Fatal(err)
	}

	i := new(states.IntroTitleState)
	e.PushState(i)

	if err = e.Run(); err != nil {
		log.Fatal(err)
	}
}
