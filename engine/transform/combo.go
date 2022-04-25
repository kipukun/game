package transform

func Chain(fs ...ChangeFunc) ChangeFunc {
	i := 0
	return func(dt float64) float64 {
		if i > len(fs)-1 {
			return 1
		}
		if fs[i](dt) >= 1.0 {
			i++
			return 1
		}
		return fs[i](dt)
	}
}

func Combine(fs ...ChangeFunc) ChangeFunc {
	progress := 0.0
	return func(dt float64) float64 {
		for _, f := range fs {
			progress = f(dt)
		}
		return progress
	}
}
