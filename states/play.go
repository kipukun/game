package states

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/engine/object"
)

type PlayState struct {
	player      object.ImageObject
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

	p.audioPlayer, err = e.Player("assets/audio/pause.wav")
	if err != nil {
		log.Fatal(err)
	}
}

func (p *PlayState) Update(e *engine.Engine) error {
	p.player.Update()
	_, h := e.Size()
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.player.SetVelocity(0, -2)
		p.audioPlayer.Play()
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
