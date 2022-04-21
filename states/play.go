package states

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kipukun/game/engine"
	"github.com/kipukun/game/engine/object"
	"github.com/kipukun/game/engine/tile"
)

type player struct {
	io    object.ImageObject
	path  object.Collection
	easer func()
	idx   int
}

func (p *player) move(dx int) {
	p.idx += dx
	if p.idx >= p.path.Len() || p.idx < 0 {
		p.idx = 0
	}
	p.easer = object.EaserTo(p.io)(p.path[p.idx])
}

type PlayState struct {
	player *player
	title  *object.Pinner
	sheet  *tile.TileSheet
	world  *ebiten.Image
	idx    int
}

func (ps *PlayState) Update(e *engine.Engine) error {
	ps.player.easer()
	return nil
}

func (ps *PlayState) Draw(e *engine.Engine, s *ebiten.Image) {
	s.DrawImage(ps.world, nil)
	s.DrawImage(ps.player.io.Image())
	s.DrawImage(ps.title.Image(e.Camera.Pos()))
	ebitenutil.DebugPrint(s, fmt.Sprintf("index: %d", ps.idx))
}

func (ps *PlayState) Init(e *engine.Engine) {
	w, _ := e.Size()

	ps.sheet = tile.NewTileSheetFromTSX(e.Asset("assets/tiles/tile_sheet.png"), e.Asset("assets/tiles/tile_set.tsx"))
	worldImg, worldObjects := tile.NewTileMapFromTMX(ps.sheet, e.Asset("assets/tiles/tile_map.tmx"))
	ps.world = worldImg

	p1 := new(player)
	worldObjects["blue_spaces"].Sort()
	p1.path = worldObjects["blue_spaces"]
	player := ebiten.NewImage(10, 10)
	player.Fill(color.White)
	p1.io = object.FromEbitenImage(player)
	p1.io.SetPosition(worldObjects["blue_spaces"][0].GetPosition())
	p1.easer = func() {}

	ps.player = p1

	title := object.FromText(e.Font(), "PLAYER 1 | ¥: 0 / $: 0", color.White)
	tx, ty := object.CenterH(w, title)
	title.SetPosition(tx, ty+40)
	ps.title = object.NewPinner(title)
}

func (ps *PlayState) Register(e *engine.Engine) {
	e.RegisterKey(ebiten.KeyBackspace, func(e *engine.Engine) {
		e.PopState()
	})
	e.RegisterHeldKey(ebiten.KeyLeft, func(e *engine.Engine) {
		x, y := ps.player.io.GetPosition()
		ps.player.io.SetPosition(x-1, y)
	})
	e.RegisterHeldKey(ebiten.KeyRight, func(e *engine.Engine) {
		x, y := ps.player.io.GetPosition()
		ps.player.io.SetPosition(x+1, y)
	})
	e.RegisterHeldKey(ebiten.KeyUp, func(e *engine.Engine) {
		x, y := ps.player.io.GetPosition()
		ps.player.io.SetPosition(x, y-1)
	})
	e.RegisterHeldKey(ebiten.KeyDown, func(e *engine.Engine) {
		x, y := ps.player.io.GetPosition()
		ps.player.io.SetPosition(x, y+1)
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
	e.RegisterKey(ebiten.KeySpace, func(e *engine.Engine) {
		ps.player.move(1)
	})
}
