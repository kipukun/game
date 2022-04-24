package translation

func lerp(v0, v1, t float64) float64 {
	return v0 + t*(v1-v0)
}

func clamp(x, min, max float64) float64 {
	if x < min {
		x = min
	} else if x > max {
		x = max
	}
	return x
}
