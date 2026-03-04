package condition

func If[T any](cond bool, then T, _else T) T {
	if cond {
		return then
	}

	return _else
}
