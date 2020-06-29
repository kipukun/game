package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	_ "image/png"
	"log"
	"math"
	"os"
)

type intro struct {
	grass, cloud *ebiten.Image
	num          int
	audioPlayer  *audio.Player
}

func (i *intro) Init(e *Engine) {
	i.grass, _, _ = ebitenutil.NewImageFromFile(`sprites\grass.png`, ebiten.FilterDefault)
	i.cloud, _, _ = ebitenutil.NewImageFromFile(`sprites\cloud.png`, ebiten.FilterDefault)
	i.num = int(math.Floor(WIDTH / TILESIZE))
	f, err := os.Open(`audio\pause.wav`)
	if err != nil {
		log.Fatal(err)
	}
	d, err := wav.Decode(e.audioCtx, f)
	if err != nil {
		log.Fatal(err)
	}
	i.audioPlayer, err = audio.NewPlayer(e.audioCtx, d)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *intro) Update(e *Engine) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		i.audioPlayer.Pause()
		i.audioPlayer.Rewind()
		i.audioPlayer.Play()
		e.PushState(new(pause))
	}
	return nil
}

func (i *intro) Draw(e *Engine, screen *ebiten.Image) {
	co := &ebiten.DrawImageOptions{}
	co.GeoM.Translate(WIDTH/2, HEIGHT/4)
	screen.DrawImage(i.cloud, co)
	for j := 0; j < i.num; j++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(j*TILESIZE), HEIGHT/1.5)
		screen.DrawImage(i.grass, op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}

type pause struct {
	text        *ebiten.Image
	audioPlayer *audio.Player
}

func (p *pause) Init(e *Engine) {
	var err error
	p.text, _, _ = ebitenutil.NewImageFromFile(`sprites\PAUSE.png`, ebiten.FilterDefault)
	f, err := os.Open(`audio\pause_exit.wav`)
	if err != nil {
		log.Fatal(err)
	}
	d, err := wav.Decode(e.audioCtx, f)
	if err != nil {
		log.Fatal(err)
	}
	p.audioPlayer, err = audio.NewPlayer(e.audioCtx, d)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *pause) Update(e *Engine) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.audioPlayer.Pause()
		p.audioPlayer.Rewind()
		p.audioPlayer.Play()
		e.PopState()
	}
	return nil
}

func (p *pause) Draw(e *Engine, screen *ebiten.Image) {
	o := &ebiten.DrawImageOptions{}
	o.GeoM.Translate(WIDTH/2, HEIGHT/2)
	screen.DrawImage(p.text, o)
}
