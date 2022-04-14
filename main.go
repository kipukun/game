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

func main() {
	e := new(engine.Engine)
	e.Assets(assets)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill)
	defer cancel()

	c, err := engine.NewConfigFromFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	err = e.Init(ctx, c)
	if err != nil {
		log.Fatal(err)
	}

	i := new(states.IntroTitleState)
	e.PushState(i)

	if err = e.Run(); err != nil {
		log.Fatal(err)
	}
}
