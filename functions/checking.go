package functions

func IsAllAboveZero(values ...int) bool {
	for _, v := range values {
		if v <= 0 {
			return false
		}
	}
	return true
}

func IsAllNonEmpty(values ...string) bool {
	for _, v := range values {
		if v == "" {
			return false
		}
	}
	return true
}
