package states

import (
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kipukun/game/engine"
)

type PlayState struct {
	cloud       *ebiten.Image
	audioPlayer *audio.Player
}

func (p *PlayState) Init(e *engine.Engine) {
	f, err := e.Asset("assets/sprites/cloud.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	p.cloud, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	fd, err := e.Asset("assets/audio/pause.wav")
	if err != nil {
		log.Fatal(err)
	}
	d, err := wav.Decode(e.AudioCtx(), fd)
	if err != nil {
		log.Fatal(err)
	}
	p.audioPlayer, err = audio.NewPlayer(e.AudioCtx(), d)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *PlayState) Update(e *engine.Engine) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.audioPlayer.Pause()
		p.audioPlayer.Rewind()
		p.audioPlayer.Play()
		// e.PushState(new(MenuState))
	}
	return nil
}

func (p *PlayState) Draw(e *engine.Engine, screen *ebiten.Image) {
	w, h := e.Size()
	co := &ebiten.DrawImageOptions{}
	co.GeoM.Translate(w/2, h/4)
	screen.DrawImage(p.cloud, co)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}
