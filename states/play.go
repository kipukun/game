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
	title  *engine.Pinner
}

func (ps *PlayState) Update(e *engine.Engine) error {
	return nil
}

func (ps *PlayState) Draw(e *engine.Engine, s *ebiten.Image) {
	s.DrawImage(ps.title.Image())
	s.DrawImage(ps.player.Image())
	ebitenutil.DebugPrint(s, "play")
}

func (ps *PlayState) Init(e *engine.Engine) {
	w, h := e.Size()
	player := ebiten.NewImage(10, 10)
	player.Fill(color.White)
	ps.player = engine.FromEbitenImage(player)
	title := engine.FromText(e, "HEALTH: <3 <3 <3", color.White)
	object.CenterH(w, title)
	tx, ty := title.GetPosition()
	title.SetPosition(tx, ty+40)
	ps.title = engine.NewPinner(e.Camera, title)

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
	e.RegisterHeldKey(ebiten.KeyA, func(e *engine.Engine) {
		x, y := e.Camera.GetPosition()
		e.Camera.SetPosition(x-1, y)
	})
	e.RegisterHeldKey(ebiten.KeyD, func(e *engine.Engine) {
		x, y := e.Camera.GetPosition()
		e.Camera.SetPosition(x+1, y)
	})
	e.RegisterHeldKey(ebiten.KeyW, func(e *engine.Engine) {
		x, y := e.Camera.GetPosition()
		e.Camera.SetPosition(x, y-1)
	})
	e.RegisterHeldKey(ebiten.KeyS, func(e *engine.Engine) {
		x, y := e.Camera.GetPosition()
		e.Camera.SetPosition(x, y+1)
	})
}
