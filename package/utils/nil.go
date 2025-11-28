package utils

func OrDefault[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}

func Ptr[T any](value T) *T {
	return &value
}

func IsNilOrBlank(s *string) bool {
	return s == nil || *s == ""
}
