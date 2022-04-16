package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine/object"
)

type Camera struct {
	object.Object
}

func NewCamera() *Camera {
	c := new(Camera)
	c.Object = object.NewEmpty(0, 0)
	return c
}

func (c *Camera) View() *ebiten.DrawImageOptions {
	o := &ebiten.DrawImageOptions{}
	x, y := c.GetPosition()
	o.GeoM.Translate(-x, -y)
	return o
}
