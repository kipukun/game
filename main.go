package main

import (
	"log"
)

const (
	WIDTH      = 320
	HEIGHT     = 240
	TILESIZE   = 16
	SAMPLERATE = 44100
)

func main() {
	e := new(Engine)
	e.Init("jayarrpeegee", WIDTH*2, HEIGHT*2)

	i := new(intro)
	e.PushState(i)

	if err := e.Run(); err != nil {
		log.Fatal(err)
	}
}
