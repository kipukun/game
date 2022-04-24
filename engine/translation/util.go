package translation

func lerp(v0, v1, t float64) float64 {
	return v0 + t*(v1-v0)
}
