package common

// Ptr returns a pointer to the provided value.
func Ptr[T any](value T) *T {
	return &value
}

// ValueOrDefault returns value when it is not the zero value, otherwise fallback.
func ValueOrDefault[T comparable](value T, fallback T) T {
	var zero T
	if value == zero {
		return fallback
	}
	return value
}

// Clamp bounds a numeric value to an inclusive range.
func Clamp[T ~int | ~int32 | ~int64 | ~float32 | ~float64](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
