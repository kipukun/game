package tile

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine/errors"
)

type TileSheet struct {
	sheet  *ebiten.Image
	dx, dy int
}

func NewTileSheet(s *ebiten.Image, dx, dy int) *TileSheet {
	ts := new(TileSheet)
	ts.sheet = s
	ts.dx, ts.dy = dx, dy
	return ts
}

func (ts *TileSheet) TileSize() (int, int) {
	return ts.dx, ts.dy
}

func (ts *TileSheet) Tile(x, y int) (*ebiten.Image, *ebiten.DrawImageOptions) {
	x = x - 1
	img := ts.sheet.SubImage(image.Rect(ts.dx*x, ts.dy*y, ts.dx*x+ts.dx, ts.dy*y+ts.dy))
	return ebiten.NewImageFromImage(img), &ebiten.DrawImageOptions{}
}

func NewTileMap(s *TileSheet, layers [][]image.Point, rowsize int) (*ebiten.Image, error) {
	var op errors.Op = "NewTileMap"

	if len(layers) < 1 {
		return nil, errors.Error(op, "need at least one layer in call to NewTileMap", nil)
	}
	tx, ty := s.TileSize()
	img := ebiten.NewImage(rowsize*tx, len(layers[0])/rowsize*ty)
	last := len(layers[0])
	for _, layer := range layers {
		if len(layer) != last {
			return nil, errors.Error(op, "all layers must be the same size", nil)
		}
		row := 0
		for i, pt := range layer {
			t, o := s.Tile(pt.X, pt.Y)
			o.GeoM.Translate(float64((i)%rowsize*tx), float64(row*ty))
			img.DrawImage(t, o)
			if (i+1)%rowsize == 0 {
				row += 1
			}
		}
		last = len(layer)
	}
	return img, nil
}
