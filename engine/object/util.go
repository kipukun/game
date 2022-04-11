package object

// CenterV centers o on s vertically and returns its new position.
func CenterV(h float64, o Object) (float64, float64) {
	_, Dy := o.Size()
	x, _ := o.GetPosition()
	o.SetPosition(x, float64(h)/2-float64(Dy)/2)
	return o.GetPosition()
}

// CenterH centers o on s horizontally.
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
