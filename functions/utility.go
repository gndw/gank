package functions

func CombineStringArray(arrays ...[]string) []string {
	result := []string{}
	for _, a := range arrays {
		result = append(result, a...)
	}
	return result
}
