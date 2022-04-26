package object

import "errors"

var (
	NegativeSizeError = errors.New("object: negative size passed")
)

// Object represents some in-game object, with some
// bounding box and an ability to update its position.
type Object interface {
	Update()
	Size() (width, height int)

	GetPosition() (x, y float64)
	GetVelocity() (dx, dy float64)
	GetAcceleration() (ddx, ddy float64)

	SetPosition(x, y float64)
	SetVelocity(dx, dy float64)
	SetAcceleration(ddx, ddy float64)
}

type obj struct {
	w, h     int
	x, y     float64
	dx, dy   float64
	ddx, ddy float64
}

// NewEmpty returns an empty Object with its width and height set to w and h, respectively.
func NewEmpty(w, h int) (Object, error) {
	o := new(obj)
	if w < 0 || h < 0 {
		return nil, NegativeSizeError
	}
	o.w = w
	o.h = h
	return o, nil
}

func (o *obj) Update() {
	o.dx += o.ddx
	o.dy += o.ddy
	o.x += o.dx
	o.y += o.dy
}

func (o *obj) Size() (int, int) {
	return o.w, o.h
}

func (o *obj) Pos() (float64, float64) {
	return o.x, o.y
}

func (o *obj) GetPosition() (x, y float64) {
	return o.x, o.y
}

func (o *obj) GetVelocity() (dx, dy float64) {
	return o.dx, o.dy
}
func (o *obj) GetAcceleration() (ddx, ddy float64) {
	return o.dx, o.dy
}

func (o *obj) SetPosition(x, y float64) {
	o.x = x
	o.y = y
}

func (o *obj) SetVelocity(dx, dy float64) {
	o.dx = dx
	o.dy = dy
}

func (o *obj) SetAcceleration(ddx, ddy float64) {
	o.ddx = ddx
	o.ddy = ddy
}
