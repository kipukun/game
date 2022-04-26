package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine/object"
)

// Camera is a special Object that determines what is shown on screen.
type Camera struct {
	object.Object
}

func NewCamera() *Camera {
	c := new(Camera)
	c.Object, _ = object.NewEmpty(0, 0)
	return c
}

// View returns what the Camera can currently see.
func (c *Camera) View() *ebiten.DrawImageOptions {
	o := &ebiten.DrawImageOptions{}
	x, y := c.GetPosition()
	o.GeoM.Translate(-x, -y)
	return o
}

func (c *Camera) Reset() {
	c.SetPosition(0, 0)
}

func (c *Camera) Pos() (float64, float64) {
	x, y := c.GetPosition()
	return float64(x), float64(y)
}
