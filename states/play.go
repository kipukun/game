package states

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/engine/object"
)

type PlayState struct {
	player      object.Object
	audioPlayer *audio.Player
}

func (p *PlayState) Init(e *engine.Engine) {
	var err error
	p.player, err = object.FromAsset(e, "assets/sprites/cloud.png")
	if err != nil {
		log.Fatal(err)
	}

	w, h := e.Size()
	p.player.SetPosition(w/2, h/2)

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
	p.player.Update()
	_, h := e.Size()
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.player.SetVelocity(0, -2)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		p.player.SetVelocity(2, 0)
	}
	if x, y := p.player.Pos(); y > h/2 {
		p.player.SetPosition(x, h/2)
	}
	p.player.SetAcceleration(0, 0.1)
	return nil
}

func (p *PlayState) Draw(e *engine.Engine, screen *ebiten.Image) {
	screen.DrawImage(p.player.Image())
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}
