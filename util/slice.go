package util

func Get[T any](slice []T, idx int) (T, bool) {
	var zero T

	if idx < 0 || idx >= len(slice) {
		return zero, false
	}
	return slice[idx], true
}

func GetOrDefault[T any](slice []T, idx int, def T) T {
	if idx < 0 || idx >= len(slice) {
		return def
	}
	return slice[idx]
}

func Set[T any](slice []T, idx int, value T) bool {
	if idx < 0 || idx >= len(slice) {
		return false
	}
	slice[idx] = value
	return true
}
