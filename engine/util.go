package engine

func head[T any](array []T) T {
	return array[len(array)-1]
}
