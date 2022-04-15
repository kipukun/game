package states

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/engine/object"
)

type PlayState struct {
	player object.ImageObject
}

func (ps *PlayState) Update(e *engine.Engine) error {
	return nil
}

func (ps *PlayState) Draw(e *engine.Engine, s *ebiten.Image) {
	s.DrawImage(ps.player.Image())
	ebitenutil.DebugPrint(s, "play")
}

func (ps *PlayState) Init(e *engine.Engine) {
	w, h := e.Size()
	player := ebiten.NewImage(10, 10)
	player.Fill(color.White)
	ps.player = object.FromEbitenImage(player)

	object.Middle(w, h, ps.player)
}

func (ps *PlayState) Register(e *engine.Engine) {
	e.RegisterKey(ebiten.KeyBackspace, func(e *engine.Engine) {
		e.PopState()
	})
	e.RegisterHeldKey(ebiten.KeyLeft, func(e *engine.Engine) {
		x, y := ps.player.GetPosition()
		ps.player.SetPosition(x-1, y)
	})
	e.RegisterHeldKey(ebiten.KeyRight, func(e *engine.Engine) {
		x, y := ps.player.GetPosition()
		ps.player.SetPosition(x+1, y)
	})
	e.RegisterHeldKey(ebiten.KeyUp, func(e *engine.Engine) {
		x, y := ps.player.GetPosition()
		ps.player.SetPosition(x, y-1)
	})
	e.RegisterHeldKey(ebiten.KeyDown, func(e *engine.Engine) {
		x, y := ps.player.GetPosition()
		ps.player.SetPosition(x, y+1)
	})
}
