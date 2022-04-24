package object

import "math"

// CenterV centers o on s vertically and returns its new position.
func CenterV(h float64, o Object) (float64, float64) {
	_, Dy := o.Size()
	x, _ := o.GetPosition()
	o.SetPosition(x, float64(h)/2-float64(Dy)/2)
	return o.GetPosition()
}

// CenterH centers o on s horizontally and returns its new position.
func CenterH(w float64, o Object) (float64, float64) {
	Dx, _ := o.Size()
	_, y := o.GetPosition()
	o.SetPosition(w/2-float64(Dx)/2, y)
	return o.GetPosition()
}

// Middle centers o in the middle of s. Equivalent to calling both CenterV and CenterH.
func Middle(w, h float64, o Object) (float64, float64) {
	CenterH(w, o)
	CenterV(h, o)
	return o.GetPosition()
}

// Offset offsets o from from by n pixels in d direction.
func Offset(o, from Object, n float64) (float64, float64) {
	fromx, _ := from.GetPosition()
	_, y := o.GetPosition()
	o.SetPosition(fromx-n, y)
	return fromx - n, y
}

// Distance returns the Euclidean distance between p and q.
func Distance(p, q Object) float64 {
	p1, p2 := p.GetPosition()
	q1, q2 := q.GetPosition()
	return math.Sqrt(math.Pow(q1-p1, 2) + math.Pow(q2-p2, 2))
}

func L1Distance(p, q Object) float64 {
	p1, p2 := p.GetPosition()
	q1, q2 := q.GetPosition()
	return math.Abs(q1-p1) + math.Abs(q2-p2)
}

func Translate(o Object, dx, dy float64) {
	x, y := o.GetPosition()
	o.SetPosition(x+dx, y+dy)
}
