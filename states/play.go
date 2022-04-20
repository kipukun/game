package states

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/engine/object"
	"github.com/kipukun/game/engine/tile"
)

type PlayState struct {
	player object.ImageObject
	title  *object.Pinner
	sheet  *tile.TileSheet
	world  object.ImageObject
}

func (ps *PlayState) Update(e *engine.Engine) error {
	return nil
}

func (ps *PlayState) Draw(e *engine.Engine, s *ebiten.Image) {
	s.DrawImage(ps.world.Image())
	s.DrawImage(ps.player.Image())
	s.DrawImage(ps.title.Image(e.Camera.Pos()))
}

func (ps *PlayState) Init(e *engine.Engine) {
	w, h := e.Size()

	ps.sheet = tile.NewTileSheetFromTSX(e.Asset("assets/tiles/tile_sheet.png"), e.Asset("assets/tiles/tile_set.tsx"))
	ps.world = object.FromEbitenImage(tile.NewTileMapFromTMX(ps.sheet, e.Asset("assets/tiles/tile_map.tmx")))
	fmt.Println(ps.world.Size())

	player := ebiten.NewImage(10, 10)
	player.Fill(color.White)
	ps.player = object.FromEbitenImage(player)
	object.Middle(w, h, ps.player)

	title := object.FromText(e.Font(), "PLAYER 1\nTROPHIES: 0\n$: 0", color.White)
	object.CenterH(w, title)
	tx, ty := title.GetPosition()
	title.SetPosition(tx, ty+40)
	ps.title = object.NewPinner(title)
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
